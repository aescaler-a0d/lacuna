/*
 * @File: find.go
 * @Date: 2019-05-30 17:32:24
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-04 20:28:29
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/angelofdeauth/lacuna/read"
)

var (
	timeout  = 100 * time.Millisecond
	attempts = 1
	poolSize = 2 // * runtime.NumCPU()
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
		fmt.Printf("Generating workers for workGenerator: %v\n\n", w.wf)
	}
	// generate 1 subnet reader for every group of workers
	s := generateSubnetReader(w)

	w.i = make(chan net.IP)
	w.o = make(chan net.IP)

	// set worker input to subnet input chan
	go func() {
		for ip := range s.data {
			w.i <- ip
		}
	}()

	// make worker error chan
	errc := make(chan error, 1)

	// read subnet error chan into worker error chan
	go func() {
		for err := range s.errc {
			if err != nil {
				errc <- err
			}
		}
	}()

	// make output from worker chans
	o := newOutput(w.o, errc)

	// create worker wait group
	var wg sync.WaitGroup

	// add poolSize to worker wait group
	wg.Add(poolSize)

	// create goroutines for worker function
	// THIS ONE RIGHT HERE OFFICER
	// for some reason this is giving me trouble.
	for j := 0; j < poolSize; j++ {
		go func(j int) {
			w.thread = j
			errc <- w.wf(w)
			wg.Done()
		}(j)
	}

	go func() {
		wg.Wait()
		fmt.Printf("Closing worker channels for worker %v\n", w.wf)
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
		fmt.Printf("Arp Filter: %v\nCount: %v\n\n", arp, len(arp))
	}

	//
	aws, err := read.ReadAWDataIntoSlice(w, debug)
	if err != nil {
		return nil, err
	}
	if debug {
		fmt.Printf("ArpWatch Filter: %v\nCount: %v\n\n", aws, len(aws))
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
	alive, err := waitForPipeline(done, debug, outSlice...)
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

func waitForPipeline(done <-chan struct{}, debug bool, o ...output) ([]net.IP, error) {
	// waitForPipeline is blocked from outputing until
	// all of the channels are closed

	if debug {
		fmt.Printf("waitForPipeline called for outputs: %v\n\n", o)
	}
	out := multiplexChans(debug, o...)

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

func multiplexChans(debug bool, channels ...output) output {
	// Count down as each channel closes. When hits zero - close ch.
	var dwg sync.WaitGroup
	dwg.Add(len(channels))
	var ewg sync.WaitGroup
	ewg.Add(len(channels))
	// The channel to output to.
	dch := make(chan net.IP, len(channels))
	ech := make(chan error, len(channels))
	mout := newOutput(dch, ech)

	// Make one go per channel.
	for i, c := range channels {
		go func(c <-chan net.IP, i int) {
			// Pump it.
			for x := range c {
				dch <- x
				fmt.Printf("Sent %v to data channel output\n", x)
			}
			// It closed.
			dwg.Done()
			fmt.Printf("Data wait group done called on mux worker %v\n", i)
		}(c.data, i)
	}

	for i, c := range channels {
		go func(c <-chan error, i int) {
			// Pump it.
			for x := range c {
				ech <- x
				fmt.Printf("Sent %v to errc channel output\n", x)
			}
			// It closed.
			ewg.Done()
			fmt.Printf("Errc wait group done called on mux worker %v\n", i)
		}(c.errc, i)
	}
	// Close the channel when the pumping is finished.
	go func() {
		// Wait for everyone to be done.
		dwg.Wait()
		// Close.
		close(dch)
		if debug {
			fmt.Printf("Close called on mux channel data")
		}

		ewg.Wait()
		close(ech)
		if debug {
			fmt.Printf("Close called on mux channel errc")
		}
	}()
	return mout
}

// func multiplexChans(done <-chan struct{}, debug bool, chans ...output) output {
// 	var wg sync.WaitGroup
// 	outdata, outerrc := make(chan net.IP, len(chans)), make(chan error, len(chans))
// 	outoutput := newOutput(outdata, outerrc)
// 	outp := func(c output) {
// 		for {
// 			select {
// 			case outdata <- c.data:
// 				outdata <- outd
// 			case oute := <-c.errc:
// 				outerrc <- oute
// 			case <-done:
// 				if debug {
// 					fmt.Println("")
// 				}
// 				return
// 			}
//
// 		}
// 		wg.Done()
// 	}
// 	wg.Add(len(chans))
// 	for _, c := range chans {
// 		go outp(c)
// 	}
//
// 	go func() {
// 		wg.Wait()
// 		close(outdata)
// 		close(outerrc)
// 	}()
// 	return outoutput
// }

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
