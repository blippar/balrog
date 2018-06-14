package browser

// HTTPConfig defines a standard HTTP server configuration
type HTTPConfig struct {
	Addr      string `json:"addr"`
	Templates string `json:"templates"`
	Dist      string `json:"dist"`
}

// SetDefaultConfig defines the default config for an HTTP server
func (h *HTTPConfig) SetDefaultConfig() {
	h.Addr = ":8000"
	h.Templates = "templates/"
	h.Dist = "dist/"
}
