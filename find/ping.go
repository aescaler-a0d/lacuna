/*
 * @File: ping.go
 * @Date: 2019-05-30 17:47:27
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-07 15:25:05
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"net"

	"github.com/angelofdeauth/lacuna/read"
	log "github.com/sirupsen/logrus"
)

func PingHosts(w workGenerator, ip net.IP) (net.IP, error) {
	if w.debug {
		log.WithFields(log.Fields{
			"Name":   w.n,
			"Worker": w.thread,
			"IP":     ip,
		}).Tracef("%v worker %v: Processing IP: %v\n", w.n, w.thread, ip)
	}
	for _, v := range w.filter {
		// check if the IP in the channel is in the filter
		if ip.Equal(v) {
			if w.debug {
				log.WithFields(log.Fields{
					"Name":   w.n,
					"Worker": w.thread,
					"IP":     ip,
				}).Debugf("%v worker %v: IP in filter: %v\n", w.n, w.thread, ip)
			}
			return ip, nil
		}
	}
	return nil, nil
}

func PingHostsRead(w workGenerator, ip net.IP) (net.IP, error) {
	if w.debug {
		log.WithFields(log.Fields{
			"Name":   w.n,
			"Worker": w.thread,
			"IP":     ip,
		}).Tracef("%v worker %v: Processing IP: %v\n", w.n, w.thread, ip)
	}

	pinger, err := read.PingHostnames(ip, attempts, interval, timeout)
	if err != nil {
		return nil, err
	}
	if pinger.PacketsSent > 0 && pinger.PacketsRecv > 0 {
		if w.debug {
			log.WithFields(log.Fields{
				"Name":   w.n,
				"Worker": w.thread,
				"IP":     ip,
			}).Debugf("%v worker %v: IP: %v responded to ping\n", w.n, w.thread, ip)
		}
		return ip, nil
	} else {
		if w.debug {
			log.WithFields(log.Fields{
				"Name":   w.n,
				"Worker": w.thread,
				"IP":     ip,
			}).Tracef("%v worker %v: IP: %v did not respond to ping\n", w.n, w.thread, ip)
		}
		return nil, nil
	}
}
