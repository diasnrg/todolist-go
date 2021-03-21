package main

import(
  "net/http"
  "html/template"
  "fmt"
  "log"
  "io/ioutil"
  "strings"
)

type todo struct{
  description string
  status bool
}

type todolist struct{
  Todos []string
}

var(
  body []byte
  err error
)

func init(){
  body,err = ioutil.ReadFile("output.txt")
  check(err)
}

func add(w http.ResponseWriter,r *http.Request){
  t,_ := template.ParseFiles("add.html")
  t.Execute(w,nil)
}

func list(w http.ResponseWriter,r *http.Request){
  body, err := ioutil.ReadFile("output.txt")
  check(err)
  todos := strings.Split(string(body),"\n")
  t,_ := template.ParseFiles("list.html")
  t.Execute(w,todolist{todos})
}

func save(w http.ResponseWriter,r *http.Request){
  description := r.FormValue("desc")
  newTodo := todo{description,false}

  bodyString := string(body) + fmt.Sprintf("%v\n",newTodo)
  err := ioutil.WriteFile("output.txt", []byte(bodyString), 0644)
  check(err)

  fmt.Printf("todo %v added\n",newTodo)
  http.Redirect(w,r,"/list/",http.StatusFound)
}

func check(err error){
  if err != nil{
    log.Fatal(err)
  }
}

func (t todo) String() string{
  return fmt.Sprintf("%v:%v",t.description,t.status)
}

func main(){
  http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
    http.Redirect(w,r,"/add/",http.StatusFound)
  })
  http.HandleFunc("/list/",list)
  http.HandleFunc("/add/",add)
  http.HandleFunc("/save/",save)
  http.ListenAndServe(":8090",nil)
}
