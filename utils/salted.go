package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	mrand "math/rand"
)

func GenSaltedPasswd(password string) (saltedPassword, salt string) {
	saltBin := make([]byte, 20)
	_, err := rand.Read(saltBin)
	if err != nil {
		// generate system random data failed, use soft random instead
		mrand.Read(saltBin)
	}
	salt = hex.EncodeToString(saltBin)
	hash := hmac.New(sha256.New, []byte(salt))
	hash.Write([]byte(password))
	saltedPasswordBin := hash.Sum(nil)
	saltedPassword = hex.EncodeToString(saltedPasswordBin)
	return
}

func VerifySaltedPassword(password, saltedPassword, salt string) bool {
	hash := hmac.New(sha256.New, []byte(salt))
	hash.Write([]byte(password))
	salted := hex.EncodeToString(hash.Sum(nil))
	return salted == saltedPassword
}
