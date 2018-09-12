package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/deis/duffle/pkg/duffle/home"
)

type credentialRemoveCmd struct {
	names []string
	home  home.Home
	out   io.Writer
}

func newCredentialRemoveCmd(w io.Writer) *cobra.Command {
	rm := &credentialRemoveCmd{out: w}

	cmd := &cobra.Command{
		Use:   "remove [NAME]",
		Short: "remove one or more credential set",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("This command requires at least 1 argument: name of credential set")
			}
			rm.names = args
			rm.home = home.Home(homePath())
			return rm.run()

		},
	}
	return cmd
}

func (rm *credentialRemoveCmd) run() error {
	var removeErrors []string
	var notFound []string
	credentialSets := findCredentialSets(rm.home.Credentials())

	for _, name := range rm.names {
		if path, ok := credentialSets[name]; ok {
			if err := removeCredentialSet(path); err != nil {
				removeErrors = append(removeErrors, fmt.Sprintf("Failed to remove credential set %s: %v", name, err))
			} else {
				fmt.Fprintf(rm.out, "Removed credential set: %s\n", name)
			}
		} else {
			notFound = append(notFound, name)
		}
	}

	if len(notFound) > 0 {
		notFoundError := fmt.Sprintf("Unable to find credential set(s): %v", strings.Join(notFound, ", "))
		removeErrors = append(removeErrors, notFoundError)
	}

	if len(removeErrors) > 0 {
		return fmt.Errorf(strings.Join(removeErrors, "\n"))
	}

	return nil
}

func removeCredentialSet(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}