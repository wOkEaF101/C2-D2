package main

import (
	"bytes"
	"crypto/rc4"
	b64 "encoding/base64"
	"encoding/gob"
	"fmt"
)

type Response struct {
	Type    string `json:"type"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

// DecryptResponse decrypts the given encrypted response using the given key
func DecryptResponse(encrypted string, key string) interface{} {
	var response Response
	decoded, _ := b64.StdEncoding.DecodeString(encrypted)
	// key = b64.StdEncoding.EncodeToString([]byte(key))
	// Create an RC4 cipher using the key
	cipher, _ := rc4.NewCipher([]byte(key))

	// Decrypt the encrypted data
	decrypted := make([]byte, len(decoded))
	cipher.XORKeyStream(decrypted, decoded)

	// Decode the decrypted data back into a JSON object
	reader := bytes.NewReader(decrypted)
	dec := gob.NewDecoder(reader)
	_ = dec.Decode(&response)
	fmt.Println(response)

	return response
}

var value string = "9TEhBTxzOkTAew2Xn7qNrV4T9Xx0ZR07wW8/ZoqTNdG215b00xIWPw2OnVF/ISAuv3OFbYcp9ubdZz8ZWNhXUJ9EaB/mLuXOa/0EUt3hyRQgBCaIIV09WhOp91M5iXkxS/c/RabZSRtp5P4ghwZRJTJPUKFMkVAO/tfaXcCw9y9FgymMZTpRdV5B38CSjo0Jiw=="
var uuid string = "ZDEwNzNiYTAtZGRiMS00OWQ1LWI0NDYtY2Y0YWZkMWE3YTQ2"

func main() {
	DecryptResponse(value, uuid)
}
