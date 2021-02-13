package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/krasimiraMilkova/cookit/pkg/comments"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type CommentService struct {
	CommentRepository comments.CommentRepository
}

var commentService *CommentService

func Get() *CommentService {
	if commentService == nil {
		commentService = &CommentService{CommentRepository: GetCommentRepository()}
	}

	return commentService
}

type Comment struct {
	Comment string `json:"comment"`
}

func (cs *CommentService) AddComment(w http.ResponseWriter, r *http.Request) {
	recipeId, err := strconv.Atoi(mux.Vars(r)["recipeId"])

	if err != nil {
		log.Print("Cannot parse recipe id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	comment := string(body)

	err = cs.CommentRepository.AddComment(recipeId, comment)

	if err != nil {
		if strings.Contains(err.Error(), "Cannot add or update") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cs *CommentService) GetComments(w http.ResponseWriter, r *http.Request) {
	recipeId, err := strconv.Atoi(mux.Vars(r)["recipeId"])

	if err != nil {
		log.Print("Cannot parse recipe id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundComments, err := cs.CommentRepository.GetComments(recipeId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(foundComments) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(foundComments)
}
