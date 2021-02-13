// Package comments provide handlers and db operations for creating and fetching recipe comments
package comments

import "net/http"

// CommentService interface provide handlers for creating and fetching recipe comments
type CommentService interface {
	// AddComment function handles comment payload for recipeId provided as a path variable
	// Returns Status BadRequest if cannot parse the recipeId,
	// Status InternalServerError if error occurs during comment insertion and
	// Status NotFound if recipe with given id does not exist and
	// Status Created if comment is successfully inserted into the db
	AddComment(w http.ResponseWriter, r *http.Request)

	// GetComments function handles requests for fetching comments by provided id as a path variable
	// Returns Status BadRequest if cannot parse the recipe id,
	// Status InternalServerError if error occurs during fetching and
	// Status NotFound if a recipe with this id or comments for it does not exist and
	// Status OK and the Comments if such are found
	GetComments(w http.ResponseWriter, r *http.Request)
}
