//go:build !ios && !android
// +build !ios,!android

package v2ray

import (
	_ "github.com/v2fly/v2ray-core/app/commander"
	_ "github.com/v2fly/v2ray-core/app/log/command"
	_ "github.com/v2fly/v2ray-core/app/proxyman/command"
	_ "github.com/v2fly/v2ray-core/app/stats/command"

	_ "github.com/v2fly/v2ray-core/app/reverse"

	_ "github.com/v2fly/v2ray-core/transport/internet/domainsocket"
)
