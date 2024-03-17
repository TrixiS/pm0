package daemon

type Process struct {
	ID   string `storm:"id"`
	CWD  string
	Bin  string
	Args []string
	PID  int
}
