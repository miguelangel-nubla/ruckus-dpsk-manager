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
	Entries tempEntries `xml:"dpsk"`
}

type tempEntries []*Dpsk

type Entries map[int]*Dpsk

type Dpsk struct {
	ID           int    `xml:"id,attr" json:"id" _dpsk_attr:"id"`
	RoleID       string `xml:"role-id,attr" json:"role-id" _dpsk_attr:"role-id"`
	Mac          string `xml:"mac,attr" json:"mac" _dpsk_attr:"mac"`
	WlansvcID    int    `xml:"wlansvc-id,attr" json:"wlansvc-id" _dpsk_attr:"wlansvc-id"`
	DvlanID      int    `xml:"dvlan-id,attr" json:"dvlan-id" _dpsk_attr:"dvlan-id"`
	User         string `xml:"user,attr" json:"user" _dpsk_attr:"user"`
	LastRekey    string `xml:"last-rekey,attr" json:"last-rekey" _dpsk_attr:"last-rekey"`
	NextRekey    string `xml:"next-rekey,attr" json:"next-rekey" _dpsk_attr:"next-rekey"`
	Expire       string `xml:"expire,attr" json:"expire" _dpsk_attr:"expire"`
	StartPoint   string `xml:"start-point,attr" json:"start-point" _dpsk_attr:"start-point"`
	Passphrase   string `xml:"passphrase,attr" json:"passphrase" _dpsk_attr:"passphrase"`
	IpAddr       string `xml:"ip-addr,attr" json:"ip-addr" _dpsk_attr:"ip-addr"`
	CurSharedNum string `xml:"cur-shared-num,attr" json:"cur-shared-num" _dpsk_attr:"cur-shared-num"`
	Usage        string `xml:"usage,attr" json:"usage" _dpsk_attr:"usage"`
}

var tagMap map[string]string

func init() {
	tagMap = createTagToFieldMap(Dpsk{}, "_dpsk_attr")
}

func createTagToFieldMap(i interface{}, tagKey string) map[string]string {
	result := make(map[string]string)
	valType := reflect.TypeOf(i)

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		tagVal := field.Tag.Get(tagKey)
		if tagVal != "" {
			result[tagVal] = field.Name
		}
	}
	return result
}

func (list *Entries) FindByWlanUser(wlanID int, username string) (*Dpsk, error) {
	for _, entry := range *list {
		if entry.User == username && entry.WlansvcID == wlanID {
			return entry, nil
		}
	}
	return &Dpsk{}, fmt.Errorf("DPSK user not found for username: %s and wlanID: %d", username, wlanID)
}

func (list *Entries) Filter(filters map[string]Filter) (Entries, error) {
	matches := make(Entries)

	for _, dpsk := range *list {
		match := true
		for tag, filter := range filters {
			fieldName, ok := tagMap[tag]
			if !ok {
				return nil, fmt.Errorf("invalid tag: %s", tag)
			}

			field := reflect.ValueOf(*dpsk).FieldByName(fieldName)

			var value string
			switch field.Kind() {
			case reflect.String:
				value = field.String()
			case reflect.Int:
				value = fmt.Sprintf("%d", field.Int())
			default:
				panic("invalid field type") // check Dpsk struct types
			}

			if filter.Test(value) == false {
				// end loop on first failed check
				match = false
				break
			}
		}

		if match {
			matches[dpsk.ID] = dpsk
		}
	}

	return matches, nil
}

func FromXml(xmlData []byte) (Entries, error) {
	var response ajaxResponse
	if err := xml.Unmarshal(xmlData, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %v", err)
	}

	dpskMap := make(Entries)
	for _, dpsk := range response.Response.Apstamgr.DpskList.Entries {
		dpskMap[dpsk.ID] = dpsk
	}

	return dpskMap, nil
}
