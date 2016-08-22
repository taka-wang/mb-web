package route

import (
	"net/http"
	"os"
	"path"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/go-zoo/bone"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var (
	// version version string
	version = "0.0.1"
	// Mux global mux
	Mux = bone.New()
	// service backend service
	//service Service
	// ip Web IP
	ip string
	// port Web port
	port string
	// pubEndpoint zmq pub endpoint
	pubEndpoint string
	// subEndpoint zmq sub endpoint
	subEndpoint string

	// defaultMbPort default modbus slave port number
	defaultMbPort string
	// minConnTimeout minimal modbus tcp connection timeout
	minConnTimeout int64
	// minPollInterval minimal modbus tcp poll interval
	minPollInterval uint64
	// maxQueueSize the size of job queue
	maxQueueSize int
	// max_workers the number of workers to start
	maxWorkers int
)

func init() {
	// before load config
	log.SetHandler(text.New(os.Stdout))
	log.SetLevel(log.DebugLevel)

	// setup config
	initConfig()
	setDefaults()
	setLogger()

	// get mbtcp default
	defaultMbPort = viper.GetString(keyTCPDefaultPort)
	minConnTimeout = viper.GetInt64(keyMinConnectionTimout)
	minPollInterval = uint64(viper.GetInt(keyPollInterval))
	maxWorkers = viper.GetInt(keyMaxWorker)
	maxQueueSize = viper.GetInt(keyMaxQueue)

	// get web defaults
	ip = viper.GetString(keyWebIP)
	port = viper.GetString(keyWebPort)
	Mux.Prefix(viper.GetString(keyWebPrefix))
	pubEndpoint = viper.GetString(keyWebPub)
	subEndpoint = viper.GetString(keyWebSub)
}

// Middleware set middleware
func Middleware() {
	/*
		s, _ := NewService()
		service = *s
		service.Start()
	*/
	// ========== general handlers ==========

	// 404 Not Found
	Mux.NotFoundFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("404 NOT FOUND!"))
	})

	// ROOT
	Mux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("version: " + version))
	}))

	// ========== one-off command ==========
	Mux.Get(mbOnceRead, http.HandlerFunc(handleMbOnceRead))
	Mux.Post(mbOnceWrite, http.HandlerFunc(handleMbOnceWrite))
	Mux.Get(mbGetTimeout, http.HandlerFunc(handleMbGetTimeout))
	Mux.Post(mbSetTimeout, http.HandlerFunc(handleMbSetTimeout))

	// ========== poll command ==========
	Mux.Get(mbGetPoll, http.HandlerFunc(handleMbGetPoll))
	Mux.Post(mbCreatePoll, http.HandlerFunc(handleMbCreatePoll))
	Mux.Put(mbUpdatePoll, http.HandlerFunc(handleMbUpdatePoll))
	Mux.Delete(mbDeletePoll, http.HandlerFunc(handleMbDeletePoll))

	Mux.Post(mbTogglePoll, http.HandlerFunc(handleMbTogglePoll))
	Mux.Get(mbGetPollHistory, http.HandlerFunc(handleMbGetPollHistory))

	Mux.Get(mbGetPolls, http.HandlerFunc(handleMbGetPolls))
	Mux.Delete(mbDeletePolls, http.HandlerFunc(handleMbDeletePolls))

	Mux.Post(mbTogglePolls, http.HandlerFunc(handleMbTogglePolls))
	Mux.Get(mbExportPolls, http.HandlerFunc(handleMbExportPolls))
	Mux.Post(mbImportPolls, http.HandlerFunc(handleMbImportPolls))

	// ========== filter command ==========
	Mux.Get(mbGetFilter, http.HandlerFunc(handleMbGetFilter))
	Mux.Post(mbCreateFilter, http.HandlerFunc(handleMbCreateFilter))
	Mux.Put(mbUpdateFilter, http.HandlerFunc(handleMbUpdateFilter))
	Mux.Delete(mbDeleteFilter, http.HandlerFunc(handleMbDeleteFilter))

	Mux.Post(mbToggleFilter, http.HandlerFunc(handleMbToggleFilter))

	Mux.Get(mbGetFilters, http.HandlerFunc(handleMbGetFilters))
	Mux.Delete(mbDeleteFilters, http.HandlerFunc(handleMbDeleteFilters))

	Mux.Post(mbToggleFilters, http.HandlerFunc(handleMbToggleFilters))
	Mux.Get(mbExportFilters, http.HandlerFunc(handleMbExportFilters))
	Mux.Post(mbImportFilters, http.HandlerFunc(handleMbImportFilters))

}

// Start server start
func Start() {
	http.ListenAndServe(ip+":"+port, Mux) // start server
}

//
// Internal
//

// initConfig int config function
func initConfig() {
	// get environment variables
	confPath := os.Getenv(envConfWeb) // config file location
	if confPath == "" {
		confPath = defaultConfigPath
	}
	endpoint := os.Getenv(envBackendEndpoint) // backend endpoint, i.e., consul url

	// setup config filename and extension
	viper.SetConfigName(keyConfigName)
	viper.SetConfigType(keyConfigType)

	// local or remote config
	if endpoint == "" {
		log.WithField("file path", confPath).Debug("Try to load 'local' config file")
		viper.AddConfigPath(confPath)
		err := viper.ReadInConfig() // read config from file
		if err != nil {
			log.WithField("file path", confPath).Warn("Fail to load 'local' config file, not found!")
		} else {
			log.WithField("file path", confPath).Info("Read 'local' config file successfully")
		}
	} else {
		log.WithFields(log.Fields{
			"endpoint":  endpoint,
			"file path": confPath,
			"file name": keyConfigName,
			"file type": keyConfigType,
		}).Debug("Try to load 'remote' config file")

		viper.AddRemoteProvider(defaultBackendName, endpoint, path.Join(confPath, keyConfigName)+"."+keyConfigType)
		err := viper.ReadRemoteConfig() // read config from backend
		if err != nil {
			log.WithFields(log.Fields{
				"err":       err,
				"endpoint":  endpoint,
				"file path": confPath,
				"file name": keyConfigName,
				"file type": keyConfigType,
			}).Error("Fail to load 'remote' config file, not found!")
		} else {
			log.WithFields(log.Fields{
				"endpoint":  endpoint,
				"file path": confPath,
				"file name": keyConfigName,
				"file type": keyConfigType,
			}).Info("Read 'remote' config file successfully")
		}
	}
}

// setDefaults set default values
func setDefaults() {
	// set default log values
	viper.SetDefault(keyLogEnableDebug, defaultLogEnableDebug)
	viper.SetDefault(keyLogToJSONFormat, defaultLogToJSONFormat)
	viper.SetDefault(keyLogToFile, defaultLogToFile)
	viper.SetDefault(keyLogFileName, defaultLogFileName)

	// set default psmbtcp values
	viper.SetDefault(keyTCPDefaultPort, defaultTCPDefaultPort)
	viper.SetDefault(keyMinConnectionTimout, defaultMinConnectionTimout)
	viper.SetDefault(keyPollInterval, defaultPollInterval)
	viper.SetDefault(keyMaxWorker, defaultMaxWorker)
	viper.SetDefault(keyMaxQueue, defaultMaxQueue)

	// set default web values
	viper.SetDefault(keyWebPub, defaultWebPub)
	viper.SetDefault(keyWebSub, defaultWebSub)
	viper.SetDefault(keyWebPrefix, defaultWebPrefix)
	viper.SetDefault(keyWebIP, defaultWebIP)
	viper.SetDefault(keyWebPort, defaultWebPort)
}

// setLogger init logger function
func setLogger() {

	writer := os.Stdout
	if viper.GetBool(keyLogToFile) {
		if f, err := os.OpenFile(viper.GetString(keyLogFileName), os.O_WRONLY|os.O_CREATE, 0755); err != nil {
			log.WithFields(log.Fields{
				"err":       err,
				"file name": viper.GetString(keyLogFileName),
			}).Error("Fail to create log file")

		} else {
			writer = f // to file
		}
	}

	// set log formatter, JSON or plain text
	if viper.GetBool(keyLogToJSONFormat) {
		log.SetHandler(json.New(writer))
	} else {
		log.SetHandler(text.New(writer))
	}

	// set debug level
	if viper.GetBool(keyLogEnableDebug) {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
