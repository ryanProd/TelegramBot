package handler_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	handler "github.com/ryanProd/TelegramBot"
)

func TestHandler(t *testing.T) {
	testRecorder := httptest.NewRecorder()
	bodyReader := strings.NewReader(`{"update_id": 12124, "message": {"message_id": 123123, "text": "Hi Welcome", 
	  "from": {"id": 23545, "username": "testbot"}, "chat": {"id": 12314, "type": "dm"}}}`)
	testRequest := httptest.NewRequest("POST", "/apis", bodyReader)
	testRequest.Header.Add("tracing-id", "123")

	handler.Handler(testRecorder, testRequest)

	if testRecorder.Result().StatusCode != 500 {
		t.Errorf("Status code returned, %d, did not match expected code %d", testRecorder.Result().StatusCode, 200)
	}
	if testRecorder.Result().Header.Get("tracing-id") != "123" {
		t.Errorf("Header value for `tracing-id`, %s, did not match expected value %s", testRecorder.Result().Header.Get("tracing-id"), "123")
	}
}
