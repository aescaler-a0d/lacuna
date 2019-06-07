/*
 * @File: subnet.go
 * @Date: 2019-06-02 01:02:57
 * @OA:   Antonio Escalera
 * @CA:   Antonio Escalera
 * @Time: 2019-06-06 13:42:51
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package read

import (
	"net"
)

type SubnetFunc func(s *net.IPNet, ip net.IP, err error) error

//func ReadSubnetIntoChans(s *net.IPNet, d <-chan struct{}, debug bool) (<-chan net.IP, <-chan error) {
func ReadSubnetIntoChan(s *net.IPNet, debug bool) (<-chan net.IP, <-chan error) {

	if debug {
		//fmt.Printf("Reading subnet %v into chan\n\n", s)
	}
	ips := make(chan net.IP)
	errc := make(chan error, 1)
	go func() { // HL
		// Close the paths channel after Walk returns.
		defer close(ips) // HL
		defer close(errc)
		// No select needed for this send, since errc is buffered.
		errc <- HostsInSubnet(s, func(s *net.IPNet, ip net.IP, err error) error {
			if err != nil {
				return err
			}
			ips <- ip
			//fmt.Printf("IP: %v read into input chan: %v\n", ip, ips)
			return nil
		})
	}()
	return ips, errc
}

func HostsInSubnet(s *net.IPNet, subFn SubnetFunc) error {
	ip, ipnet, err := net.ParseCIDR(s.String())
	if err != nil {
		err = subFn(s, ip, err)
	} else {
		err = hostsInSubnet(ipnet, ip, subFn)
	}
	return err
}

func hostsInSubnet(s *net.IPNet, ip net.IP, subFn SubnetFunc) (err error) {
	for ip := ip.Mask(s.Mask); s.Contains(ip); ip = inc(ip) {
		err = subFn(s, ip, nil)
	}
	// remove network address and broadcast address
	return err
}

func inc(ip net.IP) net.IP {
	incIP := make([]byte, len(ip))
	copy(incIP, ip)
	for j := len(incIP) - 1; j >= 0; j-- {
		incIP[j]++
		if incIP[j] > 0 {
			break
		}
	}
	return incIP
}
