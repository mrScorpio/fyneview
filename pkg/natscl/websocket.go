package natscl

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/coder/websocket"
)

type WasmNatsConnectionWrapper struct {
	ws *websocket.Conn
}

func (cw WasmNatsConnectionWrapper) Dial(network, address string) (net.Conn, error) {
	// we actually do not care about the adress given here
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	var err error
	cw.ws, _, err = websocket.Dial(ctx, "ws://"+address, nil)

	if err != nil {
		return nil, fmt.Errorf("websocket.Dial failed %w", err)
	}
	nconn := websocket.NetConn(context.Background(), cw.ws, websocket.MessageBinary)
	return nconn, nil
}

func (cw WasmNatsConnectionWrapper) SkipTLSHandshake() bool {
	return true
}
