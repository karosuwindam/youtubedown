package route

import (
	"fmt"
	"net/http"
	"strings"
	"youtubedown/webserver"
	"youtubedown/youtubedown"
)

type config struct {
	dl   *youtubedown.YouTube_Down
	flag bool
}

var dataconfig config

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello World")
}

// URL解析用
func urlAnalysis(url string) []string {
	tmp := []string{}
	for _, str := range strings.Split(url[1:], "/") {
		tmp = append(tmp, str)
	}
	return tmp
}

var Route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/mp3/", Handler: dataconfig.mp3_download},
	{Pass: "/mp3image/", Handler: dataconfig.mp3viewimage},
	{Pass: "/view/", Handler: dataconfig.mp3view},
	{Pass: "/edit/", Handler: dataconfig.mp3edit},
	{Pass: "/download", Handler: dataconfig.download},
	{Pass: "/list", Handler: dataconfig.mp3list},
	{Pass: "/health", Handler: dataconfig.health},
	{Pass: "/", Handler: viewhtml},
}

func Setup(dl *youtubedown.YouTube_Down) {
	dataconfig.flag = true
	dataconfig.dl = dl
}

func Add(homepass string, cfg *webserver.SetupServer) {
	output := Route
	if homepass != "" && homepass != "/" {
		if homepass[0:1] != "/" {
			for i, tmp := range output {
				tmp.Pass = "/" + homepass + tmp.Pass
				output[i] = tmp
			}
		} else {
			for i, tmp := range output {
				tmp.Pass = homepass + tmp.Pass
				output[i] = tmp
			}

		}
	}
	webserver.Config(cfg, output)
}
