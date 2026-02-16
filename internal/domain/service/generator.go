package service

import (
	"context"

	"github.com/ZakirAlekperov/GoTechPasport/internal/domain/entity"
)

// DocumentFormat представляет формат генерируемого документа
type DocumentFormat string

const (
	// FormatPDF генерирует PDF документ
	FormatPDF DocumentFormat = "pdf"

	// FormatDOCX генерирует Word документ
	FormatDOCX DocumentFormat = "docx"
)

// GenerateOptions опции для генерации документа
type GenerateOptions struct {
	// Format формат выходного документа
	Format DocumentFormat

	// TemplateName имя шаблона для использования
	TemplateName string

	// IncludeImages включать ли изображения (планы)
	IncludeImages bool
}

// DocumentGenerator определяет интерфейс для генерации документов
type DocumentGenerator interface {
	// Generate генерирует документ из технического паспорта
	// Возвращает данные документа и ошибку
	Generate(ctx context.Context, passport *entity.TechnicalPassport, options GenerateOptions) ([]byte, error)

	// SaveToFile сохраняет генерированный документ в файл
	SaveToFile(ctx context.Context, passport *entity.TechnicalPassport, outputPath string, options GenerateOptions) error
}
