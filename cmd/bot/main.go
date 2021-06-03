package main

import (
	"fmt"
	"github.com/unbyte/lint-bot/config"
	"github.com/unbyte/lint-bot/server"
	"net/http"
	"os"
)

func main() {
	cfgPath := "./config.yml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}
	cfg, err := config.ReadConfig(cfgPath)
	if err != nil {
		panic(err)
	}
	s, err := server.NewServer(cfg)
	if err != nil {
		panic(err)
	}
	if err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.Port), s); err != nil {
		panic(err)
	}
}
