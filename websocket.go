package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rishi-anand/fyers-go-client/api"
	fyerswatch "github.com/rishi-anand/fyers-go-client/websocket"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	YYYYMMDD = "2006-01-02"
)

func JsonToBson(message []byte) ([]byte, error) {
	reader, err := bsonrw.NewExtJSONValueReader(bytes.NewReader(message), true)
	if err != nil {
		return []byte{}, err
	}
	buf := &bytes.Buffer{}
	writer, _ := bsonrw.NewBSONValueWriter(buf)
	err = bsonrw.Copier{}.CopyDocument(writer, reader)
	if err != nil {
		return []byte{}, err
	}
	marshaled := buf.Bytes()
	return marshaled, nil
}
func main() {
	apiKey := "85VLN1I8IV-100"
	accessToken := os.Args[1]
	// mongo database connect
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	now := time.Now().UTC()

	onConnectFunc := func() {
		fmt.Println("watch subscription is connected")
	}

	onMessageFunc := func(notification api.Notification) {
		// fmt.Println(notification.Type, notification.SymbolData)
		// fmt.Println(notification.SymbolData.LowPrice, notification.SymbolData.Symbol)

		DataPack, err := json.Marshal(notification)
		datumNode, err := JsonToBson([]byte(DataPack))
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		fmt.Println(datumNode)
		// var doc interface{}
		baseString1 := strings.ReplaceAll(notification.SymbolData.Symbol, ":", "-")
		baseString2 := strings.ReplaceAll(baseString1, ".", "")
		TickerCollection := client.Database(baseString2).Collection(now.Format(YYYYMMDD))
		result, err := TickerCollection.InsertOne(context.TODO(), datumNode)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("json inserted ", result)
		}

	}

	onErrorFunc := func(err error) {
		fmt.Errorf("failed to watch | disconnected from watch. %v", err)
	}

	onCloseFunc := func() {
		fmt.Println("watch connection is closed")
	}

	cli := fyerswatch.NewNotifier(apiKey, accessToken).
		WithOnConnectFunc(onConnectFunc).
		WithOnMessageFunc(onMessageFunc).
		WithOnErrorFunc(onErrorFunc).
		WithOnCloseFunc(onCloseFunc)

		// cli.Subscribe(api.SymbolDataTick, "NSE:SBIN-EQ", "NSE:ONGC-EQ")

	symbols := []string{"NSE:SBIN-EQ", "NSE:ONGC-EQ"}
	cli.Subscribe(api.SymbolDataTick, symbols...)

}
