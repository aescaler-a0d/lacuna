/*
 * @File: workers.go
 * @Date: 2019-06-05 11:21:18
 * @OA:   Antonio Escalera
 * @CA:   Antonio Escalera
 * @Time: 2019-06-07 13:40:02
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"net"
	"sync"

	"github.com/angelofdeauth/lacuna/read"
	log "github.com/sirupsen/logrus"
)

type workFnc func(w workGenerator, ip net.IP) net.IP
type workGenerator struct {
	debug  bool
	filter []net.IP
	n      string
	s      *net.IPNet
	thread int
	wf     workFnc
}

//type output struct {
//	data <-chan net.IP
//	errc <-chan error
//}
type output struct {
	data net.IP
	errc error
}

func newWorkGenerator(debug bool, filter []net.IP, n string, s *net.IPNet, wf workFnc) workGenerator {
	w := workGenerator{
		debug:  debug,
		filter: filter,
		n:      n,
		s:      s,
		wf:     wf,
	}
	return w
}
func newOutput(data net.IP, errc error) output {
	o := output{data: data, errc: errc}
	return o
}

func generateSubnetReader(w workGenerator) (<-chan net.IP, <-chan error) {
	return read.ReadSubnetIntoChan(w.s, w.debug)
}
func generateWorkers(w workGenerator) <-chan output {

	if w.debug {
		//fmt.Printf("Generating workers for workGenerator: %v\n\n", w.wf)
		log.WithFields(log.Fields{
			"workGenerator": w.n,
		}).Infof("Generating workers for workGenerator: %v\n", w.n)
	}
	// generate 1 subnet reader for every group of workers
	i, _ := generateSubnetReader(w)

	o := make(chan output)

	// make output from worker chans

	// create worker wait group
	var wg sync.WaitGroup

	// add poolSize to worker wait group
	wg.Add(poolSize)

	// create goroutines for worker function
	// for some reason this is giving me trouble.
	for j := 0; j < poolSize; j++ {
		go func(j int) {
			w.thread = j
			for ip := range i {
				if w.debug {
					log.WithFields(log.Fields{
						"Worker": j,
						"Name":   w.n,
						"IP":     ip,
					}).Tracef("Worker %v for workFn %v processing IP: %v\n", j, w.n, ip)
				}
				filtered := w.wf(w, ip)
				if w.debug {
					log.WithFields(log.Fields{
						"Worker":   j,
						"Name":     w.n,
						"Filtered": filtered,
					}).Tracef("Worker %v for workFn %v received: %v from filter\n", j, w.n, filtered)
				}
				if filtered != nil {
					ipout := newOutput(ip, nil)
					o <- ipout
					if w.debug {
						log.WithFields(log.Fields{
							"Worker": j,
							"IP":     ip,
							"O":      o,
						}).Debugf("Worker %v sent IP: %v to output chan %v\n", j, ip, o)
					}
				}
			}
			wg.Done()
		}(j)
	}

	go func() {
		wg.Wait()
		log.WithFields(log.Fields{
			"Name": w.n,
		}).Infof("Closing worker channels for worker %v\n", w.n)
		close(o)
	}()
	return o
}
