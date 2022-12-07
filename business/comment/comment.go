package comment

import "github.com/go-pg/pg/v10"

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(*Comment) (*Comment, error)
	GetById(int64) (*Comment, error)
	All() ([]*Comment, error)
}

type commentLayer struct {
	dbConn *pg.DB
}

func NewCommentLayer(dbConn *pg.DB) commentLayer {
	return commentLayer{
		dbConn: dbConn,
	}
}

func (l commentLayer) All() ([]*Comment, error) {
	comments := make([]*Comment, 0)
	err := l.dbConn.Model(&comments).Relation("User").Select()
	return comments, err
}

func (l commentLayer) Create(cm *Comment) (*Comment, error) {
	_, err := l.dbConn.Model(cm).Insert()
	if err != nil {
		return nil, err
	}

	err = l.dbConn.Model(cm).Relation("User").Where("comment.id = ?", cm.ID).Select()
	return cm, err

}

func (l commentLayer) GetById(id int64) (*Comment, error) {
	var cm = new(Comment)
	err := l.dbConn.Model(&cm).Relation("User").Where("comment.id = ?", id).Select()

	return cm, err
}
