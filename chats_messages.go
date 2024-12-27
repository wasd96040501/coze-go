package coze

import (
	"context"
	"net/http"
)

func (r *chatMessages) List(ctx context.Context, req *ListChatsMessagesReq) (*ListChatsMessagesResp, error) {
	method := http.MethodGet
	uri := "/v3/chat/message/list"
	resp := &listChatsMessagesResp{}
	err := r.core.Request(ctx, method, uri, nil, resp,
		withHTTPQuery("conversation_id", req.ConversationID),
		withHTTPQuery("chat_id", req.ChatID),
	)
	if err != nil {
		return nil, err
	}
	result := &ListChatsMessagesResp{
		baseModel: baseModel{
			httpResponse: resp.HTTPResponse,
		},
		Messages: resp.Messages,
	}
	return result, nil
}

type chatMessages struct {
	core *core
}

func newChatMessages(core *core) *chatMessages {
	return &chatMessages{core: core}
}

// ListChatsMessagesReq represents the request to list messages
type ListChatsMessagesReq struct {
	// The Conversation ID can be viewed in the 'conversation_id' field of the Response when
	// initiating a conversation through the Chat API.
	ConversationID string `json:"conversation_id"`

	// The Chat ID can be viewed in the 'id' field of the Response when initiating a chat through the
	// Chat API. If it is a streaming response, check the 'id' field in the chat event of the Response.
	ChatID string `json:"chat_id"`
}

// ListChatsMessagesResp represents the response to list messages
type listChatsMessagesResp struct {
	baseResponse
	*ListChatsMessagesResp
}

type ListChatsMessagesResp struct {
	baseModel
	Messages []*Message `json:"data"`
}
