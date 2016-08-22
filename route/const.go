package route

// environment variable names
const (
	envConfWeb         = "CONF_WEB"   // config path
	envBackendEndpoint = "EP_BACKEND" // backend endpoint
)

// config
const (
	defaultConfigPath  = "/etc/mb-web" // environment variable backup
	defaultBackendName = "consul"      // remote backend name
	keyConfigName      = "config"      // config file name
	keyConfigType      = "toml"        // config file extension
)

// logs
const (
	keyLogEnableDebug  = "log.debug"    // enable debug flag
	keyLogToJSONFormat = "log.json"     // log to json format flag
	keyLogToFile       = "log.to_file"  // log to file flag
	keyLogFileName     = "log.filename" // log filename

	defaultLogEnableDebug  = true
	defaultLogToJSONFormat = false
	defaultLogToFile       = false
	defaultLogFileName     = "/var/log/mb-web.log"
)

// [psmbtcp]
const (
	keyTCPDefaultPort      = "psmbtcp.default_port"
	keyMinConnectionTimout = "psmbtcp.min_connection_timeout"
	keyPollInterval        = "psmbtcp.min_poll_interval"
	keyMaxWorker           = "psmbtcp.max_worker"
	keyMaxQueue            = "psmbtcp.max_queue"

	defaultTCPDefaultPort      = "502"
	defaultMinConnectionTimout = 200000
	defaultPollInterval        = 1
	defaultMaxWorker           = 6
	defaultMaxQueue            = 100
)

// [web]
const (
	keyWebPub    = "web.pub"
	keyWebSub    = "web.sub"
	keyWebPrefix = "web.prefix"
	keyWebIP     = "web.ip"
	keyWebPort   = "web.port"

	defaultWebPub    = "ipc:///tmp/to.psmb"
	defaultWebSub    = "ipc:///tmp/from.psmb"
	defaultWebPrefix = "/api/"
	defaultWebIP     = ""
	defaultWebPort   = "8080"
)
