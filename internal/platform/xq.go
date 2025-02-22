package platform

import (
	"HelloGolang/pkg/common"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func XqId() string {
	var xqId string
	xqat, found := common.MyCache.Get("xqat")
	if found {
		xqId = xqat.(string)
	} else {
		cmd := exec.Command("python", "sele.py")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error:", err)
			return ""
		}
		fmt.Println(string(output))
		output = bytes.Replace(output, []byte(`'`), []byte(`"`), -1)
		xqId = strings.ReplaceAll(strings.ReplaceAll(string(output), "\r", ""), "\n", ";")
		common.MyCache.Set("xqat", xqId, 2*time.Hour)
		fmt.Println(xqId)
	}
	return xqId
}
