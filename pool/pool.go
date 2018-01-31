package pool

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

// dynamically change the partsize depending on rate of failure, size of file
