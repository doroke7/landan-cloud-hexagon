package pkgPostgresql

type PostgresqlOrder struct {
	Field *string `json:"field,omitempty"`
	Value *string `json:"value,omitempty"`
}
