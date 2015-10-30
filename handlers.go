package main


import (
    "io"
    "log"
    "time"
    "errors"
    "strings"
    "net/http"
    "encoding/base64"
)


func Index(response http.ResponseWriter, request *http.Request) {	

    campaignName, err := getCampaignName(request)
    
    if err != nil {
        writeImage(response)
        return
    }
        
    msisdn, err := getMsisdn(request)

    if err != nil {
        writeImage(response)
        return
    }

    remoteIp := getRemoteIp(request)
    userAgent := getUserAgent(request)
    referer := getReferer(request)
    uuid, _ := getUuid(request)

    campaignLog := CampaignLog{
        CampaignId: 0, 
        Uuid: uuid, 
        Msisdn: msisdn, 
        RemoteIp: remoteIp, 
        UserAgent: userAgent, 
        Referer: referer,
    }

    db := GetConnection()
    defer db.Close()

    campaign, err := GetCampaignByName(db, campaignName)

    if err != nil {
        log.Println(err)
    
    } else {

        campaignLog.CampaignId = campaign.Id

        err = AddCampaignLog(db, campaignLog)
        if err != nil {
            log.Fatal(err)
        }
    }

    writeImage(response)
}


func Redirect(response http.ResponseWriter, request *http.Request) {
    
    backUrl, err := getBackUrl(request)
    
    if err != nil {
        log.Fatal(err)
    }

    msisdn, err := getMsisdn(request)

    if err != nil {
        redirect(response, request, "....................", backUrl)
        return
    }

    db := GetConnection()
    defer db.Close()

    defcode, err := GetDefcodeByMsisdn(db, msisdn)
    
    if err != nil {
        uuid := "...................."
    } else {
        uuid := defcode.Uuid    
    }

    redirect(response, request, uuid, backUrl)
}


func writeImage(response http.ResponseWriter){
    response.Header().Set("Content-Type", "image/gif")
    response.Header().Set("Cache-Control", "private, no-cache, no-cache=Set-Cookie, proxy-revalidate")
    response.Header().Set("Pragma", "no-cache")
    response.Header().Set("Expires", "Wed, 17 Sep 1975 21:32:10 GMT")

    output, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAP///wAAACwAAAAAAQABAAACAkQBADs=")
    io.WriteString(response, string(output))
}

func redirect(response http.ResponseWriter, request *http.Request, uuid string, backUrl string) {
    timestamp := getUnixTimestamp()

    backUrl = strings.Replace(backUrl, "$UUID", uuid, 1)
    backUrl = strings.Replace(backUrl, "$RND", timestamp, 1)

    http.Redirect(response, request, backUrl, 301)
}

func getBackUrl(request *http.Request) (string, error) {
    backUrl := request.FormValue("back_url")

    if len(backUrl) > 0 {
        return backUrl, nil
    } else {
        return "", errors.New("Invalid Back Url")
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

func getCookieMsisdn(request *http.Request) (string, error){
    cookieMsisdn, err := GetSecretCookie(request, "msisdn") 

    return cookieMsisdn, err
}

func setCookieMsisdn(response http.ResponseWriter, value string){
    SetSecretCookie(response, "msisdn", value)
}

func getRemoteIp(request *http.Request) string {
    remoteIp := strings.Split(request.RemoteAddr, ":")[0]
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

func getUnixTimestamp() int32 {
    return int32(time.Now().Unix())
}

