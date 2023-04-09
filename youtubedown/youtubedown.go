package youtubedown

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/bogem/id3v2/v2"
)

const (
	// CMD_PASS string = "/usr/bin/youtube-dl"
	CMD_PASS string = "/usr/bin/yt-dlp"
	// fillter  string = "[ffmpeg] Destination: "
	fillter     string = "[ExtractAudio] Destination: "
	fillter1    string = ".mp3"
	fillter_jpg string = ".jpg"
)

// CMD_PASSのログからタイトルを取得
func getFileTitle(tmpdata string) (string, error) {

	tmp := tmpdata
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

// ダウンロードしたファイルにタグを追加する
func addTagTitle(filename string) {
	tag, err := id3v2.Open(filename, id3v2.Options{Parse: true})
	if err != nil {
		log.Println("Error while opening mp3 file: ", err)
	}
	defer tag.Close()
	tag.SetTitle(filename[:len(filename)-len(fillter1)])
	if err := tag.Save(); err != nil {
		log.Panicln(err)
	}

}

// ダウンロードしたファイルにアートを追加する
func addTagPicture(filename, filename_j string) {
	tag, err := id3v2.Open(filename, id3v2.Options{Parse: true})
	if tag == nil || err != nil {
		log.Println("Error while opening mp3 file: ", err)
		return
	}
	defer tag.Close()
	artwork, err := ioutil.ReadFile(filename_j)
	if err != nil {
		log.Println("Error while reading artwork file", err)
		return
	}

	pic := id3v2.PictureFrame{
		Encoding:    id3v2.EncodingUTF8,
		MimeType:    "image/jpeg",
		PictureType: id3v2.PTFrontCover,
		Description: "Front cover",
		Picture:     artwork,
	}
	tag.AddAttachedPicture(pic)
	tag.Save()

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
	p, _ := os.Getwd()
	fmt.Println("mv file:", p, name)
	r, err := os.Open(name)
	if err != nil {
		log.Println(err)
		return
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
