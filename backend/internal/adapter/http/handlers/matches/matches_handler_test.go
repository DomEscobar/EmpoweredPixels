package matches

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMatchOptionsDto_HandlesPartialOptions(t *testing.T) {
	jsonStr := `{"botCount": 3, "autoStart": true}`

	var dto matchOptionsDto
	err := json.Unmarshal([]byte(jsonStr), &dto)

	assert.NoError(t, err)
	assert.Equal(t, 3, *dto.BotCount)
	assert.Equal(t, true, dto.AutoStart)
	assert.Nil(t, dto.BotPowerlevel)
}

func TestMatchOptionsDto_AllOptions(t *testing.T) {
	jsonStr := `{"isPrivate": true, "botCount": 7, "botPowerlevel": 50, "autoStart": false}`

	var dto matchOptionsDto
	err := json.Unmarshal([]byte(jsonStr), &dto)

	assert.NoError(t, err)
	assert.Equal(t, true, dto.IsPrivate)
	assert.Equal(t, 7, *dto.BotCount)
	assert.Equal(t, 50, *dto.BotPowerlevel)
	assert.Equal(t, false, dto.AutoStart)
}

func TestCreateMatch_ValidatesBotCount(t *testing.T) {
	body := map[string]any{
		"isPrivate":     false,
		"botCount":      5,
		"botPowerlevel": 50,
		"autoStart":     true,
	}
	jsonBody, _ := json.Marshal(body)

	var parsedBody map[string]any
	json.Unmarshal(jsonBody, &parsedBody)

	assert.Equal(t, false, parsedBody["isPrivate"])
	assert.Equal(t, float64(5), parsedBody["botCount"])
	assert.Equal(t, float64(50), parsedBody["botPowerlevel"])
	assert.Equal(t, true, parsedBody["autoStart"])
}

func TestJoinMatch_ValidatesInput(t *testing.T) {
	matchID := uuid.NewString()
	fighterID := uuid.NewString()

	body := map[string]string{
		"matchId":   matchID,
		"fighterId": fighterID,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/match/join", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	var parsedBody map[string]string
	json.Unmarshal(jsonBody, &parsedBody)

	assert.Equal(t, matchID, parsedBody["matchId"])
	assert.Equal(t, fighterID, parsedBody["fighterId"])
}

func TestQuickJoin_ValidatesInput(t *testing.T) {
	fighterID := uuid.NewString()

	body := map[string]string{
		"fighterId": fighterID,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/match/quick-join", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	var parsedBody map[string]string
	json.Unmarshal(jsonBody, &parsedBody)

	assert.Equal(t, fighterID, parsedBody["fighterId"])
}
