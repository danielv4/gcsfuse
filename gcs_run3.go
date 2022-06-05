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
	"github.com/mauri870/gcsfs"
	"io/fs"
	"fmt"
	"google.golang.org/api/option"
	"cloud.google.com/go/storage"
	"context"
)


func main() {

	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile("gcp_service_file.json"))
	if err != nil {
		fmt.Println(err)
	}	
	
	gfs := gcsfs.NewWithClient(client, "drive-400")
	
	files, err := fs.ReadDir(gfs, ".")	
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}	

}
