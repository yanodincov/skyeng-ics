package skyeng

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/yanodincov/skyeng-ics/internal/repository/skyeng/model"
	httphelper "github.com/yanodincov/skyeng-ics/pkg/http-helper"
)

const getScheduleURL = "https://student-cabinet-api-self-service-schedule.skyeng.ru/api/common/v1/schedule/get"

type GetScheduleSpec struct {
	Headers map[string]string
	Cookies []*http.Cookie
	UserID  int
}

type GetScheduleData struct {
	Cookies []*http.Cookie
	Lessons []model.Lesson
}

type getScheduleReq struct {
	Buffer     *struct{} `json:"buffer"`
	StudentIDs []int     `json:"studentIds"`
	Page       int       `json:"page"`
}

func (r *Repository) GetSchedule(ctx context.Context, spec GetScheduleSpec) (*GetScheduleData, error) {
	httpReq, err := httphelper.NewRequest(ctx, http.MethodPost, getScheduleURL,
		httphelper.WithHeaders(spec.Headers),
		httphelper.WithJSONBody(getScheduleReq{
			Buffer:     nil,
			Page:       1,
			StudentIDs: []int{spec.UserID},
		}),
		httphelper.WithCookies(spec.Cookies),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	httpRes, err := r.client.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", httpRes.StatusCode)
	}

	decodedBody, err := httphelper.DecodeHTTPResponseBody(httpRes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	var schedule model.Schedule
	if err = json.Unmarshal(decodedBody, &schedule); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	return &GetScheduleData{
		Cookies: httphelper.MergeCookies(httpRes, spec.Cookies),
		Lessons: schedule.Data,
	}, nil
}
