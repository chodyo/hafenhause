package hafenhause

import (
	"path"
	"strings"
)

type docpath string

const (
	BedtimePath         docpath = "bedtime"
	BedtimePeoplePath   docpath = "bedtime/people"
	BedtimeDefaultsPath docpath = "bedtime/static/defaults"
)

func (h hafenhausedb) combine(docPath docpath, docNames ...string) docpath {
	names := strings.Join(docNames, "/")
	fullPath := path.Join(h.collection, string(docPath), names)
	return docpath(fullPath)
}

// func (h hafenhausedb) docref(docPath docpath, docNames ...string) *firestore.DocumentRef {
// 	names := strings.Join(docNames, "/")
// 	fullPath := path.Join(h.collection, string(docPath), names)
// 	return h.client.Doc(fullPath)
// }
