package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Giordyfish/drand/chain"
	"github.com/Giordyfish/drand/net"
	"github.com/Giordyfish/goshimmer/client"
	"github.com/Giordyfish/goshimmer/packages/drng"
	"github.com/urfave/cli/v2"

	"github.com/iotaledger/hive.go/serializer"

	iotago "github.com/iotaledger/iota.go/v2"
)

var (
	drandClient *net.ControlClient
	api         *client.GoShimmerAPI

	dRNGInstance = uint32(1)
)

var ChrysalisAPIClient = iotago.NewNodeHTTPAPIClient("https://iota-p1.teleconsys.it", iotago.WithNodeHTTPAPIClientUserInfo(url.UserPassword("teleconsys", "tcs001")))

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

	//
	// go func() {
	// 	msgIDChrs, err := SubmitPayloadToChrysalis(context.Background(), ChrysalisAPIClient, cb.Bytes())
	// 	if err != nil {
	// 		fmt.Println("Error writing on Chrysalis Tangle: ", err)
	// 		return
	// 	}
	// 	fmt.Printf("Message written on Chrysalis Tangle, msgID %x", msgIDChrs)
	// }()
	//

	//
	go func() {
		msgIDChrs, err := SubmitPayloadToChrysalisFull(context.Background(), ChrysalisAPIClient, cb.Bytes())
		if err != nil {
			fmt.Println("Error writing on Chrysalis Tangle: ", err)
			return
		}
		fmt.Printf("Message written on Chrysalis Tangle, msgID %x", msgIDChrs)
	}()
	//

	go func() {
		msgId, err := api.BroadcastCollectiveBeacon(cb.Bytes())
		if err != nil {
			fmt.Println("Error writing on the Tangle: ", err.Error())
			return
		}
		fmt.Println("Beacon written on the Tangle with msgID: ", msgId)
	}()
}

func SubmitPayloadToChrysalis(ctx context.Context, api *iotago.NodeHTTPAPIClient, p []byte) ([]byte, error) {
	// Do not check the message because the validation would fail if
	// no parents were given. The node will first add this missing information and
	// validate the message afterwards.

	req := &iotago.RawDataEnvelope{Data: p}
	res, err := api.Do(ctx, http.MethodPost, iotago.NodeAPIRouteMessages, req, nil)
	if err != nil {
		return nil, err
	}

	messageID, err := iotago.MessageIDFromHexString(res.Header.Get("Location"))
	if err != nil {
		return nil, err
	}

	msgIDBytes := messageID[:]

	return msgIDBytes, nil
}

func SubmitPayloadToChrysalisFull(ctx context.Context, nodeHTTPAPIClient *iotago.NodeHTTPAPIClient, payload []byte) ([]byte, error) {
	// create a new node API client

	// fetch the node's info to know the min. required PoW score
	//info, err := nodeHTTPAPIClient.Info(ctx)
	//if err != nil {
	//	return nil, err
	//}

	// craft an indexation payload
	indexationPayload := &iotago.Indexation{
		Index: []byte("Teleconsys dOra"),
		Data:  payload,
	}

	//ctx, cancelFunc := context.WithTimeout(ctx, 15*time.Second)
	//defer cancelFunc()

	// build a message by fetching tips via the node API client and then do local Proof-of-Work
	msg, err := iotago.NewMessageBuilder().
		Payload(indexationPayload).
		//	Tips(ctx, nodeHTTPAPIClient).
		//	ProofOfWork(ctx, info.MinPowScore).
		Build()
	if err != nil {
		return nil, err
	}

	// submit the message to the node
	postedMsg, err := nodeHTTPAPIClient.SubmitMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	postedMsgID, _ := postedMsg.ID()
	postedMsgIDOut := postedMsgID[:]

	return postedMsgIDOut, nil
}

func SubmitPayloadToChrysalisFull2(ctx context.Context, nodeHTTPAPIClient *iotago.NodeHTTPAPIClient, payload []byte) ([]byte, error) {

	// craft an indexation payload
	indexationPayload := &iotago.Indexation{
		Index: []byte("Teleconsys dOra"),
		Data:  payload,
	}

	data, err := indexationPayload.Serialize(serializer.DeSeriModeNoValidation)
	if err != nil {
		return nil, err
	}

	req := &iotago.RawDataEnvelope{Data: data}

	res, err := nodeHTTPAPIClient.Do(ctx, http.MethodPost, iotago.NodeAPIRouteMessages, req, nil)
	if err != nil {
		return nil, err
	}

	messageID, err := iotago.MessageIDFromHexString(res.Header.Get("Location"))
	if err != nil {
		return nil, err
	}

	MsgIDOut := messageID[:]

	return MsgIDOut, nil
}
