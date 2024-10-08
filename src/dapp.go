package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"dapp/rollups"
)

var (
	infolog = log.New(os.Stderr, "[ info ]  ", log.Lshortfile)
	errlog  = log.New(os.Stderr, "[ error ] ", log.Lshortfile)
)

var (
	users        []string
	toUpperTotal int
)

func HandleAdvance(data *rollups.AdvanceResponse) error {
	dataMarshal, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("HandleAdvance: failed marshaling json: %w", err)
	}
	infolog.Println("Received advance request data", string(dataMarshal))

	metadata := data.Metadata
	sender := metadata.MsgSender
	payload := data.Payload

	sentence, err := rollups.Hex2Str(payload)
	if err != nil {
		// REPORT
		report := rollups.ReportRequest{Payload: rollups.Str2Hex(sentence)}

		_, err = rollups.SendReport(&report)
		if err != nil {
			return err
		}

	}

	users = append(users, sender)
	toUpperTotal++

	
	sentence = strings.ToUpper(sentence)
	sentence = makeStrLonger(sentence)
	// NOTICE
	notice := rollups.NoticeRequest{Payload: rollups.Str2Hex(sentence)}
	_, err = rollups.SendNotice(&notice)
	if err != nil {
		return err
	}

	return nil
}

func HandleInspect(data *rollups.InspectResponse) error {
	dataMarshal, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("HandleInspect: failed marshaling json: %w", err)
	}
	infolog.Println("Received inspect request data", string(dataMarshal))

	payload := data.Payload
	route, err := rollups.Hex2Str(payload)
	if err != nil {
		return err
	}

	responseObject := rollups.ReportRequest{}
	switch route {
	case "list":
		responseObject = rollups.ReportRequest{Payload: rollups.Str2Hex(strings.Join(users, " "))}
	case "total":
		totalNum := strconv.Itoa(toUpperTotal)
		responseObject = rollups.ReportRequest{Payload: rollups.Str2Hex(totalNum)}
	default:
		responseObject = rollups.ReportRequest{Payload: rollups.Str2Hex("route not implemented")}
	}

	_, err = rollups.SendReport(&responseObject)
	if err != nil {
		return err
	}

	return nil
}

func Handler(response *rollups.FinishResponse) error {
	var err error

	switch response.Type {
	case "advance_state":
		data := new(rollups.AdvanceResponse)
		if err = json.Unmarshal(response.Data, data); err != nil {
			return fmt.Errorf("Handler: Error unmarshaling advance:%v", err)
		}
		err = HandleAdvance(data)
	case "inspect_state":
		data := new(rollups.InspectResponse)
		if err = json.Unmarshal(response.Data, data); err != nil {
			return fmt.Errorf("Handler: Error unmarshaling inspect:,%v", err)
		}
		err = HandleInspect(data)
	}
	return err
}

func main() {
	finish := rollups.FinishRequest{Status: "accept"}

	for {
		infolog.Println("Sending finish")
		res, err := rollups.SendFinish(&finish)
		if err != nil {
			errlog.Panicln("Error: error making http request: ", err)
		}
		infolog.Println("Received finish status ", strconv.Itoa(res.StatusCode))

		if res.StatusCode == 202 {
			infolog.Println("No pending rollup request, trying again")
		} else {

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				errlog.Panicln("Error: could not read response body: ", err)
			}

			var response rollups.FinishResponse
			err = json.Unmarshal(resBody, &response)
			if err != nil {
				errlog.Panicln("Error: unmarshaling body:", err)
			}

			finish.Status = "accept"
			err = Handler(&response)
			if err != nil {
				errlog.Println(err)
				finish.Status = "reject"
			}
		}
	}
}

func makeStrLonger(word string) string {
	result:= ""
	for _, char := range word {
		if unicode.IsLetter(char) {
			b := strings.ToLower(string(char))
			count := 0
			for _, value := range b {
				count = int(value - 'a')
				
			}
			for i := 0; i < count; i++ {
				result += string(char)
			}
		}
		result += string(char)
	}
	
	return result
}
