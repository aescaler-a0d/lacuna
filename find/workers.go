/*
 * @File: workers.go
 * @Date: 2019-06-05 11:21:18
 * @OA:   Antonio Escalera
 * @CA:   Antonio Escalera
 * @Time: 2019-06-05 15:00:25
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"fmt"
	"net"
	"sync"

	"github.com/angelofdeauth/lacuna/read"
)

type workFnc func(w workGenerator, ip net.IP) net.IP
type workGenerator struct {
	debug  bool
	filter []net.IP
	i      chan net.IP
	o      chan net.IP
	s      *net.IPNet
	thread int
	wf     workFnc
}
type output struct {
	data <-chan net.IP
	errc <-chan error
}

func newWorkGenerator(debug bool, filter []net.IP, s *net.IPNet, wf workFnc) workGenerator {
	w := workGenerator{
		debug:  debug,
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
	a, b := read.ReadSubnetIntoChan(w.s, w.debug)
	var wg sync.WaitGroup
	wg.Add(2)
	d := make(chan net.IP)
	e := make(chan error)
	o := newOutput(d, e)
	go func() {
		for dat := range a {
			d <- dat
		}
		wg.Done()
	}()
	go func() {
		for err := range b {
			e <- err
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(d)
		close(e)
	}()

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

	// make worker error chan
	errc := make(chan error, 1)

	// set worker input to subnet input chan
	go func() {
		for ip := range s.data {
			w.i <- ip
		}
		close(w.i)
	}()

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
			for ip := range w.i {
				if filtered := w.wf(w, ip); filtered != nil {
					w.o <- ip
				}
			}
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
