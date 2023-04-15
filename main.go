package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"youtubedown/pyroscopesetup"
	"youtubedown/route"
	"youtubedown/webserver"
	"youtubedown/youtubedown"
)

func RootConfg() []webserver.WebConfig {
	output := []webserver.WebConfig{}
	tmp := route.Route
	output = append(output, tmp...)
	return output
}

func Config(cfg *webserver.SetupServer) error {
	// webserver.Config(cfg, RootConfg())
	dl := youtubedown.Setup()
	Stopdata.dl = &dl
	route.Setup(Stopdata.dl)
	route.Add("/", cfg)
	return nil
}

func Run(ctx context.Context) error {
	cfg, err := webserver.NewSetup()
	if err != nil {
		return err
	}
	if err := Config(cfg); err != nil {
		return err
	}
	s, err := cfg.NewServer()
	if err != nil {
		return err
	}
	go chaild(ctx, Stopdata.dl)

	return s.Run(ctx)
}

func chaild(ctx context.Context, youtubedown *youtubedown.YouTube_Down) {
	youtubedown.Run(ctx)
}

type strdata struct {
	dl *youtubedown.YouTube_Down
}

var Stopdata strdata

func EndCK() {
	for {
		if !Stopdata.dl.CKdownflag() {
			break
		}
		time.Sleep(100 * time.Microsecond)
	}

}

func main() {
	log.SetFlags(log.Llongfile | log.Flags())
	py := pyroscopesetup.Setup()
	py.Run()
	ctx := context.Background()
	fmt.Println("start")
	if err := Run(ctx); err != nil {
		EndCK()
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("end")
	EndCK()
}
