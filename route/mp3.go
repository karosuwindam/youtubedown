package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bogem/id3v2/v2"
)

type Mp3Id3Tag struct {
	Title  string `json:title`
	Artist string `json:artist`
	Album  string `json:Album`
	Year   string `json:year`
	Lyrics string `json:lyrics`
}

// Id3v2出力
func mp3edit(id string, data Mp3Id3Tag) {
	files := dataconfig.dl.Mp3ListGet()
	for _, file := range files {
		if strconv.Itoa(file.No) == id {
			//ファイルをホームフォルダにコピー
			r, err := os.Open(file.Pass + file.Name)
			if err != nil {
				log.Fatal(err)
			}
			c, err := os.Create(file.Name)
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(c, r)
			if err != nil {
				log.Fatal(err)
			}
			r.Close()
			c.Close()

			tag, err := id3v2.Open(file.Name, id3v2.Options{Parse: true})
			if err != nil {
				log.Fatal("Error while opening mp3 file: ", err)
			}
			defer tag.Close()
			tag.SetDefaultEncoding(id3v2.EncodingUTF16)
			tag.SetTitle(data.Title)
			tag.SetArtist(data.Artist)
			tag.SetAlbum(data.Album)
			tag.SetYear(data.Year)
			uslt := id3v2.UnsynchronisedLyricsFrame{
				Encoding:          id3v2.EncodingUTF16,
				Language:          "eng",
				ContentDescriptor: "",
				Lyrics:            data.Lyrics,
			}
			tag.AddUnsynchronisedLyricsFrame(uslt)
			// fmt.Println(uslt)

			if err := tag.Save(); err != nil {
				log.Panicln(err)
			}
			//ホームフォルダのファイルをコピー
			r2, err := os.Open(file.Name)
			if err != nil {
				log.Fatal(err)
			}
			c2, err := os.Create(file.Pass + file.Name)
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(c2, r2)
			if err != nil {
				log.Fatal(err)
			}
			r2.Close()
			c2.Close()
			os.Remove(file.Name)

			break
		}
	}
	fmt.Println(data)
}

// Id3v2読み取り
func mp3read(id string) Mp3Id3Tag {

	files := dataconfig.dl.Mp3ListGet()
	output := Mp3Id3Tag{}
	for _, file := range files {
		if strconv.Itoa(file.No) == id {
			tag, err := id3v2.Open(file.Pass+file.Name, id3v2.Options{Parse: true})
			if err != nil {
				log.Fatal("Error while opening mp3 file: ", err)
			}
			defer tag.Close()
			lyrics := ""
			uslfs := tag.GetFrames(tag.CommonID("Unsynchronised lyrics/text transcription"))
			for _, f := range uslfs {
				uslf, ok := f.(id3v2.UnsynchronisedLyricsFrame)
				if !ok {
					log.Fatal("Couldn't assert USLT frame")
				}

				lyrics = uslf.Lyrics
				// fmt.Println(uslf.Lyrics)
			}

			output = Mp3Id3Tag{
				Title:  tag.Title(),
				Artist: tag.Artist(),
				Album:  tag.Album(),
				Year:   tag.Year(),
				Lyrics: lyrics,
			}
			break

		}
	}
	return output
}

type Id3v2Image struct {
	name      string
	imagetype string
	image     []byte
}

// Id3v2のイメージ読み取り
func mp3readimage(id string) (Id3v2Image, error) {
	files := dataconfig.dl.Mp3ListGet()
	output := Id3v2Image{}
	for _, file := range files {
		if strconv.Itoa(file.No) == id {
			tag, err := id3v2.Open(file.Pass+file.Name, id3v2.Options{Parse: true})
			if err != nil {
				log.Println(err)
				return output, err
			}
			defer tag.Close()
			pictures := tag.GetFrames(tag.CommonID("Attached picture"))
			if len(pictures) > 0 {
				pic, ok := pictures[0].(id3v2.PictureFrame)
				if !ok {
					return output, errors.New("Couldn't assert picture frame")
				}
				output.name = pic.Description
				output.imagetype = pic.MimeType
				output.image = []byte(pic.Picture)

			} else {
				return output, errors.New("No images")
			}
			return output, nil
		}
	}
	return output, errors.New("no file data")

}

// id3タグの編集
// /edit/:id
func (dataconfig *config) mp3edit(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{}
	urldatas := urlAnalysis(r.URL.Path)
	for i, urldata := range urldatas {
		if urldata == "edit" {
			data["id"] = urldatas[i+1]
			break
		}
	}
	if r.Method == "POST" {
		var d Mp3Id3Tag
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mp3edit(data["id"], d)
		fmt.Println(r.Method, r.URL.Path, d)
	} else {
		fmt.Println(r.Method, r.URL.Path)
	}
	jsondata, err := json.Marshal(mp3read(data["id"]))
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprintf(w, "%s", jsondata)
	}
}

// id3タグの表示
// /view/:id
func (dataconfig *config) mp3view(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{}
	urldatas := urlAnalysis(r.URL.Path)
	for i, urldata := range urldatas {
		if urldata == "view" {
			data["id"] = urldatas[i+1]
			break
		}
	}

	fmt.Println(r.Method, r.URL.Path)
	jsondata, err := json.Marshal(mp3read(data["id"]))
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprintf(w, "%s", jsondata)
	}
}

// id3タグの画像表示
// /mp3image/:id
func (dataconfig *config) mp3viewimage(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{}
	urldatas := urlAnalysis(r.URL.Path)
	for i, urldata := range urldatas {
		if urldata == "mp3image" {
			data["id"] = urldatas[i+1]
			break
		}
	}

	fmt.Println(r.Method, r.URL.Path)
	image, err := mp3readimage(data["id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		// コンテントタイプ
		w.Header().Set("Content-Type", image.imagetype)
		// ファイル名
		w.Header().Set("Content-Disposition", "attachment; filename="+image.name)
		w.Header().Set("Content-Length", string(len(image.image)))
		w.Write(image.image)
		// fmt.Fprintf(w, "%s", image.image)

	}
}
