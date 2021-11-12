package cloudx

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	kratos "github.com/ory/kratos-client-go"
	"github.com/ory/x/cmdx"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
	"golang.org/x/term"
	"io"
	"strings"
	"syscall"
	"time"
)

func getLabel(attrs *kratos.UiNodeInputAttributes, node *kratos.UiNode) string {
	if attrs.Name == "password_identifier" {
		return fmt.Sprintf("%s: ", "Email")
	} else if node.Meta.Label != nil {
		return fmt.Sprintf("%s: ", node.Meta.Label.Text)
	} else if attrs.Label != nil {
		return fmt.Sprintf("%s: ", attrs.Label.Text)
	}
	return fmt.Sprintf("%s: ", attrs.Name)
}

func renderForm(stdin io.Reader, stdout io.Writer, ui kratos.UiContainer, method string, out interface{}) (err error) {
	for _, message := range ui.Messages {
		_, _ = fmt.Fprintf(stdout, "%s\n", message.Text)
	}

	for _, node := range ui.Nodes {
		for _, message := range node.Messages {
			_, _ = fmt.Fprintf(stdout, "%s\n", message.Text)
		}
	}

	values := json.RawMessage(`{}`)
	for _, node := range ui.Nodes {
		if node.Group != method {
			continue
		}

		switch node.Type {
		case "input":
			attrs := node.Attributes.UiNodeInputAttributes
			switch attrs.Type {
			case "button":
				continue
			case "submit":
				continue
			}

			if attrs.Name == "traits.consent.tos" {
				for !cmdx.AskForConfirmation(getLabel(attrs, &node), stdin, stdout) {
				}
				values, err = sjson.SetBytes(values, attrs.Name, time.Now().UTC().Format(time.RFC3339))
				if err != nil {
					return err
				}
				continue
			}

			switch attrs.Type {
			case "checkbox":
				result := cmdx.AskForConfirmation(getLabel(attrs, &node), stdin, stdout)
				var err error
				values, err = sjson.SetBytes(values, attrs.Name, result)
				if err != nil {
					return err
				}
			case "password":
				var password string
				for password == "" {
					_,_=fmt.Fprint(stdout, getLabel(attrs, &node))
					v, err := term.ReadPassword(syscall.Stdin)
					if err != nil {
						return err
					}
					password = strings.ReplaceAll(string(v), "\n", "")
					fmt.Println("")
				}

				values, err = sjson.SetBytes(values, attrs.Name, password)
				if err != nil {
					return err
				}
			default:
				var value string
				for value == "" {
					_,_ = fmt.Fprint(stdout, getLabel(attrs, &node))
					v, err := bufio.NewReader(stdin).ReadString('\n')
					if err != nil {
						return err
					}
					value = strings.ReplaceAll(v, "\n", "")
				}

				values, err = sjson.SetBytes(values, attrs.Name, value)
				if err != nil {
					return err
				}
			}
		default:
			// Do nothing
		}
	}

	values, err = sjson.SetBytes(values, "method", method)
	if err != nil {
		return err
	}

	return errors.WithStack(json.NewDecoder(bytes.NewBuffer(values)).Decode(out))
}
