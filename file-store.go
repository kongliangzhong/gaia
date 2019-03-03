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

func (fs *FileStore) leafToStr(lf Leaf) string {
    descB64 := base64.StdEncoding.EncodeToString([]byte(lf.Desc))
    contentB64 := base64.StdEncoding.EncodeToString([]byte(lf.Content))
    return lf.Id + "|" + lf.Category + "|" + lf.Tags + "|" + descB64 + "|" + contentB64
}

func (fs *FileStore) strToLeaf(str string) (lf Leaf, err error) {
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

    lf = Leaf{Id: id, Category: cate, Tags: tags, Desc: desc, Content: code}
    return
}

func (fs *FileStore) genId(lf Leaf) (id string, err error) {
    if lf.Id != "" {
        err = errors.New("id already exists")
        return
    }

    if lf.Category == "" && lf.Tags == "" {
        err = errors.New("generate id failed: category and tags both empty")
    }

    idBytes := sha1.Sum([]byte(lf.Category + lf.Tags))
    id = base64.StdEncoding.EncodeToString(idBytes[:])
    id = id[:len(id)-1] //
    return
}

func (fs *FileStore) Add(lf Leaf) error {
    if lf.Id == "" {
        id, err := fs.genId(lf)
        if err != nil {
            return err
        }
        //fmt.Println("id: ", id, "id len:", len(id))
        lf.Id = id
    }

    if err := fs.isDuplicate(lf); err != nil {
        return err
    }

    f, err := os.OpenFile(fs.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
    if err != nil {
        return err
    }
    defer f.Close()

    line := fs.leafToStr(lf)
    _, err = f.WriteString(line + "\n")
    return err
}

func (fs *FileStore) GetById(id string) (lf Leaf, err error) {
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
            return fs.strToLeaf(line)
        }
    }

    err = errors.New("can not find code-segment by id:" + id)
    return
}

func (fs *FileStore) Update(lf Leaf) error {
    newLf, err := fs.GetById(lf.Id)
    if err != nil {
        return err
    }

    if lf.Category != "" {
        newLf.Category = lf.Category
    }

    if lf.Tags != "" {
        newLf.Tags = lf.Tags
    }

    if lf.Desc != "" {
        newLf.Desc = lf.Desc
    }

    if lf.Content != "" {
        newLf.Content = lf.Content
    }

    fs.Remove(lf.Id)
    return fs.Add(newLf)
}

func (fs *FileStore) Append(id string, extraContent string) error {
    newLf, err := fs.GetById(id)
    if err != nil {
        return err
    }

    newLf.Content = strings.Trim(newLf.Content, "\n") + "\n" + strings.Trim(extraContent, "\n")
    fs.Remove(id)
    return fs.Add(newLf)
}

func (fs *FileStore) Search(category string, tagStr string) []Leaf {
    //tags := strings.Split(tagStr, ",")
    matchedLines := grepFile(fs.FilePath, category, tagStr)
    matchedLf := []Leaf{}
    for _, line := range matchedLines {
        lf, err := fs.strToLeaf(line)
        if err != nil {
            fmt.Println(err.Error())
        }
        matchedLf = append(matchedLf, lf)
    }
    return matchedLf
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

        allTagsOfLf := strings.Split(tagsInStore, ",")
        for _, t := range allTagsOfLf {
            subTs := strings.Split(t, "-")
            if len(subTs) > 1 {
                for _, subTag := range subTs {
                    allTagsOfLf = append(allTagsOfLf, subTag)
                }
            }
        }

        reqTags := strings.Split(reqTagStr, ",")

        for _, reqTag := range reqTags {
            isContains := false
            for _, tagOfLf := range allTagsOfLf {
                if tagOfLf == reqTag {
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

func (fs *FileStore) isDuplicate(lf Leaf) error {
    f, err := os.Open(fs.FilePath)
    if err != nil {
        return nil
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        lfInFile, _ := fs.strToLeaf(line)
        if lfInFile.Content == lf.Content {
            return errors.New("duplicated code content with id:" + lfInFile.Id)
        }
        if lfInFile.Id == lf.Id {
            return errors.New("duplicated id generated. category and tags is the same with code segment " + lfInFile.Id)
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
        rlf, err := fs.strToLeaf(line)
        if err != nil {
            fmt.Println(err)
            continue
        }
        stats.TotalSize ++
        cate := rlf.Category
        tagStr := rlf.Tags
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
