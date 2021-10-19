# go-services-gcs

[Google Cloud Storage](https://cloud.google.com/storage/) service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Notes

**This package has been moved to [go-storage](https://github.com/beyondstorage/go-storage/tree/master/services/gcs).**

```shell
go get go.beyondstorage.io/services/gcs/v3
```

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
