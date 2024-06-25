// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/tidwall/gjson"

	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/stringsx"
)

func (h *CommandHelper) checkSession(ctx context.Context) error {
	c, err := h.getConfig()
	if err != nil {
		return err
	}
	if c.isAuthenticated {
		return nil
	}
	if c.SessionToken == "" {
		return ErrNotAuthenticated
	}

	client, err := NewOryProjectClient()
	if err != nil {
		return err
	}

	sess, _, err := client.FrontendAPI.ToSession(ctx).XSessionToken(c.SessionToken).Execute()
	if err != nil || sess == nil {
		return stderrors.Join(err, ErrReauthenticate)
	}
	c.isAuthenticated = true

	return nil
}

func (h *CommandHelper) GetAuthenticatedConfig(ctx context.Context) (*Config, error) {
	if err := h.checkSession(ctx); err == nil {
		return h.getConfig()
	} else if stderrors.Is(err, ErrReauthenticate) {
		if h.isQuiet {
			return nil, ErrNoConfigQuiet
		}
		_, _ = fmt.Fprintf(h.VerboseErrWriter, "Your session has expired or has otherwise become invalid. Please re-authenticate to continue.\n")
	} else if stderrors.Is(err, ErrNoConfig) || stderrors.Is(err, ErrNotAuthenticated) {
		if h.isQuiet {
			return nil, ErrNoConfigQuiet
		}
	}
	if err := h.ClearConfig(); err != nil {
		return nil, err
	}

	if err := h.Authenticate(ctx); err != nil {
		return nil, err
	}

	return h.getConfig()
}

func (h *CommandHelper) signup(ctx context.Context, c *cloud.APIClient) error {
	flow, _, err := c.FrontendAPI.CreateNativeRegistrationFlow(ctx).Execute()
	if err != nil {
		return err
	}

doRegistration:
	var form cloud.UpdateRegistrationFlowWithPasswordMethod
	if err := renderForm(h.Stdin, h.pwReader, h.VerboseErrWriter, flow.Ui, "password", &form); err != nil {
		return err
	}

	signup, _, err := c.FrontendAPI.
		UpdateRegistrationFlow(ctx).
		Flow(flow.Id).
		UpdateRegistrationFlowBody(cloud.UpdateRegistrationFlowBody{UpdateRegistrationFlowWithPasswordMethod: &form}).
		Execute()
	if err != nil {
		if e, ok := err.(*cloud.GenericOpenAPIError); ok {
			switch m := e.Model().(type) {
			case *cloud.RegistrationFlow:
				flow = m
			case cloud.RegistrationFlow:
				flow = &m
			}
			_, _ = fmt.Fprintf(h.VerboseErrWriter, "\nYour account creation attempt failed. Please try again!\n\n")
			goto doRegistration
		}

		return err
	}

	sessionToken := *signup.SessionToken
	sess, _, err := c.FrontendAPI.ToSession(ctx).XSessionToken(sessionToken).Execute()
	if err != nil {
		return err
	}

	config := new(Config)
	if err := config.fromSession(sess, sessionToken); err != nil {
		return err
	}
	return h.UpdateConfig(config)
}

func (h *CommandHelper) signin(ctx context.Context, c *cloud.APIClient, sessionToken string) error {
	req := c.FrontendAPI.CreateNativeLoginFlow(ctx)
	if len(sessionToken) > 0 {
		req = req.XSessionToken(sessionToken).Aal("aal2")
	}
	flow, _, err := req.Execute()
	if err != nil {
		return err
	}

doLogin:
	var form interface{} = &cloud.UpdateLoginFlowWithPasswordMethod{}
	method := "password"
	if len(sessionToken) > 0 {
		var foundTOTP, foundLookup bool
		for _, n := range flow.Ui.Nodes {
			foundTOTP = foundTOTP || n.Group == "totp"
			foundLookup = foundLookup || n.Group == "lookup_secret"
		}
		if !foundLookup && !foundTOTP {
			return stderrors.New("only TOTP and lookup secrets are supported for two-step verification in the CLI")
		}

		method = "lookup_secret"
		if foundTOTP {
			form = &cloud.UpdateLoginFlowWithTotpMethod{}
			method = "totp"
		}
	}

	if err := renderForm(h.Stdin, h.pwReader, h.VerboseErrWriter, flow.Ui, method, form); err != nil {
		return err
	}

	var body cloud.UpdateLoginFlowBody
	switch e := form.(type) {
	case *cloud.UpdateLoginFlowWithTotpMethod:
		body.UpdateLoginFlowWithTotpMethod = e
	case *cloud.UpdateLoginFlowWithPasswordMethod:
		body.UpdateLoginFlowWithPasswordMethod = e
	default:
		panic("unexpected type")
	}

	login, _, err := c.FrontendAPI.UpdateLoginFlow(ctx).XSessionToken(sessionToken).
		Flow(flow.Id).UpdateLoginFlowBody(body).Execute()
	if err != nil {
		if e, ok := err.(*cloud.GenericOpenAPIError); ok {
			switch m := e.Model().(type) {
			case *cloud.LoginFlow:
				flow = m
			case cloud.LoginFlow:
				flow = &m
			}
			_, _ = fmt.Fprintf(h.VerboseErrWriter, "\nYour sign in attempt failed. Please try again!\n\n")
			goto doLogin
		}

		return err
	}

	sessionToken = stringsx.Coalesce(*login.SessionToken, sessionToken)
	sess, _, err := c.FrontendAPI.ToSession(ctx).XSessionToken(sessionToken).Execute()
	if err != nil {
		if e, ok := err.(interface{ Body() []byte }); ok {
			switch gjson.GetBytes(e.Body(), "error.id").String() {
			case "session_aal2_required":
				return h.signin(ctx, c, sessionToken)
			}
		}
		return err
	}
	config := new(Config)
	if err := config.fromSession(sess, sessionToken); err != nil {
		return err
	}
	return h.UpdateConfig(config)

}

func getField(i interface{}, path string) (*gjson.Result, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(i); err != nil {
		return nil, err
	}
	result := gjson.GetBytes(b.Bytes(), path)
	return &result, nil
}

func (c *Config) fromSession(session *cloud.Session, token string) error {
	email, err := getField(session.Identity.Traits, "email")
	if err != nil {
		return err
	}
	name, err := getField(session.Identity.Traits, "name")
	if err != nil {
		return err
	}

	c.Version = ConfigVersion
	c.SessionToken = token
	c.IdentityTraits = Identity{
		Email: email.String(),
		Name:  name.String(),
		ID:    uuid.FromStringOrNil(session.Identity.Id),
	}
	return nil
}

func (h *CommandHelper) Authenticate(ctx context.Context) error {
	if h.isQuiet {
		return stderrors.New("can not sign in or sign up when flag --quiet is set")
	}

	config, err := h.getConfig()
	if stderrors.Is(err, ErrNoConfig) {
		config = new(Config)
	} else if err != nil {
		return err
	}

	if len(config.SessionToken) > 0 {
		if !h.noConfirm {
			ok, err := cmdx.AskScannerForConfirmation(fmt.Sprintf("You are signed in as %q already. Do you wish to authenticate with another account?", config.IdentityTraits.Email), h.Stdin, h.VerboseErrWriter)
			if err != nil {
				return err
			} else if !ok {
				return nil
			}
			_, _ = fmt.Fprintf(h.VerboseErrWriter, "Ok, signing you out!\n")
		}

		if err := h.ClearConfig(); err != nil {
			return err
		}
	}

	c, err := NewOryProjectClient()
	if err != nil {
		return err
	}

	signIn, err := cmdx.AskScannerForConfirmation("Do you want to sign in to an existing Ory Network account?", h.Stdin, h.VerboseErrWriter)
	if err != nil {
		return err
	}

	if signIn {
		if err := h.signin(ctx, c, ""); err != nil {
			return err
		}
	} else {
		_, _ = fmt.Fprintln(h.VerboseErrWriter, "Great to have you! Let's create a new Ory Network account.")
		if err := h.signup(ctx, c); err != nil {
			return err
		}
	}

	config, err = h.getConfig()
	if err != nil {
		return err
	}
	if len(config.SessionToken) == 0 {
		return fmt.Errorf("unable to authenticate")
	}

	_, _ = fmt.Fprintf(h.VerboseErrWriter, "You are now signed in as: %s\n", config.IdentityTraits.Email)
	return nil
}

func (h *CommandHelper) ClearConfig() error {
	return h.UpdateConfig(new(Config))
}
