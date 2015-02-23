package main

import "fmt"

func ShowHelp() {
	fmt.Println(`
    gobgm [--help] <term...> [--shuffle] [-r|--rate] [--async]
    
    term:  search query
    -r,--rate:  play at playback rate
    --shuffle:  random order 
    --async:  play all track at once (max 10 songs)
    --help:  show help
    
    to stop Ctrl-c.
    `)
}
