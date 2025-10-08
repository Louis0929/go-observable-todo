package models

import "gorm.io/gorm"

// TODO: Import "gorm.io/gorm" package

// TODO: Define a Todo struct with the following fields:
// 1. Embed gorm.Model (this gives you ID, CreatedAt, UpdatedAt, DeletedAt for free)
// 2. Title field:
//    - Type: string
//    - JSON tag: "title"
//    - Binding tag: "required" (for validation)
// 3. Status field:
//    - Type: string
//    - JSON tag: "status"
//    - GORM tag: "default:pending"
//
// Remember: struct tags use backticks and look like: `json:"fieldname"`
// Multiple tags are space-separated: `json:"title" binding:"required"`

type Todo struct {
	gorm.Model
	Title  string `json:"title" binding:"required"`
	Status string `json:"status" gorm:"default:pending"`
}
