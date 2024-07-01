package matcher

import (
	// _ "MSBot/matcher/feature/at"
	_ "MSBot/matcher/feature/mofish"
	_ "MSBot/matcher/feature/roll"
	_ "MSBot/matcher/feature/sleep"

	_ "MSBot/matcher/maplestory/character"
	_ "MSBot/matcher/maplestory/help"
	_ "MSBot/matcher/maplestory/maintainer"

	"fmt"
)

func init() {
	fmt.Println("matcher init done ...")
}
