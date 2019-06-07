/*
 * @File: show.go
 * @Date: 2019-05-30 17:33:26
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-07 13:41:47
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package cmd

import (
	"fmt"
	"net"

	"github.com/angelofdeauth/lacuna/find"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var outputFormat string

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows free IP addresses in a local subnet.",
	Long: `Outputs free IP addresses in a local subnet.
Defaults to human readable output, can be configured 
using flags.`,
	Run: func(cmd *cobra.Command, args []string) {

		//initialize host variables. Should put this somewhere else in the future
		err := getHostFacts()
		if err != nil {
			fmt.Println(err)
		}

		//Shows the subnet being tested for free IPs
		showSubnet(&Subnet)
		showSubnetFreeIPs(&Subnet, ArpFile, AwFile, Debug)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	showCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "output format (default \"human-readable\")")
}

func showSubnet(s *net.IPNet) {
	log.WithFields(log.Fields{
		"Subnet": s,
	}).Info("", s)
}
func showSubnetFreeIPs(s *net.IPNet, a string, w string, debug bool) {
	ips, err := find.FreeIPs(s, a, w, debug)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Panic("Find Free IPs Returned Error\n")
	}
	log.WithFields(log.Fields{
		"Free":  ips,
		"Count": len(ips),
	}).Info("Find Free IPs Returned\n")
}
