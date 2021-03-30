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
  todos   []Todo
  err     error
)

func connectDB(){
  db, err = sql.Open("postgres","host=localhost password=postgres dbname=todolist sslmode=disable")
  if err != nil{  log.Fatal(err)  }

  err = db.Ping()
  if err != nil{  log.Fatal(err)  }
  log.Println("Connected to the database")
}

func updateTodos(){
  var id int
  var description string
  var status bool

  rows, err := db.Query("select id, description, status from todos")
  if err != nil{  log.Fatal(err)  }
  defer rows.Close()

  for rows.Next(){
    err = rows.Scan(&id, &description, &status)
    todos = append(todos,Todo{id, description, status})
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

func save(w http.ResponseWriter,r *http.Request){
  var newTodo Todo
  //decoding the data from request's body and writing it to the new todo struct, error will be when the json is not in correct format
  err = json.NewDecoder(r.Body).Decode(&newTodo)
  if err != nil{  log.Fatal(err)  }

  todos = append(todos,newTodo)
  log.Printf("todo %v added\n",newTodo)
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(newTodo)
}

func main(){
  connectDB()
  updateTodos()

  defer db.Close()
  log.Println(todos)

  r := mux.NewRouter()
  r.HandleFunc("/list/",list).Methods("GET")
  r.HandleFunc("/save/",save).Methods("POST")

  //for CORS error
  handler := cors.New(cors.Options{
    AllowedMethods: []string{"GET","POST"},
  }).Handler(r)

  http.Handle("/",r)
  http.ListenAndServe(":8090",handler)
}
