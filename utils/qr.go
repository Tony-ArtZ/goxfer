package utils

import qrcode "github.com/skip2/go-qrcode"

func PrintQR(url string) {
	qr, _ := qrcode.New(url, qrcode.Medium)
	qr.DisableBorder = true
	qrArt := qr.ToSmallString(false)
	println(qrArt)
}
