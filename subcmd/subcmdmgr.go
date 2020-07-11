package subcmd

import "fmt"

type subcmds []Subcmd

type NoMatchSubcmd struct {
}

func (*NoMatchSubcmd) Error() string {
	return "no match subcmd"
}

type InputFormatNotSupport struct {
	Limit int
	Count int
}

func (e *InputFormatNotSupport) Error() string {
	return fmt.Sprintf("input format err. try again.(%d/%d)", e.Count, e.Limit)
}

type subcmdMgr struct {
	cmds subcmds
}

func (mgr *subcmdMgr) Search(keyword string) ([]Subcmd, error) {
	result := make([]Subcmd, 0)
	for _, cmd := range mgr.cmds {
		if cmd.Match(keyword) {
			result = append(result, cmd)
		}
	}
	if len(result) == 0 {
		return []Subcmd{&Helpcmd{}}, &NoMatchSubcmd{}
	}
	return result, nil
}

var SubcmdMgr = &subcmdMgr{}

func registCmd(cmd Subcmd) {
	SubcmdMgr.cmds = append(SubcmdMgr.cmds, cmd)
}
