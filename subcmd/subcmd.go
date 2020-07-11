package subcmd

import "fmt"

type ExecMultiHintsError struct {
	Hints []string
}

func (*ExecMultiHintsError) Error() string {
	return "exec with multi hints error"
}

func CmdPrintOut(prefix, msg interface{}) {
	fmt.Printf("%s>%v\n%s>", prefix, msg, prefix)
}

func CmdPrintOutPure(msg interface{}) {
	fmt.Println(msg)
}

type Subcmd interface {
	Match(cmd string) bool

	Name() string

	Exec() error
}
