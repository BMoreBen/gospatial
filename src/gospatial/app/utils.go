package app

import (
	"bytes"
	"compress/flate"
	crand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"time"
)

const _letters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

/*=======================================*/
// Method: NewUUID
// Source: http://play.golang.org/p/4FkNSiUDMg
// Description:
//		Generates and returns a uuid
// @returns string
/*=======================================*/
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(crand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	// return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
	return fmt.Sprintf("%x%x%x%x%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

/*=======================================*/
// Method: NewAPIKey
// Description:
//		Generates apikey of desired length
// @param int length of apikey
// @returns string
/*=======================================*/
func NewAPIKey(n int) string {
	s := ""
	for i := 1; i <= n; i++ {
		s += string(_letters[rand.Intn(len(_letters))])
	}
	return s
}

/*=======================================*/
// Method: stringInSlice
// Description:
//		Loops through array of strings
//		Checks each string in array for match
//		If string match occurs returns true
// @param a {string} string to find
// @param list {[]string} array of strings to search
// @returns bool
/*=======================================*/
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

/*=======================================*/
// Method: sliceIndex
// Description:
//		Loops through array of strings
//		Checks each string in array for match
//		If string match occurs returns index
// @param value {string} string to find
// @param slice {[]string} array of strings to search
// @returns int
/*=======================================*/
func sliceIndex(value string, slice []string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

/*=======================================*/
// Methods: Compression
// Source: https://github.com/schollz/gofind/blob/master/utils.go#L146-L169
//         https://github.com/schollz/gofind/blob/master/fingerprint.go#L43-L54
// Description:
//		Compress and Decompress bytes
/*=======================================*/
func compressByte(src []byte) []byte {
	compressedData := new(bytes.Buffer)
	compress(src, compressedData, 9)
	return compressedData.Bytes()
}

func decompressByte(src []byte) []byte {
	compressedData := bytes.NewBuffer(src)
	deCompressedData := new(bytes.Buffer)
	decompress(compressedData, deCompressedData)
	return deCompressedData.Bytes()
}

func compress(src []byte, dest io.Writer, level int) {
	Info.Println("Compressing data")
	compressor, _ := flate.NewWriter(dest, level)
	compressor.Write(src)
	compressor.Close()
}

func decompress(src io.Reader, dest io.Writer) {
	Info.Println("Decompressing data")
	decompressor := flate.NewReader(src)
	io.Copy(dest, decompressor)
	decompressor.Close()
}