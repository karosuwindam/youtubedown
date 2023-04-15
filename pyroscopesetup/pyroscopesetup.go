package pyroscopesetup

import (
	"runtime"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/pyroscope-io/client/pyroscope"
)

type PyroScope struct {
	Hostname      string `env:"HOSTNAME"`
	PyroScopeFlag string `env:"PYROSCOPE_FLAG" envDefault:"TRUE"`
	AppName       string `env:"APPLICATION_NAME" envDefault:"youtubedown"`                     //pyroscopeの表示名
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"http://pyroscorpe.bookserver.home"` //pyroscopeのサーバアドレス
}

var profiletypes []pyroscope.ProfileType = []pyroscope.ProfileType{}
var tags map[string]string = map[string]string{}

func (t *PyroScope) profiledefault() {
	tmp := []pyroscope.ProfileType{
		// these profile types are enabled by default:
		pyroscope.ProfileCPU,
		pyroscope.ProfileAllocObjects,
		pyroscope.ProfileAllocSpace,
		pyroscope.ProfileInuseObjects,
		pyroscope.ProfileInuseSpace,
	}
	profiletypes = append(profiletypes, tmp...)
}

func (t *PyroScope) profileOption() {
	if true {
		tmp := []pyroscope.ProfileType{
			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		}
		profiletypes = append(profiletypes, tmp...)

	}
}

func (t *PyroScope) configTag() {
	tags["hostname"] = t.Hostname
}

func Setup() *PyroScope {
	setupdata := &PyroScope{}
	if err := env.Parse(setupdata); err != nil {
		return nil
	}
	setupdata.profiledefault()
	setupdata.profileOption()
	setupdata.configTag()
	// These 2 lines are only required if you're using mutex or block profiling
	// Read the explanation below for how to set these rates:
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)
	return setupdata

}

func (t *PyroScope) Run() {
	if strings.ToLower(t.PyroScopeFlag) == "true" {
		pyroscope.Start(pyroscope.Config{
			ApplicationName: t.AppName,
			ServerAddress:   t.ServerAddress,
			// you can provide static tags via a map:
			Tags: tags,

			ProfileTypes: profiletypes,
		})

	}

}
