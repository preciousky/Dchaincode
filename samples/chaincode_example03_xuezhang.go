/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		return t.invoke(stub,args)
	}else if function == "calculate"{
		return t.calculate(stub,args)
	}else if function == "query"{
		return t.query(stub,args)
	}else if function == "delete"{
		return t.delete(stub,args)
	}else if function == "register"{
		return t.register(stub,args)
	}else if function == "login"{
		return t.login(stub,args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"\"calculate\"")
}
func (t *SimpleChaincode) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	Account := args[0]
	Key := args[1]
	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(Account)
	if err != nil {
		return shim.Error("DB has error!")
	}
	if Avalbytes == nil {
		E := stub.PutState(Account,[]byte(Key))
		if E != nil{
			return shim.Error("the account onchain failure")
		}
		return shim.Success(nil)
	}
	return shim.Error("The account has been used")

}
func (t *SimpleChaincode) login(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	Account := args[0]
	Key := args[1]
	Avalbytes, err := stub.GetState(Account)
	if err != nil {
		return shim.Error("DB has error!")
	}
	if Avalbytes == nil {
		return shim.Error("This account is not existed ")
	}
	Aval := string(Avalbytes)
	B1 := strings.Fields(Key)
	Aval1 := strings.Fields(Aval)
	for i:=0;i<len(Aval1);i=i+1{
		if Aval1[i] !=B1[i] {
			return shim.Error("The key is wrong!")
		} 
	}
	return shim.Success(nil)
}
// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	A := args[0]
	C := args[1]
	B := args[2]
	D := A + C
	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(D)
	if err != nil {
		return shim.Error("DB has error!")
	}
	if Avalbytes == nil {
		E := stub.PutState(D,[]byte(B))
		if E != nil{
			return shim.Error("the account onchain failure")
		}
		return shim.Success(nil)
	}
	Aval := string(Avalbytes)
	B1 := strings.Fields(B)
	Aval1 := strings.Fields(Aval)
	j := 0
	for i:=0;i<len(Aval1);i=i+3{
		k:=i
		for j=0;j<3;j++{
			if Aval1[k] != B1[j]{
				break
			}
			k=k+1
		}
		if j == 3{
			return shim.Error("the account has been existed")
		}
	}
	Aval = Aval + " " + B
	err = stub.PutState(D, []byte(Aval))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) calculate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 2{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	A1 := args[0]
	A2 := args[1]
	A := A1 + A2
	Avalbytes,err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("no Transaction")
	}
	S1 := string(Avalbytes)
	B1 := strings.Fields(S1)
	Sum1 := 0
	var Tmp1 int 
	for i := 2;i < len(B1);i = i + 3{
		C1 := B1[i]
		Tmp1, _ = strconv.Atoi(C1)
		Sum1 = Sum1 + Tmp1
	} 
	Z := A2 + A1
	Zvalbytes,err := stub.GetState(Z)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Zvalbytes == nil {
		fmt.Printf("XXX")
	}
	S2 := string(Zvalbytes)
	B2 := strings.Fields(S2)
	Sum2 := 0
	var Tmp2 int 
	for i := 2;i < len(B2);i = i + 3{
		C2 := B2[i]
		Tmp2, _ = strconv.Atoi(C2)
		Sum2 = Sum2 + Tmp2
	} 
	Sum := Sum1 - Sum2
	S := "{" + strconv.Itoa(Sum) + "}"
	return shim.Success([]byte(S))
}


// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A1 := args[0]
	A2 := args[1]
	A := A1 + A2

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	B := A2 + A1
	err = stub.DelState(B)
	if err != nil{
		fmt.Printf("123")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A1 := args[0]
	A2 := args[1]
	A := A1 + A2

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Fail to get state")
	}

	if Avalbytes == nil {
		return shim.Error("Nil account")
	}
	Str := "{" + string(Avalbytes) + "}"

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(Str))
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
