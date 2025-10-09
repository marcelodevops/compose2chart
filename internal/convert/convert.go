package convert

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)


type PortMapping struct {
	HostIP        string
	HostPort      int
	ContainerPort int
	Protocol      string // tcp or udp
}

type Options struct {
	ComposeFile string
	OutDir      string
	ChartName   string
	AppVersion  string
	Version     string
}

// parsePortString handles single ports, ip:host:container, ranges, and protocols like 8080:80/tcp
func parsePortString(s string) ([]PortMapping, error) {
	protocol := "tcp"
	if strings.Contains(s, "/") {
		parts := strings.Split(s, "/")
		s = parts[0]
		protocol = strings.ToLower(parts[1])
		if protocol != "tcp" && protocol != "udp" {
			return nil, fmt.Errorf("unsupported protocol: %s", protocol)
		}
	}

	parts := strings.Split(s, ":")
	switch len(parts) {
	case 1:
		c, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid port: %s", s)
		}
		return []PortMapping{{ContainerPort: c, Protocol: protocol}}, nil
	case 2:
		host := parts[0]
		cont := parts[1]
		if strings.Contains(host, "-") || strings.Contains(cont, "-") {
			hRange := strings.Split(host, "-")
			cRange := strings.Split(cont, "-")
			if len(hRange) != 2 || len(cRange) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", s)
			}
			hStart, _ := strconv.Atoi(hRange[0])
			hEnd, _ := strconv.Atoi(hRange[1])
			cStart, _ := strconv.Atoi(cRange[0])
			cEnd, _ := strconv.Atoi(cRange[1])
			if (hEnd-hStart) != (cEnd-cStart) {
				return nil, fmt.Errorf("mismatched port ranges: %s", s)
			}
			ports := []PortMapping{}
			for i := 0; i <= (hEnd - hStart); i++ {
				ports = append(ports, PortMapping{HostPort: hStart + i, ContainerPort: cStart + i, Protocol: protocol})
			}
			return ports, nil
		}
		h, err1 := strconv.Atoi(host)
		c, err2 := strconv.Atoi(cont)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid port mapping: %s", s)
		}
		return []PortMapping{{HostPort: h, ContainerPort: c, Protocol: protocol}}, nil
	case 3:
		pm := PortMapping{HostIP: parts[0], Protocol: protocol}
		h, err1 := strconv.Atoi(parts[1])
		c, err2 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid ip:host:container mapping: %s", s)
		}
		pm.HostPort = h
		pm.ContainerPort = c
		return []PortMapping{pm}, nil
	}
	return nil, fmt.Errorf("unsupported port format: %s", s)
}

// GenerateChart reads a docker-compose file and writes a Helm chart to OutDir
func GenerateChart(opts Options) error {
	data, err := os.ReadFile(opts.ComposeFile)
	if err != nil {
		return fmt.Errorf("failed to read compose file: %w", err)
	}

	var compose map[string]interface{}
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return fmt.Errorf("failed to parse compose file: %w", err)
	}

	services, ok := compose["services"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("compose file missing services")
	}

	if err := os.MkdirAll(filepath.Join(opts.OutDir, "templates"), 0755); err != nil {
		return err
	}

	for rawName, svc := range services {
		name := sanitizeName(rawName)
		svcDef, _ := svc.(map[string]interface{})

		image := ""
		if v, ok := svcDef["image"].(string); ok {
			image = v
		}

		ports := []PortMapping{}
		if rawPorts, ok := svcDef["ports"].([]interface{}); ok {
			for _, p := range rawPorts {
				if ps, ok := p.(string); ok {
					if pms, err := parsePortString(ps); err == nil {
						ports = append(ports, pms...)
					} else {
						fmt.Fprintf(os.Stderr, "warning: %v\n", err)
					}
				}
			}
		}

		// svcValues := map[string]interface{}{
		// 	"name":   name,
		// 	"image":  image,
		// 	"ports":  ports,
		// }

		// Render Deployment YAML
		deployment := fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
spec:
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
      - name: %s
        image: %s
        ports:
%s
`,
			name, name, name, name, image, renderContainerPorts(ports))

		if err := os.WriteFile(filepath.Join(opts.OutDir, "templates", name+"-deployment.yaml"), []byte(deployment), 0644); err != nil {
			return err
		}

		// Render Service YAML if ports exist
		if len(ports) > 0 {
			service := fmt.Sprintf(`apiVersion: v1
kind: Service
metadata:
  name: %s
spec:
  selector:
    app: %s
  ports:
%s
`,
				name, name, renderServicePorts(ports))

			if err := os.WriteFile(filepath.Join(opts.OutDir, "templates", name+"-service.yaml"), []byte(service), 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func renderContainerPorts(ports []PortMapping) string {
	out := ""
	for _, pm := range ports {
		out += fmt.Sprintf("        - containerPort: %d\n          protocol: %s\n", pm.ContainerPort, strings.ToUpper(pm.Protocol))
	}
	return out
}

func renderServicePorts(ports []PortMapping) string {
	out := ""
	for i, pm := range ports {
		name := fmt.Sprintf("p%d", i)
		if pm.HostPort != 0 {
			out += fmt.Sprintf("  - name: %s\n    port: %d\n    targetPort: %d\n    protocol: %s\n", name, pm.HostPort, pm.ContainerPort, strings.ToUpper(pm.Protocol))
		} else {
			out += fmt.Sprintf("  - name: %s\n    port: %d\n    targetPort: %d\n    protocol: %s\n", name, pm.ContainerPort, pm.ContainerPort, strings.ToUpper(pm.Protocol))
		}
	}
	return out
}


// sanitizeName ensures Kubernetes resource names are DNS-1123 compliant
func sanitizeName(name string) string {
	// Lowercase, replace invalid chars with dash, trim leading/trailing dashes
	name = strings.ToLower(name)
	re := regexp.MustCompile(`[^a-z0-9-]`)
	name = re.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	if len(name) == 0 {
		name = "app"
	}
	return name
}

