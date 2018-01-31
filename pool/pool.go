package pool

var partSize = int64(1 << 20) // 1MB

type Files struct {
	maxConnections int
	failed         []Chunk
	UploadID       *string
	Vault          string
}

type Chunk struct {
	Buf      []byte
	Start    int64
	End      int64
	Checksum string
}
