const url = 'http://localhost:8090'
const list = document.querySelector('#tasks')

//{state=null} - all tasks, {state=true} - completed tasks, {state=false} - incompleted tasks
var state = null

document.querySelector('#btadd').addEventListener('click',addItemDB)
document.querySelector('#btalltasks').addEventListener('click',()=>{    state=null; getItems(state) })
document.querySelector('#btcompleted').addEventListener('click',()=>{   state=true; getItems(state) })
document.querySelector('#btincompleted').addEventListener('click',()=>{ state=false;getItems(state) })
getItems()

//get the tasks from the back end API
async function getItems(state){
  list.innerHTML = ''
  //fetch the back end api to retrieve tasks
  const response = await fetch(url+'/tasks/').then(response => response.json())

  //converting json object to map, because otherwise i can't iterate it lol
  const tasksMap = new Map(Object.entries(response))

  //iterating the tasks(map) and adding new tasks to DOM
  for(let [id,t] of tasksMap.entries()){
    if(state == null || state == t.Status){
      addItemDOM(id,t)
    }
  }
}

//create new task item from input's value and send it to the back end
async function addItemDB(){
  const inputDescription = document.querySelector('#description')
  const description = inputDescription.value

  //check for non-empty input
  if(description == ''){
    console.log('Empty string is not a valid task')
    return
  }

  inputDescription.value = ''

  const data = await fetch(url+'/add/',{
    method:'POST',
    headers:{
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({'description':description})
    })
  .then(response => response.json())
  getItems(state)
}

//delete the task from back end by fetching with it's id
async function deleteItemDB(item){
  const id = item.target.parentElement.parentElement.id
  const data = await fetch(url+'/delete/'+id,{
    method:'DELETE',
    headers:{
      'Content-Type':'application/json'
    }
  })
  getItems(state)
}

//complete/incomplete the task and fetch to back end
async function updateItemDB(item){
  const id = item.target.parentElement.parentElement.id
  const status = await fetch(url+'/update/'+id,{
    method:'POST',
    headers:{
      'Content-Type':'application/json'
    }
  }).then(response => response.json())
  updateItemDOM(id,status)
}

//updates the task's DOM depending on it's status value (true/false)
function updateItemDOM(id, status){
  const item = document.getElementById(id)
  const description = item.children[0]
  const btupdate = item.children[1].children[0]
  if(status){
    description.style.textDecoration = 'line-through'
    btupdate.classList.remove('btn-warning')
    btupdate.classList.add('btn-success')
    btupdate.textContent = '[1]'
  }else{
    description.style.textDecoration = 'none'
    btupdate.classList.remove('btn-success')
    btupdate.classList.add('btn-warning')
    btupdate.textContent = '[0]'
  }
}

//adds new task card(div) to dom, sticks task's id to div's id
function addItemDOM(id,task){
  const item = document.createElement('div')
  const description = document.createElement('div')
  const btns = document.createElement('div')
  const btdelete = document.createElement('button')
  const btupdate = document.createElement('button')

  item.id = id
  item.className += 'd-flex justify-content-between align-items-center p-2 my-4 border rounded'

  description.textContent = task.Description

  btupdate.addEventListener('click',updateItemDB)
  btupdate.classList.add('btn')

  btdelete.textContent = '[x]'
  btdelete.addEventListener('click',deleteItemDB)
  btdelete.className += 'btn btn-danger'

  item.appendChild(description)
  btns.appendChild(btupdate)
  btns.appendChild(btdelete)
  item.appendChild(btns)
  list.appendChild(item)

  updateItemDOM(id, task.Status)
}
