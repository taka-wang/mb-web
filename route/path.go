package route

//
// route path
//

const (
	pmbOnceRead       = "/mb/tcp/fc/:fc"
	pmbOnceWrite      = "/mb/tcp/fc/:fc"
	pmbGetTimeout     = "/mb/tcp/timeout"
	pmbSetTimeout     = "/mb/tcp/timeout"
	pmbCreatePoll     = "/mb/tcp/poll/:name"
	pmbUpdatePoll     = "/mb/tcp/poll/:name"
	pmbGetPoll        = "/mb/tcp/poll/:name"
	pmbDeletePoll     = "/mb/tcp/poll/:name"
	pmbTogglePoll     = "/mb/tcp/poll/:name/toggle"
	pmbGetPollHistory = "/mb/tcp/poll/:name/history"
	pmbGetPolls       = "/mb/tcp/polls"
	pmbDeletePolls    = "/mb/tcp/polls"
	pmbTogglePolls    = "/mb/tcp/polls/toggle"
	pmbImportPolls    = "/mb/tcp/polls/config"
	pmbExportPolls    = "/mb/tcp/polls/config"
	pmbCreateFilter   = "/mb/tcp/filter/:name"
	pmbUpdateFilter   = "/mb/tcp/filter/:name"
	pmbGetFilter      = "/mb/tcp/filter/:name"
	pmbDeleteFilter   = "/mb/tcp/filter/:name"
	pmbToggleFilter   = "/mb/tcp/filter/:name/toggle"
	pmbGetFilters     = "/mb/tcp/filters"
	pmbDeleteFilters  = "/mb/tcp/filters"
	pmbToggleFilters  = "/mb/tcp/filters/toggle"
	pmbImportFilters  = "/mb/tcp/filters/config"
	pmbExportFilters  = "/mb/tcp/filters/config"
)
