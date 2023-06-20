package vault

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/hashicorp/vault/api"
)

/*
keyvalue implements the KeyValue interface and allows to wrap the Vault Key-Value
v2 functions for automatic tracing and error recording.
*/
type keyvalue struct {
	config    *Config
	mountpath string
	client    *api.KVv2
}

/*
KeyValue exposes an opinionated way to interact with Vault Key-Value v2. All
functions automatically handle distributed tracing as well as error recording
within traces.
*/
type KeyValue interface {
	Delete(ctx context.Context, secretpath string) error
	DeleteMetadata(ctx context.Context, secretpath string) error
	DeleteVersions(ctx context.Context, secretpath string, versions []int) error
	Destroy(ctx context.Context, secretpath string, versions []int) error
	Get(ctx context.Context, secretpath string) (*api.KVSecret, error)
	GetMetadata(ctx context.Context, secretpath string) (*api.KVMetadata, error)
	GetVersion(ctx context.Context, secretpath string, version int) (*api.KVSecret, error)
	GetVersionsAsList(ctx context.Context, secretpath string) ([]api.KVVersionMetadata, error)
	Patch(ctx context.Context, secretpath string, data map[string]any) (*api.KVSecret, error)
	PatchMetadata(ctx context.Context, secretpath string, metadata api.KVMetadataPatchInput) error
	Put(ctx context.Context, secretpath string, data map[string]any) (*api.KVSecret, error)
	PutMetadata(ctx context.Context, secretpath string, metadata api.KVMetadataPutInput) error
	Rollback(ctx context.Context, secretpath string, toVersion int) (*api.KVSecret, error)
	Undelete(ctx context.Context, secretpath string, versions []int) error
}

/*
Delete deletes the most recent version of a secret from the KV v2 secrets engine.
To delete an older version, use DeleteVersions.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Delete(ctx context.Context, secretpath string) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Delete", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete latest version", err)
		}
	}()

	err = kv.client.Delete(ctx, secretpath)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return err
}

/*
DeleteMetadata deletes all versions and metadata of the secret at the given path.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) DeleteMetadata(ctx context.Context, secretpath string) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / DeleteMetadata", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete metadata", err)
		}
	}()

	err = kv.client.DeleteMetadata(ctx, secretpath)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return err
}

/*
DeleteVersions deletes the specified versions of a secret from the KV v2 secrets
engine. To delete the latest version of a secret, just use Delete.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) DeleteVersions(ctx context.Context, secretpath string, versions []int) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / DeleteVersions", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete versions", err)
		}
	}()

	err = kv.client.DeleteVersions(ctx, secretpath, versions)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return err
}

/*
Destroy permanently removes the specified secret versions' data from the Vault
server. If no secret exists at the given path, no action will be taken.

A list of existing versions can be retrieved using the GetVersionsAsList method.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Destroy(ctx context.Context, secretpath string, versions []int) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Destroy", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to destroy versions", err)
		}
	}()

	err = kv.client.Destroy(ctx, secretpath, versions)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return err
}

/*
Get returns the latest version of a secret from the KV v2 secrets engine.

If the latest version has been deleted, an error will not be thrown, but the Data
field on the returned secret will be nil, and the Metadata field will contain the
deletion time.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Get(ctx context.Context, secretpath string) (*api.KVSecret, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Get", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get secret", err)
		}
	}()

	s, err := kv.client.Get(ctx, secretpath)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return s, err
}

/*
GetMetadata returns the full metadata for a given secret, including a map of its
existing versions and their respective creation/deletion times, etc.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) GetMetadata(ctx context.Context, secretpath string) (*api.KVMetadata, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / GetMetadata", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get metadata", err)
		}
	}()

	md, err := kv.client.GetMetadata(ctx, secretpath)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return md, err
}

/*
GetVersion returns the data and metadata for a specific version of the given
secret.

If that version has been deleted, the Data field on the returned secret will be
nil, and the Metadata field will contain the deletion time.

GetVersionsAsList can provide a list of available versions sorted by version number,
while the response from GetMetadata contains them as a map.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) GetVersion(ctx context.Context, secretpath string, version int) (*api.KVSecret, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / GetVersion", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get latest version", err)
		}
	}()

	s, err := kv.client.GetVersion(ctx, secretpath, version)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return s, err
}

/*
GetVersionsAsList returns a subset of the metadata for each version of the secret,
sorted by version number.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) GetVersionsAsList(ctx context.Context, secretpath string) ([]api.KVVersionMetadata, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / GetVersionsAsList", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to get versions", err)
		}
	}()

	versions, err := kv.client.GetVersionsAsList(ctx, secretpath)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return versions, err
}

/*
Patch additively updates the most recent version of a key-value secret,
differentiating it from Put which will fully overwrite the previous data. Only
the key-value pairs that are new or changing need to be provided.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Patch(ctx context.Context, secretpath string, data map[string]any) (*api.KVSecret, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Patch", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to patch secret", err)
		}
	}()

	s, err := kv.client.Patch(ctx, secretpath, data)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return s, err
}

/*
PatchMetadata can be used to replace just a subset of a secret's metadata fields
at a time, as opposed to PutMetadata which is used to completely replace all fields
on the previous metadata.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) PatchMetadata(ctx context.Context, secretpath string, metadata api.KVMetadataPatchInput) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / PatchMetadata", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to patch metadata", err)
		}
	}()

	err = kv.client.PatchMetadata(ctx, secretpath, metadata)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return err
}

/*
Put inserts a key-value secret into the KV v2 secrets engine.

If the secret already exists, a new version will be created and the previous
version can be accessed with the GetVersion method. GetMetadata can provide a
list of available versions.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Put(ctx context.Context, secretpath string, data map[string]any) (*api.KVSecret, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Put", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to put secret", err)
		}
	}()

	s, err := kv.client.Put(ctx, secretpath, data)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return s, err
}

/*
PutMetadata can be used to fully replace a subset of metadata fields for a given
KV v2 secret. All fields will replace the corresponding values on the Vault server.
Any fields left as nil will reset the field on the Vault server back to its zero
value.

To only partially replace the values of these metadata fields, use PatchMetadata.

This method can also be used to create a new secret with just metadata and no
secret data yet.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) PutMetadata(ctx context.Context, secretpath string, metadata api.KVMetadataPutInput) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / PutMetadata", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to put metadata", err)
		}
	}()

	err = kv.client.PutMetadata(ctx, secretpath, metadata)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return err
}

/*
Rollback can be used to roll a secret back to a previous non-deleted / non-destroyed
version. That previous version becomes the next / newest version for the path.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Rollback(ctx context.Context, secretpath string, toVersion int) (*api.KVSecret, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Rollback", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to rollback secret", err)
		}
	}()

	s, err := kv.client.Rollback(ctx, secretpath, toVersion)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return s, err
}

/*
Undelete undeletes the given versions of a secret, restoring the data so that it
can be fetched again.

A list of existing versions can be retrieved using the GetVersionsAsList method.

It automatically handles tracing and error recording.
*/
func (kv *keyvalue) Undelete(ctx context.Context, secretpath string, versions []int) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Key-Value / Undelete", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to undelete versions", err)
		}
	}()

	err = kv.client.Undelete(ctx, secretpath, versions)
	setDefaultAttributes(span, kv.config)
	setKeyValueAttributes(span, kv.mountpath, secretpath)

	return err
}
