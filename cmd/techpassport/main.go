package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/ZakirAlekperov/GoTechPasport/internal/domain/entity"
	"github.com/ZakirAlekperov/GoTechPasport/internal/infrastructure/storage/memory"
	"github.com/ZakirAlekperov/GoTechPasport/internal/usecase/passport"
)

// App главная структура приложения
type App struct {
	fyneApp    fyne.App
	window     fyne.Window
	repo       *memory.InMemoryPassportRepository
	createUC   *passport.CreatePassportUseCase
}

func main() {
	// Инициализируем приложение
	app := &App{
		fyneApp: app.New(),
		repo:    memory.NewInMemoryPassportRepository(),
	}
	app.createUC = passport.NewCreatePassportUseCase(app.repo)
	app.window = app.fyneApp.NewWindow("Технический паспорт недвижимости")

	// Показываем главное меню
	app.showMainMenu()

	app.window.Resize(fyne.NewSize(800, 600))
	app.window.ShowAndRun()
}

// showMainMenu показывает главное меню
func (a *App) showMainMenu() {
	welcomeLabel := widget.NewLabelWithStyle(
		"Добро пожаловать в GoTechPasport!",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	infoLabel := widget.NewLabel("Система управления техническими паспортами недвижимости")
	infoLabel.Alignment = fyne.TextAlignCenter

	createButton := widget.NewButton("Создать новый паспорт", func() {
		a.showCreatePassportForm()
	})

	listButton := widget.NewButton("Просмотр паспортов", func() {
		a.showPassportsList()
	})

	content := container.NewVBox(
		widget.NewSeparator(),
		welcomeLabel,
		infoLabel,
		widget.NewSeparator(),
		createButton,
		listButton,
		widget.NewSeparator(),
	)

	a.window.SetContent(container.NewCenter(content))
}

// showCreatePassportForm показывает форму создания паспорта
func (a *App) showCreatePassportForm() {
	// Поля формы
	orgNameEntry := widget.NewEntry()
	orgNameEntry.SetPlaceHolder("Например: ГУП БТИ")

	// Адрес
	subjectEntry := widget.NewEntry()
	subjectEntry.SetPlaceHolder("Например: г. Москва")

	cityEntry := widget.NewEntry()
	cityEntry.SetPlaceHolder("Например: Москва")

	streetEntry := widget.NewEntry()
	streetEntry.SetPlaceHolder("Например: ул. Тверская")

	houseEntry := widget.NewEntry()
	houseEntry.SetPlaceHolder("Например: 1")

	// Общие сведения
	purposeEntry := widget.NewEntry()
	purposeEntry.SetPlaceHolder("Например: Жилое")

	usageEntry := widget.NewEntry()
	usageEntry.SetPlaceHolder("Например: Жилое")

	yearEntry := widget.NewEntry()
	yearEntry.SetPlaceHolder("Например: 2020")

	totalAreaEntry := widget.NewEntry()
	totalAreaEntry.SetPlaceHolder("Например: 100.5")

	livingAreaEntry := widget.NewEntry()
	livingAreaEntry.SetPlaceHolder("Например: 70.0")

	floorsEntry := widget.NewEntry()
	floorsEntry.SetPlaceHolder("Например: 2")

	// Форма
	form := container.NewVBox(
		widget.NewLabelWithStyle("Создание технического паспорта", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),

		widget.NewLabel("Наименование организации техучета:"),
		orgNameEntry,

		widget.NewSeparator(),
		widget.NewLabelWithStyle("Адрес объекта", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),

		widget.NewLabel("Субъект РФ:"),
		subjectEntry,

		widget.NewLabel("Город:"),
		cityEntry,

		widget.NewLabel("Улица:"),
		streetEntry,

		widget.NewLabel("Дом:"),
	houseEntry,

		widget.NewSeparator(),
		widget.NewLabelWithStyle("Общие сведения", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),

		widget.NewLabel("Назначение:"),
		purposeEntry,

		widget.NewLabel("Фактическое использование:"),
		usageEntry,

		widget.NewLabel("Год постройки:"),
		yearEntry,

		widget.NewLabel("Общая площадь (кв.м):"),
		totalAreaEntry,

		widget.NewLabel("Жилая площадь (кв.м):"),
		livingAreaEntry,

		widget.NewLabel("Этажность:"),
		floorsEntry,

		widget.NewSeparator(),
	)

	// Кнопки
	createBtn := widget.NewButton("Создать", func() {
		a.handleCreatePassport(
			orgNameEntry.Text,
			subjectEntry.Text,
			cityEntry.Text,
			streetEntry.Text,
			houseEntry.Text,
			purposeEntry.Text,
			usageEntry.Text,
			yearEntry.Text,
			totalAreaEntry.Text,
			livingAreaEntry.Text,
			floorsEntry.Text,
		)
	})

	cancelBtn := widget.NewButton("Отмена", func() {
		a.showMainMenu()
	})

	buttonBox := container.NewHBox(createBtn, cancelBtn)

	// Скроллинг
	scroll := container.NewVScroll(form)
	content := container.NewBorder(nil, buttonBox, nil, nil, scroll)

	a.window.SetContent(content)
}

// handleCreatePassport обрабатывает создание паспорта
func (a *App) handleCreatePassport(orgName, subject, city, street, house, purpose, usage, year, totalArea, livingArea, floors string) {
	// Простой парсинг (в реальном приложении нужна полная валидация)
	var yearInt int
	var totalAreaFloat, livingAreaFloat float64
	var floorsInt int

	fmt.Sscanf(year, "%d", &yearInt)
	fmt.Sscanf(totalArea, "%f", &totalAreaFloat)
	fmt.Sscanf(livingArea, "%f", &livingAreaFloat)
	fmt.Sscanf(floors, "%d", &floorsInt)

	input := passport.CreatePassportInput{
		ObjectType:       entity.ObjectTypeResidentialHouse,
		OrganizationName: orgName,
		Address: entity.Address{
			Subject: subject,
			City:    city,
			Street:  street,
			House:   house,
		},
		GeneralInfo: entity.GeneralInfo{
			Purpose:           purpose,
			ActualUsage:       usage,
			ConstructionYear:  yearInt,
			TotalArea:         totalAreaFloat,
			LivingArea:        livingAreaFloat,
			FloorsAboveGround: floorsInt,
		},
	}

	ctx := context.Background()
	output, err := a.createUC.Execute(ctx, input)

	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	// Успех
	msg := fmt.Sprintf("Технический паспорт успешно создан!\n\nID: %s\nАдрес: %s",
		output.Passport.ID,
		output.Passport.Address.FullAddress())

	dialog.ShowInformation("Успех", msg, a.window)
	a.showMainMenu()
}

// showPassportsList показывает список паспортов
func (a *App) showPassportsList() {
	ctx := context.Background()
	passports, err := a.repo.List(ctx)

	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	if len(passports) == 0 {
		dialog.ShowInformation("Информация", "Список паспортов пуст", a.window)
		return
	}

	// Создаем список
	var items []string
	for _, p := range passports {
		item := fmt.Sprintf("ID: %s | %s | %s",
			p.ID,
			p.Address.FullAddress(),
			p.GeneralInfo.Purpose)
		items = append(items, item)
	}

	list := widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(items[id])
		},
	)

	backBtn := widget.NewButton("Назад", func() {
		a.showMainMenu()
	})

	content := container.NewBorder(
		widget.NewLabelWithStyle("Список технических паспортов", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		backBtn,
		nil,
		nil,
		list,
	)

	a.window.SetContent(content)
}
