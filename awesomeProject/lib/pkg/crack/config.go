package crack

var (
	PortNames = map[int]string{
		21:    "ftp",
		22:    "ssh",
		135:   "wmi",
		445:   "smb",
		1433:  "mssql",
		1521:  "oracle",
		3306:  "mysql",
		3389:  "rdp",
		5432:  "postgres",
		6379:  "redis",
		11211: "memcached",
		27017: "mongodb",
	}

	SupportProtocols = map[string]bool{
		"ftp":       true,
		"ssh":       true,
		"wmi":       true,
		"wmihash":   true,
		"smb":       true,
		"mssql":     true,
		"oracle":    true,
		"mysql":     true,
		"rdp":       true,
		"postgres":  true,
		"redis":     true,
		"memcached": true,
		"mongodb":   true,
	}
)

var (
	UserMap = map[string][]string{
		//"ftp": {"ftp", "admin", "www"},
		"ftp": {"ftp"},
		//"ssh":      {"root", "oracle", "admin"},
		"ssh":       {"root", "oracle", "test"},
		"wmi":       {"administrator"},
		"wmihash":   {"administrator"},
		"smb":       {"administrator"},
		"mssql":     {"sa"},
		"oracle":    {"oracle", "system"},
		"mysql":     {"root", "test"},
		"rdp":       {"administrator"},
		"postgres":  {"postgres", "admin"},
		"redis":     {""},
		"memcached": {""},
		"mongodb":   {"admin", "root"},
	}

	TemplatePass = []string{"{user}"}

	CommonPass = []string{""}
)
