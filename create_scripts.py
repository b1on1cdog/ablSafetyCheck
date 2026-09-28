''' populate Abl scripts '''
import os

BACKUP_FILENAME = "scripts/backup_abl.sh"
RESTORE_FILENAME = "scripts/restore_abl.sh"
BINARY_NAME = "ablScTool"
REDERR = ">> /data/local/tmp/output.txt 2>&1"

socs = ["SM4450", "SM6115", "SM8250", "SM8550", "SM8650", "SM8750"]
abl_script = [
    "#!/bin/sh",
    f"cat /sdcard/rocknix_abl/abl/{BINARY_NAME} > /data/local/tmp/{BINARY_NAME}",
    f"chmod 777 /data/local/tmp/{BINARY_NAME}",
    "CMD",
    f"rm /data/local/tmp/{BINARY_NAME}"]

os.makedirs("scripts/flash_abl", exist_ok=True)
os.makedirs("scripts/backup", exist_ok=True)
os.makedirs("scripts/abl", exist_ok=True)

for soc in socs:
    script_name = f"scripts/flash_abl/{soc}.sh"
    try:
        os.remove(script_name)
    except OSError:
        pass
    flash_abl = abl_script.copy()
    flash_abl[3] = f"/data/local/tmp/{BINARY_NAME} {soc} flash {REDERR}"
    with open(script_name, "a", encoding="utf-8") as f:
        for flash_line in flash_abl:
            f.write(flash_line + "\n")

backup_script = abl_script.copy()
restore_script = abl_script.copy()

backup_script[3] = f"/data/local/tmp/{BINARY_NAME} ANY backup {REDERR}"
restore_script[3] = f"/data/local/tmp/{BINARY_NAME} ANY restore {REDERR}"
try:
    os.remove(BACKUP_FILENAME)
except OSError:
    pass
with open(BACKUP_FILENAME, "a", encoding="utf-8") as f:
    for backup_line in backup_script:
        f.write(backup_line + "\n")
try:
    os.remove(RESTORE_FILENAME)
except OSError:
    pass
with open(RESTORE_FILENAME, "a", encoding="utf-8") as f:
    for restore_line in restore_script:
        f.write(restore_line + "\n")
