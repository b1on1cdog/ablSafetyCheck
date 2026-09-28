''' populate Abl scripts '''
import os

BACKUP_FILENAME = "scripts/backup_abl.sh"
RESTORE_FILENAME = "scripts/restore_abl.sh"

socs = ["SM4450", "SM6115", "SM8250", "SM8550", "SM8650", "SM8750"]
abl_script = [
    "#!/bin/sh",
    "cat /sdcard/rocknix_abl/ablSafetyCheck > /data/local/tmp/ablSafetyCheck",
    "CMD",
    "rm /data/local/tmp/ablSafetyCheck"]

for soc in socs:
    script_name = f"scripts/flash_abl_{soc}.sh"
    os.remove(script_name)
    flash_abl = abl_script
    flash_abl[2] = f"/data/local/tmp/ablSafetyCheck {soc} flash"
    with open(script_name, "a", encoding="utf-8") as f:
        for flash_line in flash_abl:
            f.write(flash_abl + "\n")

backup_script = abl_script
restore_script = abl_script

backup_script[2] = "/data/local/tmp/ablSafetyCheck ANY backup"
restore_script[2] = "/data/local/tmp/ablSafetyCheck ANY restore"

os.remove(BACKUP_FILENAME)
with open(BACKUP_FILENAME, "a", encoding="utf-8"):
    for backup_line in backup_script:
        f.write(backup_line + "\n")

os.remove(RESTORE_FILENAME)
with open(RESTORE_FILENAME, "a", encoding="utf-8"):
    for restore_line in restore_script:
        f.write(restore_line + "\n")
