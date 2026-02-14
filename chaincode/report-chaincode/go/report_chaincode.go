package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Report struct {
	ReportID     string `json:"reportID"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	Reporter     string `json:"reporter"`
	AssignedTo   string `json:"assignedTo"`
	EvidenceHash string `json:"evidenceHash"`
}

type ReportContract struct {
	contractapi.Contract
}

// ----------------------
// ROLE 1: CITIZEN (Org1)
// ----------------------
func (rc *ReportContract) CreateReport(ctx contractapi.TransactionContextInterface,
	reportID, description, reporter, evidenceHash string) error {

	mspID, _ := ctx.GetClientIdentity().GetMSPID()
	if mspID != "Org1MSP" {
		return fmt.Errorf("only Org1 (Citizen) can create reports")
	}

	exists, _ := rc.ReportExists(ctx, reportID)
	if exists {
		return fmt.Errorf("report already exists")
	}

	report := Report{
		ReportID:     reportID,
		Description:  description,
		Status:       "Submitted",
		Reporter:     reporter,
		AssignedTo:   "",
		EvidenceHash: evidenceHash,
	}

	reportJSON, _ := json.Marshal(report)
	return ctx.GetStub().PutState(reportID, reportJSON)
}

// ----------------------
// ROLE 2: MINISTRY (Org2)
// ----------------------
func (rc *ReportContract) AssignAgency(ctx contractapi.TransactionContextInterface,
	reportID, agency string) error {

	mspID, _ := ctx.GetClientIdentity().GetMSPID()
	if mspID != "Org2MSP" {
		return fmt.Errorf("only Ministry (Org2) can assign agencies")
	}

	reportJSON, _ := ctx.GetStub().GetState(reportID)
	if reportJSON == nil {
		return fmt.Errorf("report does not exist")
	}

	var report Report
	json.Unmarshal(reportJSON, &report)

	report.AssignedTo = agency
	report.Status = "Assigned"

	updatedJSON, _ := json.Marshal(report)
	return ctx.GetStub().PutState(reportID, updatedJSON)
}

// ----------------------
// ROLE 3: AGENCY
// ----------------------
func (rc *ReportContract) UpdateStatus(ctx contractapi.TransactionContextInterface,
	reportID, newStatus string) error {

	reportJSON, _ := ctx.GetStub().GetState(reportID)
	if reportJSON == nil {
		return fmt.Errorf("report does not exist")
	}

	var report Report
	json.Unmarshal(reportJSON, &report)

	report.Status = newStatus

	updatedJSON, _ := json.Marshal(report)
	return ctx.GetStub().PutState(reportID, updatedJSON)
}

// ----------------------
// ROLE 4: AUDITOR
// ----------------------
func (rc *ReportContract) GetHistory(ctx contractapi.TransactionContextInterface,
	reportID string) (string, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(reportID)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var history []map[string]interface{}

	for resultsIterator.HasNext() {
		response, _ := resultsIterator.Next()

		var record map[string]interface{}
		if response.Value != nil {
			json.Unmarshal(response.Value, &record)
		}

		record["TxId"] = response.TxId
		record["Timestamp"] = response.Timestamp

		history = append(history, record)
	}

	historyJSON, _ := json.Marshal(history)
	return string(historyJSON), nil
}

func (rc *ReportContract) ReportExists(ctx contractapi.TransactionContextInterface,
	reportID string) (bool, error) {

	reportJSON, _ := ctx.GetStub().GetState(reportID)
	return reportJSON != nil, nil
}

func main() {
	chaincode, _ := contractapi.NewChaincode(new(ReportContract))
	chaincode.Start()
}
