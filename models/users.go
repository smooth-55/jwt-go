package models

type User struct {
	Base
	FirebaseUID string `json:"firebase_uid"`
	Username    string `json:"username" validate:"required"`
	Role        string `json:"role" `
	Email       string `json:"email" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	FullName    string `json:"full_name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Password    string `json:"-" validate:"required"`
}

// TableName gives table name of model
func (m *User) TableName() string {
	return "user"
}

type UserToUpdate struct {
	Email    string
	Username string
	Phone    string
	Role     string
	FullName string
	Address  string
}

// ToMap convert User to map
func (m User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"email":     m.Email,
		"username":  m.Username,
		"role":      m.Role,
		"phone":     m.Phone,
		"full_name": m.FullName,
		"address":   m.Address,
	}
}

// // Runs before inserting a row into table
// func (m *User) BeforeCreate(db *gorm.DB) error {
// 	id, err := uuid.NewRandom()

// 	m.ID = BINARY16(id)
// 	fmt.Println(m.ID, "--id----")
// 	b, err := json.MarshalIndent(m, "", "")
// 	if err == nil {
// 		fmt.Println(string(b))
// 	}
// 	return err
// }
