// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"os"
	"syscall"
	"time"
	"unsafe"

	efx "github.com/ohkashi/efxMfc-go"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type CompositeIndexItem struct {
	time    string
	value   float64
	change  float64
	sign    int
	chgRate float64
	volume  int
}

type MyMainWindow struct {
	*walk.MainWindow
	paintWidget *walk.CustomWidget
}

var main_wnd = &MyMainWindow{}
var kospi_data = [400]CompositeIndexItem{}
var kospi_data_len = 0
var kospi_min_value, kospi_max_value float64 = 99999, 0

func init() {
	start := time.Now()
	e1 := efx.Init("", OnRecvSysMsg)
	if e1 != nil {
		fmt.Printf("efx.Init: %s\n", e1)
		efx.Exit()
		os.Exit(1)
	}
	fmt.Printf("efxInit() elapsed time: %v\n\n", time.Since(start))
}

func main() {
	defer efx.Exit()
	efx1 := efx.NewControl(1, OnRecvData, nil, OnRecvError)
	efx.Synchonize(efx.REQ_LIMIT_MS)
	efx.SetSingleData(efx1, 0, "U")
	efx.SetSingleData(efx1, 1, "0001")
	efx.SetSingleData(efx1, 2, "60")
	efx.RequestData(efx1, "PUP02100200", "0001")
	//wait_recv_data()

	main_wnd = new(MyMainWindow)
	MainWindow{
		AssignTo: &main_wnd.MainWindow,
		Title:    "Walk: chart test",
		MinSize:  Size{320, 240},
		Size:     Size{800, 350},
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			CustomWidget{
				AssignTo:            &main_wnd.paintWidget,
				ClearsBackground:    false,
				InvalidatesOnResize: true,
				Paint:               main_wnd.drawStuff,
			},
		},
	}.Create()

	regKey, _ := registry.OpenKey(registry.CURRENT_USER,
		"Software\\Microsoft\\Windows\\CurrentVersion\\Themes\\Personalize", windows.KEY_READ)
	if regKey != 0 {
		isLightTheme, _, _ := regKey.GetIntegerValue("AppsUseLightTheme")
		regKey.Close()
		if isLightTheme == 0 {
			var useImmersiveDarkMode uint32 = 1
			windows.DwmSetWindowAttribute(windows.HWND(main_wnd.Handle()), windows.DWMWA_USE_IMMERSIVE_DARK_MODE,
				unsafe.Pointer(&useImmersiveDarkMode), uint32(unsafe.Sizeof(useImmersiveDarkMode)))
		}
	}
	main_wnd.Run()
}

var _sig_recv_data bool = false

func wait_recv_data() {
	for !_sig_recv_data {
		efx.ProcessMessage(1)
	}
	_sig_recv_data = false
}

func (mw *MyMainWindow) drawStuff(canvas *walk.Canvas, updateBounds walk.Rectangle) error {
	bounds := mw.paintWidget.ClientBoundsPixels()

	backBrush, err := walk.NewSolidColorBrush(walk.RGB(8, 8, 8))
	if err != nil {
		panic(err)
	}
	defer backBrush.Dispose()

	rectPen, err := walk.NewCosmeticPen(walk.PenSolid, walk.RGB(255, 0, 0))
	if err != nil {
		return err
	}
	defer rectPen.Dispose()

	/*if err := canvas.DrawRectangle(rectPen, bounds); err != nil {
		return err
	}*/
	if err := canvas.FillRectanglePixels(backBrush, bounds); err != nil {
		return err
	}

	linesBrush, err := walk.NewSolidColorBrush(walk.RGB(255, 8, 8))
	if err != nil {
		return err
	}
	defer linesBrush.Dispose()

	linesPen, err := walk.NewGeometricPen(walk.PenSolid, 1, linesBrush)
	if err != nil {
		return err
	}
	defer linesPen.Dispose()

	points := make([]walk.Point, kospi_data_len)
	dx := float64(bounds.Width) / float64(len(points)-1)
	kospi_cy := kospi_max_value - kospi_min_value
	cy := float64(bounds.Height)
	scale_factor := (cy - 20.0) / kospi_cy
	for i := range points {
		j := kospi_data_len - i - 1
		points[i].X = int(float64(i) * dx)
		points[i].Y = int(cy - (kospi_data[j].value-kospi_min_value)*scale_factor)
		if i > 0 {
			if err := canvas.DrawLinePixels(linesPen, points[i-1], points[i]); err != nil {
				return err
			}
		}
	}

	font, err := walk.NewFont("Consolas", 12, 0)
	if err != nil {
		return err
	}
	defer font.Dispose()

	bounds.X = 8
	bounds.X = 4
	if err := canvas.DrawTextPixels(fmt.Sprintf("kospi_data_len = %d\nkospi_cy = %.2f\nscale_factor = %.2f",
		kospi_data_len, kospi_cy, scale_factor),
		font, walk.RGB(8, 255, 8), bounds, walk.TextWordbreak); err != nil {
		return err
	}

	return nil
}

func OnRecvSysMsg(h syscall.Handle, ctrl_id uint, sys_msg int) {
	fmt.Printf("<== OnRecvSysMsg(%#x, %d, %d)\n", h, ctrl_id, sys_msg)
}

func OnRecvData(h syscall.Handle, ctrl_id uint, qry_name *string, param *string) {
	fmt.Printf("<-- OnRecvData(%#x, %d, \"%s\", \"%s\"): rq_id = %d\n", h, ctrl_id, *qry_name, *param, efx.GetRecvRqID(h))
	fmt.Printf("ReqMsgCode: %s, ReqMsg: \"%s\", RtCode: %s\n", efx.GetReqMsgCode(h), efx.GetReqMessage(h), efx.GetRtCode(h))

	switch *qry_name {
	case "PUP02100200":
		recCount := efx.GetMultiRecordCount(h, 0)
		var item CompositeIndexItem
		for i := 0; i < recCount; i++ {
			item.time = efx.GetMultiData(h, 0, i, 0)
			item.value = float64(efx.GetMultiDataFloat(h, 0, i, 1))
			item.change = float64(efx.GetMultiDataFloat(h, 0, i, 2))
			item.sign = efx.GetMultiDataInt(h, 0, i, 3)
			item.chgRate = float64(efx.GetMultiDataFloat(h, 0, i, 4))
			item.volume = efx.GetMultiDataInt(h, 0, i, 6)
			kospi_min_value = math.Min(kospi_min_value, item.value)
			kospi_max_value = math.Max(kospi_max_value, item.value)
			fmt.Printf("%s %.2f %.2f\n", item.time, item.value, item.change)
			kospi_data[kospi_data_len] = item
			kospi_data_len++
		}
		if efx.IsMoreNextData(h) {
			efx.Synchonize(efx.REQ_LIMIT_MS)
			efx.RequestNextData(h, *qry_name)
		} else {
			main_wnd.Synchronize(func() {
				main_wnd.SetTitle(fmt.Sprintf("KOSPI %.2f", kospi_data[0].value))
				main_wnd.Invalidate()
			})
		}
	default:
		fmt.Println("엥?")
	}
	fmt.Println()
	_sig_recv_data = true
}

func OnRecvError(_handle syscall.Handle, ctrl_id uint, qry_name *string, param *string) {
	fmt.Printf("<-- OnRecvError(%#x, %d, \"%s\", \"%s\")\n", _handle, ctrl_id, *qry_name, *param)
}
