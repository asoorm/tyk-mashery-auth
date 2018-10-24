package dispatcher

import (
	"context"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
	"github.com/asoorm/tyk-mashery-auth/hook"
)

type Server struct {
	ClockSkew          int64
	SharedSecret       string
	HeaderAuthKey      string
	HeaderSignatureKey string
}

func (s Server) Dispatch(ctx context.Context, obj *coprocess.Object) (*coprocess.Object, error) {

	switch obj.HookName {
	case "ValidateSignature":

		sharedSecret, ok := obj.Session.Metadata["secret"]
		if !ok {
			obj.Request.ReturnOverrides.ResponseCode = http.StatusUnauthorized
			obj.Request.ReturnOverrides.ResponseError = "user session does not contain shared secret meta"
			return obj, nil
		}

		validator := hook.Sha256{}
		validator.Init(sharedSecret, s.ClockSkew, s.HeaderAuthKey, s.HeaderSignatureKey)

		obj, err := validator.ValidateSignature(obj)
		if err != nil {
			obj.Request.ReturnOverrides.ResponseCode = http.StatusUnauthorized
			obj.Request.ReturnOverrides.ResponseError = err.Error()
		}
	default:
		log.Printf("hook not implemented %s", obj.HookName)
	}

	return obj, nil
}

func (s Server) DispatchEvent(ctx context.Context, obj *coprocess.Event) (*coprocess.EventReply, error) {
	log.Println("DispatchEvent called")

	return &coprocess.EventReply{}, nil
}
