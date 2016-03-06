package main

import (
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Index(response http.ResponseWriter, request *http.Request) {

	campaignName, err := getCampaignName(request)

	if err != nil {
		log.Println(err)
		writeImage(response)
		return
	}

	campaignLog, err := getCampaingLog(request, response)
	if err != nil {
		writeImage(response)
		return
	}

	db := GetConnection()

	campaign, err := GetCampaignByName(db, campaignName)

	if err != nil {
		log.Println(err)

	} else {

		campaignLog.CampaignId = campaign.Id

		go AddCampaignLog(db, campaignLog)
	}

	writeImage(response)
}

func Redirect(response http.ResponseWriter, request *http.Request) {

	campaignName, err := getCampaignName(request)

	if err != nil {
		log.Println(err)

		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	backUrl, err := getBackUrl(request)

	if err != nil {
		log.Println(err)

		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	guid := getDefaultGuid()

	campaignLog, err := getCampaingLog(request, response)
	if err != nil {
		redirect(response, request, guid, backUrl)
		return
	}

	db := GetConnection()

	campaign, err := GetCampaignByName(db, campaignName)

	if err != nil {
		log.Println(err)

		http.Error(response, "Partner Not Found", http.StatusBadRequest)
		return

	} else {

		campaignLog.CampaignId = campaign.Id

		go AddCampaignLog(db, campaignLog)
	}

	defcode, err := GetDefcodeByMsisdn(db, campaignLog.Msisdn)

	if err == nil {
		guid = defcode.Uuid
	}

	redirect(response, request, guid, backUrl)
}

func writeImage(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "image/gif")
	response.Header().Set("Cache-Control", "private, no-cache, no-cache=Set-Cookie, proxy-revalidate")
	response.Header().Set("Pragma", "no-cache")
	response.Header().Set("Expires", "Wed, 17 Sep 1975 21:32:10 GMT")

	output, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAP///wAAACwAAAAAAQABAAACAkQBADs=")
	io.WriteString(response, string(output))
}

func redirect(response http.ResponseWriter, request *http.Request, uuid string, backUrl string) {
	timestamp := getUnixTimestamp()

	backUrl = strings.Replace(backUrl, "$UID", uuid, 1)
	backUrl = strings.Replace(backUrl, "$RND", strconv.Itoa(timestamp), 1)

	http.Redirect(response, request, backUrl, 302)
}

func getBackUrl(request *http.Request) (string, error) {
	backUrl := request.FormValue("backurl")

	if len(backUrl) > 0 {
		return backUrl, nil
	} else {

		backUrl = request.FormValue("back_url")

		if len(backUrl) > 0 {
			return backUrl, nil
		} else {
			return "", errors.New("Invalid Back Url")
		}
	}
}

func getCampaignName(request *http.Request) (string, error) {
	campaignName := request.FormValue("campaign")

	if len(campaignName) > 0 {
		return campaignName, nil
	} else {
		return "", errors.New("Invalid campaign name")
	}
}

func getMsisdn(request *http.Request) (string, error) {
	msisdn := request.Header.Get("X-Nokia-MSISDN")
	msisdnValid := request.Header.Get("X-MSISDN-VALID")

	if msisdnValid == "YES" && msisdn != "" {
		return msisdn, nil
	} else {
		return "", errors.New("Invalid msisdn header")
	}
}

func getCookieMsisdn(request *http.Request) (string, error) {
	cookieMsisdn, err := GetSecretCookie(request, "msisdn")

	return cookieMsisdn, err
}

func setCookieMsisdn(response http.ResponseWriter, value string) {
	SetSecretCookie(response, "msisdn", value)
}

func getRemoteIp(request *http.Request) string {
	xForwardFor := request.Header.Get("X-Forwarded-For")

	var remoteIp string
	if xForwardFor != "" {
		remoteIp = strings.Split(xForwardFor, ",")[0]
	} else {
		remoteIp = strings.Split(request.RemoteAddr, ":")[0]
	}

	return remoteIp
}

func getUserAgent(request *http.Request) string {
	return request.UserAgent()
}

func getReferer(request *http.Request) string {
	return request.Referer()
}

func getUuid(request *http.Request) (string, error) {
	uuid := request.FormValue("uuid")

	if uuid != "" {
		return uuid, nil
	} else {
		return "", errors.New("uuid is empty")
	}
}

func getDefaultGuid() string {
	return "...................."
}

func getCampaingLog(request *http.Request, response http.ResponseWriter) (*CampaignLog, error) {

	msisdn, err := getMsisdn(request)

	if err != nil {

		msisdnCookie, err := getCookieMsisdn(request)
		if err != nil {
			return nil, err
		} else {
			msisdn = msisdnCookie
		}
	} else {
		setCookieMsisdn(response, msisdn)
	}

	remoteIp := getRemoteIp(request)
	userAgent := getUserAgent(request)
	referer := getReferer(request)
	uuid, _ := getUuid(request)

	campaignLog := &CampaignLog{
		CampaignId: "",
		Uuid:       uuid,
		Msisdn:     msisdn,
		RemoteIp:   remoteIp,
		UserAgent:  userAgent,
		Referer:    referer,
		Created:    time.Now(),
	}

	return campaignLog, nil
}
