package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/gl"

	"github.com/busyStone/agar/client/internal/socket"
	"github.com/busyStone/agar/client/internal/ui"
)

var (
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer

	green  float32
	touchX float32
	touchY float32

	conn *socket.Conn
)

func main() {
	app.Main(func(a app.App) {
		var sz size.Event
		for e := range a.Events() {
			switch e := app.Filter(e).(type) {
			case lifecycle.Event:
				appLifecycle(e)
			case size.Event:
				sz = e
			case paint.Event:

				debug.DrawFPS(sz)

				a.EndPaint(e)
			case touch.Event:
				touchX = e.X
				touchY = e.Y
			}
		}
	})
}

func appLifecycle(e lifecycle.Event) error {
	switch e.Crosses(lifecycle.StageVisible) {
	case lifecycle.CrossOn:
		if err := onStart(); err != nil {
			return err
		}
	case lifecycle.CrossOff:
		ui.OnStop()
	}

	return nil
}

func onStart() error {
	if conn == nil {
		var err error
		conn, err = socket.NewConnect()
		if err != nil {
			return err
		}
	}
	ui.OnStart()

	return nil
}

func newPlayer(touchX, touchY float32, sz size.Event) ui.Player {
	player := ui.Player{}

	player.Color.G = green
	player.Pos.X = touchX / float32(sz.WidthPx)
	player.Pos.Y = touchY / float32(sz.HeightPx)

	return player
}
