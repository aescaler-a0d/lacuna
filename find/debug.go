/*
 * @File: debug.go
 * @Date: 2019-05-30 17:47:33
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-07 14:10:32
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"net"

	log "github.com/sirupsen/logrus"
)

// Sends hosts in ArpHosts to output chan
func Debug(w workGenerator, ip net.IP) (net.IP, error) {
	if w.debug {
		log.WithFields(log.Fields{
			"Name":   w.n,
			"Worker": w.thread,
			"IP":     ip,
		}).Tracef("%v worker %v: Processing IP: %v\n", w.n, w.thread, ip)
	}
	return ip, nil
}
