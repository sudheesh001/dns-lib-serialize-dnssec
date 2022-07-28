package dns

import (
	"testing"
)

func mockDNSSECProof() *DNSSECProof {
	sig := &Signature{
		length:    21,
		algorithm: 1,
		labels:    1,
		ttl:       1,
		expires:   1,
		begins:    1,
		key_tag:   1,
		signature: []byte("123"),
	}
	key := &Key{
		length:     9,
		flags:      1,
		protocol:   1,
		algorithm:  1,
		public_key: []byte("abc"),
	}
	entry := &Entering{
		length:          37,
		zType:           EnteringType,
		entry_key_index: 1,
		key_sig:         *sig,
		num_keys:        1,
		keys:            []Key{*key},
	}
	exit := &Leaving{
		length:    38,
		zType:     LeavingType,
		next_name: "example.com.",
		rrtype:    RRType(TypeTXT),
		rrsig:     *sig,
	}
	zp := &ZonePair{
		entry: *entry,
		exit:  *exit,
	}

	return &DNSSECProof{
		initial_key_tag: 9,
		num_zones:       1,
		zones:           []ZonePair{*zp},
	}
}

func TestCopyDNSSEC(t *testing.T) {
	proof := mockDNSSECProof()
	proofCopy := proof.copy()
	if proof.String() != proofCopy.String() {
		t.Fatalf("copy() failed %s != %s", proof.String(), proofCopy.String())
	}
}

func TestPackUnpackDNSSEC(t *testing.T) {
	proof := mockDNSSECProof()
	compression := make(map[string]struct{})
	packedBuf := make([]byte, proof.len(0, compression))
	compressionM := make(map[string]uint16)
	off, err := proof.pack(packedBuf, 0, compressionMap{int: compressionM}, false)
	if err != nil {
		t.Fatalf("pack() failed at offset %d\n", off)
	}

	newProof := &DNSSECProof{}
	off, err = newProof.unpack(packedBuf, 0)
	if err != nil {
		t.Fatalf("unpack() failed at offset %d\n", off)
	}
	if proof.String() != newProof.String() {
		t.Fatalf("unpack() failed %s != %s", proof.String(), newProof.String())
	}
}

func TestLengthDNSSEC(t *testing.T) {
	proof := mockDNSSECProof()
	compression := make(map[string]struct{})

	l := proof.len(0, compression)

	if l != 77 {
		t.Fatalf("len() failed %d != 77", l)
	}
}
