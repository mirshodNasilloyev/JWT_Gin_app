package repository

import (
	todo_app_go "todo-app-go"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todo_app_go.User) (int, error)
	GetUser(username, password string) (todo_app_go.User, error)
}
type TodoList interface {
	Create(userId int, list todo_app_go.TodoList) (int, error)
	GetAll(userId int) ([]todo_app_go.TodoList, error)
	GetById(userId, listId int) (todo_app_go.TodoList, error)
	Update(userId, listId int, input todo_app_go.UpdateListInput) error
	Delete(userId, listId int) error
}

type TodoItem interface {
	Create(listId int, item todo_app_go.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo_app_go.TodoItem, error)
	GetItemById(userId, itemId int) (todo_app_go.TodoItem, error)
	Update(userId, itemId int, input todo_app_go.UpdateItemInput) error
	Delete(userId, itemId int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}

}
