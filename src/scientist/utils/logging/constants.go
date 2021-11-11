package logging

import "github.com/fatih/color"

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var defaltLevelNames = []string{
	"[DEBUG] ",
	"[INFO] ",
	"[WARN] ",
	"[ERROR] ",
}

var colorLevelNames = []string{
	color.CyanString("[DEBUG] "),
	color.GreenString("[INFO] "),
	color.YellowString("[WARN] "),
	color.RedString("[ERROR] "),
}
