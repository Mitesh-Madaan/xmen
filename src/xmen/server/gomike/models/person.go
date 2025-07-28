package models

type Person struct {
	ID            string `gorm:"column:id;primaryKey"`
	Name          string `gorm:"column:name;type:varchar(100);not null"`
	Kind          string `gorm:"column:kind;type:varchar(50);not null"`
	Age           int    `gorm:"column:age;type:int"`
	Description   string `gorm:"column:description;type:text"`
	Nationality   string `gorm:"column:nationality;type:varchar(100)"`
	Cloned        bool   `gorm:"column:cloned;default:false"`
	ClonedFromRef string `gorm:"column:cloned_from_ref;type:varchar(100);default:''"`
}
