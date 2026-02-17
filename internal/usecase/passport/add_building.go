package passport

import (
	"context"
	"fmt"

	"github.com/ZakirAlekperov/GoTechPasport/internal/domain/entity"
	"github.com/ZakirAlekperov/GoTechPasport/internal/domain/repository"
)

// AddBuildingInput входные данные для добавления здания
type AddBuildingInput struct {
	PassportID string
	Building   entity.Building
}

// AddBuildingOutput результат добавления здания
type AddBuildingOutput struct {
	Passport *entity.TechnicalPassport
}

// AddBuildingUseCase use case для добавления здания в паспорт
type AddBuildingUseCase struct {
	repo repository.PassportRepository
}

// NewAddBuildingUseCase создает новый use case
func NewAddBuildingUseCase(repo repository.PassportRepository) *AddBuildingUseCase {
	return &AddBuildingUseCase{
		repo: repo,
	}
}

// Execute выполняет добавление здания
func (uc *AddBuildingUseCase) Execute(ctx context.Context, input AddBuildingInput) (*AddBuildingOutput, error) {
	// Валидация входных данных
	if input.PassportID == "" {
		return nil, entity.ValidationError{Field: "passport_id", Message: "ID паспорта обязателен"}
	}

	if err := input.Building.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid building data: %w", err)
	}

	// Получаем паспорт
	passport, err := uc.repo.GetByID(ctx, input.PassportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get passport: %w", err)
	}

	// Добавляем здание
	if err := passport.AddBuilding(input.Building); err != nil {
		return nil, fmt.Errorf("failed to add building: %w", err)
	}

	// Сохраняем изменения
	if err := uc.repo.Update(ctx, passport); err != nil {
		return nil, fmt.Errorf("failed to update passport: %w", err)
	}

	return &AddBuildingOutput{
		Passport: passport,
	}, nil
}
