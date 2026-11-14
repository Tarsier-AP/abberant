package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"abberant/backend/internal/auth"
	"abberant/backend/internal/sections"
)

type sectionResponse struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Archived bool   `json:"archived"`
}

func newSectionResponse(section *sections.Section) sectionResponse {
	return sectionResponse{
		ID:       section.ID,
		Type:     section.Type,
		Title:    section.Title,
		Archived: section.Archived,
	}
}

// handleListSections is public: the frontend loads the full section registry
// (archived ones included, flagged) once at startup.
func (s *Server) handleListSections(ctx *fiber.Ctx) error {
	list, err := s.sections.List(ctx.Context())
	if err != nil {
		return err
	}

	response := make([]sectionResponse, 0, len(list))

	for _, section := range list {
		response = append(response, newSectionResponse(section))
	}

	return ctx.JSON(fiber.Map{"sections": response})
}

type createSectionRequest struct {
	Type  string `json:"type"`
	Title string `json:"title"`
}

// handleCreateSection adds a section under a server-generated slug derived
// from the title.
func (s *Server) handleCreateSection(ctx *fiber.Ctx) error {
	var req createSectionRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	sectionType := strings.TrimSpace(req.Type)
	if !sections.IsValidType(sectionType) {
		return auth.Fail(400, "Некорректный тип секции.")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return auth.Fail(400, "Введите название секции.")
	}

	if len([]rune(title)) > 120 {
		return auth.Fail(400, "Слишком длинное название секции.")
	}

	section, err := s.sections.Create(ctx.Context(), sectionType, title)
	if err != nil {
		return err
	}

	return ctx.Status(201).JSON(fiber.Map{"section": newSectionResponse(section)})
}

type updateSectionRequest struct {
	Title    *string `json:"title"`
	Archived *bool   `json:"archived"`
}

// handleUpdateSection renames a section (the slug stays immutable) and
// archives/restores it. Archiving never touches the assigned articles.
func (s *Server) handleUpdateSection(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	existing, err := s.sections.ByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	if existing == nil {
		return auth.Fail(404, "Секция не найдена.")
	}

	var req updateSectionRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return auth.Fail(400, "Введите название секции.")
		}

		if len([]rune(title)) > 120 {
			return auth.Fail(400, "Слишком длинное название секции.")
		}

		existing.Title = title
	}

	if req.Archived != nil {
		existing.Archived = *req.Archived
	}

	updated, err := s.sections.Update(ctx.Context(), existing)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"section": newSectionResponse(updated)})
}
