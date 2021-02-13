package recipes

// Recipe struct describes a recipe for cooking consisting of title, ingredients and directions
type Recipe struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	Ingredients []Ingredient `json:"ingredients"`
	Directions  string       `json:"directions"`
}
