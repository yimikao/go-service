package comment

import "tutgo/business/user"

// Comment represents an individual comment.
type Comment struct {
	ID      int64      `json:"id"`
	Comment string     `json:"comment"`
	UserID  int64      `json:"user_id"`
	User    *user.User `pg:"rel:has-one" json:"user"`
}

// NewComment contains information needed to create a new Comment.
type NewComment struct {
	Comment string `json:"comment"`
	UserID  int64  `json:"user_id"`
}
