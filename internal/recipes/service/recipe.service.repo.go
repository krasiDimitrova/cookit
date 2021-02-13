package service

import (
	"database/sql"
	"errors"
	"github.com/krasimiraMilkova/cookit/internal/db"
	"github.com/krasimiraMilkova/cookit/pkg/recipes"
	"strings"
)

type RecipeRepository struct {
	*sql.DB
}

func GetRecipeRepository() recipes.RecipeRepository {
	return &RecipeRepository{db.Get()}
}

func (recipeRepository *RecipeRepository) CreateRecipe(recipe *recipes.Recipe) error {
	if recipe.Title == "" || recipe.Directions == "" || len(recipe.Ingredients) == 0 {
		return errors.New("recipe cannot have empty fields")
	}

	result, err := recipeRepository.Exec("insert into recipes(title, directions)values(?,?);",
		recipe.Title, recipe.Directions)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	recipeId := int(id)

	for _, ingredient := range recipe.Ingredients {

		ingredientId, _ := recipeRepository.findIngredientIdByName(ingredient.Name)

		if ingredientId == 0 {
			result, _ := recipeRepository.Exec("insert ignore into ingredients(name)values(?);", ingredient.Name)
			lastInsertedId, _ := result.LastInsertId()
			ingredientId = int(lastInsertedId)
		}

		_, err = recipeRepository.Exec("insert into recipe_ingredients(recipe_id, ingredient_id, quantity, measurement)values(?,?,?,?);",
			recipeId, ingredientId, ingredient.Quantity, ingredient.Measurement)

		if err != nil {
			recipeRepository.deleteRecipeById(recipeId)
			return err
		}
	}

	return nil
}

func (recipeRepository *RecipeRepository) deleteRecipeById(id int) {
	recipeRepository.Exec("delete from recipes where id = ?;", id)
}

func (recipeRepository *RecipeRepository) findIngredientIdByName(name string) (int, error) {
	var id int
	row := recipeRepository.QueryRow("select id from ingredients where name = ?;", name)
	err := row.Scan(&id)

	if err == sql.ErrNoRows {
		return 0, err
	}

	return id, nil
}

func (recipeRepository *RecipeRepository) FindRecipesByTitle(title string) ([]recipes.RecipeSearchResult, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	titleSearch := "%" + title + "%"
	return recipeRepository.findRecipes("select id, title from recipes where LOWER(title) like LOWER(?);", titleSearch)
}

func (recipeRepository *RecipeRepository) FindRecipesByIngredients(ingredients []string) ([]recipes.RecipeSearchResult, error) {
	if len(ingredients) == 0 {
		return nil, errors.New("ingredients list cannot be empty")
	}

	args := make([]interface{}, len(ingredients))
	for i, ingredient := range ingredients {
		args[i] = ingredient
	}
	argsWildCards := strings.Repeat(",?", len(args)-1)
	query := "select id, title from recipes where id in (select distinct recipe_id from recipe_ingredients where " +
		"ingredient_id in (select id from ingredients where name in (?" + argsWildCards + ")));"
	return recipeRepository.findRecipes(query, args...)
}

func (recipeRepository *RecipeRepository) findRecipes(query string, args ...interface{}) ([]recipes.RecipeSearchResult, error) {
	var results []recipes.RecipeSearchResult

	rows, err := recipeRepository.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		result := recipes.RecipeSearchResult{}
		err = rows.Scan(&result.ID, &result.Title)

		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (recipeRepository *RecipeRepository) FindRecipeById(id int) (*recipes.Recipe, error) {
	recipe := &recipes.Recipe{}
	recipeRow := recipeRepository.QueryRow("select * from recipes where id = ?;", id)
	err := recipeRow.Scan(&recipe.ID, &recipe.Title, &recipe.Directions)

	if err == sql.ErrNoRows {
		return nil, err
	}

	ingredientRows, err := recipeRepository.Query("select ing.id, ing.name, ri.quantity, ri.measurement "+
		"from recipe_ingredients as ri "+
		"join ingredients as ing on ri.ingredient_id = ing.id "+
		"where ri.recipe_id = ?", id)
	if err != nil {
		return nil, err
	}

	defer ingredientRows.Close()
	for ingredientRows.Next() {
		ingredient := recipes.Ingredient{}
		err = ingredientRows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Quantity, &ingredient.Measurement)

		if err != nil {
			return nil, err
		}

		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}

	return recipe, nil
}
