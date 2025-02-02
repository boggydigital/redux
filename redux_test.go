package redux

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/testo"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func reduxCleanup(assets ...string) error {
	for _, asset := range assets {
		rdxPath := filepath.Join(os.TempDir(), testDir, asset+kevlar.GobExt)
		if _, err := os.Stat(rdxPath); err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if err := os.Remove(rdxPath); err != nil {
			return err
		}
	}
	return logRecordsCleanup()
}

func mockRedux(t *testing.T) *redux {
	rdx := &redux{
		dir: filepath.Join(os.TempDir(), testDir),
		akv: map[string]map[string][]string{
			"a1": {
				"k1": {"v11"},
				"k2": {"v21", "v22"},
				"k3": {"v31", "v32", "v33"},
			},
			"a2": {
				"k4": {"v41", "v42", "v43", "v44"},
				"k5": {"v51", "v52", "v53", "v54", "v55"},
			},
		},
		mtx: new(sync.Mutex),
	}

	var err error
	rdx.kv, err = kevlar.New(rdx.dir, kevlar.GobExt)

	testo.Error(t, err, false)

	return rdx
}
