package usecases

import (
	"context"
	"errors"

	"go-fake-flix/internal/common"
	"go-fake-flix/internal/modules/video/entities"
	"go-fake-flix/internal/modules/video/repository"
)

func UploadVideo(ctx context.Context, repo repository.VideoRepository, fileName, filePath string) *common.BusinessError {
	if err := repo.Add(ctx, entities.Video{FileName: fileName, FilePath: filePath}); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return common.NewBusinessError("VIDEO_ALREADY_EXISTS", "Video already exists")
		}
		return common.NewBusinessError("VIDEO_NOT_STORED", "Could not store video")
	}
	return nil
}
