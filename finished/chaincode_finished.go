/*
Copyright IBM Corp 2016 All Rights Reserved.

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

import (
	"errors"
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Package struct{
	assetId string `json:"assetId"`
	carrier string `json:"color"`
	temperature string `json:"size"`
	location string `json:"user"`
	datetime string `json:datetime`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("my_package", []byte(args[0]))
	if err != nil {
		return nil, err
	}


	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "init_package" {
		return t.init_package(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query")
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// delete - delte function 
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	var key string
	var err error

	key = args[0]
	err = stub.DelState(key)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Init resets all the things
func (t *SimpleChaincode) init_package(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("5th argument must be a non-empty string")
	}

	assetId := args[0]
	carrier := args[1]
	temp := args[2]
	location := args[3]
	datetime := args[4]

	//check if package already exists
	packageAsBytes, err := stub.GetState(assetId)
	if err != nil {
		return nil, errors.New("Failed to get package name")
	}
	res := Package{}
	json.Unmarshal(packageAsBytes, &res)
	if res.assetId == assetId{
		fmt.Println("This package arleady exists: " + assetId)
		fmt.Println(res);
		return nil, errors.New("This package arleady exists")				//all stop a marble by this name exists
	}

	//build the package json string manually
	str := `{"assetId": "` + assetId + `", "carrier": "` + carrier + `", "temperature": "` + temp + `", "location": "` + location + `", "datetime": "` + datetime + `"}`
	err = stub.PutState(assetId, []byte(str))									//store marble with id as key
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) update_package(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("5th argument must be a non-empty string")
	}

	assetId := args[0]
	carrier := args[1]
	temp := args[2]
	location := args[3]
	datetime := args[4]

	//check if package already exists
	/*packageAsBytes, err := stub.GetState(assetId)
	if err != nil {
		return nil, errors.New("Failed to get package name")
	}
	res := Package{}
	json.Unmarshal(packageAsBytes, &res)
	if res.assetId == assetId{
		fmt.Println("This package arleady exists: " + assetId)
		fmt.Println(res);
		return nil, errors.New("This package arleady exists")				//all stop a marble by this name exists
	}*/

	//build the package json string manually
	str := `{"assetId": "` + assetId + `", "carrier": "` + carrier + `", "temperature": "` + temp + `", "location": "` + location + `", "datetime": "` + datetime + `"}`
	err := stub.PutState(assetId, []byte(str))									//store marble with id as key
	if err != nil {
		return nil, err
	}

	return nil, nil
}