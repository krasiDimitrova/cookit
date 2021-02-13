package service

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/krasimiraMilkova/cookit/mocks"
	"github.com/krasimiraMilkova/cookit/pkg/recipes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestRecipeService_CreateRecipe(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mocks.NewMockRecipeRepository(mockCtrl)

	service := RecipeService{RecipeRepository: mockRepository}

	tests := []struct {
		name               string
		recipe             recipes.Recipe
		repositoryError    string
		expectedStatusCode int
	}{
		{
			name: "Successful",
			recipe: recipes.Recipe{
				Title: "Test title",
				Ingredients:
				[]recipes.Ingredient{
					{
						Name:        "Test",
						Quantity:    1,
						Measurement: "m",
					},
					{
						Name:        "Smth",
						Quantity:    1,
						Measurement: "m",
					},
				},
				Directions: "Test directions",
			},
			repositoryError:    "",
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Repository error",
			recipe:             recipes.Recipe{},
			repositoryError:    "recipe cannot have empty fields",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonRecipe, _ := json.Marshal(test.recipe)
			req, _ := http.NewRequest("POST", "/recipe", strings.NewReader(string(jsonRecipe)))
			rr := httptest.NewRecorder()

			var err error
			if test.repositoryError != "" {
				err = errors.New(test.repositoryError)
			}
			mockRepository.EXPECT().CreateRecipe(&test.recipe).Return(err)
			http.HandlerFunc(service.CreateRecipe).ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, test.expectedStatusCode)
				t.Fail()
			}
		})
	}
}

func TestRecipeService_FindRecipesByTitle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mocks.NewMockRecipeRepository(mockCtrl)

	service := RecipeService{RecipeRepository: mockRepository}

	tests := []struct {
		name               string
		title              string
		searchResults      []recipes.RecipeSearchResult
		repositoryError    string
		expectedStatusCode int
	}{
		{
			name:  "Successful",
			title: "test",
			searchResults:
			[]recipes.RecipeSearchResult{
				{
					ID:    0,
					Title: "Test recipe",
				},
				{
					ID:    5,
					Title: "Smth test",
				},
			},
			repositoryError:    "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Empty search parameter - title",
			title:              "",
			searchResults:      []recipes.RecipeSearchResult{},
			repositoryError:    "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Repository returns error",
			title:              "Some title",
			searchResults:      []recipes.RecipeSearchResult{},
			repositoryError:    "some repo error",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "Repository returns no results",
			title:              "Some title",
			searchResults:      []recipes.RecipeSearchResult{},
			repositoryError:    "",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/recipe?title="+test.title, nil)
			rr := httptest.NewRecorder()

			var err error
			if test.repositoryError != "" {
				err = errors.New(test.repositoryError)
			}
			if test.expectedStatusCode != http.StatusBadRequest {
				mockRepository.EXPECT().FindRecipesByTitle(test.title).Return(test.searchResults, err)
			}

			http.HandlerFunc(service.FindRecipesByTitle).ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, test.expectedStatusCode)
				t.Fail()
			}

			var searchResults []recipes.RecipeSearchResult
			s := rr.Body.String()
			if s != "" {
				err = json.Unmarshal([]byte(s), &searchResults)

				if err != nil {
					t.Error("error from unmarshal", err)
				}
			}

			if len(searchResults) > 0 && !reflect.DeepEqual(searchResults, test.searchResults) {
				t.Errorf("Got searchResults = %v but wanted %v", searchResults, test.searchResults)
			} else if len(test.searchResults) == 0 && len(searchResults) != 0 {
				t.Errorf("Expected [] result but got %v", searchResults)
			}
		})
	}
}

func TestRecipeService_FindRecipesByIngredients(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mocks.NewMockRecipeRepository(mockCtrl)

	service := RecipeService{RecipeRepository: mockRepository}

	tests := []struct {
		name               string
		ingredients        string
		searchResults      []recipes.RecipeSearchResult
		repositoryError    string
		expectedStatusCode int
	}{
		{
			name:        "Successful",
			ingredients: "test,smth",
			searchResults:
			[]recipes.RecipeSearchResult{
				{
					ID:    0,
					Title: "Test recipe",
				},
			},
			repositoryError:    "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Empty search parameter - ingredients",
			ingredients:        "",
			searchResults:      []recipes.RecipeSearchResult{},
			repositoryError:    "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Repository returns error",
			ingredients:        "test,smth",
			searchResults:      []recipes.RecipeSearchResult{},
			repositoryError:    "some repo error",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "Repository returns no results",
			ingredients:        "test",
			searchResults:      []recipes.RecipeSearchResult{},
			repositoryError:    "",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/recipe?ingredients="+test.ingredients, nil)
			rr := httptest.NewRecorder()

			var err error
			if test.repositoryError != "" {
				err = errors.New(test.repositoryError)
			}
			if test.expectedStatusCode != http.StatusBadRequest {
				ingredientsList := strings.Split(test.ingredients, ",")
				mockRepository.EXPECT().FindRecipesByIngredients(ingredientsList).Return(test.searchResults, err)
			}

			http.HandlerFunc(service.FindRecipesByIngredients).ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, test.expectedStatusCode)
				t.Fail()
			}

			var searchResults []recipes.RecipeSearchResult
			s := rr.Body.String()
			if s != "" {
				err = json.Unmarshal([]byte(s), &searchResults)

				if err != nil {
					t.Error("error from unmarshal", err)
				}
			}

			if len(searchResults) > 0 && !reflect.DeepEqual(searchResults, test.searchResults) {
				t.Errorf("Got searchResults = %v but wanted %v", searchResults, test.searchResults)
			} else if len(test.searchResults) == 0 && len(searchResults) != 0 {
				t.Errorf("Expected [] result but got %v", searchResults)
			}
		})
	}
}

func TestRecipeService_FindRecipeById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mocks.NewMockRecipeRepository(mockCtrl)

	service := RecipeService{RecipeRepository: mockRepository}

	tests := []struct {
		name               string
		id                 string
		recipe             recipes.Recipe
		repositoryError    string
		expectedStatusCode int
	}{
		{
			name: "Successful",
			id:   "1",
			recipe: recipes.Recipe{
				ID:    1,
				Title: "Test title",
				Ingredients:
				[]recipes.Ingredient{
					{
						Name:        "Test",
						Quantity:    1,
						Measurement: "m",
					},
				},
				Directions: "Test directions",
			},
			repositoryError:    "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Cannot parse id to number",
			id:                 "a",
			recipe:             recipes.Recipe{},
			repositoryError:    "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Repository returns no results",
			id:                 "2",
			recipe:             recipes.Recipe{},
			repositoryError:    "",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/recipe", nil)
			req = mux.SetURLVars(req, map[string]string{
				"id": test.id,
			})
			rr := httptest.NewRecorder()

			var err error
			if test.repositoryError != "" {
				err = errors.New(test.repositoryError)
			}
			if test.expectedStatusCode != http.StatusBadRequest {
				id, _ := strconv.Atoi(test.id)
				if test.expectedStatusCode == http.StatusOK {
					mockRepository.EXPECT().FindRecipeById(id).Return(&test.recipe, err)
				} else {
					mockRepository.EXPECT().FindRecipeById(id).Return(nil, err)
				}
			}

			http.HandlerFunc(service.FindRecipeById).ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, test.expectedStatusCode)
				t.Fail()
			}

			var recipe recipes.Recipe
			s := rr.Body.String()
			if s != "" {
				err = json.Unmarshal([]byte(s), &recipe)

				if err != nil {
					t.Error("error from unmarshal", err)
				}
			}

			if !reflect.DeepEqual(recipe, test.recipe) {
				t.Errorf("Got recipe = %v but wanted %v", recipe, test.recipe)
			}
		})
	}
}
