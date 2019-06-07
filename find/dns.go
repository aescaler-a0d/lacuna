/*
 * @File: dns.go
 * @Date: 2019-05-30 17:48:02
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-07 14:37:03
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"net"

	"github.com/angelofdeauth/lacuna/read"
	log "github.com/sirupsen/logrus"
)

func DnsHosts(w workGenerator, ip net.IP) (net.IP, error) {
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

func DnsHostsRead(w workGenerator, ip net.IP) (net.IP, error) {
	if w.debug {
		log.WithFields(log.Fields{
			"Name":   w.n,
			"Worker": w.thread,
			"IP":     ip,
		}).Tracef("%v worker %v: Processing IP: %v\n", w.n, w.thread, ip)
	}

	filter, err := read.DnsHostnames(w.debug, ip)
	if err != nil {
		return nil, err
	}
	if len(filter) == 0 {
		if w.debug {
			log.WithFields(log.Fields{
				"Name":   w.n,
				"Worker": w.thread,
				"IP":     ip,
			}).Tracef("%v worker %v: Returned no hosts for ip: %v\n", w.n, w.thread, ip)
		}
		return nil, nil
	}
	if w.debug {
		log.WithFields(log.Fields{
			"Name":   w.n,
			"Worker": w.thread,
			"IP":     ip,
			"Hosts":  filter,
		}).Debugf("%v worker %v: Returned %v for ip: %v\n", w.n, w.thread, filter, ip)
	}
	return ip, nil
}
