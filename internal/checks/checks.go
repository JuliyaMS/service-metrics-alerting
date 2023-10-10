package checks

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var metricsGauge = []string{"Alloc", "BuckHashSys", "Frees",
	"GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse",
	"HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups",
	"MCacheInuse", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC",
	"NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc", "RandomValue"}

func CheckType(value string) bool {
	types := []string{"gauge", "counter"}
	for _, tp := range types {
		if tp == value {
			return true
		}
	}
	return false
}

func CheckDigit(t string, val string) bool {
	if t == "gauge" {
		if _, err := strconv.ParseFloat(val, 64); err == nil {
			return true
		}
	}
	if t == "counter" {
		if _, err := strconv.ParseInt(val, 10, 64); err == nil {
			return true
		}
	}
	return false
}

func CheckFlagsServer() error {
	flags := os.Args[1:]
	if len(flags) > 0 {
		if len(flags) != 1 {
			return errors.New("incorrect count of command line arguments")
		} else {
			data := strings.Split(flags[0], "=")
			if data[0] != "-a" {
				return errors.New("incorrect flag's name")
			}
			if !checkFlagAddr(data[1]) {
				return errors.New("adress is not correct. Need address in a form host:port")
			}
		}

	}
	return nil
}

func CheckRunFlags(flag string) bool {
	runFlags := []string{"-a", "-p", "-r"}
	for _, f := range runFlags {
		if flag == f {
			return true
		}
	}
	return false
}

func CheckFlagAgent(flag string) error {
	data := strings.Split(flag, "=")
	if len(data) == 2 {
		if CheckRunFlags(data[0]) {
			if data[0] == "-a" {
				if !checkFlagAddr(data[1]) {
					return errors.New("adress is not correct. Need address in a form host:port")
				}
			} else {
				if _, err := strconv.Atoi(data[1]); err != nil {
					return fmt.Errorf("incorrect value for flag: %s", data[0])
				}
			}
		} else {
			return fmt.Errorf("incorrect flag: %s", data[0])
		}
	} else {
		return fmt.Errorf("value for flag: %s is empty", data[0])
	}
	return nil
}

func CheckFlagsAgent() error {
	flags := os.Args[1:]
	fmt.Println(flags)
	if len(flags) > 0 {
		if len(flags) > 3 {
			return errors.New("incorrect count of command line arguments")
		} else {
			for _, flag := range flags {
				if err := CheckFlagAgent(flag); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func checkFlagAddr(addr string) bool {
	data := strings.Split(addr, ":")
	return len(data) == 2
}
