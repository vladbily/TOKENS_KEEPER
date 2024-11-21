package main

import (
	"fmt"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func createTokenlist() *widget.List {
	token_list = widget.NewList(func() int {
		return len(tokens)
	}, func() fyne.CanvasObject {
		container := container.NewHBox()
		for i := 0; i < 6; i++ {
			label := widget.NewLabel("")
			container.Add(label)
		}
		return container
	}, func(lii widget.ListItemID, co fyne.CanvasObject) {
		labels := co.(*fyne.Container).Objects
		labels[0].(*widget.Label).SetText("Date: " + tokens[lii].Date)
		labels[1].(*widget.Label).SetText("Name: " + tokens[lii].Name)
		labels[2].(*widget.Label).SetText("Count: " + fmt.Sprintf("%.4f", tokens[lii].Count))
		labels[3].(*widget.Label).SetText("Current Price: " + fmt.Sprintf("%.2f", tokens[lii].Current_Price))
		labels[4].(*widget.Label).SetText("Total: " + fmt.Sprintf("%.2f", tokens[lii].Current_Price*tokens[lii].Count))
	})
	token_list.Refresh()
	return token_list
}

func createHistorylist(DB *gorm.DB) *widget.List {
	history_list := widget.NewList(
		func() int {
			return len(histories)
		},
		func() fyne.CanvasObject {
			container := container.NewHBox()
			for i := 0; i < 5; i++ {
				label := widget.NewLabel("")
				container.Add(label)
			}
			return container
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			crc := histories[lii].Count
			prc := histories[lii].Buy_Price
			tokenName := histories[lii].Name
			currentPrice, _ := getCurrentPriceFromDB(DB, tokenName)
			labels := co.(*fyne.Container).Objects
			labels[0].(*widget.Label).SetText("Date: " + histories[lii].Date)
			labels[1].(*widget.Label).SetText("Buy Price: " + fmt.Sprintf("%.2f", prc))
			labels[2].(*widget.Label).SetText("Count: " + fmt.Sprintf("%.4f", crc))
			labels[3].(*widget.Label).SetText("Difference: " + fmt.Sprintf("%.2f", crc*(currentPrice-prc)))
		},
	)

	history_list.Refresh()
	return history_list
}

func onTokenSelected(win fyne.Window, DB *gorm.DB, id widget.ListItemID) {
	detailsBar := container.NewWithoutLayout(setText("DETAILS ABOUT "+fmt.Sprint(tokens[id].Name), 25), back_btn)

	if err := DB.Where("token_id = ?", tokens[id].Id).Find(&histories).Error; err != nil {
		dialog.ShowError(err, win)
		return
	}

	name_str := setText("NAME", 20)
	date_str := setText("DATE", 20)
	avg_str := setText("AVG BUY", 20)
	count_str := setText("COUNT", 20)
	curr_str := setText("CURR", 20)
	total_str := setText("TOTAL", 20)

	setPos(name_str, 40, 20)
	setPos(date_str, 40, 50)
	setPos(avg_str, 40, 80)
	setPos(count_str, 40, 110)
	setPos(curr_str, 40, 140)
	setPos(total_str, 40, 170)

	win_title := container.NewWithoutLayout(name_str, date_str, avg_str, count_str, curr_str, total_str)
	win_title.Move(fyne.NewPos(0, 40))

	DB.Where("token_id = ?", tokens[id].Id).Find(&histories)

	history_list := createHistorylist(DB)
	token_scroll := container.NewScroll(history_list)

	token_scroll.Resize(fyne.NewSize(800, 270))
	token_scroll.Move(fyne.NewPos(0, 270))

	history_list.OnSelected = func(historyID widget.ListItemID) {
		editHistory(win, DB, historyID)
	}

	currr, _ := strconv.ParseFloat(fmt.Sprint(tokens[id].Current_Price), 64)
	countt, _ := strconv.ParseFloat(fmt.Sprint(tokens[id].Count), 64)

	var averagePrice float64
	var avg1 float64

	if !math.IsNaN(tokens[id].Avg_Price) && tokens[id].Avg_Price != 0 {
		if err := DB.Model(&TokenHistory{}).
			Select("AVG(buy_price)").
			Where("token_id = ?", tokens[id].Id).
			Scan(&averagePrice).Error; err != nil {
			return
		}
		var info Infos
		if err := DB.Where("Id = ?", tokens[id].Id).First(&info).Error; err != nil {
			return
		}

		info.Avg_Price = averagePrice
		if err := DB.Save(&info).Error; err != nil {
			return
		}

		avg1 = averagePrice
	} else {
		avg1 = tokens[id].Avg_Price
	}

	name := setText(tokens[id].Name, 20)
	date := setText(tokens[id].Date, 20)
	avg := setText(fmt.Sprintf("%.2f", avg1), 20)
	count := setText(fmt.Sprintf("%.4f", tokens[id].Count), 20)
	curr := setText(fmt.Sprintf("%.2f", tokens[id].Current_Price), 20)
	total := setText(fmt.Sprintf("%.2f", currr*countt), 20)

	setPos(name, 640, 20)
	setPos(date, 640, 50)
	setPos(avg, 640, 80)
	setPos(count, 640, 110)
	setPos(curr, 640, 140)
	setPos(total, 640, 170)

	token_info := container.NewWithoutLayout(name, date, avg, count, curr, total)
	token_info.Move(fyne.NewPos(0, 40))

	back_btn := widget.NewButton("BACK", func() {
		setContent(win, w_content)
		token_list.UnselectAll()
	})

	back_btn.Resize(fyne.NewSize(50, 35))
	back_btn.Move(fyne.NewPos(730, 1))

	detailsBox := container.NewWithoutLayout(detailsBar,
		newLine(3, 0, 40, 800, 40),
		win_title,
		newLine(2, 40, 90, 760, 90),
		token_info,
		newLine(2, 40, 120, 760, 120),
		newLine(2, 40, 150, 760, 150),
		newLine(2, 40, 180, 760, 180),
		newLine(2, 40, 210, 760, 210),
		newLine(2, 40, 240, 760, 240),
		token_info,
		newLine(3, 0, 270, 800, 270),
		token_scroll,
	)
	token_list.Refresh()
	token_list.UnselectAll()
	setContent(win, detailsBox)
}

func editHistory(win fyne.Window, DB *gorm.DB, id int) {
	editName := widget.NewEntry()
	editName.SetText(histories[id].Name)
	editDate := widget.NewEntry()
	editDate.SetText(histories[id].Date)
	editPrice := widget.NewEntry()
	editPrice.SetText(fmt.Sprint(histories[id].Buy_Price))
	editCount := widget.NewEntry()
	editCount.SetText(fmt.Sprint(histories[id].Count))

	editSave := widget.NewButton("SAVE CHANGES", func() {
		var currentCount float64
		if err := DB.Model(&TokenHistory{}).Where("Id = ?", histories[id].Id).Select("Count").Scan(&currentCount).Error; err != nil {
			fmt.Println("Error retrieving current token count:", err)
			return
		}

		editPrice_flt, err := strconv.ParseFloat(editPrice.Text, 64)
		if err != nil || editPrice_flt <= 0 {
			dialog.ShowInformation("ERROR", "The value must be greater than 0.", win)
			return
		}

		editCount_flt, err := strconv.ParseFloat(editCount.Text, 64)
		if err != nil || editCount_flt <= 0 {
			dialog.ShowInformation("ERROR", "The value must be greater than 0.", win)
			return
		}

		delta := editCount_flt - currentCount

		er := DB.Model(&TokenHistory{}).Where("Id = ?", histories[id].Id).Updates(TokenHistory{
			Date:      editDate.Text,
			Buy_Price: editPrice_flt,
			Count:     editCount_flt,
		}).Error

		if er != nil {
			fmt.Println("Error updating record:", err)
			return
		}

		if delta != 0 {
			err = DB.Model(&Infos{}).Where("Name = ?", histories[id].Name).Update("Count", gorm.Expr("Count + ?", delta)).Error
			if err != nil {
				fmt.Println("Error updating total token count in Infos:", err)
				return
			}
		}

		err = DB.Find(&tokens).Error
		if err != nil {
			fmt.Println("Error retrieving tokens:", err)
			return
		}

		setContent(win, w_content)
		token_list.Refresh()
		token_list.UnselectAll()
	})

	editSave.Resize(fyne.NewSize(200, 50))
	editSave.Move(fyne.NewPos(300, 500))

	back_btn := widget.NewButton("BACK", func() {
		setContent(win, w_content)
		token_list.UnselectAll()
	})
	back_btn.Resize(fyne.NewSize(50, 35))
	back_btn.Move(fyne.NewPos(730, 1))

	delete_btn = widget.NewButton(
		"DELETE", func() {
			dialog.ShowConfirm("DELETING TOKEN", fmt.Sprintf("TOCHNO DELET \"%s\"?", histories[id].Name),
				func(b bool) {
					if b {
						if err := DB.Delete(&TokenHistory{}, "Id = ?", histories[id].Id).Error; err != nil {
							fmt.Println("Ошибка при удалении токена из истории:", err)
							return
						}

						var count int64
						if err := DB.Model(&TokenHistory{}).Where("Name = ?", histories[id].Name).Count(&count).Error; err != nil {
							fmt.Println("Ошибка при подсчете токенов в истории:", err)
							return
						}

						if count == 0 {
							if err := DB.Delete(&Infos{}, "Name = ?", histories[id].Name).Error; err != nil {
								fmt.Println("Ошибка при удалении токена из Infos:", err)
								return
							}
						} else {
							var floatCount float64
							if err := DB.Model(&TokenHistory{}).Where("Name = ?", histories[id].Name).Select("SUM(Count)").Scan(&floatCount).Error; err != nil {
								fmt.Println("Ошибка при подсчете токенов в истории:", err)
								return
							}
							if err := DB.Model(&Infos{}).Where("Name = ?", histories[id].Name).Update("Count", floatCount).Error; err != nil {
								fmt.Println("Ошибка при обновлении количества токенов в Infos:", err)
								return
							}
						}
						DB.Find(&tokens)
						token_list.Refresh()
						setContent(win, w_content)
					}
				}, win)
		},
	)

	delete_btn.Move(fyne.NewPos(655, 1))
	delete_btn.Resize(fyne.NewSize(70, 35))

	editBar := container.NewWithoutLayout(
		setText(fmt.Sprintf("EDITING "+histories[id].Name), 25),
		back_btn)

	edit_cont := container.NewVBox(
		editBar,
		newLine(3, 0, 40, 800, 40),
		editDate,
		editPrice,
		editCount,
	)

	edit_cont.Move(fyne.NewPos(0, 0))
	edit_cont.Resize(fyne.NewSize(790, 500))

	editContainer := container.NewWithoutLayout(
		edit_cont,
		newLine(3, 0, 40, 800, 40),
		editSave,
		delete_btn,
	)
	setContent(win, editContainer)
}
