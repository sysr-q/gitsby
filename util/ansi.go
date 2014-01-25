package util

import (
	"fmt"
	"log"
	"github.com/mgutz/ansi"
)

var success string = ansi.ColorCode("green")
var failure string = ansi.ColorCode("red")
var info string = ansi.ColorCode("blue")
var reset string = ansi.ColorCode("reset")

// Raw colors!
func Green(str string) {
	fmt.Println(success + str + reset)
}

func Red(str string) {
	fmt.Println(failure + str + reset)
}

func Blue(str string) {
	fmt.Println(info + str + reset)
}

// Handy helpers
func Success(str string) {
	log.Printf("%s%c %s%s", success, '\u2714', str, reset)
}

func Successf(format string, bits ...interface{}) {
	Success(fmt.Sprintf(format, bits...))
}

func Failure(str string) {
	log.Printf("%s%c %s%s", failure, '\u2716', str, reset)
}

func Failuref(format string, bits ...interface{}) {
	Failure(fmt.Sprintf(format, bits...))
}

func Info(str string) {
	log.Printf("%s%c %s%s", info, '\u279C', str, reset)
}

func Infof(format string, bits ...interface{}) {
	Info(fmt.Sprintf(format, bits...))
}
