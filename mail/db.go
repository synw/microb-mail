package mail

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/terr"
)

var database *gorm.DB

type Mail struct {
	gorm.Model
	From    string
	To      string
	Subject string
	Content string
}

func connect() (*gorm.DB, *terr.Trace) {
	db, err := gorm.Open("sqlite3", Conf.DbAddr)
	if err != nil {
		tr := terr.New("mail.initDb", err)
		return db, tr
	}
	return db, nil
}

func InitDb() *terr.Trace {
	msgs.Status("Initializing emails database")
	db, tr := connect()
	if tr != nil {
		tr := terr.Pass("services.logs.db.initDb", tr)
		return tr
	}
	db.AutoMigrate(&Mail{})
	database = db
	return nil
}

func saveToDb(from string, to string, subject string, content string) *terr.Trace {
	entry := &Mail{
		From:    from,
		To:      to,
		Subject: subject,
		Content: content,
	}
	database.Create(entry)
	return nil
}
