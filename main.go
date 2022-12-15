package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gammazero/workerpool"
	"github.com/go-ping/ping"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const URL = "https://api.mullvad.net/www/relays/all/"

type Relay struct {
	Hostname         string  `json:"hostname"`
	CountryCode      string  `json:"country_code"`
	CountryName      string  `json:"country_name"`
	CityCode         string  `json:"city_code"`
	CityName         string  `json:"city_name"`
	Active           bool    `json:"active"`
	Owned            bool    `json:"owned"`
	Provider         string  `json:"provider"`
	Ipv4AddrIn       string  `json:"ipv4_addr_in"`
	Ipv6AddrIn       *string `json:"ipv6_addr_in"`
	NetworkPortSpeed int     `json:"network_port_speed"`
	Stboot           bool    `json:"stboot"`
	Type             string  `json:"type"`
	StatusMessages   []struct {
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"status_messages"`
	Pubkey               string  `json:"pubkey,omitempty"`
	MultihopPort         int     `json:"multihop_port,omitempty"`
	SocksName            string  `json:"socks_name,omitempty"`
	SocksPort            int     `json:"socks_port,omitempty"`
	Ipv4V2Ray            *string `json:"ipv4_v2ray,omitempty"`
	SshFingerprintSha256 string  `json:"ssh_fingerprint_sha256,omitempty"`
	SshFingerprintMd5    string  `json:"ssh_fingerprint_md5,omitempty"`

	//Added
	Ping time.Duration `json:"-"`
}

func main() {
	fmt.Println("Downloading server list...")
	res, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	relays := make([]Relay, 0)
	err = json.Unmarshal(body, &relays)
	if err != nil {
		panic(err)
	}
	countriesCount := make(map[string]int)
	for _, r := range relays {
		countriesCount[r.CountryName]++
	}
	countryNames := make([]string, 0, len(countriesCount))
	for name := range countriesCount {
		countryNames = append(countryNames, name)
	}
	sort.Strings(countryNames)

	fmt.Println("There are", len(relays), "servers")
	for i, name := range countryNames {
		fmt.Printf("%s %d \t\t", name, countriesCount[name])
		time.Sleep(50 * time.Millisecond)
		if (i+1)%3 == 0 {
			fmt.Println()
		}
	}

	var action int
	fmt.Printf("\n\nEnter 1 to ping all countries.\nEnter 2 to ping one country\nAction: ")
	_, err = fmt.Scanln(&action)
	if err != nil {
		panic(err)
	}
	if action == 0 {
		panic("Unknown action selected")
	}

	var country string
	if action == 2 {
		fmt.Printf("Country name: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		country = scanner.Text()
		if country == "" {
			panic("Input country.")
		}
	}

	if country != "" {
		filtered := make([]Relay, 0)
		for _, r := range relays {
			if strings.EqualFold(r.CountryName, country) {
				filtered = append(filtered, r)
			}
		}
		relays = filtered
	}

	relaysChan := make(chan Relay, len(relays))
	go pingAllRelays(relays, relaysChan)

	pingedRelays := make([]Relay, 0)
	for r := range relaysChan {
		pingedRelays = append(pingedRelays, r)
	}
	sort.Slice(pingedRelays, func(i, j int) bool {
		return pingedRelays[i].Ping < pingedRelays[j].Ping
	})

	fmt.Printf("\nRESULTS.....\n")
	time.Sleep(3 * time.Second)
	for i, r := range pingedRelays {
		fmt.Println(i+1, r.Hostname, r.Ipv4AddrIn, r.CityName, r.CountryName, "->", r.Ping)
		time.Sleep(20 * time.Millisecond)
	}
}

func pingAllRelays(relays []Relay, relaysChan chan Relay) {
	wp := workerpool.New(30)
	for i := range relays {
		r := relays[i]
		wp.Submit(func() {
			pingRelay(r, relaysChan)
		})
	}
	wp.StopWait()
	close(relaysChan)
}

func pingRelay(relay Relay, relayChan chan Relay) {
	fmt.Printf("Pinging %s in %s,%s..\n", relay.Ipv4AddrIn, relay.CityName, relay.CountryName)
	pinger, err := ping.NewPinger(relay.Ipv4AddrIn)
	pinger.Timeout = 3 * time.Second
	if err != nil {
		relay.Ping = 1 * time.Hour
		relayChan <- relay
		return
	}
	pinger.Count = 3
	err = pinger.Run()
	if err != nil {
		relay.Ping = 1 * time.Hour
		relayChan <- relay
		return
	}
	stats := pinger.Statistics()
	if stats.AvgRtt == 0 {
		relay.Ping = 1 * time.Hour
		relayChan <- relay
		return
	}
	relay.Ping = stats.AvgRtt
	relayChan <- relay
}
