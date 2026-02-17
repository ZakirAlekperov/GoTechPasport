package passport

import (
	"context"
	"fmt"

	"github.com/ZakirAlekperov/GoTechPasport/internal/domain/entity"
	"github.com/ZakirAlekperov/GoTechPasport/internal/domain/repository"
)

// RemoveBuildingInput входные данные для удаления здания
type RemoveBuildingInput struct {
	PassportID    string
	BuildingIndex int // индекс здания в массиве
}

// RemoveBuildingOutput результат удаления здания
type RemoveBuildingOutput struct {
	Passport *entity.TechnicalPassport
}

// RemoveBuildingUseCase use case для удаления здания из паспорта
type RemoveBuildingUseCase struct {
	repo repository.PassportRepository
}

// NewRemoveBuildingUseCase создает новый use case
func NewRemoveBuildingUseCase(repo repository.PassportRepository) *RemoveBuildingUseCase {
	return &RemoveBuildingUseCase{
		repo: repo,
	}
}

// Execute выполняет удаление здания
func (uc *RemoveBuildingUseCase) Execute(ctx context.Context, input RemoveBuildingInput) (*RemoveBuildingOutput, error) {
	// Валидация входных данных
	if input.PassportID == "" {
		return nil, entity.ValidationError{Field: "passport_id", Message: "ID паспорта обязателен"}
	}

	if input.BuildingIndex < 0 {
		return nil, entity.ValidationError{Field: "building_index", Message: "индекс здания должен быть >= 0"}
	}

	// Получаем паспорт
	passport, err := uc.repo.GetByID(ctx, input.PassportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get passport: %w", err)
	}

	// Проверяем что индекс в пределах массива
	if input.BuildingIndex >= len(passport.Buildings) {
		return nil, entity.ValidationError{
			Field:   "building_index",
			Message: "здание с таким индексом не найдено",
		}
	}

	// Удаляем здание
	passport.Buildings = append(
		passport.Buildings[:input.BuildingIndex],
		passport.Buildings[input.BuildingIndex+1:]...,
	)

	passport.AddAuditEntry("remove_building", fmt.Sprintf("Удалено здание с индексом %d", input.BuildingIndex))

	// Сохраняем изменения
	if err := uc.repo.Update(ctx, passport); err != nil {
		return nil, fmt.Errorf("failed to update passport: %w", err)
	}

	return &RemoveBuildingOutput{
		Passport: passport,
	}, nil
}
