package passport_test

import (
	"context"
	"testing"

	"github.com/ZakirAlekperov/GoTechPasport/internal/domain/entity"
	"github.com/ZakirAlekperov/GoTechPasport/internal/infrastructure/storage/memory"
	"github.com/ZakirAlekperov/GoTechPasport/internal/usecase/passport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePassportUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		input   passport.CreatePassportInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid input creates passport",
			input: passport.CreatePassportInput{
				ObjectType:       entity.ObjectTypeResidentialHouse,
				OrganizationName: "ГУП БТИ",
				Address: entity.Address{
					Subject: "г. Москва",
					City:    "Москва",
					Street:  "ул. Тверская",
					House:   "1",
				},
				GeneralInfo: entity.GeneralInfo{
					Purpose:           "Жилое",
					ActualUsage:       "Жилое",
					ConstructionYear:  2020,
					TotalArea:         100.5,
					LivingArea:        70.0,
					FloorsAboveGround: 2,
				},
			},
			wantErr: false,
		},
		{
			name: "empty object type returns error",
			input: passport.CreatePassportInput{
				OrganizationName: "ГУП БТИ",
				Address: entity.Address{
					Subject: "г. Москва",
					House:   "1",
				},
				GeneralInfo: entity.GeneralInfo{
					Purpose:          "Жилое",
					ConstructionYear: 2020,
					TotalArea:        100.5,
				},
			},
			wantErr: true,
			errMsg:  "тип объекта обязателен",
		},
		{
			name: "invalid address returns error",
			input: passport.CreatePassportInput{
				ObjectType:       entity.ObjectTypeResidentialHouse,
				OrganizationName: "ГУП БТИ",
				Address: entity.Address{
					// Subject отсутствует - невалидно
					House: "1",
				},
				GeneralInfo: entity.GeneralInfo{
					Purpose:          "Жилое",
					ConstructionYear: 2020,
					TotalArea:        100.5,
				},
			},
			wantErr: true,
			errMsg:  "субъект РФ обязателен",
		},
		{
			name: "invalid general info returns error",
			input: passport.CreatePassportInput{
				ObjectType:       entity.ObjectTypeResidentialHouse,
				OrganizationName: "ГУП БТИ",
				Address: entity.Address{
					Subject: "г. Москва",
					House:   "1",
				},
				GeneralInfo: entity.GeneralInfo{
					// Purpose отсутствует - невалидно
					ConstructionYear: 2020,
					TotalArea:        100.5,
				},
			},
			wantErr: true,
			errMsg:  "назначение объекта обязательно",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo := memory.NewInMemoryPassportRepository()
			useCase := passport.NewCreatePassportUseCase(repo)
			ctx := context.Background()

			// Act
			output, err := useCase.Execute(ctx, tt.input)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				require.NotNil(t, output)
				assert.NotNil(t, output.Passport)
				assert.NotEmpty(t, output.Passport.ID)
				assert.Equal(t, tt.input.ObjectType, output.Passport.ObjectType)
				assert.Equal(t, tt.input.OrganizationName, output.Passport.OrganizationName)

				// Проверяем что паспорт сохранен в репозитории
				saved, err := repo.GetByID(ctx, output.Passport.ID)
				require.NoError(t, err)
				assert.Equal(t, output.Passport.ID, saved.ID)
			}
		})
	}
}
