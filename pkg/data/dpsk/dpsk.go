package dpsk

import (
	"encoding/xml"
	"fmt"
	"reflect"
)

type Filter interface {
	Test(string) bool
}

// Define the struct to match the XML structure
type ajaxResponse struct {
	XMLName  xml.Name `xml:"ajax-response"`
	Response response `xml:"response"`
}

type response struct {
	Type     string       `xml:"type,attr"`
	ID       string       `xml:"id,attr"`
	Apstamgr apstamgrStat `xml:"apstamgr-stat"`
}

type apstamgrStat struct {
	DpskList dpskList `xml:"dpsk-list"`
}

type dpskList struct {
	Entries Entries `xml:"dpsk"`
}

type Entries []*Dpsk

type Dpsk struct {
	ID           int    `xml:"id,attr"`
	RoleID       string `xml:"role-id,attr"`
	Mac          string `xml:"mac,attr"`
	WlansvcID    int    `xml:"wlansvc-id,attr"`
	DvlanID      int    `xml:"dvlan-id,attr"`
	User         string `xml:"user,attr"`
	LastRekey    string `xml:"last-rekey,attr"`
	NextRekey    string `xml:"next-rekey,attr"`
	Expire       string `xml:"expire,attr"`
	StartPoint   string `xml:"start-point,attr"`
	Passphrase   string `xml:"passphrase,attr"`
	IpAddr       string `xml:"ip-addr,attr"`
	CurSharedNum string `xml:"cur-shared-num,attr"`
	Usage        string `xml:"usage,attr"`
}

func (list *Entries) FindByWlanUser(wlanID int, username string) (*Dpsk, error) {
	for _, entry := range *list {
		if entry.User == username && entry.WlansvcID == wlanID {
			return entry, nil
		}
	}
	return &Dpsk{}, fmt.Errorf("DPSK user not found for username: %s and wlanID: %d", username, wlanID)
}

func (list *Entries) Filter(filters map[string]Filter) (*Entries, error) {
	var matches Entries

	for _, dpsk := range *list {
		match := true
		for tag, filter := range filters {
			field := reflect.ValueOf(*dpsk).FieldByName(tag)
			var value string
			switch field.Kind() {
			case reflect.String:
				value = field.String()
			case reflect.Int:
				value = fmt.Sprintf("%d", field.Int())
			default:
				return nil, fmt.Errorf("unsupported type: %v", field.Kind())
			}

			if filter.Test(value) == false {
				// end loop on first failed check
				match = false
				break
			}
		}

		if match {
			matches = append(matches, dpsk)
		}
	}

	return &matches, nil
}

func FromXml(xmlData []byte) (*Entries, error) {
	var response ajaxResponse
	if err := xml.Unmarshal(xmlData, &response); err != nil {
		return &Entries{}, fmt.Errorf("error unmarshalling XML: %v", err)
	}

	return &response.Response.Apstamgr.DpskList.Entries, nil
}
