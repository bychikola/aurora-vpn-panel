package subscription

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format represents a subscription configuration format
type Format string

const (
	FormatBase64 Format = "base64"
	FormatClash  Format = "clash"
	FormatSingbox Format = "sing-box"
)

// DetectFormat определяет формат подписки по User-Agent
func DetectFormat(userAgent string) Format {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "clash") || strings.Contains(ua, "stash"):
		return FormatClash
	case strings.Contains(ua, "sing-box") || strings.Contains(ua, "sfi") || strings.Contains(ua, "sfa"):
		return FormatSingbox
	default:
		return FormatBase64
	}
}

// ProxyConfig представляет прокси-ссылку для одного протокола
type ProxyConfig struct {
	Protocol   string // vless, vmess, trojan, shadowsocks, hysteria2
	Name       string // display name
	Address    string // server host
	Port       int    // server port
	UUID       string // VLESS/VMess UUID
	Password   string // Trojan/SS password
	Flow       string // XTLS flow
	Security   string // tls, reality, xtls-vision
	SNI        string // TLS SNI
	Fingerprint string // TLS fingerprint
	Transport  string // tcp, ws, grpc, quic
	Path       string // WebSocket path
	Host       string // HTTP host
	ServiceName string // gRPC service name
	PublicKey  string // Reality public key
	ShortID    string // Reality short ID
	SpiderX    string // Reality spiderX
}

// GenerateBase64 генерирует v2ray формат (base64-encoded список ссылок)
func GenerateBase64(proxies []ProxyConfig) (string, error) {
	var links []string
	for _, p := range proxies {
		link, err := generateLink(p)
		if err != nil {
			return "", err
		}
		links = append(links, link)
	}
	payload := strings.Join(links, "\n")
	return base64.StdEncoding.EncodeToString([]byte(payload)), nil
}

// GenerateClash генерирует Clash Meta YAML конфигурацию
func GenerateClash(proxies []ProxyConfig, panelName string) (string, error) {
	type ClashProxy struct {
		Name           string `yaml:"name"`
		Type           string `yaml:"type"`
		Server         string `yaml:"server"`
		Port           int    `yaml:"port"`
		UUID           string `yaml:"uuid,omitempty"`
		Password       string `yaml:"password,omitempty"`
		Cipher         string `yaml:"cipher,omitempty"`
		Flow           string `yaml:"flow,omitempty"`
		TLS            bool   `yaml:"tls,omitempty"`
		SNI            string `yaml:"servername,omitempty"`
		Fingerprint    string `yaml:"client-fingerprint,omitempty"`
		Network        string `yaml:"network,omitempty"`
		WSOpts         *WSOpts `yaml:"ws-opts,omitempty"`
		GRPCOpts       *GRPCOpts `yaml:"grpc-opts,omitempty"`
		RealityOpts    *RealityOpts `yaml:"reality-opts,omitempty"`
		SkipCertVerify bool   `yaml:"skip-cert-verify,omitempty"`
		AlterID        int    `yaml:"alterId,omitempty"`
	}

	type WSOpts struct {
		Path    string            `yaml:"path,omitempty"`
		Headers map[string]string `yaml:"headers,omitempty"`
	}

	type GRPCOpts struct {
		ServiceName string `yaml:"grpc-service-name,omitempty"`
	}

	type RealityOpts struct {
		PublicKey string `yaml:"public-key,omitempty"`
		ShortID   string `yaml:"short-id,omitempty"`
	}

	type ProxyGroup struct {
		Name     string   `yaml:"name"`
		Type     string   `yaml:"type"`
		Proxies  []string `yaml:"proxies"`
	}

	type ClashConfig struct {
		Proxies     []ClashProxy  `yaml:"proxies"`
		ProxyGroups []ProxyGroup  `yaml:"proxy-groups"`
	}

	var clashProxies []ClashProxy
	var proxyNames []string

	for _, p := range proxies {
		cp := ClashProxy{
			Name:           p.Name,
			Server:         p.Address,
			Port:           p.Port,
			SkipCertVerify: false,
		}

		switch p.Protocol {
		case "vless":
			cp.Type = "vless"
			cp.UUID = p.UUID
			cp.Flow = p.Flow
		case "vmess":
			cp.Type = "vmess"
			cp.UUID = p.UUID
			cp.AlterID = 0
			cp.Cipher = "auto"
		case "trojan":
			cp.Type = "trojan"
			cp.Password = p.Password
		case "shadowsocks", "shadowsocks-2022":
			cp.Type = "ss"
			cp.Password = p.Password
			cp.Cipher = "2022-blake3-aes-256-gcm"
		case "hysteria2":
			cp.Type = "hysteria2"
			cp.Password = p.Password
		}

		if p.Security == "tls" || p.Security == "reality" || p.Security == "xtls-vision" {
			cp.TLS = true
		}
		cp.SNI = p.SNI
		cp.Fingerprint = p.Fingerprint

		if p.Transport == "ws" {
			cp.Network = "ws"
			cp.WSOpts = &WSOpts{
				Path: p.Path,
			}
			if p.Host != "" {
				cp.WSOpts.Headers = map[string]string{"Host": p.Host}
			}
		} else if p.Transport == "grpc" {
			cp.Network = "grpc"
			cp.GRPCOpts = &GRPCOpts{ServiceName: p.ServiceName}
		} else {
			cp.Network = p.Transport
		}

		if p.Security == "reality" && p.PublicKey != "" {
			cp.RealityOpts = &RealityOpts{
				PublicKey: p.PublicKey,
				ShortID:   p.ShortID,
			}
		}

		clashProxies = append(clashProxies, cp)
		proxyNames = append(proxyNames, p.Name)
	}

	config := ClashConfig{
		Proxies: clashProxies,
		ProxyGroups: []ProxyGroup{
			{
				Name:    fmt.Sprintf("%s Auto", panelName),
				Type:    "url-test",
				Proxies: proxyNames,
			},
		},
	}

	out, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// GenerateSingbox генерирует Sing-box JSON конфигурацию
func GenerateSingbox(proxies []ProxyConfig) (string, error) {
	// Simplified sing-box JSON generation
	type SingboxOutbound struct {
		Type       string `json:"type"`
		Tag        string `json:"tag"`
		Server     string `json:"server"`
		ServerPort int    `json:"server_port"`
		UUID       string `json:"uuid,omitempty"`
		Password   string `json:"password,omitempty"`
		Flow       string `json:"flow,omitempty"`
		TLS        *SingboxTLS `json:"tls,omitempty"`
		Transport  *SingboxTransport `json:"transport,omitempty"`
	}

	type SingboxTLS struct {
		Enabled    bool   `json:"enabled"`
		ServerName string `json:"server_name,omitempty"`
		UTLS       *SingboxUTLS `json:"utls,omitempty"`
		Reality    *SingboxReality `json:"reality,omitempty"`
	}

	type SingboxUTLS struct {
		Enabled     bool   `json:"enabled"`
		Fingerprint string `json:"fingerprint,omitempty"`
	}

	type SingboxReality struct {
		Enabled   bool   `json:"enabled"`
		PublicKey string `json:"public_key,omitempty"`
		ShortID   string `json:"short_id,omitempty"`
	}

	type SingboxTransport struct {
		Type        string `json:"type"`
		Path        string `json:"path,omitempty"`
		Host        string `json:"host,omitempty"`
		ServiceName string `json:"service_name,omitempty"`
	}

	type SingboxConfig struct {
		Outbounds []SingboxOutbound `json:"outbounds"`
	}

	var outbounds []SingboxOutbound
	for _, p := range proxies {
		ob := SingboxOutbound{
			Type:       p.Protocol,
			Tag:        p.Name,
			Server:     p.Address,
			ServerPort: p.Port,
			UUID:       p.UUID,
			Password:   p.Password,
			Flow:       p.Flow,
		}

		if p.Security != "none" {
			ob.TLS = &SingboxTLS{
				Enabled:    true,
				ServerName: p.SNI,
			}
			if p.Fingerprint != "" {
				ob.TLS.UTLS = &SingboxUTLS{
					Enabled:     true,
					Fingerprint: p.Fingerprint,
				}
			}
			if p.Security == "reality" {
				ob.TLS.Reality = &SingboxReality{
					Enabled:   true,
					PublicKey: p.PublicKey,
					ShortID:   p.ShortID,
				}
			}
		}

		if p.Transport != "tcp" {
			ob.Transport = &SingboxTransport{
				Type:        p.Transport,
				Path:        p.Path,
				Host:        p.Host,
				ServiceName: p.ServiceName,
			}
		}

		outbounds = append(outbounds, ob)
	}

	config := SingboxConfig{Outbounds: outbounds}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ─── Link generator (internal) ───

func generateLink(p ProxyConfig) (string, error) {
	switch p.Protocol {
	case "vless":
		return generateVLESSLink(p), nil
	case "vmess":
		return generateVMessLink(p)
	case "trojan":
		return generateTrojanLink(p), nil
	case "shadowsocks", "shadowsocks-2022":
		return generateSSLink(p), nil
	case "hysteria2":
		return generateHysteria2Link(p), nil
	default:
		return "", fmt.Errorf("unsupported protocol: %s", p.Protocol)
	}
}

func generateVLESSLink(p ProxyConfig) string {
	// vless://uuid@host:port?params#name
	params := fmt.Sprintf("encryption=none&security=%s&sni=%s&fp=%s&type=%s",
		p.Security, p.SNI, p.Fingerprint, p.Transport)
	if p.Flow != "" {
		params += "&flow=" + p.Flow
	}
	if p.PublicKey != "" {
		params += "&pbk=" + p.PublicKey + "&sid=" + p.ShortID
	}
	if p.Path != "" {
		params += "&path=" + p.Path
	}
	if p.Host != "" {
		params += "&host=" + p.Host
	}
	return fmt.Sprintf("vless://%s@%s:%d?%s#%s", p.UUID, p.Address, p.Port, params, urlEncode(p.Name))
}

func generateVMessLink(p ProxyConfig) (string, error) {
	// vmess://base64(json)
	config := fmt.Sprintf(
		`{"v":"2","ps":"%s","add":"%s","port":"%d","id":"%s","aid":"0","net":"%s","type":"none","host":"%s","path":"%s","tls":"%s","sni":"%s"}`,
		p.Name, p.Address, p.Port, p.UUID, p.Transport, p.Host, p.Path, p.Security, p.SNI,
	)
	return "vmess://" + base64.StdEncoding.EncodeToString([]byte(config)), nil
}

func generateTrojanLink(p ProxyConfig) string {
	params := fmt.Sprintf("security=%s&sni=%s&type=%s", p.Security, p.SNI, p.Transport)
	if p.Path != "" {
		params += "&path=" + p.Path
	}
	return fmt.Sprintf("trojan://%s@%s:%d?%s#%s", p.Password, p.Address, p.Port, params, urlEncode(p.Name))
}

func generateSSLink(p ProxyConfig) string {
	// ss://base64(method:password)@host:port#name
	method := "2022-blake3-aes-256-gcm"
	userinfo := base64.StdEncoding.EncodeToString([]byte(method + ":" + p.Password))
	return fmt.Sprintf("ss://%s@%s:%d#%s", userinfo, p.Address, p.Port, urlEncode(p.Name))
}

func generateHysteria2Link(p ProxyConfig) string {
	params := fmt.Sprintf("sni=%s&insecure=0", p.SNI)
	return fmt.Sprintf("hysteria2://%s@%s:%d?%s#%s", p.Password, p.Address, p.Port, params, urlEncode(p.Name))
}

func urlEncode(s string) string {
	return strings.ReplaceAll(s, " ", "%20")
}
