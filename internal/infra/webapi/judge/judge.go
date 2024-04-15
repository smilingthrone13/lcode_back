package judge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/pkg/struct_errors"
	"net/http"
)

const (
	waitQuery = "wait"

	submissionFields = "token,stdout,stderr,time,memory,message,status"
)

func New(cfg *config.JudgeConfig) *API {
	return &API{
		addr:   fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port),
		client: &http.Client{},
	}
}

type API struct {
	addr   string
	client *http.Client
}

func (a *API) CreateSubmission(
	ctx context.Context,
	data domain.CreateJudgeSubmission,
) (domain.JudgeSubmissionInfo, error) {
	var submissionResp createSubmissionResponse

	reqData := createSubmissionRequest{
		CreateJudgeSubmission: data,
		Fields:                submissionFields,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return domain.JudgeSubmissionInfo{}, errors.Wrap(err, "CreateSubmission judge api")
	}

	req, err := http.NewRequest("POST", a.addr+"/submissions", bytes.NewBuffer(jsonData))
	if err != nil {
		return domain.JudgeSubmissionInfo{}, errors.Wrap(err, "CreateSubmission judge api")
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add(waitQuery, "true")
	req.URL.RawQuery = q.Encode()

	resp, err := a.client.Do(req.WithContext(ctx))
	if err != nil {
		err = struct_errors.NewBaseErr("Code solving system is unavailable", err)

		return domain.JudgeSubmissionInfo{}, errors.Wrap(err, "CreateSubmission judge api")
	}

	d := json.NewDecoder(resp.Body)

	switch resp.StatusCode {
	case http.StatusCreated:
		err = d.Decode(&submissionResp)
	case http.StatusServiceUnavailable:
		err = domain.NewJudgeQueueIsFullError()
	default:
		err = struct_errors.NewInternalErr(
			fmt.Errorf("bad request to judge api with status code: %d", resp.StatusCode),
		)
	}

	if err != nil {
		return domain.JudgeSubmissionInfo{}, errors.Wrap(err, "CreateSubmission judge api")
	}

	info := domain.JudgeSubmissionInfo{
		Token:  submissionResp.Token,
		Stdout: submissionResp.Stdout,
		Stderr: submissionResp.Stderr,
		Time:   submissionResp.Time,
		Memory: submissionResp.Memory,
		Status: submissionResp.Status.ID,
	}

	return info, nil
}

func (a *API) GetAvailableLanguages(ctx context.Context) ([]domain.JudgeLanguageInfo, error) {
	var languages []domain.JudgeLanguageInfo

	req, err := http.NewRequest("GET", a.addr+"/languages", nil)
	if err != nil {
		return languages, errors.Wrap(err, "GetAvailableLanguages judge api")
	}

	resp, err := a.client.Do(req.WithContext(ctx))
	if err != nil {
		err = struct_errors.NewBaseErr("Code solving system is unavailable", err)

		return languages, errors.Wrap(err, "GetAvailableLanguages judge api")
	}

	d := json.NewDecoder(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		err = d.Decode(&languages)
	default:
		err = struct_errors.NewInternalErr(
			fmt.Errorf("bad request to judge api with status code: %d", resp.StatusCode),
		)
	}

	if err != nil {
		return languages, errors.Wrap(err, "GetAvailableLanguages judge api")
	}

	return languages, nil
}

func (a *API) GetAvailableStatuses(ctx context.Context) ([]domain.JudgeStatusInfo, error) {
	var statuses []domain.JudgeStatusInfo

	req, err := http.NewRequest("GET", a.addr+"/statuses", nil)
	if err != nil {
		return statuses, errors.Wrap(err, "GetAvailableStatuses judge api")
	}

	resp, err := a.client.Do(req.WithContext(ctx))
	if err != nil {
		err = struct_errors.NewBaseErr("Code solving system is unavailable", err)

		return statuses, errors.Wrap(err, "GetAvailableStatuses judge api")
	}

	d := json.NewDecoder(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		err = d.Decode(&statuses)
	default:
		err = struct_errors.NewInternalErr(
			fmt.Errorf("bad request to judge api with status code: %d", resp.StatusCode),
		)
	}

	if err != nil {
		return statuses, errors.Wrap(err, "GetAvailableStatuses judge api")
	}

	return statuses, nil
}
