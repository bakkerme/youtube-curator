package api

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// YTAPI provides the API globals and implements the ServerInterface
type YTAPI struct{}

// GetJobs returns all Jobs. THIS IS A STUB
func (yt *YTAPI) GetJobs(ctx echo.Context, params GetJobsParams) error {
	resp := &[]Job{
		Job{
			Finished: true,
			Id:       0,
			Running:  false,
			Type:     "youtube-dl",
		},
		Job{
			Finished: false,
			Id:       1,
			Running:  true,
			Type:     "youtube-dl",
		},
	}

	j, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return ctx.String(http.StatusOK, string(j))
}

// GetJobsSocket request
func (yt *YTAPI) GetJobsSocket(ctx echo.Context, jobID string) error {
	return fmt.Errorf("unimplemented")
}

// GetJobsByID request
func (yt *YTAPI) GetJobsByID(ctx echo.Context, jobID string) error {
	return fmt.Errorf("unimplemented")
}

// DeleteVideos deletes the videos
func (yt *YTAPI) DeleteVideos(ctx echo.Context) error {
	return fmt.Errorf("unimplemented")
}

// GetVideos request
func (yt *YTAPI) GetVideos(ctx echo.Context, params GetVideosParams) error {
	return fmt.Errorf("unimplemented")
}

// DownloadVideos starts a download Job for a video
func (yt *YTAPI) DownloadVideos(ctx echo.Context) error {
	return fmt.Errorf("unimplemented")
}

// DeleteVideoByID deletes one video ID
func (yt *YTAPI) DeleteVideoByID(ctx echo.Context, videoID string, params DeleteVideoByIDParams) error {
	return fmt.Errorf("unimplemented")
}

// GetVideoByID returns video data for a video ID
func (yt *YTAPI) GetVideoByID(ctx echo.Context, videoID string, params GetVideoByIDParams) error {
	return fmt.Errorf("unimplemented")
}

// Start sets up the API server
func Start() {
	var ytAPI YTAPI
	e := echo.New()
	RegisterHandlers(e, &ytAPI)

	e.Logger.Fatal(e.Start(":3030"))
}
