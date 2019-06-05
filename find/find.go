/*
 * @File: find.go
 * @Date: 2019-05-30 17:32:24
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-05 15:02:06
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

func FreeIPs(s *net.IPNet, a string, w string, debug bool) ([]net.IP, error) {

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
	arpw := newWorkGenerator(debug, arp, s, ArpHosts)
	// generate workers and return chans for data and err
	arpo := generateWorkers(arpw)
	// transform data
	// currently not needed
	//arpt := chanTransformer(arpo, arpw)
	// add output to outSlice
	outSlice = append(outSlice, arpo)

	awsw := newWorkGenerator(debug, aws, s, ArpWatch)
	awso := generateWorkers(awsw)
	//awst := chanTransformer(awso, awsw)
	outSlice = append(outSlice, awso)

	// wait for pipeline to error or finish
	alive, err := waitForPipeline(debug, outSlice...)
	if err != nil {
		return nil, err
	}
	if len(alive) == 0 {
		return nil, errors.New("Error: No hosts found in subnet")
	}

	dead, err := maskAliveHosts(alive, s, debug)
	if err != nil {
		return nil, err
	}
	if len(dead) == 0 {
		return nil, errors.New("Error: No free IPs found in subnet")
	}

	return dead, nil

}

func waitForPipeline(debug bool, o ...output) ([]net.IP, error) {
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
		if debug {
			fmt.Printf("waitForPipeline ranging over err: %v\n", err)
		}
	}

	if debug {
		fmt.Printf("waitForPipeline finished err chan read")
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
		if debug {
			fmt.Printf("starting mux worker: %v for data channel: %v\n", i, c.data)
		}
		go func(i int) {
			// Pump it.
			for x := range c.data {
				dch <- x
				if debug {
					fmt.Printf("Sent %v to data channel output\n", x)
				}
			}
			// It closed.
			dwg.Done()
			if debug {
				fmt.Printf("Data wait group done called on mux worker %v\n", i)
			}
		}(i)
	}

	for i, c := range channels {
		if debug {
			fmt.Printf("starting mux worker: %v for errc channel: %v\n", i, c.errc)
		}
		go func(i int) {
			// Pump it.
			for x := range c.errc {
				ech <- x
				fmt.Printf("Sent %v to errc channel output\n", x)
			}
			// It closed.
			ewg.Done()
			if debug {
				fmt.Printf("Errc wait group done called on mux worker %v\n", i)
			}
		}(i)
	}
	// Close the channel when the pumping is finished.
	go func() {
		// Wait for everyone to be done.
		dwg.Wait()
		// Close.
		close(dch)
		if debug {
			fmt.Println("Close called on mux channel data")
		}

		ewg.Wait()
		close(ech)
		if debug {
			fmt.Println("Close called on mux channel errc")
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

func maskAliveHosts(alive []net.IP, s *net.IPNet, debug bool) ([]net.IP, error) {
	sRW := newWorkGenerator(debug, nil, s, nil)
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
