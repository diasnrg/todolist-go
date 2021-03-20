package main

import(
  "net/http"
  "html/template"
  "fmt"
  "log"
  "io/ioutil"
)

type todo struct{
  description string
  status bool
}

func add(w http.ResponseWriter,r *http.Request){
  t,_ := template.ParseFiles("index.html")
  t.Execute(w,nil)
}

func save(w http.ResponseWriter,r *http.Request){
  description := r.FormValue("desc")
  newTodo := todo{description,false}

  oldBody, err := ioutil.ReadFile("output.txt")
  check(err)

  bodyString := string(oldBody) + fmt.Sprintf("%v\n",newTodo)
  err = ioutil.WriteFile("output.txt", []byte(bodyString), 0644)
  check(err)

  fmt.Printf("todo %v added\n",newTodo)

  //redirecting to the main page
  http.Redirect(w,r,"/list/",http.StatusFound)
}

func list(w http.ResponseWriter,r *http.Request){
  // body, err := ioutil.ReadFile("output.txt")
  // check(err)
  //here another template will iterate through items and show them in page
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
