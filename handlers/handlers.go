package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/disharjayanth/nginx-recipes/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewRecipeHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *RecipesHandler {
	return &RecipesHandler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

func (handler *RecipesHandler) ListRecipeHandler(c *gin.Context) {
	val, err := handler.redisClient.Get("recipes").Result()
	if err == redis.Nil {
		log.Println("Request sent to mongoDB")
		cur, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer cur.Close(handler.ctx)

		recipes := make([]models.Recipe, 0)
		for cur.Next(handler.ctx) {
			var recipe models.Recipe
			cur.Decode(&recipe)
			recipes = append(recipes, recipe)
		}

		data, err := json.Marshal(recipes)
		if err != nil {
			fmt.Println("Error marshalling to JSON:", err)
			return
		}

		handler.redisClient.Set("recipes", string(data), 0)

		c.JSON(http.StatusOK, recipes)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		log.Println("Request sent to redis")
		recipes := make([]models.Recipe, 0)
		json.Unmarshal([]byte(val), &recipes)

		c.JSON(http.StatusOK, recipes)
	}
}

func (handler *RecipesHandler) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := handler.collection.InsertOne(handler.ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("Remove recipes from redis")
	handler.redisClient.Del("recipes")

	c.JSON(http.StatusOK, recipe)
}

func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hexId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Erorr converting string to hexid:", err)
		return
	}

	if _, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": hexId,
	}, bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instruction", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags},
	}}}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	handler.redisClient.Del("recipes")

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been updated!",
	})
}

func (handler *RecipesHandler) GetOneRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	hexFromId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error from converting string id to hex:", err)
		return
	}

	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"_id": hexFromId,
	})

	var recipe models.Recipe

	if err := cur.Decode(&recipe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func (handler *RecipesHandler) DeleteOneRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	hexFromId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error converting string to hex:", err)
		return
	}

	if _, err := handler.collection.DeleteOne(handler.ctx, bson.M{"_id": hexFromId}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	handler.redisClient.Del("recipes")

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe with " + id + " has been deleted!",
	})
}

func (handler *RecipesHandler) SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipe := make([]models.Recipe, 0)

	// cur, err := handler.collection.Find(handler.ctx, bson.M{})
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// for cur.Next(handler.ctx) {
	// 	var recipe models.Recipe
	// 	cur.Decode(&recipe)
	// 	for _, recipeTag := range recipe.Tags {
	// 		if tag == recipeTag {
	// 			listOfRecipe = append(listOfRecipe, recipe)
	// 		}
	// 	}
	// }

	cur, err := handler.collection.Find(handler.ctx, bson.M{
		"tags": tag,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	for cur.Next(handler.ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		listOfRecipe = append(listOfRecipe, recipe)
	}

	if len(listOfRecipe) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No recipes found with tag " + tag,
		})
		return
	}

	c.JSON(http.StatusOK, listOfRecipe)
}
