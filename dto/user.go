package dto

type User struct {
	Email    string `json:"email" binding:"required,email,unique"`
	Password string `json:"password" binding:"required,min=8,max=50,password"`
}
