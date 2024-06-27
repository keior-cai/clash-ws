package subject

const (
	WS ProxyType = iota
)

type Subject struct {
	Proxies     []map[string]any `yaml:"proxies"`
	ProxyGroups []any            `yaml:"proxy-groups"`
	Rules       []string         `yaml:"rules"`
}

type ProxyType int

type Proxies interface {
	ProxyName() string
}
