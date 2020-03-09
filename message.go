package frontapp

import (
	"errors"
	"fmt"
)

//
// Request
//

// MessageMetadata represents the meta data of a message
type MessageMetadata struct {
	ThreadRef string            `json:"thread_ref,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

// MessageSender represents the sender of a message
type MessageSender struct {
	ContactID string `json:"contact_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Handle    string `json:"handle"`
}

// ReceiveCustomMessageReq is the ReceiveCustomMessage api request
//   cf. https://dev.frontapp.com/reference/messages-1#post_channels-channel-id-incoming-messages
type ReceiveCustomMessageReq struct {
	ChannelID   string          `json:"channel_id"`
	Sender      MessageSender   `json:"sender"`
	Subject     string          `json:"subject,omitempty"`
	Body        string          `json:"body"`
	BodyFormat  string          `json:"body_format,omitempty"`
	Metadata    MessageMetadata `json:"metadata,omitempty"`
	Attachments []string        `json:"attachments,omitempty"`
}

//
// Response
//

// MessageResp is the api response of a message creation request
//   cf. https://dev.frontapp.com/reference/messages-1#post_channels-channel-id-incoming-messages
//       https://dev.frontapp.com/reference/messages-1#import-inbox-message-1
type MessageResp struct {
	Status     string `json:"status"`
	MessageUID string `json:"message_uid"`
}

//
// Endpoints
//

// ReceiveCustomMessage posts a custom message in a custom channel
//   cf. https://dev.frontapp.com/reference/messages-1#post_channels-channel-id-incoming-messages
func (c *Client) ReceiveCustomMessage(req *ReceiveCustomMessageReq) (*MessageResp, error) {
	if req.ChannelID == "" {
		return nil, errors.New("channel id is mandatory")
	}

	path := fmt.Sprintf("channels/%s/incoming_messages", req.ChannelID)

	resp := &MessageResp{}
	if err := c.sendReq("POST", path, req, resp); err != nil {
		return nil, fmt.Errorf("failed to post custom message: %w", err)
	}
	return resp, nil
}
