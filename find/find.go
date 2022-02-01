/*
 * @File: find.go
 * @Date: 2019-05-30 17:32:24
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-07 15:42:10
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"errors"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/aescaler-a0d/lacuna/read"
	log "github.com/sirupsen/logrus"
)

var (
	timeout  = 100 * time.Millisecond
	attempts = 1
	poolSize = 2 * runtime.NumCPU()
	interval = 100 * time.Millisecond
)

func FreeIPs(s *net.IPNet, a string, w string, debug bool) ([]net.IP, error) {

	var outSlice []<-chan output

	// need to create buffered channel of size of file
	// chunk file across readers
	// and return the channel when the buffer is full.
	arp, err := read.ReadArpDataIntoSlice(a, debug)
	if err != nil {
		return nil, err
	}
	if debug {
		log.WithFields(log.Fields{
			"ArpFilter": arp,
			"Count":     len(arp),
		}).Info("Arp Filter Returned")
	}

	//
	aws, err := read.ReadAWDataIntoSlice(w, debug)
	if err != nil {
		return nil, err
	}
	if debug {
		log.WithFields(log.Fields{
			"ArpWatchFilter": aws,
			"Count":          len(aws),
		}).Info("ArpWatch Filter Returned")
	}

	dnsrw := newWorkGenerator(debug, []net.IP{}, "DNSRead", s, DnsHostsRead)
	dnsro := generateWorkers(dnsrw)
	dnsfilter, dnsrerr := waitForPipeline(debug, dnsro)
	if dnsrerr != nil {
		return nil, err
	}
	if debug {
		log.WithFields(log.Fields{
			"DNSFilter": dnsfilter,
			"Count":     len(dnsfilter),
		}).Info("DNS Filter Returned")
	}

	pingrw := newWorkGenerator(debug, []net.IP{}, "PingRead", s, PingHostsRead)
	pingro := generateWorkers(pingrw)
	pingfilter, pingrerr := waitForPipeline(debug, pingro)
	if pingrerr != nil {
		return nil, err
	}
	if debug {
		log.WithFields(log.Fields{
			"PingFilter": pingfilter,
			"Count":      len(pingfilter),
		}).Info("Ping Filter Returned")
	}

	//fakew := newWorkGenerator(debug, []net.IP{}, "Debug", s, Debug)
	//fakeo := generateWorkers(fakew)
	//outSlice = append(outSlice, fakeo)

	// create a new work generator for ArpHosts
	arpw := newWorkGenerator(debug, arp, "ArpHosts", s, ArpHosts)
	// generate workers and return chans for data and err
	arpo := generateWorkers(arpw)
	// add output to outSlice
	outSlice = append(outSlice, arpo)

	// create a new work generator for ArpWatch
	awsw := newWorkGenerator(debug, aws, "ArpWatch", s, ArpWatch)
	awso := generateWorkers(awsw)
	outSlice = append(outSlice, awso)

	// create a new work generator for DNS
	dnsw := newWorkGenerator(debug, dnsfilter, "DnsHosts", s, DnsHosts)
	dnso := generateWorkers(dnsw)
	outSlice = append(outSlice, dnso)

	// create a new work generator for Ping
	pingw := newWorkGenerator(debug, pingfilter, "PingHosts", s, PingHosts)
	pingo := generateWorkers(pingw)
	outSlice = append(outSlice, pingo)

	// wait for pipeline to error or finish
	alive, err := waitForPipeline(debug, outSlice...)
	if err != nil {
		return nil, err
	}
	if len(alive) == 0 {
		return nil, errors.New("Error: No hosts found in subnet")
	}
	log.WithFields(log.Fields{
		"Alive": alive,
		"Count": len(alive),
	}).Debug("Alive Hosts Returned")

	dead, err := maskAliveHosts(alive, s, debug)
	if err != nil {
		return nil, err
	}
	if len(dead) == 0 {
		return nil, errors.New("Error: No free IPs found in subnet")
	}

	return dead, nil
}

func waitForPipeline(debug bool, o ...<-chan output) ([]net.IP, error) {
	// waitForPipeline is blocked from outputing until
	// all of the channels are closed

	if debug {
		log.WithFields(log.Fields{
			"Output": o,
		}).Debugf("waitForPipeline called for output: %v\n", o)
	}
	out := multiplexChans(debug, o...)

	var alive []net.IP
	for ip := range out {
		if ip.errc != nil {
			return nil, ip.errc
		}
		if debug {
			log.WithFields(log.Fields{
				"IP": ip.data,
			}).Tracef("waitForPipeline ranging over IP: %v\n", ip.data)
		}
		alive = append(alive, ip.data)
	}

	return alive, nil
}

func multiplexChans(debug bool, workers ...<-chan output) <-chan output {
	// init wait groups
	var wg sync.WaitGroup
	wg.Add(len(workers))

	mout := make(chan output)

	// Make one go per channel.
	for i, c := range workers {
		if debug {
			log.WithFields(log.Fields{
				"Worker":    i,
				"InChannel": c,
			}).Debugf("starting mux worker: %v for channel: %v\n", i, c)
		}
		go func(i int, c <-chan output) {
			// Pump it.
			for x := range c {
				if x.data != nil {
					xout := newOutput(x.data, nil)
					mout <- xout
					if debug {
						log.WithFields(log.Fields{
							"Worker":    i,
							"InChannel": c,
							"Data":      x.data,
						}).Tracef("Mux worker %v: Sent data %v from channel %v to channel output\n", i, x.data, c)
					}
				}
				if x.errc != nil {
					xout := newOutput(nil, x.errc)
					mout <- xout
					if debug {
						log.WithFields(log.Fields{
							"Worker":    i,
							"InChannel": c,
							"Error":     x.errc,
						}).Tracef("Mux worker %v: Sent error %v from channel %v to channel output\n", i, x.errc, c)
					}
				}
			}
			// It closed.
			wg.Done()
			if debug {
				log.WithFields(log.Fields{
					"Worker":    i,
					"InChannel": c,
				}).Debugf("Mux worker %v servicing channel %v wait group done\n", i, c)
			}
		}(i, c)
	}
	// Close the channel when the pumping is finished.
	go func() {
		// Wait for everyone to be done.

		wg.Wait()
		// Close.
		close(mout)
		if debug {
			log.WithFields(log.Fields{}).Debug("Close called on mux channel output\n")
		}
	}()
	return mout
}

func maskAliveHosts(alive []net.IP, s *net.IPNet, debug bool) ([]net.IP, error) {

	//should rewrite this to use
	sRo, _ := generateSubnetReader(newWorkGenerator(debug, nil, "maskAliveHosts", s, nil))

	var sub []net.IP

	for ip := range sRo {
		var x bool
		for _, h := range alive {
			if h.Equal(ip) || ip == nil {
				x = true
			}
		}
		if !x {
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
