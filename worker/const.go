package worker

// dataSource
const (
	_ dataSource = iota
	upstream
	downstream
)

// [worker]
const (
	statusOK          = "ok"
	keyWebPub         = "worker.pub"
	keyWebSub         = "worker.sub"
	keyResTimeout     = "worker.timeout"
	keyMaxWorker      = "worker.max_worker"
	keyMaxQueue       = "worker.max_queue"
	defaultResTimeout = 1000
	defaultWebPub     = "ipc:///tmp/to.psmb"
	defaultWebSub     = "ipc:///tmp/from.psmb"
	defaultMaxWorker  = 6
	defaultMaxQueue   = 100
)
