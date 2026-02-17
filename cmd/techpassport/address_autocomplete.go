package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/ZakirAlekperov/GoTechPasport/internal/infrastructure/dadata"
)

// AddressAutocomplete виджет с автодополнением адресов
type AddressAutocomplete struct {
	widget.Entry
	dadataClient    *dadata.Client
	suggestionsList *widget.List
	popup           *widget.PopUp
	window          fyne.Window
	suggestions     []dadata.Suggestion
	onSelected      func(dadata.Suggestion)
	suggestFunc     func(string) ([]dadata.Suggestion, error)
}

// NewAddressAutocomplete создает новый виджет автодополнения
func NewAddressAutocomplete(
	window fyne.Window,
	dadataClient *dadata.Client,
	suggestFunc func(string) ([]dadata.Suggestion, error),
	onSelected func(dadata.Suggestion),
) *AddressAutocomplete {
	aa := &AddressAutocomplete{
		window:       window,
		dadataClient: dadataClient,
		suggestFunc:  suggestFunc,
		onSelected:   onSelected,
		suggestions:  []dadata.Suggestion{},
	}

	aa.ExtendBaseWidget(aa)

	// Создаем список подсказок
	aa.suggestionsList = widget.NewList(
		func() int {
			return len(aa.suggestions)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextWrapOff
			return label
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(aa.suggestions) {
				obj.(*widget.Label).SetText(aa.suggestions[id].Value)
			}
		},
	)

	aa.suggestionsList.OnSelected = func(id widget.ListItemID) {
		if id < len(aa.suggestions) {
			selected := aa.suggestions[id]
			aa.SetText(selected.Value)
			
			// Сразу скрываем popup
			aa.hidePopup()
			
			if aa.onSelected != nil {
				aa.onSelected(selected)
			}
		}
	}

	// Обработчик изменения текста
	aa.OnChanged = func(text string) {
		if len(text) < 2 {
			aa.hidePopup()
			return
		}

		// Получаем подсказки асинхронно
		go aa.fetchSuggestions(text)
	}

	// Обработчик потери фокуса - скрываем popup
	aa.OnFocusLost = func() {
		// Небольшая задержка чтобы успел сработать клик по списку
		go func() {
			// time.Sleep(200 * time.Millisecond)
			aa.hidePopup()
		}()
	}

	return aa
}

func (aa *AddressAutocomplete) fetchSuggestions(query string) {
	suggestions, err := aa.suggestFunc(query)
	if err != nil {
		log.Printf("Error fetching suggestions: %v", err)
		return
	}

	aa.suggestions = suggestions
	aa.suggestionsList.Refresh()

	if len(suggestions) > 0 {
		aa.showPopup()
	} else {
		aa.hidePopup()
	}
}

func (aa *AddressAutocomplete) showPopup() {
	// Скрываем предыдущий popup если есть
	aa.hidePopup()

	// Вычисляем позицию popup относительно поля ввода
	canvas := aa.window.Canvas()
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(aa)
	size := aa.Size()

	// Создаем контейнер для списка
	content := container.NewMax(aa.suggestionsList)
	
	// Создаем popup
	aa.popup = widget.NewPopUp(content, canvas)
	
	// Позиция под полем ввода
	popupPos := fyne.NewPos(pos.X, pos.Y+size.Height)
	
	// Высота popup в зависимости от количества подсказок
	popupHeight := fyne.Min(200, float32(len(aa.suggestions))*35)
	popupSize := fyne.NewSize(size.Width, popupHeight)

	aa.popup.ShowAtPosition(popupPos)
	aa.popup.Resize(popupSize)
}

func (aa *AddressAutocomplete) hidePopup() {
	if aa.popup != nil {
		aa.popup.Hide()
		aa.popup = nil
	}
}

// AddressFormDaData форма адреса с DaData
type AddressFormDaData struct {
	window         fyne.Window
	dadataClient   *dadata.Client
	regionField    *AddressAutocomplete
	cityField      *AddressAutocomplete
	streetField    *AddressAutocomplete
	houseField     *AddressAutocomplete
	buildingField  *widget.Entry
	apartmentField *widget.Entry
	fullAddressLabel *widget.Label

	// Сохраненные FIAS ID для каскадной фильтрации
	selectedRegionFiasID string
	selectedCityFiasID   string
	selectedStreetFiasID string

	// Полный выбранный адрес
	selectedAddress *dadata.Suggestion
}

// NewAddressFormDaData создает новую форму адреса
func NewAddressFormDaData(window fyne.Window) *AddressFormDaData {
	form := &AddressFormDaData{
		window:       window,
		dadataClient: dadata.NewClient(),
	}

	// Поле региона
	form.regionField = NewAddressAutocomplete(
		window,
		form.dadataClient,
		func(query string) ([]dadata.Suggestion, error) {
			return form.dadataClient.SuggestRegions(query)
		},
		func(suggestion dadata.Suggestion) {
			form.selectedRegionFiasID = suggestion.Data.RegionFiasID
			form.cityField.SetText("")
			form.streetField.SetText("")
			form.houseField.SetText("")
			log.Printf("Выбран регион: %s", suggestion.Value)
		},
	)
	form.regionField.SetPlaceHolder("Начните вводить регион...")

	// Поле города
	form.cityField = NewAddressAutocomplete(
		window,
		form.dadataClient,
		func(query string) ([]dadata.Suggestion, error) {
			return form.dadataClient.SuggestCities(query, form.selectedRegionFiasID)
		},
		func(suggestion dadata.Suggestion) {
			form.selectedCityFiasID = suggestion.Data.CityFiasID
			if form.selectedCityFiasID == "" {
				form.selectedCityFiasID = suggestion.Data.SettlementFiasID
			}
			form.streetField.SetText("")
			form.houseField.SetText("")
			log.Printf("Выбран город: %s", suggestion.Value)
		},
	)
	form.cityField.SetPlaceHolder("Начните вводить город...")

	// Поле улицы
	form.streetField = NewAddressAutocomplete(
		window,
		form.dadataClient,
		func(query string) ([]dadata.Suggestion, error) {
			return form.dadataClient.SuggestStreets(query, form.selectedCityFiasID)
		},
		func(suggestion dadata.Suggestion) {
			form.selectedStreetFiasID = suggestion.Data.StreetFiasID
			form.houseField.SetText("")
			log.Printf("Выбрана улица: %s", suggestion.Value)
		},
	)
	form.streetField.SetPlaceHolder("Начните вводить улицу...")

	// Поле дома
	form.houseField = NewAddressAutocomplete(
		window,
		form.dadataClient,
		func(query string) ([]dadata.Suggestion, error) {
			fullQuery := fmt.Sprintf("%s %s %s %s",
				form.regionField.Text,
				form.cityField.Text,
				form.streetField.Text,
				query,
			)
			return form.dadataClient.SuggestAddress(fullQuery, dadata.WithBounds("house", "house"))
		},
		func(suggestion dadata.Suggestion) {
			form.selectedAddress = &suggestion
			log.Printf("Выбран дом: %s", suggestion.Value)
		},
	)
	form.houseField.SetPlaceHolder("Номер дома...")

	// Простые поля без автодополнения
	form.buildingField = widget.NewEntry()
	form.buildingField.SetPlaceHolder("Корпус/строение")

	form.apartmentField = widget.NewEntry()
	form.apartmentField.SetPlaceHolder("Квартира")

	// Метка для отображения полного адреса
	form.fullAddressLabel = widget.NewLabel("Адрес будет отображен после сохранения")
	form.fullAddressLabel.Wrapping = fyne.TextWrapWord

	return form
}

// GetAddressData возвращает заполненные данные адреса
func (f *AddressFormDaData) GetAddressData() map[string]string {
	return map[string]string{
		"subject":   f.regionField.Text,
		"city":      f.cityField.Text,
		"street":    f.streetField.Text,
		"house":     f.houseField.Text,
		"building":  f.buildingField.Text,
		"apartment": f.apartmentField.Text,
	}
}

// GetFullAddress возвращает полный адрес строкой
func (f *AddressFormDaData) GetFullAddress() string {
	address := ""
	
	if f.regionField.Text != "" {
		address += f.regionField.Text
	}
	
	if f.cityField.Text != "" {
		if address != "" {
			address += ", "
		}
		address += f.cityField.Text
	}
	
	if f.streetField.Text != "" {
		if address != "" {
			address += ", "
		}
		address += f.streetField.Text
	}
	
	if f.houseField.Text != "" {
		if address != "" {
			address += ", "
		}
		address += "д. " + f.houseField.Text
	}
	
	if f.buildingField.Text != "" {
		if address != "" {
			address += ", "
		}
		address += "корп. " + f.buildingField.Text
	}
	
	if f.apartmentField.Text != "" {
		if address != "" {
			address += ", "
		}
		address += "кв. " + f.apartmentField.Text
	}
	
	return address
}

// UpdateFullAddressLabel обновляет метку с полным адресом
func (f *AddressFormDaData) UpdateFullAddressLabel() {
	fullAddress := f.GetFullAddress()
	if fullAddress != "" {
		f.fullAddressLabel.SetText("📍 Полный адрес: " + fullAddress)
	} else {
		f.fullAddressLabel.SetText("Адрес будет отображен после заполнения полей")
	}
}

// CreateForm создает визуальную форму
func (f *AddressFormDaData) CreateForm() fyne.CanvasObject {
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Регион *:", Widget: f.regionField},
			{Text: "Город/Населенный пункт *:", Widget: f.cityField},
			{Text: "Улица *:", Widget: f.streetField},
			{Text: "Дом *:", Widget: f.houseField},
			{Text: "Корпус/Строение:", Widget: f.buildingField},
			{Text: "Квартира:", Widget: f.apartmentField},
		},
	}

	// Кнопка сохранить
	saveButton := widget.NewButton("💾 Сохранить адрес", func() {
		f.UpdateFullAddressLabel()
	})

	return container.NewBorder(
		nil,
		container.NewVBox(
			saveButton,
			widget.NewSeparator(),
			f.fullAddressLabel,
		),
		nil, nil,
		form,
	)
}
