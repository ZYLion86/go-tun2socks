package tun2socks

import (
	"context"
	"github.com/ZYLion86/go-tun2socks/v2ray"
	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
	v2rayCore "github.com/v2fly/v2ray-core"
	vproxyman "github.com/v2fly/v2ray-core/app/proxyman"
	"log"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var lwipStack core.LWIPStack

type PacketFlow interface {
	WritePacket(packet []byte)
}

func InputPacket(data []byte) {
	lwipStack.Write(data)
}

func StartSocks(packetFlow PacketFlow, proxyHost string, proxyPort int) {
	if packetFlow != nil {
		lwipStack = core.NewLWIPStack()
		core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, uint16(proxyPort)))
		core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, uint16(proxyPort), 30*time.Second))
		core.RegisterOutputFn(func(data []byte) (int, error) {
			packetFlow.WritePacket(data)
			return len(data), nil
		})
	}
}

func StartV2Ray(packetFlow PacketFlow, configBytes []byte) {
	if packetFlow == nil {
		return
	}
	lwipStack = core.NewLWIPStack()
	var v *v2rayCore.Instance
	v, err := v2rayCore.StartInstance("json", configBytes)
	if err != nil {
		log.Fatalf("start V instance failed: %v", err)
	}

	sniffingConfig := &vproxyman.SniffingConfig{
		Enabled:             true,
		DestinationOverride: strings.Split("tls,http", ","),
	}

	debug.SetGCPercent(5)
	ctx := vproxyman.ContextWithSniffingConfig(context.Background(), sniffingConfig)
	core.RegisterTCPConnHandler(v2ray.NewTCPHandler(ctx, v))
	core.RegisterUDPConnHandler(v2ray.NewUDPHandler(ctx, v, 30*time.Second))
	core.RegisterOutputFn(func(data []byte) (int, error) {
		packetFlow.WritePacket(data)
		runtime.GC()
		debug.FreeOSMemory()
		return len(data), nil
	})
}
