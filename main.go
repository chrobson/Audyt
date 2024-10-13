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

// Audit represents pokerstars audit structure
type Audit struct {
	DateTime          time.Time `csv:"datetime"`
	Action            string    `csv:"action"`
	Tournament        string    `csv:"tournament"`
	Game              string    `csv:"game"`
	Currency          string    `csv:"currency"`
	Amount            float64   `csv:"amount"`
	Tmoney            float64   `csv:"tmoney"`
	AccruedStarsCoins float64   `csv:"accruedstarscoins"`
	Wmoney            float64   `csv:"wmoney"`
	Balance           float64   `csv:"balance"`
	TotalStarCoins    float64   `csv:"totalstarcoins"`
	TTmoney           float64   `csv:"ttmoney"`
	WWmoney           float64   `csv:"wwmoney"`
}

// DownloadAndUnzipFile downloads a zip file from a URL
func DownloadAndUnzipFile(fileURL string) (*zip.Reader, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to copy response body: %w", err)
	}

	reader := bytes.NewReader(buf.Bytes())
	return zip.NewReader(reader, int64(buf.Len()))
}

// ExtractCSVFromZip extracts and decrypts the CSV file from a zip archive
func ExtractCSVFromZip(zipReader *zip.Reader, filePassword string) (io.Reader, error) {
	for _, file := range zipReader.File {
		if file.IsEncrypted() {
			file.SetPassword(filePassword)
		}

		f, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file in zip: %w", err)
		}
		return f, nil
	}
	return nil, fmt.Errorf("no file found in zip archive")
}

// ParseFloatField parses a string field, removing commas and converting to float
func ParseFloatField(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// UnmarshalCSV converts a CSV row into an Audit struct
func UnmarshalCSV(rows [][]string) []Audit {
	var audits []Audit
	for _, row := range rows {
		datetime, _ := time.Parse("2006/01/02 3:04 PM", row[0])
		amount := ParseFloatField(row[6])
		accruedStarsCoins := ParseFloatField(row[7])
		tmoney := ParseFloatField(row[8])
		wmoney := ParseFloatField(row[9])
		balance := ParseFloatField(row[10])
		totalStarCoins := ParseFloatField(row[11])
		ttmoney := ParseFloatField(row[12])
		wwmoney := ParseFloatField(row[13])

		audit := Audit{
			DateTime:          datetime,
			Action:            row[1],
			Tournament:        row[2],
			Game:              row[3],
			Currency:          row[5],
			Amount:            amount,
			AccruedStarsCoins: accruedStarsCoins,
			Tmoney:            tmoney,
			Wmoney:            wmoney,
			Balance:           balance,
			TotalStarCoins:    totalStarCoins,
			TTmoney:           ttmoney,
			WWmoney:           wwmoney,
		}
		audits = append(audits, audit)
	}
	return audits
}

func main() {
	fileURL := "http://reports.rationalwebservices.com/reports/xxx.zip"
	filePassword := "12345"

	zipReader, err := DownloadAndUnzipFile(fileURL)
	if err != nil {
		log.Fatalf("Error downloading or unzipping file: %v", err)
	}

	csvFile, err := ExtractCSVFromZip(zipReader, filePassword)
	if err != nil {
		log.Fatalf("Error extracting CSV: %v", err)
	}

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1
	rowsPre, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
	}

	rows := rowsPre[3:]
	audits := UnmarshalCSV(rows)

	var sum float64
	for _, audit := range audits {
		if strings.Contains(audit.Action, "Chest") {
			sum += audit.Amount + (audit.AccruedStarsCoins / 100)
		}
	}

	fmt.Printf("Total sum: %.2f\n", sum)
}
