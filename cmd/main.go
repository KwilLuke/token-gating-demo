package main

import (
    "fmt"
    "os"

    "github.com/kwilteam/kwil-db/cmd/custom"
)

func main() {
    err := custom.NewCustomCmd(custom.CommonCmdConfig{
        RootCmd:     "custom-kwild",
        ProjectName: "Kwil Token Gating Example",
    }).Execute()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    os.Exit(0)
}