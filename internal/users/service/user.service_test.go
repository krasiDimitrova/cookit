package service

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/krasimiraMilkova/cookit/internal/users/auth"
	"github.com/krasimiraMilkova/cookit/mocks"
	"github.com/krasimiraMilkova/cookit/pkg/users"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mocks.NewMockUserRepository(mockCtrl)
	mockAuthenticator := mocks.NewMockUserAuthenticator(mockCtrl)

	service := UserService{UserRepository: mockRepository, UserAuthenticator: mockAuthenticator}

	tests := []struct {
		name               string
		user               users.User
		existUser          bool
		repositoryError    string
		expectedStatusCode int
	}{
		{
			name: "Successful",
			user: users.User{
				Name:     "Test",
				Email:    "test@test.com",
				Password: "test",
			},
			existUser:          false,
			repositoryError:    "",
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "Existing user",
			user: users.User{
				Name:     "Test",
				Email:    "test@test.com",
				Password: "test",
			},
			existUser:          true,
			repositoryError:    "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Repo error during create",
			user: users.User{
				Name:     "Test",
				Email:    "test@test.com",
				Password: "test",
			},
			existUser:          false,
			repositoryError:    "some error",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonUser, _ := json.Marshal(test.user)
			req, _ := http.NewRequest("POST", "/register", strings.NewReader(string(jsonUser)))
			rr := httptest.NewRecorder()

			mockRepository.EXPECT().ExistUser(test.user.Email).Return(test.existUser)

			if !test.existUser {
				mockExpect := mockRepository.EXPECT().CreateUser(&test.user)

				if test.repositoryError == "" {
					mockExpect.Return(nil)
				} else {
					mockExpect.Return(errors.New(test.repositoryError))
				}
			}

			http.HandlerFunc(service.CreateUser).ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, test.expectedStatusCode)
				t.Fail()
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mocks.NewMockUserRepository(mockCtrl)
	mockAuthenticator := mocks.NewMockUserAuthenticator(mockCtrl)

	service := UserService{UserRepository: mockRepository, UserAuthenticator: mockAuthenticator}

	tests := []struct {
		name               string
		searchEmail        string
		searchPassword     string
		user               users.User
		repositoryError    error
		authenticatorError error
		expectedStatusCode int
	}{
		{
			name:           "Successful",
			searchEmail:    "test@test.com",
			searchPassword: "test",
			user: users.User{
				ID:       1,
				Name:     "Test",
				Email:    "test@test.com",
				Password: "test",
			},
			repositoryError:    nil,
			authenticatorError: nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Repo error",
			searchEmail:        "test@test.com",
			searchPassword:     "test",
			user:               users.User{},
			repositoryError:    errors.New("empty result set"),
			authenticatorError: nil,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Authenticator error",
			searchEmail:        "test@test.com",
			searchPassword:     "test",
			user:               users.User{},
			repositoryError:    nil,
			authenticatorError: errors.New("error during create"),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			loginUser := users.User{
				Email:    test.searchEmail,
				Password: test.searchPassword,
			}
			jsonUser, _ := json.Marshal(loginUser)
			req, _ := http.NewRequest("POST", "/login", strings.NewReader(string(jsonUser)))
			rr := httptest.NewRecorder()

			mockRepository.EXPECT().FindUser(test.searchEmail, test.searchPassword).Return(&test.user, test.repositoryError)

			if test.expectedStatusCode != http.StatusNotFound {
				mockAuthenticator.EXPECT().GenerateTokenForUser(&test.user).Return("someToken", test.authenticatorError)
			}

			http.HandlerFunc(service.Login).ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.expectedStatusCode)
				t.Fail()
			}

			cookies := rr.Result().Cookies()

			if test.expectedStatusCode == http.StatusOK {
				if len(cookies) != 1 {
					t.Errorf("expected 1 cookie got %v", len(cookies))
					t.Fail()
				}

				if cookies[0].Name != auth.TokenName {
					t.Errorf("expected cookie with name %v got %v", cookies[0].Name, auth.TokenName)
					t.Fail()
				}
			} else if len(cookies) > 0 {
				t.Errorf("expected 0 cookies got %v", len(cookies))
				t.Fail()
			}
		})
	}
}
