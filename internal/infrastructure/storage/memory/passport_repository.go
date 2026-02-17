package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/ZakirAlekperov/GoTechPasport/internal/domain/entity"
)

// InMemoryPassportRepository реализация PassportRepository в памяти
// Используется для тестирования и прототипирования
type InMemoryPassportRepository struct {
	mu        sync.RWMutex
	passports map[string]*entity.TechnicalPassport
}

// NewInMemoryPassportRepository создает новый in-memory репозиторий
func NewInMemoryPassportRepository() *InMemoryPassportRepository {
	return &InMemoryPassportRepository{
		passports: make(map[string]*entity.TechnicalPassport),
	}
}

// Create сохраняет новый паспорт
func (r *InMemoryPassportRepository) Create(ctx context.Context, passport *entity.TechnicalPassport) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.passports[passport.ID]; exists {
		return fmt.Errorf("passport with ID %s already exists", passport.ID)
	}

	r.passports[passport.ID] = passport
	return nil
}

// GetByID возвращает паспорт по ID
func (r *InMemoryPassportRepository) GetByID(ctx context.Context, id string) (*entity.TechnicalPassport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	passport, exists := r.passports[id]
	if !exists {
		return nil, fmt.Errorf("passport with ID %s not found", id)
	}

	return passport, nil
}

// Update обновляет существующий паспорт
func (r *InMemoryPassportRepository) Update(ctx context.Context, passport *entity.TechnicalPassport) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.passports[passport.ID]; !exists {
		return fmt.Errorf("passport with ID %s not found", passport.ID)
	}

	r.passports[passport.ID] = passport
	return nil
}

// Delete удаляет паспорт
func (r *InMemoryPassportRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.passports[id]; !exists {
		return fmt.Errorf("passport with ID %s not found", id)
	}

	delete(r.passports, id)
	return nil
}

// List возвращает все паспорта
func (r *InMemoryPassportRepository) List(ctx context.Context) ([]*entity.TechnicalPassport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*entity.TechnicalPassport, 0, len(r.passports))
	for _, passport := range r.passports {
		result = append(result, passport)
	}

	return result, nil
}

// FindByAddress ищет паспорта по адресу
func (r *InMemoryPassportRepository) FindByAddress(ctx context.Context, address entity.Address) ([]*entity.TechnicalPassport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*entity.TechnicalPassport, 0)
	for _, passport := range r.passports {
		if addressMatches(passport.Address, address) {
			result = append(result, passport)
		}
	}

	return result, nil
}

// addressMatches проверяет совпадение адресов
func addressMatches(a, b entity.Address) bool {
	return a.Subject == b.Subject &&
		a.District == b.District &&
		a.City == b.City &&
		a.Street == b.Street &&
		a.House == b.House &&
		a.Building == b.Building
}
