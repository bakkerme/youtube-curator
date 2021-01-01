package api

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
	// "hyperfocus.systems/youtube-curator-server/videometadata"
	"net/http"
)

// YTAPI provides the API globals and implements the ServerInterface
type YTAPI struct {
	cfg *config.Config
}

// GetChannels returns all available Channels
func (yt *YTAPI) GetChannels(ctx echo.Context) error {
	ytChannels, err := getChannels(yt.cfg, &collection.YTChannelLoad{})
	if err != nil {
		return fmt.Errorf("Could not get channels. Error %s ", err)
	}

	resp, err := json.Marshal(ytChannels)
	if err != nil {
		return fmt.Errorf("Could not get channels. Error %s ", err)
	}

	return ctx.String(http.StatusOK, string(resp))
}

func getChannels(cfg *config.Config, ytcl collection.YTChannelLoader) (*[]collection.YTChannelData, error) {
	ytChannels, err := ytcl.GetAvailableYTChannels(cfg)
	if err != nil {
		return nil, err
	}

	ytChannelResponse := []collection.YTChannelData{}

	for _, ytChannel := range *ytChannels {
		ytChannelResponse = append(ytChannelResponse, collection.YTChannelData{
			IName:         ytChannel.Name(),
			IRSSURL:       ytChannel.RSSURL(),
			IChannelURL:   ytChannel.ChannelURL(),
			IArchivalMode: ytChannel.ArchivalMode(),
		})
	}

	return &ytChannelResponse, nil
}

// GetChannelByID returns a channel with the provided ID
func (yt *YTAPI) GetChannelByID(ctx echo.Context, channelID string) error {
	ytChannel, err := getChannelByID(channelID, yt.cfg, &collection.YTChannelLoad{})
	if err != nil {
		return fmt.Errorf("Could not get channels. Error %s ", err)
	}

	if ytChannel == nil {
		return ctx.String(404, "")
	}

	resp, err := json.Marshal(ytChannel)
	if err != nil {
		return fmt.Errorf("Could not get channels. Error %s ", err)
	}

	return ctx.String(http.StatusOK, string(resp))

}

func getChannelByID(id string, cfg *config.Config, ytcl collection.YTChannelLoader) (*collection.YTChannel, error) {
	ytChannels, err := ytcl.GetAvailableYTChannels(cfg)
	if err != nil {
		return nil, err
	}

	ytChannel := (*ytChannels)[id]

	if ytChannel == nil {
		return nil, nil
	}

	return &ytChannel, nil
}

// CheckChannelUpdates checks the Youtube API for updates to a Channel's Videos
func (yt *YTAPI) CheckChannelUpdates(ctx echo.Context, channelID string) error {
	return fmt.Errorf("unimplemented")
}

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
	cfg, err := config.GetConfig(&config.FileConfigProvider{})
	if err != nil {
		return err
	}

	_, err = collection.GetVideoByID(params.VideoID, cfg)
	if err != nil {
		return err
	}

	return fmt.Errorf("unimplemented")
}

// Start sets up the API server
func Start() {
	cfg, err := config.GetConfig(&config.FileConfigProvider{})
	if err != nil {
		panic(err)
	}

	ytAPI := YTAPI{
		cfg: cfg,
	}

	e := echo.New()
	RegisterHandlers(e, &ytAPI)

	e.Logger.Fatal(e.Start(":3030"))
}
