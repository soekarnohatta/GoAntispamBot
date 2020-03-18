package model

type (
	CasBan struct {
		Ok     bool   `json:"ok"`
		Result Result `json:"result"`
	}

	Result struct {
		Offenses  int      `json:"offenses"`
		Messages  []string `json:"message"`
		TimeAdded int      `json:"time_added"`
	}
)
