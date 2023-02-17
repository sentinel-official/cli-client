package types

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Config struct {
	Name      string    `json:"name"`
	Interface Interface `json:"-"`
	Peers     []Peer    `json:"-"`
}

type Interface struct {
	PrivateKey Key      `json:"-"`
	Addresses  []IPNet  `json:"-"`
	ListenPort uint16   `json:"-"`
	MTU        uint16   `json:"-"`
	DNS        []net.IP `json:"-"`
	DNSSearch  []string `json:"-"`
	PreUp      string   `json:"-"`
	PostUp     string   `json:"-"`
	PreDown    string   `json:"-"`
	PostDown   string   `json:"-"`
}

type Peer struct {
	PublicKey           Key      `json:"-"`
	PresharedKey        Key      `json:"-"`
	AllowedIPs          []IPNet  `json:"-"`
	Endpoint            Endpoint `json:"-"`
	PersistentKeepalive uint16   `json:"-"`
}

func (c *Config) ToWgQuick() string {
	var output strings.Builder
	output.WriteString("[Interface]\n")
	output.WriteString(fmt.Sprintf("PrivateKey = %s\n", c.Interface.PrivateKey.String()))

	if c.Interface.ListenPort > 0 {
		output.WriteString(fmt.Sprintf("ListenPort = %d\n", c.Interface.ListenPort))
	}

	if len(c.Interface.Addresses) > 0 {
		addresses := make([]string, len(c.Interface.Addresses))
		for i, address := range c.Interface.Addresses {
			addresses[i] = address.String()
		}

		output.WriteString(fmt.Sprintf("Address = %s\n", strings.Join(addresses, ", ")))
	}

	if len(c.Interface.DNS)+len(c.Interface.DNSSearch) > 0 {
		dns := make([]string, 0, len(c.Interface.DNS)+len(c.Interface.DNSSearch))
		for _, address := range c.Interface.DNS {
			dns = append(dns, address.String())
		}

		dns = append(dns, c.Interface.DNSSearch...)
		output.WriteString(fmt.Sprintf("DNS = %s\n", strings.Join(dns, ", ")))
	}

	if c.Interface.MTU > 0 {
		output.WriteString(fmt.Sprintf("MTU = %d\n", c.Interface.MTU))
	}

	if len(c.Interface.PreUp) > 0 {
		output.WriteString(fmt.Sprintf("PreUp = %s\n", c.Interface.PreUp))
	}
	if len(c.Interface.PostUp) > 0 {
		output.WriteString(fmt.Sprintf("PostUp = %s\n", c.Interface.PostUp))
	}
	if len(c.Interface.PreDown) > 0 {
		output.WriteString(fmt.Sprintf("PreDown = %s\n", c.Interface.PreDown))
	}
	if len(c.Interface.PostDown) > 0 {
		output.WriteString(fmt.Sprintf("PostDown = %s\n", c.Interface.PostDown))
	}

	for _, peer := range c.Peers {
		output.WriteString("\n[Peer]\n")
		output.WriteString(fmt.Sprintf("PublicKey = %s\n", peer.PublicKey.String()))

		if !peer.PresharedKey.IsZero() {
			output.WriteString(fmt.Sprintf("PresharedKey = %s\n", peer.PresharedKey.String()))
		}

		if len(peer.AllowedIPs) > 0 {
			allowedIPs := make([]string, len(peer.AllowedIPs))
			for i, address := range peer.AllowedIPs {
				allowedIPs[i] = address.String()
			}

			output.WriteString(fmt.Sprintf("AllowedIPs = %s\n", strings.Join(allowedIPs, ", ")))
		}

		if !peer.Endpoint.IsEmpty() {
			output.WriteString(fmt.Sprintf("Endpoint = %s\n", peer.Endpoint.String()))
		}

		if peer.PersistentKeepalive > 0 {
			output.WriteString(fmt.Sprintf("PersistentKeepalive = %d\n", peer.PersistentKeepalive))
		}
	}

	return output.String()
}

func (c *Config) WriteToFile(path string) error {
	return os.WriteFile(path, []byte(c.ToWgQuick()), 0600)
}
