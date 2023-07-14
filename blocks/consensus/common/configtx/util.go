package configtx

import (
	"github.com/blcvn/lib-golang-test/blocks/consensus/protoutil"
	cb "github.com/blcvn/lib-golang-test/blocks/types/common"
	"github.com/gogo/protobuf/proto"
)

// UnmarshalConfigUpdate attempts to unmarshal bytes to a *cb.ConfigUpdate
func UnmarshalConfigUpdate(data []byte) (*cb.ConfigUpdate, error) {
	configUpdate := &cb.ConfigUpdate{}
	err := proto.Unmarshal(data, configUpdate)
	if err != nil {
		return nil, err
	}
	return configUpdate, nil
}

// UnmarshalConfigUpdateFromPayload unmarshals configuration update from given payload
func UnmarshalConfigUpdateFromPayload(payload *cb.Payload) (*cb.ConfigUpdate, error) {
	configEnv, err := UnmarshalConfigEnvelope(payload.Data)
	if err != nil {
		return nil, err
	}
	configUpdateEnv, err := protoutil.EnvelopeToConfigUpdate(configEnv.LastUpdate)
	if err != nil {
		return nil, err
	}

	return UnmarshalConfigUpdate(configUpdateEnv.ConfigUpdate)
}

// UnmarshalConfigEnvelope attempts to unmarshal bytes to a *cb.ConfigEnvelope
func UnmarshalConfigEnvelope(data []byte) (*cb.ConfigEnvelope, error) {
	configEnv := &cb.ConfigEnvelope{}
	err := proto.Unmarshal(data, configEnv)
	if err != nil {
		return nil, err
	}
	return configEnv, nil
}
