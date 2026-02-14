package performancemetrics

import (
	"os"
	"syscall"
)

func getCpuTime() (float64, error) {
	_, err := os.FindProcess(os.Getpid())
	if err != nil {
		return 0, err
	}

	// Use syscall to get CPU times
	var ru syscall.Rusage
	err = syscall.Getrusage(syscall.RUSAGE_SELF, &ru)
	if err != nil {
		return 0, err
	}

	// Total CPU time in seconds
	userSec := float64(ru.Utime.Sec) + float64(ru.Utime.Usec)/1e6
	sysSec := float64(ru.Stime.Sec) + float64(ru.Stime.Usec)/1e6

	return userSec + sysSec, nil
}
