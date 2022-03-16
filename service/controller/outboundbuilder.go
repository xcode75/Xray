package controller

import (
	"encoding/json"
	"fmt"

	"github.com/xcode75/Xray/api"
	"github.com/xcode75/xraycore/common/net"
	"github.com/xcode75/xraycore/core"
	"github.com/xcode75/xraycore/infra/conf"
	
	"github.com/xcode75/xraycore/common/protocol"
	"github.com/xcode75/xraycore/common/serial"
    "github.com/xcode75/xraycore/proxy/vless"
)

//OutboundBuilder build freedom outbund config for addoutbound
func OutboundBuilder(config *Config, nodeInfo *api.NodeInfo, tag string) (*core.OutboundHandlerConfig, error) {
	outboundDetourConfig := &conf.OutboundDetourConfig{}
	outboundDetourConfig.Protocol = "freedom"
	outboundDetourConfig.Tag = tag

	// Build Send IP address
	if config.SendIP != "" {
		ipAddress := net.ParseAddress(config.SendIP)
		outboundDetourConfig.SendThrough = &conf.Address{ipAddress}
	}

	// Freedom Protocol setting
	var domainStrategy string = "Asis"
	if config.EnableDNS {
		if config.DNSType != "" {
			domainStrategy = config.DNSType
		} else {
			domainStrategy = "UseIP"
		}
	}
	proxySetting := &conf.FreedomConfig{
		DomainStrategy: domainStrategy,
	}
	
	if nodeInfo.NodeType == "dokodemo-door" {
		proxySetting.Redirect = fmt.Sprintf("0.0.0.0:%d", nodeInfo.Port-1)
	}
	
	var setting json.RawMessage
	setting, err := json.Marshal(proxySetting)
	if err != nil {
		return nil, fmt.Errorf("Marshal proxy %s config fialed: %s", nodeInfo.NodeType, err)
	}

	outboundDetourConfig.Settings = &setting
	return outboundDetourConfig.Build()
}

//OutboundBuilder build Blackhole outbund config for addoutbound
func BlackholeBuilder(config *Config) (*core.OutboundHandlerConfig, error) {
	outboundDetourConfig := &conf.OutboundDetourConfig{}
	outboundDetourConfig.Protocol = "blackhole"
	outboundDetourConfig.Tag = "block"

	// Build Send IP address
	if config.SendIP != "" {
		ipAddress := net.ParseAddress(config.SendIP)
		outboundDetourConfig.SendThrough = &conf.Address{ipAddress}
	}

	// Blackhole Protocol setting
	responses := make(map[string]string)
	responses["type"] = "http"
	var response json.RawMessage
	response, err := json.Marshal(responses)
	if err != nil {
			return nil, fmt.Errorf("Marshal Response Type %s into config fialed: %s", response, err)
	}
	
	proxySetting := &conf.BlackholeConfig{
		Response : response,
	}
	
	var setting json.RawMessage
	setting, err := json.Marshal(proxySetting)
	if err != nil {
		return nil, fmt.Errorf("Marshal blackhole settings into config fialed: %s", err)
	}
	outboundDetourConfig.Settings = &setting
	return outboundDetourConfig.Build()
}

type TrojanServerTarget struct {
	Address  string   `json:"address"`
	Port     uint16   `json:"port"`
	Password string   `json:"password"`
	Email    string   `json:"email"`
	Level    byte     `json:"level"`
	Flow     string   `json:"flow"`
}

type ShadowsocksServerTarget struct {
	Address  string `json:"address"`
	Port     uint16   `json:"port"`
	Cipher   string   `json:"method"`
	Password string   `json:"password"`
	Email    string   `json:"email"`
	Level    byte     `json:"level"`
}



//OutboundBuilder build relayoutbund config for addoutbound
func OutRelayboundBuilder(config *Config, nodeInfo *api.RelayNodeInfo , tag string, UUID string, Email string, Passwd string, UID int) (*core.OutboundHandlerConfig, error) {
		outboundDetourConfig := &conf.OutboundDetourConfig{}

		var (
			protocol      string
			streamSetting *conf.StreamConfig
			setting       json.RawMessage
		)

		var proxySetting interface{}

		if nodeInfo.NodeType == "Vless" {
				protocol = "vless"
				type VLessOutboundVnext struct {
					Address string           `json:"address"`
					Port    uint16            `json:"port"`
					Users   []json.RawMessage `json:"users"`
				}
				VlessUsers := buildRVlessUser(tag, nodeInfo , UUID, Email)
				VrawUsers := []json.RawMessage{}
				VrawUser,err := json.Marshal(&VlessUsers)
				if err != nil {
					return nil, fmt.Errorf("Marshal users %s config fialed: %s", VlessUsers, err)
				}
				VrawUsers = append(VrawUsers, VrawUser)
				proxySetting = struct {
					Vnext []*VLessOutboundVnext `json:"vnext"`
				}{
					Vnext: []*VLessOutboundVnext{&VLessOutboundVnext{
							Address: nodeInfo.Address,
							Port: uint16(nodeInfo.Port),
							Users: VrawUsers,
						},
					},
				}				
		}else if nodeInfo.NodeType == "Vmess" {
				protocol = "vmess"
				type VMessOutboundTarget struct {
					Address  string            `json:"address"`
					Port     uint16            `json:"port"`
					Users    []json.RawMessage `json:"users"`
				}
				users := buildRVmessUser(tag, UUID, Email, nodeInfo.AlterID)
				rawUsers := []json.RawMessage{}
				rawUser,err := json.Marshal(&users)
				if err != nil {
					return nil, fmt.Errorf("Marshal users %s config fialed: %s", users, err)
				}
				rawUsers = append(rawUsers, rawUser)
					
				proxySetting = struct {
					Receivers []*VMessOutboundTarget `json:"vnext"`
				}{
					Receivers: []*VMessOutboundTarget{&VMessOutboundTarget{
							Address: nodeInfo.Address,
							Port: uint16(nodeInfo.Port),
							Users: rawUsers,
						},
					},
				}				
		}else if nodeInfo.NodeType == "Trojan" {
				protocol = "trojan"	
				if nodeInfo.TLSType == "xtls" {
					proxySetting = struct {
						Servers []*TrojanServerTarget `json:"servers"`
					}{
						Servers: []*TrojanServerTarget{&TrojanServerTarget{
								Address: nodeInfo.Address,
								Port:     uint16(nodeInfo.Port),
								Password: UUID,
								Email:    fmt.Sprintf("%s_%s|%s", tag, Email, UUID),
								Level:    0,
								Flow:    nodeInfo.Flow,
							},
						},
					}
				} else {
					proxySetting = struct {
						Servers []*TrojanServerTarget `json:"servers"`
					}{
						Servers: []*TrojanServerTarget{&TrojanServerTarget{
								Address: nodeInfo.Address,
								Port:     uint16(nodeInfo.Port),
								Password: UUID,
								Email:    fmt.Sprintf("%s_%s|%s", tag, Email, UUID),
								Level:    0,
							},
						},
					}
				}
		}else if nodeInfo.NodeType == "Shadowsocks" {
				protocol = "shadowsocks"	
				proxySetting = struct {
					Servers []*ShadowsocksServerTarget `json:"servers"`
				}{
					Servers: []*ShadowsocksServerTarget{&ShadowsocksServerTarget{
							Address: nodeInfo.Address,
							Port:     uint16(nodeInfo.Port),
							Password: Passwd,
							Email:    fmt.Sprintf("%s_%s|%s", tag, Email, UUID),
							Level:    0,
							Cipher:   nodeInfo.CypherMethod,
						},
					},
				}
		}else{
			return nil, fmt.Errorf("Unsupported node type: %s, Only support: Vmess,Vless, Trojan, and Shadowsocks", nodeInfo.NodeType)
		}
		
		setting, err := json.Marshal(proxySetting)
		if err != nil {
			return nil, fmt.Errorf("Marshal proxy %s config fialed: %s", nodeInfo.NodeType, err)
		}
		
		streamSetting = new(conf.StreamConfig)
		transportProtocol := conf.TransportProtocol(nodeInfo.TransportProtocol)
		networkType, err := transportProtocol.Build()
		if err != nil {
			return nil, fmt.Errorf("convert TransportProtocol failed: %s", err)
		}
		
		if networkType == "tcp" {
			headers := make(map[string]string)
			headers["type"] = nodeInfo.HeaderType
			var header json.RawMessage
			header, err  := json.Marshal(headers)
			if err != nil {
				return nil, fmt.Errorf("Marshal Header Type %s into config fialed: %s", header, err)
			}		
			tcpSetting := &conf.TCPConfig{
				AcceptProxyProtocol: nodeInfo.ProxyProtocol,
				HeaderConfig:        header,
			}
			streamSetting.TCPSettings = tcpSetting
		} else if networkType == "websocket" {
			headers := make(map[string]string)
			headers["Host"] = nodeInfo.Host
			wsSettings := &conf.WebSocketConfig{
				AcceptProxyProtocol: nodeInfo.ProxyProtocol,
				Path:                nodeInfo.Path,
				Headers:             headers,
			}
			streamSetting.WSSettings = wsSettings
		} else if networkType == "http" {
			hosts := conf.StringList{nodeInfo.Host}
			httpSettings := &conf.HTTPConfig{
				Host: &hosts,
				Path: nodeInfo.Path,
			}
			streamSetting.HTTPSettings = httpSettings
		}else if networkType == "grpc" {
			grpcSettings := &conf.GRPCConfig{
				ServiceName: nodeInfo.ServiceName,
			}
			streamSetting.GRPCConfig = grpcSettings
		}
		
		streamSetting.Network = &transportProtocol
		
		if nodeInfo.EnableTLS{
			streamSetting.Security = nodeInfo.TLSType
			if nodeInfo.TLSType == "tls" {
				tlsSettings := &conf.TLSConfig{}
				tlsSettings.Insecure = true
				streamSetting.TLSSettings = tlsSettings
				
			} else if nodeInfo.TLSType == "xtls" {
				xtlsSettings := &conf.XTLSConfig{}
				xtlsSettings.Insecure = true
				streamSetting.XTLSSettings = xtlsSettings
			}
		}
		
		outboundDetourConfig.Tag = fmt.Sprintf("Relay_%s|%d", tag,UID)
		if config.SendIP != "" {
			ipAddress := net.ParseAddress(config.SendIP)
			outboundDetourConfig.SendThrough = &conf.Address{ipAddress}
		}
		outboundDetourConfig.Protocol = protocol
		outboundDetourConfig.StreamSetting = streamSetting
		outboundDetourConfig.Settings = &setting	
		return outboundDetourConfig.Build()
}


func buildRVmessUser(tag string, UUID string, Email string, serverAlterID int) *protocol.User {
		vmessAccount := &conf.VMessAccount{
			ID:       UUID,
			AlterIds: uint16(serverAlterID),
			Security: "auto",
		}
		return &protocol.User{
			Level:   0,
			Email:   fmt.Sprintf("%s_%s|%s", tag,Email, UUID), 
			Account: serial.ToTypedMessage(vmessAccount.Build()),
		}
}

func buildRVlessUser(tag string, nodeInfo *api.RelayNodeInfo , UUID string, Email string)  *protocol.User {
		vlessAccount := &vless.Account{
			Id:   UUID,
			Flow: nodeInfo.Flow,
			Encryption: "none",
		}
		return &protocol.User{
			Level:   0,
			Email:   fmt.Sprintf("%s_%s|%s", tag, Email, UUID),
			Account: serial.ToTypedMessage(vlessAccount),
		}
}

