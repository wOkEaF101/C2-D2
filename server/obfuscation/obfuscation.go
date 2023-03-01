package obfuscation

import (
	"C2-D2/server/models"
	"bytes"
	"crypto/rc4"
	b64 "encoding/base64"
	"encoding/gob"
	"fmt"
)

func EncryptResponse(response interface{}, key string) []byte {
	var buf bytes.Buffer
	// Convert JSON to byte array
	key = b64.StdEncoding.EncodeToString([]byte(key))
	enc := gob.NewEncoder(&buf)
	enc.Encode(response)
	// Encrypt byte array using agent token
	c, _ := rc4.NewCipher([]byte(key))
	dst := make([]byte, len(buf.Bytes()))
	c.XORKeyStream(dst, buf.Bytes())

	return dst
}

// JSON Unmarshal appears to base64 encode and manipulate the raw bytes when writing
// This function will purely handle base64 strings encrypted using the above method
// and transmitted via HTTP using JSON encoding
func DecryptResponse(encrypted string, key string) interface{} {
	var response models.Response
	// Decode from base64
	decoded, _ := b64.StdEncoding.DecodeString(encrypted)
	key = b64.StdEncoding.EncodeToString([]byte(key))
	// Create an RC4 cipher using the key
	cipher, _ := rc4.NewCipher([]byte(key))

	// Decrypt the encrypted data using the RC4 key
	decrypted := make([]byte, len(decoded))
	cipher.XORKeyStream(decrypted, decoded)

	// Decode the decrypted data and read it back into a JSON object
	reader := bytes.NewReader(decrypted)
	dec := gob.NewDecoder(reader)
	_ = dec.Decode(&response)
	fmt.Println(response)

	return response
}
