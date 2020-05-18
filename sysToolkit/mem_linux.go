package sysToolkit

import "syscall"

// GetRAMSize returns total RAM size in kB
func GetRAMSize() (int64, error) {

	// ~1kB garbage
	si := &syscall.Sysinfo_t{}

	// XXX is a raw syscall thread safe?
	e := syscall.Sysinfo(si)
	if e != nil {
		return 0, e
	}
	scale := 65536.0 // magic
	return int64(si.Totalram), nil
}
