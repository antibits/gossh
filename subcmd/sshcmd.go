package subcmd

import (
	"bytes"
	"context"
	"fmt"
	"gossh/model"
	"gossh/utils"
	"os"
	"os/exec"
	"time"
)

func init() {
	registCmd(&SshCmd{})
}

type SshCmd struct {
	hintsHosts []model.Host
}

func (cmd *SshCmd) Match(keyword string) bool {
	hosts, err := model.Config.GetHost(keyword)
	if err != nil {
		return false
	}
	cmd.hintsHosts = hosts
	return len(hosts) > 0
}

func (cmd *SshCmd) Name() string {
	if len(cmd.hintsHosts) > 0 {
		buf := bytes.NewBufferString("")
		for _, host := range cmd.hintsHosts {
			buf.WriteString(host.Host)
			buf.WriteString("\t")
		}
		return buf.String()
	}
	return ""
}

func (cmd *SshCmd) Exec() error {
	if len(cmd.hintsHosts) > 1 {
		return &ExecMultiHintsError{
			Hints: func() []string {
				hosts := make([]string, 0, len(cmd.hintsHosts))
				for _, hintHost := range cmd.hintsHosts {
					hosts = append(hosts, hintHost.Host)
				}
				return hosts
			}(),
		}
	}
	ssh(cmd.hintsHosts[0].User, cmd.hintsHosts[0].Host, cmd.hintsHosts[0].Timeout)
	return nil
}

func setSshConsoleTitle(ctx context.Context, title string) {
	go func() {
		tick := time.NewTicker(3 * time.Second)
		for {
			select {
			case <-tick.C:
				utils.SetConsoleTitle(title)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func ssh(user, host string, timeoutSec int) {
	CmdPrintOut("", fmt.Sprintf("login into %s ...", host))
	url := host
	if len(user) > 0 {
		url = fmt.Sprintf("%s@%s", user, host)
	}
	if timeoutSec == 0 {
		timeoutSec = 3
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "ssh", "-i", "~/.ssh/ssh", "-o", fmt.Sprintf("ConnectTimeout=%d", timeoutSec), url)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	setSshConsoleTitle(ctx, url)
	err := cmd.Run()
	if err != nil {
		CmdPrintOut("", fmt.Sprintf("cmd run with some err.!, %v", err))
		time.Sleep(3 * time.Second)
		os.Exit(-1)
	}
}
