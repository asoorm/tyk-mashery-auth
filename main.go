package main

import (
	"encoding/hex"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
	"github.com/asoorm/tyk-mashery-auth/dispatcher"
	"github.com/asoorm/tyk-mashery-auth/hook"
	"google.golang.org/grpc"
)

const (
	defaultListenAddress      = ":9000"
	defaultNetwork            = "tcp"
	defaultClockSkew          = 300
	defaultSharedSecret       = "4321knj8fqgm5ffq64tdzifato6fb5p5rkqze933ehivqelctivti8qs0xnzmpq3"
	defaultDebug              = false
	defaultDebugToken         = "foo"
	defaultHeaderAuthKey      = "Api-Key"
	defaultHeaderSignatureKey = "X-Signature"
)

func main() {

	allowedClockSkew := flag.Int64("skew", defaultClockSkew, "allowed clock skew in seconds")
	network := flag.String("network", defaultNetwork, "network mode e.g. tcp | unix")
	listenAddress := flag.String("listen", defaultListenAddress, "listen address e.g. :9000 | /tmp/foo.sock")
	debug := flag.Bool("debug", defaultDebug, "enable debug mode")
	debugToken := flag.String("token", defaultDebugToken, "token used for generating debug logs")
	sharedSecret := flag.String("secret", defaultSharedSecret, "shared secret for debugging")
	headerAuthKey := flag.String("header_auth", defaultHeaderAuthKey, "header location to look for auth token")
	headerSignatureKey := flag.String("header_signature", defaultHeaderSignatureKey, "header location to look for signature")

	flag.Parse()

	/*
	 * Only run if debug mode is set.
	 * Responsible for generating test curl messages to use when sending requests via the gateway to invoke signature
	 * validating gRPC middleware.
	 */
	if *debug {
		log.SetLevel(log.DebugLevel)

		myShaGenerator := hook.Sha256{}
		myShaGenerator.Init(*sharedSecret, *allowedClockSkew, *headerAuthKey, *headerSignatureKey)
		go func() {
			for {
				msgFormat := "curl http://localhost:8080/sha/get -H '%s: %s' -H '%s: %s'"

				now := time.Now().Unix()

				log.Debugf("raw: %s%s%d", *debugToken, *sharedSecret, now)
				log.Debugf("now: "+msgFormat, *headerAuthKey, *debugToken, *headerSignatureKey, hex.EncodeToString(myShaGenerator.Sha256Sum(*debugToken, now)))
				log.Debugf("-30s: "+msgFormat, *headerAuthKey, *debugToken, *headerSignatureKey, hex.EncodeToString(myShaGenerator.Sha256Sum(*debugToken, now-30)))
				log.Debugf("+30s: "+msgFormat, *headerAuthKey, *debugToken, *headerSignatureKey, hex.EncodeToString(myShaGenerator.Sha256Sum(*debugToken, now+30)))

				time.Sleep(time.Second * 30)
			}
		}()
	}

	/*
	 * If listening on unix domain socket, e.g. unix:///tmp/foo.sock, then by killing the process with ctrl+c
	 * we need to catch this signal & delete the sock file, otherwise we won't be able to re-bind to the socket.
	 */
	go func() {
		<-handleSIGINTKILL()

		log.Info("received termination signal")

		if *network == "unix" {
			log.Infof("unbinding from %s://%s", *network, *listenAddress)
			if err := os.Remove(*listenAddress); err != nil {
				log.WithError(err).Error("unable to unbind. Please delete sock file manually")
			}
		}
		os.Exit(0)
	}()

	listener, err := net.Listen(*network, *listenAddress)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("gRPC server listening on %s://%s", *network, *listenAddress)

	s := grpc.NewServer()
	coprocess.RegisterDispatcherServer(s, &dispatcher.Server{
		ClockSkew:          *allowedClockSkew,
		SharedSecret:       *sharedSecret,
		HeaderAuthKey:      *headerAuthKey,
		HeaderSignatureKey: *headerSignatureKey,
	})

	s.Serve(listener)
}

func handleSIGINTKILL() chan os.Signal {
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	return sig
}
