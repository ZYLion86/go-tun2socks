package v2ray

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/eycorsican/go-tun2socks/common/log"
	"github.com/eycorsican/go-tun2socks/core"
	vcore "github.com/v2fly/v2ray-core"
	vsession "github.com/v2fly/v2ray-core/common/session"
	vnetlegay "v2ray.com/core/common/net"
)

type tcpHandler struct {
	ctx context.Context
	v   *vcore.Instance
}

func (h *tcpHandler) handleInput(conn net.Conn, input io.ReadCloser) {
	defer func() {
		conn.Close()
		input.Close()
	}()
	io.Copy(conn, input)
}

func (h *tcpHandler) handleOutput(conn net.Conn, output io.WriteCloser) {
	defer func() {
		conn.Close()
		output.Close()
	}()
	io.Copy(output, conn)
}

func NewTCPHandler(ctx context.Context, instance *vcore.Instance) core.TCPConnHandler {
	return &tcpHandler{
		ctx: ctx,
		v:   instance,
	}
}

func (h *tcpHandler) Handle(conn net.Conn, target *net.TCPAddr) error {
	dest := vnetlegay.DestinationFromAddr(target)
	sid := vsession.NewID()
	ctx := vsession.ContextWithID(h.ctx, sid)
	c, err := vcore.Dial(ctx, h.v, dest)
	if err != nil {
		return errors.New(fmt.Sprintf("dial V proxy connection failed: %v", err))
	}
	go h.handleInput(conn, c)
	go h.handleOutput(conn, c)
	log.Infof("new proxy connection for target: %s:%s", target.Network(), target.String())
	return nil
}
