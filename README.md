# Coze Go API SDK

[![codecov](https://codecov.io/github/coze-dev/coze-go/graph/badge.svg?token=UXitaQ0wp7)](https://codecov.io/github/coze-dev/coze-go)

## Introduction

The Coze API SDK for Go is a powerful tool designed to seamlessly integrate Coze's open APIs into your Go projects.

Key Features:
- Full support for Coze open APIs and authentication APIs
- Both synchronous and streaming API calls
- Optimized streaming APIs with io.Reader interface
- Optimized list APIs with Iterator interface
- Simple and idiomatic Go API design

## Installation

```bash
go get github.com/coze-dev/coze-go
```

## Usage

### Examples

| Example                       | File                                                    |
|-------------------------------|---------------------------------------------------------|
| pat auth                      | [main.go](examples/auth/token/main.go)                  |
| oauth by web code             | [main.go](examples/auth/web_oauth/main.go)              |
| oauth by jwt flow             | [main.go](examples/auth/jwt_oauth/main.go)              |
| oauth by pkce flow            | [main.go](examples/auth/pkce_oauth/main.go)             |
| oauth by device flow          | [main.go](examples/auth/device_oauth/main.go)           |
| handle auth exception         | [main.go](examples/auth/error/main.go)                  |
| bot create, publish and chat  | [main.go](examples/bots/publish/main.go)                |
| get bot and bot list          | [main.go](examples/bots/retrieve/main.go)               |
| non-stream chat               | [main.go](examples/chats/chat/main.go)                  |
| stream chat                   | [main.go](examples/chats/chat_with_image/main.go)       |
| chat with local plugin        | [main.go](examples/chats/submit_tool_output/main.go)    |
| chat with image               | [main.go](examples/chats/chat_with_image/main.go)       |
| non-stream workflow chat      | [main.go](examples/workflows/runs/create/main.go)       |
| stream workflow chat          | [main.go](examples/workflows/runs/stream/main.go)       |
| async workflow run            | [main.go](examples/workflows/runs/async_run/main.go)    |
| conversation                  | [main.go](examples/conversations/crud/main.go)          |
| list conversation             | [main.go](examples/conversations/list/main.go)          |
| workspace                     | [main.go](examples/workspaces/list/main.go)             |
| create update delete message  | [main.go](examples/conversations/messages/crud/main.go) |
| list message                  | [main.go](examples/conversations/messages/list/main.go) |
| create update delete document | [main.go](examples/datasets/documents/crud/main.go)     |
| list documents                | [main.go](examples/datasets/documents/list/main.go)     |
| initial client                | [main.go](examples/client/init/main.go)                 |
| how to handle error           | [main.go](examples/client/error/main.go)                |
| get response log id           | [main.go](examples/client/log/main.go)                  |

### Initialize the Coze Client 

To get started, visit https://www.coze.com/open/oauth/pats (or https://www.coze.cn/open/oauth/pats for the CN environment).

Create a new token by clicking "Add Token". Configure the token name, expiration time, and required permissions. Click OK to generate your personal access token.

Important: Store your personal access token securely to prevent unauthorized access.

```go
func main() {
    // Get an access_token through personal access token or oauth.
    token := os.Getenv("COZE_API_TOKEN")
    authCli := coze.NewTokenAuth(token)
    
    /*
     * The default access is api.coze.com, but if you need to access api.coze.cn
     * please use baseUrl to configure the API endpoint to access
     */
    baseURL := os.Getenv("COZE_API_BASE")
    cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(baseURL))
}
```

### Chat

First, create a bot instance in Coze. The bot ID is the last number in the web link URL.

#### Non-Stream Chat

The SDK provides a convenient wrapper function for non-streaming chat operations. It handles polling and message retrieval automatically:

```go
func main() {
    token := os.Getenv("COZE_API_TOKEN")
    botID := os.Getenv("PUBLISHED_BOT_ID")
    uid := os.Getenv("USER_ID")
    
    authCli := coze.NewTokenAuth(token)
    cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))
    
    ctx := context.Background()
    req := &coze.CreateChatReq{
        BotID:  botID,
        UserID: uid,
        Messages: []coze.Message{
            coze.BuildUserQuestionText("What can you do?", nil),
        },
    }
    
    chat, err := cozeCli.Chats.CreateAndPoll(ctx, req, nil)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    if chat.Status == coze.ChatStatusCompleted {
        fmt.Printf("Token usage: %d\n", chat.Usage.TokenCount)
    }
}
```

#### Stream Chat

Use cozeCli.Chats.Stream() to create a streaming chat session:

```go
func main() {
    // ... initialize client as above ...
    
    ctx := context.Background()
    req := &coze.CreateChatReq{
        BotID:  botID,
        UserID: userID,
        Messages: []coze.Message{
            coze.BuildUserQuestionObjects([]coze.MessageObjectString{
                coze.NewTextMessageObject("Describe this picture"),
                coze.NewImageMessageObjectByID(imageInfo.FileInfo.ID),
            }, nil),
        },
    }
    
    resp, err := cozeCli.Chats.Stream(ctx, req)
    if err != nil {
        fmt.Println("Error starting stream:", err)
        return
    }
    defer resp.Close()
    
    for {
        event, err := resp.Recv()
        if errors.Is(err, io.EOF) {
            fmt.Println("Stream finished")
            break
        }
        if err != nil {
            fmt.Println(err)
            break
        }
        
        if event.Event == coze.ChatEventConversationMessageDelta {
            fmt.Print(event.Message.Content)
        } else if event.Event == coze.ChatEventConversationChatCompleted {
            fmt.Printf("Token usage:%d\n", event.Chat.Usage.TokenCount)
        }
    }
}
```

### Files

```go
func main() {
    // ... initialize client as above ...
    
    ctx := context.Background()
    filePath := os.Getenv("FILE_PATH")
    
    // Upload file
    uploadResp, err := cozeCli.Files.Upload(ctx, coze.NewUploadFilesReqWithPath(filePath))
    if err != nil {
        fmt.Println("Error uploading file:", err)
        return
    }
    fileInfo := uploadResp.FileInfo
    
    // Wait for file processing
    time.Sleep(time.Second)
    
    // Retrieve file
    retrievedResp, err := cozeCli.Files.Retrieve(ctx, &coze.RetrieveFilesReq{
        FileID: fileInfo.ID,
    })
    if err != nil {
        fmt.Println("Error retrieving file:", err)
        return
    }
    fmt.Println(retrievedResp.FileInfo)
}
```

### Pagination

The SDK provides an iterator interface for handling paginated results:

```go
func main() {
    // ... initialize client as above ...
    
    ctx := context.Background()
    datasetID, _ := strconv.ParseInt(os.Getenv("DATASET_ID"), 10, 64)
    
    // Use iterator to automatically retrieve next page
    documents, err := cozeCli.Datasets.Documents.List(ctx, &coze.ListDatasetsDocumentsReq{
        Size: 1,
        DatasetID: datasetID,
    })
    if err != nil {
        fmt.Println("Error fetching documents:", err)
        return
    }
    
    for documents.Next() {
        fmt.Println(documents.Current())
    }
    
    fmt.Println("has_more:", documents.HasMore())
}
```

### Error Handling

The SDK uses Go's standard error handling patterns. All API calls return an error value that should be checked:

```go
resp, err := cozeCli.Chats.Chat(ctx, req)
if err != nil {
    if cozeErr, ok := coze.AsCozeError(err); ok {
    // Handle Coze API error
    fmt.Printf("Coze API error: %s (code: %s)\n", cozeErr.ErrorMessage, cozeErr.ErrorCode)
    return
    }
}
``` 