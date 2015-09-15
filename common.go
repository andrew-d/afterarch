package afterarch

import (
	"encoding/binary"
)

var (
	trailerMagic = [4]byte{'A', 'A', '0', '1'}
	trailerSize  = int64(binary.Size(&trailer{}))
)

type trailer struct {
	Magic       [4]byte
	ArchiveSize int64
}
