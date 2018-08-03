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
	if err := os.Mkdir(path.Join(f.root, fmt.Sprintf("%d", height)), 0700); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	return nil
}

func (f *FileWAL) Save(ctx context.Context, logType LogType, data []byte) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "wal-save")
	span.SetTag("log_type", logType.String())
	defer span.Finish()

	file, err := os.Create(f.fileName(logType))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}

func (f *FileWAL) Load(ctx context.Context, logType LogType) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "wal-load")
	span.SetTag("log_type", logType.String())
	defer span.Finish()

	file, err := os.Open(f.fileName(logType))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func (f *FileWAL) fileName(logType LogType) string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	name := fmt.Sprintf("%d/%d.log", f.height, logType)

	return path.Join(f.root, name)
}
