package main

import (
	"client/internal/menu"
	"fmt"
	"net/http"
)

func main() {
	userMenu := menu.GetUserMenu()
	var tokenCookies []*http.Cookie

	for ; len(tokenCookies) == 0; {
		tokenCookies = userMenu.GetAccessToken()
	}

	quit := make(chan bool)
	recipeMenu := menu.GetRecipeMenu(tokenCookies, quit)
	recipeMenu.PrintMenu()
	recipeMenu.RecipeMenuChannel <- 3

	<-quit
	fmt.Println("Goodbye")
}
