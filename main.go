/*
 * @File: main.go
 * @Date: 2019-05-29 18:16:36
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-04 20:34:20
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package main

import "github.com/angelofdeauth/lacuna/cmd"

var (
	VERSION = "v0.0.2"
)

func main() {
	cmd.Execute(VERSION)
}
