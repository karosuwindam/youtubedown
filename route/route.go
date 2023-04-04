package route

import (
	"fmt"
	"net/http"
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

var Route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/download", Handler: dataconfig.download},
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
