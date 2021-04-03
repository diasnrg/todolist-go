const url = 'http://localhost:8090'
document.querySelector('#btgetlist').addEventListener('click',getlist)
document.querySelector('#btadd').addEventListener('click',addItemToDB)
const list = document.querySelector('#list')

async function getlist(){
  list.innerHTML = ''
  const response = await fetch(url+'/list/')
  const data = await response.json()
  const mp = new Map(Object.entries(data))
  console.log(mp)
  addItems(mp)
}

function addItemToDB(){
  const input = document.querySelector('#description')
  const description = input.value
  input.value = ''
  fetch(url+'/add/',{
    method:'POST',
    headers:{
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({'description':description})
  }).then(getlist())
}

function addItems(data){
  for(let [id,todo] of data.entries()){
    const item = document.createElement('div')
    const description = document.createElement('div')
    const btdelete = document.createElement('button')

    item.id = id
    description.innerHTML = todo.Description
    btdelete.textContent = 'delete'
    btdelete.addEventListener('click',removeItem)

    item.appendChild(description)
    item.appendChild(btdelete)
    list.appendChild(item)
  }
}

async function removeItem(item){
  const id = item.target.parentElement.id
  fetch(url+'/delete/'+id,{
    method:'DELETE',
    headers:{
      'Content-Type':'application/json'
    }
  }).then(getlist())
}
