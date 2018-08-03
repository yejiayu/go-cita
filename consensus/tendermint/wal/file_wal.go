package wal

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/opentracing/opentracing-go"
)

type FileWAL struct {
	mu sync.RWMutex

	height uint64
	root   string
}

func NewFileWAL(path string) (Interface, error) {
	if err := os.MkdirAll(path, 0700); err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	return &FileWAL{root: path}, nil
}

func (f *FileWAL) SetHeight(ctx context.Context, height uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "wal-height")
	span.SetTag("height", height)
	defer span.Finish()

	f.mu.Lock()
	defer f.mu.Unlock()
	f.height = height

	return nil
}

func (f *FileWAL) Save(ctx context.Context, logType LogType, data []byte) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "wal-save")
	span.SetTag("log_type", logType.String())
	defer span.Finish()

	file, err := os.Create(f.fileName())
	if err != nil {
		if os.IsExist(err) {
			file, err = os.Open(f.fileName())
			if err != nil {
				return err
			}
		}
		return err
	}
	defer file.Close()

	_, err = file.Write(append([]byte{byte(logType)}, data...))
	return err
}

func (f *FileWAL) Load(ctx context.Context) (LogType, []byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "wal-load")
	defer span.Finish()

	file, err := os.Open(f.fileName())
	if err != nil {
		return 0, nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, nil, err
	}

	logType := LogType(data[0])
	msg := data[0:]

	return logType, msg, nil
}

func (f *FileWAL) fileName() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return path.Join(f.root, fmt.Sprintf("%d.log", f.height))
}
