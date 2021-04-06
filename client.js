const url = 'http://localhost:8090'
const inputDescription = document.querySelector('#description')
const list = document.querySelector('#list')
document.querySelector('#btadd').addEventListener('click',addItemToDB)
document.querySelector('#btlist').addEventListener('click',getItems)
document.querySelector('#btcompleted').addEventListener('click',()=>{   getByStatus(true)   })
document.querySelector('#btuncompleted').addEventListener('click',()=>{ getByStatus(false)  })

getItems()
var todoMap = []

//get all of the todos with fetch request to db
async function updateList(){
  const response = await fetch(url+'/list/').then(response => response.json())
  //converting json object to map, 'cause otherwise i can't iterate it lol
  todoMap = new Map(Object.entries(response))
  console.log(todoMap)
}

async function getItems(){
  const response = await updateList()
  list.innerHTML = ''
  for(let [id,todo] of todoMap.entries()){
    addItemToDOM(id,todo)
  }
}

async function getByStatus(status){
  const response = await updateList()
  list.innerHTML = ''
  for(let [id,todo] of todoMap.entries()){
    if(status == todo.Status){
      addItemToDOM(id,todo)
    }
  }
}

//create new todo item from input's value and send it to db
async function addItemToDB(){
  const description = inputDescription.value
  inputDescription.value = ''
  const data = await fetch(url+'/add/',{
    method:'POST',
    headers:{
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({'description':description})
    })
    .then(response => response.json())
  getItems()
}

//delete the todo from db by fetching with it's id
async function deleteItemDB(item){
  const id = item.target.parentElement.parentElement.id
  const data = await fetch(url+'/delete/'+id,{
    method:'DELETE',
    headers:{
      'Content-Type':'application/json'
    }
  })
  getItems()
}

//complete/uncomplete todo and fetch to database
async function updateItemDB(item){
  const id = item.target.parentElement.parentElement.id
  const newStatus = await fetch(url+'/update/'+id,{
    method:'POST',
    headers:{
      'Content-Type':'application/json'
    }
  }).then(response => response.json())
  updateItemDOM(id,newStatus)
}

//updates the todo's DOM depending on it's status (true/false)
function updateItemDOM(id, status){
  const item = document.getElementById(id)
  const description = item.children[0]
  const btupdate = item.children[1].children[0]
  if(status){
    description.style.textDecoration = 'line-through'
    btupdate.classList.remove('btn-success')
    btupdate.classList.add('btn-warning')
  }else{
    description.style.textDecoration = 'none'
    btupdate.classList.remove('btn-warning')
    btupdate.classList.add('btn-success')
  }
}

//adds new todo card(div) to dom, sticks todo's id to div's id
function addItemToDOM(id,todo){
  const item = document.createElement('div')
  const description = document.createElement('div')
  const btns = document.createElement('div')
  const btdelete = document.createElement('button')
  const btupdate = document.createElement('button')

  item.id = id
  item.className += 'd-flex justify-content-between align-items-center w-50 p-2 my-4 border rounded'

  description.textContent = todo.Description

  btupdate.textContent = 'update'
  btupdate.addEventListener('click',updateItemDB)
  btupdate.classList.add('btn')

  btdelete.textContent = 'delete'
  btdelete.addEventListener('click',deleteItemDB)
  btdelete.className += 'btn btn-danger'

  item.appendChild(description)
  btns.appendChild(btupdate)
  btns.appendChild(btdelete)
  item.appendChild(btns)
  list.appendChild(item)

  updateItemDOM(id, todo.Status)
}
