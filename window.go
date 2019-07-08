// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

/* Package mjpgmon is a simple gui to view mjpeg streams from url. */

package mjpgmonitor

import (
	"image"
	_ "image/jpeg"
	"time"

	"bitbucket.org/rj/goey"
	"bitbucket.org/rj/goey/base"
)

var mainWindow *goey.Window

type display struct {
	base.Widget
	img *goey.Img
	stream *streamDecoder
}

func (d *display) frameloop(targetFPS int) {
	var count int
	var fpmsecs = targetFPS / 1000
	for {
	fpsLoop:
		timeout := time.After(time.Duration(fpmsecs) * time.Millisecond)
		for ; count < targetFPS; count++ {
			select {
			case i := <-d.stream.i:
				d.img.Height = i.Bounds().Dy()
				d.img.Width = i.Bounds().Dx()
				d.img.UpdateDimensions()
				d.img.Image = i
				continue
			case <-timeout:
				goto fpsLoop

			}
		}
	}
}

func createWindow() error {
	win, err := goey.NewWindow("mpegmonitor", renderWindow())
	if err != nil {
		return err
	}
	mainWindow = win

	// TODO: set window icon
	return nil
}

func (d *display) renderWindow() base.Widget {
	w := goey.Padding{
		Insets: goey.DefaultInsets(),
		Child: &goey.VBox{
			AlignMain: goey.MainCenter,
			AlignCross: goey.CrossCenter,
			Children: []base.Widget{
				d.img,
				&goey.Button{Text: "close", OnClick: func() {
					close(d.stream.i)
					mainWindow.Close()
				}},

			}
		}
	}
	return nil
}

