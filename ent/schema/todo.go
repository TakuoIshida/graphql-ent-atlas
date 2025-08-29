package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").
			NotEmpty().
			Comment("Todo item title"),
		field.String("description").
			Optional().
			Comment("Todo item description"),
		field.Bool("completed").
			Default(false).
			Comment("Whether the todo is completed"),
		field.Time("createdAt").
			Default(time.Now).
			Comment("When the todo was created"),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("When the todo was last updated"),
	}
}

// Edges of the Todo.
func (Todo) Edges() []ent.Edge {
	return nil
}
