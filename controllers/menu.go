package controllers

import (
	"errors"
	resturant "food-siam-si/food-siam-si-menu/internal"
	"food-siam-si/food-siam-si-menu/internal/handlers/proto"
	"food-siam-si/food-siam-si-menu/models"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func strToInt(str string) (uint, error) {
	num, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(num), nil
}

func ViewMenu(c *gin.Context) {
	RestId := c.Param("id")

	id, err := strToInt(RestId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Received id is not int"})
		return
	}

	menus, err := models.GetMenusByResturantId(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"menus": menus})

}

type AddMenuInput struct {
	UserId      uint     `json:"user_id" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Price       uint     `json:"price" binding:"required"`
	IsRecom     bool     `json:"is_recom"`
	ImageUrl    string   `json:"image_url"`
	Addons      []string `json:"addons"`
	Types       []uint32 `json:"typesId"`
}

func AddMenu(c *gin.Context) {
	RestId := c.Param("id")

	id, err := strToInt(RestId)

	var input AddMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := resturant.RestaurantClient.VerifyIdentity(c, &proto.VerifyRestaurantIdentityRequest{
		Id: uint32(id),
		User: &proto.User{
			Id: uint32(input.UserId),
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	if res.Value == false {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
		return
	}

	m := models.Menu{}
	m.Description = input.Description
	m.Price = input.Price
	m.Title = input.Title
	m.RestId = uint(id)
	m.IsRecom = input.IsRecom
	m.ImageUrl = input.ImageUrl

	var types []models.MenuType
	for _, t := range input.Types {
		types = append(types, models.MenuType{
			Id: t,
		})
	}

	m.Types = types

	_, err = m.AddMenu()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	var addons []models.MenuAddons
	for _, addon := range input.Addons {
		addons = append(addons, models.MenuAddons{
			MenuId: m.Id,
			Addons: addon,
		})
	}

	if err := models.DB.Create(&addons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Add Menu Complete"})
}

type UpdateMenuInput struct {
	AddMenuInput
	MenuId uint `json:"menu_id" binding:"required"`
}

func UpdateMenu(c *gin.Context) {
	RestId := c.Param("id")

	id, err := strToInt(RestId)

	var input UpdateMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := resturant.RestaurantClient.VerifyIdentity(c, &proto.VerifyRestaurantIdentityRequest{
		Id: uint32(id),
		User: &proto.User{
			Id: uint32(input.UserId),
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	if res.Value == false {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
		return
	}

	m, err := models.GetMenuByID(input.MenuId)

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	if m.RestId != uint(id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
		return
	}

	m.Description = input.Description
	m.Price = input.Price
	m.Title = input.Title
	m.RestId = uint(id)
	m.IsRecom = input.IsRecom
	m.ImageUrl = input.ImageUrl

	var types []models.MenuType
	for _, t := range input.Types {
		types = append(types, models.MenuType{
			Id: t,
		})
	}

	m.Types = types

	_, err = m.UpdateMenu()

	var addons []models.MenuAddons
	for _, addon := range input.Addons {
		addons = append(addons, models.MenuAddons{
			MenuId: m.Id,
			Addons: addon,
		})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	if err = models.DB.Where("menu_id = ?", input.MenuId).Delete(models.MenuAddons{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	if err = models.DB.Create(&addons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update Menu Complete"})
}

type DeleteMenuInput struct {
	UserId uint `json:"user_id" binding:"required"`
	MenuId uint `json:"menu_id" binding:"required"`
}

func DeleteMenu(c *gin.Context) {

	RestId := c.Param("id")

	id, err := strToInt(RestId)

	var input DeleteMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := resturant.RestaurantClient.VerifyIdentity(c, &proto.VerifyRestaurantIdentityRequest{
		Id: uint32(id),
		User: &proto.User{
			Id: uint32(input.UserId),
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	if res.Value == false {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
		return
	}

	m, err := models.GetMenuByID(input.MenuId)

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	if err = models.DB.Where("menu_id = ?", input.MenuId).Delete(models.MenuAddons{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	if err := models.DB.Model(&m).Association("Types").Replace([]models.MenuAddons{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	err = models.DB.Delete(&m).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delete Menu Complete"})
}

func RandomMenu(c *gin.Context) {
	RestId := c.Param("id")

	id, err := strToInt(RestId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Received id is not int"})
		return
	}

	menus, err := models.GetMenusByResturantId(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	randomIndex := rand.Intn(len(menus))

	c.JSON(http.StatusOK, gin.H{"menus": menus[randomIndex]})
}

func ViewRecommendMenu(c *gin.Context) {
	RestId := c.Param("id")

	id, err := strToInt(RestId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Received id is not int"})
		return
	}

	menus, err := models.GetRecommendMenusByResturantId(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"menus": menus})
}

type UpdateRecommendMenuInput struct {
	UserId  uint `json:"user_id" binding:"required"`
	MenuId  uint `json:"menu_id" binding:"required"`
	IsRecom bool `json:"is_recom"`
}

func UpdateRecommendMenu(c *gin.Context) {
	RestId := c.Param("id")

	id, err := strToInt(RestId)

	var input UpdateRecommendMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := resturant.RestaurantClient.VerifyIdentity(c, &proto.VerifyRestaurantIdentityRequest{
		Id: uint32(id),
		User: &proto.User{
			Id: uint32(input.UserId),
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	if res.Value == false {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
		return
	}

	m := models.Menu{}
	m.Id = input.MenuId

	err = models.DB.Model(&m).Update("is_recom", input.IsRecom).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mark Recommend Menu Complete"})
}

func GetTypes(c *gin.Context) {

	types, err := models.GetAllTypes()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"types": types})
}

func GetTypesByResturant(c *gin.Context) {

	RestId := c.Param("id")

	id, err := strToInt(RestId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Received id is not int"})
		return
	}

	types, err := models.GetTypesByResturant(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem Occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"types": types})
}
