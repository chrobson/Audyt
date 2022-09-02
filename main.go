package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexmullins/zip"
)

// Process password encrypted csv audit file which is in zip format
func GetCSVFromZipURL(fileURL, filePassword string) (ioReader io.Reader, err error) {

	resp, err := http.Get(fileURL)
	if err != nil {
		fmt.Println(err)
		//return ioReader, err
	}
	defer resp.Body.Close()

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return ioReader, err
	}

	b := bytes.NewReader(buf.Bytes())
	r, err := zip.NewReader(b, int64(b.Len()))
	if err != nil {
		return ioReader, err
	}

	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword(filePassword)
		}
		//todo error handling for bad password
		ioReader, err = f.Open()
		if err != nil {
			return ioReader, err
		}

		//return ioReader, err
	}

	return ioReader, err
}

// parse string fields for floats & remove comma thousands separataros
func parseField(s string) float64 {
	var f float64
	s = strings.ReplaceAll(s, ",", "")
	f, _ = strconv.ParseFloat(s, 64)
	return f
}

func unmarschalCSV(rows [][]string) []Audit {
	var audits []Audit
	for _, r := range rows {
		datetime, _ := time.Parse("2006/01/02 3:04 PM", r[0])
		ammount := parseField(r[6])
		accuredstarscoins := parseField(r[7])
		tmoney := parseField(r[8])
		wmoney := parseField(r[9])
		balance := parseField(r[10])
		totalstarcoins := parseField(r[11])
		ttmoney := parseField(r[12])
		wwmoney := parseField(r[13])

		audit := Audit{DateTime: datetime,
			Action:            r[1],
			Tournament:        r[2],
			Game:              r[3],
			Currency:          r[5],
			Ammount:           ammount,
			Accuredstarscoins: accuredstarscoins,
			Tmoney:            tmoney,
			Wmoney:            wmoney,
			Balance:           balance,
			Totalstarcoins:    totalstarcoins,
			TTmoney:           ttmoney,
			WWmoney:           wwmoney,
		}
		audits = append(audits, audit)
	}
	return audits

}

// Audit represents pokerstars audit structure
type Audit struct {
	DateTime          time.Time `csv:"datetime"`
	Action            string    `csv:"action"`
	Tournament        string    `csv:"tournament"`
	Game              string    `csv:"game"`
	Currency          string    `csv:"currency"`
	Ammount           float64   `csv:"ammount"`
	Tmoney            float64   `csv:"tmoney"`
	Accuredstarscoins float64   `csv:"accuredstarscoins"`
	Wmoney            float64   `csv:"wmoney"`
	Balance           float64   `csv:"balance"`
	Totalstarcoins    float64   `csv:"totalstarcoins"`
	TTmoney           float64   `csv:"ttmoney"`
	WWmoney           float64   `csv:"wwmoney"`
}

func main() {
	//todo accept only zip files from http://reports.rationalwebservices.com/reports/filename.zip
	f, err := GetCSVFromZipURL("http://reports.rationalwebservices.com/reports/xxx.zip", "12345")
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	rows_pre, _ := reader.ReadAll()
	//Skip first three rows
	rows := rows_pre[3:]
	audits := unmarschalCSV(rows)
	var sum float64
	for _, a := range audits {
		//if strings.Contains(a.Tournament, "Spin") && !strings.Contains(a.Tournament, "Max") {
		if strings.Contains(a.Action, "Chest") {
			sum = sum + a.Ammount + (a.Accuredstarscoins / 100)
		}
	}
	fmt.Println(sum)
}
