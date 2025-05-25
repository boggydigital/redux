package redux

import (
	"errors"
	"iter"
	"maps"
	"path"
	"slices"
	"time"
)

func NewReader(dir string, assets ...string) (Readable, error) {
	return newRedux(dir, assets...)
}

func (rdx *redux) MustHave(assets ...string) error {
	for _, asset := range assets {
		if !rdx.HasAsset(asset) {
			return ErrUnknownAsset(asset)
		}
	}
	return nil
}

func (rdx *redux) Keys(asset string) iter.Seq[string] {
	rdx.mtx.Lock()
	defer rdx.mtx.Unlock()

	return maps.Keys(rdx.akv[asset])
}

func (rdx *redux) Len(asset string) int {
	return len(rdx.akv[asset])
}

func (rdx *redux) HasAsset(asset string) bool {
	_, ok := rdx.akv[asset]
	return ok
}

func (rdx *redux) HasKey(asset, key string) bool {
	rdx.mtx.Lock()
	defer rdx.mtx.Unlock()

	if akr, ok := rdx.akv[asset]; ok {
		_, ok = akr[key]
		return ok
	}
	return false
}

func (rdx *redux) HasValue(asset, key, val string) bool {
	rdx.mtx.Lock()
	defer rdx.mtx.Unlock()

	if akr, ok := rdx.akv[asset]; ok {
		if kr, ok := akr[key]; ok {
			return slices.Contains(kr, val)
		}
		return false
	}
	return false
}

func (rdx *redux) GetAllValues(asset, key string) ([]string, bool) {
	if !rdx.HasAsset(asset) {
		return nil, false
	}

	rdx.mtx.Lock()
	defer rdx.mtx.Unlock()

	if rdx.akv[asset] == nil {
		return nil, false
	}

	val, ok := rdx.akv[asset][key]
	return val, ok
}

func (rdx *redux) GetLastVal(asset, key string) (string, bool) {
	if values, ok := rdx.GetAllValues(asset, key); ok && len(values) > 0 {
		return values[len(values)-1], true
	}
	return "", false
}

func (rdx *redux) ParseLastValTime(asset, key string) (time.Time, error) {
	if lvs, ok := rdx.GetLastVal(asset, key); ok && lvs != "" {
		if dt, err := time.Parse(time.RFC3339, lvs); err == nil {
			return dt, nil
		} else {
			return time.Time{}, err
		}
	} else {
		return time.Time{}, errors.New("value not found for " + path.Join(asset, key))
	}
}
