package main

// import (
//     "bufio"
//     "fmt"
//     "os"
//     "strings"
// )

// const IdLen = 27

// type CodeSegment struct {
//     Id, Category, Tags, Desc, Code string
// }

// func (cs CodeSegment) PrintToScreen() {
//     fmt.Printf("  ID: %s\nCATE: %s\nTAGS: %s\n", cs.Id, cs.Category, cs.Tags)
//     fmt.Printf("DESC: %s\n", cs.Desc)
//     codeLines := strings.Split(cs.Code, "\n")
//     for i, line := range codeLines {
//         if i == 0 {
//             fmt.Println("CONTENT:")
//             fmt.Println("      " + line)
//         } else {
//             fmt.Println("      " + line)
//         }
//     }
// }

// var CodePrefixSpace string = "          " // len:10
// func (cs CodeSegment) PrintToFile(fpath string) error {
//     f, err := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
//     if err != nil {
//         return err
//     }
//     defer f.Close()

//     f.WriteString("Id:       " + cs.Id + "\n")
//     f.WriteString("Category: " + cs.Category + "\n")
//     f.WriteString("Tags:     " + cs.Tags + "\n")
//     f.WriteString("Desc:     " + cs.Desc + "\n")

//     codeLines := strings.Split(cs.Code, "\n")
//     for i, line := range codeLines {
//         if i == 0 {
//             f.WriteString("Content:  " + line + "\n")
//         } else {
//             f.WriteString(CodePrefixSpace + line + "\n")
//         }
//     }
//     return nil
// }

// func (cs *CodeSegment) ReadFromFile(fpath string) error {
//     f, err := os.Open(fpath)
//     if err != nil {
//         return err
//     }
//     defer f.Close()

//     isCodeLine := false
//     isDescLine := false
//     scanner := bufio.NewScanner(f)
//     //var codePrefixSpace string
//     for scanner.Scan() {
//         line := scanner.Text()
//         if strings.HasPrefix(line, "Id:") {
//             cs.Id = strings.TrimSpace(line[len("Id:"):])
//         } else if strings.HasPrefix(line, "Category:") {
//             cs.Category = strings.TrimSpace(line[len("Category:"):])
//         } else if strings.HasPrefix(line, "Tags:") {
//             cs.Tags = strings.TrimSpace(line[len("Tags:"):])
//         } else if strings.HasPrefix(line, "Desc:") {
//             cs.Desc = strings.TrimSpace(line[len("Desc:"):])
//             isDescLine = true
//             isCodeLine = false
//         } else if strings.HasPrefix(line, "Content:") {
//             cs.Code = strings.TrimSpace(line[len("Content:"):])
//             isCodeLine = true
//             isDescLine = false
//         } else {
//             if isDescLine {
//                 var descLine string
//                 if strings.HasPrefix(line, CodePrefixSpace) {
//                     descLine = line[len(CodePrefixSpace):]
//                 } else {
//                     descLine = strings.TrimSpace(line)
//                 }
//                 cs.Desc = cs.Desc + "\n" + descLine
//                 //fmt.Println("cs.Desc:", cs.Desc)
//             }

//             if isCodeLine {
//                 var codeLine string
//                 if strings.HasPrefix(line, CodePrefixSpace) {
//                     codeLine = line[len(CodePrefixSpace):]
//                 } else {
//                     codeLine = strings.TrimSpace(line)
//                 }
//                 cs.Code = cs.Code + "\n" + codeLine
//                 //fmt.Println("cs.Code:", cs.Code)
//             }
//         }
//     }
//     return nil
// }
