package main

import "fmt"

func ShowHelp() {
	fmt.Println(`
gobgm [--help] <term...> [-l|--list] [--shuffle] [-r|--rate num]
    
         term:  search query
    -l,--list:  only show result, no play
    -r,--rate:  play at playback rate
    --shuffle:  random order 
       --help:  show help
    
to stop Ctrl-c.
    `)
}
