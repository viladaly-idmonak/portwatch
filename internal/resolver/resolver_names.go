package resolver

// wellKnown is a curated map of common port numbers to service names.
var wellKnown = map[uint16]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	465:  "smtps",
	587:  "submission",
	993:  "imaps",
	995:  "pop3s",
	3306: "mysql",
	5432: "postgresql",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	9200: "elasticsearch",
	27017: "mongodb",
}

// ServiceName returns a human-readable name for port, or empty string if unknown.
func ServiceName(port uint16) string {
	if name, ok := wellKnown[port]; ok {
		return name
	}
	return ""
}
