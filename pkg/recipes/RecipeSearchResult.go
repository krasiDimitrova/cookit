package recipes

// RecipeSearchResult serves as a result form the recipe search operations
type RecipeSearchResult struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}
