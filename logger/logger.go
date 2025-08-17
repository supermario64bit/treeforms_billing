package logger

import "log"

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"

	redBgWhiteText = "\033[37;41m"
)

func Success(msg string) {
	log.Println(green + "✔ " + msg + reset)
}

func Danger(msg string) {
	log.Println(red + "✘ " + msg + reset)
}

func Warning(msg string) {
	log.Println(yellow + "⚠ " + msg + reset)
}

func Info(msg string) {
	log.Println(blue + "! " + msg + reset)
}

func HighlightedDanger(msg string) {
	log.Println(redBgWhiteText + "! " + msg + reset)
}
