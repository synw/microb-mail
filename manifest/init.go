package manifest

import (
	"github.com/synw/microb/types"
)

var Service *types.Service = &types.Service{
	"mail",
	getCmds(),
	initService,
}
