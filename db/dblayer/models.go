package dblayer

import (
	"tutgo/db/models"

	"github.com/go-pg/pg/v10"
)

type CommentRepository interface {
	GetById(id int64) (*models.Comment, error)
	All() ([]*models.Comment, error)
	Create(cm *models.Comment) (*models.Comment, error)
}

type UserRepository interface {
	GetById(id int64) (*models.User, error)
	All() ([]*models.User, error)
	Create(cm *models.User) (*models.User, error)
}

type commentLayer struct {
	dbConn *pg.DB
}

func NewCommentLayer(dbConn *pg.DB) CommentRepository {
	return commentLayer{
		dbConn: dbConn,
	}
}

type userLayer struct {
	dbConn *pg.DB
}

func (l commentLayer) All() ([]*models.Comment, error) {
	comments := make([]*models.Comment, 0)
	err := l.dbConn.Model(&comments).Relation("User").Select()
	return comments, err
}

func (l commentLayer) Create(cm *models.Comment) (*models.Comment, error) {
	_, err := l.dbConn.Model(cm).Insert()
	if err != nil {
		return nil, err
	}

	err = l.dbConn.Model(cm).Relation("User").Where("comment.id = ?", cm.ID).Select()
	return cm, err

}

func (l commentLayer) GetById(id int64) (*models.Comment, error) {
	var cm = new(models.Comment)
	err := l.dbConn.Model(&cm).Relation("User").Where("comment.id = ?", id).Select()

	return cm, err
}
