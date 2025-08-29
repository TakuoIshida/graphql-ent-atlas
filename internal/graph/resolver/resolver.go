package resolver

import "github.com/ishidatakuo/graphql-ent-atlas/ent"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the root resolver
type Resolver struct {
	Client *ent.Client
}