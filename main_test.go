package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexmullins/zip"
	"github.com/stretchr/testify/assert"
)

// Mock zip file content for testing
func createMockZipServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)

		// Add a file to the zip
		f, _ := zipWriter.Create("test.csv")
		f.Write([]byte("Date,Action,Tournament\n2024/01/01 1:00 PM,Chest,Test Tournament"))

		zipWriter.Close()
		w.Header().Set("Content-Type", "application/zip")
		w.Write(buf.Bytes())
	})
	return httptest.NewServer(handler)
}

func TestDownloadAndUnzipFile_Success(t *testing.T) {
	// --- Given ---
	mockServer := createMockZipServer()
	defer mockServer.Close()

	// --- When ---
	zipReader, err := DownloadAndUnzipFile(mockServer.URL)

	// --- Then ---
	assert.NoError(t, err, "Expected no error while downloading and unzipping file")
	assert.NotNil(t, zipReader, "Expected zipReader to not be nil")
}

func TestDownloadAndUnzipFile_InvalidURL(t *testing.T) {
	// --- Given ---
	invalidURL := "http://invalid.url/doesnotexist.zip"

	// --- When ---
	_, err := DownloadAndUnzipFile(invalidURL)

	// --- Then ---
	assert.Error(t, err, "Expected error for invalid URL")
}

func TestExtractCSVFromZip_Success(t *testing.T) {
	// --- Given ---
	mockServer := createMockZipServer()
	defer mockServer.Close()

	zipReader, _ := DownloadAndUnzipFile(mockServer.URL)
	filePassword := ""

	// --- When ---
	csvFile, err := ExtractCSVFromZip(zipReader, filePassword)

	// --- Then ---
	assert.NoError(t, err, "Expected no error while extracting CSV")
	assert.NotNil(t, csvFile, "Expected csvFile to not be nil")
}

func TestExtractCSVFromZip_NoFilesInZip(t *testing.T) {
	// --- Given ---
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zipWriter := zip.NewWriter(w)
		zipWriter.Close()
	})
	mockServer := httptest.NewServer(handler)
	defer mockServer.Close()

	zipReader, _ := DownloadAndUnzipFile(mockServer.URL)
	filePassword := ""

	// --- When ---
	csvFile, err := ExtractCSVFromZip(zipReader, filePassword)

	// --- Then ---
	assert.Error(t, err, "Expected error when no files found in zip")
	assert.Nil(t, csvFile, "Expected csvFile to be nil")
}

func TestParseFloatField_ValidString(t *testing.T) {
	// --- Given ---
	input := "1,234.56"

	// --- When ---
	result := ParseFloatField(input)

	// --- Then ---
	assert.Equal(t, 1234.56, result, "Expected to parse float correctly")
}

func TestParseFloatField_InvalidString(t *testing.T) {
	// --- Given ---
	input := "invalid"

	// --- When ---
	result := ParseFloatField(input)

	// --- Then ---
	assert.Equal(t, 0.0, result, "Expected 0.0 for invalid float string")
}

func TestUnmarshalCSV_ValidData(t *testing.T) {
	// --- Given ---
	rows := [][]string{
		{"2024/01/01 1:00 PM", "Chest", "Test Tournament", "Game1", "", "USD", "100", "200", "10", "5", "500", "1000", "50", "25"},
	}

	// --- When ---
	audits := UnmarshalCSV(rows)

	// --- Then ---
	assert.Len(t, audits, 1, "Expected one audit record")
	assert.Equal(t, "Chest", audits[0].Action, "Expected correct action")
	assert.Equal(t, 100.0, audits[0].Amount, "Expected correct amount")
	assert.Equal(t, 200.0, audits[0].AccruedStarsCoins, "Expected correct stars coins")
}
