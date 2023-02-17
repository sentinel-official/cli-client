package types

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	configTemplate = strings.TrimSpace(`
{
    "api": {
        "services": [
            "StatsService"
        ],
        "tag": "api"
    },
    "inbounds": [
        {
            "listen": "127.0.0.1",
            "port": {{ .API.Port }},
            "protocol": "dokodemo-door",
            "settings": {
                "address": "127.0.0.1"
            },
            "tag": "api"
        },
        {
            "listen": "127.0.0.1",
            "port": {{ .Proxy.Port }},
            "protocol": "socks",
            "settings": {
                "ip": "127.0.0.1",
                "udp": true
            },
            "sniffing": {
                "destOverride": [
                    "http",
                    "tls"
                ],
                "enabled": true
            },
            "tag": "proxy"
        }
    ],
    "log": {
        "loglevel": "none"
    },
    "outbounds": [
        {
            "protocol": "vmess",
            "settings": {
                "vnext": [
                    {
                        "address": "{{ .VMess.Address }}",
                        "port": {{ .VMess.Port }},
                        "users": [
                            {
                                "alterId": 0,
                                "id": "{{ .VMess.ID }}"
                            }
                        ]
                    }
                ]
            },
            "streamSettings": {
                "network": "{{ .VMess.Transport }}"
            },
            "tag": "vmess"
        }
    ],
    "policy": {
        "levels": {
            "0": {
                "downlinkOnly": 0,
                "uplinkOnly": 0
            }
        },
        "system": {
            "statsOutboundDownlink": true,
            "statsOutboundUplink": true
        }
    },
    "routing": {
        "rules": [
            {
                "inboundTag": [
                    "api"
                ],
                "outboundTag": "api",
                "type": "field"
            }
        ]
    },
    "stats": {},
    "transport": {
        "dsSettings": {},
        "grpcSettings": {},
        "gunSettings": {},
        "httpSettings": {},
        "kcpSettings": {},
        "quicSettings": {
            "security": "chacha20-poly1305"
        },
        "tcpSettings": {},
        "wsSettings": {}
    }
}
	`)
)

type Config struct {
	API struct {
		Port uint16 `json:"port"`
	} `json:"api"`
	Proxy struct {
		Port uint16 `json:"port"`
	} `json:"proxy"`
	VMess struct {
		Address   string `json:"address"`
		ID        string `json:"id"`
		Port      uint16 `json:"port"`
		Transport string `json:"transport"`
	} `json:"vmess"`
}

func (c *Config) WriteToFile(dir string) error {
	t, err := template.New("config_v2ray_json").Parse(configTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, c); err != nil {
		return err
	}

	cfgFilePath := filepath.Join(dir, ConfigFileName)
	return os.WriteFile(cfgFilePath, buf.Bytes(), 0600)
}
