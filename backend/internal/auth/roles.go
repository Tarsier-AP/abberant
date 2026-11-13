package auth

// Role names stored in the free-form users.role column. New accounts default
// to RoleReader; the migration seeds also create a RoleUser test account.
const (
	RoleOwner     = "owner"
	RoleModerator = "moderator"
	RoleEditor    = "editor"
	RoleCreator   = "creator"
	RoleReader    = "reader"
	RoleUser      = "user"
)

// CanModerate reports whether the role may moderate content: see drafts of
// other authors and delete other people's comments.
func CanModerate(role string) bool {
	switch role {
	case RoleOwner, RoleModerator, RoleEditor:
		return true
	default:
		return false
	}
}

// CanEditArticles reports whether the role may create, edit and delete
// articles.
func CanEditArticles(role string) bool {
	switch role {
	case RoleOwner, RoleEditor, RoleCreator:
		return true
	default:
		return false
	}
}
