package menu

import (
	"bufio"
	"client/internal/apis"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type RecipeMenu struct {
	RecipeApi         *apis.RecipeApi
	RecipeMenuChannel chan int
	quit              chan bool
}

var recipeMenu *RecipeMenu

func GetRecipeMenu(cookies []*http.Cookie, quit chan bool) *RecipeMenu {
	if recipeMenu == nil {
		recipeMenu = &RecipeMenu{
			RecipeApi:         apis.GetRecipeApi(cookies),
			RecipeMenuChannel: make(chan int, 1),
			quit:              quit,
		}
	}

	return recipeMenu
}

// PrintMenu function handles stdin/stdout operations for main menu operations
// including creating, searching and reading the recipes as well as adding and reading comments
func (rm *RecipeMenu) PrintMenu() {
	go func() {
		for {
			switch <-rm.RecipeMenuChannel {
			case 1:
				rm.printCreateRecipe()
			case 2:
				rm.printSearch()
			case 4:
				rm.quit <- true
			default:
				rm.printMainMenu()
			}
		}
	}()
}

func (rm *RecipeMenu) printMainMenu() {
	var command int
	fmt.Print("Create recipe (1), search for recipe (2), main menu (3), quit (4): ")
	fmt.Scan(&command)

	rm.RecipeMenuChannel <- command
}

func (rm *RecipeMenu) printSearch() {
	fmt.Print("Search by title (1) or by ingredients (2) or exit this menu (3): ")

	var command int
	fmt.Scan(&command)

	var found []apis.RecipeSearchResult
	reader := bufio.NewReader(os.Stdin)
	switch command {
	case 1:
		{
			fmt.Print("Enter title or part of the title: ")
			title, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println("Failed to read title. Try again!")
				rm.RecipeMenuChannel <- 2
				return
			}

			title = strings.Trim(title, "\n")
			found, err = rm.RecipeApi.FindByTitle(title)
		}
	case 2:
		{
			fmt.Print("Enter names of ingredients (separated by ,): ")
			ingredients, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println("Failed to read ingredients. Try again!")
				rm.RecipeMenuChannel <- 2
				return
			}

			ingredients = strings.Trim(ingredients, "\n")
			found, err = rm.RecipeApi.FindByIngredients(ingredients)
		}
	default:
		rm.RecipeMenuChannel <- 3
		return
	}

	if len(found) == 0 {
		fmt.Println("No recipes found. Try again!")
		rm.RecipeMenuChannel <- 2
		return
	}

	rm.printSearchResults(found)
	rm.RecipeMenuChannel <- 3
}

func (rm *RecipeMenu) printSearchResults(searchResults []apis.RecipeSearchResult) {
	if len(searchResults) == 0 {
		fmt.Println("No recipes found")
		return
	}

	for i, sr := range searchResults {
		fmt.Println(strconv.Itoa(i) + " - " + sr.Title)
	}

	var index int

	for {
		fmt.Print("Choose recipe (enter #):")
		_, err := fmt.Scan(&index)

		if err == nil {
			break
		}

		fmt.Println("Failed to decode id. Try again!")
	}

	id := searchResults[index].ID
	rm.printRecipe(id)
}

func (rm *RecipeMenu) printRecipe(id int) {
	recipe, err := rm.RecipeApi.GetById(id)

	if err != nil {
		fmt.Println("Could not load the recipe.")
		return
	}

	fmt.Println(recipe.Title)

	for _, ingredient := range recipe.Ingredients {
		fmt.Println(ingredient.Name + " - " + strconv.Itoa(ingredient.Quantity) + " " + ingredient.Measurement)
	}

	fmt.Println(recipe.Directions)
	rm.printCommentsMenu(id)
}

func (rm *RecipeMenu) printCommentsMenu(recipeId int) {
	var command int
	for ; command != 3; {
		fmt.Println("Print comments for recipe (1), add comment to recipe (2), exit recipe (3): ")
		fmt.Scanln(&command)

		switch command {
		case 1:
			rm.printComments(recipeId)
		case 2:
			rm.printAddComment(recipeId)
		default:
			break
		}
	}
}

func (rm *RecipeMenu) printComments(recipeId int) {
	comments, err := rm.RecipeApi.GetCommentsById(recipeId)

	if err != nil {
		fmt.Println("Failed to fetch comments")
	} else if len(comments) == 0 {
		fmt.Println("No comments found")
	} else {
		fmt.Println("Comments:")
		for _, comment := range comments {
			fmt.Println(comment)
		}
	}
}

func (rm *RecipeMenu) printAddComment(recipeId int) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter comment: ")
	comment, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Cannot read comment. Try again!")
	} else {
		comment = strings.Trim(comment, "\n")
		err = rm.RecipeApi.AddCommentById(recipeId, comment)
		if err != nil {
			fmt.Println("Failed to add comment")
		}
	}
}

func (rm *RecipeMenu) printCreateRecipe() {
	fmt.Println("Create recipe")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter title: ")
	title, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Failed to load the title.")
		rm.RecipeMenuChannel <- 3
		return
	}
	title = strings.Trim(title, "\n")

	var ingredients []apis.Ingredient
	for hasNext := true; hasNext; {
		fmt.Println("Enter information for ingredient")
		ingredient := readIngredient()
		ingredients = append(ingredients, ingredient)
		fmt.Print("Has more ingredients? (y/n): ")
		var more string
		fmt.Scan(&more)
		hasNext = more == "y"
	}

	fmt.Print("Enter directions: ")
	directions, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Failed to load the directions.")
		rm.RecipeMenuChannel <- 3
		return
	}
	directions = strings.Trim(directions, "\n")

	err = rm.RecipeApi.CreateRecipe(apis.Recipe{
		Title:       title,
		Ingredients: ingredients,
		Directions:  directions,
	})

	if err != nil {
		fmt.Println("Failed to create the recipe")
	} else {
		fmt.Println("Recipe is created")
	}

	rm.RecipeMenuChannel <- 3
}

func readIngredient() apis.Ingredient {
	var name, measurement string
	var quantity int
	fmt.Print("Name: ")
	fmt.Scanln(&name)
	fmt.Print("Quantity: ")
	fmt.Scanln(&quantity)
	fmt.Print("Measurement: ")
	fmt.Scanln(&measurement)

	return apis.Ingredient{
		Name:        name,
		Quantity:    quantity,
		Measurement: measurement,
	}
}
