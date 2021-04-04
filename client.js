const url = 'http://localhost:8090'
document.querySelector('#btadd').addEventListener('click',addItemToDB)
const list = document.querySelector('#list')
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
  const input = document.querySelector('#description')
  const description = input.value
  input.value = ''
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

function addItemToDOM(id,todo){
  const item = document.createElement('div')
  const description = document.createElement('div')
  const btdelete = document.createElement('button')

  item.id = id
  description.innerHTML = todo.Description
  btdelete.textContent = 'delete'
  btdelete.addEventListener('click',removeItem)

  item.appendChild(description)
  item.appendChild(btdelete)
  btdelete.className += 'btn btn-danger'
  item.className += 'd-flex justify-content-between w-50 p-2 my-4 border border-warning align-items-center rounded'
  list.appendChild(item)
}

async function removeItem(item){
  const id = item.target.parentElement.id
  const data = await fetch(url+'/delete/'+id,{
    method:'DELETE',
    headers:{
      'Content-Type':'application/json'
    }
  })
  getlist()
}
