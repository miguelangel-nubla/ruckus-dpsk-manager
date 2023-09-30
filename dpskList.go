package main

import (
	"encoding/xml"
	"fmt"
)

// Define the struct to match the XML structure
type AjaxResponse struct {
	XMLName  xml.Name `xml:"ajax-response"`
	Response Response `xml:"response"`
}

type Response struct {
	Type      string      `xml:"type,attr"`
	ID        string      `xml:"id,attr"`
	Apstamgr  ApstamgrStat `xml:"apstamgr-stat"`
}

type ApstamgrStat struct {
	DpskList DpskList `xml:"dpsk-list"`
}

type DpskList struct {
	DpskEntries DpskEntries `xml:"dpsk"`
}

type DpskEntries []Dpsk

type Dpsk struct {
	ID          int `xml:"id,attr"`
	RoleID      string `xml:"role-id,attr"`
	Mac         string `xml:"mac,attr"`
	WlansvcID   int `xml:"wlansvc-id,attr"`
	DvlanID     int `xml:"dvlan-id,attr"`
	User        string `xml:"user,attr"`
	LastRekey   string `xml:"last-rekey,attr"`
	NextRekey   string `xml:"next-rekey,attr"`
	Expire      string `xml:"expire,attr"`
	StartPoint  string `xml:"start-point,attr"`
	Passphrase  string `xml:"passphrase,attr"`
	IpAddr      string `xml:"ip-addr,attr"`
	CurSharedNum string `xml:"cur-shared-num,attr"`
	Usage       string `xml:"usage,attr"`
}

// Custom type to encapsulate DpskEntries and provide a search method
type DpskData struct {
	Entries DpskEntries
}

// Method to search for a specific DPSK user and return its passphrase
func (d *DpskData) FindPassphrase(wlanID int, username string) (string, error) {
    for _, entry := range d.Entries {
        if entry.User == username && entry.WlansvcID == wlanID {
            return entry.Passphrase, nil
        }
    }
    return "", fmt.Errorf("DPSK user not found for username: %s and wlanID: %d", username, wlanID)
}

func NewDpskData(xmlData []byte) (*DpskData, error) {
	var response AjaxResponse
	if err := xml.Unmarshal(xmlData, &response); err != nil {
		return &DpskData{}, fmt.Errorf("error unmarshalling XML: %v", err)
	}

	return &DpskData{Entries: response.Response.Apstamgr.DpskList.DpskEntries}, nil
}
