package api

import "github.com/xcode75/xraycore/app/router"

type API interface {
	GetNodeInfo() (nodeInfo *NodeInfo, err error)
	GetRelayNodeInfo() (relaynodeInfo *RelayNodeInfo, err error)
	GetUserList() (userList *[]UserInfo, err error)
	GetRouteInfo() (routeConfig *router.Config, err error)
	ReportNodeStatus(nodeStatus *NodeStatus) (err error)
	ReportNodeOnlineUsers(onlineUser *[]OnlineUser) (err error)
	ReportUserTraffic(userTraffic *[]UserTraffic) (err error)
	Describe() ClientInfo
	GetNodeRule() (ruleList *[]DetectRule, err error)
	ReportIllegal(detectResultList *[]DetectResult) (err error)
	Debug()
}
