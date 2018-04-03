package manifest

import (
	"github.com/synw/microb-mail/mail"
	"github.com/synw/microb/libmicrob/types"
)

func getCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	return cmds
}

func initService(dev bool, start bool) error {
	mail.ParseTemplate()
	return nil
}
