package models

// Todo -> DB model
type Todo struct {
	BaseModel
	Task        string `gorm:"column:Task" json:"task"`
	IsCompleted *bool  `gorm:"column:IsCompleted" json:"is_completed"`
}

// TableName  -> returns table name of model
func (c Todo) TableName() string {
	return "Todo"
}
