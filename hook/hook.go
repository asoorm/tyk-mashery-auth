package hook

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
)

type Sha256 struct {
	sharedSecret       string
	allowedClockSkew   int64
	headerAuthKey      string
	headerSignatureKey string
}

func (s *Sha256) Init(sharedSecret string, allowedClockSkew int64, headerAuthKey string, headerSignatureKey string) {
	s.sharedSecret = sharedSecret
	s.allowedClockSkew = allowedClockSkew
	s.headerAuthKey = headerAuthKey
	s.headerSignatureKey = headerSignatureKey
}

func (s Sha256) ValidateSignature(obj *coprocess.Object) (*coprocess.Object, error) {

	//log.Info("ValidateSignature called")

	authHeader, ok := obj.Request.Headers[s.headerAuthKey]
	if !ok {
		log.Error("authorization header not present")

		return obj, errors.New("auth header not present")
	}

	xSignature, ok := obj.Request.Headers[s.headerSignatureKey]
	if !ok {
		log.Error("authorization signature not present")

		return obj, errors.New("authorization signature not present")
	}

	if err := s.validate(authHeader, xSignature); err != nil {
		// signature is not valid

		return obj, errors.New("signature is not valid")
	}

	return obj, nil
}

func (s Sha256) validate(tokenAttempt string, signatureAttempt string) error {

	now := time.Now().Unix()

	attempts := 0
	for i := int64(0); i <= s.allowedClockSkew; i++ {
		attempts++
		if hex.EncodeToString(s.Sha256Sum(tokenAttempt, now+i)) == signatureAttempt {
			//log.Info("attempts: ", attempts)
			return nil
		}

		if i == int64(0) {
			continue
		}

		attempts++
		if hex.EncodeToString(s.Sha256Sum(tokenAttempt, now-i)) == signatureAttempt {
			//log.Info("attempts: ", attempts)
			return nil
		}
	}

	//log.Info("attempts: ", attempts)
	return errors.New("invalid signature" + signatureAttempt)
}

func (s Sha256) Sha256Sum(token string, timeStamp int64) []byte {
	signature := sha256.Sum256([]byte(token + s.sharedSecret + strconv.FormatInt(timeStamp, 10)))

	return signature[:]
}
