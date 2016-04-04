// Shared code

package lib

import (
        "log"

        "github.com/nats-io/nats"
)


type NatsConnection struct {
        Nc         *nats.Conn
        Ec         *nats.EncodedConn
}

func NatsConn (urls string) (*NatsConnection) {
        // establish connection to NATS server
        nc, err := nats.Connect(urls)
        if err != nil {
                log.Fatalf("cannot reach NATS server, service will abort: ", err)
        }       
        ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	ret := &NatsConnection{Nc: nc, Ec: ec}
	return ret
}
