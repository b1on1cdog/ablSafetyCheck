Android ablSafetyCheck binary PoC<br>

# Features
- Support "flash", "backup", "verify" and "restore" commands<br>
- SoC verification (to prevent accidental mismatched abl flash)<br>
- pre-restore and post-restore checksum verification<br>
- Stdout and Stderr goes to /sdcard/rocknix_abl/output.txt<br>
- Written by a human, no LLMs involved<br>

## Disclaimer
> [!WARNING]
> Flashing a custom ABL modifies a critical component of your device's boot process.
> The authors of this project are not responsible for damage, data loss, or other issues resulting from its use.

# Compile
```
GOARCH=arm64 GOOS=linux go build -o scripts/abl/ablScTool
```
<br>

# Usage
- run create_scripts.py<br>
- copy /scripts content to /sdcard/rocknix_abl/, content should look like: <br>
    - /sdcard/rocknix_abl/backup_abl.sh<br>
    - /sdcard/rocknix_abl/restore_abl.sh<br>
    - /sdcard/rocknix_abl/abl/<br>
    - /sdcard/rocknix_abl/backup/<br>
