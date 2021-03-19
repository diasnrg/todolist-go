package main

import(
  "net/http"
  "html/template"
  "fmt"
  "log"
  "io/ioutil"
  "os"
)

type todo struct{
  description string
  status bool
}

var file *os.File
var err error

func init(){
  file, err = os.Create("output.txt")
  check(err)
}

func add(w http.ResponseWriter,r *http.Request){
  t,_ := template.ParseFiles("index.html")
  t.Execute(w,nil)
}

func save(w http.ResponseWriter,r *http.Request){
  description := r.FormValue("desc")
  newTodo := todo{description,false}

  //saving object to txt file

  //method 1
    // file, err := os.Create("output.txt")
    // check(err)
    file, err = os.Open("output.txt")
    defer file.Close()
    _,err = file.WriteString(fmt.Sprintf("%v",newTodo))
    check(err)

  // err := ioutil.WriteFile("output.txt", []byte(fmt.Sprintf("%v",newTodo)), 0644)
  // check(err)

  fmt.Printf("todo %v added\n",newTodo)

  //redirecting to the main page
  http.Redirect(w,r,"/list/",http.StatusFound)
}

func list(w http.ResponseWriter,r *http.Request){
  body, err := ioutil.ReadFile("output.txt")
  check(err)
  fmt.Println(string(body))

  //here another template will iterate through items and show them in page

}

func check(err error){
  if err != nil{
    log.Fatal(err)
  }
}

func (t todo) String() string{
  return fmt.Sprintf("%v:%v\n",t.description,t.status)
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
