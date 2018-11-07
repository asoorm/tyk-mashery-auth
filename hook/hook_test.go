package hook

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
)

func BenchmarkSha256_Sha256Sum(b *testing.B) {

	b.ReportAllocs()

	s := Sha256{}
	now := time.Now().Unix()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		_ = s.Sha256Sum("foobarbaz", now)
	}
}

func BenchmarkSha256_ValidateSignature(b *testing.B) {

	b.ReportAllocs()

	sharedSecret := "foobarbaz"
	authToken := "4321knj8fqgm5ffq64tdzifato6fb5p5rkqze933ehivqelctivti8qs0xnzmpq3"

	s := Sha256{}
	s.Init(sharedSecret, 600, "Authorization", "X-Signature")
	signatureAttempt := s.Sha256Sum(authToken, time.Now().Unix())

	coprocessObj := coprocess.Object{
		Request: &coprocess.MiniRequestObject{
			Headers: map[string]string{
				"Authorization": authToken,
				"X-Signature":   hex.EncodeToString(signatureAttempt),
			},
		},
		Session: &coprocess.SessionState{
			Metadata: map[string]string{
				"secret": sharedSecret,
			},
		},
	}

	//requestJsBytes, _ := json.MarshalIndent(coprocessObj, "", "")
	//println(string(requestJsBytes))

	for n := 0; n < b.N; n++ {
		t := Sha256{}
		t.Init(sharedSecret, 600, "Authorization", "X-Signature")
		//s.Init(sharedSecret, 600, )
		_, err := s.ValidateSignature(&coprocessObj)
		if err != nil {
			b.Log(err.Error())
			b.FailNow()
		}
	}
}
