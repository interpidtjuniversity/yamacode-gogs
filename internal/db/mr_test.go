package db

import (
	"fmt"
	"gogs.io/gogs/internal/conf"
	"testing"
)

func Test_InsertMR(t *testing.T) {
	conf.Database.Type = "sqlite3"
	conf.Database.Path = "/root/yamaHub/yamaHub-database/yamaHub.db"

	x, _ := getEngine()
	mr := MergeRequest{SourceBranch: "dev", TargetBranch: "master", RepoId: 3, RepoName: "init",
		SourceCommitId: "407eecee5b0bf77b52d6bd2db84411694f629bd0", TargetCommitId: "cd8df412af13415f40a2119745b07857baa676f9"}
	if _, err := x.Insert(&mr); err!=nil{
		fmt.Print(err)
	} else {
		fmt.Print(fmt.Sprintf("mr id is %d", mr.ID))
	}
}
