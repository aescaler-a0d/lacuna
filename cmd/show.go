/*
 * @File: show.go
 * @Date: 2019-05-30 17:33:26
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-01 20:23:41
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

		//initialize host variables. Should put this somewhere else in the future
		err := getHostFacts()
		if err != nil {
			fmt.Println(err)
		}

		//Shows the subnet being tested for free IPs
		showSubnet(&Subnet, Debug)
		aH, err := getArpHosts(&Subnet, ArpFile, Debug)
		if err != nil {
			fmt.Println(err)
		}
		wH, err := getArpWatch(&Subnet, AwFile, Debug)
		if err != nil {
			fmt.Println(err)
		}
		iH, err := getArpWitch(&Subnet, Debug)
		if err != nil {
			fmt.Println(err)
		}
		dH, err := getDnsHosts(&Subnet, Debug)
		if err != nil {
			fmt.Println(err)
		}
		pH, err := getPingHosts(&Subnet, NetwInterface, Debug)
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

		showSubnetFreeIPs(&Subnet, nH, Debug)
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

func showSubnet(s *net.IPNet, debug bool) {
	fmt.Printf("Subnet: %v\n\n", s)
}

func getArpHosts(s *net.IPNet, r string, debug bool) ([]net.IP, error) {
	arp, err := find.ArpHosts(s, r, debug)
	if err != nil {
		return []net.IP{}, err
	}
	if debug {
		fmt.Printf("ArpHosts: %v\nCount: %v\n\n", arp, len(arp))
	}
	return arp, nil
}

func getArpWatch(s *net.IPNet, a string, debug bool) ([]net.IP, error) {
	aws, err := find.ArpWatch(s, a, debug)
	if err != nil {
		return []net.IP{}, err
	}
	if debug {
		fmt.Printf("ArpWatch: %v\nCount: %v\n\n", aws, len(aws))
	}
	return aws, nil
}

func getArpWitch(s *net.IPNet, debug bool) ([]net.IP, error) {
	awi, err := find.ArpWitch(s, debug)
	if err != nil {
		return []net.IP{}, err
	}
	if debug {
		fmt.Printf("ArpWitch: %v\nCount: %v\n\n", awi, len(awi))
	}
	return awi, nil
}

func getDnsHosts(s *net.IPNet, debug bool) ([]net.IP, error) {
	dns, err := find.DnsHosts(s, debug)
	if err != nil {
		return []net.IP{}, err
	}
	if debug {
		fmt.Printf("Dns: %v\nCount: %v\n\n", dns, len(dns))
	}
	return dns, nil
}

func getPingHosts(s *net.IPNet, i string, debug bool) ([]net.IP, error) {
	pin, err := find.PingHosts(s, i, debug)
	if err != nil {
		return []net.IP{}, err
	}
	if debug {
		fmt.Printf("Ping: %v\nCount: %v\n\n", pin, len(pin))
	}
	return pin, nil
}
func showSubnetFreeIPs(s *net.IPNet, n find.NetworkHosts, debug bool) {
	ips, err := find.FreeIPs(s, n)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Free IPs: %v\n\nCount: %v\n\n", ips, len(ips))
}
