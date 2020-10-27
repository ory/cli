package newsletter

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/gochimp3"
	"github.com/ory/x/flagx"

	"github.com/ory/cli/cmd/pkg"
)

var send = &cobra.Command{
	Use:  "send <list-id>",
	Args: cobra.ExactArgs(1),
	Long: `Send a drafted campaign.

Example:

	$ MAILCHIMP_API_KEY=... \
		CIRCLE_SHA1=... \
		CIRCLE_TAG=... \ # This is set automatically in CircleCI Jobs
		CIRCLE_PROJECT_REPONAME=... \ # This is set automatically in CircleCI Jobs
		release campaign send 12345
`,
	Run: func(cmd *cobra.Command, args []string) {
		SendCampaign(args[0], flagx.MustGetBool(cmd, "dry"))
	},
}

func init() {
	Main.AddCommand(send)

	send.Flags().Bool("dry", false, "Do not actually send the campaign")
}

func SendCampaign(listID string, dry bool) {
	chimpKey := pkg.MustGetEnv("MAILCHIMP_API_KEY")
	chimp := gochimp3.New(chimpKey)
	campaignID := campaignID()

	campaigns, err := chimp.GetCampaigns(&gochimp3.CampaignQueryParams{
		Status:              "save",
		SortField:           "create_time",
		SortDir:             "DESC",
		ListId:              listID,
		ExtendedQueryParams: gochimp3.ExtendedQueryParams{Count: 100},
	})
	pkg.Check(err)

	fmt.Printf(`Looking for campaign "%s"`, campaignID)
	fmt.Println()

	for _, c := range campaigns.Campaigns {
		if c.Settings.Title == campaignID {
			if dry {
				fmt.Println("Skipping send because --dry was passed.")
				return
			}

			chimpCampaignSent, err := chimp.SendCampaign(c.ID, &gochimp3.SendCampaignRequest{
				CampaignId: c.ID,
			})
			pkg.Check(err)

			if !chimpCampaignSent {
				pkg.Fatalf("Unable to send MailChimp Campaign: %s", c.ID)
			}

			fmt.Println("Sent campaign!")
			return
		}
	}

	pkg.Fatalf(`Expected to find campaign "%s" but it could not be found.'`, campaignID)
}
