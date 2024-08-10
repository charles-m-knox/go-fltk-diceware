package main

import (
	"embed"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/pwiecz/go-fltk"

	dice "git.cmcode.dev/cmcode/go-dicewarelib"
)

//go:embed words-simple.txt
//go:embed words-complex.txt
var content embed.FS

var (
	forcePortrait  bool
	forceLandscape bool
	app            App = App{
		conf:  &AppConfig{},
		ui:    &UI{},
		words: dice.Words{},
	}
)

type App struct {
	conf *AppConfig
	ui   *UI
	// Data is stored between runs of this application in this yml config file.
	configFilePath string
	words          dice.Words
}

type AppConfig struct {
	DarkMode bool `json:"darkMode"`
	// If true, uses an extended word list
	Extra     bool   `json:"useExtendedWordList"`
	MaxLen    int    `json:"maxLen"`
	MinLen    int    `json:"minLen"`
	Separator string `json:"separator"`
	WordCount int    `json:"wordCount"`
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
	flag.Parse()
}

func main() {
	parseFlags()

	app.initDice()

	app.loadConfig()
	app.initUI()
	app.ui.theme(app.conf.DarkMode)
	app.ui.responsive()
	app.ui.upsize()
	app.setCallbacks()
	app.ui.win.End()
	app.ui.win.Show()
	go fltk.Run()

	// Create a channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-signalChan

	app.gracefulExit()
}
