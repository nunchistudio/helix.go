package nats

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/nats-io/nats.go/jetstream"
)

/*
keyvalue implements the KeyValue interface and allows to wrap the NATS JetStream
key-value functions for automatic tracing and error recording.
*/
type keyvalue struct {
	bucket string
	store  jetstream.KeyValue
}

/*
KeyValue exposes an opinionated way to interact with a NATS JetStream key-value
store. All functions are wrapped with a context because some of them automatically
do distributed tracing (by using the said context) as well as error recording
within traces.
*/
type KeyValue interface {
	Bucket(ctx context.Context) string
	Get(ctx context.Context, key string) (jetstream.KeyValueEntry, error)
	GetRevision(ctx context.Context, key string, revision uint64) (jetstream.KeyValueEntry, error)
	Create(ctx context.Context, key string, value []byte) (uint64, error)
	Put(ctx context.Context, key string, value []byte) (uint64, error)
	Update(ctx context.Context, key string, value []byte, revision uint64) (uint64, error)
	Delete(ctx context.Context, key string, opts ...jetstream.KVDeleteOpt) error
	Purge(ctx context.Context, key string, opts ...jetstream.KVDeleteOpt) error
	PurgeDeletes(ctx context.Context, opts ...jetstream.KVPurgeOpt) error
	Watch(ctx context.Context, keys string, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error)
	WatchAll(ctx context.Context, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error)
	Keys(ctx context.Context, opts ...jetstream.WatchOpt) ([]string, error)
	History(ctx context.Context, key string, opts ...jetstream.WatchOpt) ([]jetstream.KeyValueEntry, error)
	Status(ctx context.Context) (jetstream.KeyValueStatus, error)
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
func (kv *keyvalue) Get(ctx context.Context, key string) (jetstream.KeyValueEntry, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Get", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get latest value for the key", err)
		}
	}()

	entry, err := kv.store.Get(ctx, key)
	setKeyValueAttributes(span, key, jetstream.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return entry, err
}

/*
GetRevision returns a specific revision value for the key.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) GetRevision(ctx context.Context, key string, revision uint64) (jetstream.KeyValueEntry, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / GetRevision", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get revision for the key", err)
		}
	}()

	entry, err := kv.store.GetRevision(ctx, key, revision)
	setKeyValueAttributes(span, key, jetstream.KeyValueConfig{
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
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Create", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to create key", err)
		}
	}()

	revision, err := kv.store.Create(ctx, key, value)
	setKeyValueAttributes(span, key, jetstream.KeyValueConfig{
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
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Put", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to put value for the key", err)
		}
	}()

	revision, err := kv.store.Put(ctx, key, value)
	setKeyValueAttributes(span, key, jetstream.KeyValueConfig{
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
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Update", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to update value for the key", err)
		}
	}()

	revision, err := kv.store.Update(ctx, key, value, last)
	setKeyValueAttributes(span, key, jetstream.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return revision, err
}

/*
Delete will place a delete marker and leave all revisions.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Delete(ctx context.Context, key string, opts ...jetstream.KVDeleteOpt) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Delete", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete key", err)
		}
	}()

	err = kv.store.Delete(ctx, key, opts...)
	setKeyValueAttributes(span, key, jetstream.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return err

}

/*
Purge will place a delete marker and remove all previous revisions.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Purge(ctx context.Context, key string, opts ...jetstream.KVDeleteOpt) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Purge", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to purge key", err)
		}
	}()

	err = kv.store.Purge(ctx, key, opts...)
	setKeyValueAttributes(span, key, jetstream.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return err
}

/*
PurgeDeletes will remove all current delete markers.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) PurgeDeletes(ctx context.Context, opts ...jetstream.KVPurgeOpt) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / PurgeDeletes", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to purge keys", err)
		}
	}()

	err = kv.store.PurgeDeletes(ctx, opts...)
	setKeyValueAttributes(span, "", jetstream.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return err
}

/*
Watch for any updates to keys that match the keys argument which could include
wildcards. Watch will send a nil entry when it has received all initial values.
*/
func (kv *keyvalue) Watch(ctx context.Context, keys string, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error) {
	return kv.store.Watch(ctx, keys, opts...)
}

/*
WatchAll will invoke the callback for all updates.
*/
func (kv *keyvalue) WatchAll(ctx context.Context, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error) {
	return kv.store.WatchAll(ctx, opts...)
}

/*
Keys will return all keys.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Keys(ctx context.Context, opts ...jetstream.WatchOpt) ([]string, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Keys", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get all keys", err)
		}
	}()

	keys, err := kv.store.Keys(ctx, opts...)
	setKeyValueAttributes(span, "", jetstream.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return keys, err
}

/*
History will return all historical values for the key.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) History(ctx context.Context, key string, opts ...jetstream.WatchOpt) ([]jetstream.KeyValueEntry, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / History", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get history for the key", err)
		}
	}()

	entries, err := kv.store.History(ctx, key, opts...)
	setKeyValueAttributes(span, key, jetstream.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return entries, err
}

/*
Status retrieves the status and configuration of a bucket.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Status(ctx context.Context) (jetstream.KeyValueStatus, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Status", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get bucket status", err)
		}
	}()

	status, err := kv.store.Status(ctx)
	setKeyValueAttributes(span, "", jetstream.KeyValueConfig{
		Bucket: kv.bucket,
	})

	return status, err
}
