package main

var (
	VersionString   string
	RevString       string
	CompilationTime string
	CompilationDate string
)

func GetVersion() string {
	if VersionString == "" {
		VersionString = "<unknown>"
	}

	return VersionString
}

func GetRevision() string {
	if RevString == "" {
		RevString = "<unknown>"
	}

	return RevString
}

func GetCompilationTime() string {
	if CompilationTime == "" {
		CompilationTime = "<unknown>"
	}

	return CompilationTime
}

func GetCompilationDate() string {
	if CompilationDate == "" {
		CompilationDate = "<unknown>"
	}

	return CompilationDate
}
