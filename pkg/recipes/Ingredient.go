package recipes

// Ingredient struct describes a recipe ingredient with name, quantity and the quantity measurement
type Ingredient struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Quantity    uint   `json:"quantity"`
	Measurement string `json:"measurement"`
}
