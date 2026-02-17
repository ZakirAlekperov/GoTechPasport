package passport

import (
	"context"
	"fmt"
	"time"

	"github.com/ZakirAlekperov/GoTechPasport/internal/domain/entity"
	"github.com/ZakirAlekperov/GoTechPasport/internal/domain/repository"
)

// CreatePassportInput входные данные для создания паспорта
type CreatePassportInput struct {
	ObjectType       entity.ObjectType
	Address          entity.Address
	OrganizationName string
	GeneralInfo      entity.GeneralInfo
}

// CreatePassportOutput результат создания паспорта
type CreatePassportOutput struct {
	Passport *entity.TechnicalPassport
}

// CreatePassportUseCase use case для создания технического паспорта
type CreatePassportUseCase struct {
	repo repository.PassportRepository
}

// NewCreatePassportUseCase создает новый use case
func NewCreatePassportUseCase(repo repository.PassportRepository) *CreatePassportUseCase {
	return &CreatePassportUseCase{
		repo: repo,
	}
}

// Execute выполняет создание технического паспорта
func (uc *CreatePassportUseCase) Execute(ctx context.Context, input CreatePassportInput) (*CreatePassportOutput, error) {
	// Валидация входных данных
	if err := input.validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Создаем новый паспорт
	passport := entity.NewTechnicalPassport(input.ObjectType, input.Address)
	passport.OrganizationName = input.OrganizationName
	passport.GeneralInfo = input.GeneralInfo
	passport.AsOfDate = time.Now()

	// Генерируем ID (в реальном приложении можно использовать UUID)
	passport.ID = generateID()

	// Валидируем созданный паспорт
	if err := passport.IsValid(); err != nil {
		return nil, fmt.Errorf("passport validation failed: %w", err)
	}

	// Сохраняем в репозиторий
	if err := uc.repo.Create(ctx, passport); err != nil {
		return nil, fmt.Errorf("failed to save passport: %w", err)
	}

	passport.AddAuditEntry("created", "Технический паспорт создан")

	return &CreatePassportOutput{
		Passport: passport,
	}, nil
}

// validate проверяет корректность входных данных
func (input *CreatePassportInput) validate() error {
	if input.ObjectType == "" {
		return entity.ValidationError{Field: "object_type", Message: "тип объекта обязателен"}
	}

	if err := input.Address.IsValid(); err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	if input.OrganizationName == "" {
		return entity.ValidationError{Field: "organization_name", Message: "наименование организации обязательно"}
	}

	if err := input.GeneralInfo.IsValid(); err != nil {
		return fmt.Errorf("invalid general info: %w", err)
	}

	return nil
}

// generateID генерирует простой ID (в production использовать UUID)
func generateID() string {
	return fmt.Sprintf("TP-%d", time.Now().UnixNano())
}
