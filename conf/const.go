package conf

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
