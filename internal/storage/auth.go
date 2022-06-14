package storage

// Auth Storage
type AuthStorage struct {

}

// Auth Storage constructor
func NewAuthStorage() *AuthStorage {
	return &AuthStorage{}
}

// Create user
func (s *AuthStorage) Create() error {
	return nil
}
