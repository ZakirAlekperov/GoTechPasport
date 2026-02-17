package dadata_test

import (
	"testing"

	"github.com/ZakirAlekperov/GoTechPasport/internal/infrastructure/dadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SuggestRegions(t *testing.T) {
	t.Skip("Skipping integration test - requires API key")

	client := dadata.NewClient()

	// Тест поиска Москвы
	suggestions, err := client.SuggestRegions("Москв")
	require.NoError(t, err)
	require.NotEmpty(t, suggestions)

	// Проверяем что первая подсказка - Москва
	assert.Contains(t, suggestions[0].Value, "Москва")
	assert.NotEmpty(t, suggestions[0].Data.RegionFiasID)
}

func TestClient_SuggestCities(t *testing.T) {
	t.Skip("Skipping integration test - requires API key")

	client := dadata.NewClient()

	// Сначала получаем регион
	regions, err := client.SuggestRegions("Московская")
	require.NoError(t, err)
	require.NotEmpty(t, regions)

	regionFiasID := regions[0].Data.RegionFiasID

	// Теперь ищем города в этом регионе
	cities, err := client.SuggestCities("Подольск", regionFiasID)
	require.NoError(t, err)
	require.NotEmpty(t, cities)

	assert.Contains(t, cities[0].Value, "Подольск")
}

func TestClient_SuggestStreets(t *testing.T) {
	t.Skip("Skipping integration test - requires API key")

	client := dadata.NewClient()

	// Получаем Москву
	regions, err := client.SuggestRegions("Москва")
	require.NoError(t, err)
	require.NotEmpty(t, regions)

	// Ищем улицы в Москве
	streets, err := client.SuggestStreets("Тверская", regions[0].Data.CityFiasID)
	require.NoError(t, err)
	require.NotEmpty(t, streets)

	assert.Contains(t, streets[0].Value, "Тверская")
}
