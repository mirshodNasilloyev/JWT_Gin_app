package repository

import (
	"fmt"
	"strings"
	todo_app_go "todo-app-go"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ToDoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *ToDoItemPostgres {
	return &ToDoItemPostgres{
		db: db,
	}
}

func (r *ToDoItemPostgres) Create(listId int, item todo_app_go.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createQueryItem := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id", todoItemsTable)

	row := tx.QueryRow(createQueryItem, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listItemsTable)
	_, err = tx.Exec(createListItemQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

		return itemId, tx.Commit()
	}
// func (r *ToDoItemPostgres) Create(listId int, item todo_app_go.TodoItem) (int, error) {
// 	// Check if the item already exists
// 	var exists bool
// 	checkQuery := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE title=$1 AND description=$2 AND done=$3)", todoItemsTable)
// 	err := r.db.QueryRow(checkQuery, item.Title, item.Description, item.Done).Scan(&exists)
// 	if err != nil {
// 		return 0, err
// 	}
// 	if exists {
// 		return 0, fmt.Errorf("duplicate entry: item already exists")
// 	}

// 	// Begin transaction
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return 0, err
// 	}

// 	var itemId int
// 	createQueryItem := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
// 	row := tx.QueryRow(createQueryItem, item.Title, item.Description)
// 	err = row.Scan(&itemId)
// 	if err != nil {
// 		tx.Rollback()
// 		return 0, fmt.Errorf("failed to insert todo_item: %v", err)
// 	}

// 	createListItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listItemsTable)
// 	_, err = tx.Exec(createListItemQuery, listId, itemId)
// 	if err != nil {
// 		tx.Rollback()
// 		return 0, fmt.Errorf("failed to insert list_item: %v", err)
// 	}

// 	if err = tx.Commit(); err != nil {
// 		return 0, fmt.Errorf("transaction commit failed: %v", err)
// 	}
// 	return itemId, nil
// }

func (r *ToDoItemPostgres) GetAll(userId, listId int) ([]todo_app_go.TodoItem, error) {
	var items []todo_app_go.TodoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description FROM %s ti INNER JOIN %s li on li.item_id = ti.id INNER JOIN %s ul on ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2`, todoItemsTable, listItemsTable, usersListsTable)
	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ToDoItemPostgres) GetItemById(userId, itemId int) (todo_app_go.TodoItem, error) {
	var item todo_app_go.TodoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description FROM %s ti INNER JOIN %s li on li.item_id = ti.id
	 INNER JOIN %s ul on ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, listItemsTable, usersListsTable)
	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}

	return item, nil
}

func (r *ToDoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul WHERE ti.id=li.item_id AND li.list_id=ul.list_id AND ul.user_id=$1 AND ti.id=$2`,
		todoItemsTable, listItemsTable, usersListsTable)
	_, err := r.db.Exec(query, itemId, userId)
	return err
}

func (r *ToDoItemPostgres) Update(userId, itemId int, input todo_app_go.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}
	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul 
						  WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
