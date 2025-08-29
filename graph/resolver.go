package graph

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ishidatakuo/graphql-ent-atlas/ent"
)

// Resolver is the root resolver
type Resolver struct {
	Client *ent.Client
}

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

// Query resolver methods
func (r *Resolver) Todo(ctx context.Context, id string) (*ent.Todo, error) {
	todoID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid todo ID: %w", err)
	}
	return r.Client.Todo.Get(ctx, todoID)
}

func (r *Resolver) Todos(ctx context.Context) ([]*ent.Todo, error) {
	return r.Client.Todo.Query().All(ctx)
}

// Mutation resolver methods
func (r *Resolver) CreateTodo(ctx context.Context, input CreateTodoInput) (*ent.Todo, error) {
	builder := r.Client.Todo.Create().
		SetTitle(input.Title)

	if input.Description != nil {
		builder = builder.SetDescription(*input.Description)
	}

	return builder.Save(ctx)
}

func (r *Resolver) UpdateTodo(ctx context.Context, id string, input UpdateTodoInput) (*ent.Todo, error) {
	todoID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid todo ID: %w", err)
	}

	builder := r.Client.Todo.UpdateOneID(todoID)

	if input.Title != nil {
		builder = builder.SetTitle(*input.Title)
	}
	if input.Description != nil {
		builder = builder.SetDescription(*input.Description)
	}
	if input.Completed != nil {
		builder = builder.SetCompleted(*input.Completed)
	}

	return builder.Save(ctx)
}

func (r *Resolver) DeleteTodo(ctx context.Context, id string) (bool, error) {
	todoID, err := strconv.Atoi(id)
	if err != nil {
		return false, fmt.Errorf("invalid todo ID: %w", err)
	}

	err = r.Client.Todo.DeleteOneID(todoID).Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}