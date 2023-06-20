package nats

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/nats-io/nats.go"
)

/*
keyvalue implements the KeyValue interface and allows to wrap the NATS JetStream
key-value functions for automatic tracing and error recording.
*/
type keyvalue struct {
	bucket string
	store  nats.KeyValue
}

/*
KeyValue exposes an opinionated way to interact with a NATS JetStream key-value
store. All functions are wrapped with a context because some of them automatically
do distributed tracing (by using the said context) as well as error recording
within traces.
*/
type KeyValue interface {
	Bucket(ctx context.Context) string
	Get(ctx context.Context, key string) (nats.KeyValueEntry, error)
	GetRevision(ctx context.Context, key string, revision uint64) (nats.KeyValueEntry, error)
	Create(ctx context.Context, key string, value []byte) (uint64, error)
	Put(ctx context.Context, key string, value []byte) (uint64, error)
	Update(ctx context.Context, key string, value []byte, last uint64) (uint64, error)
	Delete(ctx context.Context, key string, opts ...nats.DeleteOpt) error
	Purge(ctx context.Context, key string, opts ...nats.DeleteOpt) error
	PurgeDeletes(ctx context.Context, opts ...nats.PurgeOpt) error
	Watch(ctx context.Context, keys string, opts ...nats.WatchOpt) (nats.KeyWatcher, error)
	WatchAll(ctx context.Context, opts ...nats.WatchOpt) (nats.KeyWatcher, error)
	Keys(ctx context.Context, opts ...nats.WatchOpt) ([]string, error)
	History(ctx context.Context, key string, opts ...nats.WatchOpt) ([]nats.KeyValueEntry, error)
	Status(ctx context.Context) (nats.KeyValueStatus, error)
}

/*
Bucket returns the current bucket name.
*/
func (kv *keyvalue) Bucket(ctx context.Context) string {
	return kv.store.Bucket()
}

/*
Get returns the latest value for the key.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Get(ctx context.Context, key string) (nats.KeyValueEntry, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Get", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get latest value for the key", err)
		}
	}()

	entry, err := kv.store.Get(key)
	setKeyValueAttributes(span, key, &nats.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return entry, err
}

/*
GetRevision returns a specific revision value for the key.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) GetRevision(ctx context.Context, key string, revision uint64) (nats.KeyValueEntry, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / GetRevision", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get revision for the key", err)
		}
	}()

	entry, err := kv.store.GetRevision(key, revision)
	setKeyValueAttributes(span, key, &nats.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return entry, err
}

/*
Create will add the key/value pair if it does not exist. Returns the revision
associated.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Create(ctx context.Context, key string, value []byte) (uint64, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Create", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to create key", err)
		}
	}()

	revision, err := kv.store.Create(key, value)
	setKeyValueAttributes(span, key, &nats.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return revision, err
}

/*
Put will place the new value for the key into the store. Returns the revision
associated.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Put(ctx context.Context, key string, value []byte) (uint64, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Put", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to put value for the key", err)
		}
	}()

	revision, err := kv.store.Put(key, value)
	setKeyValueAttributes(span, key, &nats.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return revision, err
}

/*
Update will update the value if the latest revision matches. Returns the revision
associated.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Update(ctx context.Context, key string, value []byte, last uint64) (uint64, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Update", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to update value for the key", err)
		}
	}()

	revision, err := kv.store.Update(key, value, last)
	setKeyValueAttributes(span, key, &nats.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return revision, err
}

/*
Delete will place a delete marker and leave all revisions.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Delete(ctx context.Context, key string, opts ...nats.DeleteOpt) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Delete", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete key", err)
		}
	}()

	err = kv.store.Delete(key, opts...)
	setKeyValueAttributes(span, key, &nats.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return err

}

/*
Purge will place a delete marker and remove all previous revisions.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Purge(ctx context.Context, key string, opts ...nats.DeleteOpt) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Purge", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to purge key", err)
		}
	}()

	err = kv.store.Purge(key, opts...)
	setKeyValueAttributes(span, key, &nats.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return err
}

/*
PurgeDeletes will remove all current delete markers.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) PurgeDeletes(ctx context.Context, opts ...nats.PurgeOpt) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / PurgeDeletes", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to purge keys", err)
		}
	}()

	err = kv.store.PurgeDeletes(opts...)
	setKeyValueAttributes(span, "", &nats.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return err
}

/*
Watch for any updates to keys that match the keys argument which could include
wildcards. Watch will send a nil entry when it has received all initial values.
*/
func (kv *keyvalue) Watch(ctx context.Context, keys string, opts ...nats.WatchOpt) (nats.KeyWatcher, error) {
	return kv.store.Watch(keys, opts...)
}

/*
WatchAll will invoke the callback for all updates.
*/
func (kv *keyvalue) WatchAll(ctx context.Context, opts ...nats.WatchOpt) (nats.KeyWatcher, error) {
	return kv.store.WatchAll(opts...)
}

/*
Keys will return all keys.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Keys(ctx context.Context, opts ...nats.WatchOpt) ([]string, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Keys", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get all keys", err)
		}
	}()

	keys, err := kv.store.Keys(opts...)
	setKeyValueAttributes(span, "", &nats.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return keys, err
}

/*
History will return all historical values for the key.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) History(ctx context.Context, key string, opts ...nats.WatchOpt) ([]nats.KeyValueEntry, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / History", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get history for the key", err)
		}
	}()

	entries, err := kv.store.History(key, opts...)
	setKeyValueAttributes(span, key, &nats.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return entries, err
}

/*
Status retrieves the status and configuration of a bucket.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Status(ctx context.Context) (nats.KeyValueStatus, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Status", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get bucket status", err)
		}
	}()

	status, err := kv.store.Status()
	setKeyValueAttributes(span, "", &nats.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return status, err
}
