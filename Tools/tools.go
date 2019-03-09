package Tools

import "os"

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func IsDir(name string) bool {
	s, err := os.Stat(name)
	if err != nil {
		return false
	}

	return s.IsDir()
}

func IsFile(name string) bool {
	s, err := os.Stat(name)
	if err != nil {
		return false
	}

	return !s.IsDir()
}
