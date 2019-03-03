package main

import (
    "io/ioutil"
    "encoding/json"
    "errors"
)

type GaiaData struct {
    AliasMap map[string]string
    IdPrefixMap map[string]string
    IdCursorMap map[string]int
    LeavesMap map[string]Leaf
}

type JsonFileStore struct {
    FilePath string
    JsonData *GaiaData
}

func (jsonStore *JsonFileStore) Add(lf Leaf) error {
    return errors.New("unimplemented")
}

func (jsonStore *JsonFileStore) AddAlias(from, to string) error {
    return errors.New("unimplemented")
}

func (jsonStore *JsonFileStore) Update(lf Leaf) error {
    return errors.New("unimplemented")
}

func (jsonStore *JsonFileStore) Append(id string, extraContent string) error {
    return errors.New("unimplemented")
}

func (jsonStore *JsonFileStore) Search(category string, tagStr string) []Leaf {
    return []Leaf{}
}

func (jsonStore *JsonFileStore) Remove(id string) error {
    return errors.New("unimplemented")
}

func (jsonStore *JsonFileStore) GetById(id string) (Leaf, error) {
    return Leaf{}, errors.New("unimplemented")
}

func (jsonStore *JsonFileStore) GetStats() Stats {
    return Stats{}
}

func (jsonStore *JsonFileStore) load() error {
    jsonStr, err := ioutil.ReadFile(jsonStore.FilePath)
    if err != nil {
        return err
    }

    err = json.Unmarshal([]byte(jsonStr), jsonStore.JsonData)
    return err
}

func (jsonStore *JsonFileStore) saveToJsonFile() error {
    bs, err := json.Marshal(jsonStore.JsonData)
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(jsonStore.FilePath, bs, 0660)
    return nil
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
