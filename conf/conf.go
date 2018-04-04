package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"os"
	"path/filepath"
)

type Conf struct {
	To       string
	Host     string
	Port     int
	User     string
	Password string
	DbAddr   string
}

func getBasePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	cp := filepath.Dir(ex)
	return cp
}

func getComChan(name string) (string, string) {
	comchan_in := "cmd:$" + name + "_in"
	comchan_out := "cmd:$" + name + "_out"
	return comchan_in, comchan_out
}

func GetServer(conf *types.Conf) (*types.WsServer, *terr.Trace) {
	comchan_in, comchan_out := getComChan(conf.Name)
	s := &types.WsServer{conf.Name, conf.Addr, conf.Key, comchan_in, comchan_out}
	return s, nil
}

func GetConf() (*Conf, *terr.Trace) {
	// set some defaults for conf
	viper.SetConfigName("mail_config")
	viper.AddConfigPath(".")
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 25)
	viper.SetDefault("user", "")
	viper.SetDefault("password", "")
	viper.SetDefault("db", "mails.sqlite")
	// get the actual conf
	err := viper.ReadInConfig()
	if err != nil {
		conf := &Conf{}
		switch err.(type) {
		case viper.ConfigParseError:
			tr := terr.New("conf.getConf", err)
			return conf, tr
		default:
			err := errors.New("Unable to locate config file")
			tr := terr.New("conf.getConf", err)
			return conf, tr
		}
	}
	conf := &Conf{
		To:       viper.Get("to").(string),
		Host:     viper.Get("host").(string),
		Port:     int(viper.Get("port").(float64)),
		User:     viper.Get("user").(string),
		Password: viper.Get("password").(string),
		DbAddr:   viper.Get("db").(string),
	}
	return conf, nil
}
