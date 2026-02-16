package repository

import (
	"context"

	"github.com/ZakirAlekperov/GoTechPasport/internal/domain/entity"
)

// PassportRepository определяет интерфейс для работы с техническими паспортами
type PassportRepository interface {
	// Create создает новый технический паспорт
	Create(ctx context.Context, passport *entity.TechnicalPassport) error

	// GetByID возвращает паспорт по ID
	GetByID(ctx context.Context, id string) (*entity.TechnicalPassport, error)

	// Update обновляет существующий паспорт
	Update(ctx context.Context, passport *entity.TechnicalPassport) error

	// Delete удаляет паспорт по ID
	Delete(ctx context.Context, id string) error

	// List возвращает список всех паспортов
	List(ctx context.Context) ([]*entity.TechnicalPassport, error)

	// FindByAddress ищет паспорта по адресу
	FindByAddress(ctx context.Context, address entity.Address) ([]*entity.TechnicalPassport, error)
}
