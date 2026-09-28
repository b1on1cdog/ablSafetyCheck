Android ablSafetyCheck binary PoC<br>

# Features
- Support "flash", "backup", "verify" and "restore" commands<br>
- SoC verification (to prevent accidental mismatched abl flash)<br>
- pre-flash checksum verification<br>
- backup include .sha256 files to confirm long term integrity<br>
- Stdout and Stderr goes to /sdcard/rocknix_abl/output.txt<br>
- Written by a human, no LLMs involved<br>

## Disclaimer
> [!WARNING]
> Flashing a custom ABL modifies a critical component of your device's boot process.<br>
> The authors of this project are not responsible for damage, data loss, or other issues resulting from its use.<br>

At the time of writing this, "restore" command is untested<br>

# Usage

- Copy rocknix_abl files to scripts/abl/<br>
- run: <br>
```
python3 create_scripts.py
GOARCH=arm64 GOOS=linux go build -o scripts/abl/ablScTool
```
- copy /scripts content to /sdcard/rocknix_abl/, content should look like: <br>
    - /sdcard/rocknix_abl/backup_abl.sh<br>
    - /sdcard/rocknix_abl/restore_abl.sh<br>
    - /sdcard/rocknix_abl/abl/<br>
    - /sdcard/rocknix_abl/backup/<br>
- launch the scripts from "Run script as root" from your handheld settings<br>
