package service

import (
	todo_app_go "todo-app-go"
	"todo-app-go/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo_app_go.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}
type TodoList interface {
	Create(userId int, list todo_app_go.TodoList) (int, error)
	GetAll(userId int) ([]todo_app_go.TodoList, error)
	GetById(userId, listId int) (todo_app_go.TodoList, error)
	Update(userId, listId int, input todo_app_go.UpdateListInput) error
	Delete(userId, listId int) error
}

type TodoItem interface {
	Create(userId, listId int, item todo_app_go.TodoItem) (int, error)
	GetAll(userId, lisId int) ([]todo_app_go.TodoItem, error)
	GetItemById(userId, itemId int) (todo_app_go.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input todo_app_go.UpdateItemInput) error
}
type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
