package recipes

// RecipeRepository interface provides functions for CRUD operations for recipe entity
type RecipeRepository interface {
	// CreateRecipe function provide insert db operation for recipe and included ingredients that do not exist yet
	// Returns an error if such occurs during the db query execution
	CreateRecipe(recipe *Recipe) error

	// FindRecipesByTitle function provide search operation for recipes by given title
	// Returns an error if such occurs during the db query execution otherwise returns a List of RecipeSearchResults
	FindRecipesByTitle(title string) ([]RecipeSearchResult, error)

	// FindRecipesByIngredients function provide search operation for recipes by given list of ingredient names
	// Returns an error if such occurs during the db query execution otherwise returns a List of RecipeSearchResults
	FindRecipesByIngredients(ingredients []string) ([]RecipeSearchResult, error)

	// FindRecipeById( function provide operation for obtaining a recipe and its ingredients for the given id
	// Returns an error if such occurs during the db query execution otherwise returns a Recipe
	FindRecipeById(id int) (*Recipe, error)
}
