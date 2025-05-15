package models

type Setting struct {
	ID       uint   `gorm:"primary_key"`
	Email    string `gorm:"unique;not null"`
	Tool     string `gorm:"default:board" json:"tool"`
	Rule     string `gorm:"default:all" json:"rule"`
	Password string `gorm:"size 100" json:"password"`
}
