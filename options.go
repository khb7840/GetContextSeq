package main

import (
	"fmt"
	"log"
	"os"
)

// Options are options
type Options struct {
	vcf    string
	fasta  string
	index string
	length int64
	output string
	ncpu int
}

func (op *Options) check() error {
	log.Print("Checking parameters")

	fastaIndex := op.fasta + ".fai"
	// set fasta index file in check() step
	op.index = fastaIndex
	_, vcferr := os.Stat(op.vcf)
	_, fastaerr := os.Stat(op.fasta)
	_, indexerr := os.Stat(fastaIndex)

	// check each parameters
	switch {
	case op.vcf == "":
		return fmt.Errorf("No input vcf")
	case os.IsNotExist(vcferr):
		return fmt.Errorf("Invalid input vcf", op.vcf)
	case op.fasta == "":
		return fmt.Errorf("No reference")
	case os.IsNotExist(fastaerr):
		return fmt.Errorf("Invalid fasta file", op.fasta)
	case os.IsNotExist(indexerr):
		return fmt.Errorf("Invalid fasta index", fastaIndex)
	case op.length < 1:
		return fmt.Errorf("Invalid length", op.length)
	case op.output == "":
		return fmt.Errorf("Invalid output path")
	case op.ncpu < 1:
		return fmt.Errorf("Invalid number of cpus", op.ncpu)
	}
	// print parameters
	fmt.Println("NUMPROCS:", op.ncpu)
	fmt.Println("INPUT VCF:", op.vcf)
	fmt.Println("REFERENCE FASTA:", op.fasta)
	fmt.Println("FASTA INDEX", op.index)
	fmt.Println("LENGTH OF CONTEXT:", op.length)
	fmt.Println("OUTPUT:", op.output)
	return nil
}
