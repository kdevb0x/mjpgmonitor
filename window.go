// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

/* Package mjpgmon is a simple gui to view mjpeg streams from url. */

package main

import (
	_ "image/jpeg"
	"time"

	"bitbucket.org/rj/goey"
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/loop"
)

var mainWindow *goey.Window

type display struct {
	base.Widget
	img    *goey.Img
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
				d.img.Height = base.Length(int32(i.Bounds().Dy()))
				d.img.Width = base.Length(int32(i.Bounds().Dx()))
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
	decoder, err := NewStreamDecoder(remoteURL, *remoteUsername, *remotePassword)
	if err != nil {
		return err
	}

	d := &display{stream: decoder}
	win, err := goey.NewWindow("mpegmonitor", spawnGUI())
	if err != nil {
		return err
	}
	mainWindow = win
	return loop.Do(func() error { return d.renderWindow() })

	// TODO: set window icon

}

func (d *display) renderWindow() error {
	w := goey.Padding{
		Insets: goey.DefaultInsets(),
		Child: &goey.VBox{
			AlignMain:  goey.MainCenter,
			AlignCross: goey.CrossCenter,
			Children: []base.Widget{
				d.img,
				&goey.Button{Text: "close", OnClick: func() {
					close(d.stream.i)
					mainWindow.Close()
				}},
			},
		},
	}
	mainWindow.SetChild(w)
	return nil
}

func spawnGUI() base.Widget {
	// TODO: spawn the gui elements and pass this to main window constructor
}
