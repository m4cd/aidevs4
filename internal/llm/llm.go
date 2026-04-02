package llm

import (
	"fmt"

	"github.com/openai/openai-go"
)

func PrintAllMessages(messages []openai.ChatCompletionMessageParamUnion) {
	for i, msg := range messages {
		switch {
		case msg.OfUser != nil:
			content := msg.OfUser.Content.OfString
			fmt.Printf("[%d] USER:\n%s\n\n", i, content)

		case msg.OfAssistant != nil:
			fmt.Printf("[%d] ASSISTANT:\n", i)
			if msg.OfAssistant.Content.OfString.Value != "" {
				fmt.Printf("  Text: %s\n", msg.OfAssistant.Content.OfString)
			}
			for _, tc := range msg.OfAssistant.ToolCalls {
				fmt.Printf("  Tool call: %s(%s) [id: %s]\n", tc.Function.Name, tc.Function.Arguments, tc.ID)
			}
			fmt.Println()

		case msg.OfTool != nil:
			fmt.Printf("[%d] TOOL RESULT [id: %s]:\n  %s\n\n", i, msg.OfTool.ToolCallID, msg.OfTool.Content.OfString)

		case msg.OfSystem != nil:
			fmt.Printf("[%d] SYSTEM:\n%s\n\n", i, msg.OfSystem.Content.OfString)
		}
	}
}

func PrintUserMessages(messages []openai.ChatCompletionMessageParamUnion) {
	for i, msg := range messages {
		if msg.OfUser != nil {
			content := msg.OfUser.Content.OfString
			fmt.Printf("[%d] USER:\n%s\n\n", i, content)

		}
	}
}

func ToToolCallParams(toolCalls []openai.ChatCompletionMessageToolCall) []openai.ChatCompletionMessageToolCallParam {
	params := make([]openai.ChatCompletionMessageToolCallParam, len(toolCalls))
	for i, tc := range toolCalls {
		params[i] = openai.ChatCompletionMessageToolCallParam{
			ID:   tc.ID,
			Type: "function",
			Function: openai.ChatCompletionMessageToolCallFunctionParam{
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			},
		}
	}
	return params
}
