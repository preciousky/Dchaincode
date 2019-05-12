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
    "bytes"
    "encoding/json"
    "fmt"
    "strconv"
    "time"
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
    } else if function == "queryPaperById" {
        return s.queryPaperById(APIstub, args)
    } else if function == "queryPapersByHolderId" {
        return s.queryPapersByHolderId(APIstub, args)
    } else if function == "releaseRank" {
		return s.releaseRank(APIstub, args)
	} else if function == "getPaperLogsById" {
        return s.getPaperLogsById(APIstub, args)
    }
    return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryPaperById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }
    paperAsBytes, _ := APIstub.GetState(args[0])
    // testString := "mytest call queryPaper Function";
    // return shim.Success([]byte(args[0]));
    return shim.Success(paperAsBytes);
}
func (s *SmartContract) queryPapersByHolderId(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
     }
    queryString := fmt.Sprintf("{\"selector\":{\"holderId\":\"%s\"}}", args[0])
    resultsIterator,err:= APIstub.GetQueryResult(queryString)
    if err!=nil{
        return shim.Error("query failed")
    }
    papers,err:=getListResult(resultsIterator)
    if err!=nil{
        return shim.Error("query failed")
    }
    return shim.Success(papers)

    //paperAsBytes, _ := APIstub.GetState(args[0])
    // return shim.Success([]byte(args[0]));
    //return shim.Success(paperAsBytes);
}
func (s *SmartContract) releaseRank(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
        return shim.Error("Incorrect number of arguments. Expecting 4")
    }
    paperAsBytes, _ := APIstub.GetState(args[0])
	paper := Paper{}

	json.Unmarshal(paperAsBytes, &paper)
    paper.StateRoleId = args[1];
    paper.StateRoleName = args[2];
    paper.CashData = args[3];
    paper.State = "2";

	paperAsBytes, _ = json.Marshal(paper)
	APIstub.PutState(args[0], paperAsBytes)

	return shim.Success(nil)
	
	// paperId := paperAsBytes["paperId"]
	// data := paperAsBytes
    // APIstub.PutState(paperAsBytes["paperId"],)
    // return shim.Success(paperAsBytes)
}
func (s *SmartContract) getPaperLogsById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }
    // it,err:= APIstub.GetHistoryForKey(args[0])
    // if err!=nil{
    //    return shim.Error(err.Error())
    // }
    // var result,_= getHistoryListResult(it)
    // return shim.Success(result)


    resultsIterator, err := APIstub.GetHistoryForKey(args[0])
    if err != nil {
        return shim.Error(err.Error())
    }
    defer resultsIterator.Close()

    // buffer is a JSON array containing historic values for the marble
    var buffer bytes.Buffer
    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        response, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        // Add a comma before array members, suppress it for the first array member
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"TxId\":")
        buffer.WriteString("\"")
        buffer.WriteString(response.TxId)
        buffer.WriteString("\"")

        buffer.WriteString(", \"Value\":")
        // if it was a delete operation on given key, then we need to set the
        //corresponding value null. Else, we will write the response.Value
        //as-is (as the Value itself a JSON marble)
        if response.IsDelete {
            buffer.WriteString("null")
        } else {
            buffer.WriteString(string(response.Value))
        }

        buffer.WriteString(", \"Timestamp\":")
        buffer.WriteString("\"")
        buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
        buffer.WriteString("\"")

        buffer.WriteString(", \"IsDelete\":")
        buffer.WriteString("\"")
        buffer.WriteString(strconv.FormatBool(response.IsDelete))
        buffer.WriteString("\"")

        buffer.WriteString("}")
        bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")

    fmt.Printf("- getHistoryForMarble returning:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
    papers := []Paper{
		Paper{StateTime:"0",PaperId:"JuhGOzz1",PaperName:"GOAness.Par453-1",Value:"1000",DDate:"2018-9-1",MDate:"2021-9-1",AcceptDate:"",DrawerId:"ZA_user1",DrawerName:"XiXiA_O Co. onene",PayerId:"ZB_user5",PayerName:"Jaga_T fiv",PayeeId:"ZA_user2",PayeeName:"YOUR_O Co. two",HolderId:"ZA_user2",HolderName:"YOUR_O Co. two",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz2",PaperName:"GOAness.Par453-2",Value:"1000",DDate:"2018-9-2",MDate:"2021-9-2",AcceptDate:"",DrawerId:"ZA_user1",DrawerName:"XiXiA_O Co. onene",PayerId:"ZB_user6",PayerName:"SISA_T six",PayeeId:"ZA_user2",PayeeName:"YOUR_O Co. two",HolderId:"ZA_user2",HolderName:"YOUR_O Co. two",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz3",PaperName:"GOAness.Par453-3",Value:"1000",DDate:"2018-9-3",MDate:"2021-9-3",AcceptDate:"",DrawerId:"ZA_user1",DrawerName:"XiXiA_O Co. onene",PayerId:"ZB_user7",PayerName:"COR_T Co. sev",PayeeId:"ZA_user3",PayeeName:"PIZZA_O LT thr",HolderId:"ZA_user3",HolderName:"PIZZA_O LT thr",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz4",PaperName:"GOAness.Par453-4",Value:"1000",DDate:"2018-9-4",MDate:"2021-9-4",AcceptDate:"",DrawerId:"ZA_user1",DrawerName:"XiXiA_O Co. onene",PayerId:"ZB_user8",PayerName:"PASS_T Co. eig",PayeeId:"ZA_user3",PayeeName:"PIZZA_O LT thr",HolderId:"ZA_user3",HolderName:"PIZZA_O LT thr",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz5",PaperName:"GOAness.Par453-5",Value:"1000",DDate:"2018-9-5",MDate:"2021-9-5",AcceptDate:"",DrawerId:"ZA_user2",DrawerName:"YOUR_O Co. twowo",PayerId:"ZB_user5",PayerName:"Jaga_T fiv",PayeeId:"ZA_user3",PayeeName:"PIZZA_O LT thr",HolderId:"ZA_user3",HolderName:"PIZZA_O LT thr",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz6",PaperName:"GOAness.Par453-6",Value:"1000",DDate:"2018-9-6",MDate:"2021-9-6",AcceptDate:"",DrawerId:"ZA_user2",DrawerName:"YOUR_O Co. twowo",PayerId:"ZB_user6",PayerName:"SISA_T six",PayeeId:"ZA_user3",PayeeName:"PIZZA_O LT thr",HolderId:"ZA_user3",HolderName:"PIZZA_O LT thr",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz7",PaperName:"GOAness.Par453-7",Value:"1000",DDate:"2018-9-7",MDate:"2021-9-7",AcceptDate:"",DrawerId:"ZA_user2",DrawerName:"YOUR_O Co. twowo",PayerId:"ZB_user7",PayerName:"COR_T Co. sev",PayeeId:"ZA_user4",PayeeName:"PePa_O fou",HolderId:"ZA_user4",HolderName:"PePa_O fou",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz8",PaperName:"GOAness.Par453-8",Value:"1000",DDate:"2018-9-8",MDate:"2021-9-8",AcceptDate:"",DrawerId:"ZA_user2",DrawerName:"YOUR_O Co. twowo",PayerId:"ZB_user8",PayerName:"PASS_T Co. eig",PayeeId:"ZA_user4",PayeeName:"PePa_O fou",HolderId:"ZA_user4",HolderName:"PePa_O fou",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz9",PaperName:"GOAness.Par453-9",Value:"1000",DDate:"2018-9-9",MDate:"2021-9-9",AcceptDate:"",DrawerId:"ZA_user3",DrawerName:"PIZZA_O LT thr",PayerId:"ZB_user5",PayerName:"Jaga_T fiv",PayeeId:"ZA_user4",PayeeName:"PePa_O fou",HolderId:"ZA_user4",HolderName:"PePa_O fou",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz10",PaperName:"GOAness.Par453-10",Value:"1000",DDate:"2018-9-10",MDate:"2021-9-10",AcceptDate:"",DrawerId:"ZA_user3",DrawerName:"PIZZA_O LT thr",PayerId:"ZB_user6",PayerName:"SISA_T six",PayeeId:"ZA_user4",PayeeName:"PePa_O fou",HolderId:"ZA_user4",HolderName:"PePa_O fou",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz11",PaperName:"GOAness.Par453-11",Value:"1000",DDate:"2018-9-11",MDate:"2021-9-11",AcceptDate:"",DrawerId:"ZA_user3",DrawerName:"PIZZA_O LT thr",PayerId:"ZB_user7",PayerName:"COR_T Co. sev",PayeeId:"ZA_user1",PayeeName:"XiXiA_O Co. one",HolderId:"ZA_user1",HolderName:"XiXiA_O Co. one",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz12",PaperName:"GOAness.Par453-12",Value:"1000",DDate:"2018-9-12",MDate:"2021-9-12",AcceptDate:"",DrawerId:"ZA_user3",DrawerName:"PIZZA_O LT thr",PayerId:"ZB_user8",PayerName:"PASS_T Co. eig",PayeeId:"ZA_user1",PayeeName:"XiXiA_O Co. one",HolderId:"ZA_user1",HolderName:"XiXiA_O Co. one",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz13",PaperName:"GOAness.Par453-13",Value:"1000",DDate:"2018-9-13",MDate:"2021-9-13",AcceptDate:"",DrawerId:"ZA_user4",DrawerName:"PePa_O fou",PayerId:"ZB_user5",PayerName:"Jaga_T fiv",PayeeId:"ZA_user1",PayeeName:"XiXiA_O Co. one",HolderId:"ZA_user1",HolderName:"XiXiA_O Co. one",Rule:"1",RuleData:"",RankInfo:"",RankerId:" ",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz14",PaperName:"GOAness.Par453-14",Value:"1000",DDate:"2018-9-14",MDate:"2021-9-14",AcceptDate:"",DrawerId:"ZA_user4",DrawerName:"PePa_O fou",PayerId:"ZB_user6",PayerName:"SISA_T six",PayeeId:"ZA_user1",PayeeName:"XiXiA_O Co. one",HolderId:"ZA_user1",HolderName:"XiXiA_O Co. one",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz15",PaperName:"GOAness.Par453-15",Value:"1000",DDate:"2018-9-15",MDate:"2021-9-15",AcceptDate:"",DrawerId:"ZA_user4",DrawerName:"PePa_O fou",PayerId:"ZB_user7",PayerName:"COR_T Co. sev",PayeeId:"ZA_user2",PayeeName:"YOUR_O Co. two",HolderId:"ZA_user2",HolderName:"YOUR_O Co. two",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz16",PaperName:"GOAness.Par453-16",Value:"1000",DDate:"2018-9-16",MDate:"2021-9-16",AcceptDate:"",DrawerId:"ZA_user4",DrawerName:"PePa_O fou",PayerId:"ZB_user8",PayerName:"PASS_T Co. eig",PayeeId:"ZA_user2",PayeeName:"YOUR_O Co. two",HolderId:"ZA_user2",HolderName:"YOUR_O Co. two",Rule:"1",RuleData:"",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz17",PaperName:"GOAness.Par453-17",Value:"1000",DDate:"2018-9-17",MDate:"2021-9-17",AcceptDate:"",DrawerId:"ZB_user5",DrawerName:"Jaga_T fiv",PayerId:"ZA_user1",PayerName:"XiXiA_O Co. one",PayeeId:"ZB_user6",PayeeName:"SISA_T six",HolderId:"ZB_user6",HolderName:"SISA_T six",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz18",PaperName:"GOAness.Par453-18",Value:"1000",DDate:"2018-9-18",MDate:"2021-9-18",AcceptDate:"",DrawerId:"ZB_user5",DrawerName:"Jaga_T fiv",PayerId:"ZA_user2",PayerName:"YOUR_O Co. two",PayeeId:"ZB_user6",PayeeName:"SISA_T six",HolderId:"ZB_user6",HolderName:"SISA_T six",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz19",PaperName:"GOAness.Par453-19",Value:"1000",DDate:"2018-9-19",MDate:"2021-9-19",AcceptDate:"",DrawerId:"ZB_user5",DrawerName:"Jaga_T fiv",PayerId:"ZA_user3",PayerName:"PIZZA_O LT thr",PayeeId:"ZB_user7",PayeeName:"COR_T Co. sev",HolderId:"ZB_user7",HolderName:"COR_T Co. sev",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz20",PaperName:"GOAness.Par453-20",Value:"1000",DDate:"2018-9-20",MDate:"2021-9-20",AcceptDate:"",DrawerId:"ZB_user5",DrawerName:"Jaga_T fiv",PayerId:"ZA_user4",PayerName:"PePa_O fou",PayeeId:"ZB_user7",PayeeName:"COR_T Co. sev",HolderId:"ZB_user7",HolderName:"COR_T Co. sev",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz21",PaperName:"GOAness.Par453-21",Value:"1000",DDate:"2018-9-21",MDate:"2021-9-21",AcceptDate:"",DrawerId:"ZB_user6",DrawerName:"SISA_T six",PayerId:"ZA_user1",PayerName:"XiXiA_O Co. one",PayeeId:"ZB_user7",PayeeName:"COR_T Co. sev",HolderId:"ZB_user7",HolderName:"COR_T Co. sev",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz22",PaperName:"GOAness.Par453-22",Value:"1000",DDate:"2018-9-22",MDate:"2021-9-22",AcceptDate:"",DrawerId:"ZB_user6",DrawerName:"SISA_T six",PayerId:"ZA_user2",PayerName:"YOUR_O Co. two",PayeeId:"ZB_user7",PayeeName:"COR_T Co. sev",HolderId:"ZB_user7",HolderName:"COR_T Co. sev",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz23",PaperName:"GOAness.Par453-23",Value:"1000",DDate:"2018-9-23",MDate:"2021-9-23",AcceptDate:"",DrawerId:"ZB_user6",DrawerName:"SISA_T six",PayerId:"ZA_user3",PayerName:"PIZZA_O LT thr",PayeeId:"ZB_user8",PayeeName:"PASS_T Co. eig",HolderId:"ZB_user8",HolderName:"PASS_T Co. eig",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz24",PaperName:"GOAness.Par453-24",Value:"1000",DDate:"2018-9-24",MDate:"2021-9-24",AcceptDate:"",DrawerId:"ZB_user6",DrawerName:"SISA_T six",PayerId:"ZA_user4",PayerName:"PePa_O fou",PayeeId:"ZB_user8",PayeeName:"PASS_T Co. eig",HolderId:"ZB_user8",HolderName:"PASS_T Co. eig",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz25",PaperName:"GOAness.Par453-25",Value:"1000",DDate:"2018-9-25",MDate:"2021-9-25",AcceptDate:"",DrawerId:"ZB_user7",DrawerName:"COR_T Co. sev",PayerId:"ZA_user1",PayerName:"XiXiA_O Co. one",PayeeId:"ZB_user8",PayeeName:"PASS_T Co. eig",HolderId:"ZB_user8",HolderName:"PASS_T Co. eig",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz26",PaperName:"GOAness.Par453-26",Value:"1000",DDate:"2018-9-26",MDate:"2021-9-26",AcceptDate:"",DrawerId:"ZB_user7",DrawerName:"COR_T Co. sev",PayerId:"ZA_user2",PayerName:"YOUR_O Co. two",PayeeId:"ZB_user8",PayeeName:"PASS_T Co. eig",HolderId:"ZB_user8",HolderName:"PASS_T Co. eig",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz27",PaperName:"GOAness.Par453-27",Value:"1000",DDate:"2018-9-27",MDate:"2021-9-27",AcceptDate:"",DrawerId:"ZB_user7",DrawerName:"COR_T Co. sev",PayerId:"ZA_user3",PayerName:"PIZZA_O LT thr",PayeeId:"ZB_user5",PayeeName:"Jaga_T fiv",HolderId:"ZB_user5",HolderName:"Jaga_T fiv",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz28",PaperName:"GOAness.Par453-28",Value:"1000",DDate:"2018-9-28",MDate:"2021-9-28",AcceptDate:"",DrawerId:"ZB_user7",DrawerName:"COR_T Co. sev",PayerId:"ZA_user4",PayerName:"PePa_O fou",PayeeId:"ZB_user5",PayeeName:"Jaga_T fiv",HolderId:"ZB_user5",HolderName:"Jaga_T fiv",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz29",PaperName:"GOAness.Par453-29",Value:"1000",DDate:"2018-9-29",MDate:"2021-9-29",AcceptDate:"",DrawerId:"ZB_user8",DrawerName:"PASS_T Co. eig",PayerId:"ZA_user1",PayerName:"XiXiA_O Co. one",PayeeId:"ZB_user5",PayeeName:"Jaga_T fiv",HolderId:"ZB_user5",HolderName:"Jaga_T fiv",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz30",PaperName:"GOAness.Par453-30",Value:"1000",DDate:"2018-9-30",MDate:"2021-9-30",AcceptDate:"",DrawerId:"ZB_user8",DrawerName:"PASS_T Co. eig",PayerId:"ZA_user2",PayerName:"YOUR_O Co. two",PayeeId:"ZB_user5",PayeeName:"Jaga_T fiv",HolderId:"ZB_user5",HolderName:"Jaga_T fiv",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz31",PaperName:"GOAness.Par453-31",Value:"1000",DDate:"2018-10-1",MDate:"2021-10-1",AcceptDate:"",DrawerId:"ZB_user8",DrawerName:"PASS_T Co. eig",PayerId:"ZA_user3",PayerName:"PIZZA_O LT thr",PayeeId:"ZB_user6",PayeeName:"SISA_T six",HolderId:"ZB_user6",HolderName:"SISA_T six",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		Paper{StateTime:"0",PaperId:"JuhGOzz32",PaperName:"GOAness.Par453-32",Value:"1000",DDate:"2018-10-2",MDate:"2021-10-2",AcceptDate:"",DrawerId:"ZB_user8",DrawerName:"PASS_T Co. eig",PayerId:"ZA_user4",PayerName:"PePa_O fou",PayeeId:"ZB_user6",PayeeName:"SISA_T six",HolderId:"ZB_user6",HolderName:"SISA_T six",Rule:"2",RuleData:"30",RankInfo:"",RankerId:"",RankerName:"",RankDate:"",State:"0",StateRoleId:"",StateRoleName:"",CashData:"",NewInfo:"1"},
		}
    i := 0
    for i < len(papers) {
        fmt.Println("i is ", i)
        paperAsBytes, _ := json.Marshal(papers[i])
        APIstub.PutState(papers[i].PaperId, paperAsBytes)
        fmt.Println("Added", papers[i])
        i = i + 1
    }
    testString := "DLUT initLedger Function calling finishes";
    return shim.Success([]byte(testString))
    // return shim.Success(nil)
}
func getListResult(resultsIterator shim.StateQueryIteratorInterface) ([]byte,error){

	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")
 
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
	   queryResponse, err := resultsIterator.Next()
	   if err != nil {
		  return nil, err
	   }
	   // Add a comma before array members, suppress it for the first array member
	   if bArrayMemberAlreadyWritten == true {
		  buffer.WriteString(",")
	   }
	   buffer.WriteString("{\"Key\":")
	   buffer.WriteString("\"")
	   buffer.WriteString(queryResponse.Key)
	   buffer.WriteString("\"")
 
	   buffer.WriteString(", \"Record\":")
	   // Record is a JSON object, so we write as-is
	   buffer.WriteString(string(queryResponse.Value))
	   buffer.WriteString("}")
	   bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
 }

 func main() {
    // Create a new Smart Contract
    err := shim.Start(new(SmartContract))
    if err != nil {
        fmt.Printf("Error creating new Smart Contract: %s", err)
    }
}

 // func (s *SmartContract) queryCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

//     if len(args) != 1 {
//         return shim.Error("Incorrect number of arguments. Expecting 1")
//     }

//     carAsBytes, _ := APIstub.GetState(args[0])
//     return shim.Success(carAsBytes)
// }



// func (s *SmartContract) createCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

//     if len(args) != 5 {
//         return shim.Error("Incorrect number of arguments. Expecting 5")
//     }

//     var car = Car{Make: args[1], Model: args[2], Colour: args[3], Owner: args[4]}

//     carAsBytes, _ := json.Marshal(car)
//     APIstub.PutState(args[0], carAsBytes)

//     return shim.Success(nil)
// }

// func (s *SmartContract) queryAllCars(APIstub shim.ChaincodeStubInterface) sc.Response {

//     startKey := "CAR0"
//     endKey := "CAR999"

//     resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
//     if err != nil {
//         return shim.Error(err.Error())
//     }
//     defer resultsIterator.Close()

//     // buffer is a JSON array containing QueryResults
//     var buffer bytes.Buffer
//     buffer.WriteString("[")

//     bArrayMemberAlreadyWritten := false
//     for resultsIterator.HasNext() {
//         queryResponse, err := resultsIterator.Next()
//         if err != nil {
//             return shim.Error(err.Error())
//         }
//         // Add a comma before array members, suppress it for the first array member
//         if bArrayMemberAlreadyWritten == true {
//             buffer.WriteString(",")
//         }
//         buffer.WriteString("{\"Key\":")
//         buffer.WriteString("\"")
//         buffer.WriteString(queryResponse.Key)
//         buffer.WriteString("\"")

//         buffer.WriteString(", \"Record\":")
//         // Record is a JSON object, so we write as-is
//         buffer.WriteString(string(queryResponse.Value))
//         buffer.WriteString("}")
//         bArrayMemberAlreadyWritten = true
//     }
//     buffer.WriteString("]")

//     fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

//     return shim.Success(buffer.Bytes())
// }

// func (s *SmartContract) changeCarOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

//     if len(args) != 2 {
//         return shim.Error("Incorrect number of arguments. Expecting 2")
//     }

//     carAsBytes, _ := APIstub.GetState(args[0])
//     car := Car{}

//     json.Unmarshal(carAsBytes, &car)
//     car.Owner = args[1]

//     carAsBytes, _ = json.Marshal(car)
//     APIstub.PutState(args[0], carAsBytes)

//     return shim.Success(nil)
// }

// The main function is only relevant in unit test mode. Only included here for completeness.
