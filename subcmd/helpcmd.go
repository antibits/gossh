package subcmd

import (
	"strings"
)

const HELP_MSG = `
------------------------------------------------------------------------------
	help		h - print this message
	new 		n - add new ssh host
	exit		e - exit
	configure	c - configure
	<host>	  - must by a ip (or substring of ip) login into saved host
------------------------------------------------------------------------------
`

func init() {
	registCmd(&Helpcmd{})
}

type Helpcmd struct {
}

func (*Helpcmd) Match(keyword string) bool {
	return strings.HasPrefix("help", keyword)
}

func (*Helpcmd) Name() string {
	return "help"
}

func (*Helpcmd) Exec() error {
	CmdPrintOutPure(HELP_MSG)
	return nil
}
