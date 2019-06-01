/*
 * @File: arpwatch.go
 * @Date: 2019-05-29 15:17:38
 * @OA:   Antonio Escalera
 * @CA:   antonioe
 * @Time: 2019-05-31 03:05:00
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package read

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"strconv"
)

type AWData struct {
	Macaddr string
	Ipaddr  net.IP
	Time    int64
	Name    string
}

func ReadAWDataIntoStruct(s string) []AWData {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	c := bytes.NewReader(b)
	r := NewFieldsReader(c)
	r.Comma = '\t'
	r.FieldsPerRecord = -1
	recs, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	a := make([]AWData, len(recs))

	for idx, val := range recs {
		a[idx].Macaddr = recs[idx][0]
		a[idx].Ipaddr = net.ParseIP(recs[idx][1])
		a[idx].Time, _ = strconv.ParseInt(recs[idx][2], 10, 32)
		switch len(val) {
		case 3:
			a[idx].Name = "UNKNOWN"
		case 4:
			a[idx].Name = recs[idx][3]
		case 0:
			break
		}
	}
	return a
}
