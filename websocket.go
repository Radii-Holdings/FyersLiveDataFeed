package main

import (
	"bytes"
	"context"
	"encoding/csv"
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
	HHMMSS   = "17:06:04 PM"
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
	now := time.Now()
	csvfile, err := os.Create("./dataset/raw/" + now.Format(YYYYMMDD) + ".csv")
	if err != nil {
		panic(err)
	} else {
		csvwriter := csv.NewWriter(csvfile)

		onConnectFunc := func() {
			fmt.Println("watch subscription is connected")
			row := []string{"Mongo Object Id", "Database Name", "collection name", "Upsert time"}
			_ = csvwriter.Write(row)
			csvwriter.Flush()

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
			// fmt.Println(datumNode)
			baseString1 := strings.ReplaceAll(notification.SymbolData.Symbol, ":", "-")
			baseString2 := strings.ReplaceAll(baseString1, ".", "")

			TickerCollection := client.Database(baseString2).Collection(now.Format(YYYYMMDD))
			result, err := TickerCollection.InsertOne(context.TODO(), datumNode)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("json inserted ", result.InsertedID)
				baseResult := fmt.Sprint(result.InsertedID)
				inter1 := strings.Trim(baseResult, "ObjectID")
				inter2 := strings.Trim(inter1, "  \" ")
				inter3 := strings.Trim(inter2, "  \" ")
				inter4 := strings.Trim(inter3, "  \"\" ")
				inter5 := strings.Trim(inter4, "  ( ) ")
				row := []string{inter5, baseString2, now.Format(YYYYMMDD), time.Now().Format(HHMMSS)}
				_ = csvwriter.Write(row)
				csvwriter.Flush()

			}

		}

		onErrorFunc := func(err error) {
			fmt.Errorf("failed to watch | disconnected from watch. %v", err)
		}

		onCloseFunc := func() {
			fmt.Println("watch connection is being closed")
			csvwriter.Flush()
			csvfile.Close()
		}

		cli := fyerswatch.NewNotifier(apiKey, accessToken).
			WithOnConnectFunc(onConnectFunc).
			WithOnMessageFunc(onMessageFunc).
			WithOnErrorFunc(onErrorFunc).
			WithOnCloseFunc(onCloseFunc)

			// cli.Subscribe(api.SymbolDataTick, "NSE:SBIN-EQ", "NSE:ONGC-EQ")

		symbols := []string{"NSE:SBIN-EQ", "NSE:ONGC-EQ"}
		// create the capped collections
		for i := 0; i < len(symbols); i++ {
			baseString1 := strings.ReplaceAll(symbols[i], ":", "-")
			baseString2 := strings.ReplaceAll(baseString1, ".", "")
			maxSize := int64(30 * 1024 * 1024)
			client.Database(baseString2).CreateCollection(context.TODO(), now.Format(YYYYMMDD), options.CreateCollection().SetCapped(true), options.CreateCollection().SetSizeInBytes(maxSize), options.CreateCollection().SetMaxDocuments(50000))

		}
		cli.Subscribe(api.SymbolDataTick, symbols...)
	}

}
