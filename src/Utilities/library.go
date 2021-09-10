package Utilities

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strconv"

	base58 "github.com/itchyny/base58-go"
)

func ReadEnvironmentVariable(variableName string, valueType reflect.Kind, defaultValue interface{}) interface{} {
	var value interface{}
	if envVariableValue := os.Getenv(variableName); envVariableValue != "" {
		if valueType == reflect.Int32 {
			value = parseStringToInteger(envVariableValue, 10, 32)
		} else if valueType == reflect.Int64 {
			value = parseStringToInteger(envVariableValue, 10, 64)
		} else {
			value = envVariableValue
		}
	} else {
		value = defaultValue
	}
	return value
}

func parseStringToInteger(stringValue string, base int, bitness int) int64 {
	var value int64
	if parsedInt, err := strconv.ParseInt(stringValue, base, bitness); err == nil && parsedInt > 0 {
		value = parsedInt
	} else {
		value = 0
	}
	return value
}

func sha256encoded(input []byte) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", errors.New("error occurred while doing Base encoding")
	}
	return string(encoded), nil
}

func GenerateShortLink(link string) (string, error) {
	urlHashBytes := sha256encoded([]byte(link))
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, error := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	if error != nil {
		return "", error
	}
	return finalString[:8], nil
}
