package daemon

type UnitModel struct {
	ID   string `storm:"id"`
	Name string
	CWD  string
	Bin  string
	Args []string
}
