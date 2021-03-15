package iterativeE

import (
	"fmt"
	"github.com/gogs/git-module"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RepoDiff(t *testing.T) {
	gitRepo, _ := git.Open("/root/yamaHub/yamaHub-repositories/interpidtjuniversity/init.git")

	diff, err := GetDiff("407eecee5b0bf77b52d6bd2db84411694f629bd0","cd8df412af13415f40a2119745b07857baa676f9", gitRepo)
	assert.Nil(t, err)
	fmt.Print(diff.Files[0].Sections[0].Lines[0].Content)
	fmt.Print("\n")
	fmt.Print(diff.Files[0].Sections[0].Lines[0].LeftLine)
	fmt.Print("\n")
	fmt.Print(diff.Files[0].Sections[0].Lines[0].RightLine)
	fmt.Print("\n")
}

