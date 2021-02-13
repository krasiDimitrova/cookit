// Package recipes provides handlers, db interactions and models for creating and fetching recipes
package recipes

import "net/http"

// RecipeService interface provide handlers for creating and searching for recipes
type RecipeService interface {
	// CreateRecipe function handles payload for creating a recipe
	// Returns Status BadRequest if cannot decode the payload,
	// Status InternalServerError if error occurs during recipe creation and
	// Status Created if recipe is successfully inserted into the db
	CreateRecipe(w http.ResponseWriter, r *http.Request)

	// FindRecipesByTitle function handles requests for fetching recipes by title
	// provided as an query parameter
	// Returns Status BadRequest if cannot decode the query parameter,
	// Status InternalServerError if error occurs during searching and
	// Status NotFound if no recipes have been found for the title and
	// Status OK and RecipeSearchResults if such are found
	FindRecipesByTitle(w http.ResponseWriter, r *http.Request)

	// FindRecipesByIngredients function handles requests for fetching recipes by list of ingredients
	// provided as an query parameter
	// Returns Status BadRequest if cannot decode the query parameter,
	// Status InternalServerError if error occurs during searching and
	// Status NotFound if no recipes have been found for the ingredients and
	// Status OK and RecipeSearchResults if such are found
	FindRecipesByIngredients(w http.ResponseWriter, r *http.Request)

	// FindRecipesByTitle function handles requests for fetching recipes by id
	// provided as a path variable
	// Returns Status BadRequest if cannot parse the recipe id,
	// Status InternalServerError if error occurs during fetching and
	// Status NotFound if a recipes with this id does not exist and
	// Status OK and the Recipe if such are found
	FindRecipeById(w http.ResponseWriter, r *http.Request)
}
