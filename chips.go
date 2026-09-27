package main

type MobileChip struct {
	SocName  string // ro.board.platform | ro.boot.product.vendor.sku
	SocID    int    // cat /sys/devices/soc0/soc_id
	SocModel string // ro.soc.model

	GPU          string
	FriendlyName string
}

// Source https://github.com/torvalds/linux/blob/fd179f8a05be3ccae366b9b96e176b51fbe54aab/include/dt-bindings/arm/qcom%2Cids.h
func GetMobileChips() []MobileChip {
	var mobileChips []MobileChip
	// Adreno 850 (8 Elite Extreme Gen 6 / SM8975)
	// TBA
	// Adreno 830
	mobileChips = append(mobileChips, MobileChip{SocID: 706, SocName: "sun", SocModel: "CQ8725S"})
	mobileChips = append(mobileChips, MobileChip{SocID: 618, SocName: "sun", SocModel: "SM8750"})  // Confirmed from Galaxy S25
	mobileChips = append(mobileChips, MobileChip{SocID: 639, SocName: "sun", SocModel: "SM8750P"}) // WiFi only
	// Adreno 750 (To-do: Add Snapdragon G3 Gen 3)
	mobileChips = append(mobileChips, MobileChip{SocID: 557, SocName: "pineapple", SocModel: "SM8650"})
	// Adreno 740
	mobileChips = append(mobileChips, MobileChip{SocID: 519, SocName: "kalama", SocModel: "SM8550"})
	mobileChips = append(mobileChips, MobileChip{SocID: 603, SocName: "kalama", SocModel: "QCS8550"}) // Confirmed from Ayn Thor
	mobileChips = append(mobileChips, MobileChip{SocID: 604, SocName: "kalama", SocModel: "QCM8550"})
	// Adreno 730
	mobileChips = append(mobileChips, MobileChip{SocID: 457, SocName: "taro", SocModel: "SM8450"}) // Confirmed in Motorola Edge Plus 2022
	// Adreno 722
	mobileChips = append(mobileChips, MobileChip{SocID: 731, SocName: "eliza", SocModel: "CQ7790M"})
	mobileChips = append(mobileChips, MobileChip{SocID: 732, SocName: "eliza", SocModel: "CQ7790S"})
	// Adreno 650
	mobileChips = append(mobileChips, MobileChip{SocID: 356, SocName: "kona", SocModel: "SM8250"})

	for i := range mobileChips {
		switch mobileChips[i].SocName {
		case "sun":
			mobileChips[i].GPU = "Adreno 830"
			mobileChips[i].FriendlyName = "Snapdragon 8 Elite"
		case "pineapple":
			mobileChips[i].GPU = "Adreno 750"
			mobileChips[i].FriendlyName = "Snapdragon 8 Gen 3"
		case "kalama":
			mobileChips[i].GPU = "Adreno 740"
			mobileChips[i].FriendlyName = "Snapdragon 8 Gen 2"
		case "taro":
			mobileChips[i].GPU = "Adreno 730"
			mobileChips[i].FriendlyName = "Snapdragon 8 Gen 1"
		case "eliza":
			mobileChips[i].GPU = "Adreno 722"
			mobileChips[i].FriendlyName = "Snapdragon 7 Gen 4"
		case "kona":
			mobileChips[i].GPU = "Adreno 650"
			mobileChips[i].FriendlyName = "Snapdragon 865"
		}
	}
	return mobileChips
}

// Return all MobileChip entries with the mobile format (ex: QCS8550 => SM8550)
func GetSimplifiedMobileChips() []MobileChip {
	mobileChips := GetMobileChips()

	for i := range mobileChips {
		switch mobileChips[i].SocName {
		case "sun":
			mobileChips[i].SocModel = "SM8750"
		case "pineapple":
			mobileChips[i].SocModel = "SM8650"
		case "kalama":
			mobileChips[i].SocModel = "SM8550"
		case "taro":
			mobileChips[i].SocModel = "SM8450"
		case "kona":
			mobileChips[i].SocModel = "SM8250"
		}
	}
	return mobileChips
}
