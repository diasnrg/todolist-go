package main

import(
  "net/http"
  "log"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/rs/cors"
  "database/sql"
  _ "github.com/lib/pq"
  "strconv"
)

type Todo struct{
  Id            int
  Description   string
  Status        bool
}

var(
  db      *sql.DB
  todos   map[int]Todo
  err     error
)

func init(){
  todos = make(map[int]Todo)
}

func connectDB(){
  db, err = sql.Open("postgres","host=localhost password=postgres dbname=todolist sslmode=disable")
  if err != nil{  log.Fatal(err)  }

  err = db.Ping()
  if err != nil{  log.Fatal(err)  }
  log.Println("connected to the database...")
}

func initTodos(){
  rows, err := db.Query("SELECT id, description, status FROM todos")
  if err != nil{  log.Fatal(err)  }
  defer rows.Close()

  var newTodo Todo
  for rows.Next(){
    err = rows.Scan(&newTodo.Id, &newTodo.Description, &newTodo.Status)
    todos[newTodo.Id] = newTodo
    if err != nil { log.Fatal(err) }
  }

  if err = rows.Err(); err != nil{
    log.Fatal(err)
  }
}

func list(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(todos)
}

func addItem(w http.ResponseWriter,r *http.Request){
  var newTodo Todo
  //decoding the data from request's body and writing it to the new todo struct, error will be when the json is not in correct format
  err = json.NewDecoder(r.Body).Decode(&newTodo)
  if err != nil{  log.Fatal(err)  }

  //insert the new Todo to database, return it's new id
  err := db.QueryRow("INSERT INTO todos(description) VALUES($1) RETURNING id",newTodo.Description).Scan(&newTodo.Id)
  if err != nil{  log.Fatal(err)  }

  todos[newTodo.Id] = newTodo
  log.Printf("{id:%v, description:%v, status:%v} added\n", newTodo.Id, newTodo.Description, newTodo.Status)
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(newTodo)
}

func deleteItem(w http.ResponseWriter,r *http.Request){
  //mux's method which takes the values after /delete/ and saves them to map 'vars'
  vars := mux.Vars(r)
  id, _ := strconv.Atoi(vars["id"])

  _, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
  if err != nil{  log.Fatal(err)  }

  delete(todos, id)
  log.Printf("todo with id=%v deleted", id)
}

func updateItem(w http.ResponseWriter,r *http.Request){
  vars := mux.Vars(r)
  id, _ := strconv.Atoi(vars["id"])
  newStatus := !todos[id].Status

  _, err := db.Exec("UPDATE todos SET Status = $2 WHERE Id = $1", id, newStatus)
  if err != nil{  log.Fatal(err)  }

  todos[id] = Todo{id,todos[id].Description,newStatus}
  log.Printf("todo with id=%v updated to %v", id, newStatus)
}

func main(){
  connectDB()
  initTodos()

  defer db.Close()
  log.Println(todos)

  r := mux.NewRouter()
  r.HandleFunc("/list/",list).Methods("GET")
  r.HandleFunc("/add/",addItem).Methods("POST")
  r.HandleFunc("/delete/{id}",deleteItem).Methods("DELETE")
  r.HandleFunc("/update/{id}",updateItem).Methods("POST")

  //for CORS error
  handler := cors.New(cors.Options{
    AllowedMethods: []string{"GET","POST","DELETE","UPDATE"},
  }).Handler(r)

  http.Handle("/",r)
  http.ListenAndServe(":8090",handler)
}
