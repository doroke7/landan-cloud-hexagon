package pkgMysql

type MysqlOrder struct {
	Field *string `json:"field,omitempty"`
	Value *string `json:"value,omitempty"`
}
