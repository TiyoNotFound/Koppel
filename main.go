package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type IPInfo struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Location string `json:"location"`
	ISP      string `json:"isp"`
}

type IPList struct {
	IPs []IPInfo `json:"ips"`
}

var rootCmd = &cobra.Command{
	Use:   "koppel",
	Short: "Koppel is an IP tracker CLI tool",
	Long:  "Koppel is a CLI tool to track IP addresses and save the information in a JSON file.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Koppel! Use 'koppel help' to see available commands.")
	},
}

var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Track an IP address",
	Long:  "Track an IP address and save its information to a JSON file.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: koppel track <ip_address>")
			return
		}

		ip := args[0]
		ipInfo := getIPInfo(ip)

		if ipInfo != nil {
			saveIPInfo(ipInfo)
			fmt.Println("IP information saved successfully.")
		} else {
			fmt.Println("Failed to retrieve IP information.")
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tracked IP addresses",
	Long:  "List all tracked IP addresses along with their information.",
	Run: func(cmd *cobra.Command, args []string) {
		ipList := loadIPList()
		if len(ipList.IPs) == 0 {
			fmt.Println("No tracked IP addresses.")
			return
		}

		renderTable(ipList)
	},
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear tracked IP list",
	Long:  "Clear all tracked IP addresses from the list.",
	Run: func(cmd *cobra.Command, args []string) {
		err := clearIPList()
		if err != nil {
			fmt.Println("Failed to clear tracked IP list:", err)
			return
		}
		fmt.Println("Tracked IP list cleared successfully.")
	},
}

func getIPInfo(ip string) *IPInfo {
	url := "http://ip-api.com/json/" + ip
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to retrieve IP information:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return nil
	}

	var ipInfo IPInfo
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		return nil
	}

	ipInfo.IP = ip
	ipInfo.Location = fmt.Sprintf("%s, %s, %s", ipInfo.City, ipInfo.Region, ipInfo.Country)

	return &ipInfo
}

func saveIPInfo(ipInfo *IPInfo) {
	ipList := loadIPList()
	ipList.IPs = append(ipList.IPs, *ipInfo)

	data, err := json.MarshalIndent(ipList, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal IP information:", err)
		return
	}

	err = ioutil.WriteFile("ip_info.json", data, 0644)
	if err != nil {
		fmt.Println("Failed to write IP information to file:", err)
		return
	}
}

func loadIPList() *IPList {
	data, err := ioutil.ReadFile("ip_info.json")
	if err != nil {
		return &IPList{}
	}

	var ipList IPList
	err = json.Unmarshal(data, &ipList)
	if err != nil {
		fmt.Println("Failed to unmarshal IP list:", err)
		return &IPList{}
	}

	return &ipList
}

func clearIPList() error {
	ipList := &IPList{}
	data, err := json.MarshalIndent(ipList, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("ip_info.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func renderTable(ipList *IPList) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP", "City", "Region", "Country", "ISP", "Location"})
	table.SetBorder(false)

	for _, ip := range ipList.IPs {
		table.Append([]string{ip.IP, ip.City, ip.Region, ip.Country, ip.ISP, ip.Location})
	}

	table.Render()
}

func main() {
	rootCmd.AddCommand(trackCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(clearCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
