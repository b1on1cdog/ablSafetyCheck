Android ablSafetyCheck binary PoC<br>

# Description
If specified SoC matches running device SoC, script is executed, otherwise is ignored<br>

# Compile
```
GOARCH=arm64 GOOS=linux go build -o ablSafetycheck
```
<br>

# Usage
- rename flash_abl.sh to UNSAFE_flash_abl.sh<br>
- create flash_abl.sh, with the next content:<br>
```
#!/bin/sh
/sdcard/rocknix_abl/SM8550/ablSafetyCheck SM8550 /sdcard/rocknix_abl/SM8550/UNSAFE_flash_abl.sh
```