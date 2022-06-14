package storage

// Auth Service interface
type Auth interface {
	Create() error
}

type Storage struct {
	Auth Auth
}

func NewStorage() *Storage {
	return &Storage{
		Auth: NewAuthStorage(),
	}
}