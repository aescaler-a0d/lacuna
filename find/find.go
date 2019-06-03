/*
 * @File: find.go
 * @Date: 2019-05-30 17:32:24
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-03 17:57:49
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/angelofdeauth/gopher/read"
)

var (
	timeout  = 100 * time.Millisecond
	attempts = 1
	poolSize = 2 * runtime.NumCPU()
	interval = 100 * time.Millisecond
)

type workFnc func(w workGenerator) error
type workGenerator struct {
	debug  bool
	done   <-chan struct{}
	filter []net.IP
	i      chan net.IP
	o      chan net.IP
	s      *net.IPNet
	thread int
	wf     workFnc
}
type output struct {
	data chan net.IP
	errc chan error
}

func newWorkGenerator(debug bool, done <-chan struct{}, filter []net.IP, s *net.IPNet, wf workFnc) workGenerator {
	w := workGenerator{
		debug:  debug,
		done:   done,
		filter: filter,
		s:      s,
		wf:     wf,
	}
	return w
}
func newOutput(data chan net.IP, errc chan error) output {
	o := output{data: data, errc: errc}
	return o
}

func generateSubnetReader(w workGenerator) output {
	o := newOutput(read.ReadSubnetIntoChan(w.s, w.done, w.debug))
	return o
}
func generateWorkers(w workGenerator) output {

	if w.debug {
		fmt.Printf("Generating workers for workGenerator: %v\n", w)
	}
	// generate 1 subnet reader for every group of workers
	s := generateSubnetReader(w)

	// split subnet input chan from subnet err chan
	i, serrc := s.data, s.errc

	// set worker input to subnet input chan
	w.i = i

	// make worker output chan
	w.o = make(chan net.IP)

	// make worker error chan
	errc := make(chan error, 1)

	// make output from worker chans
	o := newOutput(w.o, errc)

	// create worker wait group
	var wg sync.WaitGroup

	// add poolSize to worker wait group
	wg.Add(poolSize)

	// create goroutines for worker function
	for j := 0; j < poolSize; j++ {
		go func(j int) {
			w.thread = j
			errc <- w.wf(w)
			wg.Done()
		}(j)
	}

	go func() {
		for err := range serrc {
			if err != nil {
				errc <- err
			}
		}
		wg.Wait()
		close(w.o)
		close(errc)
	}()
	return o
}

func FreeIPs(s *net.IPNet, a string, w string, debug bool) ([]net.IP, error) {
	done := make(chan struct{})
	defer close(done)

	var outSlice []output

	// need to create buffered channel of size of file
	// chunk file across readers
	// and return the channel when the buffer is full.
	arp, err := read.ReadArpDataIntoSlice(a, debug)
	if err != nil {
		return nil, err
	}
	if debug {
		fmt.Printf("Arp Filter: %v\n\nCount: %i", arp, len(arp))
	}
	aws, err := read.ReadAWDataIntoSlice(w, debug)
	if err != nil {
		return nil, err
	}
	if debug {
		fmt.Printf("ArpWatch Filter: %v\n\nCount: %i", aws, len(aws))
	}

	// create a new work generator for ArpHosts
	arpw := newWorkGenerator(debug, done, arp, s, ArpHosts)
	// generate workers and return chans for data and err
	arpo := generateWorkers(arpw)
	// transform data
	// currently not needed
	//arpt := chanTransformer(arpo, arpw)
	// add output to outSlice
	outSlice = append(outSlice, arpo)

	awsw := newWorkGenerator(debug, done, aws, s, ArpWatch)
	awso := generateWorkers(awsw)
	//awst := chanTransformer(awso, awsw)
	outSlice = append(outSlice, awso)

	// wait for pipeline to error or finish
	alive, err := waitForPipeline(done, outSlice...)
	if err != nil {
		return nil, err
	}
	if len(alive) == 0 {
		return nil, errors.New("Error: No hosts found in subnet")
	}

	dead, err := maskAliveHosts(done, alive, s, debug)
	if err != nil {
		return nil, err
	}
	if len(dead) == 0 {
		return nil, errors.New("Error: No free IPs found in subnet")
	}

	return dead, nil

}

func maskAliveHosts(done <-chan struct{}, alive []net.IP, s *net.IPNet, debug bool) ([]net.IP, error) {
	sRW := newWorkGenerator(debug, done, nil, s, nil)
	sRo := generateSubnetReader(sRW)

	// blocks until no subnetreader errors are found
	for err := range sRo.errc {
		if err != nil {
			return nil, err
		}
	}

	// consumes subnetreader data into array
	var sub []net.IP
	for ip := range sRo.data {
		var x bool
		for _, h := range alive {
			if ip.Equal(h) {
				x = true
			}
		}
		if x {
			sub = append(sub, ip)
		}
	}
	return sub, nil

}

func waitForPipeline(done <-chan struct{}, o ...output) ([]net.IP, error) {
	// waitForPipeline is blocked from outputing until
	// all of the channels are closed
	out := mergeChans(done, o...)
	for err := range out.errc {
		if err != nil {
			return nil, err
		}
	}

	var alive []net.IP
	for ip := range out.data {
		alive = append(alive, ip)
	}

	return alive, nil
}

func mergeChans(done <-chan struct{}, chans ...output) output {
	var wg sync.WaitGroup
	outdata, outerrc := make(chan net.IP, len(chans)), make(chan error, len(chans))
	outoutput := newOutput(outdata, outerrc)
	outp := func(c output) {
		for {
			select {
			case outd := <-c.data:
				outdata <- outd
			case oute := <-c.errc:
				outerrc <- oute
			case <-done:
				return
			}

		}
		wg.Done()
	}
	wg.Add(len(chans))
	for _, c := range chans {
		go outp(c)
	}

	go func() {
		wg.Wait()
		close(outdata)
		close(outerrc)
	}()
	return outoutput
}

//func chanTransformer(i output, w workGenerator) output {
//	dc := make(chan net.IP)
//	ec := make(chan error, 1)
//	go func() {
//		defer close(dc)
//		defer close(ec)
//		for d := range i {
//			// Send the data to the output channel but return early
//			// if the context has been cancelled.
//			select {
//			case dc <- d.data:
//			case ec <- d.errc:
//			case <-w.done:
//				return
//			}
//		}
//	}()
//	o := newOutput(dc, ec)
//	return o
//}
