package main

import (
	"food-siam-si/food-siam-si-menu/controllers"
	resturant "food-siam-si/food-siam-si-menu/internal"
	"food-siam-si/food-siam-si-menu/models"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	models.ConnectDataBase()

	resturant.Init()

	r := gin.Default()

	menuGroup := r.Group("/menus")

	menuGroup.GET("/:id", controllers.ViewMenu)
	menuGroup.POST("/:id", controllers.AddMenu)
	menuGroup.PUT("/:id", controllers.UpdateMenu)
	menuGroup.DELETE("/:id", controllers.DeleteMenu)
	menuGroup.GET("/:id/random", controllers.RandomMenu)
	menuGroup.GET("/:id/recommend", controllers.ViewRecommendMenu)
	menuGroup.PUT("/:id/recommend", controllers.UpdateRecommendMenu)
	menuGroup.GET("/types", controllers.GetTypes)
	menuGroup.GET("/:id/types", controllers.GetTypesByResturant)

	r.Run(":" + os.Getenv("HOST_PORT"))

}
