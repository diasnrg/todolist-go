package main

import(
  "net/http"
  "log"
  "io/ioutil"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/rs/cors"
  "database/sql"
  _ "github.com/lib/pq"
)

type Todo struct{
  Description   string
  Status        bool
}

var(
  body    []byte
  todos   []Todo
  err     error
)

//initialising our todolist with data from .txt file
func init(){
  //creating the body([]byte)
  body,err = ioutil.ReadFile("output.txt")
  if err != nil{  log.Fatal(err)  }
  //converting body from json to struct([]Todo)
  err = json.Unmarshal(body,&todos)
  if err != nil{  log.Fatal(err)  }
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

  //inserting new Todo to global slice of Todos and updating the .txt file
  todos = append(todos,newTodo)
  updateTxt()

  log.Printf("todo %v added\n",newTodo)
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(newTodo)
}

func updateTxt(){
  //convert []Todo to json format([]byte)
  body,err = json.Marshal(todos)
  if err != nil{  log.Fatal(err)  }

  err := ioutil.WriteFile("output.txt",body, 0644)
  if err != nil{  log.Fatal(err)  }
}

func main(){

  db, err := sql.Open("postgres","host=localhost password=postgres dbname=todolist sslmode=disable")
  if err != nil{
    log.Fatal(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil{
    log.Fatal(err)
  }
  log.Println("Connected to the database")

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
