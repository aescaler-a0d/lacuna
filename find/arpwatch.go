/*
 * @File: arpwatch.go
 * @Date: 2019-05-30 17:47:47
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-07 13:44:10
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"net"

	log "github.com/sirupsen/logrus"
)

// Sends hosts in ArpWatch to output chan.
func ArpWatch(w workGenerator, ip net.IP) net.IP {
	if w.debug {
		log.WithFields(log.Fields{
			"Name":   w.n,
			"Worker": w.thread,
			"IP":     ip,
		}).Tracef("%v worker %v: Processing IP: %v\n", w.n, w.thread, ip)
	}
	for _, v := range w.filter {
		// check if the ip in the channel is in the filter
		if ip.Equal(v) {
			if w.debug {
				log.WithFields(log.Fields{
					"Name":   w.n,
					"Worker": w.thread,
					"IP":     ip,
				}).Debugf("%v worker %v: IP in filter: %v\n", w.n, w.thread, ip)
			}
			return ip
		}
	}
	return nil
}
