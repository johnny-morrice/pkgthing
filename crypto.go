package pkgthing

type KeyType uint16

const (
	GODLESS_KEY = KeyType(iota)
)

type SignatureBlob []byte

type KeyFingerprint []byte

type KeyReference struct {
	Type        KeyType
	Fingerprint KeyFingerprint
}

type Signature struct {
	Fingerprint KeyReference
	Data        SignatureBlob
}
