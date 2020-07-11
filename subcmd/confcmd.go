package subcmd

import (
	"bufio"
	"bytes"
	"fmt"
	"gossh/model"
	"os"
	"strings"
)

const CONF_HELP_MSG = `
------------------------------------------------------------------------------
		s <host> for search hosts;
		d <host> for delete hosts;
		q for quit configure;
		h for help configure;
------------------------------------------------------------------------------`

func init() {
	registCmd(&ConfCmd{})
}

type ConfCmd struct {
}

func (*ConfCmd) Match(keyword string) bool {
	return strings.HasPrefix("configure", keyword)
}

func (*ConfCmd) Name() string {
	return "configure"
}
func (cmd *ConfCmd) printHosts(hosts []model.Host) {
	buf := bytes.NewBufferString("")
	for _, host := range hosts {
		buf.WriteString(host.String())
		buf.WriteString("\t")
	}
	CmdPrintOut(cmd.Name(), buf.String())
}
func (cmd *ConfCmd) searchHosts(keyword string) {
	hosts, _ := model.Config.GetHost(keyword)
	cmd.printHosts(hosts)
}

func (cmd *ConfCmd) deleteHosts(keyword string) {
	if len(keyword) == 0 {
		return
	}
	hosts, _ := model.Config.GetHost(keyword)
	if len(hosts) > 0 {
		cmd.printHosts(hosts)
		fmt.Print("(y/n):")
		_y := ""
		fmt.Scanln(&_y)
		if _y == `y` {
			CmdPrintOut(cmd.Name(), "deleting hosts ...")
			model.Config.DeleteHosts(hosts)
		} else {
			CmdPrintOut(cmd.Name(), "give up deleting hosts")
		}
	} else {
		CmdPrintOut(cmd.Name(), "no hosts found")
	}
}

func (cmd *ConfCmd) Exec() error {
	CmdPrintOut(cmd.Name(), CONF_HELP_MSG)
_reconfig:
	for {
		_cmd := ""
		keyword := ""
		reader := bufio.NewReader(os.Stdin)
		line, _, _ := reader.ReadLine()
		_cmds := strings.Split(string(line), " ")
		__cmds := make([]string, 0, 2)
		for _, c := range _cmds {
			c = strings.Trim(c, " ")
			if len(c) > 0 {
				__cmds = append(__cmds, c)
			}
		}
		if len(__cmds) == 1 {
			_cmd = __cmds[0]
		} else if len(__cmds) > 1 {
			_cmd = __cmds[0]
			keyword = __cmds[1]
		}
		switch {
		case _cmd == `q`:
			break _reconfig
		case _cmd == `h`:
			CmdPrintOut(cmd.Name(), CONF_HELP_MSG)
		case _cmd == `s`:
			if len(keyword) == 0 {
				fmt.Scanln(&keyword)
			}
			cmd.searchHosts(keyword)
		case _cmd == `d`:
			if len(keyword) == 0 {
				fmt.Scanln(&keyword)
			}
			cmd.deleteHosts(keyword)
		default:
			CmdPrintOut(cmd.Name(), CONF_HELP_MSG)
		}
	}
	return nil
}
