package judge

import "lcode/internal/domain"

type status struct {
	ID          domain.JudgeStatus `json:"id"`
	Description string             `json:"description"`
}

type createSubmissionRequest struct {
	domain.CreateJudgeSubmission
	Fields string `json:"fields"`
}

type createSubmissionResponse struct {
	Stdout *string `json:"stdout"`
	Stderr *string `json:"stderr"`
	Memory int     `json:"memory"`
	Time   float64 `json:"time,string"`
	Token  string  `json:"token"`
	Status status  `json:"status"`
}
