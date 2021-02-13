package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/krasimiraMilkova/cookit/pkg/recipes"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type RecipeService struct {
	RecipeRepository recipes.RecipeRepository
}

var recipesService *RecipeService

func Get() *RecipeService {
	if recipesService == nil {
		recipesService = &RecipeService{RecipeRepository: GetRecipeRepository()}
	}

	return recipesService
}

func (rs *RecipeService) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	recipe := &recipes.Recipe{}
	err := json.NewDecoder(r.Body).Decode(recipe)

	if err != nil {
		log.Print("Error occurred when decoding recipe payload ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := rs.RecipeRepository.CreateRecipe(recipe); err != nil {
		log.Print("Error occurred when creating a recipe", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (rs *RecipeService) FindRecipesByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")

	if title == "" {
		log.Print("Expected title query parameter not present")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	results, err := rs.RecipeRepository.FindRecipesByTitle(title)

	if err != nil {
		log.Print("Error occurred when searching for recipes", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if results == nil || len(results) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func (rs *RecipeService) FindRecipesByIngredients(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ingredientsAsString := query.Get("ingredients")

	if ingredientsAsString == "" {
		log.Print("Expected ingredients query parameter not present")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ingredients := strings.Split(ingredientsAsString, ",")

	results, err := rs.RecipeRepository.FindRecipesByIngredients(ingredients)

	if err != nil {
		log.Print("Error occurred when searching for recipes ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if results == nil || len(results) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func (rs *RecipeService) FindRecipeById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Print("Cannot parse recipe id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recipe, err := rs.RecipeRepository.FindRecipeById(id)

	if recipe == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(recipe)
}
