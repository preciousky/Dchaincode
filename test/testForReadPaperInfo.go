package main

import (
    "encoding/json"
    "fmt"
    "os"
)

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

func main() {
    //writeFile()
    readFile()
}

func readFile() {

    filePtr, err := os.Open("paper_init_data.json")
    if err != nil {
        fmt.Println("Open file failed [Err:%s]", err.Error())
        //return shim.Error("Open file failed.")
    }
    defer filePtr.Close()

    var papers []Paper

    // 创建json解码器
    decoder := json.NewDecoder(filePtr)
    err = decoder.Decode(&papers)
    if err != nil {
        fmt.Println("init_paper_data_decoder_failed", err.Error())

    } else {
        fmt.Println("Decoder success")
        i := 0
        for i <len(papers){
            fmt.Println("added",papers[i].PaperName)
            i += 1
        }
    }
}


/* func writeFile() {
    personInfo := []PersonInfo{{"David", 30, true, []string{"跑步", "读书", "看电影"}}, {"Lee", 27, false, []string{"工作", "读书", "看电影"}}}

    // 创建文件
    filePtr, err := os.Create("person_info_w_test.json")
    if err != nil {
        fmt.Println("Create file failed", err.Error())
        return
    }
    defer filePtr.Close()

    // 创建Json编码器
    encoder := json.NewEncoder(filePtr)

    err = encoder.Encode(personInfo)
    if err != nil {
        fmt.Println("Encoder failed", err.Error())

    } else {
        fmt.Println("Encoder success")
    }


   // 带JSON缩进格式写文件
   //data, err := json.MarshalIndent(personInfo, "", "  ")
   //if err != nil {
   // fmt.Println("Encoder failed", err.Error())
   //
   //} else {
   // fmt.Println("Encoder success")
   //}
   //
   //filePtr.Write(data)
} */