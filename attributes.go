package stun

import (
	"fmt"
	"strconv"
)

// Attributes is list of message attributes.
type Attributes []RawAttribute

// Get returns first attribute from list by the type.
// If attribute is present the RawAttribute is returned and the
// boolean is true. Otherwise the returned RawAttribute will be
// empty and boolean will be false.
func (a Attributes) Get(t AttrType) (RawAttribute, bool) {
	for _, candidate := range a {
		if candidate.Type == t {
			return candidate, true
		}
	}
	return RawAttribute{}, false
}

// AttrType is attribute type.
type AttrType uint16

// Attributes from comprehension-required range (0x0000-0x7FFF).
const (
	AttrMappedAddress     AttrType = 0x0001 // MAPPED-ADDRESS
	AttrUsername          AttrType = 0x0006 // USERNAME
	AttrMessageIntegrity  AttrType = 0x0008 // MESSAGE-INTEGRITY
	AttrErrorCode         AttrType = 0x0009 // ERROR-CODE
	AttrUnknownAttributes AttrType = 0x000A // UNKNOWN-ATTRIBUTES
	AttrRealm             AttrType = 0x0014 // REALM
	AttrNonce             AttrType = 0x0015 // NONCE
	AttrXORMappedAddress  AttrType = 0x0020 // XOR-MAPPED-ADDRESS
)

// Attributes from comprehension-optional range (0x8000-0xFFFF).
const (
	AttrSoftware        AttrType = 0x8022 // SOFTWARE
	AttrAlternateServer AttrType = 0x8023 // ALTERNATE-SERVER
	AttrFingerprint     AttrType = 0x8028 // FINGERPRINT
)

// Attributes from RFC 5245 ICE.
const (
	AttrPriority       AttrType = 0x0024 // PRIORITY
	AttrUseCandidate   AttrType = 0x0025 // USE-CANDIDATE
	AttrICEControlled  AttrType = 0x8029 // ICE-CONTROLLED
	AttrICEControlling AttrType = 0x802A // ICE-CONTROLLING
)

// Attributes from RFC 5766 TURN.
const (
	AttrChannelNumber      AttrType = 0x000C // CHANNEL-NUMBER
	AttrLifetime           AttrType = 0x000D // LIFETIME
	AttrXORPeerAddress     AttrType = 0x0012 // XOR-PEER-ADDRESS
	AttrData               AttrType = 0x0013 // DATA
	AttrXORRelayedAddress  AttrType = 0x0016 // XOR-RELAYED-ADDRESS
	AttrEvenPort           AttrType = 0x0018 // EVEN-PORT
	AttrRequestedTransport AttrType = 0x0019 // REQUESTED-TRANSPORT
	AttrDontFragment       AttrType = 0x001A // DONT-FRAGMENT
	AttrReservationToken   AttrType = 0x0022 // RESERVATION-TOKEN
)

// Attributes from An Origin Attribute for the STUN Protocol.
const (
	AttrOrigin AttrType = 0x802F
)

// Value returns uint16 representation of attribute type.
func (t AttrType) Value() uint16 {
	return uint16(t)
}

var attrNames = map[AttrType]string{
	AttrMappedAddress:      "MAPPED-ADDRESS",
	AttrUsername:           "USERNAME",
	AttrErrorCode:          "ERROR-CODE",
	AttrMessageIntegrity:   "MESSAGE-INTEGRITY",
	AttrUnknownAttributes:  "UNKNOWN-ATTRIBUTES",
	AttrRealm:              "REALM",
	AttrNonce:              "NONCE",
	AttrXORMappedAddress:   "XOR-MAPPED-ADDRESS",
	AttrSoftware:           "SOFTWARE",
	AttrAlternateServer:    "ALTERNATE-SERVER",
	AttrFingerprint:        "FINGERPRINT",
	AttrPriority:           "PRIORITY",
	AttrUseCandidate:       "USE-CANDIDATE",
	AttrICEControlled:      "ICE-CONTROLLED",
	AttrICEControlling:     "ICE-CONTROLLING",
	AttrChannelNumber:      "CHANNEL-NUMBER",
	AttrLifetime:           "LIFETIME",
	AttrXORPeerAddress:     "XOR-PEER-ADDRESS",
	AttrData:               "DATA",
	AttrXORRelayedAddress:  "XOR-RELAYED-ADDRESS",
	AttrEvenPort:           "EVEN-PORT",
	AttrRequestedTransport: "REQUESTED-TRANSPORT",
	AttrDontFragment:       "DONT-FRAGMENT",
	AttrReservationToken:   "RESERVATION-TOKEN",
	AttrOrigin:             "ORIGIN",
}

func (t AttrType) String() string {
	s, ok := attrNames[t]
	if !ok {
		// Just return hex representation of unknown attribute type.
		return "0x" + strconv.FormatUint(uint64(t), 16)
	}
	return s
}

// RawAttribute is a Type-Length-Value (TLV) object that
// can be added to a STUN message. Attributes are divided into two
// types: comprehension-required and comprehension-optional.  STUN
// agents can safely ignore comprehension-optional attributes they
// don't understand, but cannot successfully process a message if it
// contains comprehension-required attributes that are not
// understood.
//
// TODO(ar): Decide to use pointer or non-pointer RawAttribute.
type RawAttribute struct {
	Type   AttrType
	Length uint16 // ignored while encoding
	Value  []byte
}

// AddTo adds RawAttribute to m.
func (a *RawAttribute) AddTo(m *Message) error {
	m.Add(a.Type, m.Raw)
	return nil
}

// Equal returns true if a == b.
func (a RawAttribute) Equal(b RawAttribute) bool {
	if a.Type != b.Type {
		return false
	}
	if a.Length != b.Length {
		return false
	}
	if len(b.Value) != len(a.Value) {
		return false
	}
	for i, v := range a.Value {
		if b.Value[i] != v {
			return false
		}
	}
	return true
}

func (a RawAttribute) String() string {
	return fmt.Sprintf("%s: %x", a.Type, a.Value)
}

// ErrAttributeNotFound means that attribute with provided attribute
// type does not exist in message.
const ErrAttributeNotFound Error = "Attribute not found"

// Get returns byte slice that represents attribute value,
// if there is no attribute with such type,
// ErrAttributeNotFound is returned.
func (m *Message) Get(t AttrType) ([]byte, error) {
	v, ok := m.Attributes.Get(t)
	if !ok {
		return nil, ErrAttributeNotFound
	}
	return v.Value, nil
}
