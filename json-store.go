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
    NamePrefixIdMap map[string]string
    NameIdMap map[string]string // name -> id
    NodeMap map[string]Node  // id -> node map
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
        if idPrefix, ok := jsonStore.gaiaData.NamePrefixIdMap[parts[0]]; ok {
            id = idPrefix
        } else {
            prefixUsageMap := make(map[string]bool)
            for _, v := range jsonStore.gaiaData.NamePrefixIdMap {
                prefixUsageMap[v] = true
            }

            if len(prefixUsageMap) == 256 {
                return "", errors.New("node name head is greater than 255")
            }

            for i := 0; i < 256; i++ {
                iHex := fmt.Sprintf("%02x", i)
                if !prefixUsageMap[iHex] {
                    jsonStore.gaiaData.NamePrefixIdMap[parts[0]] = iHex
                    id = iHex
                    break
                }
            }
        }

        for i := 0; i < 0xfff; i++ {
            iHex := fmt.Sprintf("%x", i)
            id = id + iHex
            _, exist := jsonStore.gaiaData.NodeMap[id]
            if !exist {
                break
            }
        }

        if _, exist := jsonStore.gaiaData.NodeMap[id]; exist {
            return "", errors.New("No more space for name prefix:" + parts[0])
        }

        return id, nil
    }

    (&node).Normalize()
    if jsonStore.gaiaData.NameIdMap[node.Name] != "" {
        return errors.New("node name exist:" + node.Name)
    }

    id, err := generateId()
    if err != nil{
        return err
    }

    fmt.Println("generate new node id:", id)
    node.Id = id
    node = node.DoTrim()
    jsonStore.gaiaData.NameIdMap[node.Name] = id
    jsonStore.gaiaData.NodeMap[id] = node

    return jsonStore.saveToFile()
}

func (jsonStore *JsonFileStore) AddAlias(from, to string) error {
    err := jsonStore.load()
    if err != nil {
        return err
    }

    from = strings.TrimSpace("from")
    to = stringds.TrimSpace("to")

    jsonStore.gaiaData.AliasMap[from] = to
    return jsonStore.saveToFile()
}

func (jsonStore *JsonFileStore) Update(node Node) error {
    node = node.DoTrim()
    old, exist := jsonStore.gaiaData.NodeMap[node.Id]
    if !exist {
        return errors.New("node with id" + node.Id + " is not exist")
    }

    if old.Name != node.Name {
        return errors.New("node's Name changed!")
    }

    jsonStore.gaiaData.NodeMap[node.Id] = node
    return jsonStore.saveToFile()
}

func (jsonStore *JsonFileStore) Append(id string, extraContent string) error {
    return errors.New("unimplemented")
}

func (jsonStore *JsonFileStore) Search(category string, keywords []string) []Node {
    res := []Node{}

    isCategoryMatch := func(node Node) bool {
        if category == "" {
            return true
        } else {
            return node.Category == category
        }
    }

    isKeywordsMatch := func(node Node) bool {
        arrayContains := func(arr []string, dest string) bool {
            if strings.TrimSpace(dest) == "" {
                return true
            }

            for _, s := range arr {
                if s == dest {
                    return true
                }
            }
            return false
        }

        nameParts := strings.Split(node.Name, "-")
        tags := strings.Split(node.Tags, ",")

        headAllowed := []string{nameParts[0]}
        headAllowed = append(headAllowed, tags...)
        tailAllowed := []string{}
        tailAllowed = append(tailAllowed, nameParts[1:]...)
        tailAllowed = append(tailAllowed, tags...)

        res := arrayContains(headAllowed, nameParts[0])
        for _, k := range nameParts[1:] {
            res = res && arrayContains(tailAllowed, k)
        }
        return res
    }

    for _, node := range jsonStore.gaiaData.NodeMap {
        if isCategoryMatch(node) && isKeywordsMatch(node) {
            res = append(res, node)
        }
    }

    return res
}

func (jsonStore *JsonFileStore) Remove(id string) error {
    fmt.Println("remove node with id:", id)
    node := jsonStore.gaiaData.NodeMap[id]
    delete(jsonStore.gaiaData.NodeMap, id)
    name := node.Name
    delete(jsonStore.gaiaData.NameIdMap, name)

    parts := strings.Split(name, "-")
    namePrefixExist := false
    for n, _ := range jsonStore.gaiaData.NameIdMap {
        if n == parts[0] {
            namePrefixExist = true
            break
        }
    }

    if !namePrefixExist {
        delete(jsonStore.gaiaData.NamePrefixIdMap, parts[0])
    }

    return jsonStore.saveToFile()
}

func (jsonStore *JsonFileStore) GetById(id string) (Node, error) {
    if node, exist := jsonStore.gaiaData.NodeMap[id]; exist {
        return node, nil
    } else {
        return Node{}, errors.New("Node with id " + id + " not found")
    }
}

func (jsonStore *JsonFileStore) GetStats() Stats {
    categories := []string{}
    tags := []string{}
    name0Arr := []string{}

    for _, node := range jsonStore.gaiaData.NodeMap {
        if !existInArray(categories, node.Category) {
            categories = append(categories, node.Category)
        }

        name0 := strings.Split(node.Name, "-")[0]
        if !existInArray(name0Arr, name0) {
            name0Arr = append(name0Arr, name0)
        }

        if node.Tags != "" {
            tagsOfNode := strings.Split(node.Tags, ",")
            for _, t := range tagsOfNode {
                if !existInArray(tags, t) {
                    tags = append(tags, t)
                }
            }
        }
    }

    return Stats{
        CategorySize: len(categories),
        NodeSize: len(jsonStore.gaiaData.NodeMap),
        TagSize: len(tags),
        Name0Size: len(name0Arr),
    }
}

func (jsonStore *JsonFileStore) GetAlias() map[string]string {
    return jsonStore.gaiaData.AliasMap
}

func (jsonStore *JsonFileStore) ListCategories() map[string][]string {
    resultMap := make(map[string][]string)
    for name, id := range jsonStore.gaiaData.NameIdMap {
        node := jsonStore.gaiaData.NodeMap[id]
        nameHead := strings.Split(name, "-")[0]
        headList, exist := resultMap[node.Category]
        if exist {
            if !existInArray(headList, nameHead) {
                headList = append(headList, nameHead)
            }
        } else {
            headList = []string{}
            headList = append(headList, nameHead)
        }
        resultMap[node.Category] = headList
    }

    return resultMap
}

func (jsonStore *JsonFileStore) ListNodes(names []string) []Node {
    resultArray := []Node{}
    namePrefix := strings.Join(names, "-")
    for _, node := range jsonStore.gaiaData.NodeMap {
        if strings.HasPrefix(node.Name, namePrefix) {
            resultArray = append(resultArray, node)
        }
    }
    return resultArray
}

func (jsonStore *JsonFileStore) FormatData() {
    newAliasMap := make(map[string]string)
    for k, v := range jsonStore.gaiaData.AliasMap {
        newAliasMap[strings.TrimSpace(k)] = strings.TrimSpace(v)
    }
    jsonStore.gaiaData.AliasMap = newAliasMap

    newNamePrefixIdMap := make(map[string]string)
    for k, v := range jsonStore.gaiaData.NamePrefixIdMap {
        newNamePrefixIdMap[strings.TrimSpace(k)] = strings.TrimSpace(v)
    }
    jsonStore.gaiaData.NamePrefixIdMap = newNamePrefixIdMap

    newNameIdMap := make(map[string]string)
    for k, v := range jsonStore.gaiaData.NameIdMap {
        newNameIdMap[strings.TrimSpace(k)] = strings.TrimSpace(v)
    }
    jsonStore.gaiaData.NameIdMap = newNameIdMap

    newNodeMap := make(map[string]Node)
    for k, v := range jsonStore.gaiaData.NodeMap {
        newNodeMap[strings.TrimSpace(k)] = v.DoTrim()
    }
    jsonStore.gaiaData.NodeMap = newNodeMap

    jsonStore.saveToFile()
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

    if jsonStore.gaiaData.NamePrefixIdMap == nil {
        jsonStore.gaiaData.NamePrefixIdMap = make(map[string]string)
    }

    if jsonStore.gaiaData.NameIdMap == nil {
        jsonStore.gaiaData.NameIdMap = make(map[string]string)
    }

    if jsonStore.gaiaData.NodeMap == nil {
        jsonStore.gaiaData.NodeMap = make(map[string]Node)
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

func existInArray(arr []string, s string) bool {
    for _, str := range arr {
        if s == str {
            return true
        }
    }
    return false
}
