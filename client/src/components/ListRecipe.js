const ListRecipe = ({ recipe }) => {
    return (
        <div className="pb-10">
            <h3 className="font-bold text-2xl md:text-4xl text-center">{recipe.name}</h3>
            <p className="p-2 flex justify-end">Tags:
            {
                recipe.tags.map((tag, index) => (
                    <div className="px-0.5" key={index}>{tag}</div>
                ))
            }
            </p>
            <h4 className="font-bold text-xl pb-1 lg:pb-4">Ingredients:</h4>
            <ul className="list-decimal my-1">
                {
                    recipe.ingredients.map((ingredient, index) => (
                        <li key={index}>{ingredient}</li>
                    ))
                }
            </ul>
            <div className="my-2">
            <h4 className="font-bold text-xl pb-1 lg:pb-4">Instructions:</h4>
            <ul className="list-disc">
                {
                    recipe.instructions.map((instruction, index) => (
                        <li key={index}>{instruction}</li>
                    ))
                }
            </ul>
            </div>
        </div>
    )
}

export default ListRecipe