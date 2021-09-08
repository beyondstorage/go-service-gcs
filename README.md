[![Build Status](https://github.com/beyondstorage/go-service-gcs/workflows/Unit%20Test/badge.svg?branch=master)](https://github.com/beyondstorage/go-service-gcs/actions?query=workflow%3A%22Unit+Test%22)
[![Integration Tests](https://teamcity.beyondstorage.io/app/rest/builds/buildType:(id:Services_Gcs_IntegrationTests)/statusIcon)](https://teamcity.beyondstorage.io/buildConfiguration/Services_Gcs_IntegrationTests)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/Xuanwo/storage/blob/master/LICENSE)
[![](https://img.shields.io/matrix/beyondstorage@go-storage:matrix.org.svg?logo=matrix)](https://matrix.to/#/#beyondstorage@go-storage:matrix.org)

# go-services-gcs

[Google Cloud Storage](https://cloud.google.com/storage/) service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get github.com/beyondstorage/go-service-gcs/v2
```

## Usage

```go
import (
	"log"

	_ "github.com/beyondstorage/go-service-gcs/v2"
	"github.com/beyondstorage/go-storage/v4/services"
)

func main() {
	store, err := services.NewStoragerFromString("gcs://bucket_name/path/to/workdir?credential=file:<absolute_path_to_token_file>&project_id=<google_cloud_project_id>")
	if err != nil {
		log.Fatal(err)
	}
	
	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/gcs) about go-service-gcs.
