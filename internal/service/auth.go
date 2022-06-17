package service

import (
	"github.com/Edbeer/restapi/config"
)

// Auth StoragePsql interface
type AuthPsql interface {
	Create() error
}

// Auth StorageRedis interface
type AuthRedis interface {
	Create() error
}

// Auth service
type AuthService struct {
	config  *config.Config
	storagePsql AuthPsql
	storageRedis AuthRedis
}

// Auth service constructor
func NewAuthService(config *config.Config, storagePsql AuthPsql, storageRedis AuthRedis) *AuthService {
	return &AuthService{config: config,storagePsql: storagePsql, storageRedis: storageRedis}
}

// Create new user
func (a *AuthService) Create() error {
	return nil
}
