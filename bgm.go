package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func main() {

	var (
		argsLen       = len(os.Args)
		rate          = 1
		err           error
		search        string
		MatsukenSamba = "http://a989.phobos.apple.com/us/r20/Music2/v4/a2/b5/1b/a2b51beb-4b0a-a9a2-a636-03c215ab2db3/mzaf_3999292651922270843.aac.m4a"
	)

	for i := 0; i < argsLen; i++ {

		if i == 0 {
			continue
		}

		v := os.Args[i]

		if i == 1 {
			search = v
		}

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

	json := <-RequestItunes(search)
	fmt.Println(json)

	filename := <-Download(MatsukenSamba)
	defer os.Remove(filename)

	fmt.Println("NOW PLAYING: {#artistname} - {#trackname}")
	out, _ := exec.Command("afplay", filename, "--rate", strconv.Itoa(rate)).CombinedOutput()
	fmt.Print(string(out))

}

/*
GetMusicFile download specified url music file, and save temp dir
@return channel tempfile name
*/
func Download(url string) <-chan string {

	fileNameChan := make(chan string)

	go func(url string) {
		fmt.Print("Loading...")
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

		fmt.Println("http status: " + response.Status)
		io.Copy(file, response.Body)
		fileNameChan <- file.Name()

	}(url)

	return fileNameChan

}

/*

*/
func RequestItunes(term string) <-chan string {
	jsonChan := make(chan string)
	itunesEndPoint := "https://itunes.apple.com/search?term=%s&country=JP&media=music&entity=song&attribute=songTerm&limit=1"

	go func(url string) {
		fmt.Print("Request Itunes...")
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		} else {
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			jsonChan <- string(contents)
		}

	}(fmt.Sprintf(itunesEndPoint, term))

	return jsonChan
}
