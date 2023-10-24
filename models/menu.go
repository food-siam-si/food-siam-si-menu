package models

import (
	"errors"
)

// title: string
// description: string
// price: number
// isRecommended: boolean
// imageUrl: string
// restaurantId: number
// addons: ???

type Menu struct {
	Id          uint         `gorm:"primary_key;auto_increment" json:"id"`
	Title       string       `gorm:"size:255;not null;" json:"title"`
	Description string       `gorm:"size:255;not null;" json:"description"`
	Price       uint         `gorm:"not null;" json:"price"`
	IsRecom     bool         `gorm:"not null;" json:"is_recom"`
	ImageUrl    string       `gorm:"size:255;" json:"image_url"`
	RestId      uint         `gorm:"not null;" json:"rest_id"`
	Addons      []MenuAddons `gorm:"foreignKey:MenuId;references:Id" json:"addons"`
}

type MenuAddons struct {
	MenuId uint   `gorm:"not null;" json:"menu_id"`
	Addons string `gorm:"size:255;not null;" json:"addons"`
}

func GetMenuByID(uid uint) (Menu, error) {

	var m Menu

	if err := DB.First(&m, uid).Error; err != nil {
		return m, errors.New("Menu not found!")
	}

	return m, nil

}

func GetMenusByResturantId(RestId uint) ([]Menu, error) {

	var m []Menu

	if err := DB.Preload("Addons").Where("rest_id = ?", RestId).Find(&m).Error; err != nil {
		return m, errors.New("Menu not found!")
	}

	return m, nil

}

func (m *Menu) AddMenu() (*Menu, error) {

	err := DB.Create(&m).Error

	if err != nil {
		return &Menu{}, err
	}

	return m, nil
}

func (m *Menu) UpdateMenu() (*Menu, error) {

	err := DB.Save(&m).Error

	if err != nil {
		return &Menu{}, err
	}

	return m, nil
}

func GetRecommendMenusByResturantId(RestId uint) ([]Menu, error) {

	var menus []Menu

	err := DB.Where("rest_id = ? AND is_recom = ?", RestId, true).Find(&menus).Error

	if err != nil {
		return menus, err
	}

	return menus, nil
}
