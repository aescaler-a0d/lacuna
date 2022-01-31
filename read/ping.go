/*
 * @File: ping.go
 * @Date: 2019-06-07 15:18:17
 * @OA:   Antonio Escalera
 * @CA:   Antonio Escalera
 * @Time: 2019-06-07 15:24:04
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package read

import (
	"net"
	"time"

	"github.com/go-ping/ping"
)

func PingHostnames(ip net.IP, attempts int, interval time.Duration, timeout time.Duration) (*ping.Pinger, error) {
	pinger, err := ping.NewPinger(ip.String())
	if err != nil {
		return nil, err
	}
	pinger.Count = attempts
	pinger.Interval = interval
	pinger.Timeout = timeout
	pinger.SetPrivileged(true)
	pinger.Run()

	return pinger, nil
}
