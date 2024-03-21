package daemon

type UnitModel struct {
	ID   string `storm:"id"`
	CWD  string
	Bin  string
	Args []string
}
