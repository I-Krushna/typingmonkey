package killswitch

import (
	"fmt"
	"os"
	"time"

	"typingmonkey/internal/input"
)

const (
	VK_CONTROL = 0x11
	VK_P       = 0x50
)

func Monitor() {
	for {
		if input.IsKeyPressed(VK_CONTROL) && input.IsKeyPressed(VK_P) {
			fmt.Println("\n\nKILLSWITCH! (╯°□°）╯︵ ┻━┻")
			os.Exit(0)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
