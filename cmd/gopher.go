/*
 * @File: gopher.go
 * @Date: 2019-05-30 17:31:09
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-01 18:43:20
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	AwFile        string
	ArpFile       string
	Debug         bool
	NetwInterface string
	Subnet        net.IPNet
	VERSION       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopher",
	Short: "Find free IP addresses using arp, ping, and DNS.",
	Long: `Gopher is a utility that can find free IP addresses
in a given subnet. It has smart defaults that just
work if only one local network is configured. If
multiple networks are configured, the first is chosen
by default. If multiple interfaces are configured 
with networks, the first numerical interface is 
chosen by default.

Gopher can parse arpwatch or arpwitch dat files
to better analyze the network for free addresses,
and will try to find them automatically. It is also
p

Gopher can also be configured as a REST api server.
The Gopher OpenAPI spec can be found here:
$$GOPHER_API_SPEC$$

In daemonized/server mode, Gopher automatically 
parses arpwatch or arpwitch data files. In addition,
if no arpwitch is configured, Gopher automatically 
installs and configures it.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//	err := getHostFacts()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	VERSION = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(getHostFacts)
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&AwFile, "awfile", "a", "/var/arpwatch/arp.dat", "arpwatch data file")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&NetwInterface, "interface", "i", "", "interface to scan on (default eno1 or eth0)")
	rootCmd.PersistentFlags().StringVarP(&ArpFile, "arpfile", "r", "/proc/net/arp", "arp data file")
	rootCmd.PersistentFlags().IPNetVarP(&Subnet, "subnet", "s", Subnet, "subnet to scan for free IPs (default interface's first ipv4 subnet)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func isIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

// for some reason go doesn't have a builtin "in" operator.
// should make a pull request...
func containsIPv4(a []net.Addr) bool {
	for _, n := range a {
		if isIPv4(n.String()) {
			return true
		}
	}
	return false
}

// returns first valid interface with an ipv4 address from default_ifaces
func getHostFirstInterface(s []string) (*net.Interface, error) {
	for _, f := range s {
		ipf, erri := net.InterfaceByName(f)
		ipa, erra := ipf.Addrs()
		if erri == nil && erra == nil && len(ipa) != 0 && containsIPv4(ipa) {
			return ipf, nil
		}
	}
	err := errors.New("Error: no valid interfaces found")
	return nil, err
}

// returns first subnet from the pased interface
func getInterfaceFirstNet(i *net.Interface) (net.IPNet, error) {
	a, err := i.Addrs()
	if err != nil {
		return net.IPNet{}, err
	}
	for _, b := range a {
		_, ipnet, err := net.ParseCIDR(b.String())
		if err != nil {
			return net.IPNet{}, err
		}
		return *ipnet, nil
	}
	return net.IPNet{}, err
}

func getHostFacts() error {
	default_ifaces := []string{"eno1", "eth0", "wlp2s0", "enp61s0u1u3u3", "enp5s0f0"}
	if (NetwInterface == "") && (Subnet.String() == "<nil>") {
		def, err := getHostFirstInterface(default_ifaces)
		if err != nil {
			return err
		}
		NetwInterface = def.Name
		Subnet, err = getInterfaceFirstNet(def)
		if err != nil {
			return err
		}
		return nil
	} else if !(NetwInterface == "") && (Subnet.String() == "<nil>") {
		def, err := net.InterfaceByName(NetwInterface)
		if err != nil {
			return err
		}
		Subnet, err = getInterfaceFirstNet(def)
		if err != nil {
			return err
		}
		return nil
	} else if (NetwInterface == "") && !(Subnet.String() == "<nil>") {
		def, err := getHostFirstInterface(default_ifaces)
		if err != nil {
			return err
		}
		NetwInterface = def.Name
		return nil
	} else if !(NetwInterface == "") && !(Subnet.String() == "<nil>") {
		return nil
	} else {
		err := errors.New("Error: This should never happen")
		return err
	}
}

// initConfig reads in config file and ENV variables if set.
//func initConfig() {
//	if cfgFile != "" {
//		// Use config file from the flag.
//		viper.SetConfigFile(cfgFile)
//	} else {
//		// Find home directory.
//		home, err := homedir.Dir()
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//
//		// Search config in home directory with name ".gopher" (without extension).
//		viper.AddConfigPath(home)
//		viper.SetConfigName(".gopher")
//	}
//
//	viper.AutomaticEnv() // read in environment variables that match
//
//	// If a config file is found, read it in.
//	if err := viper.ReadInConfig(); err == nil {
//		fmt.Println("Using config file:", viper.ConfigFileUsed())
//	}
//}
