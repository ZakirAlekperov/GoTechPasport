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

func TestAddBuildingUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*memory.InMemoryPassportRepository) string // returns passport ID
		input   func(string) passport.AddBuildingInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid building adds successfully",
			setup: func(repo *memory.InMemoryPassportRepository) string {
				p := entity.NewTechnicalPassport(
					entity.ObjectTypeResidentialHouse,
					entity.Address{Subject: "г. Москва", House: "1"},
				)
				p.ID = "TEST-1"
				p.OrganizationName = "ГУП БТИ"
				p.GeneralInfo = entity.GeneralInfo{
					Purpose:          "Жилое",
					ConstructionYear: 2020,
					TotalArea:        100.0,
				}
				repo.Create(context.Background(), p)
				return p.ID
			},
			input: func(id string) passport.AddBuildingInput {
				return passport.AddBuildingInput{
					PassportID: id,
					Building: entity.Building{
						Litera:          "А",
						Name:            "Жилой дом",
						Purpose:         "Жилое",
						TotalArea:       100.5,
						FloorsAboveGround: 2,
						ConstructionYear: 2020,
					},
				}
			},
			wantErr: false,
		},
		{
			name: "empty passport ID returns error",
			setup: func(repo *memory.InMemoryPassportRepository) string {
				return ""
			},
			input: func(id string) passport.AddBuildingInput {
				return passport.AddBuildingInput{
					PassportID: "",
					Building: entity.Building{
						Litera: "А",
						Name:   "Жилой дом",
					},
				}
			},
			wantErr: true,
			errMsg:  "ID паспорта обязателен",
		},
		{
			name: "invalid building returns error",
			setup: func(repo *memory.InMemoryPassportRepository) string {
				p := entity.NewTechnicalPassport(
					entity.ObjectTypeResidentialHouse,
					entity.Address{Subject: "г. Москва", House: "1"},
				)
				p.ID = "TEST-2"
				p.OrganizationName = "ГУП БТИ"
				p.GeneralInfo = entity.GeneralInfo{
					Purpose:          "Жилое",
					ConstructionYear: 2020,
					TotalArea:        100.0,
				}
				repo.Create(context.Background(), p)
				return p.ID
			},
			input: func(id string) passport.AddBuildingInput {
				return passport.AddBuildingInput{
					PassportID: id,
					Building: entity.Building{
						// Литера отсутствует - невалидно
						Name: "Жилой дом",
					},
				}
			},
			wantErr: true,
			errMsg:  "литера обязательна",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo := memory.NewInMemoryPassportRepository()
			passportID := tt.setup(repo)
			useCase := passport.NewAddBuildingUseCase(repo)
			ctx := context.Background()
			input := tt.input(passportID)

			// Act
			output, err := useCase.Execute(ctx, input)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				require.NotNil(t, output)
				assert.NotNil(t, output.Passport)
				assert.Len(t, output.Passport.Buildings, 1)
				assert.Equal(t, input.Building.Litera, output.Passport.Buildings[0].Litera)

				// Проверяем что изменения сохранены
				saved, err := repo.GetByID(ctx, passportID)
				require.NoError(t, err)
				assert.Len(t, saved.Buildings, 1)
			}
		})
	}
}
