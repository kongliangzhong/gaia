package main

import (
    "fmt"
    "flag"
    "os"
    "os/user"
    "strings"
)

// TODO: read and write global vars from/into store
// TODO: add alias globally for search purpose.
// TODO: generate tree by name
// TODO: search: in tree
// TODO: execute code from system clipboard
// TODO: copy code to clipboard

// TODO: read all items by skip:count

var gaiaDir = ".gaia/data/"
const dataFileName = "data.json"
var dataFilePath = ""

var subCommands = []string{
    "add",
    "alias",
    "update",
    "append",
    "merge",
    "list",
    "search",
    "remove",
    "edit",
    "exec",
}

var subCommandMap = map[string]string{
    "add": "add item",
    "alias": "add keyword alias",
    "update": "update item",
    "append": "append text to item content",
    "merge": "merge two item into one",
    "list": "list items",
    "search": "search items",
    "remove": "remove item by id",
    "edit": "edit item in vi",
    "exec": "execute item",
}

var (
    subFlag *flag.FlagSet
    isHelp bool
    id string
    oid string
    name string
    category string
    tags string
    desc string
    content string
    isExecutable bool
    mainFile string

    listCategories bool
    listTags bool
    listAlias bool
    listNames bool
)

func init() {
    usr, err := user.Current()
    if err != nil {
        panic(err)
    }

    gaiaDir = usr.HomeDir + "/" + gaiaDir
    dataFilePath = gaiaDir + dataFileName
    _, err = os.Stat(gaiaDir)
    if err != nil && os.IsNotExist(err) {
        err = os.MkdirAll(gaiaDir, 0770)
        if err != nil {
            panic(err)
        }
    }
}

// func setupFlags() {
//     for _, subCommand := range subCommands {
//         subFlag := flag.NewFlagSet(subCommand, flag.ExitOnError)
//         subFlag.String("h", "", "show help message")

//         switch (subCommand) {
//         case "add":
//             subFlag.String("i", "", "leaf id")
//             subFlag.String("n", "", "leaf name")
//             subFlag.String("c", "", "leaf category")
//             subFlag.String("t", "", "leaf tags, tag seprated by comma")
//             subFlag.String("m", "", "leaf description")
//         case "update":
//         case "append":
//         case "merge":
//         case "list":
//         case "search":
//         case "remove":
//         case "edit":
//         case "exec":

//         default:
//             fmt.Println("Unrecgonized command:", subCommand)
//         }

//     }

//     flag.String("h", "", "show usage")
//     flag.Usage = printUsage
// }

func main() {
    flag.BoolVar(&isHelp, "h", false, "show help message")
    if len(os.Args) == 1 {
        printUsage()
        os.Exit(-1)
    }

    // fmt.Println("os.Args[1]:", os.Args[1])
    // flag.Parse()
    flag.Usage = printUsage

    if os.Args[1] == "-h" || os.Args[1] == "--help" {
        printUsage()
        os.Exit(2)
    }

    subFlag = flag.NewFlagSet(os.Args[1], flag.ExitOnError)
    subFlag.Usage = func() {
        fmt.Printf("Usage: %s %s <args> \n", os.Args[0], os.Args[1])
        subFlag.PrintDefaults()
    }
    switch os.Args[1] {
    case "add":
        subFlag.StringVar(&id, "i", "", "leaf id")
        subFlag.StringVar(&name, "n", "", "leaf name")
        subFlag.StringVar(&category, "c", "", "leaf category")
        subFlag.StringVar(&tags, "t", "", "leaf tags, tag seprated by comma")
        subFlag.StringVar(&desc, "m", "", "leaf description")
        subFlag.StringVar(&content, "b", "", "leaf body content")
        subFlag.BoolVar(&isExecutable, "e", false, "is leaf executable")
        subFlag.StringVar(&mainFile, "f", "", "executable main file name")

        subFlag.Usage = func() {
            fmt.Printf("Usage: %s %s -n name -c category -b body [<other args>] \n", os.Args[0], os.Args[1])
            subFlag.PrintDefaults()
        }
    case "alias":
        subFlag.Usage = func() {
            fmt.Printf("Usage: %s %s <keyword> <target-keyword> \n", os.Args[0], os.Args[1])
            subFlag.PrintDefaults()
        }
    case "update":
        subFlag.StringVar(&id, "i", "", "leaf id")
        subFlag.StringVar(&content, "b", "", "leaf body content")
    case "append":
        subFlag.StringVar(&id, "i", "", "leaf id")
        subFlag.StringVar(&content, "b", "", "leaf append content")
    case "merge":
        subFlag.StringVar(&id, "i", "", "leaf id")
        subFlag.StringVar(&oid, "d", "", "dest leaf id")
    case "list":
        subFlag.BoolVar(&listCategories, "c", false, "list leaf categories")
        subFlag.BoolVar(&listTags, "t", false, "list leaf tags")
        subFlag.BoolVar(&listNames, "n", false, "list by name parts")
        subFlag.BoolVar(&listAlias, "a", false, "list global keyword alias")
    case "search":
        subFlag.StringVar(&category, "c", "", "list leaf categories")
    case "remove":
        subFlag.StringVar(&id, "i", "", "leaf id")
    case "edit":
        subFlag.StringVar(&id, "i", "", "leaf id")
    case "exec":
        subFlag.StringVar(&id, "i", "", "leaf id")
    default:
        fmt.Println("Unrecogniz command:", os.Args[1])
        printUsage()
        os.Exit(2)
    }

    subFlag.BoolVar(&isHelp, "h", false, "show help message")

    if len(os.Args) == 2 ||
        os.Args[2] == "-h" ||
        os.Args[2] == "--help" {
        subFlag.Usage()
        os.Exit(2)
    }

    subFlag.Parse(os.Args[2:])
    processSubCommand(os.Args[1])
}

func processSubCommand(command string) {
    op := newOperator(newJsonFileStore(dataFilePath))

    switch command {
    case "add":
        checkRequiredArg("-n", name)
        checkRequiredArg("-c", category)
        checkRequiredArg("-b", content)
        leaf := Leaf{
            Name: name,
            Category: category,
            Tags: tags,
            Desc: desc,
            Content: content,
            IsExecutable: isExecutable,
            MainFileName: mainFile,
        }
        op.Add(leaf)
    case "alias":
        aliasArgs := subFlag.Args()
        if len(aliasArgs) != 2 {
            subFlag.Usage()
            os.Exit(2)
        }
        op.AddAlias(aliasArgs[0], aliasArgs[1])

    case "update":
        // leaf := parseArgs(os.Args)
        // op.Update(leaf)
    case "append":
        // leaf := parseArgs(os.Args)
        // op.Append(leaf.Id, leaf.Content)
    case "merge":
        // ids := os.Args[2:]
        // op.Merge(ids...)
    case "list":
        if listAlias {
            op.ListAlias()
        } else if listCategories {

        } else if listTags {

        } else if listNames {

        } else {

        }
        //op.ListCates()
    // case "list-t":
    //     op.ListTags()
    case "search":
        // leaf := parseArgs(os.Args)
        // op.Search(leaf.Category, leaf.Tags)
    case "remove":
        id := os.Args[2]
        fmt.Println("Are you sure to remove code segment with id("+id+")?", "  yes|no")
        var response string
        _, err := fmt.Scanln(&response)
        if err != nil {
            fmt.Println(err)
            os.Exit(-1)
        }

        if "YES" == strings.ToUpper(response) {
            op.Remove(id)
        }
    case "edit":
        id := os.Args[2]
        op.Edit(id)
    case "exec":
        file := os.Args[2]
        op.Exec(file)
    default:
        fmt.Println("Error: wrong path")
        os.Exit(2)
    }

    if op.err != nil {
        fmt.Println("error:", op.err)
    }
}

func checkRequiredArg(argName, argValue string) {
    if strings.TrimSpace(argValue) == "" {
        fmt.Println("Missing required arg: ", argName)
        os.Exit(2)
    }
}

// func parseArgs(args []string) Leaf {
//     var ind = func(s string) int {
//         for i, a := range args {
//             if a == s {
//                 return i
//             }
//         }
//         return -1
//     }

//     var argsLen = 2
//     var getParam = func(flag string) string {
//         if ind_flag := ind(flag); ind_flag > 0 {
//             //fmt.Printf("flag:%s, index:%d ", flag, ind_flag)
//             if len(args) <= ind_flag+1 {
//                 fmt.Println("missing parameter value for ", flag)
//             }
//             argsLen += 2
//             return args[ind_flag+1]
//         }
//         return ""
//     }

//     id := getParam("-i")
//     cate := getParam("-c")
//     tagStr := getParam("-t")
//     desc := getParam("-m")
//     var content string
//     if len(args) > argsLen {
//         content = strings.Join(args[argsLen:], " ")
//     }

//     if args[1] == "search" {
//         tagStr = strings.Join(args[argsLen:], ",")
//     }

//     return Leaf{Id: id, Category: cate, Tags: tagStr, Desc: desc, Content: content}
// }

func printUsage() {
    fmt.Printf("Usage: %s <command> <args...> \n", os.Args[0])
    fmt.Println("acceptable commands are:")
    for _, command := range subCommands {
        fmt.Printf("  %-8s  %s\n", command, subCommandMap[command])
    }

    fmt.Println("")
    fmt.Println("optional args:")
    flag.PrintDefaults()
}
