package todo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Create will allow a user to create a new todo
// The supported body is {"title": "", "status": ""}
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dbUser := os.Getenv("Admin")
	dbHost := os.Getenv("LAPTOP-PIA8TBSK")
	dbPassword := os.Getenv("Admin")
	dbName := os.Getenv("ToDo")

	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Println(err.Error())
	}
	var todo CreateTodo

	json.NewDecoder(r.Body).Decode(&todo)

	if todo.Status == "" || todo.Title == "" {
		http.Error(w, "Todo request is missing status or title", http.StatusBadRequest)
	}

	invalidStatus := true
	for _, status := range allowedStatuses {
		if todo.Status == status {
			invalidStatus = false
			break
		}
	}

	if !invalidStatus {
		http.Error(w, "The provided status is not supported", http.StatusBadRequest)
	}

	insertStmt := fmt.Sprintf(`INSERT INTO todo (title, status) VALUES ('%s', '%s') RETURNING id`, todo.Title, todo.Status)

	var todoID int

	// Insert and get back newly created todo ID
	if err := db.QueryRow(insertStmt).Scan(&todoID); err != nil {
		fmt.Printf("Failed to save to db: %s", err.Error())
	}

	fmt.Printf("Todo Created -- ID: %d\n", todoID)

	newTodo := Todo{}
	db.QueryRow("SELECT id, title, status FROM todo WHERE id=$1", todoID).Scan(&newTodo.ID, &newTodo.Title, &newTodo.Status)

	jsonResp, _ := json.Marshal(newTodo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, string(jsonResp))
}

// List will provide a list of all current to-dos
func List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dbUser := os.Getenv("Admin")
	dbHost := os.Getenv("LAPTOP-PIA8TBSK")
	dbPassword := os.Getenv("Admin")
	dbName := os.Getenv("ToDo")

	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Println(err.Error())
	}

	todoList := []Todo{}

	rows, err := db.Query("SELECT id, title, status FROM todo")
	defer rows.Close()

	for rows.Next() {
		todo := Todo{}
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status); err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Failed to build todo list")
		}

		todoList = append(todoList, todo)
	}

	jsonResp, _ := json.Marshal(Todos{TodoList: todoList})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, string(jsonResp))
}

//update API to update a todo based on id

func Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	dbUser := os.Getenv("Admin")
	dbHost := os.Getenv("LAPTOP-PIA8TBSK")
	dbPassword := os.Getenv("Admin")
	dbName := os.Getenv("ToDo")

	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Println(err.Error())
	}

	id, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	var todo UpdateTodo

	json.NewDecoder(r.Body).Decode(&todo)

	row, err := db.QueryRow("SELECT title, status FROM todo WHERE id=$1", id)
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = db.Exec("UPDATE todo SET Title = ?, Status = ? WHERE Id = ?", todo.Title, todo.Status, id)

	if err != nil {
		fmt.Println("ERROR saving to db - ", err)
	}

	newTodo := Todo{}
	err = db.QueryRow("SELECT  Title, Status FROM Todo WHERE Id=?", newTodo.ID).Scan(&newTodo.Title, &newTodo.Status)
	if err != nil {
		fmt.Println("ERROR reading from db - ", err)
	}

	jsonResp, _ := json.Marshal(newTodo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, string(jsonResp))
}
