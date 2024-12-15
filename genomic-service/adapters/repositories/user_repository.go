package repositories

import (
	"errors"
	"sync"

	"github.com/minhngoc274/genomic-system/genomic-service/models"
)

// UserRepository is responsible for storing and retrieving User data.
type UserRepository struct {
	users     map[uint64]models.User
	addresses map[string]uint64
	mutex     sync.Mutex
	idCounter uint64
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:     make(map[uint64]models.User),
		addresses: make(map[string]uint64),
		idCounter: 0,
	}
}

// Save saves a user to the repository
func (r *UserRepository) Save(user models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.addresses[user.Address]; exists {
		return errors.New("user already exists")
	}

	r.idCounter++
	user.ID = r.idCounter

	r.users[user.ID] = user
	r.addresses[user.Address] = user.ID

	return nil
}

// RetrieveByAddress retrieves a user by their address
func (r *UserRepository) RetrieveByAddress(address string) (models.User, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	userID, exists := r.addresses[address]
	if !exists {
		return models.User{}, false
	}

	user, found := r.users[userID]
	return user, found
}
