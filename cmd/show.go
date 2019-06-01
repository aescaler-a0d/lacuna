/*
 * @File: show.go
 * @Date: 2019-05-30 17:33:26
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-01 13:23:02
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package cmd

import (
	"fmt"
	"net"

	"github.com/angelofdeauth/gopher/find"
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
		if Debug {
			showSubnet(&Subnet)
			showArpHosts(&Subnet, ArpFile)
			showArpWatch(&Subnet, AwFile)
			showDnsHosts(&Subnet)
			showPingHosts(&Subnet, NetwInterface)
		}
		showSubnetFreeIPs(&Subnet, AwFile, ArpFile, NetwInterface)
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
	fmt.Printf("Subnet: %v\n\n", s)
}

func showArpHosts(s *net.IPNet, r string) {
	arp := find.ArpHosts(s, r)
	for _, v := range arp {
		fmt.Println(v)
	}
	fmt.Printf("ArpHosts: %v\nCount: %v\n\n", arp, len(arp))
}

func showArpWatch(s *net.IPNet, a string) {
	aws := find.ArpWatch(s, a)
	for _, v := range aws {
		fmt.Println(v)
	}
	fmt.Printf("ArpWatch: %v\nCount: %v\n\n", aws, len(aws))
}

func showDnsHosts(s *net.IPNet) {
	dns, err := find.DnsHosts(s)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range dns {
		fmt.Println(v)
	}
	fmt.Printf("Dns: %v\nCount: %v\n\n", dns, len(dns))
}

func showPingHosts(s *net.IPNet, i string) {
	pin, err := find.PingHosts(s, i)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range pin {
		fmt.Println(v)
	}
	fmt.Printf("Ping: %v\nCount: %v\n\n", pin, len(pin))
}
func showSubnetFreeIPs(s *net.IPNet, a string, r string, i string) {
	aH := find.ArpHosts(s, r)
	//if err != nil {
	//	fmt.Println(err)
	//}
	wH := find.ArpWatch(s, a)
	//if err != nil {
	//	fmt.Println(err)
	//}
	iH := find.ArpWitch(s)

	dH, err := find.DnsHosts(s)
	if err != nil {
		fmt.Println(err)
	}
	pH, err := find.PingHosts(s, i)
	if err != nil {
		fmt.Println(err)
	}

	nH := find.NetworkHosts{
		"Arp":      aH,
		"ArpWatch": wH,
		"ArpWitch": iH,
		"Dns":      dH,
		"Ping":     pH,
	}
	ips, err := find.FreeIPs(s, nH)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Free IPs: %v\n\nCount: %v\n\n", ips, len(ips))
}
