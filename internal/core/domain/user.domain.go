package domain

type User struct {
	BaseModel
	Email string `json:"email" gorm:"type:varchar(255)"`
	Name  string `json:"name" gorm:"type:varchar(255)"`
}

func (t User) TableName() string {
	return "users"
}
