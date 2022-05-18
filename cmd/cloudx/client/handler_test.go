package client_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"
	cloud "github.com/ory/client-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestCommandHelper(t *testing.T) {
	configDir := testhelpers.NewConfigDir(t)
	email, password := testhelpers.RegisterAccount(t, configDir)
	project := testhelpers.CreateProject(t, configDir)

	stdIn := &bytes.Buffer{}
	loggedIn := &client.CommandHelper{
		ConfigLocation:   configDir,
		NoConfirm:        true,
		IsQuiet:          true,
		VerboseWriter:    io.Discard,
		VerboseErrWriter: io.Discard,
		Stdin:            bufio.NewReader(stdIn),
		Ctx:              context.Background(),
	}

	t.Run("func=GetProject", func(t *testing.T) {
		assertValidProject := func(t *testing.T, actual *cloud.Project) {
			assert.Equal(t, project, actual.Id)
			assert.NotZero(t, actual.Slug)
			assert.NotZero(t, actual.Services.Identity.Config)
			assert.NotZero(t, actual.Services.Permission.Config)
		}

		t.Run(fmt.Sprintf("is able to get project"), func(t *testing.T) {
			p, err := loggedIn.GetProject(project)
			require.NoError(t, err)
			assertValidProject(t, p)
		})

		t.Run("is not able to list projects if not authenticated and quiet flag", func(t *testing.T) {
			notLoggedIn := *loggedIn
			notLoggedIn.ConfigLocation = testhelpers.NewConfigDir(t)
			_, err := notLoggedIn.GetProject(project)
			assert.ErrorIs(t, err, client.ErrNoConfigQuiet)
		})

		t.Run("is able to get project after authenticating", func(t *testing.T) {
			notYetLoggedIn := *loggedIn
			notYetLoggedIn.ConfigLocation = testhelpers.NewConfigDir(t)
			notYetLoggedIn.IsQuiet = false
			notYetLoggedIn.PwReader = func() ([]byte, error) {
				return []byte(password), nil
			}

			stdIn.WriteString("y\n")        // Do you already have an Ory Console account you wish to use? [y/n]: y
			stdIn.WriteString(email + "\n") // Email fakeEmail()
			p, err := notYetLoggedIn.GetProject(project)
			require.NoError(t, err)
			assertValidProject(t, p)
		})
	})

}
