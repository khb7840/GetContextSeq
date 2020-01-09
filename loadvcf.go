package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

//LoadVCF loads vcf and returns variant array and records of vcf
func LoadVCF(vcfPath string) ([]Variant, [][]string) {
	// read vcf file
	vcfFile, err := os.Open(vcfPath)
	if err != nil {
		log.Println(err)
	}
	defer vcfFile.Close()
	reader := csv.NewReader(vcfFile)
	reader.Comma = '\t' // vcf format is tab-delimited
	reader.Comment = '#'
	vcfRecords, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// Extract chrom/pos/ref/alt
	output := make([]Variant, len(vcfRecords))
	for i, v := range vcfRecords {
		output[i].chrom = v[0]
		output[i].pos, _ = strconv.ParseInt(v[1], 10, 64)
		output[i].ref = v[3]
		output[i].alt = v[4]
	}
	return output, vcfRecords
}

//WriteVCFWithContext is a function
func WriteVCFWithContext(vcfRec [][]string, seqArray []string, outPath string) {
	// Create output file
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Println(err)
	}
	defer outFile.Close()
	// New CSV writer
	writer := csv.NewWriter(outFile)
	writer.Comma = '\t'
	defer writer.Flush()

	lineSlice := []string{}

	for i, val := range vcfRec {
		lineSlice = append(val, seqArray[i])
		err := writer.Write(lineSlice)
		if err != nil {
			log.Println(err)
		}
	}
}
