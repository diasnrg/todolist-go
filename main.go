package main

import(
  "net/http"
  "html/template"
  "fmt"
  "log"
  "io/ioutil"
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

  //saving object to txt file
  err := ioutil.WriteFile("output.txt", []byte(fmt.Sprintf("%v",newTodo)), 0644)
  check(err)

  //redirecting to the main page
  http.Redirect(w,r,"/",http.StatusFound)

  //method 1
    // f, err := os.Create("data.txt")
    // check(err)
    // defer f.Close()
    // _,err = f.WriteString(fmt.Sprintf("%v",newTodo))
    // check(err)
}

func check(err error){
  if err != nil{
    log.Fatal(err)
  }
}

func list(w http.ResponseWriter,r *http.Request){
  body, err := ioutil.ReadFile("output.txt")
  check(err)
  fmt.Println(string(body))

  //here another template will iterate through items and show them in page

}

func main(){
  http.HandleFunc("/",list)
  http.HandleFunc("/add/",add)
  http.HandleFunc("/save/",save)
  http.ListenAndServe(":8090",nil)
}
