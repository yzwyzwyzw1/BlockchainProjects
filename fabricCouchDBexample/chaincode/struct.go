package main

type BillStruct struct {

	ObjectType string `json:"DocType"`
	BillInfoID string `json:"BillInfoID"`
	BillInfoAmt string `json:"BillInfoAmt"`
	BillInfoType string `json:"BillInfoType"`

	BillIsseDate string `json:"BillIsseDate"`
	BillDueDate string `json:"BillDueDate"`

	HoldrAcct string `json:"HoldrAcct"`
	HoldrCmID string `json:"HoldrCmID"`

	WaitEndorseAcct string `json:"WaitEndorseAcct"`
	WaitEndorseCmID string `json:"WaitEnorseCmID"`
}

