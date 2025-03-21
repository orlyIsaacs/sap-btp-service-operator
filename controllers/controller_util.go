package controllers

import (
	"context"
	"encoding/json"
	"strings"

	"k8s.io/apimachinery/pkg/util/rand"

	v1 "k8s.io/api/authentication/v1"
)

type SecretMetadataProperty struct {
	Name      string `json:"name"`
	Container bool   `json:"container,omitempty"`
	Format    string `json:"format"`
}

type format string

const (
	TEXT    format = "text"
	JSON    format = "json"
	UNKNOWN format = "unknown"
)

func normalizeCredentials(credentialsJSON json.RawMessage) (map[string][]byte, []SecretMetadataProperty, error) {
	var credentialsMap map[string]interface{}
	err := json.Unmarshal(credentialsJSON, &credentialsMap)
	if err != nil {
		return nil, nil, err
	}

	normalized := make(map[string][]byte)
	metadata := make([]SecretMetadataProperty, 0)
	for propertyName, value := range credentialsMap {
		keyString := strings.Replace(propertyName, " ", "_", -1)
		normalizedValue, typpe, err := serialize(value)
		if err != nil {
			return nil, nil, err
		}
		metadata = append(metadata, SecretMetadataProperty{
			Name:   keyString,
			Format: string(typpe),
		})
		normalized[keyString] = normalizedValue
	}
	return normalized, metadata, nil
}

func buildUserInfo(ctx context.Context, userInfo *v1.UserInfo) string {
	log := GetLogger(ctx)
	if userInfo == nil {
		return ""
	}
	userInfoStr, err := json.Marshal(userInfo)
	if err != nil {
		log.Error(err, "failed to prepare user info")
		return ""
	}

	return string(userInfoStr)
}

func serialize(value interface{}) ([]byte, format, error) {
	if byteArrayVal, ok := value.([]byte); ok {
		return byteArrayVal, JSON, nil
	}
	if strVal, ok := value.(string); ok {
		return []byte(strVal), TEXT, nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil, UNKNOWN, err
	}
	return data, JSON, nil
}

func contains(slice []string, i string) bool {
	for _, s := range slice {
		if s == i {
			return true
		}
	}

	return false
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
