package main

import "fmt"
import "github.com/coreos/go-systemd/dbus"

func main() {
	err := _main()
	if err != nil {
		panic(err)
	}
}

func _main() error {

	conn, err := dbus.New()
	if err != nil {
		return err
	}

	defer conn.Close()

	muf, err := conn.MaskUnitFiles([]string{"systemd-journald-audit.socket"}, false, false)
	if err != nil {
		return err
	}

	if len(muf) == 0 {
		fmt.Printf("Unit already masked\n")
		return nil
	}

	for _, m := range muf {
		fmt.Printf("Masking successful: %v\n", m)
	}

	err = conn.Reload()
	if err != nil {
		return err
	}
	fmt.Printf("Daemon Reload successful\n")

	reschan := make(chan string)
	_, err = conn.RestartUnit("systemd-journald.service", "replace", reschan)
	if err != nil {
		return err
	}

	job := <-reschan
	if job != "done" {
		return fmt.Errorf("systemd-journald restart failed: %s", job)
	}
	fmt.Printf("systemd-journald restart successful\n")

	return nil
}
