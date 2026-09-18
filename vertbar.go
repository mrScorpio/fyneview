package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mrscorpio/uahelper/pkg/tagdata"
)

type VertBar struct {
	widget.BaseWidget
	Base  *canvas.Text
	Unit  *canvas.Text
	Hdr   *canvas.Text
	Plate *canvas.Rectangle
	Val   *canvas.Rectangle
	Ind   int
}

func NewVertBar() *VertBar {
	base := canvas.NewText("", color.White)
	base.Alignment = fyne.TextAlignTrailing
	hdr := canvas.NewText("", color.White)
	unit := canvas.NewText("", color.White)
	//base.TextSize = hdr.TextSize * 1.3
	//plate := canvas.NewRectangle(color.NRGBA{B: 0x11, G: 0xCC, R: 0x11, A: 0x33})
	plate := canvas.NewRectangle(theme.Color(theme.ColorNameFocus))
	//plate.CornerRadius = 6
	plate.SetMinSize(fyne.NewSize(66, 100))
	val := canvas.NewRectangle(theme.Color(theme.ColorNamePrimary))
	//val.CornerRadius = 6
	val.SetMinSize(fyne.NewSize(66, 0))
	item := &VertBar{Base: base, Hdr: hdr, Plate: plate, Unit: unit, Val: val}
	item.ExtendBaseWidget(item)
	return item
}

func (a *VertBar) CreateRenderer() fyne.WidgetRenderer {
	v := container.NewVBox(a.Base, layout.NewSpacer(), a.Val)
	m := container.NewStack(a.Plate, v)
	return widget.NewSimpleRenderer(m)
}

func (a *VertBar) SetName(name string, d *tagdata.AllTags) error {
	for i := range d.Tag {
		if name == d.Tag[i].Name {
			a.Ind = i
			a.Hdr.Text = name
			a.Unit.Text = d.Tag[i].Unit
			return nil
		}
	}
	a.Hdr.Text = "no tag"
	return fmt.Errorf("tag not found")
}

func (a *VertBar) UpdValue(buf []float32) {
	a.Base.Text = fmt.Sprint(buf[a.Ind])
	a.Val.SetMinSize(fyne.NewSize(a.Val.Size().Width, buf[a.Ind]*100))
	a.Val.Refresh()
}
