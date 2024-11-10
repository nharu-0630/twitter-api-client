package cmd

import "github.com/nharu-0630/twitter-api-client/api"

type BrutePoliticianProps struct {
	Politicians   []Politician `yaml:"politicians"`
	IntervalDays  int          `yaml:"interval_days"`
	SinceDateTime string       `yaml:"since_date_time"`
	UntilDateTime string       `yaml:"until_date_time"`
}

type BrutePoliticianCmd struct {
	CmdName string
	Props   BrutePoliticianProps
	Client  *api.Client
}

func (cmd *BrutePoliticianCmd) Execute() {
	cmd.Client = api.NewClient(
		api.ClientConfig{
			IsGuestTokenEnabled: false,
			AuthToken:           "",
			CsrfToken:           "",
		},
	)
	cmd.CmdName = "BrutePolitician_" + cmd.Props.SinceDateTime + "_" + cmd.Props.UntilDateTime
}
