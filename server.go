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

type Task struct{
  Id            int
  Description   string
  Status        bool
}

var(
  //database object
  db      *sql.DB
  //map to keep tasks in the global variable
  tasks   map[int]Task
  err     error
)

func init(){
  tasks = make(map[int]Task)
  connectDB()
  initTasks()
}

//connect to the database and establish the connection with Ping
func connectDB(){
  db, err = sql.Open("postgres","host=localhost password=postgres dbname=todolist sslmode=disable")
  if err != nil{  log.Fatal(err)  }

  err = db.Ping()
  if err != nil{  log.Fatal(err)  }
  log.Println("connected to the database...")
}

//retrieve rows from database to initialise the global map with tasks
func initTasks(){
  rows, err := db.Query("SELECT id, description, status FROM todos")
  if err != nil{  log.Fatal(err)  }
  defer rows.Close()

  var t Task
  for rows.Next(){
    err = rows.Scan(&t.Id, &t.Description, &t.Status)
    tasks[t.Id] = t
    if err != nil { log.Fatal(err) }
  }

  if err = rows.Err(); err != nil{
    log.Fatal(err)
  }
}

//encode the tasks(map) into JSON to response the client request
func getItems(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(tasks)
}

//create the new task from request's body and save it to the database
func addItem(w http.ResponseWriter,r *http.Request){
  var t Task
  //decoding the data from request's body and writing it to the new task struct
  err = json.NewDecoder(r.Body).Decode(&t)
  //error will be if the json is not in the correct format
  if err != nil{  log.Fatal(err)  }

  //insert the new task to the database, return it's new id
  err := db.QueryRow("INSERT INTO todos(description) VALUES($1) RETURNING id",t.Description).Scan(&t.Id)
  if err != nil{  log.Fatal(err)  }

  //save the new task to tasks(map)
  tasks[t.Id] = t
  log.Printf("task {id:%v, description:%v, status:%v} added\n", t.Id, t.Description, t.Status)

  //response with the new task's JSON body
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(t)
}

//delete the task using id
func deleteItem(w http.ResponseWriter,r *http.Request){
  //mux's method which takes the values after the '/delete/' and saves them to map 'vars'
  vars := mux.Vars(r)
  //retrieve the id from the vars(map) and convert it (string->int)
  id, _ := strconv.Atoi(vars["id"])

  //delete the task from the database
  _, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
  if err != nil{  log.Fatal(err)  }

  t := tasks[id]
  log.Printf("task {id:%v, description:%v, status:%v} deleted", id, t.Description, t.Status)

  //delete tasks from tasks(map)
  delete(tasks, id)

  //response with the deleted task's JSON body
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(t)
}

//update task's status (true/false)
func updateItem(w http.ResponseWriter,r *http.Request){
  //retrieve the id of the task
  vars := mux.Vars(r)
  id, _ := strconv.Atoi(vars["id"])
  //new status of the current task (negation of the previous one)
  status := !tasks[id].Status

  //updating task's status in the database
  _, err := db.Exec("UPDATE todos SET Status = $2 WHERE Id = $1", id, status)
  if err != nil{  log.Fatal(err)  }

  //update current task's status in the tasks(map)
  tasks[id] = Task{id, tasks[id].Description, status}
  log.Printf("task {id:%v, description:%v} updated to %v", id, tasks[id].Description, status)

  //response with the updated task's new status
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(status)
}

func main(){
  //close the db in the end
  defer db.Close()
  log.Printf("tasks: %v", tasks)

  //routing the requests
  r := mux.NewRouter()
  r.HandleFunc("/tasks/",getItems).Methods("GET")
  r.HandleFunc("/add/",addItem).Methods("POST")
  r.HandleFunc("/delete/{id}",deleteItem).Methods("DELETE")
  r.HandleFunc("/update/{id}",updateItem).Methods("POST")

  //wrapping the CORS header to be able to connect our API to the front end
  handler := cors.New(cors.Options{
    AllowedMethods: []string{"GET","POST","DELETE"},
  }).Handler(r)

  http.Handle("/",r)
  //listening to the port 8090
  http.ListenAndServe(":8090",handler)
}
