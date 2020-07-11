package subcmd

import (
	"gossh/model"
	"os"
	"strings"
)

func init() {
	registCmd(&ExitCmd{})
}

type ExitCmd struct {
}

func (*ExitCmd) Match(keyword string) bool {
	return strings.HasPrefix("exit", keyword)
}

func (*ExitCmd) Name() string {
	return "exit"
}

func (*ExitCmd) Exec() error {
	model.Config.Save()
	os.Exit(0)
	return nil
}
