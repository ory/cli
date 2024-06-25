// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"

	cloud "github.com/ory/client-go"
)

func (h *CommandHelper) CreateEventStream(ctx context.Context, projectID string, body cloud.CreateEventStreamBody) (*cloud.EventStream, error) {
	c, err := h.newCloudClient(ctx)
	if err != nil {
		return nil, err
	}

	stream, res, err := c.EventsAPI.CreateEventStream(ctx, projectID).CreateEventStreamBody(body).Execute()
	if err != nil {
		return nil, handleError("unable to create event stream", res, err)
	}

	return stream, nil
}

func (h *CommandHelper) UpdateEventStream(ctx context.Context, projectID, streamID string, body cloud.SetEventStreamBody) (*cloud.EventStream, error) {
	c, err := h.newCloudClient(ctx)
	if err != nil {
		return nil, err
	}

	stream, res, err := c.EventsAPI.SetEventStream(ctx, projectID, streamID).SetEventStreamBody(body).Execute()
	if err != nil {
		return nil, handleError("unable to update event stream", res, err)
	}

	return stream, nil
}

func (h *CommandHelper) DeleteEventStream(ctx context.Context, projectID, streamID string) error {
	c, err := h.newCloudClient(ctx)
	if err != nil {
		return err
	}

	res, err := c.EventsAPI.DeleteEventStream(ctx, projectID, streamID).Execute()
	if err != nil {
		return handleError("unable to delete event stream", res, err)
	}

	return nil
}

func (h *CommandHelper) ListEventStreams(ctx context.Context, projectID string) (*cloud.ListEventStreams, error) {
	c, err := h.newCloudClient(ctx)
	if err != nil {
		return nil, err
	}

	streams, res, err := c.EventsAPI.ListEventStreams(ctx, projectID).Execute()
	if err != nil {
		return nil, handleError("unable to list event streams", res, err)
	}

	return streams, nil
}
