package api

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"abberant/backend/internal/articles"
	"abberant/backend/internal/auth"
	"abberant/backend/internal/ratings"
	"abberant/backend/internal/sections"
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
	// ReaderRating is the aggregated reader vote.
	ReaderRating *readerRatingResponse `json:"readerRating,omitempty"`
}

type readerRatingResponse struct {
	Average float64 `json:"average"`
	Votes   int     `json:"votes"`
	MyVote  int     `json:"myVote,omitempty"`
}

// withReaderRating attaches the reader vote summary to the article response.
func withReaderRating(
	response articleResponse,
	stats ratings.Stats,
	myVote int,
) articleResponse {
	response.ReaderRating = &readerRatingResponse{
		Average: stats.Average,
		Votes:   stats.Votes,
		MyVote:  myVote,
	}

	return response
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
	return user != nil && auth.CanModerate(user.Role)
}

// validateArticleSection enforces the kind↔section rules for published
// material: walkthroughs need a live hub, news an optional live platform
// section, reviews no section at all. Drafts are exempt so half-finished
// work can always be saved.
func (s *Server) validateArticleSection(ctx context.Context, article *articles.Article) error {
	if article.Kind == articles.KindReview {
		if article.Section != "" {
			return auth.Fail(400, "Для обзоров раздел не указывается.")
		}

		return nil
	}

	if article.Section == "" {
		if article.Kind == articles.KindWalkthrough {
			return auth.Fail(400, "Выберите игровой хаб для прохождения.")
		}

		return nil
	}

	section, err := s.sections.ByID(ctx, article.Section)
	if err != nil {
		return err
	}

	if section == nil {
		return auth.Fail(400, "Указанный раздел не найден.")
	}

	switch article.Kind {
	case articles.KindWalkthrough:
		if section.Type != sections.TypeWalkthrough {
			return auth.Fail(400, "Указанный раздел не является хабом прохождения.")
		}
	case articles.KindNews:
		if section.Type != sections.TypeNews {
			return auth.Fail(400, "Указанный раздел не является новостной секцией.")
		}
	}

	if section.Archived {
		return auth.Fail(400, "Раздел находится в архиве.")
	}

	return nil
}

// notifyAboutPublish creates "new material" notifications for section
// subscribers when a article appears published or a draft goes live.
func (s *Server) notifyAboutPublish(ctx *fiber.Ctx, article *articles.Article, wasPublished bool) error {
	if article.Status != articles.StatusPublished || wasPublished || article.Section == "" {
		return nil
	}

	user := s.currentUser(ctx)

	_, err := s.notifications.CreateArticleEvents(ctx.Context(), article.Section, article.ID, user.ID)

	return err
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

	// The stored counter is legacy data; the real number of comments wins.
	counts, err := s.comments.CountsByArticle(ctx.Context())
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

	var (
		voteStats map[string]ratings.Stats
		myVotes   map[string]int
	)

	if user != nil {
		myVotes, err = s.ratings.VotesOf(ctx.Context(), user.ID)
		if err != nil {
			return err
		}
	}

	voteStats, err = s.ratings.StatsByArticleMap(ctx.Context())
	if err != nil {
		return err
	}

	for _, article := range list {
		item := newArticleResponse(article)

		if count, ok := counts[article.ID]; ok {
			item.Comments = count
		}

		item = withReaderRating(item, voteStats[article.ID], myVotes[article.ID])

		response = append(response, item)
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

	if article.Status == articles.StatusPublished {
		if err := s.validateArticleSection(ctx.Context(), article); err != nil {
			return err
		}
	}

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

	if err := s.notifyAboutPublish(ctx, created, false); err != nil {
		return err
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
	Section  *string `json:"section"`
	Date     *string `json:"date"`
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

	wasPublished := existing.Status == articles.StatusPublished

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

	if req.Date != nil {
		parsed, ok := articles.ParseDate(*req.Date)
		if !ok {
			return auth.Fail(400, "Некорректная дата статьи.")
		}

		existing.Date = parsed
	}

	if req.Section != nil {
		section := strings.TrimSpace(*req.Section)
		if len(section) > 64 {
			return auth.Fail(400, "Слишком длинный идентификатор раздела.")
		}

		existing.Section = section
	}

	if existing.Status == articles.StatusPublished && strings.TrimSpace(existing.Content) == "" {
		return auth.Fail(400, "Добавьте текст статьи перед публикацией.")
	}

	if existing.Status == articles.StatusPublished {
		if err := s.validateArticleSection(ctx.Context(), existing); err != nil {
			return err
		}
	}

	updated, err := s.articles.Update(ctx.Context(), existing)
	if err != nil {
		return err
	}

	if err := s.notifyAboutPublish(ctx, updated, wasPublished); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"article": newArticleResponse(updated)})
}

func canDeleteArticle(user *auth.User, article *articles.Article) bool {
	if user == nil {
		return false
	}

	if user.Role == auth.RoleOwner || user.Role == auth.RoleModerator {
		return true
	}

	return article.Author == user.Username
}

func (s *Server) handleDeleteArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	existing, err := s.articles.ByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	if existing == nil {
		return auth.Fail(404, "Статья не найдена.")
	}

	_, user, err := s.resolve(ctx)
	if err != nil {
		return err
	}

	if !canDeleteArticle(user, existing) {
		return auth.Fail(403, "Недостаточно прав для удаления этой статьи.")
	}

	removed, err := s.articles.Delete(ctx.Context(), id)
	if err != nil {
		return err
	}

	if !removed {
		return auth.Fail(404, "Статья не найдена.")
	}

	return ctx.JSON(fiber.Map{"status": "ok"})
}
