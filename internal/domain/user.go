package domain

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserValue struct {
	Name *string `json:"name,omitempty"`
}
