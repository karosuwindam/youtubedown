package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
			tag, err := id3v2.Open(file.Pass+file.Name, id3v2.Options{Parse: true})
			if err != nil {
				log.Fatal("Error while opening mp3 file: ", err)
			}
			defer tag.Close()
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

			if err := tag.Save(); err != nil {
				log.Panicln(err)
			}
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
				fmt.Println(uslf.Lyrics)
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
	jsondata, err := json.Marshal(mp3read(data["id"]))
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprintf(w, "%s", jsondata)
	}
}
