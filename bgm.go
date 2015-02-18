package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {

	var MatsukenSamba = "http://a588.phobos.apple.com/us/r20/Music/v4/82/a0/27/82a02731-c7f8-fec1-c686-dad12eca272f/mzaf_8898382402336846733.m4a"

	filename := <-GetMusicFile(MatsukenSamba)
	defer os.Remove(filename)

	fmt.Println("NOW PLAYING: Matsuken Samba 💃 💃 💃 💃 💃 💃 💃 💃 💃 💃 💃 💃 💃 💃 💃 💃 💃")
	out, _ := exec.Command("afplay", filename, "-r", "1").CombinedOutput()
	fmt.Print(string(out))
}

/*
GetMusicFile download specified url music file, and save temp dir
@return channel tempfile name
*/
func GetMusicFile(url string) <-chan string {

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
