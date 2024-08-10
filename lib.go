package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	dice "git.cmcode.dev/cmcode/go-dicewarelib"
	"github.com/adrg/xdg"
)

// Used for the config file directory and other things.
const APP_NAME = "go-fltk-diceware"

// Name of the config file.
const CONFIG_FILE = "config.json"

// Wraps around log.Println() as well as adding activity to the
// activity text buffer. Always adds a newline to the activity buffer.
func Log(v ...any) {
	log.Println(v...)
}

// Wraps around log.Printf() as well as adding activity to the
// activity text buffer. Always adds a newline to the activity buffer.
func Logf(format string, v ...any) {
	format = fmt.Sprintf("%v\n", format)
	log.Printf(format, v...)
}

func (app *App) loadConfig() {
	var err error
	if app.configFilePath == "" {
		app.configFilePath, err = xdg.SearchConfigFile(path.Join(APP_NAME, CONFIG_FILE))
		if err != nil {
			log.Printf("failed to get xdg config dir: %v", err.Error())
		}
	}

	if app.configFilePath != "" {
		bac, err := os.ReadFile(app.configFilePath)
		if err != nil {
			log.Printf("config file not readable at %v", app.configFilePath)
		}

		err = json.Unmarshal(bac, app.conf)
		if err != nil {
			log.Printf("config file %v failed to parse: %v", app.configFilePath, err.Error())
		}

		log.Printf("loaded config from %v", app.configFilePath)
	} else {
		if xdg.ConfigHome != "" {
			app.configFilePath = path.Join(xdg.ConfigHome, APP_NAME, CONFIG_FILE)
			log.Printf("using %v for config file path", app.configFilePath)
		} else {
			log.Println("unable to automatically identify any suitable config dirs; configuration will not be saved")
		}
	}
}

func (app *App) saveConfig() error {
	if app.configFilePath == "" {
		return fmt.Errorf("received empty config filename")
	}

	if app.conf == nil {
		return fmt.Errorf("config was nil")
	}

	b, err := json.Marshal(app.conf)
	if err != nil {
		return fmt.Errorf("failed to marshal app config to yaml: %v", err.Error())
	} else {
		dir, _ := filepath.Split(app.configFilePath)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			return fmt.Errorf("failed to create app config parent dir %v: %v", dir, err.Error())
		}
		err = os.WriteFile(app.configFilePath, b, 0o644)
		if err != nil {
			return fmt.Errorf("failed to save app config to %v: %v", app.configFilePath, err.Error())
		}
	}

	return nil
}

// Initializes the diceware library with the stored word lists. Can be executed
// repeatedly.
func (app *App) initDice() {
	simple, scount := dice.GetWords(content, "words-simple.txt")
	app.words.Simple = &simple
	app.words.SimpleCount = scount
	if app.conf.Extra {
		complex, ccount := dice.GetWords(content, "words-complex.txt")
		app.words.Complex = &complex
		app.words.ComplexCount = ccount
	} else {
		app.words.Complex = &map[int]string{} // zero out the ram usage
		app.words.ComplexCount = 0
	}
	Logf("loaded %v simple words and %v complex words", app.words.SimpleCount, app.words.ComplexCount)
	// log.Println(dice.GeneratePassword(&app.words, 3, " ", 64, 20, false))
}
