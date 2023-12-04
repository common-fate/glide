package table

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// Table is a lightweight wrapper around
// text/tabwriter for printing CLI tables.
type Table struct {
	tw *tabwriter.Writer
}

func New(w io.Writer) *Table {
	tw := tabwriter.NewWriter(w, 10, 1, 5, ' ', 0)
	return &Table{tw: tw}
}

func (t *Table) Columns(cols ...string) {
	var uppercase []string
	for _, c := range cols {
		uppercase = append(uppercase, strings.ToUpper(c))
	}

	tabbed := strings.Join(uppercase, "\t")
	fmt.Fprintln(t.tw, tabbed)
}

func (t *Table) Row(data ...string) {
	tabbed := strings.Join(data, "\t")
	fmt.Fprintln(t.tw, tabbed)
}

func (t *Table) Flush() error {
	return t.tw.Flush()
}
