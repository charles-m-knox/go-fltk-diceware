package main

import (
	"fmt"
	"log"
	"math"

	"github.com/pwiecz/go-fltk"
)

// If the screen is portrait or landscape, the window will be scaled
// accordingly.
const (
	WIDTH_PORTRAIT   = 100
	HEIGHT_PORTRAIT  = 150
	WIDTH_LANDSCAPE  = 150
	HEIGHT_LANDSCAPE = 100
)

// Positioning (x,y,w,h) for fltk elements
type pos struct {
	X int
	Y int
	W int
	H int
	// A reference to the UI is needed in order to perform the translation
	// according to the ui's specs
	ui *UI
}

// Buttons, inputs, widgets, etc that need to be repositioned in a
// responsive manner.
type UI struct {
	win *fltk.Window // main window

	menu *fltk.MenuBar // hidden menu bar for shortcut keys

	dark  *fltk.CheckButton // dark mode checkbox
	extra *fltk.CheckButton // "use extra words" checkbox
	max   *fltk.Input       // max output length
	min   *fltk.Input       // min output length
	out   *fltk.Input       // generated output input field
	sep   *fltk.Input       // separator character input field
	wc    *fltk.Input       // word count input field
	log   *fltk.HelpView    // shows word count and generated word length
	gen   *fltk.Button      // generate button

	// winp   pos // main window position
	darkp  pos // dark mode checkbox position
	extrap pos // "use extra words" checkbox position
	maxp   pos // max output length position
	minp   pos // min output length position
	outp   pos // generated output input field position
	sepp   pos // separator character input field position
	wcp    pos // word count input field position
	logp   pos // shows word count and generated word length (position)
	genp   pos // generate button position

	portrait        bool // portrait mode or landscape mode
	darkModeChanged bool // if true, prompts to restart after changing dark mode will not show
}

// isPortrait returns true if the screen is taller than it is wide. It returns
// false otherwise, including for square screens.
func isPortrait() (bool, error) {
	_, _, width, height := fltk.ScreenWorkArea(int(fltk.SCREEN))

	if width == 0 || height == 0 {
		return false, fmt.Errorf("received 0 for one of screen height or width")
	}

	if width > height {
		return false, nil
	}

	return true, nil
}

// Translates the widget's width/height from the original 100 or 150px base
// width/height to the window's current width/height
func (ui *UI) tr(i int, winw int, winh int, useHeight bool) int {
	if ui.portrait {
		if useHeight {
			return int(math.Round((float64(i) / float64(HEIGHT_PORTRAIT)) * float64(winh)))
		} else {
			return int(math.Round((float64(i) / float64(WIDTH_PORTRAIT)) * float64(winw)))
		}
	} else {
		if useHeight {
			return int(math.Round((float64(i) / float64(HEIGHT_LANDSCAPE)) * float64(winh)))
		} else {
			return int(math.Round((float64(i) / float64(WIDTH_LANDSCAPE)) * float64(winw)))
		}
	}
}

// Translate converts a predefined position into a scaled position based on
// the latest width & height of the window.
func (p *pos) Translate(winw, winh int) {
	p.X = p.ui.tr(p.X, winw, winh, false)
	p.Y = p.ui.tr(p.Y, winw, winh, true)
	p.W = p.ui.tr(p.W, winw, winh, false)
	p.H = p.ui.tr(p.H, winw, winh, true)
}

// Initializes the UI for the app. Call this once, only after the app config has
// been loaded.
func (app *App) initUI() {
	// fltk.SetScheme("gtk+")
	fltk.InitStyles()
	fltk.SetTooltipDelay(0.1)
	fltk.EnableTooltips()

	var err error
	app.ui.portrait, err = isPortrait()
	if err != nil {
		log.Fatalf("failed to determine screen size: %v", err.Error())
	}

	// probably could write this more intelligently later
	winw := WIDTH_LANDSCAPE
	winh := HEIGHT_LANDSCAPE
	if app.ui.portrait || forcePortrait {
		winw = WIDTH_PORTRAIT
		winh = HEIGHT_PORTRAIT
		app.ui.portrait = true
	}

	if forceLandscape {
		winw = WIDTH_LANDSCAPE
		winh = HEIGHT_LANDSCAPE
		app.ui.portrait = false
	}

	// initialize all buttons and widgets
	app.ui.win = fltk.NewWindow(winw, winh, "Diceware Password Generator FLTK")
	app.ui.menu = fltk.NewMenuBar(0, 0, 0, 0)
	app.ui.dark = fltk.NewCheckButton(0, 0, 0, 0, "&Dark Mode")
	app.ui.extra = fltk.NewCheckButton(0, 0, 0, 0, "&Extra Words")
	app.ui.max = fltk.NewInput(0, 0, 0, 0, "&Max Length")
	app.ui.min = fltk.NewInput(0, 0, 0, 0, "Mi&n Length")
	app.ui.out = fltk.NewInput(0, 0, 0, 0, "&Output")
	app.ui.sep = fltk.NewInput(0, 0, 0, 0, "&Separator")
	app.ui.wc = fltk.NewInput(0, 0, 0, 0, "&Word Count")
	app.ui.log = fltk.NewHelpView(0, 0, 0, 0, "")
	app.ui.gen = fltk.NewButton(0, 0, 0, 0, "&Generate")

	// propagate default values from config to widgets that accept them
	app.ui.dark.SetValue(app.conf.DarkMode)
	app.ui.extra.SetValue(app.conf.Extra)
	app.ui.max.SetValue(fmt.Sprint(app.conf.MaxLen))
	app.ui.min.SetValue(fmt.Sprint(app.conf.MinLen))
	app.ui.sep.SetValue(app.conf.Separator)
	app.ui.wc.SetValue(fmt.Sprint(app.conf.WordCount))

	// app.ui.dark.SetAlign(fltk.ALIGN_TOP_LEFT)
	// app.ui.extra.SetAlign(fltk.ALIGN_TOP_LEFT)
	app.ui.out.SetAlign(fltk.ALIGN_TOP_LEFT)
	app.ui.max.SetAlign(fltk.ALIGN_TOP_LEFT)
	app.ui.min.SetAlign(fltk.ALIGN_TOP_LEFT)
	app.ui.sep.SetAlign(fltk.ALIGN_TOP_LEFT)
	app.ui.wc.SetAlign(fltk.ALIGN_TOP_LEFT)
	app.ui.log.SetAlign(fltk.ALIGN_TOP_LEFT)

	app.ui.log.SetLabelSize(10)
	app.ui.log.SetLabelFont(fltk.HELVETICA)
	app.ui.log.SetValue("Output will go here")

	app.ui.dark.SetTooltip("Toggling the UI mode requires a restart, and this setting will persist to settings between app restarts.")
	app.ui.extra.SetTooltip("If enabled, a more complex word list will be used, with significantly more dictionary words to use. This is more secure, but some words may be too difficult to work with.")
	app.ui.max.SetTooltip("The maximum permissible number of characters to generate. Default=64")
	app.ui.min.SetTooltip("The minimum permissible number of characters to generate. Default=20")
	app.ui.out.SetTooltip("Generated passwords will appear here.")
	app.ui.sep.SetTooltip("The separator to place between generated words. Default is a space character. Multiple characters can be used.")
	app.ui.wc.SetTooltip("The number of words to generate. This may require experimenting with min/max length. Default=3")
	app.ui.gen.SetTooltip("Press this button to generate a password with the above settings.")

	app.ui.win.Resizable(app.ui.win)
	app.ui.win.SetXClass("gfltkdice")
}

// Sizes the window to 3x the design size, which is intentionally
// small. Use this after things have been initiated.
func (ui *UI) upsize() {
	if app.ui.portrait {
		app.ui.win.Resize(0, 0, WIDTH_PORTRAIT*3, HEIGHT_PORTRAIT*3)
	} else {
		app.ui.win.Resize(0, 0, WIDTH_LANDSCAPE*3, HEIGHT_LANDSCAPE*3)
	}
}

// Resizes and repositions all components based on the window's size.
func (ui *UI) responsive() {
	if forceLandscape || forcePortrait {
		return
	}

	winw := ui.win.W()
	winh := ui.win.H()

	if winw > winh {
		ui.portrait = false
	} else {
		ui.portrait = true
	}

	if ui.portrait {
		ui.darkp = pos{X: 50, Y: 65, W: 45, H: 15, ui: ui}
		ui.extrap = pos{X: 5, Y: 65, W: 40, H: 15, ui: ui}
		ui.genp = pos{X: 5, Y: 125, W: 90, H: 20, ui: ui}
		ui.logp = pos{X: 5, Y: 85, W: 90, H: 35, ui: ui}
		ui.maxp = pos{X: 50, Y: 45, W: 45, H: 15, ui: ui}
		ui.minp = pos{X: 5, Y: 45, W: 40, H: 15, ui: ui}
		ui.outp = pos{X: 5, Y: 5, W: 90, H: 15, ui: ui}
		ui.sepp = pos{X: 5, Y: 25, W: 40, H: 15, ui: ui}
		ui.wcp = pos{X: 50, Y: 25, W: 45, H: 15, ui: ui}
	} else {
		// landscape
		ui.darkp = pos{X: 80, Y: 45, W: 65, H: 15, ui: ui}
		ui.extrap = pos{X: 5, Y: 45, W: 70, H: 15, ui: ui}
		ui.genp = pos{X: 5, Y: 85, W: 140, H: 10, ui: ui}
		ui.logp = pos{X: 5, Y: 65, W: 140, H: 15, ui: ui}
		ui.maxp = pos{X: 120, Y: 25, W: 25, H: 15, ui: ui}
		ui.minp = pos{X: 80, Y: 25, W: 35, H: 15, ui: ui}
		ui.outp = pos{X: 5, Y: 5, W: 140, H: 15, ui: ui}
		ui.sepp = pos{X: 5, Y: 25, W: 35, H: 15, ui: ui}
		ui.wcp = pos{X: 45, Y: 25, W: 30, H: 15, ui: ui}
	}

	ui.darkp.Translate(winw, winh)
	ui.extrap.Translate(winw, winh)
	ui.maxp.Translate(winw, winh)
	ui.minp.Translate(winw, winh)
	ui.outp.Translate(winw, winh)
	ui.sepp.Translate(winw, winh)
	ui.wcp.Translate(winw, winh)
	ui.logp.Translate(winw, winh)
	ui.genp.Translate(winw, winh)

	ui.dark.Resize(ui.darkp.X, ui.darkp.Y, ui.darkp.W, ui.darkp.H)
	ui.extra.Resize(ui.extrap.X, ui.extrap.Y, ui.extrap.W, ui.extrap.H)
	ui.max.Resize(ui.maxp.X, ui.maxp.Y, ui.maxp.W, ui.maxp.H)
	ui.min.Resize(ui.minp.X, ui.minp.Y, ui.minp.W, ui.minp.H)
	ui.out.Resize(ui.outp.X, ui.outp.Y, ui.outp.W, ui.outp.H)
	ui.sep.Resize(ui.sepp.X, ui.sepp.Y, ui.sepp.W, ui.sepp.H)
	ui.wc.Resize(ui.wcp.X, ui.wcp.Y, ui.wcp.W, ui.wcp.H)
	ui.log.Resize(ui.logp.X, ui.logp.Y, ui.logp.W, ui.logp.H)
	ui.gen.Resize(ui.genp.X, ui.genp.Y, ui.genp.W, ui.genp.H)
}

const (
	DARK_COLOR_TEXT               fltk.Color = 0x9f9f9f00
	DARK_COLOR_INPUT_BG           fltk.Color = 0x20202000
	DARK_COLOR_INPUT_SELECTED_BG  fltk.Color = 0xafafaf00
	LIGHT_COLOR_TEXT              fltk.Color = 0x20030500
	LIGHT_COLOR_INPUT_BG          fltk.Color = 0xFFFFFF00
	LIGHT_COLOR_INPUT_SELECTED_BG fltk.Color = 0x00008000
)

var (
	COLOR_TEXT              fltk.Color = LIGHT_COLOR_TEXT
	COLOR_INPUT_BG          fltk.Color = LIGHT_COLOR_INPUT_BG
	COLOR_INPUT_SELECTED_BG fltk.Color = LIGHT_COLOR_INPUT_SELECTED_BG
)

// Changes the color of various widgets/states to light or dark mode.
func (ui *UI) theme(dark bool) {
	if dark {
		log.Println("dark mode activated")
		COLOR_TEXT = DARK_COLOR_TEXT
		COLOR_INPUT_BG = DARK_COLOR_INPUT_BG
		COLOR_INPUT_SELECTED_BG = DARK_COLOR_INPUT_SELECTED_BG
		fltk.SetForegroundColor(230, 230, 230)
		fltk.SetBackgroundColor(40, 40, 40)
	} else {
		log.Println("light mode activated")
		COLOR_TEXT = LIGHT_COLOR_TEXT
		COLOR_INPUT_BG = LIGHT_COLOR_INPUT_BG
		COLOR_INPUT_SELECTED_BG = LIGHT_COLOR_INPUT_SELECTED_BG
		fltk.SetBackgroundColor(192, 192, 192)
		fltk.SetForegroundColor(0, 0, 0)
		return
	}

	ui.dark.SetLabelColor(COLOR_TEXT)
	ui.extra.SetLabelColor(COLOR_TEXT)
	ui.max.SetLabelColor(COLOR_TEXT)
	ui.min.SetLabelColor(COLOR_TEXT)
	ui.out.SetLabelColor(COLOR_TEXT)
	ui.sep.SetLabelColor(COLOR_TEXT)
	ui.wc.SetLabelColor(COLOR_TEXT)
	ui.log.SetLabelColor(COLOR_TEXT)
	ui.gen.SetLabelColor(COLOR_TEXT)

	ui.dark.SetColor(COLOR_INPUT_BG)
	ui.extra.SetColor(COLOR_INPUT_BG)
	ui.max.SetColor(COLOR_INPUT_BG)
	ui.min.SetColor(COLOR_INPUT_BG)
	ui.out.SetColor(COLOR_INPUT_BG)
	ui.sep.SetColor(COLOR_INPUT_BG)
	ui.wc.SetColor(COLOR_INPUT_BG)
	ui.log.SetColor(COLOR_INPUT_BG)
	ui.gen.SetColor(COLOR_INPUT_BG)

	ui.dark.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.extra.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.max.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.min.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.out.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.sep.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.wc.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.log.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.gen.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
}
