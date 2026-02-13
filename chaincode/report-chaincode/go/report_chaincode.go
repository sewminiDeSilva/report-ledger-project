package main

import (
"encoding/json"
       "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Report struct defines the data
type Report struct {
    ReportID     string `json:"reportID"`
    Description  string `json:"description"`
    Status       string `json:"status"`
    Reporter     string `json:"reporter"`
    EvidenceHash string `json:"evidenceHash"`
}

// ReportContract defines the chaincode
type ReportContract struct {
    contractapi.Contract
}

// CreateReport - Org1 only
func (rc *ReportContract) CreateReport(ctx contractapi.TransactionContextInterface, reportID, description, reporter, evidenceHash string) error {
    mspID, err := ctx.GetClientIdentity().GetMSPID()
    if err != nil {
        return fmt.Errorf("failed to get MSPID: %v", err)
    }
    if mspID != "Org1MSP" {
        return fmt.Errorf("only Org1 can create reports")
    }

    exists, err := rc.ReportExists(ctx, reportID)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("report %s already exists", reportID)
    }

    report := Report{
        ReportID:     reportID,
        Description:  description,
        Status:       "Pending",
        Reporter:     reporter,
        EvidenceHash: evidenceHash,
    }

    reportJSON, err := json.Marshal(report)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(reportID, reportJSON)
}

// UpdateStatus - Org2 only
func (rc *ReportContract) UpdateStatus(ctx contractapi.TransactionContextInterface, reportID, newStatus string) error {
    mspID, err := ctx.GetClientIdentity().GetMSPID()
    if err != nil {
        return fmt.Errorf("failed to get MSPID: %v", err)
    }
    if mspID != "Org2MSP" {
        return fmt.Errorf("only Org2 can update status")
    }

    reportJSON, err := ctx.GetStub().GetState(reportID)
    if err != nil {
        return err
    }
    if reportJSON == nil {
        return fmt.Errorf("report %s does not exist", reportID)
    }

    var report Report
    err = json.Unmarshal(reportJSON, &report)
    if err != nil {
        return err
    }

    report.Status = newStatus
    updatedJSON, err := json.Marshal(report)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(reportID, updatedJSON)
}

// ReportExists helper
func (rc *ReportContract) ReportExists(ctx contractapi.TransactionContextInterface, reportID string) (bool, error) {
    reportJSON, err := ctx.GetStub().GetState(reportID)
    if err != nil {
        return false, err
    }
    return reportJSON != nil, nil
}

// GetHistory - Audit function
func (rc *ReportContract) GetHistory(ctx contractapi.TransactionContextInterface, reportID string) (string, error) {
    resultsIterator, err := ctx.GetStub().GetHistoryForKey(reportID)
    if err != nil {
        return "", err
    }
    defer resultsIterator.Close()

    var history []map[string]interface{}
    for resultsIterator.HasNext() {
        response, err := resultsIterator.Next()
        if err != nil {
            return "", err
        }
        var record map[string]interface{}
        if response.Value != nil {
            json.Unmarshal(response.Value, &record)
        } else {
            record = map[string]interface{}{"Deleted": true}
        }
        record["TxId"] = response.TxId
        record["Timestamp"] = response.Timestamp
        history = append(history, record)
    }

    historyJSON, err := json.Marshal(history)
    if err != nil {
        return "", err
    }

    return string(historyJSON), nil
}

// Main
func main() {
    chaincode, err := contractapi.NewChaincode(new(ReportContract))
    if err != nil {
        fmt.Printf("Error create chaincode: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting chaincode: %s", err.Error())
    }
}
