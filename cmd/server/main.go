package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	_ "github.com/lib/pq"

	"github.com/ishidatakuo/graphql-ent-atlas/ent"
	"github.com/ishidatakuo/graphql-ent-atlas/internal/graph"
)

const defaultPort = "8090"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// データベース接続
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:password@localhost:5432/todoapp?sslmode=disable"
	}

	// Entクライアント初期化
	client, err := ent.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	// スキーマ自動マイグレーション（開発環境用）
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// GraphQLスキーマを設定
	schema := graphql.MustParseSchema(graph.Schema, &graph.Resolver{
		Client: client,
	})

	// Chi routerの設定
	r := chi.NewRouter()

	// CORS設定
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8090"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	// ミドルウェア設定
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// GraphQL endpoint
	r.Handle("/graphql", &relay.Handler{Schema: schema})

	// GraphQL Playground
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>GraphQL Playground</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/css/index.css" />
</head>
<body>
    <div id="root"></div>
    <script src="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/js/middleware.js"></script>
    <script>
        GraphQLPlayground.init(document.getElementById('root'), {
            endpoint: '/graphql'
        })
    </script>
</body>
</html>`
		w.Write([]byte(html))
	})

	// REST API エンドポイント（互換性のため）
	r.Get("/todos", getTodos(client))
	r.Post("/todos", createTodo(client))
	r.Get("/todos/{id}", getTodo(client))
	r.Put("/todos/{id}", updateTodo(client))
	r.Delete("/todos/{id}", deleteTodo(client))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// シンプルなヘルスチェック
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("🚀 GraphQL server running on http://localhost:%s/", port)
	log.Printf("📊 GraphQL Playground: http://localhost:%s/", port)
	log.Printf("🔗 GraphQL endpoint: http://localhost:%s/graphql", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// REST API handler functions

// TodoRequest represents the request body for creating/updating todos
type TodoRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

// TodoResponse represents the response for todos
type TodoResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Completed   bool    `json:"completed"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func entTodoToResponse(todo *ent.Todo) TodoResponse {
	var description *string
	if todo.Description != "" {
		description = &todo.Description
	}

	return TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}
}

func getTodos(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := client.Todo.Query().All(r.Context())
		if err != nil {
			http.Error(w, "Failed to get todos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var response []TodoResponse
		for _, todo := range todos {
			response = append(response, entTodoToResponse(todo))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func createTodo(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TodoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if req.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		builder := client.Todo.Create().SetTitle(req.Title)
		if req.Description != nil {
			builder = builder.SetDescription(*req.Description)
		}

		todo, err := builder.Save(r.Context())
		if err != nil {
			http.Error(w, "Failed to create todo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(entTodoToResponse(todo))
	}
}

func getTodo(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid todo ID", http.StatusBadRequest)
			return
		}

		todo, err := client.Todo.Get(r.Context(), id)
		if err != nil {
			if ent.IsNotFound(err) {
				http.Error(w, "Todo not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to get todo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entTodoToResponse(todo))
	}
}

func updateTodo(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid todo ID", http.StatusBadRequest)
			return
		}

		var req TodoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		builder := client.Todo.UpdateOneID(id)
		if req.Title != "" {
			builder = builder.SetTitle(req.Title)
		}
		if req.Description != nil {
			builder = builder.SetDescription(*req.Description)
		}
		if req.Completed != nil {
			builder = builder.SetCompleted(*req.Completed)
		}

		todo, err := builder.Save(r.Context())
		if err != nil {
			if ent.IsNotFound(err) {
				http.Error(w, "Todo not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to update todo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entTodoToResponse(todo))
	}
}

func deleteTodo(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid todo ID", http.StatusBadRequest)
			return
		}

		err = client.Todo.DeleteOneID(id).Exec(r.Context())
		if err != nil {
			if ent.IsNotFound(err) {
				http.Error(w, "Todo not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to delete todo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
