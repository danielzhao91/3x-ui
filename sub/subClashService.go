package sub

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/logger"
	"github.com/mhsanaei/3x-ui/v2/util/random"
	"github.com/mhsanaei/3x-ui/v2/web/service"
	"github.com/mhsanaei/3x-ui/v2/xray"
)

//go:embed acl4ssr_full.tpl
var acl4ssrFullTpl string

// SubClashService handles Clash subscription configuration generation and management.
type SubClashService struct {
	ruleSet    string

	inboundService service.InboundService
	SubService     *SubService
}

// NewSubClashService creates a new Clash subscription service with the given configuration.
func NewSubClashService(ruleSet string, subService *SubService) *SubClashService {
	return &SubClashService{
		ruleSet:    ruleSet,
		SubService: subService,
	}
}

// GetClash generates a Clash subscription configuration for the given subscription ID and host.
func (s *SubClashService) GetClash(subId string, host string) (string, string, error) {
	inbounds, err := s.SubService.getInboundsBySubId(subId)
	if err != nil || len(inbounds) == 0 {
		return "", "", err
	}

	var header string
	var traffic xray.ClientTraffic
	var clientTraffics []xray.ClientTraffic
	var proxies []map[string]any

	// Prepare Inbounds
	for _, inbound := range inbounds {
		clients, err := s.inboundService.GetClients(inbound)
		if err != nil {
			logger.Error("SubClashService - GetClients: Unable to get clients from inbound")
		}
		if clients == nil {
			continue
		}
		if len(inbound.Listen) > 0 && inbound.Listen[0] == '@' {
			listen, port, streamSettings, err := s.SubService.getFallbackMaster(inbound.Listen, inbound.StreamSettings)
			if err == nil {
				inbound.Listen = listen
				inbound.Port = port
				inbound.StreamSettings = streamSettings
			}
		}

		for _, client := range clients {
			if client.Enable && client.SubID == subId {
				clientTraffics = append(clientTraffics, s.SubService.getClientTraffics(inbound.ClientStats, client.Email))
				newProxies := s.getProxies(inbound, client, host)
				proxies = append(proxies, newProxies...)
			}
		}
	}

	if len(proxies) == 0 {
		return "", "", nil
	}

	// Prepare statistics
	for index, clientTraffic := range clientTraffics {
		if index == 0 {
			traffic.Up = clientTraffic.Up
			traffic.Down = clientTraffic.Down
			traffic.Total = clientTraffic.Total
			if clientTraffic.ExpiryTime > 0 {
				traffic.ExpiryTime = clientTraffic.ExpiryTime
			}
		} else {
			traffic.Up += clientTraffic.Up
			traffic.Down += clientTraffic.Down
			if traffic.Total == 0 || clientTraffic.Total == 0 {
				traffic.Total = 0
			} else {
				traffic.Total += clientTraffic.Total
			}
			if clientTraffic.ExpiryTime != traffic.ExpiryTime {
				traffic.ExpiryTime = 0
			}
		}
	}

	// Format proxies section
	proxiesSection := "proxies:\n"
	for _, proxy := range proxies {
		proxiesSection += s.formatProxy(proxy) + "\n"
	}

	// Replace proxies section in template
	templateLines := strings.Split(acl4ssrFullTpl, "\n")
	var resultLines []string
	proxiesSectionFound := false
	
	for _, line := range templateLines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "proxies:" {
			proxiesSectionFound = true
			resultLines = append(resultLines, proxiesSection[:len(proxiesSection)-1]) // Remove trailing newline
		} else if proxiesSectionFound && (strings.HasPrefix(trimmedLine, "proxy-groups:") || strings.HasPrefix(trimmedLine, "rules:")) {
			proxiesSectionFound = false
			resultLines = append(resultLines, line)
		} else if !proxiesSectionFound {
			resultLines = append(resultLines, line)
		}
	}
	
	// If proxies section was not found, append it
	if !proxiesSectionFound {
		resultLines = append(resultLines, proxiesSection[:len(proxiesSection)-1]) // Remove trailing newline
	}

	header = fmt.Sprintf("upload=%d; download=%d; total=%d; expire=%d", traffic.Up, traffic.Down, traffic.Total, traffic.ExpiryTime/1000)
	return strings.Join(resultLines, "\n"), header, nil
}

// formatProxy formats a proxy map to a string in the required format
func (s *SubClashService) formatProxy(proxy map[string]any) string {
	// Start with the proxy name and opening brace
	result := "  - {name: " + proxy["name"].(string)
	
	// Add all other fields
	for key, value := range proxy {
		if key == "name" {
			continue // Already added
		}
		
		// Format the value based on its type
		var formattedValue string
		switch v := value.(type) {
		case string:
			formattedValue = v
		case int:
			formattedValue = fmt.Sprintf("%d", v)
		case bool:
			formattedValue = fmt.Sprintf("%t", v)
		default:
			formattedValue = fmt.Sprintf("%v", v)
		}
		
		// Add the field to the result
		result += ", " + key + ": " + formattedValue
	}
	
	// Close the brace
	result += "}"
	
	return result
}

func (s *SubClashService) getProxies(inbound *model.Inbound, client model.Client, host string) []map[string]any {
	var proxies []map[string]any
	stream := s.streamData(inbound.StreamSettings)

	externalProxies, ok := stream["externalProxy"].([]any)
	if !ok || len(externalProxies) == 0 {
		externalProxies = []any{
			map[string]any{
				"forceTls": "same",
				"dest":     host,
				"port":     float64(inbound.Port),
				"remark":   "",
			},
		}
	}

	delete(stream, "externalProxy")

	for _, ep := range externalProxies {
		extPrxy := ep.(map[string]any)
		inbound.Listen = extPrxy["dest"].(string)
		inbound.Port = int(extPrxy["port"].(float64))
		newStream := stream
		switch extPrxy["forceTls"].(string) {
		case "tls":
			if newStream["security"] != "tls" {
				newStream["security"] = "tls"
				newStream["tlsSettings"] = map[string]any{}
			}
		case "none":
			if newStream["security"] != "none" {
				newStream["security"] = "none"
				delete(newStream, "tlsSettings")
			}
		}

		var proxy map[string]any
		switch inbound.Protocol {
		case "vmess":
			proxy = s.genVmess(inbound, newStream, client, extPrxy["remark"].(string))
		case "vless":
			proxy = s.genVless(inbound, newStream, client, extPrxy["remark"].(string))
		case "trojan":
			proxy = s.genTrojan(inbound, newStream, client, extPrxy["remark"].(string))
		case "shadowsocks":
			proxy = s.genShadowsocks(inbound, newStream, client, extPrxy["remark"].(string))
		}

		if proxy != nil {
			proxies = append(proxies, proxy)
		}
	}

	return proxies
}

func (s *SubClashService) streamData(stream string) map[string]any {
	var streamSettings map[string]any
	json.Unmarshal([]byte(stream), &streamSettings)
	security, _ := streamSettings["security"].(string)
	switch security {
	case "tls":
		streamSettings["tlsSettings"] = s.tlsData(streamSettings["tlsSettings"].(map[string]any))
	case "reality":
		streamSettings["realitySettings"] = s.realityData(streamSettings["realitySettings"].(map[string]any))
	}
	delete(streamSettings, "sockopt")

	// remove proxy protocol
	network, _ := streamSettings["network"].(string)
	switch network {
	case "tcp":
		streamSettings["tcpSettings"] = s.removeAcceptProxy(streamSettings["tcpSettings"])
	case "ws":
		streamSettings["wsSettings"] = s.removeAcceptProxy(streamSettings["wsSettings"])
	case "httpupgrade":
		streamSettings["httpupgradeSettings"] = s.removeAcceptProxy(streamSettings["httpupgradeSettings"])
	}
	return streamSettings
}

func (s *SubClashService) removeAcceptProxy(setting any) map[string]any {
	netSettings, ok := setting.(map[string]any)
	if ok {
		delete(netSettings, "acceptProxyProtocol")
	}
	return netSettings
}

func (s *SubClashService) tlsData(tData map[string]any) map[string]any {
	tlsData := make(map[string]any, 1)
	tlsClientSettings, _ := tData["settings"].(map[string]any)

	tlsData["serverName"] = tData["serverName"]
	tlsData["alpn"] = tData["alpn"]
	if allowInsecure, ok := tlsClientSettings["allowInsecure"].(bool); ok {
		tlsData["allowInsecure"] = allowInsecure
	}
	if fingerprint, ok := tlsClientSettings["fingerprint"].(string); ok {
		tlsData["fingerprint"] = fingerprint
	}
	return tlsData
}

func (s *SubClashService) realityData(rData map[string]any) map[string]any {
	rltyData := make(map[string]any, 1)
	rltyClientSettings, _ := rData["settings"].(map[string]any)

	rltyData["show"] = false
	rltyData["publicKey"] = rltyClientSettings["publicKey"]
	rltyData["fingerprint"] = rltyClientSettings["fingerprint"]
	rltyData["mldsa65Verify"] = rltyClientSettings["mldsa65Verify"]

	// Set random data
	rltyData["spiderX"] = "/" + random.Seq(15)
	shortIds, ok := rData["shortIds"].([]any)
	if ok && len(shortIds) > 0 {
		rltyData["short-id"] = shortIds[random.Num(len(shortIds))].(string)
	} else {
		rltyData["short-id"] = ""
	}
	serverNames, ok := rData["serverNames"].([]any)
	if ok && len(serverNames) > 0 {
		rltyData["serverName"] = serverNames[random.Num(len(serverNames))].(string)
	} else {
		rltyData["serverName"] = ""
	}

	return rltyData
}

func (s *SubClashService) genVmess(inbound *model.Inbound, stream map[string]any, client model.Client, remark string) map[string]any {
	proxy := map[string]any{
		"name":        s.SubService.genRemark(inbound, client.Email, remark),
		"type":        "vmess",
		"server":      inbound.Listen,
		"port":        inbound.Port,
		"uuid":        client.ID,
		"alterId":     0,
		"cipher":      client.Security,
		"udp":         true,
		"skip-cert-verify": false,
	}

	// Network settings
	network, _ := stream["network"].(string)
	proxy["network"] = network

	switch network {
	case "tcp":
		// TCP settings
		if tcpSettings, ok := stream["tcpSettings"].(map[string]any); ok {
			if header, ok := tcpSettings["header"].(map[string]any); ok {
				if headerType, ok := header["type"].(string); ok && headerType != "" && headerType != "none" {
					proxy["headers"] = map[string]any{
						"type": headerType,
					}
				}
			}
		}
	case "ws":
		// WebSocket settings
		if wsSettings, ok := stream["wsSettings"].(map[string]any); ok {
			if path, ok := wsSettings["path"].(string); ok && path != "" {
				proxy["ws-opts"] = map[string]any{
					"path": path,
				}
			}
			if headers, ok := wsSettings["headers"].(map[string]any); ok {
				if host, ok := headers["Host"].(string); ok && host != "" {
					if wsOpts, ok := proxy["ws-opts"].(map[string]any); ok {
						wsOpts["headers"] = map[string]any{
							"Host": host,
						}
					} else {
						proxy["ws-opts"] = map[string]any{
							"headers": map[string]any{
								"Host": host,
							},
						}
					}
				}
			}
		}
	case "httpupgrade":
		// HTTPUpgrade settings
		if httpupgradeSettings, ok := stream["httpupgradeSettings"].(map[string]any); ok {
			if path, ok := httpupgradeSettings["path"].(string); ok && path != "" {
				proxy["httpupgrade-opts"] = map[string]any{
					"path": path,
				}
			}
			if host, ok := httpupgradeSettings["host"].(string); ok && host != "" {
				if httpupgradeOpts, ok := proxy["httpupgrade-opts"].(map[string]any); ok {
					httpupgradeOpts["host"] = host
				} else {
					proxy["httpupgrade-opts"] = map[string]any{
						"host": host,
					}
				}
			}
		}
	}

	// TLS settings
	if security, ok := stream["security"].(string); ok && security == "tls" {
		proxy["tls"] = true
		if tlsSettings, ok := stream["tlsSettings"].(map[string]any); ok {
			if serverName, ok := tlsSettings["serverName"].(string); ok && serverName != "" {
				proxy["servername"] = serverName
			}
			if alpn, ok := tlsSettings["alpn"].([]any); ok && len(alpn) > 0 {
				var alpnStr []string
				for _, a := range alpn {
					if str, ok := a.(string); ok {
						alpnStr = append(alpnStr, str)
					}
				}
				proxy["alpn"] = alpnStr
			}
			if tlsClientSettings, ok := tlsSettings["settings"].(map[string]any); ok {
				if allowInsecure, ok := tlsClientSettings["allowInsecure"].(bool); ok {
					proxy["skip-cert-verify"] = allowInsecure
				}
				if fingerprint, ok := tlsClientSettings["fingerprint"].(string); ok && fingerprint != "" {
					proxy["fingerprint"] = fingerprint
				}
			}
		}
	}

	return proxy
}

func (s *SubClashService) genVless(inbound *model.Inbound, stream map[string]any, client model.Client, remark string) map[string]any {
	proxy := map[string]any{
		"name":        s.SubService.genRemark(inbound, client.Email, remark),
		"type":        "vless",
		"server":      inbound.Listen,
		"port":        inbound.Port,
		"uuid":        client.ID,
		"udp":         true,
		"skip-cert-verify": false,
	}

	// Network settings
	network, _ := stream["network"].(string)
	proxy["network"] = network

	// Add encryption for VLESS from inbound settings
	var inboundSettings map[string]any
	json.Unmarshal([]byte(inbound.Settings), &inboundSettings)
	if encryption, ok := inboundSettings["encryption"].(string); ok {
		proxy["encryption"] = encryption
	}

	switch network {
	case "tcp":
		// TCP settings
		if tcpSettings, ok := stream["tcpSettings"].(map[string]any); ok {
			if header, ok := tcpSettings["header"].(map[string]any); ok {
				if headerType, ok := header["type"].(string); ok && headerType != "" && headerType != "none" {
					proxy["headers"] = map[string]any{
						"type": headerType,
					}
				}
			}
		}
		// Flow settings for TCP
		if client.Flow != "" {
			proxy["flow"] = client.Flow
		}
	case "ws":
		// WebSocket settings
		if wsSettings, ok := stream["wsSettings"].(map[string]any); ok {
			if path, ok := wsSettings["path"].(string); ok && path != "" {
				proxy["ws-opts"] = map[string]any{
					"path": path,
				}
			}
			if headers, ok := wsSettings["headers"].(map[string]any); ok {
				if host, ok := headers["Host"].(string); ok && host != "" {
					if wsOpts, ok := proxy["ws-opts"].(map[string]any); ok {
						wsOpts["headers"] = map[string]any{
							"Host": host,
						}
					} else {
						proxy["ws-opts"] = map[string]any{
							"headers": map[string]any{
								"Host": host,
							},
						}
					}
				}
			}
		}
	case "httpupgrade":
		// HTTPUpgrade settings
		if httpupgradeSettings, ok := stream["httpupgradeSettings"].(map[string]any); ok {
			if path, ok := httpupgradeSettings["path"].(string); ok && path != "" {
				proxy["httpupgrade-opts"] = map[string]any{
					"path": path,
				}
			}
			if host, ok := httpupgradeSettings["host"].(string); ok && host != "" {
				if httpupgradeOpts, ok := proxy["httpupgrade-opts"].(map[string]any); ok {
					httpupgradeOpts["host"] = host
				} else {
					proxy["httpupgrade-opts"] = map[string]any{
						"host": host,
					}
				}
			}
		}
	}

	// TLS settings
	if security, ok := stream["security"].(string); ok && security == "tls" {
		proxy["tls"] = true
		if tlsSettings, ok := stream["tlsSettings"].(map[string]any); ok {
			if serverName, ok := tlsSettings["serverName"].(string); ok && serverName != "" {
				proxy["servername"] = serverName
			}
			if alpn, ok := tlsSettings["alpn"].([]any); ok && len(alpn) > 0 {
				var alpnStr []string
				for _, a := range alpn {
					if str, ok := a.(string); ok {
						alpnStr = append(alpnStr, str)
					}
				}
				proxy["alpn"] = alpnStr
			}
			if tlsClientSettings, ok := tlsSettings["settings"].(map[string]any); ok {
				if allowInsecure, ok := tlsClientSettings["allowInsecure"].(bool); ok {
					proxy["skip-cert-verify"] = allowInsecure
				}
				if fingerprint, ok := tlsClientSettings["fingerprint"].(string); ok && fingerprint != "" {
					proxy["fingerprint"] = fingerprint
				}
			}
		}
	}

	// Reality settings
	if security, ok := stream["security"].(string); ok && security == "reality" {
		proxy["tls"] = true
		proxy["reality"] = true
		if realitySettings, ok := stream["realitySettings"].(map[string]any); ok {
			if serverName, ok := realitySettings["serverName"].(string); ok && serverName != "" {
				proxy["servername"] = serverName
			}
			if publicKey, ok := realitySettings["publicKey"].(string); ok && publicKey != "" {
				proxy["public-key"] = publicKey
			}
			if shortId, ok := realitySettings["shortId"].(string); ok && shortId != "" {
				proxy["short-id"] = shortId
			}
			if spiderX, ok := realitySettings["spiderX"].(string); ok && spiderX != "" {
				proxy["spiderX"] = spiderX
			}
		}
		if client.Flow != "" {
			proxy["flow"] = client.Flow
		}
	}

	return proxy
}

func (s *SubClashService) genTrojan(inbound *model.Inbound, stream map[string]any, client model.Client, remark string) map[string]any {
	proxy := map[string]any{
		"name":        s.SubService.genRemark(inbound, client.Email, remark),
		"type":        "trojan",
		"server":      inbound.Listen,
		"port":        inbound.Port,
		"password":    client.Password,
		"udp":         true,
		"skip-cert-verify": false,
	}

	// Network settings
	network, _ := stream["network"].(string)
	proxy["network"] = network

	switch network {
	case "tcp":
		// TCP settings
		if tcpSettings, ok := stream["tcpSettings"].(map[string]any); ok {
			if header, ok := tcpSettings["header"].(map[string]any); ok {
				if headerType, ok := header["type"].(string); ok && headerType != "" && headerType != "none" {
					proxy["headers"] = map[string]any{
						"type": headerType,
					}
				}
			}
		}
	case "ws":
		// WebSocket settings
		if wsSettings, ok := stream["wsSettings"].(map[string]any); ok {
			if path, ok := wsSettings["path"].(string); ok && path != "" {
				proxy["ws-opts"] = map[string]any{
					"path": path,
				}
			}
			if headers, ok := wsSettings["headers"].(map[string]any); ok {
				if host, ok := headers["Host"].(string); ok && host != "" {
					if wsOpts, ok := proxy["ws-opts"].(map[string]any); ok {
						wsOpts["headers"] = map[string]any{
							"Host": host,
						}
					} else {
						proxy["ws-opts"] = map[string]any{
							"headers": map[string]any{
								"Host": host,
							},
						}
					}
				}
			}
		}
	case "httpupgrade":
		// HTTPUpgrade settings
		if httpupgradeSettings, ok := stream["httpupgradeSettings"].(map[string]any); ok {
			if path, ok := httpupgradeSettings["path"].(string); ok && path != "" {
				proxy["httpupgrade-opts"] = map[string]any{
					"path": path,
				}
			}
			if host, ok := httpupgradeSettings["host"].(string); ok && host != "" {
				if httpupgradeOpts, ok := proxy["httpupgrade-opts"].(map[string]any); ok {
					httpupgradeOpts["host"] = host
				} else {
					proxy["httpupgrade-opts"] = map[string]any{
						"host": host,
					}
				}
			}
		}
	}

	// TLS settings
	if security, ok := stream["security"].(string); ok && security == "tls" {
		proxy["tls"] = true
		if tlsSettings, ok := stream["tlsSettings"].(map[string]any); ok {
			if serverName, ok := tlsSettings["serverName"].(string); ok && serverName != "" {
				proxy["servername"] = serverName
			}
			if alpn, ok := tlsSettings["alpn"].([]any); ok && len(alpn) > 0 {
				var alpnStr []string
				for _, a := range alpn {
					if str, ok := a.(string); ok {
						alpnStr = append(alpnStr, str)
					}
				}
				proxy["alpn"] = alpnStr
			}
			if tlsClientSettings, ok := tlsSettings["settings"].(map[string]any); ok {
				if allowInsecure, ok := tlsClientSettings["allowInsecure"].(bool); ok {
					proxy["skip-cert-verify"] = allowInsecure
				}
				if fingerprint, ok := tlsClientSettings["fingerprint"].(string); ok && fingerprint != "" {
					proxy["fingerprint"] = fingerprint
				}
			}
		}
	}

	// Reality settings
	if security, ok := stream["security"].(string); ok && security == "reality" {
		proxy["tls"] = true
		proxy["reality"] = true
		if realitySettings, ok := stream["realitySettings"].(map[string]any); ok {
			if serverName, ok := realitySettings["serverName"].(string); ok && serverName != "" {
				proxy["servername"] = serverName
			}
			if publicKey, ok := realitySettings["publicKey"].(string); ok && publicKey != "" {
				proxy["public-key"] = publicKey
			}
			if shortId, ok := realitySettings["shortId"].(string); ok && shortId != "" {
				proxy["short-id"] = shortId
			}
			if spiderX, ok := realitySettings["spiderX"].(string); ok && spiderX != "" {
				proxy["spiderX"] = spiderX
			}
		}
		if client.Flow != "" {
			proxy["flow"] = client.Flow
		}
	}

	return proxy
}

func (s *SubClashService) genShadowsocks(inbound *model.Inbound, stream map[string]any, client model.Client, remark string) map[string]any {
	// Get method from inbound settings
	var inboundSettings map[string]any
	json.Unmarshal([]byte(inbound.Settings), &inboundSettings)
	method, _ := inboundSettings["method"].(string)

	proxy := map[string]any{
		"name":     s.SubService.genRemark(inbound, client.Email, remark),
		"type":     "ss",
		"server":   inbound.Listen,
		"port":     inbound.Port,
		"cipher":   method,
		"password": client.Password,
		"udp":      true,
	}

	// For 2022 protocols, combine server password and client password
	if strings.HasPrefix(method, "2022") {
		if serverPassword, ok := inboundSettings["password"].(string); ok {
			proxy["password"] = fmt.Sprintf("%s:%s", serverPassword, client.Password)
		}
	}

	// Network settings
	network, _ := stream["network"].(string)
	if network != "tcp" {
		proxy["plugin"] = "v2ray-plugin"
		pluginOpts := fmt.Sprintf("mode=%s", network)
		
		switch network {
		case "ws":
			if wsSettings, ok := stream["wsSettings"].(map[string]any); ok {
				if path, ok := wsSettings["path"].(string); ok && path != "" {
					pluginOpts += fmt.Sprintf(";path=%s", path)
				}
				if headers, ok := wsSettings["headers"].(map[string]any); ok {
					if host, ok := headers["Host"].(string); ok && host != "" {
						pluginOpts += fmt.Sprintf(";host=%s", host)
					}
				}
			}
		case "httpupgrade":
			if httpupgradeSettings, ok := stream["httpupgradeSettings"].(map[string]any); ok {
				if path, ok := httpupgradeSettings["path"].(string); ok && path != "" {
					pluginOpts += fmt.Sprintf(";path=%s", path)
				}
				if host, ok := httpupgradeSettings["host"].(string); ok && host != "" {
					pluginOpts += fmt.Sprintf(";host=%s", host)
				}
			}
		}
		
		// TLS settings
		if security, ok := stream["security"].(string); ok && security == "tls" {
			pluginOpts += ";tls"
			if tlsSettings, ok := stream["tlsSettings"].(map[string]any); ok {
				if serverName, ok := tlsSettings["serverName"].(string); ok && serverName != "" {
					pluginOpts += fmt.Sprintf(";servername=%s", serverName)
				}
				if tlsClientSettings, ok := tlsSettings["settings"].(map[string]any); ok {
					if allowInsecure, ok := tlsClientSettings["allowInsecure"].(bool); ok && allowInsecure {
						pluginOpts += ";allowInsecure=true"
					}
				}
			}
		}
		
		proxy["plugin-opts"] = pluginOpts
	}

	return proxy
}
