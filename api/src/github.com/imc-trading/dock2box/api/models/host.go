package models

import "gopkg.in/mgo.v2/bson"

type Host struct {
	ID         bson.ObjectId `field:"id" json:"id,omitempty" bson:"_id"`
	Hostname   string        `field:"hostname" json:"hostname,omitempty" bson:"hostname"`
	Build      bool          `field:"build" json:"build" bson:"build"`
	Debug      bool          `field:"debug" json:"debug" bson:"debug"`
	GPT        bool          `field:"gpt" json:"gpt" bson:"gpt"`
	TagID      bson.ObjectId `field:"tagId" json:"tagId,omitempty" bson:"tagId"`
	Tag        *Tag          `json:"tag,omitempty"`
	Image      *Image        `json:"image,omitempty"`
	BootTag    *Tag          `json:"bootTag,omitempty"`
	BootImage  *Image        `json:"bootImage,omitempty"`
	KOpts      string        `field:"kOpts" json:"kOpts,omitempty" bson:"kOpts"`
	TenantID   bson.ObjectId `field:"tenantId" json:"tenantId,omitempty" bson:"tenantId"`
	Tenant     *Tenant       `json:"tenant,omitempty"`
	Labels     []string      `field:"labels" json:"labels,omitempty" bson:"labels"`
	SiteID     bson.ObjectId `field:"siteId" json:"siteId,omitempty" bson:"siteId"`
	Site       *Site         `json:"site,omitempty"`
	Interfaces *[]Interface  `json:"interfaces,omitempty"`
	Links      *[]Link       `json:"links,omitempty"`
}
