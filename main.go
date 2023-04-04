package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func cmdbbb() {
	cmddata := []string{
		"/usr/bin/youtube-dl",
		"-x",
		"--audio-format",
		"mp3",
		"https://www.youtube.com/watch?v=nTPI23Nd0dg",
	}
	cmd := exec.Command("python3", cmddata...)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var stdout1 bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout1
	cmd.Stderr = &stderr
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println(i, stdout1.String())
			time.Sleep(time.Second)

		}
	}()
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	cmd.Wait()
	fmt.Println(stdout1.String())
	if stderr.String() != "" {
		fmt.Println(stderr.String())
	} else {
		fmt.Println("no err")
	}

}

func main() {
	// cmdbbb()
	// return
	ctx := context.Background()
	// url := "https://www.youtube.com/watch?v=nTPI23Nd0dg"
	// stopdata.dl.Add(url)
	// url = "https://www.youtube.com/watch?v=nTPI23Nd0d"
	// stopdata.dl.Add(url)
	// time.Sleep(1 * time.Second)
	// fmt.Println("Shutdown now")
	// return
	fmt.Println("start")
	if err := Run(ctx); err != nil {
		EndCK()
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("end")
	EndCK()
}
