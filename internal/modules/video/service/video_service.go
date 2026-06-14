package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"go-fake-flix/internal/common"
	"go-fake-flix/internal/infrastructure/filestore"
	"go-fake-flix/internal/modules/video/entities"
	"go-fake-flix/internal/modules/video/repository"
)

type VideoService interface {
	UploadVideo(ctx context.Context, fileName string, reader io.Reader) (int64, *common.BusinessError)
	GetVideo(ctx context.Context, id int64) (*entities.Video, *common.BusinessError)
}

type Service struct {
	repo      repository.VideoRepository
	filestore filestore.FileStore
}

func New(repo repository.VideoRepository, fs filestore.FileStore) *Service {
	return &Service{repo: repo, filestore: fs}
}

func (s *Service) UploadVideo(ctx context.Context, fileName string, reader io.Reader) (int64, *common.BusinessError) {
	id := s.repo.NextID()

	fileID := fmt.Sprintf("%d.mp4", id)
	if err := s.filestore.Save(ctx, fileID, reader); err != nil {
		return 0, common.NewBusinessError("VIDEO_FILE_SAVE_FAILED", "Could not save video file")
	}

	video := entities.Video{
		ID:       id,
		FileName: fileName,
		FilePath: s.filestore.GetPath(fileID),
	}
	if err := s.repo.Add(ctx, video); err != nil {
		return 0, common.NewBusinessError("VIDEO_NOT_STORED", "Could not store video")
	}
	return video.ID, nil
}

func (s *Service) GetVideo(ctx context.Context, id int64) (*entities.Video, *common.BusinessError) {
	v, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, common.NewBusinessError("VIDEO_NOT_FOUND", "Video not found")
		}
		return nil, common.NewBusinessError("VIDEO_FETCH_FAILED", "Could not fetch video")
	}
	return &v, nil
}
