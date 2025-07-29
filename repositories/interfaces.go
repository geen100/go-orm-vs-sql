package repositories

import "go-orm-vs-sql/models"

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	EmailExists(email string) (bool, error)
	UsernameExists(username string) (bool, error)
}
