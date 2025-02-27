package watchers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/gomp/database"
	"github.com/aditya-K2/gomp/render"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/ui/notify"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/gomp/views"
	"github.com/fhs/gompd/v2/mpd"
)

var (
	currentSong mpd.Attrs
	start       bool = true
	ctime       time.Time
	status      mpd.Attrs
)

func OnConfigChange() {
	render.DrawCover(currentSong, false)
}

func Skip() bool {
	skip := false
	if _status, err := client.Conn.Status(); err != nil {
		skip = true
	} else {
		if status["state"] == "pause" && _status["state"] == "play" {
			skip = true
		}
		status = _status
	}
	return skip
}

func Init() {
	config.OnConfigChange = OnConfigChange
	database.SetDB(config.Config.DBPath)
	database.Read()
	database.Start()
	if c, err := client.Conn.CurrentSong(); err != nil {
		notify.Send("Couldn't get current song from MPD")
	} else {
		currentSong = c
		database.Send(currentSong, time.Second*0)
	}
	render.Rendr = render.NewRenderer()
	ctime = time.Now()
}
func StartRectWatcher() {
	redrawInterval := config.Config.RedrawInterval

	// Wait Until the ImagePreviewer is drawn
	// Ensures that cover art is not drawn before the UI is rendered.
	// Ref Issue: #39
	drawCh := make(chan bool)
	go func() {
		for ui.ImgX == 0 && ui.ImgY == 0 {
			ui.ImgX, ui.ImgY, ui.ImgW, ui.ImgH = ui.Ui.ImagePreviewer.GetRect()
		}
		drawCh <- true
	}()

	go func() {
		// Waiting for the draw channel
		draw := <-drawCh
		if draw {
			go func() {
				for {
					ImgX, ImgY, ImgW, ImgH := ui.Ui.ImagePreviewer.GetRect()
					if start {
						fmt.Println(ui.ImgX, ui.ImgY, ui.ImgW, ui.ImgH)
						render.DrawCover(currentSong, true)
						start = false
					}
					if ImgX != ui.ImgX || ImgY != ui.ImgY ||
						ImgW != ui.ImgW || ImgH != ui.ImgH {
						ui.ImgX = ImgX
						ui.ImgY = ImgY
						ui.ImgW = ImgW
						ui.ImgH = ImgH
						render.DrawCover(currentSong, false)
					}
					time.Sleep(time.Millisecond * time.Duration(redrawInterval))
				}
			}()
		}
	}()
}

func StartPlaylistWatcher() {
	var err error
	if views.PView.Playlist == nil {
		if views.PView.Playlist, err = client.Conn.PlaylistInfo(-1, -1); err != nil {
			utils.Print("RED", "Watcher couldn't get the current Playlist.\n")
			panic(err)
		}
	}

	nt, addr := utils.GetNetwork(config.Config.NetworkType, config.Config.Port, config.Config.NetworkAddress)
	w, err := mpd.NewWatcher(nt, addr, "")
	if err != nil {
		utils.Print("RED", "Could Not Start Watcher.\n")
		utils.Print("GREEN", "Please check your MPD Info in config File.\n")
		panic(err)
	}

	go func() {
		for err := range w.Error {
			notify.Send(err.Error())
		}
	}()

	go func() {
		for subsystem := range w.Event {
			if subsystem == "playlist" {
				if views.PView.Playlist, err = client.Conn.PlaylistInfo(-1, -1); err != nil {
					utils.Print("RED", "Watcher couldn't get the current Playlist.\n")
					panic(err)
				}
			} else if subsystem == "player" {
				if c, cerr := client.Conn.CurrentSong(); cerr != nil {
					utils.Print("RED", "Watcher couldn't get the current Playlist.\n")
					panic(err)
				} else {
					_ctime := time.Now()
					if !Skip() {
						database.Send(currentSong, _ctime.Sub(ctime))
					}
					ctime = _ctime
					currentSong = c
					render.DrawCover(c, false)
				}
			}
		}
	}()
}

func StartMPListener() {
	views.MPView.FSlice = []string{}
	mch := make(chan database.SubPayload)
	database.Subscribe(mch)
	go func() {
		for {
			sp := <-mch
			views.MPView.FMap = sp.Fmap
			views.MPView.FSlice = sp.Slice
		}
	}()
}

func ProgressFunction() (string, string, string, float64) {
	_currentAttributes := currentSong
	var song, top, text string
	var percentage float64
	song = config.Config.Colors.PBarTrack.String() +
		_currentAttributes["Title"] + "[-:-:-] - " + config.Config.Colors.PBarArtist.String() +
		_currentAttributes["Artist"] + "\n"
	_status, err := client.Conn.Status()
	el, err1 := strconv.ParseFloat(_status["elapsed"], 8)
	du, err := strconv.ParseFloat(_status["duration"], 8)
	top = fmt.Sprintf("[[::i] %s [-:-:-]Shuffle: %s Repeat: %s Volume: %s ]",
		utils.FormatString(_status["state"]),
		utils.FormatString(_status["random"]),
		utils.FormatString(_status["repeat"]),
		_status["volume"])
	if du != 0 {
		percentage = el / du * 100
		if err == nil && err1 == nil {
			text = utils.StrTime(el) + "/" + utils.StrTime(du) +
				"(" + strconv.FormatFloat(percentage, 'f', 2, 32) + "%" + ")"
		} else {
			text = ""
		}
	} else {
		text = "   ---:---"
		percentage = 0
	}
	if percentage > 100 {
		percentage = 0
	}
	return song, top, text, percentage
}

func DBCheck() {
	database.Add(currentSong["file"], time.Now().Sub(ctime))
	database.Write()
}
