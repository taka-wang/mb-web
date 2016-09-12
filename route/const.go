package route

const (
	keyTCPDefaultPort          = "psmbtcp.default_port"
	keyMinConnectionTimout     = "psmbtcp.min_connection_timeout"
	keyPollInterval            = "psmbtcp.min_poll_interval"
	keyWebPrefix               = "route.prefix"
	keyWebIP                   = "route.ip"
	keyWebPort                 = "route.port"
	keyWebVersion              = "route.version"
	keyWebIAM                  = "route.iam"
	defaultTCPDefaultPort      = "502"
	defaultMinConnectionTimout = 200000
	defaultPollInterval        = 1
	defaultWebPrefix           = "/api/"
	defaultWebIP               = ""
	defaultWebPort             = "8080"
	defaultWebVersion          = "0.0.5"
	defaultWebIAM              = "web"
)
