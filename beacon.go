package main

import (
	"fmt"

	"github.com/Giordyfish/drand/chain"
	"github.com/Giordyfish/drand/net"
	"github.com/Giordyfish/goshimmer/client"
	"github.com/Giordyfish/goshimmer/packages/drng"
	"github.com/urfave/cli/v2"
)

var (
	drandClient *net.ControlClient
	api         *client.GoShimmerAPI

	dRNGInstance = uint32(1)
)

var goshimmerAPIurl = &cli.StringFlag{
	Name:  "goshimmerAPIurl",
	Value: "http://127.0.0.1:8080",
	Usage: "The address of the goshimmer API",
}

var instanceID = &cli.UintFlag{
	Name:  "instanceID",
	Value: 1,
	Usage: "The instanceID of the dRNG",
}

func getCoKey(client *net.ControlClient) ([]byte, error) {
	resp, err := client.ChainInfo()
	if err != nil {
		return nil, err
	}
	return resp.PublicKey, nil
}

func beaconCallback(b *chain.Beacon) {
	coKey, err := getCoKey(drandClient)
	if err != nil {
		fmt.Println("Error writing on the Tangle: ", err.Error())
		return
	}
	cb := drng.NewCollectiveBeaconPayload(
		dRNGInstance,
		b.Round,
		b.PreviousSig,
		b.Signature,
		b.Message,
		coKey)

	go func() {
		msgId, err := api.BroadcastCollectiveBeacon(cb.Bytes())
		if err != nil {
			fmt.Println("Error writing on the Tangle: ", err.Error())
			return
		}
		fmt.Println("Beacon written on the Tangle with msgID: ", msgId)
	}()
}
