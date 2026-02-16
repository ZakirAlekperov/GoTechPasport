package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Создаем приложение Fyne
	myApp := app.New()
	myWindow := myApp.NewWindow("Технический паспорт недвижимости")

	// Временный UI для проверки
	welcomeLabel := widget.NewLabel("Добро пожаловать в GoTechPasport!")
	infoLabel := widget.NewLabel("Приложение для генерации технических паспортов")

	createButton := widget.NewButton("Создать новый паспорт", func() {
		fmt.Println("Создание нового паспорта...")
		// TODO: Реализовать логику создания
	})

	openButton := widget.NewButton("Открыть существующий", func() {
		fmt.Println("Открытие существующего паспорта...")
		// TODO: Реализовать логику открытия
	})

	content := container.NewVBox(
		welcomeLabel,
		infoLabel,
		widget.NewSeparator(),
		createButton,
		openButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.ShowAndRun()
}
