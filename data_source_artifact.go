package main

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
)

//type FileInfo struct {
//	uri         string
//	downloadUri string
//	repo        string
//	path        string
//	size        string
//	checksums   struct {
//		md5    string
//		sha1   string
//		sha256 string
//	}
//	originalChecksums struct {
//		md5    string
//		sha1   string
//		sha256 string
//	}
//}

func dataSourceArtifact() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArtifactRead,

		Schema: map[string]*schema.Schema{
			"repository_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"checksums": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceArtifactRead(d *schema.ResourceData, m interface{}) error {
	return resourceArtifactRead(d, m)
}

func resourceArtifactRead(d *schema.ResourceData, m interface{}) error {
	resp, err := http.Get("https://artifacts.amfamlabs.com/api/storage/lambda/propinc/ingest/replicate-2.30.0.zip")
	if err != nil {
		log.Fatal(err)
		//
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//var f FileInfo
	//err := json.NewDecoder(resp.Body).Decode(&f)
	var v map[string]interface{}
	//jsonerr := json.NewDecoder(resp.Body).Decode(&v)
	jsonerr := json.Unmarshal(body, &v)
	if jsonerr != nil {
		log.Fatal(jsonerr)
	}
	for k, vv := range v {
		switch k {
		case "repo":
			d.Set("repo", vv)
		case "path":
			d.Set("path", vv)
			//default:
			//	return
			//
		default:
		}
	}
	//d.Set("repo", f.repo)
	//d.Set("path", f.path)
	return nil

}
