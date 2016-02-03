package binary_test

import (
	"fmt"
	"github.com/toshaf/remora/comms"
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

func Test_Duplex(t *testing.T) {
	dir := path.Join(os.TempDir(), "remora-test")

	p := binary.NewProvider(dir)

	sin, sout, err := p.Server("test", comms.Duplex)
	if err != nil {
		t.Error(err)
	}
	if sin == nil {
		t.Error("sin nil")
	}
	if sout == nil {
		t.Error("sout nil")
	}
	defer sin.Close()
	defer sout.Close()

	cin, cout, err := p.Client("test", comms.Duplex)
	if err != nil {
		t.Error(err)
	}
	if cin == nil {
		t.Error("cin nil")
	}
	if cout == nil {
		t.Error("cout nil")
	}
	defer cin.Close()
	defer cout.Close()

	// server
	go func() {
		var req Req
		err := sin.Recv(&req)
		if err != nil {
			t.Error(err)
		}

		err = sout.Send(Res{Id: 456, To: req.Id})
	}()

	// client
	err = cout.Send(Req{Id: 123})
	if err != nil {
		t.Error(err)
	}

	var res Res
	err = cin.Recv(&res)
	if err != nil {
		t.Error(err)
	}

	if res != (Res{Id: 456, To: 123}) {
		t.Errorf("Wrong response: %v", res)
	}
}

func Test_Server_connection_deletes_existing_files(t *testing.T) {
	ioutil.WriteFile(path.Join(os.TempDir(), "remora-test/test.in"), []byte{42}, 0666)
	ioutil.WriteFile(path.Join(os.TempDir(), "remora-test/test.out"), []byte{42}, 0666)

	provider := binary.NewProvider(path.Join(os.TempDir(), "remora-test"))

	in, out, err := provider.Server("test", comms.Duplex)
	if err != nil {
		t.Error(err)
	}
	if in == nil {
		t.Error("in nil")
	}
	if out == nil {
		t.Error("out nil")
	}
}

func Benchmark_binary_comms(b *testing.B) {
	dir := path.Join(os.TempDir(), "remora-test")
	provider := binary.NewProvider(dir)

	// start the server
	gate := make(chan struct{})
	go func() {
		in, out, err := provider.Server("echo", comms.Duplex)
		if err != nil {
			panic(err)
		}

		close(gate)

		for {
			var msg int
			err := in.Recv(&msg)
			if err != nil {
				panic(err)
			}

			err = out.Send(msg)
			if err != nil {
				panic(err)
			}
		}
	}()

	// wait for server
	<-gate

	// set up client
	in, out, err := provider.Client("echo", comms.Duplex)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = out.Send(i)
		if err != nil {
			panic(err)
		}

		var r int
		err = in.Recv(&r)
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
