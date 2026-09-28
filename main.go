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
	ABL_A_Backup  = "/sdcard/rocknix_abl/backup/abl_a.img"
	ABL_B_Backup  = "/sdcard/rocknix_abl/backup/abl_b.img"
)

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

// to-do: create a zip so user can just drop a single file
func ablBackup() {
	oprint("abl_backup: starting...\n")
	f, ferr := os.OpenFile(outputLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	ablA := exec.Command("dd", "if=/dev/block/by-name/abl_a", "of="+ABL_A_Backup, "bs=1M")
	ablB := exec.Command("dd", "if=/dev/block/by-name/abl_b", "of="+ABL_B_Backup, "bs=1M")

	if ferr == nil {
		ablA.Stdout, ablA.Stderr, ablB.Stdout, ablB.Stderr = f, f, f, f
	}
	e1 := ablA.Run()
	e2 := ablB.Run()

	if ferr == nil {
		f.Close()
	}
	if e1 != nil || e2 != nil {
		oprint("abl_backup error : (abl_a : %v | abl_b : %v)\n", e1, e2)
		return
	}

	oprint("abl_backup: creating checksum files...\n")
	hashA, errA := checksumFile(ABL_A_Backup)
	hashB, errB := checksumFile(ABL_B_Backup)

	hashAS, _ := checksumFile("/dev/block/by-name/abl_a")
	hashBS, _ := checksumFile("/dev/block/by-name/abl_b")

	oprint("abl_a checksum : %v\n", ternary(hashA == hashAS, "OK", "FAIL"))
	oprint("abl_b checksum : %v\n", ternary(hashB == hashBS, "OK", "FAIL"))

	if errA == nil && errB == nil {
		os.WriteFile(ABL_A_Backup+".sha256", []byte(hashA+"  "+filepath.Base(ABL_A_Backup)), 0644)
		os.WriteFile(ABL_B_Backup+".sha256", []byte(hashB+"  "+filepath.Base(ABL_B_Backup)), 0644)
	}
	oprint("abl_backup: success\n")
}

// untested
func ablRestore() {
	oprint("abl_restore: running untested function\n")
	oprint("abl_restore: starting...\n")

	// user might place backups without checksum, so i'll deliberately allow the process to continue
	hashA, a1h := checksumFile(ABL_A_Backup)
	hashB, a2h := checksumFile(ABL_B_Backup)

	verifyFile(ABL_A_Backup)
	verifyFile(ABL_B_Backup)

	f, ferr := os.OpenFile(outputLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	ablA := exec.Command("dd", "if="+ABL_A_Backup, "of=/dev/block/by-name/abl_a", "bs=1M")
	ablB := exec.Command("dd", "if="+ABL_B_Backup, "of=/dev/block/by-name/abl_b", "bs=1M")
	if ferr == nil {
		ablA.Stdout, ablA.Stderr, ablB.Stdout, ablB.Stderr = f, f, f, f
	}
	e1 := ablA.Run()
	e2 := ablB.Run()
	if ferr == nil {
		f.Close()
	}

	if a1h == nil && a2h == nil {
		hashAS, _ := checksumFile("/dev/block/by-name/abl_a")
		hashBS, _ := checksumFile("/dev/block/by-name/abl_b")
		oprint("abl_restore abl_a checksum : %v (%v)\n", ternary(hashA == hashAS, "OK", "FAIL"), hashAS)
		oprint("abl_restore abl_b checksum : %v (%v)\n", ternary(hashB == hashBS, "OK", "FAIL"), hashBS)
	}
	oprint("abl_restore : %v\n", ternary(e1 == nil && e2 == nil, "success", "error"))
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

func checksumFileWithLength(filepath string, bytesToRead int64) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		oprint("Unable to open %v : %v\n", filepath, err)
		return "", err
	}

	buf := make([]byte, bytesToRead)

	_, err = file.Read(buf)
	if err != nil {
		oprint("Unable to allocate %v : %v\n", filepath, err)
		return "", err
	}
	h := sha256.New()
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil)), nil
}

func verifyFile(filepath string) bool {
	checksumPath := filepath + ".sha256"
	checksumData, cerr := os.ReadFile(checksumPath)
	if cerr != nil {
		oprint("verify_file: unable to read %v checksum file : %v\n", checksumPath, cerr)
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

func verifyAbl(soc string) bool {
	if len(soc) < 3 {
		oprint("verify_abl error : Missing SoC\n")
		return false
	}
	ablElf := fmt.Sprintf("/sdcard/rocknix_abl/abl/abl_signed-%v.elf", soc)
	felf, err := os.Open(ablElf)
	if err == nil {
		fi, err := felf.Stat()
		if err == nil {
			h1, e1 := checksumFileWithLength("/dev/block/by-name/abl_a", fi.Size())
			h2, e2 := checksumFileWithLength("/dev/block/by-name/abl_b", fi.Size())
			ablSum, _ := checksumFile(ablElf)

			if e1 != nil || e2 != nil {
				oprint("verify_abl error: (abl_a : %v | abl_b : %v )\n", e1, e2)
				return false
			}
			oprint("verify_abl abl_a checksum : %v (%v)\n", ternary(h1 == ablSum, "OK", "FAIL"), h1)
			oprint("verify_abl abl_b checksum : %v (%v)\n", ternary(h2 == ablSum, "OK", "FAIL"), h2)
			return true
		}
		oprint("verify_abl error : %v\n", err)
		return false
	}
	oprint("verify_abl error : %v\n", err)
	return false
}

func ablFlash(soc string) {
	oprint("abl_flash: starting...\n")
	if len(soc) < 3 || soc == "ANY" {
		oprint("abl_flash error : Missing SoC, aborting...\n")
		return
	}
	f, ferr := os.OpenFile(outputLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	ablElf := fmt.Sprintf("/sdcard/rocknix_abl/abl/abl_signed-%v.elf", soc)
	if !verifyFile(ablElf) {
		oprint("Aborting operation, unable to confirm %v integrity..\n", ablElf)
		return
	}
	ablA := exec.Command("dd", "if="+ablElf, "of=/dev/block/by-name/abl_a", "bs=1M")
	ablB := exec.Command("dd", "if="+ablElf, "of=/dev/block/by-name/abl_b", "bs=1M")
	if ferr == nil {
		ablA.Stdout, ablA.Stderr, ablB.Stdout, ablB.Stderr = f, f, f, f
	}
	e1 := ablA.Run()
	e2 := ablB.Run()

	if ferr == nil {
		f.Close()
	}

	if e1 != nil || e2 != nil {
		oprint("abl_flash error: (abl_a: %v | abl_b: %v)\n", e1, e2)
		return
	}
	verifyAbl(soc)

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
		// We would normally use chip.SocModel, but in order to avoid a control flow issue
		// from dropping the wrong one, i'll prioritize the parameter
		ablFlash(expectedChip)
		return
	case "backup":
		ablBackup()
		return
	case "restore":
		ablRestore()
		return
	case "verify":
		if ignoreChipset {
			oprint("Chipset verification is off, aborting....\n")
			return
		}
		verifyFile(fmt.Sprintf("/sdcard/rocknix_abl/abl/abl_signed-%v.elf", chip.SocModel))
		return
	case "verify_abl":
		if ignoreChipset {
			oprint("Chipset verification is off, aborting....\n")
			return
		}
		verifyAbl(chip.SocModel)
		return
	}

	f, ferr := os.OpenFile(outputLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	cmd := exec.Command(shellScript)
	if ferr == nil {
		defer f.Close()
		cmd.Stdout, cmd.Stderr = f, f
	}
	cmd.Run()
	oprint("Operation finished\n")
}
