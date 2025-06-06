package redux

import (
	"encoding/gob"
	"errors"
	"github.com/boggydigital/kevlar"
	"io"
	"sync"
)

func ErrUnknownAsset(asset string) error {
	return errors.New("unknown redux asset " + asset)
}

type redux struct {
	dir    string
	kv     kevlar.KeyValues
	assets []string
	akv    map[string]map[string][]string
	lmt    map[string]int64
	mtx    *sync.Mutex
}

func newRedux(dir string, assets ...string) (*redux, error) {
	kv, err := kevlar.New(dir, kevlar.GobExt)
	if err != nil {
		return nil, err
	}

	assetKeyValues := make(map[string]map[string][]string)
	amts := make(map[string]int64)
	for _, asset := range assets {
		if assetKeyValues[asset], err = loadAsset(kv, asset); err != nil {
			return nil, err
		}
		amts[asset] = kv.LogModTime(asset)
	}

	return &redux{
		kv:     kv,
		dir:    dir,
		assets: assets,
		akv:    assetKeyValues,
		lmt:    amts,
		mtx:    new(sync.Mutex),
	}, nil
}

func loadAsset(kv kevlar.KeyValues, asset string) (map[string][]string, error) {

	if !kv.Has(asset) {
		return make(map[string][]string), nil
	}

	arc, err := kv.Get(asset)
	if err != nil {
		return nil, err
	}
	defer arc.Close()

	var nkv map[string][]string
	if arc != nil {
		if err := gob.NewDecoder(arc).Decode(&nkv); err == io.EOF {
			// empty reduction - do nothing, it'll be initialized below
		} else if err != nil {
			return nil, err
		}
	}

	if nkv == nil {
		nkv = make(map[string][]string)
	}

	return nkv, nil
}
