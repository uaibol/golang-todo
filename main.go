package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"encoding/json"

	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:pass1234@tcp(127.0.0.1:3306)/gotodo")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()

	//sm := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false,
		AllowedHeaders:   []string{"*"},
	})

	handler := c.Handler(router)
	
	router.HandleFunc("/", func(resp http.ResponseWriter, _ *http.Request) {

        fmt.Fprint(resp, "Home Page!")
    })
	router.HandleFunc("/todos", getTodos).Methods("GET")
	router.HandleFunc("/todos", createTodos).Methods("POST")
	router.HandleFunc("/todo/{id}", getTodo).Methods("GET")
	router.HandleFunc("/todo/{id}", updateTodo).Methods("PUT")
	router.HandleFunc("/todo/{id}", deleteTodo).Methods("DELETE")
	fmt.Println("Listen and Serving")
	http.ListenAndServe(":8081", handler)

}

func getTodos(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.URL.Path != "/todos" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var todos []Todo
	result, err := db.Query("SELECT id, title from todo")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var todo Todo

		err := result.Scan(&todo.ID, &todo.Title)
		if err != nil {
			panic(err.Error())
		}
		//fmt.Println(todo, "added")
		todos = append(todos, todo)
	}
	json.NewEncoder(w).Encode(todos)

}
func createTodos(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Path != "/todos" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	

	stmt, err := db.Prepare("INSERT INTO todo(title) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}  
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}  
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]  
	_, err = stmt.Exec(title)
	if err != nil {
		panic(err.Error())
	}  
	fmt.Fprintf(w, "Todo created")

}
func getTodo(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, title FROM todo WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var todo Todo
	for result.Next(){
		err := result.Scan(&todo.ID, &todo.Title)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(todo)


}
func updateTodo(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE todo SET title = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["title"]
	_, err = stmt.Exec(newTitle, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Todo with ID = %s updated", params["id"])

}
func deleteTodo(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM todo WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Todo with ID = %s deleted", params["id"])

}
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "custom 404")
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	//(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
    // (*w).Header().Set("Access-Control-Expose-Headers", "*")
    //(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
    // //(*w).Header().Set("Content-Security-Policy", "default-src *; style-src 'self' 'unsafe-inline'; font-src 'self' data:; script-src 'self' 'unsafe-inline' 'unsafe-eval' 127.0.0.1:8081")
}
