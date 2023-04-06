package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
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

func cmd() {
	// [youtube] IiXdKjLKknA: Downloading webpage
	// [download] Destination: Connecting Happy!!-IiXdKjLKknA.webm
	// [download] 100% of 4.43MiB in 01:12
	// [ffmpeg] Destination: Connecting Happy!!-IiXdKjLKknA.mp3
	// Deleting original file Connecting Happy!!-IiXdKjLKknA.webm (pass -k to keep)
	var fillter string = "[ffmpeg] Destination: "
	data := "[youtube] IiXdKjLKknA: Downloading webpage"
	data += "\r\n[download] Destination: Connecting Happy!!-IiXdKjLKknA.webm"
	data += "\r\n[download] 100% of 4.43MiB in 01:12"
	data += "\r\n[ffmpeg] Destination: Connecting Happy!!-IiXdKjLKknA.mp3"
	data += "\r\nDeleting original file Connecting Happy!!-IiXdKjLKknA.webm (pass -k to keep)"
	tmp := data
	ary := strings.Split(tmp, "\r")
	for i := 0; i < len(ary); i++ {
		if strings.Index(ary[i], fillter) > 0 {
			tmp = ary[i]
			break
		}
	}
	ary = strings.Split(tmp, "\n")
	for i := 0; i < len(ary); i++ {
		if strings.Index(ary[i], fillter) > 0 {
			tmp = ary[i]
			break
		}
	}
	fmt.Println(data)
	if i := strings.Index(tmp, fillter); i > 0 {
		fmt.Println(tmp[len(fillter)+i:])

	}
}

func main() {
	log.SetFlags(log.Llongfile | log.Flags())
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
