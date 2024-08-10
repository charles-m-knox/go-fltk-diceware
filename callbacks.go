package main

import (
	"log"
	"os"

	dice "git.cmcode.dev/cmcode/go-dicewarelib"
	"github.com/pwiecz/go-fltk"
)

func (app *App) setCallbacks() {
	app.ui.win.SetCallback(app.gracefulExit)
	app.ui.win.SetResizeHandler(func() { app.ui.responsive() })

	app.ui.menu.AddEx("Quit", fltk.CTRL+'q', app.gracefulExit, 0)

	app.darkCB()
	app.genCB()
	app.extraCB()
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

// Generates passwords according to the requirements when the "Generate" button
// is clicked.
func (app *App) genCB() {
	app.ui.gen.SetCallback(func() {
		r := dice.GeneratePassword(&app.words, 3, " ", 64, 20, false)
		app.ui.out.SetValue(r)
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
