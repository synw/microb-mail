package manifest

import (
	"github.com/synw/microb-mail/mail"
	//"github.com/synw/microb/msgs"
	"github.com/synw/microb/events"
	"github.com/synw/microb/types"
	"github.com/synw/terr"
)

func getCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["mails"] = mails()
	return cmds
}

func initService(dev bool, start bool) *terr.Trace {
	tr := mail.Init(dev)
	if tr != nil {
		tr = tr.Pass()
		return tr
	}
	tr = mail.InitDb()
	if tr != nil {
		tr = tr.Pass()
		return tr
	}
	mail.ParseTemplate()
	return nil
}

func mails() *types.Cmd {
	cmd := &types.Cmd{Name: "mails", Exec: runMails}
	return cmd
}

func runMails(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	// this function will be run on command call
	var resp []interface{}
	resp = append(resp, "Last 10 mails sent:")
	mails, tr := mail.GetMails()
	if tr != nil {
		tr = tr.Add("manifest.Cmds", "Can not get mails")
		events.Error("mail", "Can not get emails", tr)
		cmd.Status = "error"
		c <- cmd
	}
	for _, mail := range mails {
		row := mail.From + " : " + mail.Subject
		resp = append(resp, row)
	}
	// the command will return "Hello world"
	cmd.ReturnValues = resp
	cmd.Status = "success"
	c <- cmd
}
