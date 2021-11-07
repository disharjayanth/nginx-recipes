package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/disharjayanth/nginx-recipes/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesHandler *handlers.RecipesHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		fmt.Println("Error connecting to mongodb server:", err)
		return
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongoDB server")

	collection := client.Database("nginxRecipes").Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping()
	log.Println(status)

	recipesHandler = handlers.NewRecipeHandler(ctx, collection, redisClient)
}

func main() {
	router := gin.Default()

	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost/api/recipes"},
	// 	AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	router.Use(cors.Default())

	router.GET("/recipes", recipesHandler.ListRecipeHandler)
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipe/:id", recipesHandler.GetOneRecipeHandler)
	router.GET("/recipe", recipesHandler.SearchRecipeHandler)
	router.PUT("/recipe/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipe/:id", recipesHandler.DeleteOneRecipeHandler)

	router.Run(":3000")
}

// // Reading recipes from recipe.json file and inserting them to mongoDB
// ctx := context.Background()
// 	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
// 	if err != nil {
// 		fmt.Println("Error connecting to mongodb server:", err)
// 		return
// 	}

// 	sb, err := ioutil.ReadFile("recipes.json")
// 	if err != nil {
// 		fmt.Println("Error reading recipes.json file:", err)
// 		return
// 	}

// 	if err := json.Unmarshal(sb, &recipes); err != nil {
// 		fmt.Println("Error unmarshalling:", err)
// 		return
// 	}

// 	var listOfRecipes []interface{}
// 	for _, recipe := range recipes {
// 		listOfRecipes = append(listOfRecipes, recipe)
// 	}

// 	res, err := client.Database("nginxRecipes").Collection("recipes").InsertMany(ctx, listOfRecipes)
// 	if err != nil {
// 		fmt.Println("Error inserting many docs:", err)
// 		return
// 	}

// 	fmt.Println("Result of inserting many docs:", len(res.InsertedIDs))
