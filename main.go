package main

import(
  "net/http"
  "html/template"
)

type todo struct{
  descrition string
  status bool
}

func create(w http.ResponseWriter,req *http.Request){
  t,_ := template.ParseFiles("index.html")
  t.Execute(w,nil)
}

func main(){
  http.HandleFunc("/",create)
  http.ListenAndServe(":8090",nil)
}
