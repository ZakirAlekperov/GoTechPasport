package entity

// Building представляет здание или сооружение в составе объекта
type Building struct {
	// Литера (буквенное обозначение на плане)
	Litera string `json:"litera"`

	// Наименование (жилой дом, сарай, гараж и т.п.)
	Name string `json:"name"`

	// Год ввода в эксплуатацию
	CommissionYear int `json:"commission_year"`

	// Материал стен
	WallMaterial string `json:"wall_material"`

	// Общая площадь (кв.м)
	TotalArea float64 `json:"total_area"`

	// Площадь застройки (кв.м)
	BuildArea float64 `json:"build_area"`

	// Высота (м)
	Height float64 `json:"height"`

	// Объем (куб.м)
	Volume float64 `json:"volume"`

	// Инвентаризационная стоимость (руб)
	InventoryValue float64 `json:"inventory_value"`
}

// IsValid проверяет корректность данных здания
func (b *Building) IsValid() error {
	if b.Litera == "" {
		return ValidationError{Field: "litera", Message: "литера обязательна"}
	}

	if b.Name == "" {
		return ValidationError{Field: "name", Message: "наименование обязательно"}
	}

	if b.CommissionYear < 1800 || b.CommissionYear > 2100 {
		return ValidationError{Field: "commission_year", Message: "некорректный год ввода в эксплуатацию"}
	}

	if b.TotalArea < 0 {
		return ValidationError{Field: "total_area", Message: "площадь не может быть отрицательной"}
	}

	if b.Height < 0 {
		return ValidationError{Field: "height", Message: "высота не может быть отрицательной"}
	}

	if b.Volume < 0 {
		return ValidationError{Field: "volume", Message: "объем не может быть отрицательным"}
	}

	return nil
}
