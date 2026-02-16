package service

import (
	"context"

	"github.com/ZakirAlekperov/GoTechPasport/internal/domain/entity"
)

// ValidationResult результат валидации
type ValidationResult struct {
	// Valid указывает, прошла ли валидация успешно
	Valid bool

	// Errors список ошибок валидации
	Errors []entity.ValidationError

	// Warnings предупреждения (некритичные проблемы)
	Warnings []string
}

// Validator определяет интерфейс для валидации сущностей
type Validator interface {
	// ValidatePassport выполняет полную валидацию технического паспорта
	ValidatePassport(ctx context.Context, passport *entity.TechnicalPassport) ValidationResult

	// ValidateAddress валидирует адрес по правилам ГАР
	ValidateAddress(ctx context.Context, address entity.Address) ValidationResult
}
