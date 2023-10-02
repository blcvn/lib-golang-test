/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoutil

import (
	"errors"

	cb "github.com/blcvn/lib-golang-test/consensus/blocks/types/common"
	"github.com/golang/protobuf/proto"
)

// MarshalOrPanic serializes a protobuf message and panics if this
// operation fails
func MarshalOrPanic(pb proto.Message) []byte {
	data, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	return data
}

// IsConfigBlock validates whenever given block contains configuration
// update transaction
func IsConfigBlock(block *cb.Block) bool {
	envelope, err := ExtractEnvelope(block, 0)
	if err != nil {
		return false
	}

	payload, err := UnmarshalPayload(envelope.Payload)
	if err != nil {
		return false
	}

	if payload.Header == nil {
		return false
	}

	hdr, err := UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return false
	}

	return cb.HeaderType(hdr.Type) == cb.HeaderType_CONFIG
}

// ExtractEnvelope retrieves the requested envelope from a given block and
// unmarshals it
func ExtractEnvelope(block *cb.Block, index int) (*cb.Envelope, error) {
	if block.Data == nil {
		return nil, errors.New("block data is nil")
	}

	envelopeCount := len(block.Data.Data)
	if index < 0 || index >= envelopeCount {
		return nil, errors.New("envelope index out of bounds")
	}
	marshaledEnvelope := block.Data.Data[index]
	envelope, err := GetEnvelopeFromBlock(marshaledEnvelope)
	err = WithMessagef(err, "block data does not carry an envelope at index %d", index)
	return envelope, err
}

// GetEnvelopeFromBlock gets an envelope from a block's Data field.
func GetEnvelopeFromBlock(data []byte) (*cb.Envelope, error) {
	// Block always begins with an envelope
	var err error
	env := &cb.Envelope{}
	if err = proto.Unmarshal(data, env); err != nil {
		return nil, Wrap(err, "error unmarshalling Envelope")
	}

	return env, nil
}

// ChannelHeader returns the *cb.ChannelHeader for a given *cb.Envelope.
func ChannelHeader(env *cb.Envelope) (*cb.ChannelHeader, error) {
	if env == nil {
		return nil, errors.New("Invalid envelope payload. can't be nil")
	}

	envPayload, err := UnmarshalPayload(env.Payload)
	if err != nil {
		return nil, err
	}

	if envPayload.Header == nil {
		return nil, errors.New("header not set")
	}

	if envPayload.Header.ChannelHeader == nil {
		return nil, errors.New("channel header not set")
	}

	chdr, err := UnmarshalChannelHeader(envPayload.Header.ChannelHeader)
	if err != nil {
		return nil, WithMessage(err, "error unmarshalling channel header")
	}

	return chdr, nil
}

// EnvelopeToConfigUpdate is used to extract a ConfigUpdateEnvelope from an envelope of
// type CONFIG_UPDATE
func EnvelopeToConfigUpdate(configtx *cb.Envelope) (*cb.ConfigUpdateEnvelope, error) {
	configUpdateEnv := &cb.ConfigUpdateEnvelope{}
	_, err := UnmarshalEnvelopeOfType(configtx, cb.HeaderType_CONFIG_UPDATE, configUpdateEnv)
	if err != nil {
		return nil, err
	}
	return configUpdateEnv, nil
}

// UnmarshalEnvelopeOfType unmarshals an envelope of the specified type,
// including unmarshalling the payload data
func UnmarshalEnvelopeOfType(envelope *cb.Envelope, headerType cb.HeaderType, message proto.Message) (*cb.ChannelHeader, error) {
	payload, err := UnmarshalPayload(envelope.Payload)
	if err != nil {
		return nil, err
	}

	if payload.Header == nil {
		return nil, errors.New("envelope must have a Header")
	}

	chdr, err := UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, err
	}

	if chdr.Type != int32(headerType) {
		return nil, Errorf("invalid type %s, expected %s", cb.HeaderType(chdr.Type), headerType)
	}

	err = proto.Unmarshal(payload.Data, message)
	err = Wrapf(err, "error unmarshalling message for type %s", headerType)
	return chdr, err
}
