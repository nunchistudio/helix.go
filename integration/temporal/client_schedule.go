package temporal

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/telemetry/trace"

	"go.temporal.io/sdk/client"
)

/*
scheduleclient implements the ScheduleClient interface and allows to wrap the
Temporal scheduling functions for automatic tracing and error recording.
*/
type scheduleclient struct {
	config *Config
	client client.ScheduleClient
}

/*
ScheduleClient exposes an opinionated way to interact with Temporal scheduling
capabilities.
*/
type ScheduleClient interface {
	Create(ctx context.Context, options client.ScheduleOptions) (ScheduleHandle, error)
	List(ctx context.Context, options client.ScheduleListOptions) (client.ScheduleListIterator, error)
	Handle(ctx context.Context, scheduleID string) ScheduleHandle
}

/*
Create creates a new workflow schedule.

It automatically handles tracing and error recording.
*/
func (sc *scheduleclient) Create(ctx context.Context, opts client.ScheduleOptions) (ScheduleHandle, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Schedule / Create", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to create schedule", err)
		}
	}()

	h, err := sc.client.Create(ctx, opts)
	setDefaultAttributes(span, sc.config)

	sh := &schedulehandle{
		config:  sc.config,
		handler: h,
	}

	if h != nil {
		sh.id = h.GetID()
		setScheduleAttributes(span, sh.id)
	}

	return sh, err
}

/*
List returns an interator to list all schedules.

It automatically handles tracing and error recording.
*/
func (sc *scheduleclient) List(ctx context.Context, opts client.ScheduleListOptions) (client.ScheduleListIterator, error) {
	return sc.client.List(ctx, opts)
}

/*
Handle returns a schedule handler, allowing to manage a worfklow schedule.

It automatically handles tracing and error recording.
*/
func (sc *scheduleclient) Handle(ctx context.Context, scheduleID string) ScheduleHandle {
	h := sc.client.GetHandle(ctx, scheduleID)
	sh := &schedulehandle{
		id:      scheduleID,
		config:  sc.config,
		handler: h,
	}

	return sh
}

/*
schedulehandle implements the ScheduleHandle interface and allows to wrap the
Temporal scheduling functions for automatic tracing and error recording.
*/
type schedulehandle struct {
	id      string
	config  *Config
	handler client.ScheduleHandle
}

/*
ScheduleHandle exposes an opinionated way to interact with Temporal scheduling
capabilities.
*/
type ScheduleHandle interface {
	GetID(ctx context.Context) string
	Delete(ctx context.Context) error
	Backfill(ctx context.Context, options client.ScheduleBackfillOptions) error
	Update(ctx context.Context, options client.ScheduleUpdateOptions) error
	Describe(ctx context.Context) (*client.ScheduleDescription, error)
	Trigger(ctx context.Context, options client.ScheduleTriggerOptions) error
	Pause(ctx context.Context, options client.SchedulePauseOptions) error
	Unpause(ctx context.Context, options client.ScheduleUnpauseOptions) error
}

/*
GetID returns the schedule ID asssociated with this handle.

It automatically handles tracing and error recording.
*/
func (sh *schedulehandle) GetID(ctx context.Context) string {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Schedule / GetID", humanized))
	defer span.End()

	id := sh.handler.GetID()
	setDefaultAttributes(span, sh.config)
	setScheduleAttributes(span, sh.id)

	return id
}

/*
Delete deletes the schedule.

It automatically handles tracing and error recording.
*/
func (sh *schedulehandle) Delete(ctx context.Context) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Schedule / Delete", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete schedule", err)
		}
	}()

	err = sh.handler.Delete(ctx)
	setDefaultAttributes(span, sh.config)
	setScheduleAttributes(span, sh.id)

	return err
}

/*
Backfill backfills the schedule by going though the specified time periods and
taking Actions as if that time passed by right now, all at once.

It automatically handles tracing and error recording.
*/
func (sh *schedulehandle) Backfill(ctx context.Context, opts client.ScheduleBackfillOptions) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Schedule / Backfill", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to backfill schedule", err)
		}
	}()

	err = sh.handler.Backfill(ctx, opts)
	setDefaultAttributes(span, sh.config)
	setScheduleAttributes(span, sh.id)

	return err
}

/*
Update updates the schedule. If two Update calls are made in parallel to the same
Schedule there is the potential for a race condition.

It automatically handles tracing and error recording.
*/
func (sh *schedulehandle) Update(ctx context.Context, opts client.ScheduleUpdateOptions) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Schedule / Update", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to update schedule", err)
		}
	}()

	err = sh.handler.Update(ctx, opts)
	setDefaultAttributes(span, sh.config)
	setScheduleAttributes(span, sh.id)

	return err
}

/*
Describe fetches the schedule's description from the server.

It automatically handles tracing and error recording.
*/
func (sh *schedulehandle) Describe(ctx context.Context) (*client.ScheduleDescription, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Schedule / Describe", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to describe schedule", err)
		}
	}()

	desc, err := sh.handler.Describe(ctx)
	setDefaultAttributes(span, sh.config)
	setScheduleAttributes(span, sh.id)

	return desc, err
}

/*
Trigger triggers an Action to be taken immediately. Will override the schedules
default policy with the one specified here. If overlap is SCHEDULE_OVERLAP_POLICY_UNSPECIFIED
the schedule policy will be used.

It automatically handles tracing and error recording.
*/
func (sh *schedulehandle) Trigger(ctx context.Context, opts client.ScheduleTriggerOptions) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Schedule / Trigger", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to trigger schedule", err)
		}
	}()

	err = sh.handler.Trigger(ctx, opts)
	setDefaultAttributes(span, sh.config)
	setScheduleAttributes(span, sh.id)

	return err
}

/*
Pause pauses the schedule. It will also overwrite the schedules current note with
the new note.

It automatically handles tracing and error recording.
*/
func (sh *schedulehandle) Pause(ctx context.Context, opts client.SchedulePauseOptions) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Schedule / Pause", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to pause schedule", err)
		}
	}()

	err = sh.handler.Pause(ctx, opts)
	setDefaultAttributes(span, sh.config)
	setScheduleAttributes(span, sh.id)

	return err
}

/*
Unpause unpauses the schedule. It will also overwrite the schedules current note
with the new note.

It automatically handles tracing and error recording.
*/
func (sh *schedulehandle) Unpause(ctx context.Context, opts client.ScheduleUnpauseOptions) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Schedule / Unpause", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to unpause schedule", err)
		}
	}()

	err = sh.handler.Unpause(ctx, opts)
	setDefaultAttributes(span, sh.config)
	setScheduleAttributes(span, sh.id)

	return err
}
