package main

import (
	"flag"

	"mobilda"
)

var (
	cd  = flag.String("config-dir", "./etc", "Path to config file dir")
	env = flag.String("env", "prod", "Config file environment")
)

func main() {
	flag.Parse()
	app, err := mobilda.NewApplication(*cd, *env)
	if err != nil {
		panic(err)
	}

	//Run application
	app.Run()
}
