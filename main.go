package main

import(
  "net/http"
  "html/template"
  "fmt"
  "log"
  "io/ioutil"
  "encoding/json"
)

type Todo struct{
  Description   string  `json:description`
  Status        bool    `json:status`
}

type TodoList struct{
  Todos         []Todo
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

func add(w http.ResponseWriter,r *http.Request){
  t,_ := template.ParseFiles("add.html")
  t.Execute(w,nil)
}

func list(w http.ResponseWriter,r *http.Request){
  //updating global variables
  body, err = ioutil.ReadFile("output.txt")
  if err != nil{
    log.Fatal(err)
  }
  err = json.Unmarshal(body,&todos)
  if err != nil{
    log.Fatal(err)
  }

  //parsing template
  t,_ := template.ParseFiles("list.html")
  t.Execute(w,TodoList{todos})
}

func save(w http.ResponseWriter,r *http.Request){
  //creating new object with form data
  description := r.FormValue("desc")
  newTodo := Todo{description,false}

  //inserting new Todo to global slice of Todos and updating body([]byte)
  todos = append(todos,newTodo)
  body,err = json.Marshal(todos)
  if err != nil{  log.Fatal(err)  }

  err := ioutil.WriteFile("output.txt",body, 0644)
  if err != nil{  log.Fatal(err)  }

  fmt.Printf("todo %v added\n",newTodo)
  http.Redirect(w,r,"/list/",http.StatusFound)
}

func main(){
  http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
    http.Redirect(w,r,"/add/",http.StatusFound)   })
  http.HandleFunc("/list/",list)
  http.HandleFunc("/add/",add)
  http.HandleFunc("/save/",save)
  http.ListenAndServe(":8090",nil)
}
