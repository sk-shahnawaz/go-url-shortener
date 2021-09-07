package utilities

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"

	base58 "github.com/itchyny/base58-go"
	"github.com/jackc/pgx/v4/pgxpool"
)

var store map[string]string = make(map[string]string)

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

func GenerateShortLink(link string, dbClient *pgxpool.Pool) (string, error) {
	urlHashBytes := sha256encoded([]byte(link))
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, error := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	if error != nil {
		return "", error
	}
	if useInMemoryDb := ReadEnvironmentVariable("USE_IN_MEMORY_DB", reflect.String, "Y"); strings.ToUpper(useInMemoryDb.(string)) == "N" {
		query := fmt.Sprintf(`INSERT INTO urls (original, shortened) VALUES ('%s', '%s')`, link, finalString[:8])
		_, err := dbClient.Exec(context.Background(), query)
		if err != nil {
			return "", err
		}
	} else {
		store[link] = finalString[:8]
	}
	return finalString[:8], nil

}

func ResolveShortenedLink(resolvable string, dbClient *pgxpool.Pool) (string, error) {
	if useInMemoryDb := ReadEnvironmentVariable("USE_IN_MEMORY_DB", reflect.String, "Y"); strings.ToUpper(useInMemoryDb.(string)) == "N" {
		query := fmt.Sprintf(`SELECT original FROM urls WHERE shortened = '%s'`, resolvable)
		rows, err := dbClient.Query(context.Background(), query)
		if err != nil {
			return "", err
		}
		if rows != nil {
			var original string
			for rows.Next() {
				err := rows.Scan(&original)
				if err != nil {
					return "", err
				}
			}
			return original, nil
		}
		return "", nil
	} else {
		for link, shortenedLink := range store {
			if shortenedLink == resolvable {
				return link, nil
			}
		}
		return "", errors.New("no entry found")
	}
}
