package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func checkRunFlags(flag string) bool {
	runFlags := []string{"-a", "-p", "-r"}
	for _, f := range runFlags {
		if flag == f {
			return true
		}
	}
	return false
}

func checkFlagAgent(flag string) error {
	data := strings.Split(flag, "=")
	if len(data) == 2 {
		if checkRunFlags(data[0]) {
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

func checkFlagsAgent() error {
	flags := os.Args[1:]
	fmt.Println(flags)
	if len(flags) > 0 {
		if len(flags) > 3 {
			return errors.New("incorrect count of command line arguments")
		} else {
			for _, flag := range flags {
				if err := checkFlagAgent(flag); err != nil {
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
