package graph

const Schema = `
	schema {
		query: Query
		mutation: Mutation
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

	type Todo {
		id: ID!
		title: String!
		description: String
		completed: Boolean!
		createdAt: Time!
		updatedAt: Time!
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

	scalar Time
`