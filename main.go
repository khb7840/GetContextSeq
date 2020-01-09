// getcontext_toy.go
/* get context sequence of given variants
   Author: Hyunbin Kim (khb7840@gmail.com)
   Modified: 2020. 01. 09 */

package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"sync"
)

var param = new(Options)

func init() {
	// Parsing flags in init() function.
	// init() automatically works before main()
	flag.StringVar(
		&param.vcf, "vcf", "", "Input file in vcf format",
	)
	flag.StringVar(
		&param.vcf, "v", "", "Input file in vcf format",
	)
	flag.StringVar(
		&param.fasta, "fasta", "",
		"Reference genome in fasta format",
	)
	flag.StringVar(
		&param.fasta, "f", "",
		"Reference genome in fasta format",
	)
	flag.Int64Var(
		&param.length, "length", 1,
		"Length of context sequence in integer",
	)
	flag.Int64Var(
		&param.length, "l", 1,
		"Length of context sequence in integer",
	)
	flag.StringVar(
		&param.output, "output", "output.vcf", "Output",
	)
	flag.StringVar(
		&param.output, "out", "output.vcf", "Output",
	)
	flag.StringVar(
		&param.output, "o", "output.vcf", "Output",
	)
	flag.IntVar(
		&param.ncpu, "ncpu", 1, "Number of CPUs",
	)
	flag.IntVar(
		&param.ncpu, "n", 1, "Number of CPUs",
	)
}

func main() {
	// 01. Initializing
	log.Print("Initializing")
	flag.Parse()
	if param.ncpu <= runtime.NumCPU() {
		runtime.GOMAXPROCS(param.ncpu)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	var wg sync.WaitGroup
	// 02. Check parameters
	paramError := param.check()
	// if parameters are not given, print error log and exit
	if paramError != nil {
		log.Fatal(paramError)
	}
	// 03. Read input files
	// 03-1 fai file
	log.Print("Parsing .fai file")
	indexMap, _ := FaiToIndexMap(param.index)
	// 03-2 fasta file
	log.Print("Reading fasta file")
	fastaFile, _ := os.Open(param.fasta)
	// 03-3 vcf file
	log.Print("Loading VCF")
	variantArray, variantRecords := LoadVCF(param.vcf)

	//04. Calculate byte index for each variant
	// empty array to save SeqOffsets
	offsetArray := make([]SeqOffset, len(variantArray))
	contextSeqArray := make([]string, len(variantArray))
	// iteration with goroutine
	for i, variant := range variantArray {
		wg.Add(1)
		go func(variant Variant, i int) {
			defer wg.Done()
			offsetArray[i] = VariantOffset(variant, indexMap, param.length)
		}(variant, i)
	}
	wg.Wait()

	//05. Get context sequence with byte indices
	log.Println("Getting context sequence")
	for i, offset := range offsetArray {
		contextSeqArray[i] = GetContextSeq(fastaFile, offset)
	}
	WriteVCFWithContext(variantRecords, contextSeqArray, param.output)
	log.Println("Saved:", param.output)
	fastaFile.Close()
}
