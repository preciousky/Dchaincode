/*
 * this is the chaincode for graducation, the file is coded on windows10, 
 * I will copy it to the specified dir on centos-7 after it is finished.
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 * os is added for init network data by read a json file
 */

import (
	// "bytes"
	"encoding/json"
	"fmt"
	// "strconv"
	// "os"
 
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)
 
// Define the Smart Contract structure
type SmartContract struct {
}
 
// Define the paper structure, with 26 properties. Structure tags are used by encoding/json library
type Paper struct {
    StateTime string `json:"stateTime"`
	PaperId string `json:"paperId"`
	PaperName string `json:"paperName"`
	Value string `json:"value"`
	DDate string `json:"dDate"`
	MDate string `json:"mDate"`
	AcceptDate string `json:"acceptDate"`
	DrawerId string `json:"drawerId"`
	DrawerName string `json:"drawerName"`
	PayerId string `json:"payerId"`
	PayerName string `json:"payerName"`
	PayeeId string `json:"payeeId"`
	PayeeName string `json:"payeeName"`
	HolderId string `json:"holderId"`
	HolderName string `json:"holderName"`
	Rule string `json:"rule"`
	RuleData string `json:"ruleData"`
	RankInfo string `json:"rankInfo"`
	RankerId string `json:"rankerId"`
	RankerName string `json:"rankerName"`
	RankDate string `json:"rankDate"`
	State string `json:"state"`
	StateRoleId string `json:"stateRoleId"`
	StateRoleName string `json:"stateRoleName"`
	CashData string `json:"cashData"`
	NewInfo string `json:"newInfo"`
}

/**
 * The Init method is called when the Smart Contract "DLUTchaincode" is instantiated 
 * by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


/**
 * The Invoke method is called as a result of an application request 
 * to run the Smart Contract "DLUTchaincode"
 * The calling application program has also specified the particular smart contract function to be called, 
 * with arguments.
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryPaper" {
		return s.queryPaper(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryPaper(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	paperAsBytes, _ := APIstub.GetState(args[0])
	// testString := "mytest call queryPaper Function";
	// return shim.Success([]byte(args[0]));
	return shim.Success(paperAsBytes);
}
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	papers := []Paper{
		Paper{StateTime:"0",PaperId:"zzzzz1",PaperName:"票据名称01",Value:"1000",DDate:"2015/9/1",MDate:"2019/6/30",AcceptDate:"2018/9/28",DrawerId:"drawerIdz",DrawerName:"出票方名称",PayerId:"payerIdz",PayerName:"受票方名称",PayeeId:"paperIdz",PayeeName:"收款方名称",HolderId:"holderIdz",HolderName:"持有者名称",Rule:"",RuleData:"",RankInfo:"",RankerId:"rankerIdz",RankerName:"评级方名称",RankDate:"2019/9/1",State:"",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:""},
		Paper{StateTime:"0",PaperId:"zzzzz2",PaperName:"票据名称02",Value:"1001",DDate:"2015/9/1",MDate:"2019/6/30",AcceptDate:"2018/9/28",DrawerId:"drawerIdz",DrawerName:"出票方名称",PayerId:"payerIdz",PayerName:"受票方名称",PayeeId:"paperIdz",PayeeName:"收款方名称",HolderId:"holderIdz",HolderName:"持有者名称",Rule:"",RuleData:"",RankInfo:"",RankerId:"rankerIdz",RankerName:"评级方名称",RankDate:"2019/9/2",State:"",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:""},
		Paper{StateTime:"0",PaperId:"zzzzz3",PaperName:"票据名称03",Value:"1002",DDate:"2015/9/1",MDate:"2019/6/30",AcceptDate:"2018/9/28",DrawerId:"drawerIdz",DrawerName:"出票方名称",PayerId:"payerIdz",PayerName:"受票方名称",PayeeId:"paperIdz",PayeeName:"收款方名称",HolderId:"holderIdz",HolderName:"持有者名称",Rule:"",RuleData:"",RankInfo:"",RankerId:"rankerIdz",RankerName:"评级方名称",RankDate:"2019/9/3",State:"",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:""},
		Paper{StateTime:"0",PaperId:"zzzzz4",PaperName:"票据名称04",Value:"1003",DDate:"2015/9/1",MDate:"2019/6/30",AcceptDate:"2018/9/28",DrawerId:"drawerIdz",DrawerName:"出票方名称",PayerId:"payerIdz",PayerName:"受票方名称",PayeeId:"paperIdz",PayeeName:"收款方名称",HolderId:"holderIdz",HolderName:"持有者名称",Rule:"",RuleData:"",RankInfo:"",RankerId:"rankerIdz",RankerName:"评级方名称",RankDate:"2019/9/4",State:"",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:""},
		Paper{StateTime:"0",PaperId:"zzzzz5",PaperName:"票据名称05",Value:"1004",DDate:"2015/9/1",MDate:"2019/6/30",AcceptDate:"2018/9/28",DrawerId:"drawerIdz",DrawerName:"出票方名称",PayerId:"payerIdz",PayerName:"受票方名称",PayeeId:"paperIdz",PayeeName:"收款方名称",HolderId:"holderIdz",HolderName:"持有者名称",Rule:"",RuleData:"",RankInfo:"",RankerId:"rankerIdz",RankerName:"评级方名称",RankDate:"2019/9/5",State:"",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:""},
		}
	i := 0
	for i < len(papers) {
		fmt.Println("i is ", i)
		paperAsBytes, _ := json.Marshal(papers[i])
		APIstub.PutState(papers[i].PaperId, paperAsBytes)
		fmt.Println("Added", papers[i])
		i = i + 1
	}
	testString := "mytest initLedger Function calling finishes";
	return shim.Success([]byte(testString))
	// return shim.Success(nil)
}
















// func (s *SmartContract) queryCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}

// 	carAsBytes, _ := APIstub.GetState(args[0])
// 	return shim.Success(carAsBytes)
// }



// func (s *SmartContract) createCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 5 {
// 		return shim.Error("Incorrect number of arguments. Expecting 5")
// 	}

// 	var car = Car{Make: args[1], Model: args[2], Colour: args[3], Owner: args[4]}

// 	carAsBytes, _ := json.Marshal(car)
// 	APIstub.PutState(args[0], carAsBytes)

// 	return shim.Success(nil)
// }

// func (s *SmartContract) queryAllCars(APIstub shim.ChaincodeStubInterface) sc.Response {

// 	startKey := "CAR0"
// 	endKey := "CAR999"

// 	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	defer resultsIterator.Close()

// 	// buffer is a JSON array containing QueryResults
// 	var buffer bytes.Buffer
// 	buffer.WriteString("[")

// 	bArrayMemberAlreadyWritten := false
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return shim.Error(err.Error())
// 		}
// 		// Add a comma before array members, suppress it for the first array member
// 		if bArrayMemberAlreadyWritten == true {
// 			buffer.WriteString(",")
// 		}
// 		buffer.WriteString("{\"Key\":")
// 		buffer.WriteString("\"")
// 		buffer.WriteString(queryResponse.Key)
// 		buffer.WriteString("\"")

// 		buffer.WriteString(", \"Record\":")
// 		// Record is a JSON object, so we write as-is
// 		buffer.WriteString(string(queryResponse.Value))
// 		buffer.WriteString("}")
// 		bArrayMemberAlreadyWritten = true
// 	}
// 	buffer.WriteString("]")

// 	fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

// 	return shim.Success(buffer.Bytes())
// }

// func (s *SmartContract) changeCarOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 2 {
// 		return shim.Error("Incorrect number of arguments. Expecting 2")
// 	}

// 	carAsBytes, _ := APIstub.GetState(args[0])
// 	car := Car{}

// 	json.Unmarshal(carAsBytes, &car)
// 	car.Owner = args[1]

// 	carAsBytes, _ = json.Marshal(car)
// 	APIstub.PutState(args[0], carAsBytes)

// 	return shim.Success(nil)
// }

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}