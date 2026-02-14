
## 1. System Overview

This project implements a permissioned blockchain-based reporting system using Hyperledger Fabric.

The system demonstrates how role-based access control, immutability, auditability, and evidence integrity can be enforced using smart contracts.

Only report metadata and evidence hashes are stored on-chain. Actual evidence files remain encrypted in off-chain storage.

---

## 2. Architecture Roles

The system models four organizational roles:

- Org1 – Citizen / Report Creator
- Org2 – Ministry (Central Authority)
- Org3 – Investigation Agency
- Org4 – Auditor

Each organization has restricted permissions enforced by MSP identity validation inside the smart contract.

---

## 3. Demonstration Scenario (End-to-End Flow)

This integrated demo demonstrates the full lifecycle of a report:

### Step 1 – Report Submission (Org1)

- Org1 submits a report.
- Evidence is stored off-chain.
- SHA-256 hash of evidence is stored on-chain.
- Status is set to "Submitted".

Security Property Demonstrated:
- Role-based creation
- Evidence integrity via hashing

---

### Step 2 – Report Assignment (Org2)

- Ministry (Org2) assigns the report to an investigation agency.
- Status changes to "Assigned".

Security Property Demonstrated:
- Separation of duties
- Controlled workflow governance

---

### Step 3 – Investigation Update (Org3)

- Assigned agency updates investigation status.
- Status changes to "Under Investigation" or "Completed".

Security Property Demonstrated:
- Restricted modification rights
- Authorized state transition enforcement

---

### Step 4 – Audit History Retrieval (Org4)

- Auditor retrieves full transaction history using GetHistory.
- All changes are shown with TxID and timestamp.

Security Property Demonstrated:
- Immutability
- Full auditability
- Tamper-evidence

---

## 4. Security Alignment

This system demonstrates the following security guarantees:

- Integrity:
  Reports cannot be deleted or overwritten.
  Evidence hash ensures tamper detection.

- Authorization:
  Only specific organizations can perform specific operations.
  MSP identity checks enforce role restrictions.

- Auditability:
  Full transaction history is permanently recorded.
  No silent deletion is possible.

- Confidentiality:
  Sensitive evidence is stored off-chain.
  Only metadata and hashes are on-chain.

---

## 5. Proof Mapping

| Proof | What It Demonstrates |
|-------|----------------------|
| Proof 1 | Immutable report lifecycle |
| Proof 2 | Evidence tampering detection |
| Proof 3 | Endorsement policy enforcement |
| Proof 4 | Integrated multi-role workflow |

---

## 6. Conclusion

This integrated demonstration validates that a permissioned blockchain network can enforce governance, integrity, and accountability in a multi-organizational reporting system.
