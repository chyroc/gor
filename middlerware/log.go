package middlerware

import (
	"fmt"
	"time"

	"github.com/Chyroc/gor"
)

// https://stackoverflow.com/questions/2616906/how-do-i-output-coloured-text-to-a-linux-terminal
func formatColor(color, msg interface{}) string {
	format := "\033[%dm%v\033[39m"
	if color == "red" {
		return fmt.Sprintf(format, 31, msg)
	} else if color == "green" {
		return fmt.Sprintf(format, 32, msg)
	} else if color == "yellow" {
		return fmt.Sprintf(format, 33, msg)
	} else if color == "gray" {
		return fmt.Sprintf(format, 90, msg)
	} else if color == "white" {
		return fmt.Sprintf(format, 97, msg)
	}
	return ""
}

func formatColorCode(code int) string {
	switch {
	case code >= 200 && code < 300:
		return formatColor("gray", code)
	case code >= 300 && code < 400:
		return formatColor("white", code)
	case code >= 400 && code < 500:
		return formatColor("yellow", code)
	default:
		return formatColor("red", code)
	}
}

// Logger default gor Logger middleware
var Logger = func(req *gor.Req, res *gor.Res, next gor.Next) {
	startTime := time.Now()
	next()
	fmt.Printf("✨ %s %s %s %s\n", formatColor("green", req.Method), req.OriginalURL, formatColorCode(res.StatusCode), time.Since(startTime))
}
