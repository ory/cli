package newsletter

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ory/cli/view"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/spf13/cobra"

	"github.com/ory/gochimp3"
	_ "github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"

	"github.com/ory/cli/cmd/pkg"
)

var draft = &cobra.Command{
	Use:   "draft list-id path/to/tag-message path/to/changelog.md",
	Args:  cobra.ExactArgs(3),
	Short: "Creates a draft release notification via the Mailchimp Campaign / Newsletter API",
	Long: `TL;DR

	$ git tag -l --format='%(contents)' v0.0.103 > tag-message.txt
	$ # run changelog generator > changelog.md
	$ MAILCHIMP_API_KEY=... \
		CIRCLE_SHA1=... \
		CIRCLE_TAG=... \ # This is set automatically in CircleCI Jobs
		CIRCLE_PROJECT_REPONAME=... \ # This is set automatically in CircleCI Jobs
		release campaign draft \
			--segment-id ... \ # optional - e.g. only to people interested in ORY Hydra
			list-id-1234123 \
			./tag-message.md \
			./changelog.md

To send out a release newsletter you need to specify an API Key for Mailchimp
(https://admin.mailchimp.com/account/api/) using the MAILCHIMP_API_KEY environment variable:

	export MAILCHIMP_API_KEY=...

Additionally, these CI environment variables are expected to be set as well:

	$CIRCLE_PROJECT_REPONAME (e.g. hydra)
	$CIRCLE_TAG (e.g. v1.4.5-beta.1)
	$CIRCLE_SHA1

If you want to send only to a segment within that list, add the Segment ID as well:

	release notify --segment 1234 ...
`,
	Run: func(cmd *cobra.Command, args []string) {
		listID := args[0]
		tagMessagePath := args[1]
		changelogPath := args[2]

		tagMessageRaw, err := ioutil.ReadFile(tagMessagePath)
		pkg.Check(err)
		changelogRaw, err := ioutil.ReadFile(changelogPath)
		pkg.Check(err)

		chimpCampaign, err := Draft(listID, flagx.MustGetInt(cmd, "segment"), tagMessageRaw, changelogRaw)
		pkg.Check(err)

		fmt.Printf(`Created campaign "%s" (%s)`, chimpCampaign.Settings.Title, chimpCampaign.ID)
		fmt.Println()

		fmt.Println("Campaign drafted")
	},
}

func Draft(listID string, segmentID int, tagMessageRaw, changelogRaw []byte) (*gochimp3.CampaignResponse, error) {
	tag := pkg.GitHubTag()

	var repoName string
	ghRepo := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
	if len(ghRepo) == 2 {
		repoName = ghRepo[1]
	} else {
		repoName = pkg.MustGetEnv("CIRCLE_PROJECT_REPONAME")
	}

	caser := cases.Title(language.AmericanEnglish)
	projectName := "ORY " + caser.String(strings.ToLower(repoName))

	changelog := renderMarkdown(changelogRaw)
	tagMessage := renderMarkdown(tagMessageRaw)

	chimpKey := pkg.MustGetEnv("MAILCHIMP_API_KEY")
	var segmentOptions *gochimp3.CampaignCreationSegmentOptions
	if segmentID > 0 {
		var payload struct {
			Options gochimp3.CampaignCreationSegmentOptions `json:"options"`
		}
		newMailchimpRequest(chimpKey, fmt.Sprintf("/lists/%s/segments/%d", listID, segmentID), &payload)
		segmentOptions = &payload.Options
		segmentOptions.SavedSegmentId = segmentID
	}

	brandColor := "#5528FF"
	switch strings.ToLower(repoName) {
	case "oathkeeper":
		brandColor = "#BD2FEF"
	case "hydra":
		brandColor = "#FF6A85"
	case "kratos":
		brandColor = "#FF9800"
	case "keto":
		brandColor = "#1DE9B6"
	}

	var body bytes.Buffer

	t, err := template.New("mail-body.html").Parse(string(view.MailBody))
	pkg.Check(err)

	pkg.Check(t.Execute(&body, struct {
		Version     string
		GitTag      string
		ProjectName string
		RepoName    string
		Changelog   template.HTML
		Message     template.HTML
		BrandColor  string
	}{
		Version:     tag,
		GitTag:      tag,
		ProjectName: projectName,
		RepoName:    repoName,
		Changelog:   changelog,
		Message:     tagMessage,
		BrandColor:  brandColor,
	}))

	chimp := gochimp3.New(chimpKey)
	chimpTemplate, err := chimp.CreateTemplate(&gochimp3.TemplateCreationRequest{
		Name: substr(fmt.Sprintf("%s %s Announcement", projectName, tag), 0, 49),
		Html: body.String(),
	})
	if err != nil {
		return nil, err
	}

	chimpCampaign, err := chimp.CreateCampaign(&gochimp3.CampaignCreationRequest{
		Type:       gochimp3.CAMPAIGN_TYPE_REGULAR,
		Recipients: gochimp3.CampaignCreationRecipients{ListId: listID, SegmentOptions: segmentOptions},
		Settings: gochimp3.CampaignCreationSettings{
			Title:        campaignID(),
			SubjectLine:  fmt.Sprintf("%s %s has been released!", projectName, tag),
			FromName:     "ORY",
			ReplyTo:      "office@ory.sh",
			Authenticate: true,
			FbComments:   false,
			TemplateId:   chimpTemplate.ID,
		},
		Tracking: gochimp3.CampaignTracking{
			Opens:      true,
			HtmlClicks: true,
			TextClicks: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return chimpCampaign, err
}

func init() {
	Main.AddCommand(draft)
	draft.Flags().Int("segment", 0, "The Mailchimp Segment ID")
}
