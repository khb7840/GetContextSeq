package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Variant is a struct which contains the position info and
type Variant struct {
	chrom string
	pos   int64
	ref   string
	alt   string
}

// Concat is a method
func (v Variant) Concat() string {
	var output string
	output = fmt.Sprintf("%s:%d%s>%s", v.chrom, v.pos, v.ref, v.alt)
	return output
}

// ContigIndex is a struct
type ContigIndex struct {
	chrom       string
	length      int64 // number of bases
	start       int64 // byte index of sequence begins
	end         int64 // byte index of sequence ends
	basePerLine int64
	bytePerLine int64
}

// SetEnd is a method
func (ci *ContigIndex) SetEnd() {
	ciLineNum := ci.length / ci.basePerLine
	ciRemainder := ci.length % ci.basePerLine
	ci.end = ci.start + (ciLineNum * ci.bytePerLine) + ciRemainder - 1
}

// SeqOffset is a struct which contains the return value of VariantOffset()
type SeqOffset struct {
	offset  int64
	end     int64
	byteLen int64
	NNstart string
	NNend   string
}

// FaiToIndexMap returns ContigIndex struct from the path of .fai file
func FaiToIndexMap(fai string) (map[string]ContigIndex, []string) {
	// declare contigIndex struct
	var contigIndex ContigIndex
	var fastaIndexKeys []string
	fastaIndexMap := make(map[string]ContigIndex)

	// Read .fai file
	faiFile, err := os.Open(fai)
	if err != nil {
		log.Println(err)
	}
	defer faiFile.Close()

	reader := csv.NewReader(faiFile)
	reader.Comma = '\t' // fai format is tab-delimited

	dat, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// Save to fastaIndexArray
	for _, each := range dat {
		contigIndex.chrom = each[0]
		contigIndex.length, _ = strconv.ParseInt(each[1], 10, 64)
		contigIndex.start, _ = strconv.ParseInt(each[2], 10, 64)
		contigIndex.basePerLine, _ = strconv.ParseInt(each[3], 10, 64)
		contigIndex.bytePerLine, _ = strconv.ParseInt(each[4], 10, 64)
		contigIndex.SetEnd()
		fastaIndexMap[contigIndex.chrom] = contigIndex
		fastaIndexKeys = append(fastaIndexKeys, each[0])
	}
	// Return
	log.Printf("Parsed fasta index: %s\n", fai)
	return fastaIndexMap, fastaIndexKeys
}

// GetContextSeq is a function
func GetContextSeq(fastaFile *os.File, sos SeqOffset) string {
	var contextSeq string
	byteArray := make([]byte, sos.byteLen)
	fastaFile.Seek(sos.offset, 0)
	fastaFile.Read(byteArray)
	contextSeq = string(byteArray)
	contextSeq = strings.Replace(contextSeq, "\n", "", -1)
	contextSeq = strings.ToUpper(contextSeq)
	contextSeq = sos.NNstart + contextSeq + sos.NNend
	return contextSeq
}

// VariantOffset is a function
func VariantOffset(v Variant, faiMap map[string]ContigIndex,
	length int64) SeqOffset {

	output := SeqOffset{}

	// 0 based base position
	baseStart := v.pos - length - 1
	baseEnd := v.pos + length - 1
	// 0 based byte index of contig (chromosome)
	contigStart := faiMap[v.chrom].start // 0 based
	contigEnd := faiMap[v.chrom].end

	// byte index of given sequence in fasta file
	var fastaStart int64
	var fastaEnd int64
	// "NN" repeated string to append
	var NNstart string
	var NNend string
	var Nlen int64

	// consider newline characters for baseStart
	if baseStart < 0 {
		Nlen = -baseStart
		NNstart = strings.Repeat("N", int(Nlen))
		fastaStart = contigStart
	} else {
		startNewlineNum := baseStart / faiMap[v.chrom].basePerLine
		startRemainder := baseStart % faiMap[v.chrom].basePerLine
		fastaStart = (faiMap[v.chrom].bytePerLine * startNewlineNum) +
			startRemainder + contigStart
		NNstart = ""
	}
	// consider newline characters for baseEnd
	endNewlineNum := baseEnd / faiMap[v.chrom].basePerLine
	endRemainder := baseEnd % faiMap[v.chrom].basePerLine
	fastaEnd = (faiMap[v.chrom].bytePerLine * endNewlineNum) +
		endRemainder + contigStart
	if baseEnd >= faiMap[v.chrom].length {
		Nlen := baseEnd - faiMap[v.chrom].length + 1
		NNend = strings.Repeat("N", int(Nlen))
		fastaEnd = contigEnd
	} else {
		NNend = ""
	}

	// set values to output
	output.offset = fastaStart
	output.end = fastaEnd
	output.byteLen = fastaEnd - fastaStart + 1
	output.NNstart = NNstart
	output.NNend = NNend

	return output
}
