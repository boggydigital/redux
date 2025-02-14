package redux

import (
	"github.com/boggydigital/wits"
	"io"
	"maps"
	"slices"
)

func (rdx *redux) Export(w io.Writer, keys ...string) error {

	sortedAssets := slices.Sorted(maps.Keys(rdx.akv))

	skv := make(wits.SectionKeyValues)

	rdx.mtx.Lock()
	defer rdx.mtx.Unlock()

	for _, key := range keys {
		skv[key] = make(wits.KeyValues)
		for _, asset := range sortedAssets {
			if values := rdx.akv[asset][key]; len(values) > 0 {
				skv[key][asset] = values
			}
		}
	}

	return skv.Write(w)
}
