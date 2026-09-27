package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// replace "output.txt" to "/sdcard/rocknix_abl/output.txt" in production
const outputLogFile = "output.txt"

// func oprint(soc string, format string, a ...any) {
func oprint(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	fmt.Print(s)
	f, err := os.OpenFile(outputLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		f.Write([]byte(s))
		f.Close()
	}
}

/*
func backup() {
	ablA := exec.Command("dd", "if=/dev/block/by-name/abl_a", `of="/sdcard/rocknix_abl/abl_a.img"`, "bs=1M")
	ablB := exec.Command("dd", "if=/dev/block/by-name/abl_b", `of="/sdcard/rocknix_abl/abl_b.img"`, "bs=1M")
}


func flash(){

}
*/

func main() {
	data, err := os.ReadFile("/sys/devices/soc0/soc_id")
	if err != nil || len(data) == 0 {
		oprint("Unable to read SocID!\n")
		return
	}
	socIDStr := strings.ReplaceAll(string(data), "\n", "")

	socID, err := strconv.Atoi(socIDStr)
	if err != nil {
		oprint("Unable to parse SocID : %v (%v)\n", socIDStr, len(socIDStr))
		return
	}

	chips := GetSimplifiedMobileChips()

	args := os.Args
	if len(args) < 3 {
		oprint("Missing arguments!\n")
		return
	}

	expectedChip := args[1]
	shellScript := args[2]

	var chip MobileChip
	for i := range chips {
		if chips[i].SocID == socID {
			chip = chips[i]
			break
		}
	}

	if chip.SocModel != expectedChip {
		oprint("Wrong chip (running : %v, expected : %v)\n", chip.SocModel, expectedChip)
		return
	}

	oprint("Chipset verified : %v (%v)\n", chip.SocModel, chip.FriendlyName)
	f, ferr := os.OpenFile(outputLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	cmd := exec.Command(shellScript)
	if ferr == nil {
		defer f.Close()
		cmd.Stdout = f
		cmd.Stderr = f
	}
	//cmd.Stdin = os.Stdin
	cmd.Run()
	oprint("Operation finished\n")
}
