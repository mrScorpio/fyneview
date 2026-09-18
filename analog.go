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

type Analog struct {
	widget.BaseWidget
	Base  *canvas.Text
	Unit  *canvas.Text
	Hdr   *canvas.Text
	Plate *canvas.Rectangle
	Ind   int
}

func NewAnalog() *Analog {
	base := canvas.NewText("", color.White)
	base.Alignment = fyne.TextAlignTrailing
	hdr := canvas.NewText("", color.White)
	unit := canvas.NewText("", color.White)
	base.TextSize = hdr.TextSize * 1.3
	//plate := canvas.NewRectangle(color.NRGBA{B: 0x11, G: 0xCC, R: 0x11, A: 0x33})
	plate := canvas.NewRectangle(theme.Color(theme.ColorNameFocus))
	plate.CornerRadius = 6
	item := &Analog{Base: base, Hdr: hdr, Plate: plate, Unit: unit}
	item.ExtendBaseWidget(item)
	return item
}

func (a *Analog) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewVSplit(a.Hdr, container.NewHBox(layout.NewSpacer(), a.Base, a.Unit))
	c.SetOffset(0.2)
	m := container.NewStack(a.Plate, c)
	return widget.NewSimpleRenderer(m)
}

func (a *Analog) SetName(name string, d *tagdata.AllTags) error {
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

func (a *Analog) UpdValue(buf []float32) {
	a.Base.Text = fmt.Sprint(buf[a.Ind])
	a.Base.Refresh()
}
