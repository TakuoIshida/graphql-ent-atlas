package resolver

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ishidatakuo/graphql-ent-atlas/ent"
	"github.com/ishidatakuo/graphql-ent-atlas/internal/graph/model"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.CreateTodoInput) (*ent.Todo, error) {
	builder := r.Client.Todo.Create().
		SetTitle(input.Title)

	if input.Description != nil {
		builder = builder.SetDescription(*input.Description)
	}

	return builder.Save(ctx)
}

// UpdateTodo is the resolver for the updateTodo field.
func (r *mutationResolver) UpdateTodo(ctx context.Context, id string, input model.UpdateTodoInput) (*ent.Todo, error) {
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

// DeleteTodo is the resolver for the deleteTodo field.
func (r *mutationResolver) DeleteTodo(ctx context.Context, id string) (bool, error) {
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

// Todo is the resolver for the todo field.
func (r *queryResolver) Todo(ctx context.Context, id string) (*ent.Todo, error) {
	todoID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid todo ID: %w", err)
	}

	return r.Client.Todo.Get(ctx, todoID)
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*ent.Todo, error) {
	return r.Client.Todo.Query().All(ctx)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }