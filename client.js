const url = 'http://localhost:8090'
document.querySelector('.btgetlist').addEventListener('click',getlist)
document.querySelector('.btadd').addEventListener('click',addTodo)

async function addTodo(){
  const description = document.querySelector('.description').value
  const response = await fetch(url+'/save/',{
    method:'POST',
    headers:{
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({'description':description,'status':false})
  })
}

async function getlist(){
  const response = await fetch(url+'/list/')
  const data = await response.json()
  document.querySelector('.list').innerHTML = objectToString(data)
  console.log(data)
}

function objectToString(data){
  result = []
  data.forEach(obj => {
    result.push(obj.Description)
  })
  return result
}
