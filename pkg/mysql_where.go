package pkg

type MysqlWhere struct {
	Field    *string `json:"field,omitempty"`
	Operator *string `json:"operator,omitempty"`
	Value    any     `json:"value,omitempty"`
}
