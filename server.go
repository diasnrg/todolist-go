package main

import(
  "net/http"
  "log"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/rs/cors"
  "database/sql"
  _ "github.com/lib/pq"
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
  log.Println("Connected to the database")
}

func updateTodos(){
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

func add(w http.ResponseWriter,r *http.Request){
  var newTodo Todo
  //decoding the data from request's body and writing it to the new todo struct, error will be when the json is not in correct format
  err = json.NewDecoder(r.Body).Decode(&newTodo)
  if err != nil{  log.Fatal(err)  }

  //insert the new Todo to database, return it's new id
  err := db.QueryRow("INSERT INTO todos(description) VALUES($1) RETURNING id",newTodo.Description).Scan(&newTodo.Id)
  if err != nil{  log.Fatal(err)  }

  log.Printf("{id:%v, description:%v, status:%v} added\n", newTodo.Id, newTodo.Description, newTodo.Status)
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(newTodo)
}

func delete(w http.ResponseWriter,r *http.Request){
  //mux's method which takes the values after /delete/ and saves them to map 'vars'
  vars := mux.Vars(r)
  id := vars["id"]

  _, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
  if err != nil{  log.Fatal(err)  }
  log.Printf("todo with id=%v deleted", id)
}

func main(){
  connectDB()
  updateTodos()

  defer db.Close()
  log.Println(todos)

  r := mux.NewRouter()
  r.HandleFunc("/list/",list).Methods("GET")
  r.HandleFunc("/add/",add).Methods("POST")
  r.HandleFunc("/delete/{id}",delete).Methods("DELETE")

  //for CORS error
  handler := cors.New(cors.Options{
    AllowedMethods: []string{"GET","POST","DELETE"},
  }).Handler(r)

  http.Handle("/",r)
  http.ListenAndServe(":8090",handler)
}
