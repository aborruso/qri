package cmd

import (
	"testing"
)

func TestGetComplete(t *testing.T) {
	streams, in, out, errs := NewTestIOStreams()
	setNoColor(true)

	f, err := NewTestFactory(nil)
	if err != nil {
		t.Errorf("error creating new test factory: %s", err)
		return
	}

	cases := []struct {
		args []string
		path string
		refs []string
		err  string
	}{
		{[]string{}, "", []string{}, ""},
		{[]string{"one arg"}, "", []string{"one arg"}, ""},
		{[]string{"commit", "peer/ds"}, "commit", []string{"peer/ds"}, ""},
		{[]string{"peer/ds_two", "peer/ds"}, "", []string{"peer/ds_two", "peer/ds"}, ""},
		{[]string{"foo", "peer/ds"}, "", []string{"foo", "peer/ds"}, ""},
		{[]string{"structure"}, "structure", []string{}, ""},
	}

	for i, c := range cases {
		opt := &GetOptions{
			IOStreams: streams,
		}

		opt.Complete(f, c.args)

		if c.err != errs.String() {
			t.Errorf("case %d, error mismatch. Expected: '%s', Got: '%s'", i, c.err, errs.String())
			ioReset(in, out, errs)
			continue
		}

		if !testSliceEqual(c.refs, opt.Refs) {
			t.Errorf("case %d, opt.Refs not set correctly. Expected: '%s', Got: '%s'", i, c.refs, opt.Refs)
			ioReset(in, out, errs)
			continue
		}

		if c.path != opt.Path {
			t.Errorf("case %d, opt.Path not set correctly. Expected: '%s', Got: '%s'", i, c.path, opt.Path)
			ioReset(in, out, errs)
			continue
		}

		if opt.DatasetRequests == nil {
			t.Errorf("case %d, opt.DatasetRequests not set.", i)
			ioReset(in, out, errs)
			continue
		}
		ioReset(in, out, errs)
	}
}
