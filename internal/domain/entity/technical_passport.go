package entity

import "time"

// TechnicalPassport представляет технический паспорт объекта недвижимости
// Это основная агрегатная сущность (Аggregate Root) в домене
type TechnicalPassport struct {
	// Уникальный идентификатор
	ID string `json:"id"`

	// Тип объекта (жилой дом, квартира и т.п.)
	ObjectType ObjectType `json:"object_type"`

	// Адрес объекта
	Address Address `json:"address"`

	// Наименование организации технического учета (ОТИ)
	OrganizationName string `json:"organization_name"`

	// Инвентарный номер
	InventoryNumber string `json:"inventory_number,omitempty"`

	// Кадастровый номер
	CadastralNumber string `json:"cadastral_number,omitempty"`

	// Дата создания паспорта
	CreatedDate time.Time `json:"created_date"`

	// Дата последнего обновления
	UpdatedDate time.Time `json:"updated_date"`

	// Дата, по состоянию на которую составлен паспорт
	AsOfDate time.Time `json:"as_of_date"`

	// ======================== Разделы паспорта ========================

	// 1. Общие сведения
	GeneralInfo GeneralInfo `json:"general_info"`

	// 2. Состав объекта (здания, сооружения)
	Buildings []Building `json:"buildings"`

	// 3. Правообладатели
	Owners []Owner `json:"owners"`

	// 4. Ситуационный план (ссылка на файл или данные)
	SituationPlanPath string `json:"situation_plan_path,omitempty"`

	// 5. Благоустройство
	Utilities Utilities `json:"utilities"`

	// 6-7. Поэтажные планы и экспликация
	FloorPlans []string `json:"floor_plans,omitempty"`  // пути к файлам
	Explication []Room `json:"explication"`

	// История изменений
	AuditLog []AuditEntry `json:"audit_log,omitempty"`
}

// NewTechnicalPassport создает новый технический паспорт
func NewTechnicalPassport(objectType ObjectType, address Address) *TechnicalPassport {
	now := time.Now()
	return &TechnicalPassport{
		ObjectType:  objectType,
		Address:     address,
		CreatedDate: now,
		UpdatedDate: now,
		AsOfDate:    now,
		Buildings:   []Building{},
		Owners:      []Owner{},
		FloorPlans:  []string{},
		Explication: []Room{},
		AuditLog:    []AuditEntry{},
	}
}

// IsValid проверяет корректность заполнения паспорта
func (tp *TechnicalPassport) IsValid() error {
	// Проверка адреса
	if err := tp.Address.IsValid(); err != nil {
		return err
	}

	// Проверка общих сведений
	if err := tp.GeneralInfo.IsValid(); err != nil {
		return err
	}

	// Проверка зданий
	if len(tp.Buildings) == 0 {
		return ValidationError{Field: "buildings", Message: "должно быть хотя бы одно здание"}
	}

	for i, building := range tp.Buildings {
		if err := building.IsValid(); err != nil {
			return ValidationError{
				Field:   "buildings[" + string(rune(i)) + "]",
				Message: err.Error(),
			}
		}
	}

	// Проверка владельцев
	if len(tp.Owners) == 0 {
		return ValidationError{Field: "owners", Message: "должен быть хотя бы один правообладатель"}
	}

	for i, owner := range tp.Owners {
		if err := owner.IsValid(); err != nil {
			return ValidationError{
				Field:   "owners[" + string(rune(i)) + "]",
				Message: err.Error(),
			}
		}
	}

	// Проверка благоустройства
	if err := tp.Utilities.IsValid(); err != nil {
		return err
	}

	// Проверка экспликации
	for i, room := range tp.Explication {
		if err := room.IsValid(); err != nil {
			return ValidationError{
				Field:   "explication[" + string(rune(i)) + "]",
				Message: err.Error(),
			}
		}
	}

	return nil
}

// AddBuilding добавляет здание в состав объекта
func (tp *TechnicalPassport) AddBuilding(building Building) error {
	if err := building.IsValid(); err != nil {
		return err
	}

	tp.Buildings = append(tp.Buildings, building)
	tp.UpdatedDate = time.Now()
	tp.AddAuditEntry("add_building", "Добавлено здание: " + building.Name)

	return nil
}

// AddOwner добавляет правообладателя
func (tp *TechnicalPassport) AddOwner(owner Owner) error {
	if err := owner.IsValid(); err != nil {
		return err
	}

	tp.Owners = append(tp.Owners, owner)
	tp.UpdatedDate = time.Now()
	tp.AddAuditEntry("add_owner", "Добавлен правообладатель")

	return nil
}

// AddRoom добавляет помещение в экспликацию
func (tp *TechnicalPassport) AddRoom(room Room) error {
	if err := room.IsValid(); err != nil {
		return err
	}

	tp.Explication = append(tp.Explication, room)
	tp.UpdatedDate = time.Now()
	tp.AddAuditEntry("add_room", "Добавлено помещение: " + room.RoomNumber)

	return nil
}

// AddAuditEntry добавляет запись в историю изменений
func (tp *TechnicalPassport) AddAuditEntry(action, description string) {
	entry := AuditEntry{
		Timestamp:   time.Now(),
		Action:      action,
		Description: description,
	}
	tp.AuditLog = append(tp.AuditLog, entry)
}

// CalculateTotalArea вычисляет общую площадь по всем зданиям
func (tp *TechnicalPassport) CalculateTotalArea() float64 {
	total := 0.0
	for _, building := range tp.Buildings {
		total += building.TotalArea
	}
	return total
}

// CalculateTotalInventoryValue вычисляет общую инвентаризационную стоимость
func (tp *TechnicalPassport) CalculateTotalInventoryValue() float64 {
	total := 0.0
	for _, building := range tp.Buildings {
		total += building.InventoryValue
	}
	return total
}
