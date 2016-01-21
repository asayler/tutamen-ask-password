package main

import (
	"bytes"
	"fmt"
	"github.com/go-ini/ini"
	"net"
	"os"
	"path/filepath"
)

const (
	SYSTEMD_ASK_PASSWORD_DIR  = "/run/systemd/ask-password"
	SYSTEMD_ASK_PATTERN       = "ask.*"
)

func main() {

	err := os.Chdir(SYSTEMD_ASK_PASSWORD_DIR)
	if err != nil {
		fmt.Printf("unable to enter '%s': %s\n", SYSTEMD_ASK_PASSWORD_DIR, err);
		os.Exit(1)
	}

	matches, err := filepath.Glob(SYSTEMD_ASK_PATTERN)
	if err != nil {
		fmt.Printf("unable to search for files, glob error: %s\n", err);
		os.Exit(1)
	}

	for _,ask_file := range matches {

		fmt.Println("matched path:", ask_file);
		socket := parse_socket(ask_file)
		if socket != "" {
			fmt.Println("Socket =", socket)
			password, err := TutamenGetPassword()
			if err == nil {
				write_password(socket, password)
			} else {
				fmt.Println("Unable to get secret:", err)
			}
		}
	}
}

func parse_socket(ask_path string) (socket string) {

	cfg, err := ini.Load(ask_path)
	if err != nil {
		//fmt.Println("error parsing '", ask_path, "': ", err)
		return
	}

	section, err := cfg.GetSection("Ask")
	if err != nil {
		//fmt.Println("error getting [Ask] section from '", ask_path, "': ", err)
		return
	}

	key, err := section.GetKey("Socket")
	if err != nil {
		fmt.Println(err)
		return
	}

	socket = key.String()
	return
}

func write_password(socket_path string, password string) {

	cxn, err := net.Dial("unixgram", socket_path)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		defer cxn.Close()
	}

	var pwbuf bytes.Buffer
	pwbuf.WriteString("+")
	pwbuf.WriteString(password)
	raw := pwbuf.Bytes()

	numw, err := cxn.Write(raw)

	if err != nil {
		fmt.Println(err)

	} else if numw != len(raw) {
		fmt.Println("Only write %d of %d bytes", numw, len(raw))
	}

	return
}
