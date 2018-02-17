package pool

// request represents the necessary info for making a call to an API.
// type request struct {
// 	method string
// 	path   string
// 	body   []byte
// 	out    interface{}
// 	echan  chan error
// }

// type requestQ struct {
// 	queue []request
// 	mux   sync.Mutex
// }

// try increasing the number of connections - see what happens to overall rate
// try increasing the chunk size - see what happens to the overall rate
// whatever diff is bigger, take it

// defaults for s3
// https://docs.aws.amazon.com/cli/latest/topic/s3-config.html

// New creates a new connection pool. It returns a scheduler to add requests
// to the queue.
// func New() func(string, string, []byte, interface{}) chan error {
// 	var rq requestQ
// 	var chunks []chunk
// 	// p := sync.Pool{} // TODO: use something else - this isn't reliable

// 	for i := 0; i < MaxConnections; i++ {
// 		// b := newBackend(apiURL, JWT)
// 		// p.Put(b)
// 	}

// 	return func(method string, path string, body []byte, out interface{}) chan error {
// 		echan := make(chan error)
// 		r := request{
// 			method: method,
// 			path:   path,
// 			body:   body,
// 			out:    out,
// 			echan:  echan,
// 		}
// 		rq.mux.Lock()
// 		rq.queue = append(rq.queue, r)
// 		rq.mux.Unlock()

// 		go execute(&rq, &p, echan)
// 		return echan
// 	}
// }

// when pulling new files off the queue, determine chunk size
// when starting a new chunk, determine # connections

// func initiateUpload() {
// 	// determine part size
// 	// initiate multipart
// }

// func chunkFile() {
// 	// determine part size
// 	// create Chunks
// }

// func uploadChunk(svc glacieriface.GlacierAPI, chunk *Chunk) {
// 	// time the chunk upload
// }

// func completeUpload() {
// 	// complete multipart
// 	// record upload rate / chunk
// }

// func execute(rq *requestQ, p *sync.Pool, echan chan error) {
// 	rq.mux.Lock()
// 	if len(rq.queue) == 0 {
// 		rq.mux.Unlock()
// 		return
// 	}

// 	if client, ok := p.Get().(something); ok {
// 		req := rq.queue[0]
// 		rq.queue = rq.queue[1:]
// 		rq.mux.Unlock()

// 		err := client.Call(req.method, req.path, req.body, req.out)
// 		p.Put(client)

// 		if err != nil {
// 			req.echan <- err
// 			return
// 		}
// 		req.echan <- nil
// 		go execute(rq, p, echan)
// 	} else {
// 		rq.mux.Unlock()
// 	}
// }

// func newBackend(apiURL string, JWT string) something {
// 	// log := &hoth.DefaultLogger{}

// 	// b, err := hoth.GetBackend(apiURL, "", "", "mc", log)
// 	if err != nil {
// 		log.Fatalf("Unable to create backend client for %s | %s", apiURL, err)
// 	}
// 	b.SetJWT(JWT)

// 	return b
// }
