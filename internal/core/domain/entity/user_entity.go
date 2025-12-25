package entity

type UserEntity struct {
	ID       int64 
	Name     string 
	Email    string
}

type UserEntityWithPassword struct {
	ID       int64 
	Name     string 
	Email    string
	Password string
}

type UpdatePasswordEntity struct {
	CurrentPassword string
	NewPassword     string
	ConfirmPassword string
}

type RegisterUserEntity struct {
	Name     string 
	Email    string
	Role     string
	Password string
	ConfirmPassword string
}