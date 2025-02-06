package redux

import (
	"github.com/boggydigital/testo"
	"slices"
	"strings"
	"testing"
)

func TestRedux_Export(t *testing.T) {
	rdx := mockRedux(t)

	sb := &strings.Builder{}
	testo.EqualValues(t, sb.Len(), 0)
	testo.Error(t, rdx.Export(sb, slices.Collect(rdx.Keys("a1"))...), false)
	testo.CompareInt64(t, int64(sb.Len()), 0, testo.Greater)
}
