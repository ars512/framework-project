package models

type Product struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	Name        string   `json:"name" gorm:"not null"`
	Description string   `json:"description" gorm:"not null"`
	Price       float64  `json:"price" gorm:"not null"`
	Stock       int      `json:"stock" gorm:"not null"`
	BrandID     uint     `json:"brand_id" gorm:"not null"`
	Brand       Brand    `json:"brand" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CategoryID  uint     `json:"category_id" gorm:"not null"`
	Category    Category `json:"category" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
