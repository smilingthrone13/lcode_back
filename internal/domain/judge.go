package domain

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
