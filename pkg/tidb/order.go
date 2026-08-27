package pkgTidb

type TidbOrder struct {
	Field *string `json:"field,omitempty"`
	Value *string `json:"value,omitempty"`
}
