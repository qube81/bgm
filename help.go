package main

import "fmt"

func ShowHelp() {
	fmt.Println(`
gobgm [--help] <term...> [-l|--list] [--shuffle] [-r|--rate] [--async]
    
         term:  search query
    -l,--list:  only show result, no play
    -r,--rate:  play at playback rate
    --shuffle:  random order 
      --async:  play all track at once (max 10 songs)
       --help:  show help
    
to stop Ctrl-c.
    `)
}
