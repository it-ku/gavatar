# gavatar

# example

```
package main

import (
	"fmt"
	"github.com/it-ku/gavatar"
	"image/color"
)

var colors = []uint32{
	0xff6200, 0x42c58e, 0x5a8de1, 0x785fe0,
}
func main(){
	ag := gavatar.NewAvatarGenerate("runtime/fonts/JetBrainsMono-Bold.ttf")
	ag.SetBackgroundColorHex(colors[3])
	ag.SetFrontColor(color.White)
	ag.SetFontSize(64)
	if err := ag.GenerateImage("O", "./outCn.png"); err != nil {
		fmt.Println(err)
		return
	}
}
```
