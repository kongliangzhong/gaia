package main

import (
    "bufio"
    "crypto/sha1"
    "encoding/base64"
    "errors"
    "fmt"
    "os"
    "strings"
)

type FileStore struct {
    FilePath string
}

func (fs *FileStore) nodeToStr(node Node) string {
    descB64 := base64.StdEncoding.EncodeToString([]byte(node.Desc))
    contentB64 := base64.StdEncoding.EncodeToString([]byte(node.Content))
    return node.Id + "|" + node.Category + "|" + node.Tags + "|" + descB64 + "|" + contentB64
}

func (fs *FileStore) strToNode(str string) (node Node, err error) {
    flds := strings.Split(str, "|")
    if len(flds) != 5 {
        err = errors.New("parse segemnt str failed: " + str)
        return
    }
    id := flds[0]
    cate := flds[1]
    tags := flds[2]
    desc := flds[3]
    code := flds[4]

    descBs, err := base64.StdEncoding.DecodeString(desc)
    if err != nil {
        return
    }

    codeBs, err := base64.StdEncoding.DecodeString(code)
    if err != nil {
        return
    }

    desc = string(descBs)
    code = string(codeBs)

    node = Node{Id: id, Category: cate, Tags: tags, Desc: desc, Content: code}
    return
}

func (fs *FileStore) genId(node Node) (id string, err error) {
    if node.Id != "" {
        err = errors.New("id already exists")
        return
    }

    if node.Category == "" && node.Tags == "" {
        err = errors.New("generate id failed: category and tags both empty")
    }

    idBytes := sha1.Sum([]byte(node.Category + node.Tags))
    id = base64.StdEncoding.EncodeToString(idBytes[:])
    id = id[:len(id)-1] //
    return
}

func (fs *FileStore) Add(node Node) error {
    if node.Id == "" {
        id, err := fs.genId(node)
        if err != nil {
            return err
        }
        //fmt.Println("id: ", id, "id len:", len(id))
        node.Id = id
    }

    if err := fs.isDuplicate(node); err != nil {
        return err
    }

    f, err := os.OpenFile(fs.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
    if err != nil {
        return err
    }
    defer f.Close()

    line := fs.nodeToStr(node)
    _, err = f.WriteString(line + "\n")
    return err
}

func (fs *FileStore) GetById(id string) (node Node, err error) {
    if len(id) < IdLen {
        err = errors.New("invalid id:" + id)
        return
    }

    f, err := os.Open(fs.FilePath)
    if err != nil {
        return
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, id) {
            return fs.strToNode(line)
        }
    }

    err = errors.New("can not find code-segment by id:" + id)
    return
}

func (fs *FileStore) Update(node Node) error {
    newNode, err := fs.GetById(node.Id)
    if err != nil {
        return err
    }

    if node.Category != "" {
        newNode.Category = node.Category
    }

    if node.Tags != "" {
        newNode.Tags = node.Tags
    }

    if node.Desc != "" {
        newNode.Desc = node.Desc
    }

    if node.Content != "" {
        newNode.Content = node.Content
    }

    fs.Remove(node.Id)
    return fs.Add(newNode)
}

func (fs *FileStore) Append(id string, extraContent string) error {
    newNode, err := fs.GetById(id)
    if err != nil {
        return err
    }

    newNode.Content = strings.Trim(newNode.Content, "\n") + "\n" + strings.Trim(extraContent, "\n")
    fs.Remove(id)
    return fs.Add(newNode)
}

func (fs *FileStore) Search(category string, tagStr string) []Node {
    //tags := strings.Split(tagStr, ",")
    matchedLines := grepFile(fs.FilePath, category, tagStr)
    matchedNode := []Node{}
    for _, line := range matchedLines {
        node, err := fs.strToNode(line)
        if err != nil {
            fmt.Println(err.Error())
        }
        matchedNode = append(matchedNode, node)
    }
    return matchedNode
}

func grepFile(file string, reqCate string, reqTagStr string) []string {
    var categoryMatch = func(cateInStore string, cateReq string) bool {
        if cateReq == "" {
            return true
        }

        cateInStore = strings.ToUpper(cateInStore)
        cateReq = strings.ToUpper(cateReq)

        cates := strings.Split(cateInStore, "-")
        for _, cate := range cates {
            if cate == cateReq {
                return true
            }
        }
        return false
    }

    var tagsMatch = func(tagsInStore string, reqTagStr string) bool {
        if reqTagStr == "" {
            return true
        }

        tagsInStore = strings.ToUpper(tagsInStore)
        reqTagStr = strings.ToUpper(reqTagStr)

        allTagsOfNode := strings.Split(tagsInStore, ",")
        for _, t := range allTagsOfNode {
            subTs := strings.Split(t, "-")
            if len(subTs) > 1 {
                for _, subTag := range subTs {
                    allTagsOfNode = append(allTagsOfNode, subTag)
                }
            }
        }

        reqTags := strings.Split(reqTagStr, ",")

        for _, reqTag := range reqTags {
            isContains := false
            for _, tagOfNode := range allTagsOfNode {
                if tagOfNode == reqTag {
                    isContains = true
                }
            }

            if !isContains {
                return false
            }
        }
        return true
    }

    //cate = strings.ToUpper(cate)
    res := []string{}
    f, err := os.Open(file)
    if err != nil {
        fmt.Println(err)
        return res
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        flds := strings.Split(line, "|")
        cateInStore := flds[1]
        tagStr := flds[2]
        tagStr = cateInStore + "," + tagStr
        if categoryMatch(cateInStore, reqCate) && tagsMatch(tagStr, reqTagStr) {
            res = append(res, line)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
    return res
}

func (fs *FileStore) Remove(id string) error {
    if len(id) < IdLen {
        return errors.New("Invalid id, id is too short")
    }

    f, err := os.Open(fs.FilePath)
    if err != nil {
        return err
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    fLines := []string{}
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, id) {
            continue
        }
        fLines = append(fLines, line)
    }

    oldFilePath := fs.FilePath + ".old"
    os.Remove(oldFilePath)
    err = os.Rename(fs.FilePath, oldFilePath)
    if err != nil {
        fmt.Println(err.Error())
    }

    newFile, err := os.OpenFile(fs.FilePath, os.O_CREATE|os.O_WRONLY, 0660)
    if err != nil {
        return err
    }
    defer newFile.Close()
    for _, line := range fLines {
        newFile.WriteString(line + "\n")
    }

    return nil
}

func (fs *FileStore) isDuplicate(node Node) error {
    f, err := os.Open(fs.FilePath)
    if err != nil {
        return nil
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        nodeInFile, _ := fs.strToNode(line)
        if nodeInFile.Content == node.Content {
            return errors.New("duplicated code content with id:" + nodeInFile.Id)
        }
        if nodeInFile.Id == node.Id {
            return errors.New("duplicated id generated. category and tags is the same with code segment " + nodeInFile.Id)
        }
    }
    return nil
}

func (fs *FileStore) GetStats() Stats {
    stats := Stats{
        AllCates:    []string{},
        AllTags:     []string{},
        CateTagsMap: map[string][]string{},
        CateNumMap:  map[string]int{},
        TagCatesMap: map[string][]string{},
        TagNumMap:   map[string]int{},
    }

    f, err := os.Open(fs.FilePath)
    if err != nil {
        fmt.Println(err)
        return stats
    }

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        rnode, err := fs.strToNode(line)
        if err != nil {
            fmt.Println(err)
            continue
        }
        stats.TotalSize ++
        cate := rnode.Category
        tagStr := rnode.Tags
        tagsArr := strings.Split(tagStr, ",")

        cateSize := stats.CateNumMap[cate]
        stats.CateNumMap[cate] = cateSize + 1

        if !ArrContains(stats.AllCates, cate) {
            stats.AllCates = append(stats.AllCates, cate)
        }

        for _, t := range tagsArr {
            if !ArrContains(stats.AllTags, t) {
                stats.AllTags = append(stats.AllTags, t)
            }

            tagsOfCate := stats.CateTagsMap[cate]
            if !ArrContains(tagsOfCate, t) {
                tagsOfCate = append(tagsOfCate, t)
                stats.CateTagsMap[cate] = tagsOfCate
            }

            catesOfTag := stats.TagCatesMap[t]
            if !ArrContains(catesOfTag, cate) {
                catesOfTag = append(catesOfTag, cate)
                stats.TagCatesMap[t] = catesOfTag
            }

            tagSize := stats.TagNumMap[t]
            stats.TagNumMap[t] = tagSize + 1
        }
    }

    return stats
}
