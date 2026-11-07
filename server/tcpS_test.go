package server

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"bufio"
)

func TestHandleRespInput(t *testing.T) {

	reader := bufio.NewReader(bytes.NewReader([]byte("*1\r\n$4\r\nPING\r\n")))

	s := Server{
		store: &Memory{
			mu:    sync.RWMutex{},
			KVMap: make(map[string]string),
		},
	}

	got, err := s.handleResp(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq := bytes.Equal(got, []byte("+PONG\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, "+PONG\r\n")
	}

}

func TestHandleRespEncoding(t *testing.T) {

	s := Server{
		store: &Memory{
			mu:    sync.RWMutex{},
			KVMap: make(map[string]string),
		},
	}

	reader := bufio.NewReader(bytes.NewReader([]byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n")))

	got, err := s.handleResp(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(got, []byte("$-1\r\n")) {
		t.Errorf("Invalid encoding = %q, Want = %q", got, []byte("$-1\r\n"))
	}

}

func TestHandleRespEcho(t *testing.T) {
	reader := bufio.NewReader(bytes.NewReader([]byte("*3\r\n$4\r\nECHO\r\n$5\r\nHello\r\n$5\r\nWorld\r\n")))

	s := Server{
		store: &Memory{
			mu:    sync.RWMutex{},
			KVMap: make(map[string]string),
		},
	}

	got, err := s.handleResp(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq := bytes.Equal(got, []byte("+Hello World\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, "+Hello World\r\n")
	}

}

func TestHandleRespSet(t *testing.T) {
	reader := bufio.NewReader(bytes.NewReader([]byte("*3\r\n$3\r\nSET\r\n$4\r\nName\r\n$4\r\nJohn\r\n")))

	s := Server{
		store: &Memory{
			mu:    sync.RWMutex{},
			KVMap: make(map[string]string),
		},
	}

	got, err := s.handleResp(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq := bytes.Equal(got, []byte("+OK\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, "+OK\r\n")
	}
}

func TestKeySetExpiryTimestamp(t *testing.T) {
	//we want to store keys
	// we receive a reader with options PX EX
	setReader := bufio.NewReader(bytes.NewReader([]byte("*5\r\n$3\r\nSET\r\n$4\r\nName\r\n$4\r\nJohn\r\n$2\r\nEX\r\n$1\r\n5\r\n")))
	s := Server{
		store: &Memory{
			mu:       sync.RWMutex{},
			KVMap:    make(map[string]string),
			Expirity: make(map[string]time.Time),
		},
	}

	got, err := s.handleResp(setReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq := bytes.Equal(got, []byte("+OK\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, "+OK\r\n")
	}

	getReader := bufio.NewReader(bytes.NewReader([]byte("*2\r\n$3\r\nGET\r\n$4\r\nName\r\n")))

	// esperamos n tiempo
	time.Sleep(5 * time.Second)
	// hacemos un get y visualizamos que se borra
	got, err = s.handleResp(getReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq = bytes.Equal(got, []byte("$-1\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, "$-1\r\n")
	}
}

func TestTTLKeyExpirity(t *testing.T) {

	s := Server{
		store: &Memory{
			mu:       sync.RWMutex{},
			KVMap:    make(map[string]string),
			Expirity: make(map[string]time.Time),
		},
	}

	setReader := bufio.NewReader(bytes.NewReader([]byte("*5\r\n$3\r\nSET\r\n$4\r\nName\r\n$4\r\nJohn\r\n$2\r\nEX\r\n$1\r\n5\r\n")))
	//the key is added
	got, err := s.handleResp(setReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq := bytes.Equal(got, []byte("+OK\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, "+OK\r\n")
	}

	time.Sleep(3 * time.Second)

	ttlReader := bufio.NewReader(bytes.NewReader([]byte("*2\r\n$3\r\nTTL\r\n$4\r\nName\r\n")))

	got, err = s.handleResp(ttlReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(bytes.Equal(got, []byte(":1\r\n")) || bytes.Equal(got, []byte(":2\r\n"))) {
		t.Errorf("TTL = %d, want between :1\r\n and :2\r\n", got)

	}

}

func TestSetExpire(t *testing.T) {

	s := Server{
		store: &Memory{
			mu:       sync.RWMutex{},
			KVMap:    make(map[string]string),
			Expirity: make(map[string]time.Time),
		},
	}

	setReader := bufio.NewReader(bytes.NewReader([]byte("*3\r\n$3\r\nSET\r\n$4\r\nName\r\n$4\r\nJohn\r\n")))

	got, err := s.handleResp(setReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq := bytes.Equal(got, []byte("+OK\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, "+OK\r\n")
	}

	expireReader := bufio.NewReader(bytes.NewReader([]byte("*3\r\n$6\r\nEXPIRE\r\n$4\r\nName\r\n$1\r\n4\r\n")))

	got, err = s.handleResp(expireReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq = bytes.Equal(got, []byte(":1\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, ":1\r\n")
	}

	ttlReader := bufio.NewReader(bytes.NewReader([]byte("*2\r\n$3\r\nTTL\r\n$4\r\nName\r\n")))
	got, err = s.handleResp(ttlReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(bytes.Equal(got, []byte(":3\r\n")) || bytes.Equal(got, []byte(":4\r\n"))) {
		t.Errorf("TTL = %d, want between :3\r\n and :4\r\n", got)
	}
}

func TestExpireSweep(t *testing.T) {

	s := Server{
		store: &Memory{
			mu:       sync.RWMutex{},
			KVMap:    make(map[string]string),
			Expirity: make(map[string]time.Time),
		},
	}

	setReader := bufio.NewReader(bytes.NewReader([]byte("*5\r\n$3\r\nSET\r\n$4\r\nName\r\n$4\r\nJohn\r\n$2\r\nPX\r\n$2\r\n50\r\n")))

	got, err := s.handleResp(setReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq := bytes.Equal(got, []byte("+OK\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, "+OK\r\n")
	}

	time.Sleep(time.Millisecond * 100)
	s.store.ExpireSweep()

	_, exist := s.store.KVMap["Name"]
	if exist {
		t.Errorf("the key should have been deleted by the expire sweep")
	}
}

func TestConcurrentSetGetStress(t *testing.T) {
	s := Server{
		store: &Memory{
			mu:       sync.RWMutex{},
			KVMap:    make(map[string]string),
			Expirity: make(map[string]time.Time),
		},
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			setReader := bufio.NewReader(bytes.NewReader([]byte("*3\r\n$3\r\nSET\r\n$4\r\nName\r\n$4\r\nJohn\r\n")))
			s.handleResp(setReader)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			getReader := bufio.NewReader(bytes.NewReader([]byte("*2\r\n$3\r\nGET\r\n$4\r\nName\r\n")))
			s.handleResp(getReader)
		}()
	}
	wg.Wait()
}

func TestDeleteKey(t *testing.T) {

	s := Server{
		store: &Memory{
			mu:       sync.RWMutex{},
			KVMap:    make(map[string]string),
			Expirity: make(map[string]time.Time),
		},
	}

	setReader1 := bufio.NewReader(bytes.NewReader([]byte("*5\r\n$3\r\nSET\r\n$4\r\nkey1\r\n$4\r\nfooo\r\n$2\r\nEX\r\n$1\r\n2\r\n")))

	got, err := s.handleResp(setReader1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq := bytes.Equal(got, []byte("+OK\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, "+OK\r\n")
	}

	setReader2 := bufio.NewReader(bytes.NewReader([]byte("*3\r\n$3\r\nSET\r\n$4\r\nkey2\r\n$4\r\nfooo\r\n")))

	got, err = s.handleResp(setReader2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq = bytes.Equal(got, []byte("+OK\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, "+OK\r\n")
	}

	delReader := bufio.NewReader(bytes.NewReader([]byte("*4\r\n$3\r\nDEL\r\n$4\r\nkey1\r\n$4\r\nkey2\r\n$4\r\nkey3\r\n")))

	got, err = s.handleResp(delReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq = bytes.Equal(got, []byte(":2\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, ":2\r\n")
	}
}

func TestExistsKey(t *testing.T) {

	s := Server{
		store: &Memory{
			mu:       sync.RWMutex{},
			KVMap:    make(map[string]string),
			Expirity: make(map[string]time.Time),
		},
	}

	setReader1 := bufio.NewReader(bytes.NewReader([]byte("*5\r\n$3\r\nSET\r\n$4\r\nkey1\r\n$4\r\nfooo\r\n$2\r\nPX\r\n$2\r\n50\r\n")))

	got, err := s.handleResp(setReader1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq := bytes.Equal(got, []byte("+OK\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, "+OK\r\n")
	}

	time.Sleep(1 * time.Second)

	setReader2 := bufio.NewReader(bytes.NewReader([]byte("*3\r\n$3\r\nSET\r\n$4\r\nkey2\r\n$4\r\nfooo\r\n")))

	got, err = s.handleResp(setReader2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	eq = bytes.Equal(got, []byte("+OK\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, "+OK\r\n")
	}

	exReader := bufio.NewReader(bytes.NewReader([]byte("*4\r\n$6\r\nEXISTS\r\n$4\r\nkey1\r\n$4\r\nkey2\r\n$4\r\nkey5\r\n")))

	got, err = s.handleResp(exReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eq = bytes.Equal(got, []byte(":1\r\n"))
	if !eq {
		t.Errorf("Response = %q, Want = %q", got, ":1\r\n")
	}

}
