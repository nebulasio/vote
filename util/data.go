package util

import (
    "encoding/json"
    "fmt"
)

func writeDataToFile(r string, outputPath string) {
    if err := WriteFile(outputPath, r); err == nil {
        Print("success. ->", outputPath)
    } else {
        PrintError(err)
    }
}

func CreateContractData(data interface{}) map[string]interface{} {
    d, err := json.Marshal(data)
    if err != nil {
        PrintError(err)
    }
    return map[string]interface{}{"data": string(d)}
}

func SerializeDataToFile(data interface{}, outputPath string) {
    m := CreateContractData(data)
    SerializeDataListToFile([]map[string]interface{}{m}, outputPath)
}

func SerializeDataListToFile(data []map[string]interface{}, outputPath string) {
    r, err := json.Marshal(data)
    if err != nil {
        PrintError(err)
    }
    writeDataToFile(string(r), outputPath)
}

func DeserializeDataFile(file string) ([]map[string]interface{}, error) {
    text, err := ReadFile(file)
    if err != nil {
        return nil, err
    }

    var r []map[string]interface{}
    err = json.Unmarshal([]byte(text), &r)
    if err != nil {
        return nil, err
    }
    return r, nil
}

func DeserializeData(data string) map[string]interface{} {
    var r map[string]interface{}
    err := json.Unmarshal([]byte(data), &r)
    if err != nil {
        PrintError("Data error. ")
    }
    return r
}

func VerifyData(data map[string]interface{}) (string, string) {
    action, ok := data["action"]
    if !ok {
        PrintError("Data error.")
    }

    detail, ok := data["detail"]
    if !ok {
        PrintError("Data error.")
    }

    v := ""
    switch action {
    case ActionRemoveSignee:
        VerifyAddress(detail.(string))

    case ActionAddSignee:
        VerifyAddress(detail.(string))

    case ActionReplaceSignee:
        verifyReplaceManagerData(detail.(map[string]interface{}))

    case ActionUpdateRules:
        VerifySendRules(detail.(map[string]interface{}))

    case ActionUpdateConstitution:
        VerifyConstitution(detail.(map[string]interface{}))

    case ActionSend:
        v = verifySendNasData(detail.(map[string]interface{}))

    case ActionVote:
        v = verifyVoteData(detail.(map[string]interface{}))

    default:
        PrintError("Action", action, "is not supported.")
    }
    return action.(string), v
}

func verifyReplaceManagerData(data map[string]interface{}) {
    oldAddress, ok := data["oldAddress"]
    if !ok {
        PrintError("oldAddress is empty. ")
    }

    newAddress, ok := data["newAddress"]
    if !ok {
        PrintError("newAddress is empty. ")
    }

    VerifyAddress(oldAddress.(string))
    VerifyAddress(newAddress.(string))
    if oldAddress == newAddress {
        PrintError("Data error. ")
    }
}

func verifySendNasData(item map[string]interface{}) string {
    id, ok := item["id"]
    if !ok || IsEmptyString(id.(string)) {
        PrintError("tx.id is empty. ")
    }
    to, ok := item["to"]
    if !ok {
        PrintError("tx.to is empty. ")
    }
    VerifyAddress(to.(string))
    value, ok := item["value"]
    if !ok {
        PrintError("tx.value is empty. ")
    }
    VerifyNumber(value.(string))
    return id.(string)
}

func verifyVoteData(item map[string]interface{}) string {
    id, ok := item["id"]
    if !ok || IsEmptyString(id.(string)) {
        PrintError("vote.id is empty. ")
    }
    _, ok = item["content"]
    if !ok {
        PrintError("vote.content is empty. ")
    }

    action, ok := item["approvedAction"]
    if ok {
        verifyVoteApprovedAction(action.(map[string]interface{}))
    }

    p, ok := item["proportionOfApproved"]
    if !ok {
        PrintError("vote.proportionOfApproved is empty. ")
    }
    VerifyProportions(p.(string))

    return id.(string)
}

func verifyVoteApprovedAction(action map[string]interface{}) {
    // TODO:
}

func VerifyConstitution(data map[string]interface{}) {
    ver, ok := data["version"]
    if !ok {
        PrintError("version is empty. ")
    }
    VerifyNumber(ver.(string))

    t, ok := data["proportionOfSigners"]
    if !ok {
        PrintError("proportionOfSigners is empty. ")
    }

    p := t.(map[string]interface{})
    ks := []interface{}{"updateConstitution", "updateSendRules", "addSignee", "removeSignee", "replaceSignee", "vote"}
    n := 0
    for k, v := range p {
        if Contains(ks, k) {
            n++
        }
        VerifyProportions(v.(string))
    }
    if n != len(ks) {
        PrintError("Constitution data error. ")
    }
}

func VerifySendRules(data map[string]interface{}) {
    ver, ok := data["version"]
    if !ok {
        PrintError("version is empty. ")
    }
    VerifyNumber(ver.(string))

    t, ok := data["rules"]
    if !ok {
        PrintError("rules is empty. ")
    }

    rules := t.([]interface{})
    if len(rules) <= 0 {
        PrintError("rules is empty. ")
    }

    v := 0.0
    for _, i := range rules {
        r := i.(map[string]interface{})
        p, ok := r["proportionOfSigners"]
        if !ok {
            PrintError("proportionOfSigners is empty. ")
        }
        VerifyProportions(p.(string))

        t, ok := r["startValue"]
        if !ok {
            PrintError("startValue is empty. ")
        }
        startValue := ParseFloat(t.(string))
        if v == -1 || startValue != v {
            PrintError("Rules error. ", startValue, v)
        }

        e, ok := r["endValue"]
        if !ok {
            PrintError("endValue is empty. ")
        }
        if e != Infinity {
            v = ParseFloat(e.(string))
            if startValue >= v {
                PrintError("Rules error. ")
            }
        } else {
            v = -1
        }
    }
    if v != -1 {
        PrintError("Rules error. ")
    }
}

func PrintData(data map[string]interface{}) {
    action, _ := data["action"]
    fmt.Println("\n============", action, "============")
    detail, _ := data["detail"]
    d, _ := json.MarshalIndent(detail, "", "  ")
    fmt.Println(string(d))
}
