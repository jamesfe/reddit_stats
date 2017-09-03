package data_types

type FilterSettings struct {
	SubredditListFile string `json:"list_file"` // a list of subreddits we will use to filter by
}

type InputFilterSettings struct {
	OutputDirectory string `json:"output_dir"` // the output directory to which our matching lines are written
}

type LongevitySettings struct {
	MinDays int `json:"min_days"`
}

type TargetSettings struct {
	RandomSample bool `json:"random_sample"`
	/* true if we should sample randomly, false if we should just search the_donald */
}

type AnalysisSettings struct {
	AnalysisTypes          []string `json:"analysis_types"`
	AnalysisMap            map[string]bool
	LongevityConfiguration LongevitySettings `json:"longevity"`
	TargetConfig           TargetSettings    `json:"target"`
}

type ProfileSettings struct {
	UserListFile string `json:"target_user_list_file"`
	// A file containing all the users we are interested in.
	FilteredDataSource string `json:"filtered_data_source"`
	// A data source for filtered data only, used to find all the usernames we will be studying.
}

type StatsConfiguration struct {
	DataSource               string              `json:"data_source"`
	CheckInterval            int                 `json:"check_interval"` // Print a note to screen every this many lines
	MaxLines                 int                 `json:"max_lines"`      // Analyze up to this many lines
	CpuProfile               string              `json:"cpu_profile"`    // Collect CPU profiling data
	FilterConfiguration      FilterSettings      `json:"filter_settings",omitempty`
	InputFilterConfiguration InputFilterSettings `json:"input_filter_settings",omitempty`
	AnalysisConfiguration    AnalysisSettings    `json:"analysis_settings"`
	ProfileConfiguration     ProfileSettings     `json:"profile_settings"`
}
