package XRIT

// VCID2Name are GOES 16/17/18 VCIDs names
var VCID2Name = map[int]string{
	0:  "Admin Text",
	1:  "Mesoscale",
	2:  "GOES-ABI", // Band 2
	6:  "GOES15",
	7:  "GOES-ABI", // Band 7
	8:  "GOES-ABI", // Band 8
	9:  "GOES-ABI", // Band 8
	13: "GOES-ABI", // Band 13
	14: "GOES-ABI", // Band 14
	15: "GOES-ABI", // Band 15
	17: "GOES17",
	20: "EMWIN",
	21: "EMWIN",
	22: "EMWIN",
	23: "NWS",
	24: "NHC",
	25: "GOES16-JPG",
	26: "INTL",
	30: "DCS",
	31: "DCS",
	32: "DCS",
	60: "Himawari",
	63: "IDLE",
}

/*
    0             Imagery    Admin Text Messages
    1             Imagery    Mesoscale (ch. 2, 7, 13)
    2             Imagery    Band 2 - Red
    6             Imagery    GOES-15
    7             Imagery    Band 7 - Shortwave Window
    8             Imagery    Band 8
    9             Imagery    Band 9 - Mid-Level Trop
   13             Imagery    Band 13
   14             Imagery    Band 14 - IR
   15             Imagery    Band 15
   20             EMWIN      Priority
   21             EMWIN      Graphics
   22             EMWIN      Other
   23             Imagery    NWS Products
   24             Imagery    NHC Graphics Products
   25             Imagery    GOES-R JPG Products
   26             Imagery    International Graphics Products
   30             DCS        DCS Admin
   31             DCS        DCS Data
   60             Imagery    Himawari

*/
