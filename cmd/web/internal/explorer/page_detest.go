package explorer

import (
	"fmt"
	"math/rand"
	"net/http"
)

func (e *Explorer) testPage(rw http.ResponseWriter, r *http.Request) {
	data := e.frontData(rw, r).
		Css([]string{}).
		JS([]string{
			"https://unpkg.com/vuex@next",
			"https://unpkg.com/vue-router@next",
			"https://cdn.jsdelivr.net/npm/vue3-sfc-loader/dist/vue3-sfc-loader.js",
		}).
		ModuleJS([]string{
			"/static/test/main.js?v0.1." + fmt.Sprintf("%d", rand.Intn(10000)),
		}).
		Title("TestPageTitle")

	e.ExecuteTemplate(rw, r, "test", data)
}
