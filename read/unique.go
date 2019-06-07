/*
 * @File: unique.go
 * @Date: 2019-06-07 12:53:02
 * @OA:   Antonio Escalera
 * @CA:   Antonio Escalera
 * @Time: 2019-06-07 12:53:09
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */
package read

import "net"

func unique(inSlice []net.IP) []net.IP {
	keys := make(map[string]bool)
	list := []net.IP{}
	for _, entry := range inSlice {
		if _, value := keys[entry.String()]; !value {
			keys[entry.String()] = true
			list = append(list, entry)
		}
	}
	return list
}
