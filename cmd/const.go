package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

var back_img *canvas.Image

var w_content *fyne.Container
var create_window *fyne.Container
var s_content *fyne.Container

var token_list *widget.List

var save_btn *widget.Button
var save_set_btn *widget.Button
var delete_btn *widget.Button
var back_btn *widget.Button
var back_set_btn *widget.Button
