package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/gl"

	"github.com/busyStone/agar/client/internal/ui"
)

var (
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer

	green             float32
	touchXOld, touchX float32
	touchYOld, touchY float32
)

func main() {
	app.Main(func(a app.App) {
		var sz size.Event
		for e := range a.Events() {
			switch e := app.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					ui.OnStart()
				case lifecycle.CrossOff:
					ui.OnStop()
				}
			case size.Event:
				sz = e
				touchX = float32(sz.WidthPx / 2)
				touchY = float32(sz.HeightPx / 2)
			case paint.Event:

				green += 0.01
				if green > 1 {
					green = 0
				}

				player1 := newPlayer(touchX, touchY, sz)
				player2 := newPlayer(touchX+float32(sz.WidthPx/10), touchY+float32(sz.HeightPx/10), sz)
				ui.OnPaint([]ui.Player{player1, player2})

				debug.DrawFPS(sz)

				a.EndPaint(e)
			case touch.Event:
				touchX = e.X
				touchY = e.Y
			}
		}
	})
}

func newPlayer(touchX, touchY float32, sz size.Event) ui.Player {
	player := ui.Player{}

	player.Color.G = green
	player.Pos.X = touchX / float32(sz.WidthPx)
	player.Pos.Y = touchY / float32(sz.HeightPx)

	return player
}
