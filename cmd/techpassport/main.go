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
	fyneApp  fyne.App
	window   fyne.Window
	repo     *memory.InMemoryPassportRepository
	createUC *passport.CreatePassportUseCase

	// Текущий редактируемый паспорт
	currentPassport *entity.TechnicalPassport

	// Поля форм для всех вкладок
	generalFields   *GeneralInfoFields
	addressFields   *AddressFields
	buildingsList   *widget.List
	buildings       []entity.Building
	ownersList      *widget.List
	owners          []entity.Owner
	roomsList       *widget.List
	rooms           []entity.Room
	utilitiesFields *UtilitiesFields
}

// GeneralInfoFields поля общих сведений
type GeneralInfoFields struct {
	orgName         *widget.Entry
	purpose         *widget.Entry
	usage           *widget.Entry
	year            *widget.Entry
	totalArea       *widget.Entry
	livingArea      *widget.Entry
	floors          *widget.Entry
	undergroundFloors *widget.Entry
}

// AddressFields поля адреса
type AddressFields struct {
	subject      *widget.Entry
	district     *widget.Entry
	city         *widget.Entry
	cityDistrict *widget.Entry
	street       *widget.Entry
	house        *widget.Entry
	building     *widget.Entry
	apartment    *widget.Entry
}

// UtilitiesFields поля благоустройства
type UtilitiesFields struct {
	waterCentral    *widget.Entry
	waterAutonomous *widget.Entry
	sewerageCentral *widget.Entry
	sewerageAutonomous *widget.Entry
	heatingCentral  *widget.Entry
	heatingAutonomous *widget.Entry
	gasCentral      *widget.Entry
	gasAutonomous   *widget.Entry
	electricityCentral *widget.Entry
	electricityAutonomous *widget.Entry
}

func main() {
	// Инициализируем приложение
	app := &App{
		fyneApp:  app.New(),
		repo:     memory.NewInMemoryPassportRepository(),
		buildings: []entity.Building{},
		owners:    []entity.Owner{},
		rooms:     []entity.Room{},
	}
	app.createUC = passport.NewCreatePassportUseCase(app.repo)
	app.window = app.fyneApp.NewWindow("Технический паспорт недвижимости")

	// Создаем новый пустой паспорт для редактирования
	app.initNewPassport()

	// Создаем интерфейс с вкладками
	tabs := app.createTabs()

	// Создаем тулбар с кнопками
	toolbar := app.createToolbar()

	// Основной контейнер
	content := container.NewBorder(toolbar, nil, nil, nil, tabs)

	app.window.SetContent(content)
	app.window.Resize(fyne.NewSize(900, 700))
	app.window.ShowAndRun()
}

// initNewPassport инициализирует новый паспорт
func (a *App) initNewPassport() {
	a.currentPassport = entity.NewTechnicalPassport(
		entity.ObjectTypeResidentialHouse,
		entity.Address{},
	)
}

// createToolbar создает панель инструментов
func (a *App) createToolbar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(widget.NewIcon(nil), func() {
			a.createNewPassport()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(widget.NewIcon(nil), func() {
			a.savePassport()
		}),
		widget.NewToolbarSpacer(),
	)
}

// createTabs создает все вкладки паспорта
func (a *App) createTabs() *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Общие сведения", a.createGeneralInfoTab()),
		container.NewTabItem("Адрес", a.createAddressTab()),
		container.NewTabItem("Состав объекта", a.createBuildingsTab()),
		container.NewTabItem("Правообладатели", a.createOwnersTab()),
		container.NewTabItem("Экспликация", a.createRoomsTab()),
		container.NewTabItem("Благоустройство", a.createUtilitiesTab()),
	)
}

// createGeneralInfoTab создает вкладку "Общие сведения"
func (a *App) createGeneralInfoTab() fyne.CanvasObject {
	a.generalFields = &GeneralInfoFields{
		orgName:           widget.NewEntry(),
		purpose:           widget.NewEntry(),
		usage:             widget.NewEntry(),
		year:              widget.NewEntry(),
		totalArea:         widget.NewEntry(),
		livingArea:        widget.NewEntry(),
		floors:            widget.NewEntry(),
		undergroundFloors: widget.NewEntry(),
	}

	a.generalFields.orgName.SetPlaceHolder("Например: ГУП БТИ")
	a.generalFields.purpose.SetPlaceHolder("Например: Жилое")
	a.generalFields.usage.SetPlaceHolder("Например: Жилое")
	a.generalFields.year.SetPlaceHolder("Например: 2020")
	a.generalFields.totalArea.SetPlaceHolder("Например: 100.5")
	a.generalFields.livingArea.SetPlaceHolder("Например: 70.0")
	a.generalFields.floors.SetPlaceHolder("Например: 2")
	a.generalFields.undergroundFloors.SetPlaceHolder("Например: 0")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Организация техучета:", Widget: a.generalFields.orgName},
			{Text: "Назначение:", Widget: a.generalFields.purpose},
			{Text: "Фактическое использование:", Widget: a.generalFields.usage},
			{Text: "Год постройки:", Widget: a.generalFields.year},
			{Text: "Общая площадь (кв.м):", Widget: a.generalFields.totalArea},
			{Text: "Жилая площадь (кв.м):", Widget: a.generalFields.livingArea},
			{Text: "Этажей надземных:", Widget: a.generalFields.floors},
			{Text: "Этажей подземных:", Widget: a.generalFields.undergroundFloors},
		},
	}

	return container.NewVScroll(form)
}

// createAddressTab создает вкладку "Адрес"
func (a *App) createAddressTab() fyne.CanvasObject {
	a.addressFields = &AddressFields{
		subject:      widget.NewEntry(),
		district:     widget.NewEntry(),
		city:         widget.NewEntry(),
		cityDistrict: widget.NewEntry(),
		street:       widget.NewEntry(),
		house:        widget.NewEntry(),
		building:     widget.NewEntry(),
		apartment:    widget.NewEntry(),
	}

	a.addressFields.subject.SetPlaceHolder("Например: г. Москва")
	a.addressFields.district.SetPlaceHolder("Например: Центральный АО")
	a.addressFields.city.SetPlaceHolder("Например: Москва")
	a.addressFields.cityDistrict.SetPlaceHolder("Например: Тверской район")
	a.addressFields.street.SetPlaceHolder("Например: ул. Тверская")
	a.addressFields.house.SetPlaceHolder("Например: 1")
	a.addressFields.building.SetPlaceHolder("Например: корп. 2")
	a.addressFields.apartment.SetPlaceHolder("Например: 10")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Субъект РФ:", Widget: a.addressFields.subject},
			{Text: "Административный район:", Widget: a.addressFields.district},
			{Text: "Город:", Widget: a.addressFields.city},
			{Text: "Район города:", Widget: a.addressFields.cityDistrict},
			{Text: "Улица:", Widget: a.addressFields.street},
			{Text: "Дом:", Widget: a.addressFields.house},
			{Text: "Строение/корпус:", Widget: a.addressFields.building},
			{Text: "Квартира:", Widget: a.addressFields.apartment},
		},
	}

	return container.NewVScroll(form)
}

// createBuildingsTab создает вкладку "Состав объекта"
func (a *App) createBuildingsTab() fyne.CanvasObject {
	a.buildingsList = widget.NewList(
		func() int { return len(a.buildings) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(fmt.Sprintf("%s - %s (%.2f кв.м)",
				a.buildings[id].Litera,
				a.buildings[id].Name,
				a.buildings[id].TotalArea))
		},
	)

	addBtn := widget.NewButton("Добавить здание", func() {
		dialog.ShowInformation("В разработке", "Форма добавления здания будет реализована на следующем этапе", a.window)
	})

	return container.NewBorder(nil, addBtn, nil, nil, a.buildingsList)
}

// createOwnersTab создает вкладку "Правообладатели"
func (a *App) createOwnersTab() fyne.CanvasObject {
	a.ownersList = widget.NewList(
		func() int { return len(a.owners) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			owner := a.owners[id]
			name := owner.FullName
			if name == "" {
				name = owner.CompanyName
			}
			obj.(*widget.Label).SetText(fmt.Sprintf("%s - %s", name, owner.RightType))
		},
	)

	addBtn := widget.NewButton("Добавить правообладателя", func() {
		dialog.ShowInformation("В разработке", "Форма добавления владельца будет реализована на следующем этапе", a.window)
	})

	return container.NewBorder(nil, addBtn, nil, nil, a.ownersList)
}

// createRoomsTab создает вкладку "Экспликация"
func (a *App) createRoomsTab() fyne.CanvasObject {
	a.roomsList = widget.NewList(
		func() int { return len(a.rooms) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			room := a.rooms[id]
			obj.(*widget.Label).SetText(fmt.Sprintf("Литера %s, эт. %s, пом. %s - %s (%.2f кв.м)",
				room.Litera, room.Floor, room.RoomNumber, room.Purpose, room.Area))
		},
	)

	addBtn := widget.NewButton("Добавить помещение", func() {
		dialog.ShowInformation("В разработке", "Форма добавления помещения будет реализована на следующем этапе", a.window)
	})

	return container.NewBorder(nil, addBtn, nil, nil, a.roomsList)
}

// createUtilitiesTab создает вкладку "Благоустройство"
func (a *App) createUtilitiesTab() fyne.CanvasObject {
	a.utilitiesFields = &UtilitiesFields{
		waterCentral:          widget.NewEntry(),
		waterAutonomous:       widget.NewEntry(),
		sewerageCentral:       widget.NewEntry(),
		sewerageAutonomous:    widget.NewEntry(),
		heatingCentral:        widget.NewEntry(),
		heatingAutonomous:     widget.NewEntry(),
		gasCentral:            widget.NewEntry(),
		gasAutonomous:         widget.NewEntry(),
		electricityCentral:    widget.NewEntry(),
		electricityAutonomous: widget.NewEntry(),
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Водоснабжение (централизованное):", Widget: a.utilitiesFields.waterCentral},
			{Text: "Водоснабжение (автономное):", Widget: a.utilitiesFields.waterAutonomous},
			{Text: "Канализация (централизованная):", Widget: a.utilitiesFields.sewerageCentral},
			{Text: "Канализация (автономная):", Widget: a.utilitiesFields.sewerageAutonomous},
			{Text: "Отопление (централизованное):", Widget: a.utilitiesFields.heatingCentral},
			{Text: "Отопление (автономное):", Widget: a.utilitiesFields.heatingAutonomous},
			{Text: "Газоснабжение (централизованное):", Widget: a.utilitiesFields.gasCentral},
			{Text: "Газоснабжение (автономное):", Widget: a.utilitiesFields.gasAutonomous},
			{Text: "Электроснабжение (централизованное):", Widget: a.utilitiesFields.electricityCentral},
			{Text: "Электроснабжение (автономное):", Widget: a.utilitiesFields.electricityAutonomous},
		},
	}

	return container.NewVScroll(form)
}

// createNewPassport создает новый паспорт
func (a *App) createNewPassport() {
	dialog.ShowConfirm("Новый паспорт", 
		"Создать новый паспорт? Несохраненные данные будут потеряны.", 
		func(confirmed bool) {
			if confirmed {
				a.initNewPassport()
				a.clearAllFields()
				dialog.ShowInformation("Успех", "Новый паспорт создан", a.window)
			}
		}, a.window)
}

// savePassport сохраняет паспорт
func (a *App) savePassport() {
	// Собираем данные из полей
	if err := a.collectDataFromFields(); err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	// Создаем input для use case
	input := passport.CreatePassportInput{
		ObjectType:       a.currentPassport.ObjectType,
		OrganizationName: a.currentPassport.OrganizationName,
		Address:          a.currentPassport.Address,
		GeneralInfo:      a.currentPassport.GeneralInfo,
	}

	ctx := context.Background()
	output, err := a.createUC.Execute(ctx, input)

	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	// Обновляем текущий паспорт
	a.currentPassport = output.Passport

	msg := fmt.Sprintf("Технический паспорт успешно сохранен!\n\nID: %s\nАдрес: %s",
		output.Passport.ID,
		output.Passport.Address.FullAddress())

	dialog.ShowInformation("Успех", msg, a.window)
}

// collectDataFromFields собирает данные из всех полей формы
func (a *App) collectDataFromFields() error {
	// Организация
	a.currentPassport.OrganizationName = a.generalFields.orgName.Text

	// Адрес
	a.currentPassport.Address = entity.Address{
		Subject:      a.addressFields.subject.Text,
		District:     a.addressFields.district.Text,
		City:         a.addressFields.city.Text,
		CityDistrict: a.addressFields.cityDistrict.Text,
		Street:       a.addressFields.street.Text,
		House:        a.addressFields.house.Text,
		Building:     a.addressFields.building.Text,
		Apartment:    a.addressFields.apartment.Text,
	}

	// Общие сведения
	var year, floors, undergroundFloors int
	var totalArea, livingArea float64

	fmt.Sscanf(a.generalFields.year.Text, "%d", &year)
	fmt.Sscanf(a.generalFields.floors.Text, "%d", &floors)
	fmt.Sscanf(a.generalFields.undergroundFloors.Text, "%d", &undergroundFloors)
	fmt.Sscanf(a.generalFields.totalArea.Text, "%f", &totalArea)
	fmt.Sscanf(a.generalFields.livingArea.Text, "%f", &livingArea)

	a.currentPassport.GeneralInfo = entity.GeneralInfo{
		Purpose:           a.generalFields.purpose.Text,
		ActualUsage:       a.generalFields.usage.Text,
		ConstructionYear:  year,
		TotalArea:         totalArea,
		LivingArea:        livingArea,
		FloorsAboveGround: floors,
		FloorsUnderground: undergroundFloors,
	}

	// TODO: Utilities

	return nil
}

// clearAllFields очищает все поля формы
func (a *App) clearAllFields() {
	// Общие сведения
	a.generalFields.orgName.SetText("")
	a.generalFields.purpose.SetText("")
	a.generalFields.usage.SetText("")
	a.generalFields.year.SetText("")
	a.generalFields.totalArea.SetText("")
	a.generalFields.livingArea.SetText("")
	a.generalFields.floors.SetText("")
	a.generalFields.undergroundFloors.SetText("")

	// Адрес
	a.addressFields.subject.SetText("")
	a.addressFields.district.SetText("")
	a.addressFields.city.SetText("")
	a.addressFields.cityDistrict.SetText("")
	a.addressFields.street.SetText("")
	a.addressFields.house.SetText("")
	a.addressFields.building.SetText("")
	a.addressFields.apartment.SetText("")

	// Очищаем списки
	a.buildings = []entity.Building{}
	a.owners = []entity.Owner{}
	a.rooms = []entity.Room{}

	a.buildingsList.Refresh()
	a.ownersList.Refresh()
	a.roomsList.Refresh()
}
