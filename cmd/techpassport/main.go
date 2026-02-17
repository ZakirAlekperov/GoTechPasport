package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
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
	addBuildingUC *passport.AddBuildingUseCase
	removeBuildingUC *passport.RemoveBuildingUseCase

	// Текущий редактируемый паспорт
	currentPassport *entity.TechnicalPassport

	// Поля форм для всех вкладок
	generalFields   *GeneralInfoFields
	addressForm     *AddressFormDaData // Новая форма с DaData
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
	orgName           *widget.Entry
	purpose           *widget.Entry
	usage             *widget.Entry
	year              *widget.Entry
	totalArea         *widget.Entry
	livingArea        *widget.Entry
	floors            *widget.Entry
	undergroundFloors *widget.Entry
}

// UtilitiesFields поля благоустройства
type UtilitiesFields struct {
	waterCentral          *widget.Entry
	waterAutonomous       *widget.Entry
	sewerageCentral       *widget.Entry
	sewerageAutonomous    *widget.Entry
	heatingCentral        *widget.Entry
	heatingAutonomous     *widget.Entry
	gasCentral            *widget.Entry
	gasAutonomous         *widget.Entry
	electricityCentral    *widget.Entry
	electricityAutonomous *widget.Entry
}

func main() {
	// Инициализируем приложение
	myApp := app.New()

	// Устанавливаем светлую тему
	myApp.Settings().SetTheme(theme.LightTheme())

	app := &App{
		fyneApp:   myApp,
		repo:      memory.NewInMemoryPassportRepository(),
		buildings: []entity.Building{},
		owners:    []entity.Owner{},
		rooms:     []entity.Room{},
	}
	app.createUC = passport.NewCreatePassportUseCase(app.repo)
	app.addBuildingUC = passport.NewAddBuildingUseCase(app.repo)
	app.removeBuildingUC = passport.NewRemoveBuildingUseCase(app.repo)
	app.window = app.fyneApp.NewWindow("Технический паспорт недвижимости")

	// Создаем новый пустой паспорт для редактирования
	app.initNewPassport()

	// Создаем меню
	menu := app.createMenu()
	app.window.SetMainMenu(menu)

	// Создаем интерфейс с вкладками
	tabs := app.createTabs()

	// Создаем тулбар с кнопками
	toolbar := app.createToolbar()

	// Основной контейнер
	content := container.NewBorder(toolbar, nil, nil, nil, tabs)

	app.window.SetContent(content)
	app.window.Resize(fyne.NewSize(1000, 700))
	app.window.ShowAndRun()
}

// initNewPassport инициализирует новый паспорт
func (a *App) initNewPassport() {
	a.currentPassport = entity.NewTechnicalPassport(
		entity.ObjectTypeResidentialHouse,
		entity.Address{},
	)
}

// createMenu создает главное меню
func (a *App) createMenu() *fyne.MainMenu {
	fileMenu := fyne.NewMenu("Файл",
		fyne.NewMenuItem("Новый", func() {
			a.createNewPassport()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Сохранить", func() {
			a.savePassport()
		}),
		fyne.NewMenuItem("Открыть...", func() {
			a.openPassport()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Экспорт в PDF", func() {
			dialog.ShowInformation("В разработке", "Функция экспорта будет реализована на следующем этапе", a.window)
		}),
		fyne.NewMenuItem("Экспорт в Word", func() {
			dialog.ShowInformation("В разработке", "Функция экспорта будет реализована на следующем этапе", a.window)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Выход", func() {
			a.fyneApp.Quit()
		}),
	)

	helpMenu := fyne.NewMenu("Справка",
		fyne.NewMenuItem("О программе", func() {
			dialog.ShowInformation("О GoTechPasport",
				"GoTechPasport v0.2\n\nПриложение для генерации технических паспортов недвижимости.\n\nСогласно Приказу Минэкономразвития РФ от 17 августа 2006 г. № 244\n\nИнтеграция с DaData для валидации адресов\n\n© 2026 Zakir Alekperov",
				a.window)
		}),
	)

	return fyne.NewMainMenu(fileMenu, helpMenu)
}

// createToolbar создает панель инструментов
func (a *App) createToolbar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			a.createNewPassport()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			a.savePassport()
		}),
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			a.openPassport()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			dialog.ShowInformation("Справка",
				"Используйте вкладки для заполнения разделов паспорта.\nВкладка 'Адрес' использует подсказки DaData для предотвращения ошибок.\nНачните вводить регион, затем город, улицу и дом.\nНажмите кнопку 'Сохранить адрес' чтобы увидеть полный адрес.\n\nНажмите 'Сохранить' в меню для сохранения всего паспорта.",
				a.window)
		}),
	)
}

// createTabs создает все вкладки паспорта
func (a *App) createTabs() *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Общие сведения", a.createGeneralInfoTab()),
		container.NewTabItem("Адрес 🔍", a.createAddressTab()),
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

// createAddressTab создает вкладку "Адрес" с DaData
func (a *App) createAddressTab() fyne.CanvasObject {
	a.addressForm = NewAddressFormDaData(a.window)

	infoLabel := widget.NewLabel("💡 Используйте подсказки DaData для точного ввода адреса. Нажмите 'Сохранить адрес' для отображения полного адреса.")
	infoLabel.Wrapping = fyne.TextWrapWord

	// CreateForm теперь возвращает fyne.CanvasObject с кнопкой и меткой
	formContent := a.addressForm.CreateForm()

	return container.NewBorder(
		infoLabel,
		nil, nil, nil,
		container.NewVScroll(formContent),
	)
}

// createBuildingsTab создает вкладку "Состав объекта"
func (a *App) createBuildingsTab() fyne.CanvasObject {
	a.buildingsList = widget.NewList(
		func() int { return len(a.buildings) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(fmt.Sprintf("Лит. %s - %s (%.2f кв.м, %d г.)",
				a.buildings[id].Litera,
				a.buildings[id].Name,
				a.buildings[id].TotalArea,
				a.buildings[id].CommissionYear))
		},
	)

	addBtn := widget.NewButton("Добавить здание", func() {
		a.showAddBuildingDialog()
	})

	removeBtn := widget.NewButton("Удалить", func() {
		if len(a.buildings) == 0 {
			dialog.ShowInformation("Информация", "Список зданий пуст", a.window)
			return
		}
		// TODO: реализовать выбор здания для удаления
		dialog.ShowInformation("В разработке", "Выберите здание в списке для удаления", a.window)
	})

	info := widget.NewLabel("Список зданий и сооружений в составе объекта")
	buttonBox := container.NewHBox(addBtn, removeBtn)

	return container.NewBorder(info, buttonBox, nil, nil, a.buildingsList)
}

// showAddBuildingDialog показывает диалог добавления здания
func (a *App) showAddBuildingDialog() {
	// Поля формы
	litera := widget.NewEntry()
	litera.SetPlaceHolder("Например: А")

	name := widget.NewEntry()
	name.SetPlaceHolder("Например: Жилой дом")

	commissionYear := widget.NewEntry()
	commissionYear.SetPlaceHolder("Например: 2020")

	totalArea := widget.NewEntry()
	totalArea.SetPlaceHolder("Например: 100.5")

	height := widget.NewEntry()
	height.SetPlaceHolder("Например: 6.0")

	wallMaterial := widget.NewEntry()
	wallMaterial.SetPlaceHolder("Например: Кирпич")

	// Форма
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Литера *:", Widget: litera},
			{Text: "Наименование *:", Widget: name},
			{Text: "Год ввода в эксплуатацию *:", Widget: commissionYear},
			{Text: "Общая площадь (кв.м) *:", Widget: totalArea},
			{Text: "Высота (м):", Widget: height},
			{Text: "Материал стен:", Widget: wallMaterial},
		},
	}

	// Диалог
	dlg := dialog.NewCustomConfirm(
		"Добавить здание",
		"Добавить",
		"Отмена",
		container.NewVScroll(form),
		func(confirmed bool) {
			if confirmed {
				// Парсим данные
				var year int
				var area, heightVal float64

				fmt.Sscanf(commissionYear.Text, "%d", &year)
				fmt.Sscanf(totalArea.Text, "%f", &area)
				fmt.Sscanf(height.Text, "%f", &heightVal)

				building := entity.Building{
					Litera:         litera.Text,
					Name:           name.Text,
					CommissionYear: year,
					TotalArea:      area,
					Height:         heightVal,
					WallMaterial:   wallMaterial.Text,
				}

				// Добавляем в локальный список (если паспорт еще не сохранен)
				if a.currentPassport.ID == "" {
					// Паспорт еще не сохранен - добавляем в локальный список
					a.buildings = append(a.buildings, building)
					a.buildingsList.Refresh()
					dialog.ShowInformation("Успех", "Здание добавлено", a.window)
				} else {
					// Паспорт сохранен - используем use case
					input := passport.AddBuildingInput{
						PassportID: a.currentPassport.ID,
						Building:   building,
					}

					ctx := context.Background()
					output, err := a.addBuildingUC.Execute(ctx, input)

					if err != nil {
						dialog.ShowError(err, a.window)
						return
					}

					// Обновляем текущий паспорт и список
					a.currentPassport = output.Passport
					a.buildings = output.Passport.Buildings
					a.buildingsList.Refresh()
					dialog.ShowInformation("Успех", "Здание добавлено в паспорт", a.window)
				}
			}
		},
		a.window,
	)

	dlg.Resize(fyne.NewSize(500, 400))
	dlg.Show()
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

	info := widget.NewLabel("Список правообладателей объекта недвижимости")

	return container.NewBorder(info, addBtn, nil, nil, a.ownersList)
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

	info := widget.NewLabel("Экспликация помещений (расшифровка площадей)")

	return container.NewBorder(info, addBtn, nil, nil, a.roomsList)
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

	a.utilitiesFields.waterCentral.SetPlaceHolder("0.0")
	a.utilitiesFields.waterAutonomous.SetPlaceHolder("0.0")
	a.utilitiesFields.sewerageCentral.SetPlaceHolder("0.0")
	a.utilitiesFields.sewerageAutonomous.SetPlaceHolder("0.0")
	a.utilitiesFields.heatingCentral.SetPlaceHolder("0.0")
	a.utilitiesFields.heatingAutonomous.SetPlaceHolder("0.0")
	a.utilitiesFields.gasCentral.SetPlaceHolder("0.0")
	a.utilitiesFields.gasAutonomous.SetPlaceHolder("0.0")
	a.utilitiesFields.electricityCentral.SetPlaceHolder("0.0")
	a.utilitiesFields.electricityAutonomous.SetPlaceHolder("0.0")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Водоснабжение (центр. кв.м):", Widget: a.utilitiesFields.waterCentral},
			{Text: "Водоснабжение (автон. кв.м):", Widget: a.utilitiesFields.waterAutonomous},
			{Text: "Канализация (центр. кв.м):", Widget: a.utilitiesFields.sewerageCentral},
			{Text: "Канализация (автон. кв.м):", Widget: a.utilitiesFields.sewerageAutonomous},
			{Text: "Отопление (центр. кв.м):", Widget: a.utilitiesFields.heatingCentral},
			{Text: "Отопление (автон. кв.м):", Widget: a.utilitiesFields.heatingAutonomous},
			{Text: "Газоснабжение (центр. кв.м):", Widget: a.utilitiesFields.gasCentral},
			{Text: "Газоснабжение (автон. кв.м):", Widget: a.utilitiesFields.gasAutonomous},
			{Text: "Электроснабжение (центр. кв.м):", Widget: a.utilitiesFields.electricityCentral},
			{Text: "Электроснабжение (автон. кв.м):", Widget: a.utilitiesFields.electricityAutonomous},
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

	// Добавляем все здания из локального списка
	for _, building := range a.buildings {
		buildingInput := passport.AddBuildingInput{
			PassportID: a.currentPassport.ID,
			Building:   building,
		}
		a.addBuildingUC.Execute(ctx, buildingInput)
	}

	msg := fmt.Sprintf("Технический паспорт успешно сохранен!\n\nID: %s\nАдрес: %s\nЗданий: %d",
		output.Passport.ID,
		output.Passport.Address.FullAddress(),
		len(a.buildings))

	dialog.ShowInformation("Успех", msg, a.window)
}

// openPassport открывает список паспортов
func (a *App) openPassport() {
	ctx := context.Background()
	passports, err := a.repo.List(ctx)

	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	if len(passports) == 0 {
		dialog.ShowInformation("Информация", "Список паспортов пуст.\nСоздайте новый паспорт.", a.window)
		return
	}

	// Создаем список для выбора
	var items []string
	for _, p := range passports {
		item := fmt.Sprintf("%s - %s", p.ID, p.Address.FullAddress())
		items = append(items, item)
	}

	// Показываем список
	list := widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(items[id])
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		// TODO: Загрузить паспорт в форму
		dialog.ShowInformation("В разработке",
			"Загрузка паспорта будет реализована на следующем этапе",
			a.window)
	}

	dialog.ShowCustom("Открыть паспорт", "Закрыть", list, a.window)
}

// collectDataFromFields собирает данные из всех полей формы
func (a *App) collectDataFromFields() error {
	// Организация
	a.currentPassport.OrganizationName = a.generalFields.orgName.Text

	// Адрес из DaData формы
	addressData := a.addressForm.GetAddressData()
	a.currentPassport.Address = entity.Address{
		Subject: addressData["subject"],
		City:    addressData["city"],
		Street:  addressData["street"],
		House:   addressData["house"],
		Building: addressData["building"],
		Apartment: addressData["apartment"],
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

	// TODO: Собрать данные о Utilities

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

	// Адрес - создаем новую форму
	a.addressForm = NewAddressFormDaData(a.window)

	// Utilities
	a.utilitiesFields.waterCentral.SetText("")
	a.utilitiesFields.waterAutonomous.SetText("")
	a.utilitiesFields.sewerageCentral.SetText("")
	a.utilitiesFields.sewerageAutonomous.SetText("")
	a.utilitiesFields.heatingCentral.SetText("")
	a.utilitiesFields.heatingAutonomous.SetText("")
	a.utilitiesFields.gasCentral.SetText("")
	a.utilitiesFields.gasAutonomous.SetText("")
	a.utilitiesFields.electricityCentral.SetText("")
	a.utilitiesFields.electricityAutonomous.SetText("")

	// Очищаем списки
	a.buildings = []entity.Building{}
	a.owners = []entity.Owner{}
	a.rooms = []entity.Room{}

	a.buildingsList.Refresh()
	a.ownersList.Refresh()
	a.roomsList.Refresh()
}
