package redux

import (
	"io"
)

type Readable interface {
	MustHave(assets ...string) error
	Keys(asset string) []string
	HasAsset(asset string) bool
	HasKey(asset, key string) bool
	HasValue(asset, key, val string) bool
	GetAllValues(asset, key string) ([]string, bool)
	GetLastVal(asset, key string) (string, bool)
	FileModTime() (int64, error)
	RefreshReader() (Readable, error)
	MatchAsset(asset string, terms []string, scope []string, options ...MatchOption) []string
	Match(query map[string][]string, options ...MatchOption) []string
	Sort(ids []string, desc bool, sortBy ...string) ([]string, error)
	Export(w io.Writer, keys ...string) error
}

type Writeable interface {
	Readable
	AddValues(asset, key string, values ...string) error
	BatchAddValues(asset string, keyValues map[string][]string) error
	ReplaceValues(asset, key string, values ...string) error
	BatchReplaceValues(asset string, keyValues map[string][]string) error
	CutKeys(asset string, keys ...string) error
	CutValues(asset, key string, values ...string) error
	BatchCutValues(asset string, keyValues map[string][]string) error
	RefreshWriter() (Writeable, error)
}
