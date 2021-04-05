const url = 'http://localhost:8090'
const inputDescription = document.querySelector('#description')
const list = document.querySelector('#list')
document.querySelector('#btadd').addEventListener('click',addItemToDB)
getlist()

async function getlist(){
  list.innerHTML = ''
  const response = await fetch(url+'/list/').then(response => response.json())
  const data = new Map(Object.entries(response))
  console.log(data)
  for(let [id,todo] of data.entries()){
    addItemToDOM(id,todo)
  }
}

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
    getlist()
}

async function deleteItemDB(item){
  const id = item.target.parentElement.parentElement.id
  const data = await fetch(url+'/delete/'+id,{
    method:'DELETE',
    headers:{
      'Content-Type':'application/json'
    }
  })
  getlist()
}

async function updateItemDB(item){
  const id = item.target.parentElement.parentElement.id
  const data = await fetch(url+'/update/'+id,{
    method:'POST',
    headers:{
      'Content-Type':'application/json'
    }
  })

}

function updateItemDOM(id, status){
  //change the uptade button and make/remove line-through
  // const item = document.querySelector('#${id}')
  // const description = item.childNotes[1]
  // console.log(description)
  // if(status){
  //   description.style.textDecoration = 'line-through'
  //   btupdate.classList.add('btn-warning')
  // }else{
  //   description.style.textDecoration = 'none'
  //   btupdate.classList.add('btn-success')
  // }
}

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
  //don't work yet
  updateItemDOM(id, status)
}
