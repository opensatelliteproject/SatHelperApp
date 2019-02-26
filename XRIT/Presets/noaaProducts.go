package Presets

import (
	"github.com/opensatelliteproject/SatHelperApp/XRIT/NOAAProductID"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/PacketData"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/ScannerSubProduct"
)

var NOAAProducts = map[int]PacketData.NOAAProduct{}

func init() {
	NOAAProducts[NOAAProductID.NOAA_TEXT] = PacketData.MakeNOAAProductWithName(NOAAProductID.NOAA_TEXT, "NOAA Text")

	NOAAProducts[NOAAProductID.OTHER_SATELLITES_1] = PacketData.MakeNOAAProductWithSubProductsAndName(NOAAProductID.OTHER_SATELLITES_1, "Other Satellites", map[int]PacketData.NOAASubProduct{
		0: PacketData.MakeSubProduct(0, "None"),
		1: PacketData.MakeSubProduct(1, "Infrared Full Disk"),
		2: PacketData.MakeSubProduct(3, "Visible Full Disk"),
	})

	NOAAProducts[NOAAProductID.WEATHER_DATA] = PacketData.MakeNOAAProductWithName(NOAAProductID.WEATHER_DATA, "Weather Data")

	NOAAProducts[NOAAProductID.DCS] = PacketData.MakeNOAAProductWithName(NOAAProductID.DCS, "DCS")

	NOAAProducts[NOAAProductID.HRIT_EMWIN] = PacketData.MakeNOAAProductWithName(NOAAProductID.HRIT_EMWIN, "HRIT EMWIN Text")

	NOAAProducts[NOAAProductID.ABI_RELAY] = PacketData.MakeNOAAProductWithSubProductsAndName(NOAAProductID.ABI_RELAY, "ABI RELAY", map[int]PacketData.NOAASubProduct{
		ScannerSubProduct.NONE:                         PacketData.MakeSubProduct(ScannerSubProduct.NONE, "None"),
		ScannerSubProduct.INFRARED_FULLDISK:            PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_FULLDISK, "Infrared Fulldisk"),
		ScannerSubProduct.INFRARED_NORTHERN:            PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_NORTHERN, "Infrared Northern Hemisphere"),
		ScannerSubProduct.INFRARED_SOUTHERN:            PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_SOUTHERN, "Infrared Southern Hemisphere"),
		ScannerSubProduct.INFRARED_UNITEDSTATES:        PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_UNITEDSTATES, "Infrared United States"),
		ScannerSubProduct.INFRARED_AREA_OF_INTEREST:    PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_AREA_OF_INTEREST, "Infrared Area of Interest"),
		ScannerSubProduct.VISIBLE_FULLDISK:             PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_FULLDISK, "Visible Full Disk"),
		ScannerSubProduct.VISIBLE_NORTHERN:             PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_NORTHERN, "Visible Northern Hemisphere"),
		ScannerSubProduct.VISIBLE_SOUTHERN:             PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_SOUTHERN, "Visible Southern Hemisphere"),
		ScannerSubProduct.VISIBLE_UNITEDSTATES:         PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_UNITEDSTATES, "Visible United States"),
		ScannerSubProduct.VISIBLE_AREA_OF_INTEREST:     PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_AREA_OF_INTEREST, "Visible Area of Interest"),
		ScannerSubProduct.WATERVAPOUR_FULLDISK:         PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_FULLDISK, "Water Vapour Full Disk"),
		ScannerSubProduct.WATERVAPOUR_NORTHERN:         PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_NORTHERN, "Water Vapour Northern Hemisphere"),
		ScannerSubProduct.WATERVAPOUR_SOUTHERN:         PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_SOUTHERN, "Water Vapour Southern Hemisphere"),
		ScannerSubProduct.WATERVAPOUR_UNITEDSTATES:     PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_UNITEDSTATES, "Water Vapour United States"),
		ScannerSubProduct.WATERVAPOUR_AREA_OF_INTEREST: PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_AREA_OF_INTEREST, "Water Vapour Area of Interest"),
	})

	NOAAProducts[NOAAProductID.GOES13_ABI] = PacketData.MakeNOAAProductWithSubProductsAndName(NOAAProductID.GOES13_ABI, "GOES 13 ABI", map[int]PacketData.NOAASubProduct{
		ScannerSubProduct.NONE:                         PacketData.MakeSubProduct(ScannerSubProduct.NONE, "None"),
		ScannerSubProduct.INFRARED_FULLDISK:            PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_FULLDISK, "Infrared Full Disk"),
		ScannerSubProduct.INFRARED_NORTHERN:            PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_NORTHERN, "Infrared Northern Hemisphere"),
		ScannerSubProduct.INFRARED_SOUTHERN:            PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_SOUTHERN, "Infrared Southern Hemisphere"),
		ScannerSubProduct.INFRARED_UNITEDSTATES:        PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_UNITEDSTATES, "Infrared United States"),
		ScannerSubProduct.INFRARED_AREA_OF_INTEREST:    PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_AREA_OF_INTEREST, "Infrared Area of Interest"),
		ScannerSubProduct.VISIBLE_FULLDISK:             PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_FULLDISK, "Visible Full Disk"),
		ScannerSubProduct.VISIBLE_NORTHERN:             PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_NORTHERN, "Visible Northern Hemisphere"),
		ScannerSubProduct.VISIBLE_SOUTHERN:             PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_SOUTHERN, "Visible Southern Hemisphere"),
		ScannerSubProduct.VISIBLE_UNITEDSTATES:         PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_UNITEDSTATES, "Visible United States"),
		ScannerSubProduct.VISIBLE_AREA_OF_INTEREST:     PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_AREA_OF_INTEREST, "Visible Area of Interest"),
		ScannerSubProduct.WATERVAPOUR_FULLDISK:         PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_FULLDISK, "Water Vapour Full Disk"),
		ScannerSubProduct.WATERVAPOUR_NORTHERN:         PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_NORTHERN, "Water Vapour Northern Hemisphere"),
		ScannerSubProduct.WATERVAPOUR_SOUTHERN:         PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_SOUTHERN, "Water Vapour Southern Hemisphere"),
		ScannerSubProduct.WATERVAPOUR_UNITEDSTATES:     PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_UNITEDSTATES, "Water Vapour United States"),
		ScannerSubProduct.WATERVAPOUR_AREA_OF_INTEREST: PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_AREA_OF_INTEREST, "Water Vapour Area of Interest"),
	})

	NOAAProducts[NOAAProductID.GOES15_ABI] = PacketData.MakeNOAAProductWithSubProductsAndName(NOAAProductID.GOES15_ABI, "GOES 15 ABI", map[int]PacketData.NOAASubProduct{
		ScannerSubProduct.NONE:                         PacketData.MakeSubProduct(ScannerSubProduct.NONE, "None"),
		ScannerSubProduct.INFRARED_FULLDISK:            PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_FULLDISK, "Infrared Full Disk"),
		ScannerSubProduct.INFRARED_NORTHERN:            PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_NORTHERN, "Infrared Northern Hemisphere"),
		ScannerSubProduct.INFRARED_SOUTHERN:            PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_SOUTHERN, "Infrared Southern Hemisphere"),
		ScannerSubProduct.INFRARED_UNITEDSTATES:        PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_UNITEDSTATES, "Infrared United States"),
		ScannerSubProduct.INFRARED_AREA_OF_INTEREST:    PacketData.MakeSubProduct(ScannerSubProduct.INFRARED_AREA_OF_INTEREST, "Infrared Area of Interest"),
		ScannerSubProduct.VISIBLE_FULLDISK:             PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_FULLDISK, "Visible Full Disk"),
		ScannerSubProduct.VISIBLE_NORTHERN:             PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_NORTHERN, "Visible Northern Hemisphere"),
		ScannerSubProduct.VISIBLE_SOUTHERN:             PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_SOUTHERN, "Visible Southern Hemisphere"),
		ScannerSubProduct.VISIBLE_UNITEDSTATES:         PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_UNITEDSTATES, "Visible United States"),
		ScannerSubProduct.VISIBLE_AREA_OF_INTEREST:     PacketData.MakeSubProduct(ScannerSubProduct.VISIBLE_AREA_OF_INTEREST, "Visible Area of Interest"),
		ScannerSubProduct.WATERVAPOUR_FULLDISK:         PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_FULLDISK, "Water Vapour Full Disk"),
		ScannerSubProduct.WATERVAPOUR_NORTHERN:         PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_NORTHERN, "Water Vapour Northern Hemisphere"),
		ScannerSubProduct.WATERVAPOUR_SOUTHERN:         PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_SOUTHERN, "Water Vapour Southern Hemisphere"),
		ScannerSubProduct.WATERVAPOUR_UNITEDSTATES:     PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_UNITEDSTATES, "Water Vapour United States"),
		ScannerSubProduct.WATERVAPOUR_AREA_OF_INTEREST: PacketData.MakeSubProduct(ScannerSubProduct.WATERVAPOUR_AREA_OF_INTEREST, "Water Vapour Area of Interest"),
	})

	NOAAProducts[NOAAProductID.GOES16_ABI] = PacketData.MakeNOAAProductWithSubProductsAndName(NOAAProductID.GOES16_ABI, "GOES 16 ABI", map[int]PacketData.NOAASubProduct{
		0:  PacketData.MakeSubProduct(0, "None"),
		1:  PacketData.MakeSubProduct(1, "Channel 1"),
		2:  PacketData.MakeSubProduct(2, "Channel 2"),
		3:  PacketData.MakeSubProduct(3, "Channel 3"),
		4:  PacketData.MakeSubProduct(4, "Channel 4"),
		5:  PacketData.MakeSubProduct(5, "Channel 5"),
		6:  PacketData.MakeSubProduct(6, "Channel 6"),
		7:  PacketData.MakeSubProduct(7, "Channel 7"),
		8:  PacketData.MakeSubProduct(8, "Channel 8"),
		9:  PacketData.MakeSubProduct(9, "Channel 9"),
		10: PacketData.MakeSubProduct(10, "Channel 10"),
		11: PacketData.MakeSubProduct(11, "Channel 11"),
		12: PacketData.MakeSubProduct(12, "Channel 12"),
		13: PacketData.MakeSubProduct(13, "Channel 13"),
		14: PacketData.MakeSubProduct(14, "Channel 14"),
		15: PacketData.MakeSubProduct(15, "Channel 15"),
		16: PacketData.MakeSubProduct(16, "Channel 16"),
	})

	NOAAProducts[NOAAProductID.GOES17_ABI] = PacketData.MakeNOAAProductWithSubProductsAndName(NOAAProductID.GOES17_ABI, "GOES 17 ABI", map[int]PacketData.NOAASubProduct{
		0:  PacketData.MakeSubProduct(0, "None"),
		1:  PacketData.MakeSubProduct(1, "Channel 1"),
		2:  PacketData.MakeSubProduct(2, "Channel 2"),
		3:  PacketData.MakeSubProduct(3, "Channel 3"),
		4:  PacketData.MakeSubProduct(4, "Channel 4"),
		5:  PacketData.MakeSubProduct(5, "Channel 5"),
		6:  PacketData.MakeSubProduct(6, "Channel 6"),
		7:  PacketData.MakeSubProduct(7, "Channel 7"),
		8:  PacketData.MakeSubProduct(8, "Channel 8"),
		9:  PacketData.MakeSubProduct(9, "Channel 9"),
		10: PacketData.MakeSubProduct(10, "Channel 10"),
		11: PacketData.MakeSubProduct(11, "Channel 11"),
		12: PacketData.MakeSubProduct(12, "Channel 12"),
		13: PacketData.MakeSubProduct(13, "Channel 13"),
		14: PacketData.MakeSubProduct(14, "Channel 14"),
		15: PacketData.MakeSubProduct(15, "Channel 15"),
		16: PacketData.MakeSubProduct(16, "Channel 16"),
	})

	NOAAProducts[NOAAProductID.HIMAWARI8_ABI] = PacketData.MakeNOAAProductWithSubProductsAndName(NOAAProductID.HIMAWARI8_ABI, "Himawari ABI", map[int]PacketData.NOAASubProduct{
		0:  PacketData.MakeSubProduct(0, "None"),
		1:  PacketData.MakeSubProduct(1, "Channel 1"),
		2:  PacketData.MakeSubProduct(2, "Channel 2"),
		3:  PacketData.MakeSubProduct(3, "Channel 3"),
		4:  PacketData.MakeSubProduct(4, "Channel 4"),
		5:  PacketData.MakeSubProduct(5, "Channel 5"),
		6:  PacketData.MakeSubProduct(6, "Channel 6"),
		7:  PacketData.MakeSubProduct(7, "Channel 7"),
		8:  PacketData.MakeSubProduct(8, "Channel 8"),
		9:  PacketData.MakeSubProduct(9, "Channel 9"),
		10: PacketData.MakeSubProduct(10, "Channel 10"),
		11: PacketData.MakeSubProduct(11, "Channel 11"),
		12: PacketData.MakeSubProduct(12, "Channel 12"),
		13: PacketData.MakeSubProduct(13, "Channel 13"),
		14: PacketData.MakeSubProduct(14, "Channel 14"),
		15: PacketData.MakeSubProduct(15, "Channel 15"),
		16: PacketData.MakeSubProduct(16, "Channel 16"),
	})

	NOAAProducts[NOAAProductID.EMWIN] = PacketData.MakeNOAAProductWithName(NOAAProductID.EMWIN, "EMWIN")
}

func GetProductById(productId int) PacketData.NOAAProduct {
	val, ok := NOAAProducts[productId]
	if !ok {
		return PacketData.MakeNOAAProduct(productId)
	}

	return val
}
