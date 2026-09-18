package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/alexiusacademia/fynesimplechart"
	"github.com/mrscorpio/fyneview/pkg/natscl"
	"github.com/mrscorpio/uahelper/pkg/tagdata"
)

//go:embed assets/natserv
var natsServAddr string

//go:embed assets/img/base.png
var baseImg []byte

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	d := new(tagdata.AllTags)

	nc, err := natscl.NewNats(natsServAddr)
	if err != nil {
		log.Println(err)
	}
	mes := widget.NewLabel(fmt.Sprint(err))

	err = nc.SendCmd(66, d)
	if err != nil {
		log.Println(err)
	}

	hello := widget.NewLabel("")
	hello.Alignment = fyne.TextAlignCenter
	output := canvas.NewText("время", color.NRGBA{B: 0xCC, G: 0xCC, R: 0xCC, A: 0xDF})

	base := canvas.NewImageFromResource(fyne.NewStaticResource("base.png", baseImg))
	base.FillMode = canvas.ImageFillContain
	base.SetMinSize(fyne.NewSize(600, 500))

	// bad simple chart
	nodes := []fynesimplechart.Node{}
	plot := fynesimplechart.NewPlot(nodes, "online")
	plot.ShowLine = true
	plot.LineWidth = 2
	chart := fynesimplechart.NewGraphWidget([]fynesimplechart.Plot{*plot})
	chart.SetChartTitle("rt")
	chart.Resize(fyne.NewSize(400, 300))

	ai := make(map[string]*Analog)
	aiList := []string{"PT908", "AT901", "TT903", "FT901", "TE521"}
	for _, v := range aiList {
		ai[v] = NewAnalog()
	}

	bar := make(map[string]*widget.ProgressBar)
	barList := []string{"ST50", "FV401_2", "FC424"}
	for _, v := range barList {
		bar[v] = widget.NewProgressBar()
	}

	vbar := make(map[string]*VertBar)
	vList := []string{"VT507", "VT508", "VT509"}
	for _, v := range vList {
		vbar[v] = NewVertBar()
	}
	aindex := make(map[string]int)

	go func() {
		ticker := time.NewTicker(time.Second)
		subscribed := false
		for range ticker.C {
			if nc == nil {
				continue
			}
			if len(nc.OnlineBuf) == 0 {
				continue
			}
			if !subscribed {
				log.Println(len(d.Tag[1].V))
				for _, v := range d.Tag[1].V {
					nodes = append(nodes, *fynesimplechart.NewNode(float32(len(nodes)), v))
				}
				err = nc.GetCurrent()
				if err != nil {
					log.Println(err)
				}
				subscribed = true
				for k := range ai {
					ai[k].SetName(k, d)
				}
				for k := range vbar {
					vbar[k].SetName(k, d)
				}
				for k := range bar {
					for i, v := range d.Tag {
						if k == v.Name {
							aindex[k] = i
							break
						}
					}
					if k == "ST50" {
						bar[k].Max = 30000
					}

					bar[k].TextFormatter = func() string {
						return fmt.Sprintf("%.0f            %.2f %s            %.0f", bar[k].Min, nc.OnlineBuf[aindex[k]], d.Tag[aindex[k]].Unit, bar[k].Max)
					}
					bar[k].Refresh()
				}

				continue
			}

			for i := range nc.OnlineBuf {
				d.AddV(i, nc.OnlineBuf[i])
			}
			d.AddT(nc.TimeBuf, true)

			fyne.Do(func() {
				for k := range ai {
					ai[k].UpdValue(nc.OnlineBuf)
				}
				for k := range vbar {
					vbar[k].UpdValue(nc.OnlineBuf)
				}
				output.Refresh()
				for k := range bar {
					bar[k].Value = float64(nc.OnlineBuf[aindex[k]])
					bar[k].Refresh()
				}

				//myBind.Reload()

				//output.TextSize = w.Content().Size().Width / 20
				//bar.SetValue(float64(nc.OnlineBuf[1] / 100))
				//nodes = append(nodes, *fynesimplechart.NewNode(float32(len(nodes)), nc.OnlineBuf[1]))
				//chart.Plots[0].Nodes = nodes
				//chart.Refresh()
			})

		}
	}()

	vibr := container.NewHBox(
		vbar["VT507"],
		vbar["VT508"],
		vbar["VT509"],
	)

	barNames := container.NewVBox(
		widget.NewLabel("обороты"),
		widget.NewLabel("гашетка"),
		widget.NewLabel("расход"),
	)

	bars := container.NewVBox(
		bar["ST50"],
		bar["FV401_2"],
		bar["FC424"],
	)

	dwn := container.NewBorder(container.NewPadded(vibr),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Hi!", func() {
				formatted := time.Now().Format("Time: 03:04:05")
				hello.SetText(formatted)
				if err != nil {
					log.Println(err)
				}
			}),
			layout.NewSpacer(),
		),
		barNames, nil, bars)

	//wf := container.NewGridWithRows(2, container.NewPadded(vibr), dwn)

	line1 := container.NewGridWithColumns(6, ai["PT908"], ai["AT901"], ai["TT903"])
	line2 := container.NewGridWithColumns(6, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), ai["TE521"])
	line3 := container.NewGridWithColumns(6, layout.NewSpacer(), layout.NewSpacer(), ai["FT901"], layout.NewSpacer())
	grid := container.NewGridWithRows(11, layout.NewSpacer(), line1, line2, layout.NewSpacer(), line3)

	center := container.NewStack(base, grid)

	main := container.NewBorder(mes, dwn, nil, nil, center)

	w.SetContent(main)

	w.ShowAndRun()
}
