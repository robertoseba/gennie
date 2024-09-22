package output

type Color string

const (
	Reset   Color = "\033[0m"
	Red     Color = "\033[31m"
	Green   Color = "\033[32m"
	Yellow  Color = "\033[33m"
	Blue    Color = "\033[34m"
	Magenta Color = "\033[35m"
	Cyan    Color = "\033[36m"
	Gray    Color = "\033[37m"
	White   Color = "\033[97m"

	BgWhite  Color = "\033[107m"
	BgCyan   Color = "\033[46m"
	BgYellow Color = "\033[43m"
)

func (c Color) String() string {
	return string(c)
}
