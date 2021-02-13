package comments

// CommentRepository interface provides functions for CRUD operations for recipe comments
type CommentRepository interface {
	// AddComment provides an insert operation for the given recipe id and comment
	// Returns an error if such occurs during db query execution
	AddComment(recipeId int, comment string) error

	// GetComments provides a fetch operation for recipe comments for provided recipeId
	// Returns an error such occurs during db query execution otherwise returns the found comments
	GetComments(recipeId int) ([]string, error)
}
