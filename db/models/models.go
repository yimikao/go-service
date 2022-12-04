package models

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	ID      int64  `json:"id"`
	Comment string `json:"comment"`
	UserID  int64  `json:"user_id"`
	User    *User  `pg:"rel:has-one" json:"user"`
}
