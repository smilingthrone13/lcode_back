package domain

import "lcode/pkg/struct_errors"

type JudgeStatus int

const (
	InQueue JudgeStatus = iota + 1
	Processing
	Accepted
	WrongAnswer
	TimeLimitExceeded
	CompilationError
	RuntimeSIGSEV
	RuntimeSIGXFSZ
	RuntimeSIGFPE
	RuntimeSIGABRT
	RuntimeNZEC
	RuntimeOther
	InternalError
	ExecFormatError
)

type LanguageType int

const (
	TypeScript_3_7_4 LanguageType = 74
	NodeJS_12_14_0   LanguageType = 63
)

type CreateJudgeSubmission struct {
	SourceCode     string       `json:"source_code"`
	LanguageID     LanguageType `json:"language_id"`
	Stdin          string       `json:"stdin"`
	ExpectedOutput string       `json:"expected_output"`
	CpuTimeLimit   float64      `json:"cpu_time_limit"`
	MemoryLimit    int          `json:"memory_limit"`
}

type JudgeSubmissionInfo struct {
	Token  string      `json:"token"`
	Stdout *string     `json:"stdout"`
	Stderr *string     `json:"stderr"`
	Time   float64     `json:"time"`
	Memory int         `json:"memory"`
	Status JudgeStatus `json:"status"`
}

// errors
type JudgeQueueIsFullError struct {
	struct_errors.BaseError
}

func NewJudgeQueueIsFullError() *JudgeQueueIsFullError {
	e := &JudgeQueueIsFullError{}
	e.SetCode("judge.queue_is_full")
	//e.SetMsg("Unknown error")
	e.SetErr("Judge queue is full", nil)

	return e
}
