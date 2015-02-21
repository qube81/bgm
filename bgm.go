package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
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
	argsLen        int
	rate           = 1
	err            error
	search         string
	resultCount    int
	nowPlayingFile string
)

func main() {
	RegisterExitProcess()

	search = os.Args[1]

	for i := 2; i < argsLen; i++ {

		v := os.Args[i]

		if i+1 < argsLen {
			if v == "--rate" || v == "-r" {
				rate, err = strconv.Atoi(os.Args[i+1])
				if err != nil {
					rate = 1
				}
				i++
			}
		}

	}

	itunes := <-RequestItunes(search)
	resultCount = itunes.ResultCount

	for i, music := range itunes.Results {
		Info(music, i+1)
		Play(<-Download(music.PreviewURL))
	}

}

func Play(fileName string) {
	defer os.Remove(fileName)
	nowPlayingFile = fileName
	out, _ := exec.Command("afplay", fileName, "--rate", strconv.Itoa(rate)).CombinedOutput()
	fmt.Print(string(out))
}

func Info(music Music, num int) {
	fmt.Printf("♪ (%d/%d)\n", num, resultCount)
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

		file, err := ioutil.TempFile(os.TempDir(), "tmp_bgm_")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		io.Copy(file, response.Body)
		fileNameChan <- file.Name()

	}(url)

	return fileNameChan

}

func RequestItunes(term string) <-chan ItunesResult {
	resultChan := make(chan ItunesResult)

	itunesEndPoint := "https://itunes.apple.com/search?term=%s&country=JP&media=music&limit=200"

	go func(url string) {
		fmt.Print("Request Itunes...")
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)

		} else {
			defer response.Body.Close()
			fmt.Println("http status: " + response.Status + "\n")

			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}

			var data ItunesResult

			json.Unmarshal([]byte(contents), &data)

			resultChan <- data
		}

	}(fmt.Sprintf(itunesEndPoint, term))

	return resultChan
}

func RegisterExitProcess() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			os.Remove(nowPlayingFile)
			os.Exit(1)
		}
	}()
}
