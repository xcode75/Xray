package api


// API config
type Config struct {
	APIHost      string  `mapstructure:"ApiHost"`
	NodeID       int     `mapstructure:"NodeID"`
	Key          string  `mapstructure:"ApiKey"`
	Timeout      int     `mapstructure:"Timeout"`
	SpeedLimit   float64 `mapstructure:"SpeedLimit"`
	DeviceLimit  int     `mapstructure:"DeviceLimit"`
	RuleListPath string  `mapstructure:"RuleListPath"`
}

// Node status
type NodeStatus struct {
	CPU    float64
	Mem    float64
	Disk   float64
	Uptime int
}

type NodeInfo struct {
	NodeType          string // Must be V2ray, Trojan, and Shadowsocks
	NodeID            int
	Port              int
	SpeedLimit        uint64 // Bps
	AlterID           int
	TransportProtocol string
	Host              string
	Path              string
	EnableTLS         bool
	TLSType           string
	CypherMethod      string
	ServiceName       string
	HeaderType        string
	AllowInsecure     bool
	Relay			  int
	ListenIP          string
	ProxyProtocol     bool
	Sniffing          bool
}

type RelayNodeInfo struct {
	NodeType          string // Must be V2ray, Trojan, and Shadowsocks
	NodeID            int
	Port              int
	SpeedLimit        uint64 // Bps
	AlterID           int
	TransportProtocol string
	Host              string
	Path              string
	EnableTLS         bool
	TLSType           string
	CypherMethod      string
	ServiceName       string
	HeaderType        string
	AllowInsecure     bool
	Address           string
	Relay			  int
	ListenIP          string
	ProxyProtocol     bool
	Sniffing          bool
	Flow              string
}

type UserInfo struct {
	UID           int
	Email         string
	Passwd        string
	Port          int
	SpeedLimit    uint64 
	DeviceLimit   int
	UUID          string
}

type OnlineUser struct {
	UID int
	IP  string
}

type UserTraffic struct {
	UID      int
	Email    string
	Upload   int64
	Download int64
}

type ClientInfo struct {
	APIHost  string
	NodeID   int
	Key      string
}

type DetectRule struct {
	ID      int
	Pattern string
}

type DetectResult struct {
	UID    int
	RuleID int
}
