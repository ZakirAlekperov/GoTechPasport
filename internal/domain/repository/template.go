package repository

import "context"

// TemplateRepository определяет интерфейс для работы с шаблонами документов
type TemplateRepository interface {
	// LoadTemplate загружает шаблон по имени
	LoadTemplate(ctx context.Context, name string) ([]byte, error)

	// ListTemplates возвращает список доступных шаблонов
	ListTemplates(ctx context.Context) ([]string, error)
}
