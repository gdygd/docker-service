package memory

const (
	STATE_OFF = 0
	STATE_ON  = 1
)

type SystemInfo struct {
	SvrUtc    int64 `json:"svrutc"`
	DbSvrComm int   `json:"dbstate"`
	RdSvrComm int   `json:"rdstate"`
}
