package main


import (
    "io"
    "net/http"
    "encoding/base64"
)


func Index(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type","image/gif")

	output, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAP///wAAACwAAAAAAQABAAACAkQBADs=")
    io.WriteString(response, string(output))
}