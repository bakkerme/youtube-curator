// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

// Job defines model for Job.
type Job struct {
	Finished bool    `json:"finished"`
	Id       float32 `json:"id"`
	Running  bool    `json:"running"`
	Type     string  `json:"type"`
}

// Video defines model for Video.
type Video struct {
	Creator     string `json:"creator"`
	Description string `json:"description"`
	Duration    string `json:"duration"`
	FileType    string `json:"fileType"`
	Id          string `json:"id"`
	Path        string `json:"path"`
	PublishedAt string `json:"publishedAt"`
	Title       string `json:"title"`
}

// Deleted defines model for deleted.
type Deleted struct {
	Id *string `json:"id,omitempty"`
}

// Error defines model for error.
type Error struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// GetJobsParams defines parameters for GetJobs.
type GetJobsParams struct {

	// Filter by job status
	Status *string `json:"status,omitempty"`
}

// DeleteVideosJSONBody defines parameters for DeleteVideos.
type DeleteVideosJSONBody struct {
	VideoID string `json:"videoID"`
}

// GetVideosParams defines parameters for GetVideos.
type GetVideosParams struct {
	ChannelID *string `json:"channelID,omitempty"`
}

// DownloadVideosJSONBody defines parameters for DownloadVideos.
type DownloadVideosJSONBody struct {
	ChannelID  *[]string `json:"channelID,omitempty"`
	PlaylistID *[]string `json:"playlistID,omitempty"`
	VideoID    *[]string `json:"videoID,omitempty"`
}

// DeleteVideoByIDParams defines parameters for DeleteVideoByID.
type DeleteVideoByIDParams struct {

	// A video ID
	VideoID *string `json:"videoID,omitempty"`
}

// GetVideoByIDParams defines parameters for GetVideoByID.
type GetVideoByIDParams struct {

	// A video ID to get data from
	VideoID string `json:"videoID"`
}

// DeleteVideosRequestBody defines body for DeleteVideos for application/json ContentType.
type DeleteVideosJSONRequestBody DeleteVideosJSONBody

// DownloadVideosRequestBody defines body for DownloadVideos for application/json ContentType.
type DownloadVideosJSONRequestBody DownloadVideosJSONBody
