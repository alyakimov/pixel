package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func SetSecretCookie(response http.ResponseWriter, name string, value string) {

	secret := viper.GetString("cookie.secret")
	version := viper.GetInt("cookie.version")
	keyVersion := viper.GetInt("cookie.key_version")

	expiration := time.Now().Add(365 * 24 * time.Hour)
	signedValue, err := createSignedValue(secret, name, value, version, keyVersion)

	if err == nil {
		cookie := http.Cookie{Name: name, Value: signedValue, Expires: expiration}
		http.SetCookie(response, &cookie)
	}
}

func GetSecretCookie(request *http.Request, name string) (string, error) {

	cookie, err := request.Cookie(name)

	if err != nil {
		return "", errors.New("cookie not found")
	}

	cookieValue := cookie.Value

	secret := viper.GetString("cookie.secret")
	maxAgeDays := viper.GetInt("cookie.max_age_days")
	version := viper.GetInt("cookie.version")

	value, err := decodeSignedValue(secret, name, cookieValue, maxAgeDays, version)

	return value, err
}

func createSignedValue(secret string, name string, value string, version int, keyVersion int) (string, error) {

	timestamp := strconv.Itoa(getUnixTimestamp())

	data := []byte(value)
	value = base64.StdEncoding.EncodeToString(data)

	if version == 2 {
		toSign := strings.Join([]string{
			"2",
			formatField(strconv.Itoa(keyVersion)),
			formatField(timestamp),
			formatField(name),
			formatField(value),
			"",
		}, "|")

		signature := createSignatureV2(secret, toSign)

		return fmt.Sprintf("%s%s", toSign, signature), nil

	} else {
		return "", errors.New(fmt.Sprintf("Unsupported version: %s", version))
	}
}

func formatField(value string) string {
	return fmt.Sprintf("%d:%s", len(value), value)
}

func consumeField(s string) (string, string) {

	b := strings.SplitN(s, ":", 2)
	length, rest := b[0], string(b[1])

	n, _ := strconv.ParseInt(length, 10, 0)
	fieldValue := rest[:n]

	if rest[n:n+1] != "|" {
		panic("malformed v2 signed value field")
	}

	rest = string(rest[n+1:])

	return fieldValue, rest
}

func decodeFieldsV2(value string) (int, string, string, string, string) {

	rest := value[2:]
	keyVersionStr, rest := consumeField(rest)
	keyVersion, _ := strconv.ParseInt(keyVersionStr, 10, 0)

	timestamp, rest := consumeField(rest)
	nameField, rest := consumeField(rest)
	valueField, passedSign := consumeField(rest)

	return int(keyVersion), timestamp, nameField, valueField, passedSign
}

func getVersion(value string) int {

	signedValueVersionRe := regexp.MustCompile("^([1-9][0-9]*)|(.*)$")
	match := signedValueVersionRe.FindStringSubmatch(value)

	version := 1

	if match != nil {
		numVersion, err := strconv.ParseInt(match[0], 10, 0)

		if err == nil {

			version = int(numVersion)

			if version > 999 {
				version = 1
			}
		}
	}

	return version
}

func decodeSignedValue(secret string, name string, value string, maxAgeDays int, minVersion int) (string, error) {

	version := getVersion(value)

	if version < minVersion {
		return "", errors.New(fmt.Sprintf("Unsupported min_version %d", minVersion))
	}

	if version == 2 {
		signedValue, err := decodeSignedValueV2(secret, name, value, maxAgeDays)

		return signedValue, err
	} else {
		return "", errors.New(fmt.Sprintf("Unsupported version %d", version))
	}
}

func decodeSignedValueV2(secret string, name string, value string, maxAgeDays int) (string, error) {
	_, timestampString, nameField, valueField, passedSign := decodeFieldsV2(value)

	signdedString := value[:len(value)-len(passedSign)]
	expectedSig := createSignatureV2(secret, signdedString)

	if passedSign != expectedSig {
		return "", errors.New("Invalid sign")
	}

	if nameField != name {
		return "", errors.New("Invalid name field")
	}

	timestamp, _ := strconv.ParseInt(timestampString, 10, 0)

	if int(timestamp) < getUnixTimestamp()-maxAgeDays*86400 {
		return "", errors.New("The signature has expired.")
	}

	decodeValueField, err := base64.StdEncoding.DecodeString(valueField)

	if err != nil {
		return "", errors.New("Base64 decode error")
	}

	return string(decodeValueField), nil
}

func getSignatureKeyVersion(value string) (int, error) {
	version := getVersion(value)

	if version < 2 {
		return 0, errors.New("")
	}

	keyVersion, _, _, _, _ := decodeFieldsV2(value)

	return keyVersion, nil
}

func createSignatureV2(secret string, message string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))

	return hex.EncodeToString(h.Sum(nil))
}

func getUnixTimestamp() int {
	return int(time.Now().Unix())
}
