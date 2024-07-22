package auth

import (
    "encoding/hex"
    "testing"
)

func TestMakeRefreshToken(t *testing.T) {
    token, err := MakeRefreshToken()

    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if token == "" {
        t.Errorf("Expected non-empty token, got an empty string")
    }
    _, decodingErr := hex.DecodeString(token)
    if decodingErr != nil {
        t.Errorf("Expected valid hexadecimal token, got %v", decodingErr)
    }
}
