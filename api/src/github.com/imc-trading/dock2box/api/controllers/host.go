package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/imc-trading/dock2box/api/models"
)

type HostController struct {
	database  string
	schemaURI string
	session   *mgo.Session
	baseURI   string
	envelope  bool
	hateoas   bool
}

func NewHostController(s *mgo.Session, b string, e bool, h bool) *HostController {
	return &HostController{
		database:  "d2b",
		schemaURI: "file://schemas/host.json",
		session:   s,
		baseURI:   b,
		envelope:  e,
		hateoas:   h,
	}
}

func (c *HostController) SetDatabase(database string) {
	c.database = database
}

func (c *HostController) SetSchemaURI(uri string) {
	c.schemaURI = uri + "/host.json"
}

func (c *HostController) CreateIndex() {
	index := mgo.Index{
		Key:    []string{"host"},
		Unique: true,
	}

	if err := c.session.DB(c.database).C("hosts").EnsureIndex(index); err != nil {
		panic(err)
	}
}

func (c *HostController) All(w http.ResponseWriter, r *http.Request) {
	// Get allowed key names
	es := reflect.ValueOf(models.Host{})
	keys, _ := structTags(es, "field", "bson")
	kinds, _ := structKinds(es, "field")

	// Query
	cond := bson.M{}
	qry := r.URL.Query()
	for k, v := range qry {
		if k == "envelope" || k == "embed" || k == "sort" || k == "hateoas" || k == "fields" {
			continue
		}
		//		fmt.Println(k, kinds[k]
		if _, ok := keys[k]; !ok {
			jsonError(w, r, fmt.Sprintf("Incorrect key used in query: %s", k), http.StatusBadRequest, c.envelope)
			return
		} else if kinds[k] == reflect.Bool {
			switch strings.ToLower(v[0]) {
			case "true":
				cond[keys[k]] = true
			case "false":
				cond[keys[k]] = false
			default:
				jsonError(w, r, fmt.Sprintf("Incorrect value for key used in query needs to be either true/false: %s", k), http.StatusBadRequest, c.envelope)
				return
			}
		} else if bson.IsObjectIdHex(v[0]) {
			cond[keys[k]] = bson.ObjectIdHex(v[0])
		} else {
			cond[keys[k]] = v[0]
		}
	}

	// Sort
	sort := []string{}
	if _, ok := qry["sort"]; ok {
		for _, str := range strings.Split(qry["sort"][0], ",") {
			op := ""
			k := str
			switch str[0] {
			case '+':
				op = "+"
				k = str[1:len(str)]
			case '-':
				op = "-"
				k = str[1:len(str)]
			}
			if _, ok := keys[k]; !ok {
				jsonError(w, r, fmt.Sprintf("Incorrect key used in sort: %s", k), http.StatusBadRequest, c.envelope)
				return
			}
			sort = append(sort, op+keys[k])
		}
	}

	// Fields
	fields := bson.M{}
	if _, ok := qry["fields"]; ok {
		for _, k := range strings.Split(qry["fields"][0], ",") {
			if _, ok := keys[k]; !ok {
				jsonError(w, r, fmt.Sprintf("Incorrect key used in fields: %s", k), http.StatusBadRequest, c.envelope)
				return
			}
			fields[keys[k]] = 1
		}
	}

	// Initialize empty struct list
	s := []models.Host{}

	// Get all entries
	if err := c.session.DB(c.database).C("hosts").Find(cond).Sort(sort...).Select(fields).All(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Embed related data
	if r.URL.Query().Get("embed") == "true" {
		for i, v := range s {
			// Get Tag
			if err := c.session.DB(c.database).C("tags").FindId(v.TagID).One(&s[i].Tag); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			// Get Tenant
			if err := c.session.DB(c.database).C("tenants").FindId(v.TenantID).One(&s[i].Tenant); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			// Get Site
			if err := c.session.DB(c.database).C("sites").FindId(v.SiteID).One(&s[i].Site); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
	}

	// HATEOAS Links
	hateoas := c.hateoas
	switch strings.ToLower(r.URL.Query().Get("hateoas")) {
	case "true":
		hateoas = true
	case "false":
		hateoas = false
	}
	if hateoas == true {
		for i, v := range s {
			links := []models.Link{}

			if v.TagID != "" {
				links = append(links, models.Link{
					HRef:   c.baseURI + "/tags/" + v.TagID.Hex(),
					Rel:    "self",
					Method: "GET",
				})
			}

			if v.TenantID != "" {
				links = append(links, models.Link{
					HRef:   c.baseURI + "/tenants/" + v.TenantID.Hex(),
					Rel:    "self",
					Method: "GET",
				})
			}

			if v.SiteID != "" {
				links = append(links, models.Link{
					HRef:   c.baseURI + "/sites/" + v.SiteID.Hex(),
					Rel:    "self",
					Method: "GET",
				})
			}

			s[i].Links = &links
		}
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}

func (c *HostController) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get object id
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.Host{}

	// Get entry
	if err := c.session.DB(c.database).C("hosts").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Embed related data
	if r.URL.Query().Get("embed") == "true" {
		// Get Tag
		if err := c.session.DB(c.database).C("tags").FindId(s.TagID).One(&s.Tag); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Get Tenant
		if err := c.session.DB(c.database).C("tenants").FindId(s.TenantID).One(&s.Tenant); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Get Site
		if err := c.session.DB(c.database).C("sites").FindId(s.SiteID).One(&s.Site); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	// HATEOAS Links
	hateoas := c.hateoas
	switch strings.ToLower(r.URL.Query().Get("hateoas")) {
	case "true":
		hateoas = true
	case "false":
		hateoas = false
	}
	if hateoas == true {
		links := []models.Link{}

		if s.TagID != "" {
			links = append(links, models.Link{
				HRef:   c.baseURI + "/tags/" + s.TagID.Hex(),
				Rel:    "self",
				Method: "GET",
			})
		}

		if s.TenantID != "" {
			links = append(links, models.Link{
				HRef:   c.baseURI + "/tenants/" + s.TenantID.Hex(),
				Rel:    "self",
				Method: "GET",
			})
		}

		if s.SiteID != "" {
			links = append(links, models.Link{
				HRef:   c.baseURI + "/sites/" + s.SiteID.Hex(),
				Rel:    "self",
				Method: "GET",
			})
		}

		s.Links = &links
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}

func (c *HostController) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct
	s := models.Host{}

	// Decode JSON into struct
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		jsonError(w, r, "Failed to deconde JSON: "+err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	// Set ID
	s.ID = bson.NewObjectId()

	// Validate input using JSON Schema
	docLoader := gojsonschema.NewGoLoader(s)
	schemaLoader := gojsonschema.NewReferenceLoader(c.schemaURI)

	res, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	if !res.Valid() {
		var errors []string
		for _, e := range res.Errors() {
			errors = append(errors, fmt.Sprintf("%s: %s", e.Context().String(), e.Description()))
		}
		jsonError(w, r, errors, http.StatusInternalServerError, c.envelope)
		return
	}

	// Insert entry
	if err := c.session.DB(c.database).C("hosts").Insert(s); err != nil {
		jsonError(w, r, "Insert: "+err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusCreated, c.envelope)
}

func (c *HostController) Update(w http.ResponseWriter, r *http.Request) {
	// Get ID
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get object id
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.Host{}

	// Decode JSON into struct
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		jsonError(w, r, "Failed to deconde JSON: "+err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	// Validate input using JSON Schema
	docLoader := gojsonschema.NewGoLoader(s)
	schemaLoader := gojsonschema.NewReferenceLoader(c.schemaURI)

	res, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		jsonError(w, r, "Failed to load schema: "+err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	if !res.Valid() {
		var errors []string
		for _, e := range res.Errors() {
			errors = append(errors, fmt.Sprintf("%s: %s", e.Context().String(), e.Description()))
		}
		jsonError(w, r, errors, http.StatusInternalServerError, c.envelope)
		return
	}

	// Update entry
	if err := c.session.DB(c.database).C("hosts").UpdateId(oid, s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}

func (c *HostController) Delete(w http.ResponseWriter, r *http.Request) {
	// Get ID
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get object id
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.Host{}

	// Get entry
	if err := c.session.DB(c.database).C("hosts").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Remove entry
	if err := c.session.DB(c.database).C("hosts").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}
