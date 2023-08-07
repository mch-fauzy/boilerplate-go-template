package oauth

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Create Access Token is method to generate unique access_token
// https://github.com/bshaffer/oauth2-server-php/blob/5a0c8000d4763b276919e2106f54eddda6bc50fa/src/OAuth2/ResponseType/AccessToken.php#L133-L163
func generateAccessToken() (string, error) {
	random := make([]string, 4)
	for i := 0; i < 4; i++ {
		b := make([]byte, 10)
		_, err := rand.Read(b)
		if err != nil {
			return "", err
		}
		random = append(random, fmt.Sprintf("%x", b))
	}

	currentTimestamp := time.Now().UnixNano() / int64(time.Millisecond)
	random = append(random, strconv.Itoa(int(currentTimestamp)))

	s := sha512.New()
	s.Write([]byte(strings.Join(random, "")))

	accessToken := fmt.Sprintf("%x", s.Sum(nil))

	return string(accessToken[0:40]), nil
}
