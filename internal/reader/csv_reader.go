package reader

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/iambaangkok/Card-Maker/internal/config"
)

const (
	extension = ".csv"
)

type CSVReader interface {
	Read(filePath string) CSVFile
}

type CSVReaderImpl struct {
	Config config.Config
}

func (c CSVReaderImpl) Read(fileName string) CSVFile {
		
	filePath := c.Config.CSV.Dir + "/" + fileName + extension

	log.Println("reading from " + filePath)
    f, err := os.Open(filePath)
    if err != nil {
        log.Println("unable to read input file", err)
		return CSVFile{}
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        log.Println("unable to parse file as CSV", err)
		return CSVFile{}
    }

    return CSVFile{
		Name: fileName,
		Records: records,
	}
}