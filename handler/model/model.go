package model

type RegisterUserRequest struct {
	Email    string
	Password string
	Name     string
}

type UpdateUserRequest struct {
	Email string
	Name  string
}

type CreateTaskRequest struct {
	UserID      int
	Title       string
	Description string
}

type UpdateTaskRequest struct {
	Title       string
	Description string
	Status      string
	UserID      int
}
