package apple

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"fmt"

	"github.com/blacktop/go-macho/pkg/codesign/types"
)

const (
	cdVersion        = 0x20500
	cdFlagAdhoc      = 0x2
	cdFlagRuntime    = 0x10000
	execSegMain      = 0x1
	execSegAllowUnsg = 0x10
	slotCMS          = types.SlotType(0x10000)
	pageShift        = 12
)

type codeInfo struct {
	code     []byte
	textOff  uint64
	textSize uint64
	exec     bool
	sdk      uint32
}

func buildSignature(ci codeInfo, s *Signer, o MachOOptions) ([]byte, []byte, error) {
	adhoc := s == nil || s.Key == nil
	var chain []*x509.Certificate
	if !adhoc {
		chain = s.Chain
	}
	reqBlob, err := types.CreateRequirements(o.Identifier, chain, adhoc)
	if err != nil {
		return nil, nil, err
	}
	reqBytes, err := reqBlob.Bytes()
	if err != nil {
		return nil, nil, err
	}

	slots := 2
	var entBlob, derBlob types.Blob
	var entBytes, derBytes []byte
	allowUnsigned := false
	if len(o.Entitlements) > 0 && ci.exec {
		v, err := ParsePlist(o.Entitlements)
		if err != nil {
			return nil, nil, fmt.Errorf("entitlements: %w", err)
		}
		der, err := PlistDER(v)
		if err != nil {
			return nil, nil, err
		}
		if d, ok := v.(map[string]any); ok {
			allowUnsigned, _ = d["get-task-allow"].(bool)
		}
		entBlob = types.NewBlob(types.MAGIC_EMBEDDED_ENTITLEMENTS, o.Entitlements)
		derBlob = types.NewBlob(types.MAGIC_EMBEDDED_ENTITLEMENTS_DER, der)
		entBytes, _ = entBlob.Bytes()
		derBytes, _ = derBlob.Bytes()
		slots = 7
	} else if len(o.ResourcesHash) > 0 {
		slots = 3
	}

	special := make([][]byte, slots)
	for i := range special {
		special[i] = make([]byte, sha256.Size)
	}
	if o.InfoPlist != nil {
		h := sha256.Sum256(o.InfoPlist)
		special[0] = h[:]
	}
	h := sha256.Sum256(reqBytes)
	special[1] = h[:]
	if slots >= 3 && len(o.ResourcesHash) > 0 {
		special[2] = o.ResourcesHash
	}
	if slots >= 7 {
		h5 := sha256.Sum256(entBytes)
		h7 := sha256.Sum256(derBytes)
		special[4] = h5[:]
		special[6] = h7[:]
	}

	pages := (len(ci.code) + (1 << pageShift) - 1) >> pageShift
	var flags uint32
	if adhoc {
		flags |= cdFlagAdhoc
	}
	if o.HardenedRuntime {
		flags |= cdFlagRuntime
	}
	var execFlags uint64
	if ci.exec {
		execFlags |= execSegMain
		if allowUnsigned {
			execFlags |= execSegAllowUnsg
		}
	}
	team := ""
	if !adhoc {
		team = TeamID(chain[0])
	}

	identOff := 96
	teamOff := 0
	hashOff := identOff + len(o.Identifier) + 1
	if team != "" {
		teamOff = hashOff
		hashOff += len(team) + 1
	}
	hashOff += slots * sha256.Size

	var runtime uint32
	if o.HardenedRuntime {
		runtime = ci.sdk
	}

	var cd bytes.Buffer
	w := func(v any) { binary.Write(&cd, binary.BigEndian, v) }
	w(uint32(cdVersion))
	w(flags)
	w(uint32(hashOff))
	w(uint32(identOff))
	w(uint32(slots))
	w(uint32(pages))
	w(uint32(len(ci.code)))
	cd.Write([]byte{sha256.Size, 2, 0, pageShift})
	w(uint32(0))
	w(uint32(0))
	w(uint32(teamOff))
	w(uint32(0))
	w(uint64(0))
	w(ci.textOff)
	w(ci.textSize)
	w(execFlags)
	w(runtime)
	w(uint32(0))
	cd.WriteString(o.Identifier)
	cd.WriteByte(0)
	if team != "" {
		cd.WriteString(team)
		cd.WriteByte(0)
	}
	for i := slots - 1; i >= 0; i-- {
		cd.Write(special[i])
	}
	for i := 0; i < pages; i++ {
		end := min((i+1)<<pageShift, len(ci.code))
		h := sha256.Sum256(ci.code[i<<pageShift : end])
		cd.Write(h[:])
	}
	cdBlob := types.NewBlob(types.MAGIC_CODEDIRECTORY, cd.Bytes())
	cdBytes, err := cdBlob.Bytes()
	if err != nil {
		return nil, nil, err
	}

	var cms []byte
	if !adhoc {
		if cms, err = SignCMS(cdBytes, s.Key, s.Chain, s.Timestamp); err != nil {
			return nil, nil, err
		}
	}

	sb := types.NewSuperBlob(types.MAGIC_EMBEDDED_SIGNATURE)
	sb.AddBlob(types.CSSLOT_CODEDIRECTORY, cdBlob)
	sb.AddBlob(types.CSSLOT_REQUIREMENTS, reqBlob)
	if slots >= 7 {
		sb.AddBlob(types.CSSLOT_ENTITLEMENTS, entBlob)
		sb.AddBlob(types.CSSLOT_ENTITLEMENTS_DER, derBlob)
	}
	sb.AddBlob(slotCMS, types.NewBlob(types.MAGIC_BLOBWRAPPER, cms))
	var out bytes.Buffer
	if err := sb.Write(&out, binary.BigEndian); err != nil {
		return nil, nil, err
	}
	return out.Bytes(), cdBytes, nil
}
