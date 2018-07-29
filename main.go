package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func main() {
	// server listen address
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = "0.0.0.0:8899"
	}

	// cama notification endpoint
	endpoint := os.Getenv("CAMA_ENDPOINT")
	if endpoint == "" {
		log.Errorln("cama notification endpoint must provided")
		os.Exit(1)
	}

	router := gin.Default()
	router.POST("/webhook", func(c *gin.Context) {
		var msg AlertMessage

		err := c.BindJSON(&msg)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": " successful receive alert notification message!"})

		log.Println("Received msg:", msg)

		log.Println("Sending alert event to CAMA")
		for _, alert := range msg.Alerts {
			event := buildEvent(alert)
			if err := sendEvent(endpoint, event); err != nil {
				log.Errorf("send event to CAMA", err)
				continue
			}
		}

	})

	router.Run(listen)
}

func sendEvent(endpoint string, event *Event) error {
	addr, err := net.ResolveTCPAddr("tcp4", endpoint)
	if err != nil {
		log.Errorf("resolve tcp addr failed. Error: %v", err)
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Errorf("DialTCP to %s failed. Error %s: ", addr.String(), err)
		return err
	}

	defer conn.Close()

	body, err := json.Marshal(event)
	if err != nil {
		log.Errorf("Marshal event failed. Error %s: ", err)
		return err
	}

	//var head string
	//msgLength := len(message)
	//switch {
	//case msgLength < 10:
	//	head = fmt.Sprintf("000%d", msgLength)
	//case msgLength >= 10 && msgLength < 100:
	//	head = fmt.Sprintf("00%d", msgLength)
	//case msgLength >= 100 && msgLength < 1000:
	//	head = fmt.Sprintf("0%d", msgLength)
	//case msgLength >= 1000 && msgLength < 10000:
	//	head = fmt.Sprintf("%d", msgLength)
	//default:
	//	head = "9999"
	//}
	header, err := buildHeader(len(body))
	if err != nil {
		log.Errorf("build header: %v", err)
		return err
	}

	byteArray := [][]byte{
		header,
		body,
	}

	content := bytes.Join(byteArray, nil)

	data, err := utf8ToGbk(content)
	if err != nil {
		log.Errorf("transform data from utf8 to gbk: %v", err)
		return err
	}

	if _, err := conn.Write(data); err != nil {
		log.Errorf("write message to conn failed. Error %s: ", err)
		return err
	}

	log.Infof("sent alert message to %s with content %s", addr, string(data))

	return nil
}

func buildHeader(msgLen int) ([]byte, error) {
	now := time.Now()

	byteArray := [][]byte{
		[]byte(fmt.Sprintf("%08d", msgLen+192)),
		[]byte("DOCKER"),
		[]byte(strings.Repeat("0", 32)),
		[]byte(fmt.Sprintf("%4d%02d%02d", now.Year(), now.Month(), now.Day())),
		[]byte(fmt.Sprintf("%02d%02d%02d", now.Hour(), now.Minute(), now.Second())),
		[]byte("0000000 "),
		[]byte("GDOCKER2018061500000000"),
		[]byte("DOCKER"),
		[]byte("1"),
		[]byte("200001"),
		[]byte("1.0.0"),
		[]byte("01    "),
		[]byte("01"),
		[]byte(strings.Repeat(" ", 12)),
		[]byte("0"),
		[]byte("03"),
		[]byte("0"),
		[]byte("0"),
		[]byte("00"),
		[]byte(strings.Repeat(" ", 10)),
		[]byte(strings.Repeat(" ", 23)),
		[]byte("json  "),
		[]byte(strings.Repeat(" ", 34)),
	}

	return bytes.Join(byteArray, nil), nil
}

func buildEvent(alert Alert) *Event {
	event := new(Event)
	event.ID = alert.Labels["eventID"]
	event.Type = alert.Labels["eventType"]
	event.Level = alert.Labels["eventLevel"]
	event.SourceID = alert.Labels["eventSourceID"]
	event.SourceType = alert.Labels["eventSourceType"]
	event.NodeName = alert.Labels["nodeName"]
	event.NodeIP = alert.Labels["nodeIP"]
	event.FirstTime = alert.StartsAt
	event.LastTime = time.Now().Format("20060101010000")
	event.AlertTimes = 0
	event.AlertKeyType = alert.Labels["alertKeyType"]
	event.AlertKey = alert.Labels["alertKey"]
	event.AlertValue = alert.Labels["alertValue"]
	event.AlertThreshold = alert.Labels["alertThreshold"]
	event.AlertMsg = alert.Annotations["description"]

	return event
}

func utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return d, nil
}
