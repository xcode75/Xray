package all

import (
	// The following are necessary as they register handlers in their init functions.

	// Required features. Can't remove unless there is replacements.
	// _ "github.com/xcode75/xraycore/app/dispatcher"
	_ "github.com/xcode75/xraycore/app/proxyman/inbound"
	_ "github.com/xcode75/xraycore/app/proxyman/outbound"

	// Default commander and all its services. This is an optional feature.
	_ "github.com/xcode75/xraycore/app/commander"
	_ "github.com/xcode75/xraycore/app/log/command"
	_ "github.com/xcode75/xraycore/app/proxyman/command"
	_ "github.com/xcode75/xraycore/app/stats/command"

	// Other optional features.
	_ "github.com/xcode75/xraycore/app/dns"
	_ "github.com/xcode75/xraycore/app/log"
	_ "github.com/xcode75/xraycore/app/policy"
	_ "github.com/xcode75/xraycore/app/reverse"
	_ "github.com/xcode75/xraycore/app/router"
	_ "github.com/xcode75/xraycore/app/stats"

	// Inbound and outbound proxies.
	_ "github.com/xcode75/xraycore/proxy/blackhole"
	_ "github.com/xcode75/xraycore/proxy/dns"
	_ "github.com/xcode75/xraycore/proxy/dokodemo"
	_ "github.com/xcode75/xraycore/proxy/freedom"
	_ "github.com/xcode75/xraycore/proxy/http"
	_ "github.com/xcode75/xraycore/proxy/mtproto"
	_ "github.com/xcode75/xraycore/proxy/shadowsocks"
	_ "github.com/xcode75/xraycore/proxy/socks"
	_ "github.com/xcode75/xraycore/proxy/trojan"
	_ "github.com/xcode75/xraycore/proxy/vless/inbound"
	_ "github.com/xcode75/xraycore/proxy/vless/outbound"
	_ "github.com/xcode75/xraycore/proxy/vmess/inbound"
	_ "github.com/xcode75/xraycore/proxy/vmess/outbound"

	// Transports
	_ "github.com/xcode75/xraycore/transport/internet/domainsocket"
	_ "github.com/xcode75/xraycore/transport/internet/http"
	_ "github.com/xcode75/xraycore/transport/internet/kcp"
	_ "github.com/xcode75/xraycore/transport/internet/quic"
	_ "github.com/xcode75/xraycore/transport/internet/tcp"
	_ "github.com/xcode75/xraycore/transport/internet/tls"
	_ "github.com/xcode75/xraycore/transport/internet/udp"
	_ "github.com/xcode75/xraycore/transport/internet/websocket"
	_ "github.com/xcode75/xraycore/transport/internet/xtls"

	// Transport headers
	_ "github.com/xcode75/xraycore/transport/internet/headers/http"
	_ "github.com/xcode75/xraycore/transport/internet/headers/noop"
	_ "github.com/xcode75/xraycore/transport/internet/headers/srtp"
	_ "github.com/xcode75/xraycore/transport/internet/headers/tls"
	_ "github.com/xcode75/xraycore/transport/internet/headers/utp"
	_ "github.com/xcode75/xraycore/transport/internet/headers/wechat"
	_ "github.com/xcode75/xraycore/transport/internet/headers/wireguard"

	// JSON & TOML & YAML
	_ "github.com/xcode75/xraycore/main/json"
	_ "github.com/xcode75/xraycore/main/toml"
	_ "github.com/xcode75/xraycore/main/yaml"

	// Load config from file or http(s)
	_ "github.com/xcode75/xraycore/main/confloader/external"

	// Commands
	_ "github.com/xcode75/xraycore/main/commands/all"
)
