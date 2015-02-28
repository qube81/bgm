package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func init() {
	RegisterExitProcess()
	if envvar := os.Getenv("GOMAXPROCS"); envvar == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	if _, err := exec.LookPath(PlayCmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {

	if len(os.Args) == 1 {
		fmt.Println("please input term")
		return
	}

	if os.Args[1] == "--help" {
		ShowHelp()
		return
	}

	query, rate, shuffle, list := ProcessArgs()

	params := DefaultITunesRequestParams()
	params.Term = query

	result := <-SearchMusic(params)

	if shuffle {
		result.Results.Shuffle()
	}

	if list {
		InfoAll(result.Results)
		return
	}

	for i, music := range result.Results {
		Info(music, i+1, result.Count)
		Play(<-Download(music.PreviewURL), rate)
	}

}

func InfoAll(musics []Music) {
	for i, music := range musics {
		Info(music, i+1, len(musics))
	}
}

func ProcessArgs() (query string, rate string, shuffle bool, list bool) {

	rate = "1"
	shuffle = false
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
		default:
			fmt.Printf("unknown option %s\n", v)
			os.Exit(1)
		}

	}

	return
}
