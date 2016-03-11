package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Campaign struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	PartnerId bson.ObjectId `json:"partnerId" bson:"partnerId"`
	Name      string        `json:"name" bson:"name"`
	Status    int           `json:"status" bson:"status"`
	Created   time.Time     `json:"created" bson:"created"`
	Updated   time.Time     `json:"updated" bson:"updated"`
}

type CampaignLog struct {
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	CampaignId bson.ObjectId `json:"campaignId" bson:"campaignId"`
	Uuid       string        `json:"uuid" bson:"uuid"`
	Msisdn     string        `json:"msisdn" bson:"msisdn"`
	RemoteIp   string        `json:"remoteIp" bson:"remoteIp"`
	UserAgent  string        `json:"userAgent" bson:"userAgent"`
	Referer    string        `json:"referer" bson:"referer"`
	Created    time.Time     `json:"created" bson:"created"`
}

type Defcode struct {
	Msisdn int    `json:"msisdn" bson:"msisdn"`
	Uuid   string `json:"uuid" bson:"uuid"`
}

func GetCampaignByName(db *mgo.Session, name string) (*Campaign, error) {
	session := db.Copy()
	defer session.Close()

	campaign := Campaign{}
	campaigns := session.DB("test").C("campaigns")
	err := campaigns.Find(bson.M{"status": 1, "name": name}).One(&campaign)

	return &campaign, err
}

func AddCampaignLog(db *mgo.Session, campaignLog *CampaignLog) error {
	session := db.Copy()
	defer session.Close()

	campaignsLog := session.DB("test").C("campaigns_log")
	err := campaignsLog.Insert(&campaignLog)

	return err
}

func GetDefcodeByMsisdn(db *mgo.Session, msisdn string) (*Defcode, error) {
	session := db.Copy()
	defer session.Close()

	defcode := Defcode{}
	defcodes := session.DB("test").C("defcodes")
	err := defcodes.Find(bson.M{"msisdn": msisdn}).One(&defcode)

	return &defcode, err
}
