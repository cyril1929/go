package main

import (
	"context"
	"fmt"
	goos "os"
	"os/signal"
	"path"
	"path/filepath"
	"time"

	"github.com/hack-pad/hackpadfs/os"
	"github.com/johnstarich/go/gowerline/internal/status"
)

const appName = "gowerline"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), goos.Interrupt)
	defer cancel()
	err := run(ctx)
	if err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	fs := os.NewFS()
	cacheDir, err := goos.UserCacheDir()
	if err != nil {
		return err
	}
	absoluteCacheDir, err := filepath.Abs(cacheDir)
	if err != nil {
		return err
	}
	cacheSubPath, err := fs.FromOSPath(absoluteCacheDir)
	if err != nil {
		return err
	}
	appCacheSubPath := path.Join(cacheSubPath, appName)
	cacheFS, err := fs.Sub(appCacheSubPath)
	if err != nil {
		return err
	}

	statusLine := status.Line{
		Segments: []status.Segment{
			{
				Font:            status.Font{Foreground: "#797aac", Background: "#121212"},
				GenerateContent: weatherStatus,
				Name:            "weather",
				Separator:       status.Separator{Font: status.Font{Foreground: "#121212"}},
			},
			{
				Font:            status.Font{Foreground: "#f3e6d8", Background: "#121212"},
				GenerateContent: batteryStatus,
				Name:            "battery",
				Separator:       status.Separator{Font: status.Font{Foreground: "#f3e6d8", Background: "#121212"}},
			},
			{
				Font: status.Font{Foreground: "#9e9e9e", Background: "#303030"},
				GenerateContent: func(ctx status.Context) (time.Duration, error) {
					const dateFormat = "Mon Jan _2"
					fmt.Fprint(ctx.Writer, time.Now().Format(dateFormat))
					return 0, nil
				},
				Name:      "date",
				Separator: status.Separator{Font: status.Font{Foreground: "#303030", Background: "#121212"}, FullArrow: true},
			},
			{
				Font: status.Font{Foreground: "#d0d0d0", Background: "#303030", Bold: true},
				GenerateContent: func(ctx status.Context) (time.Duration, error) {
					const timeFormat = "3:04 PM"
					fmt.Fprint(ctx.Writer, time.Now().Format(timeFormat))
					return 0, nil
				},
				Name:      "time",
				Separator: status.Separator{Font: status.Font{Foreground: "#626262", Background: "#303030"}},
			},
		},
	}
	return statusLine.Status(ctx, goos.Stdout, cacheFS)
}
