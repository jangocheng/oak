package oak

import (
	"image"
	"os"
	"path/filepath"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
)

var (

	// The init channel communicates between
	// initializing goroutines for when significant
	// steps in initialization have been reached
	// initCh = make(chan bool)
	// currently unused

	//
	transitionCh = make(chan bool)

	// The Scene channel receives a signal
	// when a scene's .loop() function should
	// be called.
	sceneCh = make(chan bool)

	// The skip scene channel receives a debug
	// signal to forcibly go to the next
	// scene.
	skipSceneCh = make(chan bool)

	// The quit channel receives a signal when
	// the program should stop.
	quitCh = make(chan bool)

	// The draw channel receives a signal when
	// drawing should cease (or resume)
	drawCh = make(chan bool)

	// The debug reset channel represents
	// when the debug console should forget the
	// commands that have been sent to it.
	debugResetCh = make(chan bool)

	// The viewport channel controls when new
	// viewport positions should be drawn
	viewportCh = make(chan [2]int)

	debugResetInProgress bool

	// ScreenWidth is the width of the screen
	ScreenWidth int
	// ScreenHeight is the height of the screen
	ScreenHeight int

	// FrameRate is the current logical frame rate.
	// Changing this won't directly effect frame rate, that
	// requires changing the LogicTicker, but it will take
	// effect next scene
	FrameRate int

	// DrawFrameRate is the equivalent to FrameRate for
	// the rate at which the screen is drawn.
	DrawFrameRate int

	zeroPoint = image.Point{0, 0}
)

// Init initializes the oak engine.
// It spawns off an event loop of several goroutines
// and loops through scenes after initialization.
func Init(firstScene string) {
	dlog.SetLogger(dlog.NewLogger())
	dlog.CreateLogFile()

	initConf()

	// Set variables from conf file
	lvl, err := dlog.ParseDebugLevel(conf.Debug.Level)
	dlog.SetDebugLevel(lvl)
	// We are intentionally using the lvl value before checking error,
	// because we can only log errors through dlog itself anyway

	// We do this knowing that the default debug level when SetDebugLevel fails
	// is ERROR, so this will be recorded.
	dlog.ErrorCheck(err)
	dlog.SetDebugFilter(conf.Debug.Filter)

	dlog.Info("Oak Init Start")

	ScreenWidth = conf.Screen.Width
	ScreenHeight = conf.Screen.Height
	FrameRate = conf.FrameRate
	DrawFrameRate = conf.DrawFrameRate
	SetLang(conf.Language)

	wd, _ := os.Getwd()

	render.SetFontDefaults(wd, conf.Assets.AssetPath, conf.Assets.FontPath,
		conf.Font.Hinting, conf.Font.Color, conf.Font.File, conf.Font.Size,
		conf.Font.DPI)
	// END of loading variables from configuration

	SeedRNG(DefaultSeed)

	imageDir := filepath.Join(wd,
		conf.Assets.AssetPath,
		conf.Assets.ImagePath)
	audioDir := filepath.Join(wd,
		conf.Assets.AssetPath,
		conf.Assets.AudioPath)

	dlog.Info("Init Scene Loop")
	go sceneLoop(firstScene)
	dlog.Info("Init asset load")
	render.SetAssetPaths(imageDir)
	go loadAssets(imageDir, audioDir)
	dlog.Info("Init Console")
	go debugConsole(debugResetCh, skipSceneCh, os.Stdin)
	dlog.Info("Init Main Driver")
	InitDriver(lifecycleLoop)
}
