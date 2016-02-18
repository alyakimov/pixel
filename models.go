package main


import (
    "time"
    "os"
    "database/sql"
    "encoding/json"    
)

type Campaign struct {
    Id int
    PartnerId int
    Name string
    Status int
    Created time.Time
    Updated time.Time
}


type CampaignLog struct {
    CampaignId  int     `json:"campaignId"`
    Uuid        string  `json:"uuid"`
    Msisdn      string  `json:"msisdn"`
    RemoteIp    string  `json:"remoteIp"`
    UserAgent   string  `json:"userAgent"`
    Referer     string  `json:"referer"`
}


type Defcode struct {
    Msisdn int
    Uuid string
}


func GetCampaignByName(db *sql.DB, name string) (*Campaign, error) {
    const query = "SELECT id, partner_id, name, status, created, updated FROM campaigns where status = 1 AND name = ?"

    var retval Campaign
    err := db.QueryRow(query, name).Scan(&retval.Id, &retval.PartnerId, &retval.Name, &retval.Status, &retval.Created, &retval.Updated)    

    return &retval, err
}

func AddCampaignLog(db *sql.DB, campaignLog *CampaignLog) error {
    const query = "INSERT INTO campaigns_log (campaign_id, uuid, msisdn, remote_ip, user_agent, referer, created) VALUES (?, ?, ?, ?, ?, ?, NOW())"

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }

    _, err = stmt.Exec(
        campaignLog.CampaignId, 
        campaignLog.Uuid, 
        campaignLog.Msisdn, 
        campaignLog.RemoteIp, 
        campaignLog.UserAgent, 
        campaignLog.Referer,
    )

    return err
}

func AddCampaignLogIntoFile(filename string, campaignLog *CampaignLog) error {
    
    data, err := json.Marshal(campaignLog)

    if err != nil {
        return err
    }

    f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        return err
    }

    defer f.Close()

    _, err = f.WriteString(string(data) + "'\n")

    return err

}

func GetDefcodeByMsisdn(db *sql.DB, msisdn string) (*Defcode, error) {
    const query = "SELECT uuid, msisdn FROM defcodes where msisdn = ?"

    var retval Defcode
    err := db.QueryRow(query, msisdn).Scan(&retval.Uuid, &retval.Msisdn)    

    return &retval, err
}

