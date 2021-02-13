package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type RecipeApi struct {
	*http.Client
	cookies []*http.Cookie
}

func GetRecipeApi(cookies []*http.Cookie) *RecipeApi {
	return &RecipeApi{
		Client:  &http.Client{},
		cookies: cookies,
	}
}

type Recipe struct {
	Title       string       `json:"title"`
	Ingredients []Ingredient `json:"ingredients"`
	Directions  string       `json:"directions"`
}

type Ingredient struct {
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
	Measurement string `json:"measurement"`
}

type RecipeSearchResult struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// CreateRecipe function sends recipe creation request to the server
// Returns error if such occurs
func (ra *RecipeApi) CreateRecipe(recipe Recipe) error {
	payloadBuffer := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuffer).Encode(recipe)

	if err != nil {
		return err
	}

	request, _ := http.NewRequest("POST", serverUrl+"/api/v1/recipe", payloadBuffer)
	for _, cookie := range ra.cookies {
		request.AddCookie(cookie)
	}
	response, err := ra.Client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if !(response.StatusCode >= 200 && response.StatusCode < 300) {
		return errors.New("failed to create recipe")
	}

	return nil
}

// FindByTitle function sends search by title request to the server
// Returns error if such occurs or the obtained search results
func (ra *RecipeApi) FindByTitle(title string) ([]RecipeSearchResult, error) {
	return ra.findRecipes(serverUrl + "/api/v1/recipe?title=" + title)
}

// FindByIngredients function sends search by ingredients request to the server
// Returns error if such occurs or the obtained search results
func (ra *RecipeApi) FindByIngredients(ingredients string) ([]RecipeSearchResult, error) {
	return ra.findRecipes(serverUrl + "/api/v1/recipe?ingredients=" + ingredients)
}

func (ra *RecipeApi) findRecipes(url string) ([]RecipeSearchResult, error) {
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	for _, cookie := range ra.cookies {
		request.AddCookie(cookie)
	}
	response, err := ra.Client.Do(request)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	defer response.Body.Close()

	if !(response.StatusCode >= 200 && response.StatusCode < 300) {
		return nil, errors.New("failed to find recipes")
	}

	var recipes []RecipeSearchResult
	err = json.NewDecoder(response.Body).Decode(&recipes)
	if err != nil {
		fmt.Print(err)
		return nil, errors.New("failed to decode search result")
	}

	return recipes, nil
}

// GetById function sends a get recipe request to the server
// Returns error if such occurs or the obtained recipe
func (ra *RecipeApi) GetById(id int) (*Recipe, error) {
	url := serverUrl + "/api/v1/recipe/" + strconv.Itoa(id)
	request, _ := http.NewRequest("GET", url, nil)
	for _, cookie := range ra.cookies {
		request.AddCookie(cookie)
	}
	response, err := ra.Client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if !(response.StatusCode >= 200 && response.StatusCode < 300) {
		return nil, errors.New("failed to fetch recipe")
	}

	var recipe *Recipe
	err = json.NewDecoder(response.Body).Decode(&recipe)

	if err != nil {
		return nil, errors.New("failed to decode recipe")
	}

	return recipe, nil
}

// GetCommentsById function sends requests for retrieving recipe comments to the server
// Returns error if such occurs or the obtained comments
func (ra *RecipeApi) GetCommentsById(id int) ([]string, error) {
	url := serverUrl + "/api/v1/recipe/" + strconv.Itoa(id) + "/comment"
	request, _ := http.NewRequest("GET", url, nil)
	for _, cookie := range ra.cookies {
		request.AddCookie(cookie)
	}
	response, err := ra.Client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if !(response.StatusCode >= 200 && response.StatusCode < 300) {
		if response.StatusCode == 404 {
			return nil, nil
		}
		return nil, errors.New("failed to fetch comments")
	}

	var comments []string
	err = json.NewDecoder(response.Body).Decode(&comments)

	if err != nil {
		return nil, errors.New("failed to decode comments")
	}

	return comments, nil
}

// AddCommentById function sends post request for creating a recipe comments to the server
// Returns error if such occurs
func (ra *RecipeApi) AddCommentById(recipeId int, comment string) error {
	payloadBuffer := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuffer).Encode(comment)

	if err != nil {
		return err
	}

	url := serverUrl + "/api/v1/recipe/" + strconv.Itoa(recipeId) + "/comment"
	request, _ := http.NewRequest("POST", url, payloadBuffer)
	for _, cookie := range ra.cookies {
		request.AddCookie(cookie)
	}
	response, err := ra.Client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if !(response.StatusCode >= 200 && response.StatusCode < 300) {
		return errors.New("failed to add comment")
	}

	return nil
}
