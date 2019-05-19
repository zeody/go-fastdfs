package go_fastdfs

import (
	"fmt"
	"testing"
)

func TestDefault(t *testing.T) {
	client := Default()
	file := client.uploadFile("NO20180913155948.gif", "gif", nil)
	fmt.Printf("%v", file)
}
