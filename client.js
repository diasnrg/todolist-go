const url = 'http://localhost:8090'
document.querySelector('.btgetlist').addEventListener('click',getlist)

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
