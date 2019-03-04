package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "encoding/json"
    "errors"
    "strings"
)

type GaiaData struct {
    AliasMap map[string]string
    IdPrefixMap map[string]string
    NameIdMap map[string]string
    LeavesMap map[string]Node
}

type JsonFileStore struct {
    FilePath string
    gaiaData *GaiaData
}

func newJsonFileStore(dataFilePath string) *JsonFileStore {
    jsonStore := &JsonFileStore{dataFilePath, &GaiaData{}}
    jsonStore.load()

    return jsonStore
}

func (jsonStore *JsonFileStore) Add(node Node) error {
    generateId := func() (string, error) {
        parts := strings.Split(node.Name, "-")
        id := ""
        if idPrefix, ok := jsonStore.gaiaData.IdPrefixMap[parts[0]]; ok {
            id = idPrefix
        } else {
            prefixUsageMap := make(map[string]bool)
            for _, v := range jsonStore.gaiaData.IdPrefixMap {
                prefixUsageMap[v] = true
            }

            if len(prefixUsageMap) == 256 {
                return "", errors.New("node name head is greater than 255")
            }

            for i := 0; i < 256; i++ {
                iHex := fmt.Sprintf("%02x", i)
                if !prefixUsageMap[iHex] {
                    jsonStore.gaiaData.IdPrefixMap[parts[0]] = iHex
                    id = iHex
                    break
                }
            }
        }

        for i := 0; i < 0xfff; i++ {
            iHex := fmt.Sprintf("%x", i)
            id = id + iHex
            _, exist := jsonStore.gaiaData.LeavesMap[id]
            if !exist {
                break
            }
        }

        if _, exist := jsonStore.gaiaData.LeavesMap[id]; exist {
            return "", errors.New("No more space for name prefix:" + parts[0])
        }

        return id, nil
    }

    if jsonStore.gaiaData.NameIdMap[node.Name] != "" {
        return errors.New("node name exist:" + node.Name)
    }

    id, err := generateId()
    if err != nil{
        return err
    }

    node.Id = id
    jsonStore.gaiaData.NameIdMap[node.Name] = id
    jsonStore.gaiaData.LeavesMap[id] = node

    return jsonStore.saveToFile()
}

func (jsonStore *JsonFileStore) AddAlias(from, to string) error {
    err := jsonStore.load()
    if err != nil {
        return err
    }

    jsonStore.gaiaData.AliasMap[from] = to
    return jsonStore.saveToFile()
}

func (jsonStore *JsonFileStore) Update(node Node) error {
    return errors.New("unimplemented")
}

func (jsonStore *JsonFileStore) Append(id string, extraContent string) error {
    return errors.New("unimplemented")
}

func (jsonStore *JsonFileStore) Search(category string, tagStr string) []Node {
    return []Node{}
}

func (jsonStore *JsonFileStore) Remove(id string) error {
    return errors.New("unimplemented")
}

func (jsonStore *JsonFileStore) GetById(id string) (Node, error) {
    return Node{}, errors.New("unimplemented")
}

func (jsonStore *JsonFileStore) GetStats() Stats {
    return Stats{}
}

func (jsonStore *JsonFileStore) GetAlias() map[string]string {
    return jsonStore.gaiaData.AliasMap
}

func (jsonStore *JsonFileStore) load() error {
    jsonStrBytes, err := ioutil.ReadFile(jsonStore.FilePath)
    if err != nil {
        if os.IsNotExist(err) {
            ioutil.WriteFile(jsonStore.FilePath, []byte(""), 0660)
            err = nil
        } else {
            return err
        }
    }

    if len(jsonStrBytes) > 0 {
        err = json.Unmarshal(jsonStrBytes, jsonStore.gaiaData)
    }

    // write data backup file:
    ioutil.WriteFile(jsonStore.FilePath + ".bk", jsonStrBytes, 0660)

    // init map if empty in gaiaData
    if jsonStore.gaiaData.AliasMap == nil {
        jsonStore.gaiaData.AliasMap = make(map[string]string)
    }

    if jsonStore.gaiaData.IdPrefixMap == nil {
        jsonStore.gaiaData.IdPrefixMap = make(map[string]string)
    }

    if jsonStore.gaiaData.NameIdMap == nil {
        jsonStore.gaiaData.NameIdMap = make(map[string]string)
    }

    if jsonStore.gaiaData.LeavesMap == nil {
        jsonStore.gaiaData.LeavesMap = make(map[string]Node)
    }

    return err
}

func (jsonStore *JsonFileStore) saveToFile() error {
    bs, err := json.MarshalIndent(jsonStore.gaiaData, "", "  ")
    if err != nil {
        return err
    }

    return ioutil.WriteFile(jsonStore.FilePath, bs, 0660)
}

// func main() {
//     aliasMap := make(map[string]string)
//     aliasMap["javascript"] = "js"
//     aliasMap["typescript"] = "ts"
//     aliasMap["golang"] = "go"

//     gaiaData := &GaiaData{AliasMap: aliasMap}
//     err := gaiaData.SaveToJsonFile("/tmp/1234.json")
//     if err != nil {
//         fmt.Println("error:", err)
//     }

//     gaiaData2 := Load("/tmp/1234.json")
//     fmt.Println("unmarshalData:", gaiaData2)
// }
