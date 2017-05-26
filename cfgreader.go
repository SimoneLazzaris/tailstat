package main
import (
	"github.com/BurntSushi/toml"
	"fmt"
)

type metric struct {
	Name  string
	Match string
	Counter int
}

type Config struct {
	Logfile string
	Prefix  string
	Deltat  int
	Graphite string
	Metrics map[string] metric
}

var Cfg Config

func initCfg() {
	if _, err := toml.DecodeFile("/usr/local/etc/tailstat.toml", &Cfg); err != nil {
		fmt.Println(err)
	}
}
