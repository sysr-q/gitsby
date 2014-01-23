package github

import (
	"fmt"
	"github.com/hoisie/web"
	_ "github.com/plausibility/gitsby/git"
)

func Hook(ctx *web.Context) {
	p := NewPayload([]byte(ctx.Params["payload"]))
	fmt.Println(">>>", p.Url())
}
