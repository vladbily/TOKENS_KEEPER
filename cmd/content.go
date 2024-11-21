package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func createWcontent(win fyne.Window, DB *gorm.DB, DB_A *gorm.DB) fyne.CanvasObject {
	if err := DB.Find(&tokens).Error; err != nil {
		return nil
	}

	btn_add := widget.NewButton("ADD TOKEN", func() {
		setContent(win, create_window)
	})
	btn_add.Resize(fyne.NewSize(200, 50))
	btn_add.Move(fyne.NewPos(300, 500))

	menubar := container.NewHBox(
		setText("TOKENS", 25),
	)

	for i := range tokens {
		newPrice := parser(tokens[i].Name, DB_A)
		tokens[i].Current_Price = newPrice

		var info Infos
		if err := DB.Where("name = ?", tokens[i].Name).First(&info).Error; err != nil {
			newInfo := Infos{Name: tokens[i].Name, Current_Price: newPrice}
			if err := DB.Create(&newInfo).Error; err != nil {
				fmt.Println("Error creating new record:", err)
			}
		} else {
			info.Current_Price = newPrice
			if err := DB.Save(&info).Error; err != nil {
				fmt.Println("Error updating record:", err)
			}
		}
	}

	tS := totalSum(*DB)
	total_sum := checkTotal(tS)
	ref_btn := widget.NewButton("REFRESH", func() {
		tS = totalSum(*DB)
		total_sum.Text = fmt.Sprintf("%.2f", tS)
		total_sum.Refresh()
	})

	setting_btn := widget.NewButton("SETTINGS", func() {
		setContent(win, s_content)
	})

	setting_btn.Resize(fyne.NewSize(80, 40))
	setting_btn.Move(fyne.NewPos(700, 535))

	ref_btn.Resize(fyne.NewSize(80, 40))
	ref_btn.Move(fyne.NewPos(10, 535))

	token_list = createTokenlist()

	token_scroll := container.NewScroll(token_list)

	token_scroll.Resize(fyne.NewSize(790, 450))
	token_scroll.Move(fyne.NewPos(0, 40))

	token_list.OnSelected = func(id widget.ListItemID) {
		onTokenSelected(win, DB, id)
	}

	total_txt := setText("TOTAL: ", 20)
	total_txt.Move(fyne.NewPos(300, 1))

	w_content = container.NewWithoutLayout(
		menubar,
		newLine(3, 0, 40, 800, 40),
		token_scroll,
		btn_add,
		total_sum,
		total_txt,
		ref_btn,
		setting_btn,
	)
	return w_content
}

func createCreateContent(win fyne.Window, DB *gorm.DB, DB_A *gorm.DB) fyne.CanvasObject {
	back_btn = widget.NewButton("BACK", func() {
		DB.Find(&tokens)
		setContent(win, w_content)
		token_list.UnselectAll()
	})

	back_btn.Resize(fyne.NewSize(50, 35))
	back_btn.Move(fyne.NewPos(730, 1))

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("TOKEN NAME")
	dateEntry := widget.NewEntry()
	dateEntry.SetPlaceHolder("DD.MM.YYYY")
	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("PRICE")
	countEntry := widget.NewEntry()
	countEntry.SetPlaceHolder("COUNT")

	save_btn = widget.NewButton("SAVE", func() {
		price_flt, _ := strconv.ParseFloat(priceEntry.Text, 64)
		count_flt, _ := strconv.ParseFloat(countEntry.Text, 64)
		var exApi API
		result := DB_A.Where("Id = ?", 1).First(&exApi)
		if result.RowsAffected > 0 {
			if nameEntry.Text != "" && dateEntry.Text != "" && price_flt > 0 && count_flt > 0 {
				var existingToken Infos
				result := DB.Where("Name = ?", nameEntry.Text).First(&existingToken)
				if result.RowsAffected > 0 {

					history := TokenHistory{
						TokenId:       existingToken.Id,
						Name:          existingToken.Name,
						Date:          dateEntry.Text,
						Buy_Price:     price_flt,
						Count:         count_flt,
						Current_Price: parser(existingToken.Name, DB_A),
					}
					DB.Create(&history)

					newCount := existingToken.Count
					additionalCount, _ := strconv.ParseFloat(countEntry.Text, 64)
					existingToken.Count = newCount + additionalCount
					DB.Save(&existingToken)
				} else {
					token := Infos{
						Name:          nameEntry.Text,
						Date:          dateEntry.Text,
						Avg_Price:     price_flt,
						Count:         count_flt,
						Current_Price: parser(nameEntry.Text, DB_A),
					}
					DB.Create(&token)
					history := TokenHistory{
						TokenId:       token.Id,
						Name:          token.Name,
						Date:          dateEntry.Text,
						Buy_Price:     price_flt,
						Count:         count_flt,
						Current_Price: parser(token.Name, DB_A),
					}
					DB.Create(&history)
				}

				DB.Find(&histories)
				DB.Find(&tokens)

				token_list.Refresh()

				setContent(win, w_content)

				token_list.UnselectAll()

				nameEntry.Text = ""
				nameEntry.Refresh()
				dateEntry.Text = ""
				dateEntry.Refresh()
				priceEntry.Text = ""
				priceEntry.Refresh()
				countEntry.Text = ""
				countEntry.Refresh()
			} else {
				notification := dialog.NewInformation("ERROR", "The entered parameters are incorrect.", win)
				notification.Show()
			}
		} else {
			notification := dialog.NewInformation("ERROR", "Enter your API in the settings.", win)
			notification.Show()
		}
	})

	save_btn.Resize(fyne.NewSize(200, 50))
	save_btn.Move(fyne.NewPos(300, 500))

	save_cont := container.NewWithoutLayout(
		save_btn,
	)

	back_btn = widget.NewButton("BACK", func() {
		nameEntry.Text = ""
		nameEntry.Refresh()
		dateEntry.Text = ""
		dateEntry.Refresh()
		priceEntry.Text = ""
		priceEntry.Refresh()
		countEntry.Text = ""
		countEntry.Refresh()
		token_list.UnselectAll()
		setContent(win, w_content)
	})

	back_btn.Resize(fyne.NewSize(50, 35))
	back_btn.Move(fyne.NewPos(730, 1))

	create_text := container.NewVBox(setText("FILL IN THE INFORMATION", 25))
	create_text.Resize(fyne.NewSize(800, 600))
	create_text.Move(fyne.NewPos(0, 0))

	createbar := container.NewWithoutLayout(
		create_text,
		back_btn,
	)

	cont_wind := container.NewVBox(
		createbar,
		newLine(3, 0, 40, 800, 40),
		nameEntry,
		dateEntry,
		priceEntry,
		countEntry,
	)

	cont_wind.Resize(fyne.NewSize(790, 500))
	cont_wind.Move(fyne.NewPos(0, 0))

	create_window = container.NewWithoutLayout(
		cont_wind,
		save_cont,
	)
	return create_window
}

func createSettingContent(win fyne.Window, DB_A *gorm.DB) fyne.CanvasObject {
	back_set_btn = widget.NewButton("BACK", func() {
		setContent(win, w_content)
	})

	back_set_btn.Resize(fyne.NewSize(50, 35))
	back_set_btn.Move(fyne.NewPos(730, 1))

	create_text := container.NewVBox(setText("SETTINGS", 25))
	create_text.Resize(fyne.NewSize(800, 600))
	create_text.Move(fyne.NewPos(0, 0))

	settingbar := container.NewWithoutLayout(
		create_text,
		back_set_btn,
	)

	APIentry := widget.NewEntry()
	var exApi API
	result := DB_A.Where("Id = ?", 1).First(&exApi)
	if result.RowsAffected > 0 {
		APIentry.Text = (exApi.Api)
		APIentry.SetPlaceHolder("SET YOUR API")
	} else {
		APIentry.SetPlaceHolder("SET YOUR API")
	}

	save_set_btn = widget.NewButton("SAVE", func() {
		input := APIentry.Text
		var exApi API
		result := DB_A.Where("Id = ?", 1).First(&exApi)
		if input != "" {
			if result.RowsAffected > 0 {
				exApi.Api = input
				DB_A.Save(&exApi)
			} else {
				api := API{
					Api: input,
				}
				DB_A.Create(&api)
			}
			notification := dialog.NewInformation("NOTIFICATION", "Restart the application to work correctly.", win)
			notification.Show()
			setContent(win, w_content)
		} else {
			notification := dialog.NewInformation("ERROR", "The input fields should not be empty.", win)
			notification.Show()
		}
	})

	save_set_btn.Resize(fyne.NewSize(200, 50))
	save_set_btn.Move(fyne.NewPos(300, 500))

	save_cont := container.NewWithoutLayout(
		save_set_btn,
	)

	s_con := container.NewVBox(
		settingbar,
		newLine(3, 0, 40, 800, 40),
		APIentry,
	)
	s_con.Resize(fyne.NewSize(790, 500))
	s_con.Move(fyne.NewPos(0, 0))

	s_content = container.NewWithoutLayout(
		s_con,
		save_cont,
	)
	return s_content
}
