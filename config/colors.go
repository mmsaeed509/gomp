package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/aditya-K2/gomp/utils"
	"github.com/gdamore/tcell/v2"
)

var (
	ColorError = func(s string) {
		_s := fmt.Sprintf("Wrong Color Provided: %s", s)
		utils.Print("RED", _s)
		os.Exit(-1)
	}
	DColors = map[string]tcell.Color{
		"Black":   tcell.ColorBlack,
		"Maroon":  tcell.ColorMaroon,
		"Green":   tcell.ColorGreen,
		"Olive":   tcell.ColorOlive,
		"Navy":    tcell.ColorNavy,
		"Purple":  tcell.ColorPurple,
		"Teal":    tcell.ColorTeal,
		"Silver":  tcell.ColorSilver,
		"Gray":    tcell.ColorGray,
		"Red":     tcell.ColorRed,
		"Lime":    tcell.ColorLime,
		"Yellow":  tcell.ColorYellow,
		"Blue":    tcell.ColorBlue,
		"Fuchsia": tcell.ColorFuchsia,
		"Aqua":    tcell.ColorAqua,
		"White":   tcell.ColorWhite,
	}
)

type Color struct {
	Foreground string `mapstructure:"foreground"`
	Bold       bool   `mapstructure:"bold"`
	Italic     bool   `mapstructure:"italic"`
}

type Colors struct {
	Artist        Color `mapstructure:"artist"`
	Album         Color `mapstructure:"album"`
	Track         Color `mapstructure:"track"`
	File          Color `mapstructure:"file"`
	Folder        Color `mapstructure:"folder"`
	Timestamp     Color `mapstructure:"timestamp"`
	MatchedTitle  Color `mapstructure:"matched_title"`
	MatchedFolder Color `mapstructure:"matched_folder"`
	PBarArtist    Color `mapstructure:"pbar_artist"`
	PBarTrack     Color `mapstructure:"pbar_track"`
	Null          Color
}

func (c Color) Color() tcell.Color {
	if strings.HasPrefix(c.Foreground, "#") && len(c.Foreground) == 7 {
		return tcell.GetColor(c.Foreground)
	} else if val, ok := DColors[c.Foreground]; ok {
		return val
	} else {
		ColorError(c.Foreground)
		return tcell.ColorBlack
	}
}

func (c Color) String() string {
	style := ""
	if c.Bold {
		style += "b"
	}
	if c.Italic {
		style += "i"
	}
	checkColor := func(s string) string {
		var res string
		if _, ok := DColors[s]; ok {
			res = strings.ToLower(s)
		} else if strings.HasPrefix(s, "#") && len(s) == 7 {
			res = s
		} else {
			ColorError(s)
		}
		return res
	}
	foreground := checkColor(c.Foreground)
	return fmt.Sprintf("[%s::%s]", foreground, style)
}

func NewColors() *Colors {
	return &Colors{
		Artist: Color{
			Foreground: "Purple",
			Bold:       false,
			Italic:     false,
		},
		Album: Color{
			Foreground: "Yellow",
			Bold:       false,
			Italic:     false,
		},
		Track: Color{
			Foreground: "Green",
			Bold:       false,
			Italic:     false,
		},
		Timestamp: Color{
			Foreground: "Red",
			Bold:       false,
			Italic:     true,
		},
		File: Color{
			Foreground: "Blue",
			Bold:       true,
			Italic:     false,
		},
		Folder: Color{
			Foreground: "Yellow",
			Bold:       true,
			Italic:     false,
		},
		MatchedFolder: Color{
			Foreground: "Blue",
			Bold:       true,
			Italic:     true,
		},
		MatchedTitle: Color{
			Foreground: "Yellow",
			Bold:       true,
			Italic:     true,
		},
		PBarArtist: Color{
			Foreground: "Blue",
			Bold:       true,
			Italic:     false,
		},
		PBarTrack: Color{
			Foreground: "Green",
			Bold:       true,
			Italic:     true,
		},
		Null: Color{
			Foreground: "White",
			Bold:       true,
			Italic:     false,
		},
	}
}
