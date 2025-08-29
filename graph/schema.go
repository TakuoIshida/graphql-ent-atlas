package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/introspection"
	"github.com/ishidatakuo/graphql-ent-atlas/ent"
)

// NewExecutableSchema creates an executable schema from our resolver
func NewExecutableSchema(cfg Config) graphql.ExecutableSchema {
	return &executableSchema{
		resolvers: cfg.Resolvers,
	}
}

// Config is the configuration for the GraphQL schema
type Config struct {
	Resolvers *Resolver
}

type executableSchema struct {
	resolvers *Resolver
}

func (e *executableSchema) Schema() *graphql.Schema {
	return parsedSchema
}

func (e *executableSchema) Complexity(typeName, field string, childComplexity int, rawArgs map[string]interface{}) (int, bool) {
	return 0, false
}

func (e *executableSchema) Exec(ctx context.Context) graphql.ResponseHandler {
	return graphql.NewExecutor(e.Schema()).Exec
}

// Schema definition string
const schemaString = `
type Todo {
  id: ID!
  title: String!
  description: String
  completed: Boolean!
}

input CreateTodoInput {
  title: String!
  description: String
}

input UpdateTodoInput {
  title: String
  description: String
  completed: Boolean
}

type Query {
  todo(id: ID!): Todo
  todos: [Todo!]!
}

type Mutation {
  createTodo(input: CreateTodoInput!): Todo!
  updateTodo(id: ID!, input: UpdateTodoInput!): Todo!
  deleteTodo(id: ID!): Boolean!
}
`

var parsedSchema = graphql.MustParseSchema(schemaString, &schemaResolver{})

type schemaResolver struct{}

func (r *schemaResolver) Todo(ctx context.Context, args struct{ ID string }) (*Todo, error) {
	resolver := getResolver(ctx)
	entTodo, err := resolver.Todo(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return entTodoToGraphQL(entTodo), nil
}

func (r *schemaResolver) Todos(ctx context.Context) ([]*Todo, error) {
	resolver := getResolver(ctx)
	entTodos, err := resolver.Todos(ctx)
	if err != nil {
		return nil, err
	}

	todos := make([]*Todo, len(entTodos))
	for i, entTodo := range entTodos {
		todos[i] = entTodoToGraphQL(entTodo)
	}
	return todos, nil
}

func (r *schemaResolver) CreateTodo(ctx context.Context, args struct{ Input CreateTodoInput }) (*Todo, error) {
	resolver := getResolver(ctx)
	entTodo, err := resolver.CreateTodo(ctx, args.Input)
	if err != nil {
		return nil, err
	}
	return entTodoToGraphQL(entTodo), nil
}

func (r *schemaResolver) UpdateTodo(ctx context.Context, args struct {
	ID    string
	Input UpdateTodoInput
}) (*Todo, error) {
	resolver := getResolver(ctx)
	entTodo, err := resolver.UpdateTodo(ctx, args.ID, args.Input)
	if err != nil {
		return nil, err
	}
	return entTodoToGraphQL(entTodo), nil
}

func (r *schemaResolver) DeleteTodo(ctx context.Context, args struct{ ID string }) (bool, error) {
	resolver := getResolver(ctx)
	return resolver.DeleteTodo(ctx, args.ID)
}

// Todo represents a GraphQL Todo type
type Todo struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Completed   bool    `json:"completed"`
}

func entTodoToGraphQL(entTodo *ent.Todo) *Todo {
	var description *string
	if entTodo.Description != "" {
		description = &entTodo.Description
	}

	return &Todo{
		ID:          strconv.Itoa(entTodo.ID),
		Title:       entTodo.Title,
		Description: description,
		Completed:   entTodo.Completed,
	}
}

func getResolver(ctx context.Context) *Resolver {
	return ctx.Value("resolver").(*Resolver)
}