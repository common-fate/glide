package gconfig

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

// CLIPrompt prompts the user to enter a value for the config varsiable
// in a CLI context. If the config variable implements Defaulter, the
// default value is returned and the user is not prompted for any input.
// @TODO I think that this cli prompt should actually be defined elsewhere like in the cli cmd
// gconfig should be a pure package concerned with providing the API to read and write configs
// CLI IO is a layer built ontop of gconfig using the public API
func (f *Field) CLIPrompt() error {
	grey := color.New(color.FgHiBlack)
	msg := f.Key()
	if f.Usage() != "" {
		msg = f.Usage() + " " + grey.Sprintf("(%s)", msg)
	}

	// @TODO work out how to integrate the optional prompt here
	// you shoudl be able to choose to set or unset
	// if you choose to set, it should use a default if it exists
	// By design, we can't have an optional secret as they are mutually exclusive.
	var p survey.Prompt
	// if this value is a secret, use a password prompt to key the secret out of the terminal history
	if f.IsSecret() {
		if f.Get() != "" {
			confMsg := msg + " would you like to update this secret?"
			p = &survey.Confirm{
				Message: confMsg,
			}
			var doUpdate bool
			err := survey.AskOne(p, &doUpdate)
			if err != nil {
				return err
			}
			if !doUpdate {
				return nil
			}
		}
		p = &survey.Password{
			Message: msg,
		}
	} else {
		p = &survey.Input{
			Message: msg,
			Default: f.Get(),
		}
	}
	var val string
	err := survey.AskOne(p, &val)
	if err != nil {
		return err
	}
	// set the value.
	return f.Set(val)
}
