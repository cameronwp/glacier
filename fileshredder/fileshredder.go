package fileshredder

const startingPartSize = int64(1 << 20) // 1MB

// rip up files, watch how long they take, change part sizes, # pool connections
