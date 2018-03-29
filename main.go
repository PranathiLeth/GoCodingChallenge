package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/gorilla/mux"
	 _ "github.com/lib/pq"
	"github.com/PranathiLeth/GoCodingChallenge/todo"

)

// Status :=
func Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("Status Request Received")
	w.WriteHeader(200)
	fmt.Fprint(w, "OK\n")
}

func main() {
	router := httprouter.New()
	router.GET("/", Status)
	router.POST("/todos", todo.Create)
	router.GET("/todos", todo.List)
	router.PUT("/todos/{todoID}", todo.Update)

	log.Println("Starting server...")


	log.Fatal(http.ListenAndServe(":8080", router))
}
