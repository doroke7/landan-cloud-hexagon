package pkgInput

type Filter struct {
	Field    *string `json:"field,omitempty"`
	Operator *string `json:"operator,omitempty"`
	Value    any     `json:"value,omitempty"` // Value 如果解析出來數字 預設是 float64
}
