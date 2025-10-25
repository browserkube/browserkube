package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
)

const (
	videoSize  = "1920x1080"
	frameRate  = "30"
	displayNum = "2"
	codec      = "test-codec"
	out        = "video.mp4"
	filePath   = "/home/videos"
)

func prepareTestCtx(t *testing.T) *cli.Command {
	t.Helper()

	cmd := &cli.Command{
		Name: "test",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: flagVideoSize, Value: videoSize},
			&cli.StringFlag{Name: flagFrameRate, Value: frameRate},
			&cli.StringFlag{Name: flagDisplayNum, Value: displayNum},
			&cli.StringFlag{Name: flagCodec, Value: codec},
			&cli.StringFlag{Name: flagFilePath, Value: filePath},
		},
	}

	// Setting values to simulate flag usage
	for _, flag := range cmd.Flags {
		strFlag, ok := flag.(*cli.StringFlag)
		if ok {
			cmd.Set(strFlag.Name, strFlag.Value)
		}
	}

	return cmd
}

func Test_SetConfig(t *testing.T) {
	cmd := prepareTestCtx(t)

	cfg := getConfig(cmd)
	args := buildArgs(cfg)

	expected := []string{
		"-y", "-nostdin", "-f", "x11grab", "-draw_mouse", "0", "-video_size", videoSize,
		"-r", frameRate, "-i", fmt.Sprintf(":%s.0", displayNum), "-codec:v", codec, "-pix_fmt", "yuv420p", "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2", "/home/videos/video.mp4",
	}
	assert.Equal(t, expected, args)
}

func Test_shellAsync(t *testing.T) {
	cmd, err := shellAsync("go", "version")
	assert.NoError(t, err)
	assert.Empty(t, cmd.Err)
}

func Test_shellSync(t *testing.T) {
	stOut, stErr, err := shellSync("go", "version")
	assert.NoError(t, err)
	assert.Empty(t, stErr)
	assert.Contains(t, stOut, "go version")
}
