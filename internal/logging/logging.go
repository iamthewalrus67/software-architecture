package logging

import (
	"fmt"
	"time"
)

func Log(msg string) {
	Logf("%s", msg)
}

func Logf(format string, a ...any) {
	currentTime := time.Now()
	fmt.Printf("[%s] %s\n", currentTime.Format("2006.01.02 15:04:05"), fmt.Sprintf(format, a...))
}
