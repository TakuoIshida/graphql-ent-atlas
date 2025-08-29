package model

// CreateTodoInput represents the input for creating a todo
type CreateTodoInput struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

// UpdateTodoInput represents the input for updating a todo
type UpdateTodoInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}