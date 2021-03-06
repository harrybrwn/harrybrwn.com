package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"harrybrown.com/pkg/log"
	"harrybrown.com/pkg/web"
)

// Debug cooresponds with the debug flag
var Debug = false

func init() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	BoolFlag(&Debug, "debug", "turn on debugging options")

	web.TemplateDir = "templates/"
	web.BaseTemplates = []string{"/index.html", "/nav.html"} // included in all pages
}

// Routes is a list of all the app's routes
var Routes = []web.Route{
	&web.Page{
		Title:     "Harry Brown",
		Template:  "pages/home.html",
		RoutePath: "/",
		RequestHook: func(self *web.Page, w http.ResponseWriter, r *http.Request) {
			self.Data = &struct{ Age string }{Age: getAge()}
		},
		HotReload: Debug,
	},
	&web.Page{
		Title:     "Freelancing",
		Template:  "pages/freelance.html",
		RoutePath: "/freelance",
	},
	&web.Page{
		Title:     "Resume",
		Template:  "pages/resume.html",
		RoutePath: "/resume",
		Data:      getResume("./static/data/resume.json"),
	},
	web.NewRoute("/static/", NewFileServer("static")), // handle file server
	web.NewNestedRoute("/api", apiroutes...).SetHandler(&web.JSONRoute{
		Static: func() interface{} { return info{Error: "Not implimented"} },
	}),
	// web.NewRoute("/github", http.RedirectHandler("https://github.com/harrybrwn", 301)),
}

var apiroutes = []web.Route{
	web.APIRoute("info", func(w http.ResponseWriter, r *http.Request) interface{} {
		return info{Age: time.Since(bday).Hours() / 24 / 365}
	}),
	web.StaticAPIRoute("testing", func() interface{} { return map[string]string{"testing": "testing the api"} }),
}

var bday = time.Date(1998, time.August, 4, 4, 0, 0, 0, time.UTC)

func getAge() string {
	age := time.Since(bday).Hours() / 24 / 365
	return fmt.Sprintf("%d", int(age))
}

func getResume(file string) *resumeContent {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return nil
	}
	c := &resumeContent{}

	if err = json.Unmarshal(b, c); err != nil {
		log.Println(err)
		return nil
	}
	return c
}

type resumeContent struct {
	Experience []resumeItem
	Education  []resumeItem
}

type resumeItem struct {
	Name, Title, Date, Content string
	BulletPoints               []string
}

type info struct {
	Age   float64 `json:"age,omitempty"`
	Error string  `json:"error,omitempty"`
}
