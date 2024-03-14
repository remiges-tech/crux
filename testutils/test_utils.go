package testutils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/remiges-tech/alya/wscutils"
)

func MarshalJson(data any) []byte {
	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Fatal("error marshaling")
	}
	return jsonData
}

func ReadJsonFromFile(filepath string) ([]byte, error) {
	// var err error
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("testFile path is not exist")
	}
	defer file.Close()
	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

type TestCasesStruct struct {
	Name             string
	RequestPayload   wscutils.Request
	ExpectedHttpCode int
	TestJsonFile     string
	ExpectedResult   *wscutils.Response
	Url              string
}
