//go:generate mockgen -source storage.go -destination mock/psql_storage_mock.go -package mock
package psql

import (
	"context"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Auth StoragePsql interface
type AuthPsql interface {
	Register(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	FindUsersByName(ctx context.Context, name string, pq *utils.PaginationQuery) (*entity.UsersList, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*entity.UsersList, error)
	FindUserByEmail(ctx context.Context, user *entity.User) (*entity.User, error)
}

// News StoragePsql interface
type NewsPsql interface {
	Create(ctx context.Context, news *entity.News) (*entity.News, error)
	Update(ctx context.Context, news *entity.News) (*entity.News, error)
	GetNews(ctx context.Context, pq *utils.PaginationQuery) (*entity.NewsList, error)
	GetNewsByID(ctx context.Context, newsID uuid.UUID) (*entity.NewsBase, error)
	SearchNews(ctx context.Context, title string, pq *utils.PaginationQuery) (*entity.NewsList, error)
	Delete(ctx context.Context, newsID uuid.UUID) error
}

// Comments storage interface
type CommentsPsql interface {
	Create(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
	Update(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
	GetByID(ctx context.Context, commentID uuid.UUID) (*entity.CommentBase, error)
	GetAllByNewsID(ctx context.Context, newsID uuid.UUID, pq *utils.PaginationQuery) (*entity.CommentsList, error)
	Delete(ctx context.Context, commentID uuid.UUID) error
}

type Storage struct {
	Auth     *AuthStorage
	News     *NewsStorage
	Comments *CommentsStorage
}

func NewStorage(psql *sqlx.DB) *Storage {
	return &Storage{
		Auth:     NewAuthStorage(psql),
		News:     NewNewsStorage(psql),
		Comments: NewCommentsStorage(psql),
	}
}
