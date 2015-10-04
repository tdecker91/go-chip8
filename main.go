package main

import (
	"fmt"
	"os"
	"flag"
	"tysondecker.com/tdecker/chip8/chip8"
	"tysondecker.com/tdecker/chip8/screen"
)

/**
 * Prints the given message and exits the program
 * @param  {string}
 * @return {nil}
 */
func error(message string) {
	fmt.Println(message)
	os.Exit(1)
}

/**
 * Parses the command line arguments for the romfile.
 * prints the help messages if -h/-help is passed in
 * @return {string} path to romfile
 */
func parseArgs() string {

	romPathPtr := flag.String("rom", "", "Required: path to the romfile. (ex -rom=\"path/to/rom.file\")")
	helpPtr := flag.Bool("h", false, "prints the usage")

	flag.Parse()

	if *helpPtr || *romPathPtr == "" {
		flag.Usage()
		os.Exit(1)
	}

	return *romPathPtr

}

func main() {

	romPath := parseArgs()
	chip := chip8.NewChip8()

	chip.LoadRom(romPath)
	display := new(screen.Screen)
	display.Init()


	// Emulation Loop
	for chip.Running {
		
		chip.EmulateCycle()

		if chip.DrawFlag {
			chip.DumpScreen()
		}

		// Store keypress...

	}
}