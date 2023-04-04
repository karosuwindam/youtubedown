package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

func mp3Output(root string, name string, w http.ResponseWriter) {
	pass := root + name
	fmt.Println("Download pass:", pass)
	if file, err := os.Open(pass); err != nil {
		log.Println(err)
	} else {
		defer file.Close()
		buf := make([]byte, 1024)
		var buffer []byte
		for {
			n, err := file.Read(buf)
			if n == 0 {
				break
			}
			if err != nil {
				// Readエラー処理
				break
			}
			buffer = append(buffer, buf[:n]...)
		}
		// ファイル名
		w.Header().Set("Content-Disposition", "attachment; filename="+name)
		// コンテントタイプ
		w.Header().Set("Content-Type", "application/mpeg")
		// ファイルの長さ
		w.Header().Set("Content-Length", string(len(buffer)))
		// bodyに書き込み
		w.Write(buffer)
	}
}

// /mp3/:idの結果
func (dataconfig *config) mp3_download(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{}
	urldatas := urlAnalysis(r.URL.Path)
	for i, urldata := range urldatas {
		if urldata == "mp3" {
			data["id"] = urldatas[i+1]
			break
		}
	}
	files := dataconfig.dl.Mp3ListGet()
	for _, file := range files {
		if strconv.Itoa(file.No) == data["id"] {
			mp3Output(file.Pass, file.Name, w)
			fmt.Fprintln(w, file.Pass, file.Name)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "error id")
}

type M3List struct {
	No   int    `json:no`
	Name string `json:name`
}

func (dataconfig *config) mp3list(w http.ResponseWriter, r *http.Request) {
	files := dataconfig.dl.Mp3ListGet()
	output := []M3List{}
	for _, file := range files {
		tmp := M3List{No: file.No, Name: file.Name}
		output = append(output, tmp)
	}
	jsondata, err := json.Marshal(output)
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprintf(w, "%s", jsondata)
	}
}
