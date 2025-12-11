package domain

type OpenHabType string

const (
	OpenHabTypeItemCommand OpenHabType = "ItemCommandEvent"
)

type OpenHabPayloadType string

const (
	OpenHabPayloadTypeItemCommandPercent OpenHabPayloadType = "Percent"
	OpenHabPayloadTypeItemCommandHSB     OpenHabPayloadType = "HSB"
)

type OpenHabPayload struct {
	Type  OpenHabPayloadType `json:"type"`
	Value string             `json:"value"`
}

type OpenHabMessage struct {
	Type    OpenHabType    `json:"type"`
	Topic   string         `json:"topic"`
	Payload OpenHabPayload `json:"payload"`
}
