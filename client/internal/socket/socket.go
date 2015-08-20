package socket

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	ws "github.com/gorilla/websocket"
)

const (
	_host         = "ws://127.0.0.1:8080"
	_connectURL   = _host + "/connect"
	_readBufSize  = 1024
	_writeBufSize = 1024
)

// Conn 封装 ws.Conn
type Conn struct {
	*ws.Conn
}

// NewConnect 创建新连接
func NewConnect() (*Conn, error) {
	u, err := url.Parse(_connectURL)
	if err != nil {
		return nil, err
	}

	rawConn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return nil, err
	}

	wsHeaders := http.Header{
		"Origin": {_host},
		// your milage may differ
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}

	wsConn, resp, err := ws.NewClient(rawConn, u, wsHeaders, _readBufSize, _writeBufSize)
	if err != nil {
		return nil, fmt.Errorf("websocket.NewClient Error: %s\nResp:%+v", err, resp)
	}

	return &Conn{wsConn}, nil
}
