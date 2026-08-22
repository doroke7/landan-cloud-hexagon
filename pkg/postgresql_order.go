package pkg

type PostgresqlOrder struct {
	Field *string `json:"field,omitempty"`
	Value *string `json:"value,omitempty"`
}
