package api

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"

	"abberant/backend/internal/auth"
	"abberant/backend/internal/media"
)

type mediaUploadResponse struct {
	// URL is the site-relative address editors store in articles.image or
	// insert into markdown content.
	URL string `json:"url"`
	// Srcset lists every generated width variant as a srcset attribute value.
	Srcset string `json:"srcset,omitempty"`
	// Width and Height describe the derivative behind URL.
	Width  int `json:"width"`
	Height int `json:"height"`
}

// handleUploadMedia processes a multipart upload into media storage.
// Query/form kind: "cover"|"body" (editors only) or "avatar" (any session).
func (s *Server) handleUploadMedia(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	kind := media.Kind(strings.ToLower(strings.TrimSpace(ctx.FormValue("kind"))))
	if kind == "" {
		kind = media.KindCover
	}

	if kind != media.KindAvatar && !auth.CanEditArticles(user.Role) {
		return auth.Fail(403, "Недостаточно прав для управления статьями.")
	}

	if kind != media.KindCover && kind != media.KindBody && kind != media.KindAvatar {
		return auth.Fail(400, "Некорректный тип загрузки.")
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return auth.Fail(400, "Выберите файл изображения.")
	}

	if file.Size > media.MaxUpload {
		return auth.Fail(400, "Файл слишком большой. Максимальный размер — 10 МБ.")
	}

	opened, err := file.Open()
	if err != nil {
		return auth.Fail(400, "Не удалось прочитать файл.")
	}

	defer opened.Close()

	data, err := io.ReadAll(opened)
	if err != nil {
		return auth.Fail(400, "Не удалось прочитать файл.")
	}

	var result *media.Result

	if kind == media.KindAvatar {
		result, err = s.media.PutAvatar(data)
	} else {
		result, err = s.media.Put(data)
	}

	if err != nil {
		// Store errors are user-facing messages with an internal prefix.
		message := strings.TrimPrefix(err.Error(), "media: ")

		if len(message) > 0 {
			message = string(unicode.ToUpper([]rune(message)[0])) + message[1:]
		}

		return auth.Fail(400, message)
	}

	main := result.MainURL()

	width, height := 0, 0

	srcset := make([]string, 0, len(result.Files))

	for _, file := range result.Files {
		if file.Original {
			continue
		}

		srcset = append(srcset, file.URL+" "+strconv.Itoa(file.Width)+"w")

		if file.URL == main {
			width, height = file.Width, file.Height
		}
	}

	return ctx.Status(201).JSON(mediaUploadResponse{
		URL:    main,
		Srcset: strings.Join(srcset, ", "),
		Width:  width,
		Height: height,
	})
}

// handleMediaFile serves stored media with immutable caching: identical
// content always maps to the same URL, so responses never need revalidation.
func (s *Server) handleMediaFile(ctx *fiber.Ctx) error {
	path := ctx.Params("*")

	if !media.ValidMediaPath(path) {
		return auth.Fail(404, "Файл не найден.")
	}

	full := filepath.Join(s.media.Dir(), filepath.FromSlash(path))

	data, err := os.ReadFile(full)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return auth.Fail(404, "Файл не найден.")
		}

		return err
	}

	name := filepath.Base(path)

	switch {
	case strings.HasSuffix(name, ".webp"):
		ctx.Set("Content-Type", "image/webp")
	case strings.HasSuffix(name, ".png"):
		ctx.Set("Content-Type", "image/png")
	case strings.HasSuffix(name, ".jpg"):
		ctx.Set("Content-Type", "image/jpeg")
	}

	ctx.Set("Cache-Control", "public, max-age=31536000, immutable")

	return ctx.Send(data)
}
