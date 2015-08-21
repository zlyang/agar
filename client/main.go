package main

import (
	"encoding/json"
	"strconv"
	"strings"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/gl"

	"github.com/busyStone/agar/client/internal/ui"
	"github.com/busyStone/agar/conn"
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

	wsConn          = conn.DefaultConn
	wsConnWriteChan = make(chan []byte, 256)
	player          conn.S2CSelfInfo
	mClients        = make(map[string]conn.Logic)

	paintPlayers []ui.Player
	sz           size.Event
)

func main() {
	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch e := app.Filter(e).(type) {
			case lifecycle.Event:
				if err := appLifecycle(e); err != nil {
					return
				}
			case size.Event:
				sz = e
			case paint.Event:
				players := paintPlayers

				ui.OnPaint(players)
				debug.DrawFPS(sz)

				a.EndPaint(e)
			case touch.Event:

				action := conn.C2SAction{
					Type: conn.CUActionType,
					ID:   player.ID,
				}

				if e.Y > touchY {
					action.Action = "D"
				} else if e.Y < touchY {
					action.Action = "U"
				}

				if e.X > touchX {
					action.Action += "R"
				} else if e.X < touchX {
					action.Action += "L"
				}
				if action.Action != "" {
					data, err := json.Marshal(action)
					if err == nil {
						wsConnWriteChan <- data
					}
				}

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

		// 启动 读取 goroutine
		go wsConn.ReadPump(handleMessage, nil)
		go wsConn.WritePump(wsConnWriteChan)

	case lifecycle.CrossOff:
		// 关闭 读取 goroutine

		ui.OnStop()
	}

	return nil
}

func onStart() error {
	err := wsConn.Client()
	if err != nil {
		return err
	}

	ui.OnStart()

	return nil
}

// MessageType 用于确定消息类型
type MessageType struct {
	Type int
}

func handleMessage(message []byte) {
	var mt MessageType
	if err := json.Unmarshal(message, &mt); err != nil {
		return
	}

	switch mt.Type {
	case conn.CDAllClientsType: // 所有用户的坐标信息返回，其中第一个是分配给客户端的帐号
		fallthrough
	case conn.CDUpdateClientsType: // 有更新的用户坐标信息返回 有可能有新用户
		info := conn.S2CClientInfo{}
		if err := json.Unmarshal(message, &info); err != nil {
			return
		}
		for i := 0; i < len(info.Clients); i++ {
			mClients[info.Clients[i].Name] = info.Clients[i]
		}
	case conn.CDDeleteClientType: // 有用户断线
		info := conn.S2CDeleteClient{}
		if err := json.Unmarshal(message, &info); err != nil {
			return
		}
		if _, ok := mClients[info.Name]; ok {
			delete(mClients, info.Name)
		}
	case conn.CDSelfClientType: // 连接时返回自己的信息
		if err := json.Unmarshal(message, &player); err != nil {
			return
		}
	default:
		return
	}

	// 构造 players
	var players []ui.Player
	for _, v := range mClients {
		players = append(players, newPlayer(&v, sz))
	}

	paintPlayers = players
}

func newPlayer(client *conn.Logic, sz size.Event) ui.Player {
	player := ui.Player{}

	R, G, B := paseColor(client.Color)

	player.Color.R = R
	player.Color.G = G
	player.Color.B = B

	player.Pos.X = float32(client.Position.X) / float32(sz.WidthPx)
	player.Pos.Y = float32(client.Position.Y) / float32(sz.HeightPx)

	return player
}

func paseColor(color string) (R, G, B float32) {
	if color == "" {
		return
	}

	intColor, err := strconv.ParseInt(strings.TrimPrefix(color, "#"), 16, 32)
	if err != nil {
		return
	}

	R = float32(intColor>>16) / float32(255)
	G = float32((intColor&0x0FF00)>>8) / float32(255)
	B = float32(intColor&0x0FF) / float32(255)
	return
}
