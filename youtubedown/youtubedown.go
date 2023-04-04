package youtubedown

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	fillter  string = "[ffmpeg]"
	fillter1 string = ".mp3"
)

func download(url string) (string, error) {
	cmddata := []string{
		"/usr/bin/youtube-dl",
		"-x",
		"--audio-format",
		"mp3",
	}
	if url != "" {
		cmddata = append(cmddata, url)
	}
	cmd := exec.Command("python3", cmddata...)
	fmt.Println("Download Start for", url)
	out, err := cmd.Output()
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
	ary = strings.Split(tmp, " ")
	for i := 0; i < len(ary); i++ {
		if n := strings.Index(ary[i], fillter1); n > 0 {
			tmp = ary[i][:n+len(fillter1)]
			break
		}
	}

	return tmp, nil

}

// ダウンロードしたファイルを特定に移動する
func mvfolder(name string) {
	///フォルダ作成
	_, err := os.Stat(DOWNLOAD_FOLDER)
	if os.IsNotExist(err) {
		fmt.Println("file does not exist")
		if err := os.Mkdir(DOWNLOAD_FOLDER, 0775); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Create folder", DOWNLOAD_FOLDER)
		}
	}
	//ファイル移動
	// if err := os.Rename(name, DOWNLOAD_FOLDER+"/"+name); err != nil {
	// 	log.Fatal(err)
	// }
	// 大本ファイル
	r, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	//出力先
	// 作成するファイル
	dest := DOWNLOAD_FOLDER + "/" + name

	c, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(c, r)
	if err != nil {
		log.Fatal(err)
	}
	r.Close()
	os.Remove(name)

}
