package main


import (
    "time"
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
    CampaignId int
    Uuid string
    Msisdn string
    RemoteIp string
    UserAgent string
    Referer string
}


type Defcode struct {
    Msisdn int
    Uuid string
}
