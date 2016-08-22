package route

type timeoutReq struct {
	Data int64 `json:"timeout,omitempty"`
}

type timeoutRes struct {
	Status string `json:"status"`
	Data   int64  `json:"timeout,omitempty"`
}

type simpleRes struct {
	Status string `json:"status"`
}
