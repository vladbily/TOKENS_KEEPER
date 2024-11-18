package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func setContent(win fyne.Window, currentContainer *fyne.Container) {
	win.SetContent(container.New(layout.NewMaxLayout(), back_img, currentContainer))
}

func setText(text string, size int8) fyne.CanvasObject {
	txt := canvas.NewText(text, color.Black)
	txt.TextSize = float32(size)
	txt.TextStyle.Bold = true
	txt.TextStyle.Italic = true
	return txt
}

func setPos(text fyne.CanvasObject, x float32, y float32) {
	text.Move(fyne.NewPos(x, y))
}

func newLine(size int, x1 float32, y1 float32, x2 float32, y2 float32) fyne.CanvasObject {
	line := canvas.NewLine(color.Black)
	line.StrokeWidth = float32(size)
	line.Position1 = fyne.NewPos(x1, y1)
	line.Position2 = fyne.NewPos(x2, y2)
	return line
}

/*func checkDif(sum float64) fyne.CanvasObject {
	if sum > 0 {
		total_dif := canvas.NewText(fmt.Sprintf("%.4f", sum), color.RGBA{54, 235, 48, 255})
		total_dif.TextSize = 20
		total_dif.TextStyle.Italic = true
		total_dif.TextStyle.Bold = true
		total_dif.Move(fyne.NewPos(630, 1))
		return total_dif
	} else {
		total_dif := canvas.NewText(fmt.Sprintf("%.4f", sum), color.RGBA{228, 62, 36, 255})
		total_dif.TextSize = 20
		total_dif.TextStyle.Italic = true
		total_dif.TextStyle.Bold = true
		total_dif.Move(fyne.NewPos(630, 1))
		return total_dif
	}
}*/

func checkTotal(sum float64) *canvas.Text {
	total := canvas.NewText(fmt.Sprintf("%.2f", sum), color.RGBA{238, 228, 36, 255})
	total.TextSize = 20
	total.TextStyle.Italic = true
	total.TextStyle.Bold = true
	total.Move(fyne.NewPos(370, 1))
	return total
}
