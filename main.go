package main

import (
	"fmt"

	"github.com/aquilax/go-perlin"
	"github.com/nsf/termbox-go"
	// "time"
	// "math/rand"
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"net/http"

	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var basePath = "/"
var baseListener = ":8080"

type Tile struct {
	X, Y, TilePrint int
	Type            int     `json:"-"`
	NoiseValue      float64 `json:"-"`
}

type Map struct {
	alpha float64
	beta  float64
	n     int
	div   float64
	seed  int64
	Tiles map[string]Tile
}

type MapsWatch struct {
	Maps map[int]Map
}

var allMaps = MapsWatch{}

const (
	// alpha       = 2.
	// beta        = 2.
	// n           = 3
	// div 		= 50
	// seed  int64 = 100

	defaultHeight = 50
	defaultWidth  = 50
)

func getMap(mapId int, width int, height int, startX int, startY int) {

	theMap, ok := allMaps.Maps[mapId]
	if !ok {
		//Clean Up
		panic("Map does not exist")
	}

	// check if we have generated this part of the map before
	w, h := termbox.Size()

	w = width + startX
	h = height + startY

	p := perlin.NewPerlin(theMap.alpha, theMap.beta, theMap.n, theMap.seed)
	_, _, _ = w, p, strconv.Itoa
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for y := startY; y < h; y++ {
		for x := startX; x < w; x++ {

			tileId := strconv.FormatInt(int64(x), 10) + ":" + strconv.FormatInt(int64(y), 10)
			tile, okT := theMap.Tiles[tileId]
			char := rune(' ')

			if !okT {
				tile = Tile{}

				noiseFloat := p.Noise2D(float64(x)/theMap.div, float64(y)/theMap.div)
				noise := noiseFloat
				os.Getenv("GG_MAP_PATH")
				noiseFloat = math.Abs(noiseFloat*theMap.div) + 1
				// fmt.Printf("%0.0f\t%0.0f\t%0.4f\n", x, y, noise)
				// noise = noise%8

				if noiseFloat <= 1 {
					//water

					tile.TilePrint = 1

					noiseFloat = 75
					//Blue
				} else if noiseFloat > 1 && noiseFloat <= 4 {
					//Beach
					tile.TilePrint = 2
					noiseFloat = 70
					//Green
				} else if noiseFloat > 4 && noiseFloat <= 6 {
					// 1 off beach
					tile.TilePrint = 2
					noiseFloat = 145
				} else if noiseFloat > 6 && noiseFloat <= 7 {
					//Mountanous
					tile.TilePrint = 3
					noiseFloat = 78
				} else if noiseFloat > 7 && noiseFloat <= 15 {
					noiseFloat = 3
					tile.TilePrint = 3
				} else if noiseFloat > 15 && noiseFloat <= 25 {
					noiseFloat = 23
					tile.TilePrint = 3
				} else if noiseFloat > 25 && noiseFloat <= 50 {
					noiseFloat = 250
					tile.TilePrint = 3
				} else {
					//Mountain
					tile.TilePrint = 4
					noiseFloat = 250
				}
				_ = char

				_ = termbox.SetOutputMode(termbox.Output256)

				tile.X = x
				tile.Y = y
				tile.Type = int(noiseFloat)
				// tile.TilePrint = 1
				tile.NoiseValue = noise
				theMap.Tiles[tileId] = tile
				// termbox.SetCell(x, y, char, termbox.Attribute(y), termbox.Attribute(tile.Type))

			}
			termbox.SetCell(x, y, char, termbox.Attribute(y), termbox.Attribute(tile.Type))
		}
	}
	termbox.Flush()
}

func generateNewMap() {
	//Get new Map ID
	newMap := Map{2, 2, 3, 50, 1, map[string]Tile{}}
	newMap.Tiles = make(map[string]Tile)

	max := len(allMaps.Maps)
	// fmt.Println(max)
	// panic("NO")
	allMaps.Maps[max+1] = newMap

	getMap(max+1, defaultWidth, defaultHeight, 0, 0)
}

func handler(w http.ResponseWriter, r *http.Request) {

	uri := r.URL.Path[1:]
	if basePath != "/" {
		uri = strings.Replace(uri, basePath[1:], "", -1)
	}

	mapId, err := strconv.Atoi(uri)
	if err != nil {

		// out, err := json.Marshal(allMaps)
		// if err != nil {
		// 	panic (err)
		// }
		// fmt.Fprintf(w, "%s", out)

		fmt.Fprintf(w, "Cannot use: %s!", r.URL.Path[1:])
		// fmt.Println(i1)
	} else {
		_, ok := allMaps.Maps[mapId]
		if !ok {
			//Clean Up
			// panic("Map does not exist");
			fmt.Fprintf(w, "No Map with id: %s", r.URL.Path[1:])
		} else {

			getMap(mapId, defaultWidth, defaultHeight, 0, 0)

			//Clean Up
			theMap2, _ := allMaps.Maps[mapId]

			out, err := json.Marshal(theMap2)
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "%s", out)
		}

	}

}

func main() {

	if "" != os.Getenv("GG_MAP_PATH") {
		basePath = os.Getenv("GG_MAP_PATH")
	}

	if "" != os.Getenv("GG_MAP_LISTEN") {
		baseListener = os.Getenv("GG_MAP_LISTEN")
	}

	fmt.Println("basePath: ", basePath)
	fmt.Println("baseListener: ", baseListener)
	http.HandleFunc(basePath, handler)
	allMaps.Maps = make(map[int]Map)
	newMap := Map{2, 2, 3, 50, 1, map[string]Tile{}}
	max := 0
	fmt.Println(max)
	allMaps.Maps[max+1] = newMap
	//Clean Up

	getMap(max+1, defaultWidth, defaultHeight, 0, 0)

	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		go http.ListenAndServe(baseListener, nil)
	} else {
		http.ListenAndServe(baseListener, nil)
	}

	// termbox.SetOut
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	//Clean Up

	// generateNewMap()
loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
			// default:
			// 	draw()
			// time.Sleep(10 * time.Millisecond)
		}
	}
}
