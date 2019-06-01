/*
 * @File: main.go
 * @Date: 2019-05-29 18:16:36
 * @OA:   antonioe
 * @CA:   antonioe
 * @Time: 2019-05-30 17:33:46
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package main

import "github.com/angelofdeauth/gopher/cmd"

var (
	VERSION = "v0.0.2"
)

func main() {
	cmd.Execute(VERSION)
}
