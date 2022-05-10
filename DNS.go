package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"strings"
)

type DNSRecord struct {
	Host    string `json:"host"`
	Address net.IP `json:"address"`
}

type DNS struct {
	records []DNSRecord
}

func NewDNS() *DNS {
	return &DNS{
		records: make([]DNSRecord, 0),
	}
}

func (d *DNS) ReadFromJsonFile(jsonFile string) error {
	var recs []DNSRecord

	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &recs)
	if err != nil {
		return err
	}

	d.records = recs
	return nil
}

func (d *DNS) Resolve(host string) net.IP {
	for _, v := range d.records {
		if strings.EqualFold(host, v.Host) {
			return v.Address
		}
	}

	if ips, err := net.LookupIP(host); err != nil {
		return net.IP{0, 0, 0, 0}
	} else {
		return ips[0]
	}
}
