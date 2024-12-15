package repositories

import (
	"errors"
	"sync"

	"github.com/minhngoc274/genomic-system/genomic-service/models"
)

// GeneticDataRepository is responsible for storing and retrieving Genetic data.
type GeneticDataRepository struct {
	data  map[string]models.GeneticData
	mutex sync.Mutex
}

// NewGeneticDataRepository creates a new instance of GeneticDataRepository
func NewGeneticDataRepository() *GeneticDataRepository {
	return &GeneticDataRepository{
		data: make(map[string]models.GeneticData),
	}
}

// Create stores or updates GeneticData in the repository
func (r *GeneticDataRepository) Create(geneticData models.GeneticData) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, exists := r.data[geneticData.FileID]; exists {
		return errors.New("file ID already exists")
	}
	r.data[geneticData.FileID] = geneticData
	return nil
}

// Retrieve retrieves GeneticData by its FileID
func (r *GeneticDataRepository) Retrieve(fileID string) (models.GeneticData, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	data, exists := r.data[fileID]
	return data, exists
}

// Update update data
func (r *GeneticDataRepository) Update(geneticData models.GeneticData) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, exists := r.data[geneticData.FileID]; !exists {
		return errors.New("file ID not exists")
	}
	r.data[geneticData.FileID] = geneticData
	return nil
}
