package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func readInput(filename string) (puzzle, words []string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	lastlen := -1
	wordsnow := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Replace(scanner.Text(), " ", "", -1)
		if wordsnow { // Puzzle is done, all lines now are words
			words = append(words, line)
			continue
		}
		// if last line is different size then this line it must be words
		if lastlen == len(line) || lastlen == -1 {
			puzzle = append(puzzle, line)
		} else {
			words = append(words, line)
			wordsnow = true
		}
		lastlen = len(line)
	}
	return
}

func findWord(findme string, puzzle []string, output [][]byte) [][]byte {
	for y, row := range puzzle {
		for x := range row {

			// look up right
			type direction struct{ xoff, yoff int }
			d := []direction{
				// right
				{xoff: 1, yoff: 0},
				// left
				{xoff: -1, yoff: 0},
				// down
				{xoff: 0, yoff: 1},
				// up
				{xoff: 0, yoff: -1},
				// down right
				{xoff: 1, yoff: 1},
				// down left
				{xoff: -1, yoff: 1},
				// up right
				{xoff: 1, yoff: -1},
				// up left
				{xoff: -1, yoff: -1},
			}
			for _, dir := range d {
				for li, letter := range findme {
					var xoff, yoff int
					// Ugly ass nested ifs to figure what direction to go
					if dir.xoff == 1 {
						xoff = x + li
					} else if dir.xoff == -1 {
						xoff = x - li
					} else if dir.xoff == 0 {
						xoff = x
					}
					if dir.yoff == 1 {
						yoff = y + li
					} else if dir.yoff == -1 {
						yoff = y - li
					} else if dir.yoff == 0 {
						yoff = y
					}
					// The word will not fit in the puzzle stop looking
					if xoff >= len(puzzle) || xoff < 0 {
						break
					}
					// The word will not fit in the puzzle stop looking
					if yoff >= len(row) || yoff < 0 {
						break
					}
					// The letter doesn't match stop looking
					if rune(puzzle[yoff][xoff]) != letter {
						break
					}
					// We found all the letters in this direction!
					if li == len(findme)-1 { // FOUND IT
						// Work backwards and add them to the output
						for i := range findme {
							var nxoff, nyoff int
							if dir.xoff == 1 {
								nxoff = xoff - i
							} else if dir.xoff == -1 {
								nxoff = xoff + i
							} else if dir.xoff == 0 {
								nxoff = xoff
							}
							if dir.yoff == 1 {
								nyoff = yoff - i
							} else if dir.yoff == -1 {
								nyoff = yoff + i
							} else if dir.yoff == 0 {
								nyoff = yoff
							}
							output[nyoff][nxoff] = puzzle[nyoff][nxoff]
						}
						return output
					}
				}
			}

		}
	}
	return output
}

func main() {
	puzzle, words, err := readInput("/tmp/data1")
	if err != nil {
		log.Fatal(err)
	}
	_ = words

	// Make an empty output
	output := make([][]byte, len(puzzle))
	for y, row := range puzzle {
		output[y] = make([]byte, len(row))
		for x := range output[y] {
			output[y][x] = []byte(" ")[0]
		}
	}

	for _, findme := range words {
		// We reuse the output adding more to it each time
		// This probably fucks running each word as a go routine
		output = findWord(findme, puzzle, output)
	}

	// Display output
	for _, o := range output {
		for _, c := range o {
			//Output has spaces between each letter
			fmt.Printf("%s ", string(c))
		}
		fmt.Printf("\n")
	}
}
