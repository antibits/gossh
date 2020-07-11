package subcmd

import (
	"fmt"
	"gossh/model"
	"regexp"
	"strings"
)

const DEFAULT_USER_NAME = "wps"
const ALLOW_FORMAT_ERR_LIMIT = 5

var HOST_PATTERN, _ = regexp.Compile(`^(([1-9][0-9]?)|(1[0-9]{2})|(2[0-4][0-9])|(25[0-5]))(\.(([1-9][0-9]?)|(1[0-9]{2})|(2[0-4][0-9])|(25[0-5]))){3}$`)

func init() {
	registCmd(&NewhostCmd{})
}

type NewhostCmd struct {
}

func (cmd *NewhostCmd) Match(keyword string) bool {
	return strings.HasPrefix("new", keyword)
}

func (cmd *NewhostCmd) Name() string {
	return "new"
}

func (cmd *NewhostCmd) Exec() error {
	count := 0
	CmdPrintOut(cmd.Name(), "user(wps):")
_resetuser:
	if count >= ALLOW_FORMAT_ERR_LIMIT {
		return nil
	}
	user := ""
	fmt.Scanln(&user)
	if len(user) == 0 {
		user = DEFAULT_USER_NAME
	}
	if strings.ContainsAny(user, " ;\\") {
		count++
		CmdPrintOut(cmd.Name(), &InputFormatNotSupport{ALLOW_FORMAT_ERR_LIMIT, count})
		goto _resetuser
	}

	count = 0
	CmdPrintOut(cmd.Name(), "host:")
_resethost:
	if count >= ALLOW_FORMAT_ERR_LIMIT {
		return nil
	}
	host := ""
	fmt.Scanln(&host)
	if len(host) == 0 || !HOST_PATTERN.MatchString(host) {
		count++
		CmdPrintOut(cmd.Name(), &InputFormatNotSupport{ALLOW_FORMAT_ERR_LIMIT, count})
		goto _resethost
	}

	CmdPrintOut(cmd.Name(), "timeout(3):")
	timeout := 0
	fmt.Scanln(&timeout)
	model.Config.NewHost(user, host, timeout)
	// added success connect to host immidiately
	sshCmds, _ := SubcmdMgr.Search(host)
	if len(sshCmds) == 1 {
		return sshCmds[0].Exec()
	}
	return nil
}
