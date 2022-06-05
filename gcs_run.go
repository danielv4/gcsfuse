/*
 * gcs_run.go
 *
 * Copyright 2022 Daniel Vanderloo
 */
/*
 * This file is part of Cgofuse.
 *
 * It is licensed under the MIT license. The full license text can be found
 * in the License.txt file at the root of this project.
 */

package main

import (
	"context"
    "fmt"
	//"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	
	"encoding/json"
	"os"
	//"path"
)





type Credentials struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}


type Gcs struct {
	
	client      *storage.Client
	bucket      string 
	credentials Credentials
	ctx         context.Context
}


type GcsFileObject struct {
	
	IsDir         bool
	Name          string
	Size          int
	LastModified  time.Time
}







func (self *Gcs) ReadDir(dirname string) ([]GcsFileObject, error) {

	var arr []GcsFileObject
	var err error
	
	
	input := storage.Query{}
	input.Delimiter = "/"
	if dirname != "/" {
		input.Prefix = dirname[1:] + "/"
	}
	
	it := self.client.Bucket(self.bucket).Objects(self.ctx, &input)
	
	
	//objs, err := self.client.Bucket(self.bucket).List(self.ctx, &input)
    //if err != nil {
    //    return arr, err
    //}
	
	fmt.Printf("%+v\n", it)
	

	
	return arr, err
}


func NewClient(bucketName string, credentialsFile string) (*Gcs, error) {
	
	var credentials Credentials
	gs := new(Gcs)
	ctx := context.Background()
	
	fp, err := os.Open(credentialsFile)
    if err != nil {
        return gs, err
    }

    parser := json.NewDecoder(fp)
    if err = parser.Decode(&credentials); err != nil {
        return gs, err
    }
	
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("gcp_service_file.json"))
	if err != nil {
		return gs, err
	}	
	
	
	gs.credentials = credentials
	gs.bucket = bucketName
	gs.ctx = ctx
	gs.client = client
	

	return gs, err
}



func main() {

	gs, err := NewClient("bucket2", "gcp_service_file.json")
	if err != nil {
		fmt.Println(err)
	}	
	

	entries, err := gs.ReadDir("/")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", entries)

	
	
	it := gs.client.Buckets(gs.ctx, gs.credentials.ProjectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("unable to list Buckets")
		}
		fmt.Printf("Bucket: %v\n", battrs.Name)
	}		
	

}
