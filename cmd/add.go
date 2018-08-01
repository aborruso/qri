package cmd

import (
	"github.com/qri-io/qri/lib"
	"github.com/qri-io/qri/repo"
	"github.com/spf13/cobra"
)


// TODO: Tests.


// NewAddCommand creates an add command
func NewAddCommand(f Factory, ioStreams IOStreams) *cobra.Command {
	o := &AddOptions{IOStreams: ioStreams}
	cmd := &cobra.Command{
		Use:        "add",
		Short:      "Add a dataset",
		Annotations: map[string]string{
			"group": "dataset",
		},
		Long: `
Add retrieves a dataset owned by another peer and adds it to your repo.`,
		Example: `  add a dataset named their_data, owned by other_peer:
  $ qri add other_peer/their_data`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(f); err != nil {
				return err
			}
			if err := o.Run(args); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}

// AddOptions encapsulates state for the add command
type AddOptions struct {
	IOStreams
	DatasetRequests *lib.DatasetRequests
}

// Complete adds any missing configuration that can only be added just before calling Run
func (o *AddOptions) Complete(f Factory) (err error) {
	if o.DatasetRequests, err = f.DatasetRequests(); err != nil {
		return
	}
	return nil
}

// Run adds another peer's dataset to this user's repo
func (o *AddOptions) Run(args []string) error {
	for _, arg := range args {
		ref, err := parseCmdLineDatasetRef(arg)
		if err != nil {
			return err
		}

		res := repo.DatasetRef{}
		if err = o.DatasetRequests.Add(&ref, &res); err != nil {
			return err
		}

		printDatasetRefInfo(o.Out, 1, res)
		printInfo(o.Out, "Successfully added dataset %s", ref)
	}

	return nil
}
