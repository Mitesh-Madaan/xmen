package models

type Animal struct {
	ID            string `gorm:"column:id;primaryKey"`
	Name          string `gorm:"column:name;type:varchar(255);not null"`
	Kind          string `gorm:"column:kind;type:varchar(100);not null"`
	Age           int    `gorm:"column:age;type:int;not null"`
	Description   string `gorm:"column:description;type:text"`
	Breed         string `gorm:"column:breed;type:varchar(255)"`
	Cloned        bool   `gorm:"column:cloned;type:boolean;default:false"`
	ClonedFromRef string `gorm:"column:cloned_from_ref;type:varchar(100);default:''"`
}
