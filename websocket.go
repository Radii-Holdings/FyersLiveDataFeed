package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rishi-anand/fyers-go-client/api"
	fyerswatch "github.com/rishi-anand/fyers-go-client/websocket"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	YYYYMMDD       = "2006-01-02"
	HHMMSS         = "17:06:04 PM"
	DATAPATH_RAW   = "./dataset/raw/"
	DATAPATH_BACUP = "./dataset/bacup/"
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

// file copier function here

func copy_toPath(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
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
	filename := DATAPATH_RAW + now.Format(YYYYMMDD) + ".csv"
	if _, err := os.Stat(filename); err == nil {
		// file exists
		copy_filename := DATAPATH_BACUP + strings.ReplaceAll(now.Format(time.ANSIC), ":", "-") + ".csv"
		fmt.Println("file is being bscked up", copy_filename, filename)
		copy_toPath(filename, copy_filename)

	}
	csvfile, err := os.Create(filename)
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
				inter6 := strings.Trim(inter5, "  \"\"\" ")
				t := time.Now()
				row := []string{inter6, baseString2, now.Format(YYYYMMDD), strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())}
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
