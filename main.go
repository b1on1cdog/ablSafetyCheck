package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	outputLogFile = "/sdcard/rocknix_abl/output.txt"
	ABL_A_Backup  = "/sdcard/rocknix_abl/abl_a.img"
	ABL_B_Backup  = "/sdcard/rocknix_abl/abl_b.img"
)

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

func ternary(cond bool, str1 string, str2 string) string {
	if cond {
		return str1
	}
	return str2
}

// /system/bin/dd
// to-do: create a zip so user can just drop a single file
func ablBackup() {
	f, ferr := os.OpenFile(outputLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	ablA := exec.Command("dd", "if=/dev/block/by-name/abl_a", "of="+ABL_A_Backup, "bs=1M")
	ablB := exec.Command("dd", "if=/dev/block/by-name/abl_b", "of="+ABL_B_Backup, "bs=1M")
	if ferr == nil {
		defer f.Close()
		ablA.Stdout, ablA.Stderr, ablB.Stdout, ablB.Stderr = f, f, f, f
	}
	err := ablA.Run()
	if err != nil {
		oprint("abl_a backup error : %v\n", err)
		return
	}
	err2 := ablB.Run()
	if err2 == nil {
		oprint("abl_b backup error : %v\n", err)
		return
	}
	hashA, errA := checksumFile(ABL_A_Backup)
	hashB, errB := checksumFile(ABL_B_Backup)

	hashAS, _ := checksumFile("/dev/block/by-name/abl_a")
	hashBS, _ := checksumFile("/dev/block/by-name/abl_b")

	oprint("abl_a checksum : %v\n", ternary(hashA == hashAS, "OK", "FAIL"))
	oprint("abl_b checksum : %v\n", ternary(hashB == hashBS, "OK", "FAIL"))

	hashAF := hashA + " " + filepath.Base(ABL_A_Backup)
	hashBF := hashB + " " + filepath.Base(ABL_B_Backup)
	if errA == nil && errB == nil {
		os.WriteFile(hashAF, []byte(hashA), 0644)
		os.WriteFile(hashBF, []byte(hashB), 0644)
	}
	oprint("abl_backup: success\n")
}

// untested
// to-do: add checksum verification before flashing
func ablRestore() {
	f, ferr := os.OpenFile(outputLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	ablA := exec.Command("dd", "if="+ABL_A_Backup, "of=/dev/block/by-name/abl_a", "bs=1M")
	ablB := exec.Command("dd", "if="+ABL_B_Backup, "of=/dev/block/by-name/abl_b", "bs=1M")
	if ferr == nil {
		defer f.Close()
		ablA.Stdout, ablA.Stderr, ablB.Stdout, ablB.Stderr = f, f, f, f
	}
	ablA.Run()
	ablB.Run()
}

func checksumFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		oprint("Unable to open %v : %v\n", filepath, err)
		return "", err
	}

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		oprint("Unable to calculate file hash: %v\n", err)
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func verifyFile(filepath string) bool {
	checksumPath := filepath + ".sha256"
	checksumData, cerr := os.ReadFile(checksumPath)
	if cerr != nil {
		oprint("Unable to read %v checksum file : %v\n", checksumPath, cerr)
		return false
	}
	expectedHash := strings.Split(string(checksumData), " ")[0]

	currentHash, err := checksumFile(filepath)
	if err != nil {
		return false
	}

	if strings.EqualFold(currentHash, expectedHash) {
		oprint("%v : checksum OK\n", filepath)
		return true
	}
	oprint("%v : checksum FAIL (value: %v, expected: %v)\n", filepath, currentHash, expectedHash)
	return false
}

// untested
func ablFlash(soc string) {
	f, ferr := os.OpenFile(outputLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	ablElf := fmt.Sprintf("/sdcard/rocknix_abl/%v/abl_signed-%v.elf", soc, soc)
	if !verifyFile(ablElf) {
		oprint("Aborting operation, unable to confirm %v integrity..\n", ablElf)
		return
	}
	// dd if="/sdcard/rocknix_abl/SM8550/abl_signed-SM8550.elf" of=/dev/block/by-name/abl_a bs=1M
	// dd if="/sdcard/rocknix_abl/SM8550/abl_signed-SM8550.elf" of=/dev/block/by-name/abl_b bs=1M
	ablA := exec.Command("dd", "if="+ablElf, "of=/dev/block/by-name/abl_a", "bs=1M")
	ablB := exec.Command("dd", "if="+ablElf, "of=/dev/block/by-name/abl_b", "bs=1M")
	if ferr == nil {
		ablA.Stdout, ablA.Stderr, ablB.Stdout, ablB.Stderr = f, f, f, f
		defer f.Close()
	}
	e1 := ablA.Run()
	e2 := ablB.Run()
	if e1 != nil || e2 != nil {
		oprint("abl_flash: error\n")
		return
	}

	h1, e1 := checksumFile("/dev/block/by-name/abl_a")
	h2, e2 := checksumFile("/dev/block/by-name/abl_b")
	ablSum, _ := checksumFile(ablElf)
	if e1 != nil || e2 != nil {
		oprint("abl_flash: error\n")
		return
	}
	oprint("abl_a dst checksum : %v\n", ternary(h1 == ablSum, "OK", "FAIL"))
	oprint("abl_b dst checksum : %v\n", ternary(h2 == ablSum, "OK", "FAIL"))
	oprint("abl_flash: success\n")
}

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
	ignoreChipset := expectedChip == "ANY"

	var chip MobileChip
	for i := range chips {
		if chips[i].SocID == socID {
			chip = chips[i]
			break
		}
	}

	if chip == (MobileChip{}) {
		oprint("Unable to find MobileChip for SocID : %v, some operations are unavailable\n", socID)
		ignoreChipset = true
	}

	if chip.SocModel != expectedChip && !ignoreChipset {
		oprint("Wrong chip (running : %v, expected : %v)\n", chip.SocModel, expectedChip)
		return
	}

	if !ignoreChipset {
		oprint("Chipset verified : %v (%v)\n", chip.SocModel, chip.FriendlyName)
	}

	switch strings.ToLower(shellScript) {
	case "flash":
		if ignoreChipset {
			oprint("Chipset verification is off, aborting....\n")
			return
		}
		ablFlash(expectedChip)
		return
	case "backup":
		ablBackup()
		return
	case "restore":
		ablRestore()
		return
	case "verify":
		verifyFile(fmt.Sprintf("/sdcard/rocknix_abl/%v/abl_signed-%v.elf", chip.SocModel, chip.SocModel))
		return
	}

	f, ferr := os.OpenFile(outputLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	cmd := exec.Command(shellScript)
	if ferr == nil {
		defer f.Close()
		cmd.Stdout = f
		cmd.Stderr = f
	}
	cmd.Run()
	oprint("Operation finished\n")
}
