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
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(aa.suggestions[id].Value)
		},
	)

	aa.suggestionsList.OnSelected = func(id widget.ListItemID) {
		if id < len(aa.suggestions) {
			selected := aa.suggestions[id]
			aa.SetText(selected.Value)
			if aa.popup != nil {
				aa.popup.Hide()
			}
			if aa.onSelected != nil {
				aa.onSelected(selected)
			}
		}
	}

	// Обработчик изменения текста
	aa.OnChanged = func(text string) {
		if len(text) < 2 {
			if aa.popup != nil {
				aa.popup.Hide()
			}
			return
		}

		// Получаем подсказки
		go aa.fetchSuggestions(text)
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
	} else if aa.popup != nil {
		aa.popup.Hide()
	}
}

func (aa *AddressAutocomplete) showPopup() {
	if aa.popup == nil {
		// Создаем контейнер с ограниченной высотой
		content := container.NewMax(aa.suggestionsList)
		aa.popup = widget.NewPopUp(content, aa.window.Canvas())
	}

	// Располагаем popup под полем ввода
	canvasPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(&aa.Entry)
	aa.popup.ShowAtPosition(fyne.NewPos(
		canvasPos.X,
		canvasPos.Y+aa.Size().Height,
	))

	// Устанавливаем размер
	aa.popup.Resize(fyne.NewSize(
		aa.Size().Width,
		fyne.Min(200, float32(len(aa.suggestions))*40),
	))
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

	return form
}

// GetAddressData возвращает заполненные данные адреса
func (f *AddressFormDaData) GetAddressData() map[string]string {
	return map[string]string{
		"subject":  f.regionField.Text,
		"city":     f.cityField.Text,
		"street":   f.streetField.Text,
		"house":    f.houseField.Text,
		"building": f.buildingField.Text,
		"apartment": f.apartmentField.Text,
	}
}

// CreateForm создает визуальную форму
func (f *AddressFormDaData) CreateForm() *widget.Form {
	return &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Регион *:", Widget: f.regionField},
			{Text: "Город/Населенный пункт *:", Widget: f.cityField},
			{Text: "Улица *:", Widget: f.streetField},
			{Text: "Дом *:", Widget: f.houseField},
			{Text: "Корпус/Строение:", Widget: f.buildingField},
			{Text: "Квартира:", Widget: f.apartmentField},
		},
	}
}
