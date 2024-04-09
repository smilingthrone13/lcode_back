package judge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"lcode/internal/domain"
	"lcode/pkg/struct_errors"
	"net/http"
)

const (
	waitQuery = "wait"

	submissionFields = "token,stdout,stderr,time,memory,message,status"
)

func New(host string, port int) *API {
	return &API{
		addr:   fmt.Sprintf("%s:%d", host, port),
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
