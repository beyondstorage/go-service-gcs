package gcs

import (
	"context"
	"fmt"
	"io"

	gs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/aos-dev/go-storage/v3/pkg/iowrap"
	. "github.com/aos-dev/go-storage/v3/types"
)

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	o = s.newObject(false)
	o.Mode = ModeRead
	o.ID = s.getAbsPath(path)
	o.Path = path
	return o
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	err = s.bucket.Object(rp).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	input := &objectPageStatus{
		prefix: s.getAbsPath(path),
	}

	var nextFn NextObjectFunc

	switch {
	case opt.ListMode.IsDir():
		input.delimiter = "/"
		nextFn = s.nextObjectPageByDir
	case opt.ListMode.IsPrefix():
		nextFn = s.nextObjectPageByPrefix
	default:
		return nil, fmt.Errorf("invalid list mode")
	}

	return NewObjectIterator(ctx, nextFn, input), nil
}

func (s *Storage) metadata(ctx context.Context, opt pairStorageMetadata) (meta *StorageMeta, err error) {
	meta = NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	return
}

func (s *Storage) nextObjectPageByDir(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	it := s.bucket.Objects(ctx, &gs.Query{
		Prefix:    input.prefix,
		Delimiter: input.delimiter,
	})

	remaining := 200
	for remaining > 0 {
		object, err := it.Next()
		if err == iterator.Done {
			return IterateDone
		}
		if err != nil {
			return err
		}

		// Prefix is set only for ObjectAttrs which represent synthetic "directory
		// entries" when iterating over buckets using Query.Delimiter. See
		// ObjectIterator.Next. When set, no other fields in ObjectAttrs will be
		// populated.
		if object.Prefix != "" {
			o := s.newObject(true)
			o.ID = object.Prefix
			o.Path = s.getRelPath(object.Prefix)
			o.Mode |= ModeDir

			page.Data = append(page.Data, o)

			remaining -= 1
			continue
		}

		o, err := s.formatFileObject(object)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
		remaining -= 1
	}

	return nil
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	it := s.bucket.Objects(ctx, &gs.Query{
		Prefix: input.prefix,
	})

	remaining := 200
	for remaining > 0 {
		object, err := it.Next()
		if err == iterator.Done {
			return IterateDone
		}
		if err != nil {
			return err
		}

		o, err := s.formatFileObject(object)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
		remaining -= 1
	}

	return nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)

	var rc io.ReadCloser

	object := s.bucket.Object(rp)
	if opt.HasSseCustomerKey {
		object = object.Key(opt.SseCustomerKey)
	}
	rc, err = object.NewReader(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		cerr := rc.Close()
		if cerr != nil {
			err = cerr
		}
	}()

	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}

	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	rp := s.getAbsPath(path)

	attr, err := s.bucket.Object(rp).Attrs(ctx)
	if err != nil {
		return nil, err
	}

	return s.formatFileObject(attr)
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	rp := s.getAbsPath(path)

	object := s.bucket.Object(rp)
	if opt.HasSseCustomerKey {
		object = object.Key(opt.SseCustomerKey)
	}
	w := object.NewWriter(ctx)
	defer func() {
		cerr := w.Close()
		if cerr != nil {
			err = cerr
		}
	}()

	w.Size = size
	if opt.HasContentMd5 {
		// FIXME: we need to check value's encoding type.
		w.MD5 = []byte(opt.ContentMd5)
	}
	if opt.HasStorageClass {
		w.StorageClass = opt.StorageClass
	}
	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	return io.Copy(w, r)
}
