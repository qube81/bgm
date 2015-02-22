# gobgm

gobgm - iTunes preview player

## description

gobgm is based on a ruby tool "[bgm](https://github.com/hitode909/bgm)" and written in golang.
Download music preview files from [iTunes Search API](https://www.apple.com/itunes/affiliates/resources/documentation/itunes-store-web-service-search-api.html), and play it via `afplay` command on OSX.


## installation

```
go get github.com/qube81/gobgm
```

## usage

```
gobgm [term...] [--shuffle] [-r|--rate] [--async]

     term:  search query
-r,--rate:  play at playback rate
--shuffle:  random order 
  --async:  play all track at once (max 10 songs)
   --help:  show help
  
 to stop Ctrl-c.
```


## example

```
gobgm SEKAI NO OWARI

gobgm SEKAI NO OWARI -r 0.5

gobgm SEKAI NO OWARI --shuffle

gobgm 落語 --async
```

	
