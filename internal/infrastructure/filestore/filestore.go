package filestore

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileStore interface {
	Save(ctx context.Context, id string, reader io.Reader) error
	ReadFile(ctx context.Context, filePath string) ([]byte, error)
	ReadFileRange(ctx context.Context, filePath string, start int64, end int64) ([]byte, error)
	GetPath(id string) string
	Delete(id string) error
}

type LocalStore struct {
	basePath string
}

func NewLocalStore(basePath string) *LocalStore {
	return &LocalStore{basePath: basePath}
}


func (s *LocalStore) Save(ctx context.Context, id string, reader io.Reader) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	
	path := s.GetPath(id)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}
	
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	
	if _, err := io.Copy(f, reader); err != nil {
		_ = f.Close()
		return fmt.Errorf("write file: %w", err)
	}
	
	if err := f.Close(); err != nil {
		return fmt.Errorf("close file: %w", err)
	}
	return nil
}

func (s *LocalStore) ReadFile(ctx context.Context, filePath string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return os.ReadFile(filePath)
}

func (s *LocalStore) ReadFileRange(ctx context.Context, filePath string, start int64, end int64) ([]byte, error){
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file")
	}

	defer f.Close();

	buf := make([]byte, end-start+1)
	_, err = f.ReadAt(buf, start)
	if err != nil {
		return nil, fmt.Errorf("failed to read the file")
	}
	return buf, nil
}

func (s *LocalStore) GetPath(id string) string {
	return filepath.Join(s.basePath, id)
}

func (s *LocalStore) Delete(id string) error {
	return os.Remove(s.GetPath(id))
}
