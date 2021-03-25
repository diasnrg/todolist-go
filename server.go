package main

import(
  "net/http"
  "log"
  "io/ioutil"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/rs/cors"
)

type Todo struct{
  Description   string  `json:description`
  Status        bool    `json:status`
}

var(
  body    []byte
  todos   []Todo
  err     error
)

func init(){
  //initialising our todolist with data from .txt file
  body,err = ioutil.ReadFile("output.txt")
  if err != nil{  log.Fatal(err)  }

  err = json.Unmarshal(body,&todos)
  if err != nil{  log.Fatal(err)  }
}

func list(w http.ResponseWriter,r *http.Request){
  //updating global variables
  body, err = ioutil.ReadFile("output.txt")
  if err != nil{  log.Fatal(err)  }

  err = json.Unmarshal(body,&todos)
  if err != nil{  log.Fatal(err)  }

  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(todos)
}

func save(w http.ResponseWriter,r *http.Request){
  //decoding the data from request body and writing it to todo struct
  var newTodo Todo
  err = json.NewDecoder(r.Body).Decode(&newTodo)
  // error will be when the json is not in correct format
  if err != nil{  log.Fatal(err)  }

  //inserting new Todo to global slice of Todos and updating body([]byte)
  todos = append(todos,newTodo)
  body,err = json.Marshal(todos)
  if err != nil{  log.Fatal(err)  }

  err := ioutil.WriteFile("output.txt",body, 0644)
  if err != nil{  log.Fatal(err)  }

  log.Printf("todo %v added\n",newTodo)
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(newTodo)
}

func main(){
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
