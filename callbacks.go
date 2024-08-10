package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	dice "git.cmcode.dev/cmcode/go-dicewarelib"
	"github.com/atotto/clipboard"
	"github.com/pwiecz/go-fltk"
)

func (app *App) setCallbacks() {
	app.ui.win.SetCallback(app.gracefulExit)
	app.ui.win.SetResizeHandler(func() { app.ui.responsive() })

	app.ui.menu.AddEx("Copy", fltk.CTRL+fltk.SHIFT+'c', app.copy, 0)
	app.ui.menu.AddEx("Quit", fltk.CTRL+'q', app.gracefulExit, 0)
	app.ui.menu.AddEx("Generate 1", fltk.CTRL+'r', app.gen, 0)
	app.ui.menu.AddEx("Generate 2", fltk.CTRL+fltk.ENTER_KEY, app.gen, 0)
	app.ui.menu.AddEx("Help", fltk.F1, app.help, 0)

	app.darkCB()
	app.genCB()
	app.extraCB()
	app.sepCB()
	app.minCB()
	app.maxCB()
	app.wcCB()
}

func (app *App) help() {
	fltk.MessageBox("Help", "Generates relatively secure passwords that meet most website requirements.\nKeyboard shortcuts:\nCtrl+Shift+C: Copy to clipboard\nCtrl+R and Ctrl+Enter: Generate new password\nCtrl+Q: Quit\nF1: Help")
}

// Enables/disables dark mode.
func (app *App) darkCB() {
	app.ui.dark.SetCallback(func() {
		app.conf.DarkMode = !app.conf.DarkMode
		if !app.ui.darkModeChanged {
			fltk.MessageBox("App Restart Required", "In order for this to take effect, the theme change won't take effect until the application is restarted.")

			app.ui.darkModeChanged = true
		}
		app.ui.dark.SetValue(app.conf.DarkMode)
	})
}

// Enables/disables extra word usage.
func (app *App) extraCB() {
	app.ui.extra.SetCallback(func() {
		if !app.conf.Extra {
			confirmed := fltk.ChoiceDialog("Loading the extra words into memory will increase the RAM usage of this program. Proceed?", "Yes", "Cancel")
			if confirmed == 1 {
				app.ui.extra.SetValue(app.conf.Extra)
				return
			}
		}
		app.conf.Extra = !app.conf.Extra

		// fltk.MessageBox("App Restart Required", "In order for this to take effect, this application must be restarted.")
		app.ui.extra.SetValue(app.conf.Extra)
		app.initDice()
	})
}

// Copies the last-shown output value to the clipboard.
func (app *App) copy() {
	v := app.ui.out.Value()
	if v == "" {
		return
	}

	err := clipboard.WriteAll(v)
	if err != nil {
		log.Printf("failed to copy value to clipboard: %v", err.Error())
		return
	}

	app.ui.log.SetValue(fmt.Sprintf("Copied password with length %v to clipboard", len(v)))
}

// A standalone function that generates passwords according to the requirements.
func (app *App) gen() {
	r := dice.GeneratePassword(
		&app.words,
		app.conf.WordCount,
		app.conf.Separator,
		app.conf.MaxLen,
		app.conf.MinLen,
		app.conf.Extra,
	)
	app.ui.out.SetValue(r)
	app.ui.log.SetValue(fmt.Sprintf("Currently generated password length: %v", len(r)))
}

// Generates passwords according to the requirements when the "Generate" button
// is clicked.
func (app *App) genCB() {
	app.ui.gen.SetCallback(func() { app.gen() })
}

// Updates the separator when the user changes the separator input field.
func (app *App) sepCB() {
	app.ui.sep.SetCallback(func() { app.conf.Separator = app.ui.sep.Value() })
}

// Updates the min length when the user changes the min input field.
func (app *App) minCB() {
	app.ui.min.SetCallback(func() {
		m := app.ui.min.Value()
		if m == "" {
			return
		}
		i, err := strconv.ParseInt(m, 10, 64)
		if err != nil {
			return
		}
		app.conf.MinLen = int(i)
	})
}

// Updates the max length when the user changes the max input field.
func (app *App) maxCB() {
	app.ui.max.SetCallback(func() {
		m := app.ui.max.Value()
		if m == "" {
			return
		}
		i, err := strconv.ParseInt(m, 10, 64)
		if err != nil {
			return
		}
		app.conf.MaxLen = int(i)
	})
}

// Updates the word count when the user changes the word count input field.
func (app *App) wcCB() {
	app.ui.wc.SetCallback(func() {
		m := app.ui.wc.Value()
		if m == "" {
			return
		}
		i, err := strconv.ParseInt(m, 10, 64)
		if err != nil {
			return
		}
		app.conf.WordCount = int(i)
	})
}

// Called when the app attempts to exit, such as the window closing or ctrl+c
// interrupt signal on the command line.
func (app *App) gracefulExit() {
	Log("closing app and saving config, please wait a moment...")
	err := app.saveConfig()
	if err != nil {
		log.Printf("failed to save config: %v", err.Error())
	}

	Log("done, exiting now.")
	os.Exit(0)
}
