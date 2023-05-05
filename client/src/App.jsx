import './App.css'

function App() {
  return (
    <div>
        <h1>Hello world</h1>
        <button onClick={async () => {
          const response = await fetch('/users')
          const data = await response.json()
          console.log(data)
        }}>Data</button>
    </div>
  )
}

export default App
