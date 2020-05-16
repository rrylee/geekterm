package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	geekhub "github.com/rrylee/geekterm"
	"github.com/spf13/viper"
)

var (
	geekhubDir string
)

var (
	configFile = flag.String("config-file", "", "yaml config file")
	cookie     = flag.String("cookie", "", "geekhub cookie")
	logLevel   = flag.Int("log-level", -1, "log level")
)

// Main entry point.
func main() {
	flag.Parse()

	v := viper.New()
	if *configFile != "" {
		v.SetConfigFile(*configFile)
		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}
	}

	cfg := initConfig(v)

	//Start the application.
	geekhub.Setup(cfg)
	geekhub.Draw()
	geekhub.Keybinds()
	geekhub.WatchUpgrade()

	if err := geekhub.Run(); err != nil {
		fmt.Printf("Error running application: %s\n", err)
	}
}

func initConfig(v *viper.Viper) *geekhub.Config {
	geekhubDir = getUserHomeDir() + "/.geekhub/"
	if _, err := os.Stat(geekhubDir); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(geekhubDir, 0644)
			if err != nil {
				panic(err)
			}
		}
	}

	cfg := &geekhub.Config{}
	if v.GetString("log-file") != "" {
		cfg.LogFile = v.GetString("log-file")
	} else {
		cfg.LogFile = geekhubDir + "log.txt"
	}

	cfg.LogLevel = v.GetInt("log-level")
	if *logLevel >= 0 {
		cfg.LogLevel = *logLevel
	}

	if *cookie != "" {
		cfg.Cookie = *cookie
	} else if v.GetString("cookie") != "" {
		cfg.Cookie = v.GetString("cookie")
	}

	return cfg
}

func getUserHomeDir() string {
	var home string
	switch runtime.GOOS {
	case "windows":
		home, _ = os.LookupEnv("LOCALAPPDATA")
	case "linux":
		home, _ = os.LookupEnv("HOME")
		break
	case "darwin":
		home, _ = os.LookupEnv("HOME")
		break
	}
	return home
}
