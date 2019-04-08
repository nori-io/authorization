package database

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"os"

	"golang.org/x/crypto/scrypt"
)

var (
	User     = []byte("zenithar")
	Password = []byte("mypassw0rd")
)

const (
	// scrypt is used for strong keys
	// these are the recommended scrypt parameters
	scryptN      = 16384
	scryptR      = 8
	scryptP      = 1
	scryptKeyLen = 32
)

type Token struct {
	Username string `json:"user"`
	Password string `json:"pwd"`
	Salt     string `json:"salt"`
}

func randbytes(count int) ([]byte, error) {
	// Generate a salt
	salt := make([]byte, count)
	_, err := rand.Read(salt)
	return salt, err
}

func hmac_sha256(in, salt []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, salt)
	_, err := mac.Write(in)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

func enc_scrypt(in, salt []byte) ([]byte, error) {
	return scrypt.Key(in, salt, scryptN, scryptR, scryptP, scryptKeyLen)
}

func HashPassword(password, salt []byte) ([]byte, error) {
	peppered, _ := hmac_sha256(Password, []byte(os.Getenv("TS_PEPPER")))
	cur, _ := enc_scrypt(peppered, salt)

	return cur, nil
}

func Authenticate(password, salt, hash []byte) (bool, error) {
	h, _ := HashPassword(password, salt)

	if subtle.ConstantTimeCompare(h, hash) != 1 {
		return false, fmt.Errorf("Invalid username or password")
	}

	return true, nil
}

/*func main() {
	salt, _ := randbytes(16)
	cur, _ := HashPassword(Password, salt)

	t := Token{
		Username: string(User),
		Password: base64.StdEncoding.EncodeToString(cur),
		Salt:     base64.StdEncoding.EncodeToString(salt),
	}

	if ok, _ := Authenticate(Password, salt, cur); ok {
		fmt.Println("OK")
	}

	payload, _ := json.Marshal(t)
	fmt.Println(string(payload))
}*/