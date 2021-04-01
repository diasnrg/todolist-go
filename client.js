const url = 'http://localhost:8090'
document.querySelector('#btgetlist').addEventListener('click',getlist)
document.querySelector('#btadd').addEventListener('click',addTodo)

async function addTodo(){
  const input = document.querySelector('#description')
  const description = input.value
  const response = await fetch(url+'/add/',{
    method:'POST',
    headers:{
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({'description':description})
  })
  input.value = ''
  getlist()
}

async function getlist(){
  const response = await fetch(url+'/list/')
  const data = await response.json()
  document.querySelector('#list').innerHTML = objectToString(data)
}

function objectToString(data){
  result = []
  Array.from(data).forEach(obj => {
    result.push(obj.Description)
  })
  return result
}
