package sftputil

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"

	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/logger"
	"github.com/pkg/sftp"
)

type TailFileOptions struct {
	Path           string
	Callback       func(line string)
	IgnoreExisting bool
	Client         *sftp.Client
}

func TailFile(ctx context.Context, options *TailFileOptions) {
	file, err := options.Client.OpenFile(options.Path, os.O_RDONLY)
	if err != nil {
		logger.Error(ctx, err, "failed to open file")
		return
	}
	defer func() {
		logger.Info(ctx, "closing file")
		err := file.Close()
		if err != nil {
			logger.Error(ctx, err, "failed to close file")
		}
	}()

	info, err := file.Stat()
	if err != nil {
		panic(err)
	}
	oldSize := info.Size()
	reader := bufio.NewReader(file)

	reachedEndOnce := false

	for {
		for line, _, err := reader.ReadLine(); err != io.EOF; line, _, err = reader.ReadLine() {
			if !options.IgnoreExisting || (options.IgnoreExisting && reachedEndOnce) {
				options.Callback(string(line))
			}
		}

		reachedEndOnce = true

		pos, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			logger.Error(ctx, err, "failed to seek file")
			return
		}

		for {
			time.Sleep(time.Second)
			newinfo, err := file.Stat()
			if err != nil {
				panic(err)
			}
			newSize := newinfo.Size()
			if newSize != oldSize {
				if newSize < oldSize {
					file.Seek(0, 0)
				} else {
					file.Seek(pos, io.SeekStart)
				}
				reader = bufio.NewReader(file)
				oldSize = newSize
				break
			}

		}
	}

}
