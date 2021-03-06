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
    CategoryIdMap map[string]string
    BranchIdMap map[string]string
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
    (&node).Normalize(jsonStore.gaiaData.AliasMap)
    if node.Name == "" {
        return errors.New("node name is empty")
    }
    if jsonStore.gaiaData.NameIdMap[node.Name] != "" {
        return errors.New("node name exist:" + node.Name)
    }

    id, err := jsonStore.generateId(node.Name)
    if err != nil{
        return err
    }

    fmt.Println("generate new node id:", id)
    node.Id = id
    node.Category = strings.Split(node.Name, "-")[0]
    jsonStore.gaiaData.NameIdMap[node.Name] = id
    jsonStore.gaiaData.NodeMap[id] = node

    return jsonStore.saveToFile()
}

func (jsonStore *JsonFileStore) AddAlias(from, to string) error {
    from = strings.ToLower(strings.TrimSpace(from))
    to = strings.ToLower(strings.TrimSpace(to))

    jsonStore.gaiaData.AliasMap[from] = to
    return jsonStore.saveToFile()
}

func (jsonStore *JsonFileStore) RemoveAlias(keyword string) error {
    keyword = strings.ToLower(strings.TrimSpace(keyword))
    delete(jsonStore.gaiaData.AliasMap, keyword)
    return jsonStore.saveToFile()
}

func (jsonStore *JsonFileStore) Update(node Node) error {
    (&node).Normalize(jsonStore.gaiaData.AliasMap)
    old, exist := jsonStore.gaiaData.NodeMap[node.Id]
    if !exist {
        return errors.New("node with id" + node.Id + " is not exist")
    }

    oldBranch := old.GetBranch()
    newBranch := node.GetBranch()

    if oldBranch != newBranch {
        return errors.New("can not do update, node's branch changed!")
    }

    jsonStore.gaiaData.NodeMap[node.Id] = node
    return jsonStore.saveToFile()
}

func (jsonStore *JsonFileStore) Append(id string, extraContent string) error {
    node, exist := jsonStore.gaiaData.NodeMap[id]
    if !exist {
        return errors.New("node with id " + id + " not exists")
    }

    oldContent := strings.TrimSpace(node.Content)
    node.Content = oldContent + "\n\n" + strings.TrimSpace(extraContent)
    jsonStore.gaiaData.NodeMap[id] = node

    return jsonStore.saveToFile()
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

    // TODO: in this function i implemented search AND,
    // maybe a search OR function should be implemented.
    isKeywordsMatch := func(node Node) bool {
        arrayContains := func(arr []string, dest string) bool {
            if strings.TrimSpace(dest) == "" {
                return true
            }

            for _, s := range arr {
                if strings.HasPrefix(s, dest) {
                    return true
                }
            }
            return false
        }

        nameParts := strings.Split(node.Name, "-")
        tags := strings.Split(node.Tags, ",")

        nodeKeywords := nameParts
        nodeKeywords = append(nodeKeywords, tags...)

        isMatch := true
        for _, k := range keywords {
            isMatch = isMatch && arrayContains(nodeKeywords, k)
        }
        return isMatch
    }

    for _, node := range jsonStore.gaiaData.NodeMap {
        if isCategoryMatch(node) && isKeywordsMatch(node) {
            res = append(res, node)
        }
    }

    return res
}

func (jsonStore *JsonFileStore) Remove(id string) error {
    node := jsonStore.gaiaData.NodeMap[id]
    delete(jsonStore.gaiaData.NodeMap, id)
    name := node.Name
    delete(jsonStore.gaiaData.NameIdMap, name)

    parts := strings.Split(name, "-")
    namePrefixExist := false
    for n, _ := range jsonStore.gaiaData.NameIdMap {
        if strings.HasPrefix(n, parts[0]) {
            namePrefixExist = true
            break
        }
    }

    if !namePrefixExist {
        delete(jsonStore.gaiaData.CategoryIdMap, parts[0])
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
    for _, node := range jsonStore.gaiaData.NodeMap {
        if !existInArray(categories, node.Category) {
            categories = append(categories, node.Category)
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
    }
}

func (jsonStore *JsonFileStore) GetAlias() map[string]string {
    return jsonStore.gaiaData.AliasMap
}

func (jsonStore *JsonFileStore) ListCategories() map[string][]string {
    resultMap := make(map[string][]string)
    for name, id := range jsonStore.gaiaData.NameIdMap {
        node := jsonStore.gaiaData.NodeMap[id]
        parts := strings.Split(name, "-")
        middlePart := ""
        if len(parts) > 1 {
            middlePart = parts[1]
        } else {
            continue
        }
        middlePartList, exist := resultMap[node.Category]
        if exist {
            if !existInArray(middlePartList, middlePart) {
                middlePartList = append(middlePartList, middlePart)
            }
        } else {
            middlePartList = []string{ middlePart }
        }
        resultMap[node.Category] = middlePartList
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

func (jsonStore *JsonFileStore) ReplaceAlias(strArr []string) []string {
    replacedArr := []string{}
    for _, s := range strArr {
        s = strings.ToLower(strings.TrimSpace(s))
        if jsonStore.gaiaData.AliasMap[s] != "" {
            replacedArr = append(replacedArr, jsonStore.gaiaData.AliasMap[s])
        } else {
            replacedArr = append(replacedArr, s)
        }
    }
    return replacedArr
}

func (jsonStore *JsonFileStore) FormatData() error {
    newAliasMap := make(map[string]string)
    for k, v := range jsonStore.gaiaData.AliasMap {
        newAliasMap[strings.TrimSpace(k)] = strings.TrimSpace(v)
    }
    jsonStore.gaiaData.AliasMap = newAliasMap

    newCategoryIdMap := make(map[string]string)
    for k, v := range jsonStore.gaiaData.CategoryIdMap {
        newCategoryIdMap[strings.TrimSpace(k)] = strings.TrimSpace(v)
    }
    jsonStore.gaiaData.CategoryIdMap = newCategoryIdMap

    newNameIdMap := make(map[string]string)
    for k, v := range jsonStore.gaiaData.NameIdMap {
        newNameIdMap[strings.TrimSpace(k)] = strings.TrimSpace(v)
    }
    jsonStore.gaiaData.NameIdMap = newNameIdMap

    newNodeMap := make(map[string]Node)
    for k, v := range jsonStore.gaiaData.NodeMap {
        (&v).Normalize(jsonStore.gaiaData.AliasMap)
        newNodeMap[strings.TrimSpace(k)] = v
    }
    jsonStore.gaiaData.NodeMap = newNodeMap

    return jsonStore.saveToFile()
}

func (jsonStore *JsonFileStore) ReorgAllData() error {
    jsonStore.gaiaData.CategoryIdMap = map[string]string{
        "self": "0",
        "math": "1",
        "algorithm": "2",
        "utils": "3",
        "os": "4",
        "tools": "5",
        "lang": "6",
        "others": "7",
        "archive": "a",
        "cache": "c",
        "db": "d",
        "eth": "e",
        "frontend": "f",
    }
    jsonStore.gaiaData.BranchIdMap = map[string]string{}
    jsonStore.gaiaData.NameIdMap = map[string]string{}

    oldNodeMap := jsonStore.gaiaData.NodeMap
    jsonStore.gaiaData.NodeMap = map[string]Node{}

    for _, node := range oldNodeMap {
        node.Id = ""
        err := jsonStore.Add(node)
        if err != nil {
            fmt.Println("err:", err)
            // return err
            continue
        }
    }

    return jsonStore.saveToFile()
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

    if jsonStore.gaiaData.CategoryIdMap == nil {
        jsonStore.gaiaData.CategoryIdMap = make(map[string]string)
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

func (jsonStore *JsonFileStore) generateId(nodeName string) (string, error) {
    parts := strings.Split(nodeName, "-")
    idPrefix, ok := jsonStore.gaiaData.CategoryIdMap[parts[0]]
    if !ok {
        prefixUsageMap := make(map[string]bool)
        for _, v := range jsonStore.gaiaData.CategoryIdMap {
            prefixUsageMap[v] = true
        }

        if len(prefixUsageMap) >= 16 {
            return "", errors.New("out of category space, max allowed: 16")
        }

        for i := 0; i < 16; i++ {
            idPrefix = fmt.Sprintf("%x", i)
            if !prefixUsageMap[idPrefix] {
                jsonStore.gaiaData.CategoryIdMap[parts[0]] = idPrefix
                break
            }
        }
    }

    if len(parts) == 1 {
        return idPrefix, nil
    }

    branch := parts[0] + "-" + parts[1]
    branchIdPrefix, exist := jsonStore.gaiaData.BranchIdMap[branch]
    if !exist {
        branchPrefixUsageMap := make(map[string]bool)
        for _, v := range jsonStore.gaiaData.BranchIdMap {
            branchPrefixUsageMap[v] = true
        }

        for i := 0; i < 16; i++ {
            branchIdPrefix = idPrefix + fmt.Sprintf("%x", i)
            if !branchPrefixUsageMap[branchIdPrefix] {
                jsonStore.gaiaData.BranchIdMap[branch] = branchIdPrefix
                break
            }
            branchIdPrefix = ""
        }

        if branchIdPrefix == "" {
            return "", errors.New("out of branch space in category, max allowed: 16")
        }
    }

    if len(parts) == 2 {
        return branchIdPrefix, nil
    }

    idTail := ""
    for j := 0; j < 256; j++ {
        idTail = fmt.Sprintf("%02x", j)
        tmpId := branchIdPrefix + idTail
        _, exist := jsonStore.gaiaData.NodeMap[tmpId]
        if !exist {
            break
        }
        idTail = ""
    }
    if idTail == "" {
        return "", errors.New("No more space in branch: " + branch + "; max allowed: 256")
    }

    id := branchIdPrefix + idTail
    return id, nil
}

func existInArray(arr []string, s string) bool {
    for _, str := range arr {
        if s == str {
            return true
        }
    }
    return false
}
