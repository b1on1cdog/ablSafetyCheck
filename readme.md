ablSafetyCheck binary PoC<br>

# Description
If specified SoC matches running device SoC, script is executed, otherwise is ignored<br>

# Compile
```
GOARCH=arm64 GOOS=linux go build -o ablSafetycheck
```
<br>

# Usage
rename flash_abl.sh to UNSAFE_flash_abl.sh<br>
create flash_abl.sh, with the next content:<br>
```
/sdcard/rocknix_abl/SM8550/ablSafetyCheck SM8550 /sdcard/rocknix_abl/SM8550/UNSAFE_flash_abl.sh
```

# Considerations
- Using a shared ablSafetyCheck (ex: /sdcard/rocknix_abl/ablSafetyCheck) instead of a copy per SoC can save near 15mb of space<br>
- If desired, the flashing code can be easily embedded into the binary to avoid having a UNSAFE_flash_abl.sh per SoC folder<br/>
- This code was compiled and tested from ADB Shell, i didn't actually try to flash something with it<br/>
- Having UNSAFE_flash_abl.sh in a different dir than flash_abl.sh might be a good idea, so user does not accidentally bypass the protection by chosing the wrong script<br/>
- The reason why "chips.go" is more complex than it should is because is a copy-paste from another library i'm writing<br/>
