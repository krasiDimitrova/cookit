package service

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/krasimiraMilkova/cookit/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestCommentService_AddComment(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mocks.NewMockCommentRepository(mockCtrl)

	service := CommentService{CommentRepository: mockRepository}

	tests := []struct {
		name               string
		comment            string
		recipeId           string
		repositoryError    string
		expectedStatusCode int
	}{
		{
			name:               "Successful",
			comment:            "Some comment",
			recipeId:           "1",
			repositoryError:    "",
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Cannot parse id",
			comment:            "Smth",
			recipeId:           "a",
			repositoryError:    "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Recipe not found",
			comment:            "Some comment",
			recipeId:           "2",
			repositoryError:    "Cannot add or update",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Repository error",
			comment:            "Some comment",
			recipeId:           "2",
			repositoryError:    "some other error",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonComment, _ := json.Marshal(test.comment)
			req, _ := http.NewRequest("POST", "/recipe/"+test.recipeId+"/comment", strings.NewReader(string(jsonComment)))
			req = mux.SetURLVars(req, map[string]string{
				"recipeId": test.recipeId,
			})
			rr := httptest.NewRecorder()

			var err error
			if test.repositoryError != "" {
				err = errors.New(test.repositoryError)
			}

			if test.expectedStatusCode != http.StatusBadRequest {
				id, _ := strconv.Atoi(test.recipeId)
				mockRepository.EXPECT().AddComment(id, string(jsonComment)).Return(err)
			}
			http.HandlerFunc(service.AddComment).ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, test.expectedStatusCode)
				t.Fail()
			}
		})
	}
}

func TestCommentService_GetComments(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mocks.NewMockCommentRepository(mockCtrl)

	service := CommentService{CommentRepository: mockRepository}

	tests := []struct {
		name               string
		comments           []string
		recipeId           string
		repositoryError    string
		expectedStatusCode int
	}{
		{
			name:               "Successful",
			comments:           []string{"Some comment", "Another"},
			recipeId:           "1",
			repositoryError:    "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Cannot parse id",
			comments:           []string{},
			recipeId:           "a",
			repositoryError:    "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "No comments found",
			comments:           []string{},
			recipeId:           "2",
			repositoryError:    "",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Repository error",
			comments:           []string{},
			recipeId:           "2",
			repositoryError:    "some error",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/recipe/"+test.recipeId+"/comment", nil)
			req = mux.SetURLVars(req, map[string]string{
				"recipeId": test.recipeId,
			})
			rr := httptest.NewRecorder()

			var err error
			if test.repositoryError != "" {
				err = errors.New(test.repositoryError)
			}

			if test.expectedStatusCode != http.StatusBadRequest {
				id, _ := strconv.Atoi(test.recipeId)
				mockRepository.EXPECT().GetComments(id).Return(test.comments, err)
			}
			http.HandlerFunc(service.GetComments).ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, test.expectedStatusCode)
				t.Fail()
			}

			var resultComments []string
			s := rr.Body.String()
			if s != "" {
				err = json.Unmarshal([]byte(s), &resultComments)

				if err != nil {
					t.Error("error from unmarshal", err)
				}
			}

			if len(resultComments) > 0 && !reflect.DeepEqual(resultComments, test.comments) {
				t.Errorf("Got searchResults = %v but wanted %v", resultComments, test.comments)
			} else if len(test.comments) == 0 && len(resultComments) != 0 {
				t.Errorf("Expected [] result but got %v", resultComments)
			}
		})
	}
}
