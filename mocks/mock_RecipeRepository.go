// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/recipes/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	recipes "github.com/krasimiraMilkova/cookit/pkg/recipes"
	reflect "reflect"
)

// MockRecipeRepository is a mock of RecipeRepository interface
type MockRecipeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRecipeRepositoryMockRecorder
}

// MockRecipeRepositoryMockRecorder is the mock recorder for MockRecipeRepository
type MockRecipeRepositoryMockRecorder struct {
	mock *MockRecipeRepository
}

// NewMockRecipeRepository creates a new mock instance
func NewMockRecipeRepository(ctrl *gomock.Controller) *MockRecipeRepository {
	mock := &MockRecipeRepository{ctrl: ctrl}
	mock.recorder = &MockRecipeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRecipeRepository) EXPECT() *MockRecipeRepositoryMockRecorder {
	return m.recorder
}

// CreateRecipe mocks base method
func (m *MockRecipeRepository) CreateRecipe(recipe *recipes.Recipe) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRecipe", recipe)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRecipe indicates an expected call of CreateRecipe
func (mr *MockRecipeRepositoryMockRecorder) CreateRecipe(recipe interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRecipe", reflect.TypeOf((*MockRecipeRepository)(nil).CreateRecipe), recipe)
}

// FindRecipesByTitle mocks base method
func (m *MockRecipeRepository) FindRecipesByTitle(title string) ([]recipes.RecipeSearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindRecipesByTitle", title)
	ret0, _ := ret[0].([]recipes.RecipeSearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindRecipesByTitle indicates an expected call of FindRecipesByTitle
func (mr *MockRecipeRepositoryMockRecorder) FindRecipesByTitle(title interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindRecipesByTitle", reflect.TypeOf((*MockRecipeRepository)(nil).FindRecipesByTitle), title)
}

// FindRecipesByIngredients mocks base method
func (m *MockRecipeRepository) FindRecipesByIngredients(ingredients []string) ([]recipes.RecipeSearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindRecipesByIngredients", ingredients)
	ret0, _ := ret[0].([]recipes.RecipeSearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindRecipesByIngredients indicates an expected call of FindRecipesByIngredients
func (mr *MockRecipeRepositoryMockRecorder) FindRecipesByIngredients(ingredients interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindRecipesByIngredients", reflect.TypeOf((*MockRecipeRepository)(nil).FindRecipesByIngredients), ingredients)
}

// FindRecipeById mocks base method
func (m *MockRecipeRepository) FindRecipeById(id int) (*recipes.Recipe, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindRecipeById", id)
	ret0, _ := ret[0].(*recipes.Recipe)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindRecipeById indicates an expected call of FindRecipeById
func (mr *MockRecipeRepositoryMockRecorder) FindRecipeById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindRecipeById", reflect.TypeOf((*MockRecipeRepository)(nil).FindRecipeById), id)
}
