package service

import (
	"context"

	"go-fake-flix/internal/common"
	"go-fake-flix/internal/infrastructure/filestore"
)

type MediaService struct {
	videoService VideoService
	fileStorage filestore.FileStore
}

type VideoRange struct {
	Start int64
	End   int64
}

func NewMediaService(videoService VideoService, fileStorage filestore.FileStore) *MediaService {
	return &MediaService{
		videoService: videoService,
		fileStorage:  fileStorage,
	}
}

func (m *MediaService) GetMedia(ctx context.Context, videoID int64, r *VideoRange) ([]byte, *common.BusinessError) {
	video, busErr := m.videoService.GetVideo(ctx, videoID)
	if busErr != nil {
		return nil, busErr
	}

	if r == nil {
		data, err := m.fileStorage.ReadFile(ctx, video.FilePath)
		if err != nil {
			return nil, common.NewBusinessError("VIDEO_READ_FAIL", err.Error())
		}
		return data, nil
	}

	c, err := m.fileStorage.ReadFileRange(ctx, video.FilePath, r.Start, r.End)
	if err != nil {
		return nil, common.NewBusinessError("VIDEO_READ_FAIL", err.Error())
	}
	return c, nil
}
