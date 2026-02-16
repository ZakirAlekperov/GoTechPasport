package entity

import "time"

// Owner представляет правообладателя объекта недвижимости
type Owner struct {
	// Дата записи
	EntryDate time.Time `json:"entry_date"`

	// Тип лица (физическое/юридическое)
	PersonType PersonType `json:"person_type"`

	// ФИО (для физических лиц)
	FullName string `json:"full_name,omitempty"`

	// Паспортные данные (для физических лиц)
	PassportData string `json:"passport_data,omitempty"`

	// Наименование организации (для юридических лиц)
	CompanyName string `json:"company_name,omitempty"`

	// ИНН
	TIN string `json:"tin,omitempty"`

	// Вид права (собственность, аренда и т.п.)
	RightType string `json:"right_type"`

	// Правоустанавливающий документ
	RightDocument string `json:"right_document"`

	// Доля в праве (например, "1/2", "1" для единоличного собственника)
	Share string `json:"share"`
}

// IsValid проверяет корректность данных владельца
func (o *Owner) IsValid() error {
	if o.EntryDate.IsZero() {
		return ValidationError{Field: "entry_date", Message: "дата записи обязательна"}
	}

	if o.PersonType == "" {
		return ValidationError{Field: "person_type", Message: "тип лица обязателен"}
	}

	if o.PersonType == PersonTypeIndividual && o.FullName == "" {
		return ValidationError{Field: "full_name", Message: "ФИО обязательно для физического лица"}
	}

	if o.PersonType == PersonTypeLegal && o.CompanyName == "" {
		return ValidationError{Field: "company_name", Message: "наименование организации обязательно для юридического лица"}
	}

	if o.RightType == "" {
		return ValidationError{Field: "right_type", Message: "вид права обязателен"}
	}

	if o.RightDocument == "" {
		return ValidationError{Field: "right_document", Message: "правоустанавливающий документ обязателен"}
	}

	if o.Share == "" {
		return ValidationError{Field: "share", Message: "доля в праве обязательна"}
	}

	return nil
}
