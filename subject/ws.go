package subject

type Ws struct {
	Name      string    `json:"name"`
	ProxyType ProxyType `json:"type"`
	Username  string    `json:"username"`
	Method    string    `json:"method"`
	PasswdLen int       `json:"passwd-len"`
}

func (w Ws) ProxyName() string {
	return w.Name
}

func (w Ws) Type() ProxyType {
	return WS
}
