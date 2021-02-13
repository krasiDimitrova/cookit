package routes

import (
	"github.com/gorilla/mux"
	cs "github.com/krasimiraMilkova/cookit/internal/comments/service"
	rs "github.com/krasimiraMilkova/cookit/internal/recipes/service"
	"github.com/krasimiraMilkova/cookit/internal/users/auth"
	us "github.com/krasimiraMilkova/cookit/internal/users/service"
	"net/http"
)

func Handlers() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(CommonMiddleware)

	userService := us.Get()

	router.HandleFunc("/register", userService.CreateUser).Methods("POST")
	router.HandleFunc("/login", userService.Login).Methods("POST")

	jwtAuthenticator := auth.GetAuthenticator()
	authenticatedSubrouter := router.PathPrefix("/api/v1").Subrouter()
	authenticatedSubrouter.Use(jwtAuthenticator.VerifyJWT)

	recipeService := rs.Get()

	authenticatedSubrouter.HandleFunc("/recipe", recipeService.CreateRecipe).Methods("POST")
	authenticatedSubrouter.HandleFunc("/recipe/{id}", recipeService.FindRecipeById).Methods("GET")
	authenticatedSubrouter.HandleFunc("/recipe", recipeService.FindRecipesByTitle).Queries("title", "{title}").Methods("GET")
	authenticatedSubrouter.HandleFunc("/recipe", recipeService.FindRecipesByIngredients).Queries("ingredients", "{ingredients}").Methods("GET")

	commentService := cs.Get()
	authenticatedSubrouter.HandleFunc("/recipe/{recipeId}/comment", commentService.GetComments).Methods("GET")
	authenticatedSubrouter.HandleFunc("/recipe/{recipeId}/comment", commentService.AddComment).Methods("POST")

	return router
}

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
