package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	YOUTUBE_URL = "www.youtube.com"
)

type HealthMesage struct {
	Downflag bool     `json:downflag`
	DownUrl  []string `json:downurl`
	DownCmd  string   `json:downcmd`
}

func (dataconfig *config) health(w http.ResponseWriter, r *http.Request) {
	msg := HealthMesage{
		Downflag: dataconfig.dl.CKdownflag(),
		DownUrl:  dataconfig.dl.ReadURL(),
		DownCmd:  dataconfig.dl.ReadCMD(),
	}
	jsondata, err := json.Marshal(msg)
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprintf(w, "%s", jsondata)
	}
}

func (dataconfig *config) download(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		tmp := map[string]string{}
		if err := r.ParseForm(); err != nil {
			log.Println(err)
		} else {
			// url := "http://www.youtube.com/watch?v=nTPI23Nd0dg"
			for k, v := range r.Form {
				tmp[k] = v[0]
			}
			url := tmp["url"]
			if strings.Index(url, YOUTUBE_URL) > 0 {
				fmt.Printf("%v %v", r.Method, r.URL.Path)
				for k, v := range tmp {
					fmt.Printf("," + k + "=" + v)
				}
				fmt.Printf("\n")
				dataconfig.dl.Add(url)
				fmt.Fprintln(w, "start download", url)
				return
			}
		}
		if len(tmp) > 0 {
			fmt.Printf("%v %v", r.Method, r.URL.Path)
			for k, v := range tmp {
				fmt.Printf("," + k + "=" + v)
			}
			fmt.Printf("\n")
		} else {
			fmt.Println(r.Method, r.URL.Path)
		}
	} else {
		fmt.Println(r.Method, r.URL.Path)
	}
	fmt.Fprintln(w, "output")
}

func (dataconfig *config) mp3_download(w http.ResponseWriter, r *http.Request) {

}
