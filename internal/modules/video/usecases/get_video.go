package usecases

import (
	"context"
	"errors"

	"go-fake-flix/internal/common"
	"go-fake-flix/internal/modules/video/entities"
	"go-fake-flix/internal/modules/video/repository"
)

type VideoRange struct {
	start int
	end   int
}

func GetVideo(ctx context.Context, repo repository.VideoRepository, id string) (*entities.Video, *common.BusinessError) {
	v, err := repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, common.NewBusinessError("VIDEO_NOT_FOUND", "Video not found  ")
		}
		return nil, common.NewBusinessError("VIDEO_FETCH_FAILED", "Could not fetch video")
	}
	return &v, nil
}
