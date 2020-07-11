package subcmd

import (
	"bytes"
	"fmt"
	"gossh/model"
	"gossh/utils"
)

const SLOGAN = "******* hello gooooooooooo ssh *******"

func Run() {
	CmdPrintOut("", SLOGAN)
	for {
		utils.SetConsoleTitle("------ go ssh by zhangyushuang ------")
		inputCmd := ""
		fmt.Scanln(&inputCmd)

		cmds, err := SubcmdMgr.Search(inputCmd)
		if err != nil {
			if _, ok := err.(*NoMatchSubcmd); ok {
				CmdPrintOut("", "input cmd is not a supportted command. see h(help):")
			} else {
				panic(-1)
			}
		}
		isNewhostcmd, isConfcmd := false, false
		if len(cmds) > 1 {
			buf := bytes.NewBufferString("")
			for _, cmd := range cmds {
				buf.WriteString(cmd.Name())
				buf.WriteString("\t")
			}
			CmdPrintOut("", buf.String())
		} else if len(cmds) == 1 {
			if err := cmds[0].Exec(); err != nil {
				if multiHintsErr, ok := err.(*ExecMultiHintsError); ok {
					buf := bytes.NewBufferString("")
					for _, hint := range multiHintsErr.Hints {
						buf.WriteString(hint)
						buf.WriteString("\t")
					}
					CmdPrintOut("", buf.String())
				} else {
					CmdPrintOut("", err)
				}
			} else {
				CmdPrintOut("", SLOGAN)
			}
			_, isNewhostcmd = cmds[0].(*NewhostCmd)
			_, isConfcmd = cmds[0].(*ConfCmd)
		}
		if isNewhostcmd || isConfcmd {
			model.Config.Save()
		} else {
			// multi win process sync configuration
			model.Config.Load()
		}
	}
}
