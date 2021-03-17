package main

import(
  "net/http"
  "html/template"
  "fmt"
)

type todo struct{
  descrition string
  status bool
}

func add(w http.ResponseWriter,r *http.Request){
  t,_ := template.ParseFiles("index.html")
  t.Execute(w,nil)
}

func save(w http.ResponseWriter,r *http.Request){
  description := r.FormValue("desc")
  newTodo := todo{description,false}
  fmt.Println(newTodo)
}

func main(){
  http.HandleFunc("/add/",add)
  http.HandleFunc("/save/",save)
  http.ListenAndServe(":8090",nil)

}
