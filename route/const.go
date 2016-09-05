package route

const (
	// version version string
	version = "0.0.1"

	// [psmbtcp]
	keyTCPDefaultPort          = "psmbtcp.default_port"
	keyMinConnectionTimout     = "psmbtcp.min_connection_timeout"
	keyPollInterval            = "psmbtcp.min_poll_interval"
	defaultTCPDefaultPort      = "502"
	defaultMinConnectionTimout = 200000
	defaultPollInterval        = 1

	// [route]
	keyWebPrefix     = "route.prefix"
	keyWebIP         = "route.ip"
	keyWebPort       = "route.port"
	defaultWebPrefix = "/api/"
	defaultWebIP     = ""
	defaultWebPort   = "8080"
	iam              = "web" // who am I
)
