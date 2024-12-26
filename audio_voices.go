package coze

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (r *audioVoices) Clone(ctx context.Context, req *CloneAudioVoicesReq) (*CloneAudioVoicesResp, error) {
	path := "/v1/audio/voices/clone"
	if req.File == nil {
		return nil, fmt.Errorf("file is required")
	}

	fields := map[string]string{
		"voice_name":   req.VoiceName,
		"audio_format": req.AudioFormat.String(),
	}

	// Add other fields
	if req.Language != nil {
		fields["language"] = req.Language.String()
	}
	if req.VoiceID != nil {
		fields["voice_id"] = *req.VoiceID
	}
	if req.PreviewText != nil {
		fields["preview_text"] = *req.PreviewText
	}
	if req.Text != nil {
		fields["text"] = *req.Text
	}
	resp := &cloneAudioVoicesResp{}
	if err := r.core.UploadFile(ctx, path, req.File, req.VoiceName, fields, resp); err != nil {
		return nil, err
	}
	resp.Data.setHTTPResponse(resp.HTTPResponse)
	return resp.Data, nil
}

func (r *audioVoices) List(ctx context.Context, req *ListAudioVoicesReq) (*NumberPaged[Voice], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged[Voice](
		func(request *PageRequest) (*PageResponse[Voice], error) {
			uri := "/v1/audio/voices"
			resp := &ListAudioVoicesResp{}
			err := r.core.Request(ctx, http.MethodGet, uri, nil, resp,
				withHTTPQuery("page_num", strconv.Itoa(request.PageNum)),
				withHTTPQuery("page_size", strconv.Itoa(request.PageSize)),
				withHTTPQuery("filter_system_voice", strconv.FormatBool(req.FilterSystemVoice)))
			if err != nil {
				return nil, err
			}
			return &PageResponse[Voice]{
				HasMore: len(resp.Data.VoiceList) >= request.PageSize,
				Data:    resp.Data.VoiceList,
				LogID:   resp.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

type audioVoices struct {
	core *core
}

func newVoice(core *core) *audioVoices {
	return &audioVoices{core: core}
}

// Voice represents the voice model
type Voice struct {
	VoiceID                string `json:"voice_id"`
	Name                   string `json:"name"`
	IsSystemVoice          bool   `json:"is_system_voice"`
	LanguageCode           string `json:"language_code"`
	LanguageName           string `json:"language_name"`
	PreviewText            string `json:"preview_text"`
	PreviewAudio           string `json:"preview_audio"`
	AvailableTrainingTimes int    `json:"available_training_times"`
	CreateTime             int    `json:"create_time"`
	UpdateTime             int    `json:"update_time"`
}

// CloneAudioVoicesReq represents the request for cloning a voice
type CloneAudioVoicesReq struct {
	VoiceName   string
	File        io.Reader
	AudioFormat AudioFormat
	Language    *LanguageCode
	VoiceID     *string
	PreviewText *string
	Text        *string
}

// cloneAudioVoicesResp represents the response for cloning a voice
type cloneAudioVoicesResp struct {
	baseResponse
	Data *CloneAudioVoicesResp `json:"data"`
}

// CloneAudioVoicesResp represents the response for cloning a voice
type CloneAudioVoicesResp struct {
	baseModel
	VoiceID string `json:"voice_id"`
}

// ListAudioVoicesReq represents the request for listing voices
type ListAudioVoicesReq struct {
	FilterSystemVoice bool `json:"filter_system_voice,omitempty"`
	PageNum           int  `json:"page_num"`
	PageSize          int  `json:"page_size"`
}

// ListAudioVoicesResp represents the response for listing voices
type ListAudioVoicesResp struct {
	baseResponse
	Data struct {
		VoiceList []*Voice `json:"voice_list"`
	} `json:"data"`
}
