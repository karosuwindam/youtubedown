package youtubedown

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	DOWNLOAD_FOLDER = "download"
)

type YouTube_Down struct {
	url_chan     []string
	cmd_out      string
	flag         bool
	downloadflag bool
	mu           sync.Mutex
}

func Setup() YouTube_Down {
	Folderdata.Setup("./" + DOWNLOAD_FOLDER + "/")
	Folderdata.Read("")
	var youtubedown YouTube_Down
	youtubedown.url_chan = []string{}
	youtubedown.flag = true
	return youtubedown
}

func (youtubedown *YouTube_Down) changedownflag(flag bool) {
	youtubedown.mu.Lock()
	youtubedown.flag = flag
	youtubedown.mu.Unlock()
}

func (youtubedown *YouTube_Down) CKdownflag() bool {
	youtubedown.mu.Lock()
	tmp := youtubedown.flag
	youtubedown.mu.Unlock()
	return tmp
}

func (youtubedown *YouTube_Down) ReadURL() []string {
	youtubedown.mu.Lock()
	tmp := youtubedown.url_chan
	youtubedown.mu.Unlock()
	return tmp

}

func (youtubedown *YouTube_Down) writecmd(out string) {
	if out == "" {
		return
	}
	youtubedown.mu.Lock()
	youtubedown.cmd_out = out
	youtubedown.mu.Unlock()
}

func (youtubedown *YouTube_Down) ReadCMD() string {
	youtubedown.mu.Lock()
	tmp := youtubedown.cmd_out
	youtubedown.mu.Unlock()
	return tmp

}

func (youtubedown *YouTube_Down) Add(url string) {
	youtubedown.mu.Lock()
	youtubedown.url_chan = append(youtubedown.url_chan, url)
	youtubedown.mu.Unlock()
}

func (youtubedown *YouTube_Down) download(ctx context.Context, url string) (string, error) {
	var err error = nil
	cmddata := []string{
		"/usr/bin/youtube-dl",
		"-x",
		"--audio-format",
		"mp3",
	}
	if url != "" {
		cmddata = append(cmddata, url)
	}

	ctx1, cancel1 := context.WithCancel(ctx)

	cmd := exec.Command("python3", cmddata...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Start()

	fmt.Println("Download Start for", url)
	go func(ctx1 context.Context) {
		for {
			select {
			case <-ctx1.Done():
				youtubedown.writecmd("")
			default:
				youtubedown.writecmd(stdout.String())
			}
			time.Sleep(100 * time.Microsecond)
		}
	}(ctx1)
	cmd.Wait()
	out := stdout.String()
	if stderr.String() != "" {
		err = errors.New(stderr.String())
	}

	cancel1()
	if err != nil {
		return "", err
	}
	tmp := string(out)
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
	if i := strings.Index(tmp, fillter); i > 0 {
		tmp = tmp[len(fillter)+i:]
	} else {
		return "", errors.New("Not name " + fillter)
	}
	if i := strings.Index(tmp, fillter1); i > 0 {
		tmp = tmp[:i+len(fillter1)]
	} else {
		return "", errors.New("Not name " + fillter1)
	}

	return tmp, nil

}

type FileData struct {
	No   int
	Name string
	Pass string
}

func (youtubedown *YouTube_Down) Mp3ListGet() []FileData {
	youtubedown.mu.Lock()
	Folderdata.Read("")
	tmp := Folderdata.Data
	youtubedown.mu.Unlock()
	output := []FileData{}
	count := 1
	for _, file := range tmp {
		if strings.Index(file.Name, fillter1) > 0 {
			outtmp := FileData{
				No:   count,
				Name: file.Name,
				Pass: file.RootPath,
			}
			output = append(output, outtmp)
			count++
		}
	}

	return output
}

func (youtubedown *YouTube_Down) Run(ctx context.Context) (string, error) {
	if !youtubedown.flag {
		return "", errors.New("Not Youtube down SetUp")
	}
	fmt.Println("download loop start")
	for {
		select {
		case <-ctx.Done():
			return "", nil
		default:
			var url string
			youtubedown.changedownflag(true)
			if len(youtubedown.url_chan) > 0 {
				youtubedown.mu.Lock()
				url = youtubedown.url_chan[0]
				// youtubedown.url_chan = youtubedown.url_chan[1:]
				youtubedown.mu.Unlock()

			} else {
				youtubedown.changedownflag(false)
			}
			if url != "" {
				// if str, err := download(url); err != nil {
				if str, err := youtubedown.download(ctx, url); err != nil {
					log.Println(url, ":", err)
				} else {
					mvfolder(str)
					fmt.Println(str)
				}
				youtubedown.mu.Lock()
				tmp := []string{}
				for _, tmpurl := range youtubedown.url_chan {
					if tmpurl != url {
						tmp = append(tmp, tmpurl)
					}
				}
				youtubedown.url_chan = tmp
				Folderdata.Read("")
				youtubedown.mu.Unlock()
			}
		}
		time.Sleep(100 * time.Microsecond)
	}
}
