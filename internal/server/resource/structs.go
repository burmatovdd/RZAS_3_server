package resource

type Info struct {
	Domain string `json:"domain"`
	Waf    bool   `json:"waf"`
	SSL    bool   `json:"SSL"`
	Active bool   `json:"active"`
	Owner  string `json:"owner"`
}

type Resource struct {
	Domain string `json:"domain"`
}
