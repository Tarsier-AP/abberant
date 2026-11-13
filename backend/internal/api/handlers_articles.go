package api

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"abberant/backend/internal/articles"
	"abberant/backend/internal/auth"
)

type articleResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Image    string `json:"image"`
	Date     string `json:"date"`
	Comments int    `json:"comments"`
	Kind     string `json:"kind"`
	Content  string `json:"content"`
	Author   string `json:"author"`
	Status   string `json:"status"`
	Rating   int    `json:"rating,omitempty"`
	Platform string `json:"platform,omitempty"`
	Verdict  string `json:"verdict,omitempty"`
	Section  string `json:"section,omitempty"`
}

func newArticleResponse(article *articles.Article) articleResponse {
	return articleResponse{
		ID:       article.ID,
		Title:    article.Title,
		Image:    article.Image,
		Date:     article.Date.Format("2006-01-02"),
		Comments: article.Comments,
		Kind:     article.Kind,
		Content:  article.Content,
		Author:   article.Author,
		Status:   article.Status,
		Rating:   article.Rating,
		Platform: article.Platform,
		Verdict:  article.Verdict,
		Section:  article.Section,
	}
}

func canModerateArticles(user *auth.User) bool {
	return user != nil && (user.Role == "editor" || user.Role == "owner")
}

func (s *Server) handleListArticles(ctx *fiber.Ctx) error {
	_, user, err := s.resolve(ctx)
	if err != nil {
		return err
	}

	list, err := s.articles.All(ctx.Context())
	if err != nil {
		return err
	}

	if !canModerateArticles(user) {
		visible := make([]*articles.Article, 0, len(list))

		for _, article := range list {
			published := article.Status == articles.StatusPublished
			ownDraft := user != nil && article.Author == user.Username

			if published || ownDraft {
				visible = append(visible, article)
			}
		}

		list = visible
	}

	response := make([]articleResponse, 0, len(list))

	for _, article := range list {
		response = append(response, newArticleResponse(article))
	}

	return ctx.JSON(fiber.Map{"articles": response})
}

type createArticleRequest struct {
	ID       *string `json:"id"`
	Title    string  `json:"title"`
	Image    string  `json:"image"`
	Date     string  `json:"date"`
	Comments *int    `json:"comments"`
	Kind     string  `json:"kind"`
	Content  string  `json:"content"`
	Status   string  `json:"status"`
	Rating   *int    `json:"rating"`
	Platform string  `json:"platform"`
	Verdict  string  `json:"verdict"`
	Section  string  `json:"section"`
}

func (s *Server) handleCreateArticle(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	var req createArticleRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	article, err := buildArticleFromRequest(&req)
	if err != nil {
		return err
	}

	article.Author = user.Username

	requestedID := ""
	if req.ID != nil {
		requestedID = strings.TrimSpace(*req.ID)
	}

	if requestedID != "" {
		if len(requestedID) > 64 {
			return auth.Fail(400, "Слишком длинный идентификатор статьи.")
		}

		article.ID = requestedID
	} else {
		article.ID, err = s.articles.NextUserID(ctx.Context())
		if err != nil {
			return err
		}
	}

	created, err := s.articles.Create(ctx.Context(), article)
	if err != nil {
		// The requested id is already taken: fall back to a generated one.
		if !articles.IsUniqueViolation(err) || requestedID == "" {
			return err
		}

		article.ID, err = s.articles.NextUserID(ctx.Context())
		if err != nil {
			return err
		}

		created, err = s.articles.Create(ctx.Context(), article)
		if err != nil {
			return err
		}
	}

	return ctx.Status(201).JSON(fiber.Map{"article": newArticleResponse(created)})
}

func buildArticleFromRequest(req *createArticleRequest) (*articles.Article, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, auth.Fail(400, "Введите заголовок статьи.")
	}

	if len([]rune(title)) > 300 {
		return nil, auth.Fail(400, "Слишком длинный заголовок статьи.")
	}

	kind := req.Kind
	if kind == "" {
		kind = articles.KindNews
	}

	if !articles.ValidKind(kind) {
		return nil, auth.Fail(400, "Некорректный тип материала.")
	}

	status := req.Status
	if status == "" {
		status = articles.StatusDraft
	}

	if !articles.ValidStatus(status) {
		return nil, auth.Fail(400, "Некорректный статус статьи.")
	}

	image := strings.TrimSpace(req.Image)
	if image == "" {
		image = articles.FallbackCover
	}

	if !articles.ValidImage(image) {
		return nil, auth.Fail(400, "Некорректное изображение обложки.")
	}

	comments := 0
	if req.Comments != nil {
		if *req.Comments < 0 {
			return nil, auth.Fail(400, "Некорректное число комментариев.")
		}

		comments = *req.Comments
	}

	rating := 0
	if req.Rating != nil && *req.Rating != 0 && !articles.ValidRating(*req.Rating) {
		return nil, auth.Fail(400, "Оценка должна быть от 1 до 5.")
	}

	if req.Rating != nil {
		rating = *req.Rating
	}

	date := time.Now()
	if req.Date != "" {
		parsed, ok := articles.ParseDate(req.Date)
		if !ok {
			return nil, auth.Fail(400, "Некорректная дата статьи.")
		}

		date = parsed
	}

	section := strings.TrimSpace(req.Section)
	if len(section) > 64 {
		return nil, auth.Fail(400, "Слишком длинный идентификатор раздела.")
	}

	return &articles.Article{
		Title:    title,
		Image:    image,
		Date:     date,
		Comments: comments,
		Kind:     kind,
		Content:  req.Content,
		Status:   status,
		Rating:   rating,
		Platform: strings.TrimSpace(req.Platform),
		Verdict:  strings.TrimSpace(req.Verdict),
		Section:  section,
	}, nil
}

type updateArticleRequest struct {
	Title    *string `json:"title"`
	Image    *string `json:"image"`
	Kind     *string `json:"kind"`
	Content  *string `json:"content"`
	Status   *string `json:"status"`
	Rating   *int    `json:"rating"`
	Platform *string `json:"platform"`
	Verdict  *string `json:"verdict"`
}

func (s *Server) handleUpdateArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	existing, err := s.articles.ByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	if existing == nil {
		return auth.Fail(404, "Статья не найдена.")
	}

	var req updateArticleRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return auth.Fail(400, "Введите заголовок статьи.")
		}

		if len([]rune(title)) > 300 {
			return auth.Fail(400, "Слишком длинный заголовок статьи.")
		}

		existing.Title = title
	}

	if req.Image != nil {
		image := strings.TrimSpace(*req.Image)
		if image == "" {
			image = articles.FallbackCover
		}

		if !articles.ValidImage(image) {
			return auth.Fail(400, "Некорректное изображение обложки.")
		}

		existing.Image = image
	}

	if req.Kind != nil {
		if !articles.ValidKind(*req.Kind) {
			return auth.Fail(400, "Некорректный тип материала.")
		}

		existing.Kind = *req.Kind
	}

	if req.Content != nil {
		existing.Content = *req.Content
	}

	if req.Status != nil {
		if !articles.ValidStatus(*req.Status) {
			return auth.Fail(400, "Некорректный статус статьи.")
		}

		existing.Status = *req.Status
	}

	// Review-only fields always follow the payload: a missing value clears them
	// when the material is switched back to a news item.
	if req.Rating != nil && *req.Rating != 0 && !articles.ValidRating(*req.Rating) {
		return auth.Fail(400, "Оценка должна быть от 1 до 5.")
	}

	if req.Rating != nil {
		existing.Rating = *req.Rating
	} else {
		existing.Rating = 0
	}

	if req.Platform != nil {
		existing.Platform = strings.TrimSpace(*req.Platform)
	} else {
		existing.Platform = ""
	}

	if req.Verdict != nil {
		existing.Verdict = strings.TrimSpace(*req.Verdict)
	} else {
		existing.Verdict = ""
	}

	if existing.Status == articles.StatusPublished && strings.TrimSpace(existing.Content) == "" {
		return auth.Fail(400, "Добавьте текст статьи перед публикацией.")
	}

	updated, err := s.articles.Update(ctx.Context(), existing)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"article": newArticleResponse(updated)})
}

func (s *Server) handleDeleteArticle(ctx *fiber.Ctx) error {
	removed, err := s.articles.Delete(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}

	if !removed {
		return auth.Fail(404, "Статья не найдена.")
	}

	return ctx.JSON(fiber.Map{"status": "ok"})
}
