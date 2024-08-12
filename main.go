package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pwiecz/go-fltk"

	dice "git.cmcode.dev/cmcode/go-dicewarelib"
)

//go:embed words-simple.txt
//go:embed words-complex.txt
var content embed.FS

// The version of the application; set at build time via:
//
//	`go build -ldflags "-X main.version=1.2.3" main.go`
//
//nolint:revive
var version string = "dev"

// Flag for showing the version and subsequently quitting.
var flagVersion bool

var (
	// If true, the app will always be rendered in portrait mode
	forcePortrait bool
	// If true, the app will always be rendered in landscape mode
	forceLandscape bool
	// app contains the shared state that is required for the entire app to
	// function.
	app App = App{
		conf:  &AppConfig{},
		ui:    &UI{},
		words: dice.Words{},
	}
)

// App contains the shared state that is required for the entire app to
// function.
type App struct {
	// The configuration for the entire app, loaded and saved to the XDG config.
	conf *AppConfig
	// All of the UI elements for this app are contained in the UI struct.
	ui *UI
	// Data is stored between runs of this application in this yml config file.
	configFilePath string
	// The dictionary of diceware words, provided by the diceware lib.
	words dice.Words
}

type AppConfig struct {
	// If true, the app will start in dark mode
	DarkMode bool `json:"darkMode"`
	// If true, uses an extended word list
	Extra bool `json:"useExtendedWordList"`
	// The maximum permissible generated output length
	MaxLen int `json:"maxLen"`
	// The minimum permissible generated output length
	MinLen int `json:"minLen"`
	// The separator character (s) to place between generated words
	Separator string `json:"separator"`
	// The number of words to generate
	WordCount int `json:"wordCount"`
}

func parseFlags() {
	flag.BoolVar(&forcePortrait, "portrait", false, "force portrait orientation for the interface")
	flag.BoolVar(&forceLandscape, "landscape", false, "force landscape orientation for the interface")
	flag.StringVar(&app.configFilePath, "f", "", "the config file to write to, instead of the default provided by XDG config directories")
	flag.StringVar(&app.conf.Separator, "s", " ", "the character(s) to place between each word")
	flag.IntVar(&app.conf.MaxLen, "max", 64, "the longest permissible length of generated passwords")
	flag.IntVar(&app.conf.MinLen, "min", 20, "the least permissible length of generated passwords")
	flag.IntVar(&app.conf.WordCount, "wc", 3, "the number of words to generate")
	flag.BoolVar(&app.conf.Extra, "extra", false, "if true, more complicated permutations of words will be used")
	flag.BoolVar(&flagVersion, "v", false, "print version and exit")
	flag.Parse()
}

func main() {
	parseFlags()
	if flagVersion {
		//nolint:forbidigo
		fmt.Println(version)
		os.Exit(0)
	}

	app.loadConfig()
	app.initDice()
	app.initUI()
	app.ui.theme(app.conf.DarkMode)
	app.ui.responsive()
	app.ui.upsize()
	app.setCallbacks()
	app.ui.win.End()
	app.ui.win.Show()
	go fltk.Run()
	// start with an initial password populated in the output field
	app.gen()

	// Channel that receives OS signals, like ctrl+c to interrupt
	var sc chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	// Block until a signal is received
	<-sc

	app.gracefulExit()
}
