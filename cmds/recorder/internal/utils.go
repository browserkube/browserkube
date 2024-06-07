package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/urfave/cli/v2"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
)

const (
	getDisplayCmd = `(
cd /tmp/.X11-unix && 
for x in X*; 
do echo ":${x#X}"; 
done
)
`
	ffmpegCmd = "ffmpeg"

	noDisplayFound = "exit status 2"
)

type Config struct {
	VideoSize         string
	FrameRate         string
	DisplayNum        string
	Codec             string
	FileName          string
	SaveVideoEndpoint string
	SessionID         string
	FilePath          string
}

func getConfig(ctx *cli.Context) (*Config, error) {
	return &Config{
		VideoSize:  ctx.String(flagVideoSize),
		FrameRate:  ctx.String(flagFrameRate),
		DisplayNum: ctx.String(flagDisplayNum),
		Codec:      ctx.String(flagCodec),
		FileName:   "video.mp4",
		FilePath:   ctx.String(flagFilePath),
	}, nil
}

func buildArgs(cfg *Config) []string {
	return []string{
		"-y", "-nostdin", "-f", "x11grab", "-draw_mouse", "0",
		"-video_size", cfg.VideoSize,
		"-r", cfg.FrameRate,
		"-i", ":" + cfg.DisplayNum + ".0",
		"-codec:v", cfg.Codec,
		"-pix_fmt", "yuv420p",
		"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2",
		path.Join(cfg.FilePath, cfg.FileName),
	}
}

func shellSync(command string, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func shellAsync(command string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(command, args...)
	cmd.SysProcAttr = sysProcAttrSetPgid
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	go func() {
		_, _ = io.Copy(os.Stdout, stdout)
	}()
	go func() {
		_, _ = io.Copy(os.Stderr, stderr)
	}()
	return cmd, err
}

func waitForDisplay() error {
displayWaitLoop:
	for timeout := time.After(time.Second * 60); ; {
		select {
		case <-timeout:
			return errors.New("timeout waiting for display")
		default:
			_, _, err := shellSync("sh", "-c", getDisplayCmd)
			if err != nil {
				if err.Error() == noDisplayFound {
					continue
				}
				return fmt.Errorf("unable to list displays: %w", err)
			}
			break displayWaitLoop
		}
	}

	return nil
}
