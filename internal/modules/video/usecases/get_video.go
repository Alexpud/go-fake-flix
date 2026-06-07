package usecases

import (
	"context"
	"errors"
	"net/http"

	"go-fake-flix/internal/common"
	"go-fake-flix/internal/modules/video/entities"
	"go-fake-flix/internal/modules/video/repository"
)

type VideoRange struct {
	start int
	end   int
}

func GetVideo(ctx context.Context, repo repository.VideoRepository, id string) (*entities.Video, *common.AppError) {
	v, err := repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, common.NewAppError("VIDEO_NOT_FOUND", "Video not found", http.StatusNotFound)
		}
		return nil, common.NewAppError("VIDEO_FETCH_FAILED", "Could not fetch video", http.StatusInternalServerError)
	}
	return &v, nil
}
