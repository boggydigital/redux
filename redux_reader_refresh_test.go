package redux

import (
	"github.com/boggydigital/testo"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	testDir            = "redux_test"
	logRecordsFilename = "_log.gob"
)

// copied from kevlar
func logRecordsCleanup() error {
	logPath := filepath.Join(os.TempDir(), testDir, logRecordsFilename)
	if _, err := os.Stat(logPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if err := os.Remove(logPath); err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(os.TempDir(), testDir))
}

func TestRedux_FileModTime(t *testing.T) {
	start := time.Now().UTC().Unix()

	wrdx, err := NewWriter(filepath.Join(os.TempDir(), testDir), "test")
	testo.Error(t, err, false)

	rdx := wrdx.(*redux)
	testo.Nil(t, rdx, false)

	// first test: compare unmodified redux mod time
	// expected result: mod time should be less than start of the test

	rmt, err := rdx.FileModTime()
	testo.Error(t, err, false)
	testo.CompareInt64(t, rmt, start, testo.Less)

	// second test: add a value and compare redux mod time
	// expected result: mod time should be greater or equal than start of the test

	testo.Error(t, rdx.AddValues("test", "k1", "v1"), false)

	rmt, err = rdx.FileModTime()
	testo.Error(t, err, false)
	testo.CompareInt64(t, rmt, start, testo.GreaterOrEqual)

	// cleanup

	testo.Error(t, rdx.CutValues("test", "k1", "v1"), false)
	err = rdx.kv.Cut("test")
	testo.Error(t, err, false)

	testo.Error(t, logRecordsCleanup(), false)
}

func TestRedux_RefreshReader(t *testing.T) {
	wrdx, err := NewWriter(filepath.Join(os.TempDir(), testDir), "test")
	testo.Error(t, err, false)

	rdx := wrdx.(*redux)
	testo.Nil(t, rdx, false)

	// first test: set modTime to force Refresh and try RefreshReader
	// expected result: redux is refreshed and modTime is updated

	testo.Nil(t, rdx.kv, false)

	testo.Error(t, rdx.kv.Set("test", strings.NewReader("test")), false)
	err = rdx.kv.Cut("test")
	testo.Error(t, err, false)

	rrdx, err := rdx.RefreshReader()
	testo.Error(t, err, false)

	var ok bool
	rdx, ok = rrdx.(*redux)
	testo.EqualValues(t, ok, true)

	mt, err := rdx.FileModTime()
	testo.Error(t, err, false)
	testo.CompareInt64(t, mt, -1, testo.Greater)

	// second time: don't change modTime and try to RefreshReader again
	// expected result: no refresh is necessary and modTime is unchanged

	startModTime := mt

	rrdx, err = rdx.RefreshReader()
	testo.Error(t, err, false)

	rdx, ok = rrdx.(*redux)
	testo.EqualValues(t, ok, true)

	newMt, err := rdx.FileModTime()
	testo.Error(t, err, false)
	testo.EqualValues(t, newMt, startModTime)

	testo.Error(t, logRecordsCleanup(), false)
}
