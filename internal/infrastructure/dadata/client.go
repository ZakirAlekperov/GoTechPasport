package dadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	APIURL = "https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address"
	APIKey = "3f921259458d51b26e1aff0e74be9a6ac5c14c19"
)

// Client клиент для работы с DaData API
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient создает новый клиент DaData
func NewClient() *Client {
	return &Client{
		apiKey: APIKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SuggestionRequest запрос на получение подсказок
type SuggestionRequest struct {
	Query            string              `json:"query"`
	Count            int                 `json:"count,omitempty"`
	FromBound        *Bound              `json:"from_bound,omitempty"`
	ToBound          *Bound              `json:"to_bound,omitempty"`
	Locations        []Location          `json:"locations,omitempty"`
	RestrictValue    bool                `json:"restrict_value,omitempty"`
}

// Bound граница уровня детализации
type Bound struct {
	Value string `json:"value"`
}

// Location ограничение по региону
type Location struct {
	RegionFiasID string `json:"region_fias_id,omitempty"`
	CityFiasID   string `json:"city_fias_id,omitempty"`
	AreaFiasID   string `json:"area_fias_id,omitempty"`
}

// SuggestionResponse ответ с подсказками
type SuggestionResponse struct {
	Suggestions []Suggestion `json:"suggestions"`
}

// Suggestion подсказка адреса
type Suggestion struct {
	Value             string      `json:"value"`
	UnrestrictedValue string      `json:"unrestricted_value"`
	Data              AddressData `json:"data"`
}

// AddressData детальные данные адреса
type AddressData struct {
	PostalCode       string `json:"postal_code"`
	Country          string `json:"country"`
	RegionFiasID     string `json:"region_fias_id"`
	RegionKladrID    string `json:"region_kladr_id"`
	RegionWithType   string `json:"region_with_type"`
	RegionType       string `json:"region_type"`
	RegionTypeFull   string `json:"region_type_full"`
	Region           string `json:"region"`
	AreaFiasID       string `json:"area_fias_id"`
	AreaKladrID      string `json:"area_kladr_id"`
	AreaWithType     string `json:"area_with_type"`
	AreaType         string `json:"area_type"`
	AreaTypeFull     string `json:"area_type_full"`
	Area             string `json:"area"`
	CityFiasID       string `json:"city_fias_id"`
	CityKladrID      string `json:"city_kladr_id"`
	CityWithType     string `json:"city_with_type"`
	CityType         string `json:"city_type"`
	CityTypeFull     string `json:"city_type_full"`
	City             string `json:"city"`
	CityArea         string `json:"city_area"`
	CityDistrict     string `json:"city_district"`
	SettlementFiasID string `json:"settlement_fias_id"`
	SettlementKladrID string `json:"settlement_kladr_id"`
	SettlementWithType string `json:"settlement_with_type"`
	SettlementType   string `json:"settlement_type"`
	SettlementTypeFull string `json:"settlement_type_full"`
	Settlement       string `json:"settlement"`
	StreetFiasID     string `json:"street_fias_id"`
	StreetKladrID    string `json:"street_kladr_id"`
	StreetWithType   string `json:"street_with_type"`
	StreetType       string `json:"street_type"`
	StreetTypeFull   string `json:"street_type_full"`
	Street           string `json:"street"`
	HouseFiasID      string `json:"house_fias_id"`
	HouseKladrID     string `json:"house_kladr_id"`
	HouseType        string `json:"house_type"`
	HouseTypeFull    string `json:"house_type_full"`
	House            string `json:"house"`
	BlockType        string `json:"block_type"`
	BlockTypeFull    string `json:"block_type_full"`
	Block            string `json:"block"`
	FlatType         string `json:"flat_type"`
	FlatTypeFull     string `json:"flat_type_full"`
	Flat             string `json:"flat"`
	FiasID           string `json:"fias_id"`
	FiasLevel        string `json:"fias_level"`
	KladrID          string `json:"kladr_id"`
}

// SuggestAddress получает подсказки адресов
func (c *Client) SuggestAddress(query string, opts ...RequestOption) ([]Suggestion, error) {
	req := &SuggestionRequest{
		Query: query,
		Count: 10,
	}

	// Применяем опции
	for _, opt := range opts {
		opt(req)
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", APIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Token "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result SuggestionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Suggestions, nil
}

// RequestOption опция для запроса
type RequestOption func(*SuggestionRequest)

// WithRegionFilter фильтр по региону
func WithRegionFilter(regionFiasID string) RequestOption {
	return func(req *SuggestionRequest) {
		req.Locations = append(req.Locations, Location{
			RegionFiasID: regionFiasID,
		})
	}
}

// WithCityFilter фильтр по городу
func WithCityFilter(cityFiasID string) RequestOption {
	return func(req *SuggestionRequest) {
		req.Locations = append(req.Locations, Location{
			CityFiasID: cityFiasID,
		})
	}
}

// WithBounds установка границ детализации
func WithBounds(from, to string) RequestOption {
	return func(req *SuggestionRequest) {
		if from != "" {
			req.FromBound = &Bound{Value: from}
		}
		if to != "" {
			req.ToBound = &Bound{Value: to}
		}
	}
}

// WithCount количество подсказок
func WithCount(count int) RequestOption {
	return func(req *SuggestionRequest) {
		req.Count = count
	}
}

// SuggestRegions получает подсказки только регионов
func (c *Client) SuggestRegions(query string) ([]Suggestion, error) {
	return c.SuggestAddress(query, 
		WithBounds("region", "region"),
		WithCount(20),
	)
}

// SuggestCities получает подсказки городов в регионе
func (c *Client) SuggestCities(query string, regionFiasID string) ([]Suggestion, error) {
	opts := []RequestOption{
		WithBounds("city", "city"),
		WithCount(20),
	}
	if regionFiasID != "" {
		opts = append(opts, WithRegionFilter(regionFiasID))
	}
	return c.SuggestAddress(query, opts...)
}

// SuggestStreets получает подсказки улиц в городе
func (c *Client) SuggestStreets(query string, cityFiasID string) ([]Suggestion, error) {
	opts := []RequestOption{
		WithBounds("street", "street"),
		WithCount(20),
	}
	if cityFiasID != "" {
		opts = append(opts, WithCityFilter(cityFiasID))
	}
	return c.SuggestAddress(query, opts...)
}

// SuggestHouses получает подсказки домов на улице
func (c *Client) SuggestHouses(query string, streetFiasID string) ([]Suggestion, error) {
	return c.SuggestAddress(query,
		WithBounds("house", "house"),
		WithCount(20),
	)
}
