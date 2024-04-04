package daemon

type UnitID uint32

type UnitModel struct {
	ID            UnitID `storm:"id,increment"`
	Name          string
	CWD           string
	Bin           string
	Args          []string
	RestartsCount uint32
}
