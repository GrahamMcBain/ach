package main

import (
	"log"
	"os"
	"time"

	"github.com/moov-io/ach"
)

func main() {
	// Example transfer to write an ACH MTE file acknowledging a credit
	// Important: All financial institutions are different and will require registration and exact field values.

	// Set originator bank ODFI and destination Operator for the financial institution
	// this is the funding/receiving source of the transfer
	fh := ach.NewFileHeader()
	fh.ImmediateDestination = "031300012" // Routing Number of the ACH Operator or receiving point to which the file is being sent
	fh.ImmediateOrigin = "231380104"      // Routing Number of the ACH Operator or sending point that is sending the file
	fh.FileCreationDate = time.Now()      // Today's Date
	fh.ImmediateDestinationName = "Federal Reserve Bank"
	fh.ImmediateOriginName = "My Bank Name"

	// BatchHeader identifies the originating entity and the type of transactions contained in the batch
	bh := ach.NewBatchHeader()
	bh.ServiceClassCode = 225            // ACH credit pushes money out, 225 debits/pulls money in.
	bh.CompanyName = "Merchant with ATM" // Merchant with the ATM
	bh.CompanyIdentification = fh.ImmediateOrigin
	bh.StandardEntryClassCode = "MTE"
	bh.CompanyEntryDescription = "CASH WITHDRAW" // will be on receiving accounts statement // TODO
	bh.EffectiveEntryDate = time.Now()           // Date physical money was received
	bh.ODFIIdentification = "23138010"           // Originating Routing Number

	// Identifies the receivers account information
	// can be multiple entry's per batch
	entry := ach.NewEntryDetail()
	// Identifies the entry as a debit and credit entry AND to what type of account (Savings, DDA, Loan, GL)
	entry.TransactionCode = 27
	entry.SetRDFI("031300012")             // Receivers bank transit routing number
	entry.DFIAccountNumber = "744-5678-99" // Receivers bank account number
	entry.Amount = 10000                   // Amount of transaction with no decimal. One dollar and eleven cents = 111
	entry.SetOriginalTraceNumber("031300010000001")
	entry.SetReceivingCompany("JANE DOE")
	entry.SetTraceNumber(bh.ODFIIdentification, 1)

	addenda02 := ach.NewAddenda02()
	// NACHA rules example: 200509*321 East Market Street*Anytown*VA\
	addenda02.TerminalIdentificationCode = "200509"
	addenda02.TerminalLocation = "321 East Market Street"
	addenda02.TerminalCity = "ANYTOWN"
	addenda02.TerminalState = "VA" // TODO(adam): validate?

	addenda02.TransactionSerialNumber = "123456" // Generated by Terminal, used for audits
	addenda02.TransactionDate = "1224"
	addenda02.TraceNumber = entry.TraceNumber
	entry.Addenda02 = addenda02
	entry.AddendaRecordIndicator = 1

	// build the batch
	batch := ach.NewBatchMTE(bh)
	batch.AddEntry(entry)
	if err := batch.Create(); err != nil {
		log.Fatalf("Unexpected error building batch: %s\n", err)
	}

	// build the file
	file := ach.NewFile()
	file.SetHeader(fh)
	file.AddBatch(batch)
	if err := file.Create(); err != nil {
		log.Fatalf("Unexpected error building file: %s\n", err)
	}

	// write the file to std out. Anything io.Writer
	w := ach.NewWriter(os.Stdout)
	if err := w.Write(file); err != nil {
		log.Fatalf("Unexpected error: %s\n", err)
	}
	w.Flush()
}
