package declarative

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/a-h/templ"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Entity interface{}

type EntityView struct {
	db *gorm.DB
	r  Entity
}

func (g *EntityView) Headers() []string {
	headers := make([]string, 0)
	t := reflect.ValueOf(g.r).Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		headers = append(headers, f.Tag.Get("label"))
	}
	return headers
}

func (g *EntityView) Title() string {
	return reflect.ValueOf(g.r).Type().Name()
}

type Entry struct {
	Key   string
	Value string
}

func Fields(r Entity) []*Entry {
	values := make([]*Entry, 0)
	e := reflect.ValueOf(r)
	t := e.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		for _, a := range strings.Split(f.Tag.Get("gorm"), ";") {
			if strings.HasPrefix(a, "column") {
				values = append(values, &Entry{Key: strings.Split(a, ":")[1], Value: e.Field(i).String()})
			}
		}
	}
	return values
}

func (g *EntityView) Entities() []Entity {
	es := []Entity{}
	elemType := reflect.TypeOf(g.r)
	c := reflect.New(reflect.SliceOf(elemType))
	entitySlice := c.Interface()
	_ = g.db.Find(entitySlice)
	c = reflect.Indirect(c)
	for i := 0; i < c.Len(); i++ {
		val := c.Index(i).Interface()
		es = append(es, val.(Entity))
	}
	return es
}

func (e *EntityView) handler(db *gorm.DB, f func(db *gorm.DB, e Entity, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(db, e.r, w, r)
	}
}
func Setup(entities ...interface{}) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(entities...)

	for _, e := range entities {
		g := &EntityView{db: db, r: e}
		fmt.Println("Adding...", "/"+strings.ToLower(g.Title()))
		http.Handle("/"+strings.ToLower(g.Title()), templ.Handler(View(db, g)))
		http.Handle("/"+strings.ToLower(g.Title())+"/add", g.handler(db, addRow))
		http.Handle("/"+strings.ToLower(g.Title())+"/rem", g.handler(db, remRow))
	}
	return db
}
