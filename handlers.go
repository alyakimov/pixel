package main


import (
    "io"
    "log"
    "errors"
    "strings"
    "net/http"
    "encoding/base64"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/spf13/viper"
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

    dsl := viper.GetString("connections.onepixel")

    connect, err := sql.Open("mysql", dsl)
    defer connect.Close()
    
    if err != nil {
        log.Fatal(err)
    }

    campaign, err := GetCampaignByName(connect, campaignName)

    if err == nil {

        campaignLog.CampaignId = campaign.Id

        err = AddCampaignLog(connect, campaignLog)
        if err != nil {
            log.Fatal(err)
        }
    }

    writeImage(response)
}


func Redirect(response http.ResponseWriter, request *http.Request) {   

    // guidNotFound := "...................."
    
    backUrl, err := getBackUrl(request)
    
    if err != nil {
        log.Fatal(err)
    }

    msisdn, err := getMsisdn(request)

    if err != nil {
        log.Fatal(err)
    }

    dsl := viper.GetString("connections.onepixel")

    connect, err := sql.Open("mysql", dsl)
    defer connect.Close()
    
    if err != nil {
        log.Fatal(err)
    }

    stmt, err := connect.Prepare("SELECT uuid FROM defcodes where msisdn = ?")
    if err != nil {
        log.Fatal(err)
    }
    row := stmt.QueryRow(msisdn)

    var uuid string    
    row.Scan(&uuid)

    backUrl = strings.Replace(backUrl, "$UUID", uuid, 1)
    backUrl = strings.Replace(backUrl, "$RND", "12345", 1)

    http.Redirect(response, request, backUrl, 301)
}


func writeImage(response http.ResponseWriter){
    response.Header().Set("Content-Type", "image/gif")
    response.Header().Set("Cache-Control", "private, no-cache, no-cache=Set-Cookie, proxy-revalidate")
    response.Header().Set("Pragma", "no-cache")
    response.Header().Set("Expires", "Wed, 17 Sep 1975 21:32:10 GMT")

    output, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAP///wAAACwAAAAAAQABAAACAkQBADs=")
    io.WriteString(response, string(output))
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

