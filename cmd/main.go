package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	a := app.NewWithID("com.tokens.ru")
	a.Settings().SetTheme(theme.DarkTheme())
	win := a.NewWindow("TOKENS")

	back_img = canvas.NewImageFromFile("BACK.jpg")

	win.CenterOnScreen()
	win.Resize(fyne.NewSize(800, 600))
	win.SetFixedSize(true)

	DB, _ := gorm.Open(sqlite.Open("tokens.db"), &gorm.Config{})
	DB.AutoMigrate(&Infos{})
	DB.Find(&tokens)

	DB_A, _ := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})
	DB_A.AutoMigrate(&API{})

	initializeDatabase(DB)
	createWcontent(win, DB, DB_A)
	createCreateContent(win, DB, DB_A)
	createSettingContent(win, DB_A)
	setContent(win, w_content)
	win.ShowAndRun()
}
