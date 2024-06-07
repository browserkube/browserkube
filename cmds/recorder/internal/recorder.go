package internal

import (
	"fmt"
	"os"
	"syscall"

	"github.com/urfave/cli/v2"
)

const (
	flagVideoSize  = "video-size"
	flagFrameRate  = "frame-rate"
	flagDisplayNum = "display-num"
	flagCodec      = "codec"
	flagFilePath   = "file-path"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:      "recorder",
		Usage:     "record x11 screen",
		Action:    Record,
		Reader:    os.Stdin,
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagVideoSize,
				Aliases: []string{"s"},
				Value:   "1360x1020",
				Usage:   "video resolution",
			},
			&cli.StringFlag{
				Name:    flagFrameRate,
				Aliases: []string{"r"},
				Value:   "12",
				Usage:   "video frame rate",
			},
			&cli.StringFlag{
				Name:    flagDisplayNum,
				Aliases: []string{"d"},
				Value:   "99",
				Usage:   "x11 display number",
			},
			&cli.StringFlag{
				Name:    flagCodec,
				Aliases: []string{"c"},
				Value:   "libx264",
				Usage:   "video encoder codec",
			},
			&cli.StringFlag{
				Name:    flagFilePath,
				Aliases: []string{"f"},
				Usage:   "internal path to a file",
			},
		},
	}
}

func Record(ctx *cli.Context) error {
	cfg, err := getConfig(ctx)
	if err != nil {
		return fmt.Errorf("unable to parse config: %w", err)
	}

	if err := waitForDisplay(); err != nil {
		return fmt.Errorf("unable to wait for display: %w", err)
	}

	cmd, err := shellAsync(ffmpegCmd, buildArgs(cfg)...)
	if err != nil {
		return fmt.Errorf("unable to run ffmpeg: %w", err)
	}

	// create a child of our command which is ffmpeg
	child, err := sysGetPgid(cmd.Process.Pid)
	if err != nil {
		return fmt.Errorf("unable to get pgid: %w", err)
	}

	<-ctx.Done()

	// kill our ffmpeg process manually and wait for graceful shutdown
	if err = sysKill(-(child), syscall.SIGINT); err != nil {
		return fmt.Errorf("unable kill child process: %w", err)
	}
	var status syscall.WaitStatus
	_, _ = sysWaitFor(-1, &status, 0, nil)

	return nil
}
