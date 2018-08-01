package cmd

import (
	"regexp"

	"github.com/qri-io/qri/lib"
	"github.com/qri-io/qri/repo"
	"github.com/spf13/cobra"
)

// NewGetCommand creates a new `qri search` command that searches for datasets
func NewGetCommand(f Factory, ioStreams IOStreams) *cobra.Command {
	o := &GetOptions{IOStreams: ioStreams}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get elements of qri datasets",
		Long:  ``,
		Annotations: map[string]string{
			"group": "dataset",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(f, args); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&o.Format, "format", "f", "yaml", "set output format [json, yaml]")
	cmd.Flags().BoolVar(&o.Concise, "concise", false, "print output without indentation, only applies to json format")

	return cmd
}

// GetOptions encapsulates state for the get command
type GetOptions struct {
	IOStreams

	Refs    []string
	Path    string
	Format  string
	Concise bool

	DatasetRequests *lib.DatasetRequests
}

// isDatasetField checks if a string is a dataset field or not
var isDatasetField = regexp.MustCompile("(?i)commit|structure|body|meta|viz|transform")

// Complete adds any missing configuration that can only be added just before calling Run
func (o *GetOptions) Complete(f Factory, args []string) (err error) {
	if len(args) > 0 {
		if isDatasetField.MatchString(args[0]) {
			o.Path = args[0]
			args = args[1:]
		}
	}
	o.Refs = args
	o.DatasetRequests, err = f.DatasetRequests()
	return
}

// Run executes the get command
func (o *GetOptions) Run() (err error) {
	var refs []repo.DatasetRef
	for _, refstr := range o.Refs {
		ref, err := repo.ParseDatasetRef(refstr)
		if err != nil {
			return err
		}
		refs = append(refs, ref)
	}

	p := &lib.SelectParams{
		Path:    o.Path,
		Format:  o.Format,
		Concise: o.Concise,
	}

	res := []byte{}
	if err = o.DatasetRequests.Select(p, &res); err != nil {
		return err
	}

	_, err = o.Out.Write(res)
	return err
}
