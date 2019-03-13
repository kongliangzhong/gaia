package main

import (
    "fmt"
    "flag"
    "os"
    "os/user"
    "strings"
    "io/ioutil"
)

// TODO: execute code from system clipboard
// TODO: copy code to clipboard

// TODO: read all items by skip:count
// TODO: generate all as a pdf file.

// CARD NOTE
// card note chain:  exchange card.
// Link between note.

var gaiaDir = ".gaia/"
const dataFileName = "data/data.json"
var dataFilePath = ""
var codeBase = "codebase/"

var subCommands = []string{
    "add",
    "get",
    "alias",
    "append",
    "merge",
    "list",
    "search",
    "remove",
    "edit",
    "exec",
    "stats",
    "admin",
}

var subCommandMap = map[string]string{
    "add": "add item",
    "get": "get item by id",
    "alias": "add keyword alias",
    "append": "append text to item content",
    "merge": "merge two item into one",
    "list": "list items",
    "search": "search items",
    "remove": "remove item by id",
    "edit": "edit item in vi",
    "exec": "execute item",
    "stats": "stats info",
    "admin": "admin",
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
    executable bool
    mainFile string
    inputFile string

    listCategories bool
    listTags bool
    listAlias bool
    listNames bool
    countStats bool
    onlyContent bool
    isFormat bool
    isRemove bool
)

func init() {
    usr, err := user.Current()
    if err != nil {
        panic(err)
    }

    gaiaDir = usr.HomeDir + "/" + gaiaDir
    dataFilePath = gaiaDir + dataFileName
    codeBase = gaiaDir + codeBase
    _, err = os.Stat(gaiaDir)
    if err != nil && os.IsNotExist(err) {
        err = os.MkdirAll(gaiaDir, 0770)
        if err != nil {
            panic(err)
        }
    }
}

func main() {
    flag.BoolVar(&isHelp, "h", false, "show help message")
    if len(os.Args) == 1 {
        printUsage()
        os.Exit(-1)
    }

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
        subFlag.StringVar(&id, "i", "", "node id")
        subFlag.StringVar(&name, "n", "", "node name")
        subFlag.StringVar(&category, "c", "", "node category")
        subFlag.StringVar(&tags, "t", "", "node tags, tag seprated by comma")
        subFlag.StringVar(&desc, "d", "", "node description")
        subFlag.StringVar(&content, "b", "", "node body content")
        subFlag.BoolVar(&executable, "e", false, "is node executable")
        subFlag.StringVar(&mainFile, "m", "", "executable main file name")
        subFlag.StringVar(&inputFile, "f", "", "node body content input file")

        subFlag.Usage = func() {
            fmt.Printf("Usage: %s %s -n name -c category -b body [<other args>] \n", os.Args[0], os.Args[1])
            subFlag.PrintDefaults()
        }
    case "get":
        subFlag.StringVar(&id, "i", "", "node id")
        subFlag.BoolVar(&onlyContent, "c", false, "only print content")
    case "alias":
        subFlag.BoolVar(&isRemove, "r", false, "remove alias")
        subFlag.Usage = func() {
            fmt.Printf("Usage: %s %s <keyword> <target-keyword> \n", os.Args[0], os.Args[1])
            fmt.Printf("       %s %s -r <keyword> \n", os.Args[0], os.Args[1])
            subFlag.PrintDefaults()
        }
    case "append":
        subFlag.StringVar(&id, "i", "", "node id")
        subFlag.StringVar(&content, "b", "", "node append content")
    case "merge":
        subFlag.StringVar(&id, "i", "", "node id")
        subFlag.StringVar(&oid, "d", "", "dest node id")
    case "list":
        subFlag.BoolVar(&listCategories, "c", false, "list node categories")
        subFlag.BoolVar(&listTags, "t", false, "list node tags")
        subFlag.BoolVar(&listNames, "n", false, "list by name parts")
        subFlag.BoolVar(&listAlias, "a", false, "list global keyword alias")
    case "search":
        subFlag.StringVar(&category, "c", "", "search in certain category")
    case "remove":
        subFlag.StringVar(&id, "i", "", "node id")
    case "edit":
        subFlag.StringVar(&id, "i", "", "node id")
    case "exec":
        subFlag.StringVar(&id, "i", "", "node id")
    case "stats":
        subFlag.BoolVar(&countStats, "n", false, "count stats")
    case "admin":
        subFlag.BoolVar(&isFormat, "f", false, "format all data")
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
    // fmt.Println("dataFilePath:", dataFilePath)

    switch command {
    case "add":
        checkRequiredArg("-n", name)
        if executable {
            checkRequiredArg("-m", mainFile)
        }
        if content == "" {
            if inputFile == "" {
                fmt.Println("-b and -f can not both be empty")
                os.Exit(2)
            }

            contentBs, _ := ioutil.ReadFile(inputFile)
            content = string(contentBs)
        }
        if strings.TrimSpace(content) == "" {
            fmt.Println("node content is empty!")
            os.Exit(2)
        }

        node := Node{
            Name: name,
            Category: strings.Split(name, "-")[0],
            Tags: tags,
            Desc: desc,
            Content: content,
            Executable: executable,
            ExecFile: mainFile,
        }
        op.Add(node)
    case "get":
        if id == "" && len(subFlag.Args()) > 0 {
            id = subFlag.Args()[0]
        }
        op.Get(id, onlyContent)
    case "alias":
        aliasArgs := subFlag.Args()

        if isRemove {
            if len(aliasArgs) != 1 {
                subFlag.Usage()
                os.Exit(2)
            }
            op.RemoveAlias(aliasArgs[0])
        } else {
            if len(aliasArgs) != 2 {
                subFlag.Usage()
                os.Exit(2)
            }
            op.AddAlias(aliasArgs[0], aliasArgs[1])
        }
    case "append":
        extraContent := subFlag.Args()[0]
        op.Append(id, extraContent)
    case "merge":
        ids := subFlag.Args()
        op.Merge(ids)
    case "list":
        if listAlias {
            op.ListAlias()
        } else if listCategories {
            op.ListCates()
        } else if listTags {
            op.ListTags()
        } else if listNames {
            op.ListNodes(subFlag.Args())
        } else {
            subFlag.Usage()
            os.Exit(2)
        }
    case "search":
        op.Search(category, subFlag.Args())
    case "remove":
        if id == "" && len(subFlag.Args()) > 0 {
            id = subFlag.Args()[0]
        }

        if strings.TrimSpace(id) == "" {
            subFlag.Usage()
            os.Exit(2)
        }

        fmt.Println("Are you sure to remove node with id "+ id +"?", "  yes|no")
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
        // fmt.Println("id:", id)
        if id == "" && len(subFlag.Args()) > 0 {
            id = subFlag.Args()[0]
        }

        op.Edit(id)
    case "exec":
        file := os.Args[2]
        op.Exec(file)
    case "stats":
        op.Stats()
    case "admin":
        if isFormat {
            op.FormatData()
        }
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
