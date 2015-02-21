package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

type Music struct {
	ArtistName     string `json:"artistName"`
	TrackName      string `json:"trackName"`
	PreviewURL     string `json:"previewUrl"`
	CollectionName string `json:"CollectionName"`
	TrackViewURL   string `json:"TrackViewURL"`
}

type ItunesResult struct {
	ResultCount int `json:"resultCount"`
	Results     []Music
}

var (
	downloadedFiles []string
)

func main() {
	RegisterExitProcess()

	search, rate, shuffle := ProcessArgs()

	itunes := <-RequestITunes(search)

	if shuffle {
		Shuffle(&itunes.Results)
	}

	for i, music := range itunes.Results {
		Info(music, i+1, itunes.ResultCount)
		Play(<-Download(music.PreviewURL), rate)
	}

}

func ProcessArgs() (search string, rate string, shuffle bool) {

	argsLen := len(os.Args)
	if argsLen == 1 {
		fmt.Println("please input term")
		os.Exit(1)
	}

	search = os.Args[1]
	rate = "1"
	shuffle = false

	hasOption := false

	for i := 2; i < argsLen; i++ {

		v := os.Args[i]

		if !hasOption {
			if v[0:1] == "-" {
				hasOption = true
			} else {
				search = search + " " + v
				continue
			}
		}

		switch v {
		case "--rate", "-r":
			if i+1 < argsLen {
				rate = os.Args[i+1]
				i++
			}
		case "--shuffle":
			shuffle = true
		}

	}

	return
}

func Play(fileName string, rate string) {
	defer os.Remove(fileName)

	out, _ := exec.Command("afplay", fileName, "--rate", rate).CombinedOutput()
	fmt.Print(string(out))
}

func Info(music Music, num int, total int) {
	fmt.Printf("♪ (%d/%d)\n", num, total)
	fmt.Printf("# %s - %s / %s\n", music.TrackName, music.ArtistName, music.CollectionName)
	fmt.Printf("%s\n", music.TrackViewURL)
	fmt.Println()
}

func Download(url string) <-chan string {

	fileNameChan := make(chan string)

	go func(url string) {
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		file, err := ioutil.TempFile(os.TempDir(), "bgm_tmp")
		downloadedFiles = append(downloadedFiles, file.Name())

		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		io.Copy(file, response.Body)

		fileNameChan <- file.Name()

	}(url)

	return fileNameChan

}

func RequestITunes(term string) <-chan ItunesResult {
	resultChan := make(chan ItunesResult)

	params := url.Values{}
	params.Add("term", term)
	params.Add("country", "JP")
	params.Add("media", "music")
	params.Add("limit", "200")

	itunesEndPoint := "https://itunes.apple.com/search/"

	go func(endPoint string) {
		fmt.Println("Request iTunes search API...")
		fmt.Print()
		response, err := http.Get(endPoint + "?" + params.Encode())

		if err != nil {
			log.Fatal(err)
			os.Exit(1)

		} else {
			defer response.Body.Close()

			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			var data ItunesResult
			json.Unmarshal([]byte(contents), &data)

			resultChan <- data
		}

	}(itunesEndPoint)

	return resultChan
}

func Shuffle(a *[]Music) {

	musics := *a

	rand.Seed(time.Now().UnixNano())
	for i := range musics {
		j := rand.Intn(i + 1)
		musics[i], musics[j] = musics[j], musics[i]
	}
}

func RegisterExitProcess() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			for _, file := range downloadedFiles {
				os.Remove(file)
			}
			os.Exit(1)
		}
	}()
}
