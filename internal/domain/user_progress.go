package domain

type (
	ProgressType   string
	StatisticsType string
)

const (
	ProgressCompleted  ProgressType = "completed"
	ProgressInProgress ProgressType = "in_progress"

	StatisticCategory   StatisticsType = "category"
	StatisticDifficulty StatisticsType = "difficulty"
)

type (
	ProgressData struct {
		Status  string   `json:"status" db:"status"`
		TaskIDs []string `json:"task_ids" db:"task_ids"`
	}

	UserProgress struct {
		Progress []ProgressData `json:"progress_data"`
	}

	StatisticData struct {
		Param      string `json:"param" db:"param"`
		CountDone  int    `json:"count_done" db:"count_done"`
		CountTotal int    `json:"count_total" db:"count_total"`
	}

	UserStatistic struct {
		Type       StatisticsType  `json:"statistic_type"`
		Statistics []StatisticData `json:"statistic_data"`
	}
)

type (
	GetStatisticsDTO struct {
		UserID string         `json:"-"`
		Type   StatisticsType `json:"type"`
	}
)
