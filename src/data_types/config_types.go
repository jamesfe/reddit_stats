package data_types

type FilterSettings struct {
	SubredditListFile string `json:"list_file"`
}

type InputFilterSettings struct {
	OutputDirectory string `json:"output_dir"`
}

type AnalysisSettings struct {
	AnalysisTypes []string `json:"analysis_types"`
	AnalysisMap   map[string]bool
}

type StatsConfiguration struct {
	DataSource               string              `json:"data_source"`
	CheckInterval            int                 `json:"check_interval"`
	MaxLines                 int                 `json:"max_lines"`
	CpuProfile               string              `json:"cpu_profile"`
	FilterConfiguration      FilterSettings      `json:"filter_settings",omitempty`
	InputFilterConfiguration InputFilterSettings `json:"input_filter_settings",omitempty`
	AnalysisConfiguration    AnalysisSettings    `json:"analysis_settings"`
}
