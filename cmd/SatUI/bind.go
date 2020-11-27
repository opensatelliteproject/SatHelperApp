// +build !linux,amd64
// +build !darwin,amd64
// +build !windows,amd64

package main

func Asset(name string) ([]byte, error) {
	return []byte{}, nil
}

func AssetDir(name string) ([]string, error) {
	return []string{}, nil
}

func RestoreAssets(dir, name string) error {
	return nil
}
