package binary_test

import (
	"fmt"
	"github.com/toshaf/remora/comms/binary"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type Req struct {
	Id int
}

type Res struct {
	Id, To int
}

func Test_Connection(t *testing.T) {
	dir := path.Join(os.TempDir(), "remora-test")

	p := binary.NewProvider(dir)

	server, err := p.Server("test")
	if err != nil {
		t.Error(err)
	}
	if server == nil {
		t.Error("server nil")
	}

	client, err := p.Client("test")
	if err != nil {
		t.Error(err)
	}
	if client == nil {
		t.Error("client nil")
	}
	defer server.Close()
	defer client.Close()

	// server
	go func() {
		var req Req
		err := server.Recv(&req)
		if err != nil {
			t.Error(err)
		}

		err = server.Send(Res{Id: 456, To: req.Id})
	}()

	// client
	err = client.Send(Req{Id: 123})
	if err != nil {
		t.Error(err)
	}

	var res Res
	err = client.Recv(&res)
	if err != nil {
		t.Error(err)
	}

	if res != (Res{Id: 456, To: 123}) {
		t.Errorf("Wrong response: %v", res)
	}
}

func Test_Idempotent_Close(t *testing.T) {
	fmt.Fprintf(os.Stderr, "creating provider...")
	provider := binary.NewProvider(path.Join(os.TempDir(), "remora-test"))
	fmt.Fprintf(os.Stderr, "done\ncreating server...")
	server, err := provider.Server("test")
	if err != nil {
		t.Error(err)
	}
	fmt.Fprintf(os.Stderr, "done\nclosing once...")
	err = server.Close()
	if err != nil {
		t.Error(err)
	}
	fmt.Fprintf(os.Stderr, "done\nclosing twice...")
	err = server.Close()
	if err != nil {
		t.Error(err)
	}
	fmt.Fprintf(os.Stderr, "done\n")
}

func Test_Server_connection_deletes_existing_files(t *testing.T) {
	ioutil.WriteFile(path.Join(os.TempDir(), "remora-test/test.in"), []byte{42}, 0666)
	ioutil.WriteFile(path.Join(os.TempDir(), "remora-test/test.out"), []byte{42}, 0666)

	provider := binary.NewProvider(path.Join(os.TempDir(), "remora-test"))

	server, err := provider.Server("test")
	if err != nil {
		t.Error(err)
	}
	if server == nil {
		t.Error("server nil")
	}
}

func Benchmark_binary_comms(b *testing.B) {
	dir := path.Join(os.TempDir(), "remora-test")
	provider := binary.NewProvider(dir)

	// start the server
	gate := make(chan struct{})
	go func() {
		server, err := provider.Server("echo")
		if err != nil {
			panic(err)
		}

		close(gate)

		for {
			var msg int
			err := server.Recv(&msg)
			if err != nil {
				panic(err)
			}

			err = server.Send(msg)
			if err != nil {
				panic(err)
			}
		}
	}()

	// wait for server
	<-gate

	// set up client
	client, err := provider.Client("echo")
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = client.Send(i)
		if err != nil {
			panic(err)
		}

		var r int
		err = client.Recv(&r)
		if err != nil {
			panic(err)
		}

		if r != i {
			panic(fmt.Sprintf("Expected %d, got %d", i, r))
		}
	}
}

// for comparison with piped echo
func Benchmark_direct_echo(b *testing.B) {

	out := make(chan int)
	in := make(chan int)

	// start the "server"
	go func() {
		for {
			in <- <-out
		}
	}()

	// run the "client"
	for i := 0; i < b.N; i++ {
		out <- i
		r := <-in

		if r != i {
			panic(fmt.Sprintf("Expected %d, got %d", i, r))
		}
	}
}
