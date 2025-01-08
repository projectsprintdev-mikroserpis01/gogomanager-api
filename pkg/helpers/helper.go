package helpers

import (
	"bufio"
	"database/sql/driver"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"

	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/log"
)

func Contains(search string, words []string) bool {
	for _, word := range words {
		if search == word {
			return true
		}
	}

	return false
}

func ReadFile(filepath string, separator string) ([]string, error) {
	file, err := os.Open(filepath)

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[HELPERS][ReadFiles] failed to open file")
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	results := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()

		results = append(results, line)
	}

	if err := scanner.Err(); err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[HELPERS][ReadFiles] failed to scan file")
		return nil, err
	}

	return results, nil
}

func CheckRowsAffected(rows int64) error {
	if rows == 0 {
		return domain.ErrNotFound
	}

	if rows > 1 {
		return fmt.Errorf("weird behaviour. rows affected : %v", rows)
	}

	return nil
}

func GenerateRandomString(lenght int) string {
	alphaNumRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	randomRune := make([]rune, lenght)

	for i := 0; i < lenght; i++ {
		randomRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}

	return string(randomRune)
}

// Helper function to convert struct fields into a slice of interface{}
func StructToSlice(i interface{}) []interface{} {
	var result []interface{}
	val := reflect.ValueOf(i)

	for i := 0; i < val.NumField(); i++ {
		// if the field is of type entity, skip it
		if strings.Contains(val.Field(i).Type().String(), "entity") {
			continue
		}

		result = append(result, val.Field(i).Interface())
	}
	return result
}

// Helper function to convert slice of interface{} into a slice of driver.Value
func ConvertToDriverValue(values []interface{}) []driver.Value {
	driverValues := make([]driver.Value, len(values))
	for i, v := range values {
		driverValues[i] = driver.Value(v)
	}
	return driverValues
}
