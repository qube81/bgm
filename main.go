package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

const (
	asyncLimit = 10
)

func main() {
	RegisterExitProcess()

	if envvar := os.Getenv("GOMAXPROCS"); envvar == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	if len(os.Args) == 1 {
		fmt.Println("please input term")
		return
	}

	if os.Args[1] == "--help" {
		ShowHelp()
		return
	}

	query, rate, shuffle, async, list := ProcessArgs()

	params := ITunesRequestParams{
		Term:    query,
		Country: "JP",
		Lang:    "ja_jp",
		Media:   "music",
		Limit:   "200",
	}

	result := <-SearchMusic(params)

	if shuffle {
		ShuffleMusic(&result.Musics)
	}

	if async {
		if len(result.Musics) > asyncLimit {
			result.Musics = result.Musics[0:asyncLimit]
		}
	}

	if list {
		InfoAll(result.Musics)
		return
	}

	if async {
		PlayAll(result.Musics, rate)
	} else {
		for i, music := range result.Musics {
			Info(music, i+1, result.Count)
			Play(<-Download(music.PreviewURL), rate)
		}
	}

}

func InfoAll(musics []Music) {
	for i, music := range musics {
		Info(music, i+1, len(musics))
	}
}

func PlayAll(musics []Music, rate string) {

	var files []string
	wait := new(sync.WaitGroup)

	for i, music := range musics {

		Info(music, i+1, asyncLimit)

		wait.Add(1)
		go func(music Music) {
			files = append(files, <-Download(music.PreviewURL))
			wait.Done()
		}(music)

	}
	wait.Wait()

	for _, f := range files {
		wait.Add(1)
		go func(fileName string) {
			Play(fileName, rate)
			wait.Done()
		}(f)
	}
	wait.Wait()
}

func ProcessArgs() (query string, rate string, shuffle bool, async bool, list bool) {

	rate = "1"
	shuffle = false
	async = false
	list = false

	query = os.Args[1]

	hasOption := false

	for i := 2; i < len(os.Args); i++ {

		v := os.Args[i]

		if !hasOption {
			if v[0:1] == "-" {
				hasOption = true
			} else {
				query = query + " " + v
				continue
			}
		}

		switch v {
		case "--rate", "-r":
			if i+1 < len(os.Args) {
				rate = os.Args[i+1]
				i++
			}
		case "--shuffle":
			shuffle = true
		case "--list", "-l":
			list = true
		case "--async":
			async = true
		}

	}

	return
}
