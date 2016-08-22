package route

// path
const (
	mbOnceRead   = "/mb/tcp/fc/:fc"
	mbOnceWrite  = "/mb/tcp/fc/:fc"
	mbGetTimeout = "/mb/tcp/timeout"
	mbSetTimeout = "/mb/tcp/timeout"

	mbCreatePoll = "/mb/tcp/poll/:name"
	mbUpdatePoll = "/mb/tcp/poll/:name"
	mbGetPoll    = "/mb/tcp/poll/:name"
	mbDeletePoll = "/mb/tcp/poll/:name"

	mbTogglePoll     = "/mb/tcp/poll/:name/toggle"
	mbGetPollHistory = "/mb/tcp/poll/:name/history"

	mbGetPolls    = "/mb/tcp/polls"
	mbDeletePolls = "/mb/tcp/polls"

	mbTogglePolls = "/mb/tcp/polls/toggle"
	mbImportPolls = "/mb/tcp/polls/config"
	mbExportPolls = "/mb/tcp/polls/config"

	mbCreateFilter = "/mb/tcp/filter/:name"
	mbUpdateFilter = "/mb/tcp/filter/:name"
	mbGetFilter    = "/mb/tcp/filter/:name"
	mbDeleteFilter = "/mb/tcp/filter/:name"

	mbToggleFilter = "/mb/tcp/filter/:name/toggle"

	mbGetFilters    = "/mb/tcp/filters"
	mbDeleteFilters = "/mb/tcp/filters"

	mbToggleFilters = "/mb/tcp/filters/toggle"
	mbImportFilters = "/mb/tcp/filters/config"
	mbExportFilters = "/mb/tcp/filters/config"
	// Poll data
	mbData = "mbtcp.data"
)
