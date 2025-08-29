package graph

import (
	"context"
	"fmt"
	"strconv"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/ishidatakuo/graphql-ent-atlas/ent"
)

// Resolver is the root resolver
type Resolver struct {
	Client *ent.Client
}

// Todo GraphQL resolver
type TodoResolver struct {
	todo *ent.Todo
}

func (r *TodoResolver) ID() graphql.ID {
	return graphql.ID(strconv.Itoa(r.todo.ID))
}

func (r *TodoResolver) Title() string {
	return r.todo.Title
}

func (r *TodoResolver) Description() *string {
	if r.todo.Description == "" {
		return nil
	}
	return &r.todo.Description
}

func (r *TodoResolver) Completed() bool {
	return r.todo.Completed
}

func (r *TodoResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.todo.CreatedAt}
}

func (r *TodoResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.todo.UpdatedAt}
}

// Input types
type CreateTodoInput struct {
	Title       string
	Description *string
}

type UpdateTodoInput struct {
	Title       *string
	Description *string
	Completed   *bool
}

// Query resolvers
func (r *Resolver) Todo(ctx context.Context, args struct{ ID graphql.ID }) (*TodoResolver, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, fmt.Errorf("invalid todo ID: %w", err)
	}

	todo, err := r.Client.Todo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &TodoResolver{todo: todo}, nil
}

func (r *Resolver) Todos(ctx context.Context) ([]*TodoResolver, error) {
	todos, err := r.Client.Todo.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	var resolvers []*TodoResolver
	for _, todo := range todos {
		resolvers = append(resolvers, &TodoResolver{todo: todo})
	}

	return resolvers, nil
}

// Mutation resolvers
func (r *Resolver) CreateTodo(ctx context.Context, args struct{ Input CreateTodoInput }) (*TodoResolver, error) {
	builder := r.Client.Todo.Create().
		SetTitle(args.Input.Title)

	if args.Input.Description != nil {
		builder = builder.SetDescription(*args.Input.Description)
	}

	todo, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &TodoResolver{todo: todo}, nil
}

func (r *Resolver) UpdateTodo(ctx context.Context, args struct {
	ID    graphql.ID
	Input UpdateTodoInput
}) (*TodoResolver, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, fmt.Errorf("invalid todo ID: %w", err)
	}

	builder := r.Client.Todo.UpdateOneID(id)

	if args.Input.Title != nil {
		builder = builder.SetTitle(*args.Input.Title)
	}
	if args.Input.Description != nil {
		builder = builder.SetDescription(*args.Input.Description)
	}
	if args.Input.Completed != nil {
		builder = builder.SetCompleted(*args.Input.Completed)
	}

	todo, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &TodoResolver{todo: todo}, nil
}

func (r *Resolver) DeleteTodo(ctx context.Context, args struct{ ID graphql.ID }) (bool, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return false, fmt.Errorf("invalid todo ID: %w", err)
	}

	err = r.Client.Todo.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}