package translate

import (
	"fmt"
	"testing"
)

func TestFreeGoogle(t *testing.T) {
	var freeGoogle FreeGoogle

	translate, err := freeGoogle.Translate(&TranslateReq{
		From: "auto",
		Text: []string{
			"你好",
			"你好世界！",
		},
		To: "en",
	})
	if err != nil {
		return
	}

	for _, result := range translate {
		fmt.Println(*result)
	}

	fmt.Println()
}
