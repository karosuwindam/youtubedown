package pyroscopesetup

import (
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/pyroscope-io/client/pyroscope"
)

type PyroScope struct {
	PyroScopeFlag string `env:"PYROSCOPE_FLAG" envDefault:"TRUE"`
	AppName       string `env:"APPLICATION_NAME" envDefault:"youtubedown"`                     //pyroscopeの表示名
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"http://pyroscorpe.bookserver.home"` //pyroscopeのサーバアドレス
}

func Setup() {
	setupdata := &PyroScope{}
	if err := env.Parse(setupdata); err != nil {
		return
	}
	if strings.ToLower(setupdata.PyroScopeFlag) == "true" {
		pyroscope.Start(pyroscope.Config{
			ApplicationName: setupdata.AppName,
			ServerAddress:   setupdata.ServerAddress,
		})

	}
}
