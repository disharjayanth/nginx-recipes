import { useEffect, useState } from 'react'
import axios from 'axios'
import ListRecipe from './components/ListRecipe'

function App() {
  const [recipes, setRecipes] = useState([])

  useEffect(() => {
    async function fetchRecipes() {
      const res = await axios.get('/api/recipes')
      const data = await res.data
      setRecipes(data)
    }

    fetchRecipes()
  }, [])

  return (
    <div className="App">
      {recipes.map((recipe, index) => (
        <ListRecipe key={index} recipe={recipe} />
      ))}
    </div>
  )
}

export default App
