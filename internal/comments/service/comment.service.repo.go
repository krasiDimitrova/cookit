package service

import (
	"database/sql"
	"errors"
	"github.com/krasimiraMilkova/cookit/internal/db"
	"github.com/krasimiraMilkova/cookit/pkg/comments"
)

type CommentRepository struct {
	*sql.DB
}

func GetCommentRepository() comments.CommentRepository {
	return &CommentRepository{db.Get()}
}

func (commentRepository *CommentRepository) AddComment(recipeId int, comment string) error {
	if comment == "" {
		return errors.New("comment cannot be empty")
	}

	_, err := commentRepository.Exec("insert into comments(recipe_id, comment)values(?,?);", recipeId, comment)

	if err != nil {
		return err
	}

	return nil
}

func (commentRepository *CommentRepository) GetComments(recipeId int) ([]string, error) {
	resultRows, err := commentRepository.Query("select comment from comments where recipe_id = ?", recipeId)

	if err != nil {
		return nil, err
	}

	var foundComments []string

	defer resultRows.Close()
	for resultRows.Next() {
		var comment string
		err = resultRows.Scan(&comment)

		if err == nil {
			foundComments = append(foundComments, comment)
		}
	}

	return foundComments, nil
}
