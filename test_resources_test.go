/*
Gapstone is a Go binding for the Capstone disassembly library. For examples,
try reading the *_test.go files.

	Library Author: Nguyen Anh Quynh
	Binding Author: Ben Nagy
	License: BSD style - see LICENSE file for details
    (c) 2013 COSEINC. All Rights Reserved.
*/

package gapstone

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"testing"
)

// Maintain the expected version and sanity checks manually, so we can verify
// against the installed C lib. Not foolproof, but should save 90% of accidents
const (
	expectedMaj = 5
	expectedMin = 0
)

type sanityCheck struct {
	insMax int
	regMax int
	grpMax int
}

type sanityChecks map[int]sanityCheck

func (s *sanityChecks) Maj() int { return expectedMaj }
func (s *sanityChecks) Min() int { return expectedMin }

// Remember the all the constants CONST are direct refs to C.CONST, so in
// combination with these we should be _fairly_ sure we're getting the
// disassembly capstone expects to provide.
var checks = sanityChecks{
	CS_ARCH_ARM64: sanityCheck{
		regMax: 342,
		insMax: 1289,
		grpMax: 156,
	},
	CS_ARCH_ARM: sanityCheck{
		regMax: 111,
		insMax: 472,
		grpMax: 160,
	},
	CS_ARCH_MIPS: sanityCheck{
		regMax: 137,
		insMax: 626,
		grpMax: 161,
	},
	CS_ARCH_PPC: sanityCheck{
		regMax: 344,
		insMax: 1727,
		grpMax: 143,
	},
	CS_ARCH_SPARC: sanityCheck{
		regMax: 88,
		insMax: 279,
		grpMax: 135,
	},
	CS_ARCH_SYSZ: sanityCheck{
		regMax: 117,
		insMax: 2346,
		grpMax: 157,
	},
	CS_ARCH_X86: sanityCheck{
		regMax: 246,
		insMax: 1524,
		grpMax: 170,
	},
	CS_ARCH_XCORE: sanityCheck{
		regMax: 26,
		insMax: 121,
		grpMax: 2,
	},
}

type option struct {
	ty    uint
	value uint
}

type platform struct {
	arch    int
	mode    int
	options []option
	code    string
	comment string
}

type platforms []platform

var skipPlatforms = []string{
	"MOS65XX",
	"eBPF",
}

func shouldSkipPlatform(platform string) bool {
	return slices.Contains(skipPlatforms, platform)
}

var address = uint64(0x1000)

var (
	armCode     = "\x86\x48\x60\xf4\x4d\x0f\xe2\xf4\xED\xFF\xFF\xEB\x04\xe0\x2d\xe5\x00\x00\x00\x00\xe0\x83\x22\xe5\xf1\x02\x03\x0e\x00\x00\xa0\xe3\x02\x30\xc1\xe7\x00\x00\x53\xe3\x00\x02\x01\xf1\x05\x40\xd0\xe8\xf4\x80\x00\x00"
	armCode2    = "\xd1\xe8\x00\xf0\xf0\x24\x04\x07\x1f\x3c\xf2\xc0\x00\x00\x4f\xf0\x00\x01\x46\x6c"
	thumbCode   = "\x60\xf9\x1f\x04\xe0\xf9\x4f\x07\x70\x47\x00\xf0\x10\xe8\xeb\x46\x83\xb0\xc9\x68\x1f\xb1\x30\xbf\xaf\xf3\x20\x84\x52\xf8\x23\xf0"
	thumbCode2  = "\x4f\xf0\x00\x01\xbd\xe8\x00\x88\xd1\xe8\x00\xf0\x18\xbf\xad\xbf\xf3\xff\x0b\x0c\x86\xf3\x00\x89\x80\xf3\x00\x8c\x4f\xfa\x99\xf6\xd0\xff\xa2\x01"
	thumbMClass = "\xef\xf3\x02\x80"
	armV8       = "\xe0\x3b\xb2\xee\x42\x00\x01\xe1\x51\xf0\x7f\xf5"
	arm64Code   = "\x09\x00\x38\xd5\xbf\x40\x00\xd5\x0c\x05\x13\xd5\x20\x50\x02\x0e\x20\xe4\x3d\x0f\x00\x18\xa0\x5f\xa2\x00\xae\x9e\x9f\x37\x03\xd5\xbf\x33\x03\xd5\xdf\x3f\x03\xd5\x21\x7c\x02\x9b\x21\x7c\x00\x53\x00\x40\x21\x4b\xe1\x0b\x40\xb9\x20\x04\x81\xda\x20\x08\x02\x8b\x10\x5b\xe8\x3c"
	x86Code64   = "\x55\x48\x8b\x05\xb8\x13\x00\x00\xe9\xea\xbe\xad\xde\xff\x25\x23\x01\x00\x00\xe8\xdf\xbe\xad\xde\x74\xff"
	x86Code16   = "\x8d\x4c\x32\x08\x01\xd8\x81\xc6\x34\x12\x00\x00\x05\x23\x01\x00\x00\x36\x8b\x84\x91\x23\x01\x00\x00\x41\x8d\x84\x39\x89\x67\x00\x00\x8d\x87\x89\x67\x00\x00\xb4\xc6\x66\xe9\xb8\x00\x00\x00\x67\xff\xa0\x23\x01\x00\x00\x66\xe8\xcb\x00\x00\x00\x74\xfc"
	x86Code32   = "\x8d\x4c\x32\x08\x01\xd8\x81\xc6\x34\x12\x00\x00\x05\x23\x01\x00\x00\x36\x8b\x84\x91\x23\x01\x00\x00\x41\x8d\x84\x39\x89\x67\x00\x00\x8d\x87\x89\x67\x00\x00\xb4\xc6\xe9\xea\xbe\xad\xde\xff\xa0\x23\x01\x00\x00\xe8\xdf\xbe\xad\xde\x74\xff"
	mipsCode    = "\x0C\x10\x00\x97\x00\x00\x00\x00\x24\x02\x00\x0c\x8f\xa2\x00\x00\x34\x21\x34\x56"
	mipsCode2   = "\x56\x34\x21\x34\xc2\x17\x01\x00"
	mips32R6M   = "\x00\x07\x00\x07\x00\x11\x93\x7c\x01\x8c\x8b\x7c\x00\xc7\x48\xd0"
	mips32R6    = "\xec\x80\x00\x19\x7c\x43\x22\xa0"
	mips64SD    = "\x70\x00\xb2\xff"
	ppcCode     = "\x43\x20\x0c\x07\x41\x56\xff\x17\x80\x20\x00\x00\x80\x3f\x00\x00\x10\x43\x23\x0e\xd0\x44\x00\x80\x4c\x43\x22\x02\x2d\x03\x00\x80\x7c\x43\x20\x14\x7c\x43\x20\x93\x4f\x20\x00\x21\x4c\xc8\x00\x21\x40\x82\x00\x14"
	ppcCode2    = "\x10\x60\x2a\x10\x10\x64\x28\x88\x7c\x4a\x5d\x0f"
	sysZCode    = "\xed\x00\x00\x00\x00\x1a\x5a\x0f\x1f\xff\xc2\x09\x80\x00\x00\x00\x07\xf7\xeb\x2a\xff\xff\x7f\x57\xe3\x01\xff\xff\x7f\x57\xeb\x00\xf0\x00\x00\x24\xb2\x4f\x00\x78\xec\x18\x00\x00\xc1\x7f"
	sparcCode   = "\x80\xa0\x40\x02\x85\xc2\x60\x08\x85\xe8\x20\x01\x81\xe8\x00\x00\x90\x10\x20\x01\xd5\xf6\x10\x16\x21\x00\x00\x0a\x86\x00\x40\x02\x01\x00\x00\x00\x12\xbf\xff\xff\x10\xbf\xff\xff\xa0\x02\x00\x09\x0d\xbf\xff\xff\xd4\x20\x60\x00\xd4\x4e\x00\x16\x2a\xc2\x80\x03"
	sparcV9Code = "\x81\xa8\x0a\x24\x89\xa0\x10\x20\x89\xa0\x1a\x60\x89\xa0\x00\xe0"
	xcoreCode   = "\xfe\x0f\xfe\x17\x13\x17\xc6\xfe\xec\x17\x97\xf8\xec\x4f\x1f\xfd\xec\x37\x07\xf2\x45\x5b\xf9\xfa\x02\x06\x1b\x10\x09\xfd\xec\xa7"
)

var (
	basicX86Code16      = "\x8d\x4c\x32\x08\x01\xd8\x81\xc6\x34\x12\x00\x00"
	basicX86Code32      = "\xba\xcd\xab\x00\x00\x8d\x4c\x32\x08\x01\xd8\x81\xc6\x34\x12\x00\x00"
	basicX86Code64      = "\x55\x48\x8b\x05\xb8\x13\x00\x00"
	basicArmCode        = "\xED\xFF\xFF\xEB\x04\xe0\x2d\xe5\x00\x00\x00\x00\xe0\x83\x22\xe5\xf1\x02\x03\x0e\x00\x00\xa0\xe3\x02\x30\xc1\xe7\x00\x00\x53\xe3"
	basicArmCode2       = "\x10\xf1\x10\xe7\x11\xf2\x31\xe7\xdc\xa1\x2e\xf3\xe8\x4e\x62\xf3"
	basicArmV8          = "\xe0\x3b\xb2\xee\x42\x00\x01\xe1\x51\xf0\x7f\xf5"
	basicThumbMClass    = "\xef\xf3\x02\x80"
	basicThumbCode      = "\x70\x47\xeb\x46\x83\xb0\xc9\x68"
	basicThumbCode2     = "\x4f\xf0\x00\x01\xbd\xe8\x00\x88\xd1\xe8\x00\xf0"
	basicMipsCode       = "\x0C\x10\x00\x97\x00\x00\x00\x00\x24\x02\x00\x0c\x8f\xa2\x00\x00\x34\x21\x34\x56"
	basicMipsCode2      = "\x56\x34\x21\x34\xc2\x17\x01\x00"
	basicMips32R6M      = "\x00\x07\x00\x07\x00\x11\x93\x7c\x01\x8c\x8b\x7c\x00\xc7\x48\xd0"
	basicMips32R6       = "\xec\x80\x00\x19\x7c\x43\x22\xa0"
	basicArm64Code      = "\x21\x7c\x02\x9b\x21\x7c\x00\x53\x00\x40\x21\x4b\xe1\x0b\x40\xb9"
	basicPPCCode        = "\x80\x20\x00\x00\x80\x3f\x00\x00\x10\x43\x23\x0e\xd0\x44\x00\x80\x4c\x43\x22\x02\x2d\x03\x00\x80\x7c\x43\x20\x14\x7c\x43\x20\x93\x4f\x20\x00\x21\x4c\xc8\x00\x21"
	basicPPCCode2       = "\x10\x60\x2a\x10\x10\x64\x28\x88\x7c\x4a\x5d\x0f"
	basicSparcCode      = "\x80\xa0\x40\x02\x85\xc2\x60\x08\x85\xe8\x20\x01\x81\xe8\x00\x00\x90\x10\x20\x01\xd5\xf6\x10\x16\x21\x00\x00\x0a\x86\x00\x40\x02\x01\x00\x00\x00\x12\xbf\xff\xff\x10\xbf\xff\xff\xa0\x02\x00\x09\x0d\xbf\xff\xff\xd4\x20\x60\x00\xd4\x4e\x00\x16\x2a\xc2\x80\x03"
	basicSparcV9Code    = "\x81\xa8\x0a\x24\x89\xa0\x10\x20\x89\xa0\x1a\x60\x89\xa0\x00\xe0"
	basicSysZCode       = "\xed\x00\x00\x00\x00\x1a\x5a\x0f\x1f\xff\xc2\x09\x80\x00\x00\x00\x07\xf7\xeb\x2a\xff\xff\x7f\x57\xe3\x01\xff\xff\x7f\x57\xeb\x00\xf0\x00\x00\x24\xb2\x4f\x00\x78"
	basicXcoreCode      = "\xfe\x0f\xfe\x17\x13\x17\xc6\xfe\xec\x17\x97\xf8\xec\x4f\x1f\xfd\xec\x37\x07\xf2\x45\x5b\xf9\xfa\x02\x06\x1b\x10"
	basicM68KCode       = "\xd4\x40\x87\x5a\x4e\x71\x02\xb4\xc0\xde\xc0\xde\x5c\x00\x1d\x80\x71\x12\x01\x23\xf2\x3c\x44\x22\x40\x49\x0e\x56\x54\xc5\xf2\x3c\x44\x00\x44\x7a\x00\x00\xf2\x00\x0a\x28"
	basicTMS320C64XCode = "\x01\xac\x88\x40\x81\xac\x88\x43\x00\x00\x00\x00\x02\x90\x32\x96\x02\x80\x46\x9e\x05\x3c\x83\xe6\x0b\x0c\x8b\x24"
	basicM680XCode      = "\x06\x10\x19\x1a\x55\x1e\x01\x23\xe9\x31\x06\x34\x55\xa6\x81\xa7\x89\x7f\xff\xa6\x9d\x10\x00\xa7\x91\xa6\x9f\x10\x00\x11\xac\x99\x10\x00\x39"
	basicEVMCode        = "\x60\x61"
)

var (
	detailX86Code16   = "\x8d\x4c\x32\x08\x01\xd8\x81\xc6\x34\x12\x00\x00"
	detailX86Code32   = "\x8d\x4c\x32\x08\x01\xd8\x81\xc6\x34\x12\x00\x00"
	detailX86Code64   = "\x55\x48\x8b\x05\xb8\x13\x00\x00"
	detailArmCode     = "\xED\xFF\xFF\xEB\x04\xe0\x2d\xe5\x00\x00\x00\x00\xe0\x83\x22\xe5\xf1\x02\x03\x0e\x00\x00\xa0\xe3\x02\x30\xc1\xe7\x00\x00\x53\xe3"
	detailArmCode2    = "\x10\xf1\x10\xe7\x11\xf2\x31\xe7\xdc\xa1\x2e\xf3\xe8\x4e\x62\xf3"
	detailThumbCode   = "\x70\x47\xeb\x46\x83\xb0\xc9\x68"
	detailThumbCode2  = "\x4f\xf0\x00\x01\xbd\xe8\x00\x88\xd1\xe8\x00\xf0"
	detailThumbMClass = "\xef\xf3\x02\x80"
	detailArmV8       = "\xe0\x3b\xb2\xee\x42\x00\x01\xe1\x51\xf0\x7f\xf5"
	detailMipsCode    = "\x0C\x10\x00\x97\x00\x00\x00\x00\x24\x02\x00\x0c\x8f\xa2\x00\x00\x34\x21\x34\x56\x00\x80\x04\x08"
	detailMipsCode2   = "\x56\x34\x21\x34\xc2\x17\x01\x00"
	detailMips32R6M   = "\x00\x07\x00\x07\x00\x11\x93\x7c\x01\x8c\x8b\x7c\x00\xc7\x48\xd0"
	detailMips32R6    = "\xec\x80\x00\x19\x7c\x43\x22\xa0"
	detailArm64Code   = "\x09\x00\x38\xd5\xbf\x40\x00\xd5\x0c\x05\x13\xd5\x20\x50\x02\x0e\x20\xe4\x3d\x0f\x00\x18\xa0\x5f\xa2\x00\xae\x9e\x9f\x37\x03\xd5\xbf\x33\x03\xd5\xdf\x3f\x03\xd5\x21\x7c\x02\x9b\x21\x7c\x00\x53\x00\x40\x21\x4b\xe1\x0b\x40\xb9\x20\x04\x81\xda\x20\x08\x02\x8b\x10\x5b\xe8\x3c"
	detailPPCCode     = "\x80\x20\x00\x00\x80\x3f\x00\x00\x10\x43\x23\x0e\xd0\x44\x00\x80\x4c\x43\x22\x02\x2d\x03\x00\x80\x7c\x43\x20\x14\x7c\x43\x20\x93\x4f\x20\x00\x21\x4c\xc8\x00\x21\x40\x82\x00\x14"
	detailPPCCode2    = "\x10\x60\x2a\x10\x10\x64\x28\x88\x7c\x4a\x5d\x0f"
	detailSparcCode   = "\x80\xa0\x40\x02\x85\xc2\x60\x08\x85\xe8\x20\x01\x81\xe8\x00\x00\x90\x10\x20\x01\xd5\xf6\x10\x16\x21\x00\x00\x0a\x86\x00\x40\x02\x01\x00\x00\x00\x12\xbf\xff\xff\x10\xbf\xff\xff\xa0\x02\x00\x09\x0d\xbf\xff\xff\xd4\x20\x60\x00\xd4\x4e\x00\x16\x2a\xc2\x80\x03"
	detailSparcV9Code = "\x81\xa8\x0a\x24\x89\xa0\x10\x20\x89\xa0\x1a\x60\x89\xa0\x00\xe0"
	detailSysZCode    = "\xed\x00\x00\x00\x00\x1a\x5a\x0f\x1f\xff\xc2\x09\x80\x00\x00\x00\x07\xf7\xeb\x2a\xff\xff\x7f\x57\xe3\x01\xff\xff\x7f\x57\xeb\x00\xf0\x00\x00\x24\xb2\x4f\x00\x78"
	detailXcoreCode   = "\xfe\x0f\xfe\x17\x13\x17\xc6\xfe\xec\x17\x97\xf8\xec\x4f\x1f\xfd\xec\x37\x07\xf2\x45\x5b\xf9\xfa\x02\x06\x1b\x10"
	detailM68KCode    = "\xd4\x40\x87\x5a\x4e\x71\x02\xb4\xc0\xde\xc0\xde\x5c\x00\x1d\x80\x71\x12\x01\x23\xf2\x3c\x44\x22\x40\x49\x0e\x56\x54\xc5\xf2\x3c\x44\x00\x44\x7a\x00\x00\xf2\x00\x0a\x28"
	detailM680XCode   = "\x06\x10\x19\x1a\x55\x1e\x01\x23\xe9\x31\x06\x34\x55\xa6\x81\xa7\x89\x7f\xff\xa6\x9d\x10\x00\xa7\x91\xa6\x9f\x10\x00\x11\xac\x99\x10\x00\x39"
)

var basicTests = platforms{
	{
		CS_ARCH_X86,
		CS_MODE_16,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicX86Code16,
		"X86 16bit (Intel syntax)",
	},
	{
		CS_ARCH_X86,
		CS_MODE_32,
		[]option{{CS_OPT_SYNTAX, CS_OPT_SYNTAX_ATT}, {CS_OPT_DETAIL, CS_OPT_ON}},
		basicX86Code32,
		"X86 32bit (ATT syntax)",
	},
	{
		CS_ARCH_X86,
		CS_MODE_32,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicX86Code32,
		"X86 32 (Intel syntax)",
	},
	{
		CS_ARCH_X86,
		CS_MODE_32,
		[]option{{CS_OPT_SYNTAX, CS_OPT_SYNTAX_MASM}, {CS_OPT_DETAIL, CS_OPT_ON}},
		basicX86Code32,
		"X86 32 (MASM syntax)",
	},
	{
		CS_ARCH_X86,
		CS_MODE_64,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicX86Code64,
		"X86 64 (Intel syntax)",
	},
	{
		CS_ARCH_ARM,
		CS_MODE_ARM,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicArmCode,
		"ARM",
	},
	{
		CS_ARCH_ARM,
		CS_MODE_THUMB,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicThumbCode2,
		"THUMB-2",
	},
	{
		CS_ARCH_ARM,
		CS_MODE_ARM,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicArmCode2,
		"ARM: Cortex-A15 + NEON",
	},
	{
		CS_ARCH_ARM,
		CS_MODE_THUMB,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicThumbCode,
		"THUMB",
	},
	platform{
		CS_ARCH_ARM,
		CS_MODE_THUMB + CS_MODE_MCLASS,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicThumbMClass,
		"Thumb-MClass",
	},
	platform{
		CS_ARCH_ARM,
		CS_MODE_ARM + CS_MODE_V8,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicArmV8,
		"Arm-V8",
	},
	{
		CS_ARCH_MIPS,
		CS_MODE_MIPS32 + CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicMipsCode,
		"MIPS-32 (Big-endian)",
	},
	{
		CS_ARCH_MIPS,
		CS_MODE_MIPS64 + CS_MODE_LITTLE_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicMipsCode2,
		"MIPS-64-EL (Little-endian)",
	},
	{
		CS_ARCH_MIPS,
		CS_MODE_MIPS32R6 + CS_MODE_MICRO + CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicMips32R6M,
		"MIPS-32R6 | Micro (Big-endian)",
	},
	{
		CS_ARCH_MIPS,
		CS_MODE_MIPS32R6 + CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicMips32R6,
		"MIPS-32R6 (Big-endian)",
	},
	platform{
		CS_ARCH_ARM64,
		CS_MODE_ARM,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicArm64Code,
		"ARM-64",
	},
	platform{
		CS_ARCH_PPC,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicPPCCode,
		"PPC-64",
	},
	platform{
		CS_ARCH_PPC,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}, {CS_OPT_SYNTAX, CS_OPT_SYNTAX_NOREGNAME}},
		basicPPCCode,
		"PPC-64, print register with number only",
	},
	platform{
		CS_ARCH_PPC,
		CS_MODE_BIG_ENDIAN + CS_MODE_QPX,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicPPCCode2,
		"PPC-64 + QPX",
	},
	platform{
		CS_ARCH_SPARC,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicSparcCode,
		"Sparc",
	},
	platform{
		CS_ARCH_SPARC,
		CS_MODE_BIG_ENDIAN + CS_MODE_V9,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicSparcV9Code,
		"SparcV9",
	},
	platform{
		CS_ARCH_SYSZ,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicSysZCode,
		"SystemZ",
	},
	platform{
		CS_ARCH_XCORE,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicXcoreCode,
		"XCore",
	},
	platform{
		CS_ARCH_M68K,
		CS_MODE_BIG_ENDIAN + CS_MODE_M68K_040,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicM68KCode,
		"M68K",
	},
	platform{
		CS_ARCH_TMS320C64X,
		0,
		[]option{},
		basicTMS320C64XCode,
		"TMS320C64x",
	},
	platform{
		CS_ARCH_M680X,
		CS_MODE_M680X_6809,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicM680XCode,
		"M680X_M6809",
	},
	platform{
		CS_ARCH_EVM,
		0,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		basicEVMCode,
		"EVM",
	},
}

// Honestly, these are _almost_ identical, but it's just easier to maintain
// them as a separate list and not mess about modifying the slice in the test
// code.
var detailTests = platforms{
	{
		CS_ARCH_X86,
		CS_MODE_16,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailX86Code16,
		"X86 16bit (Intel syntax)",
	},
	{
		CS_ARCH_X86,
		CS_MODE_32,
		[]option{{CS_OPT_SYNTAX, CS_OPT_SYNTAX_ATT}, {CS_OPT_DETAIL, CS_OPT_ON}},
		detailX86Code32,
		"X86 32bit (ATT syntax)",
	},
	{
		CS_ARCH_X86,
		CS_MODE_32,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailX86Code32,
		"X86 32 (Intel syntax)",
	},
	{
		CS_ARCH_X86,
		CS_MODE_64,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailX86Code64,
		"X86 64 (Intel syntax)",
	},
	{
		CS_ARCH_ARM,
		CS_MODE_ARM,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailArmCode,
		"ARM",
	},
	{
		CS_ARCH_ARM,
		CS_MODE_THUMB,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailThumbCode2,
		"THUMB-2",
	},
	{
		CS_ARCH_ARM,
		CS_MODE_ARM,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailArmCode2,
		"ARM: Cortex-A15 + NEON",
	},
	{
		CS_ARCH_ARM,
		CS_MODE_THUMB,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailThumbCode,
		"THUMB",
	},
	platform{
		CS_ARCH_ARM,
		CS_MODE_THUMB + CS_MODE_MCLASS,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailThumbMClass,
		"Thumb-MClass",
	},
	platform{
		CS_ARCH_ARM,
		CS_MODE_ARM + CS_MODE_V8,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailArmV8,
		"Arm-V8",
	},
	{
		CS_ARCH_MIPS,
		CS_MODE_32 + CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailMipsCode,
		"MIPS-32 (Big-endian)",
	},
	{
		CS_ARCH_MIPS,
		CS_MODE_64 + CS_MODE_LITTLE_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailMipsCode2,
		"MIPS-64-EL (Little-endian)",
	},
	{
		CS_ARCH_MIPS,
		CS_MODE_MIPS32R6 + CS_MODE_MICRO + CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailMips32R6M,
		"MIPS-32R6 | Micro (Big-endian)",
	},
	{
		CS_ARCH_MIPS,
		CS_MODE_MIPS32R6 + CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailMips32R6,
		"MIPS-32R6 (Big-endian)",
	},
	platform{
		CS_ARCH_ARM64,
		CS_MODE_ARM,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailArm64Code,
		"ARM-64",
	},
	platform{
		CS_ARCH_PPC,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailPPCCode,
		"PPC-64",
	},
	platform{
		CS_ARCH_PPC,
		CS_MODE_BIG_ENDIAN + CS_MODE_QPX,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailPPCCode2,
		"PPC-64 + QPX",
	},
	platform{
		CS_ARCH_SPARC,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailSparcCode,
		"Sparc",
	},
	platform{
		CS_ARCH_SPARC,
		CS_MODE_BIG_ENDIAN + CS_MODE_V9,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailSparcV9Code,
		"SparcV9",
	},
	platform{
		CS_ARCH_SYSZ,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailSysZCode,
		"SystemZ",
	},
	platform{
		CS_ARCH_XCORE,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailXcoreCode,
		"XCore",
	},
	platform{
		CS_ARCH_M68K,
		CS_MODE_BIG_ENDIAN + CS_MODE_M68K_040,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailM68KCode,
		"M68K",
	},
	platform{
		CS_ARCH_M680X,
		CS_MODE_M680X_6809,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		detailM680XCode,
		"M680X_M6809",
	},
}

var armTests = platforms{
	platform{
		CS_ARCH_ARM,
		CS_MODE_ARM,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		armCode,
		"ARM",
	},
	platform{
		CS_ARCH_ARM,
		CS_MODE_THUMB,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		thumbCode,
		"Thumb",
	},
	platform{
		CS_ARCH_ARM,
		CS_MODE_THUMB,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		armCode2,
		"Thumb-mixed",
	},
	platform{
		CS_ARCH_ARM,
		CS_MODE_THUMB,
		[]option{
			{CS_OPT_DETAIL, CS_OPT_ON},
			{CS_OPT_SYNTAX, CS_OPT_SYNTAX_NOREGNAME},
		},
		thumbCode2,
		"Thumb-2 & register named with numbers",
	},
	platform{
		CS_ARCH_ARM,
		CS_MODE_THUMB + CS_MODE_MCLASS,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		thumbMClass,
		"Thumb-MClass",
	},
	platform{
		CS_ARCH_ARM,
		CS_MODE_ARM + CS_MODE_V8,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		armV8,
		"Arm-V8",
	},
}

var arm64Tests = platforms{
	platform{
		CS_ARCH_ARM64,
		CS_MODE_ARM,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		arm64Code,
		"ARM-64",
	},
}

var mips_tests = platforms{
	platform{
		CS_ARCH_MIPS,
		CS_MODE_32 + CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		mipsCode,
		"MIPS-32 (Big-endian)",
	},
	platform{
		CS_ARCH_MIPS,
		CS_MODE_64 + CS_MODE_LITTLE_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		mipsCode2,
		"MIPS-64-EL (Little-endian)",
	},
	{
		CS_ARCH_MIPS,
		CS_MODE_MIPS32R6 + CS_MODE_MICRO + CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		mips32R6M,
		"MIPS-32R6 | Micro (Big-endian)",
	},
	{
		CS_ARCH_MIPS,
		CS_MODE_MIPS32R6 + CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		mips32R6,
		"MIPS-32R6 (Big-endian)",
	},
	{
		CS_ARCH_MIPS,
		CS_MODE_MIPS64 + CS_MODE_MIPS2 + CS_MODE_LITTLE_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		mips64SD,
		"MIPS-64-EL + Mips II (Little-endian)",
	},
	{
		CS_ARCH_MIPS,
		CS_MODE_MIPS64 + CS_MODE_LITTLE_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		mips64SD,
		"MIPS-64-EL (Little-endian)",
	},
}

var x86Tests = platforms{
	platform{
		CS_ARCH_X86,
		CS_MODE_16,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		x86Code16,
		"X86 16bit (Intel syntax)",
	},
	platform{
		CS_ARCH_X86,
		CS_MODE_32,
		[]option{{CS_OPT_SYNTAX, CS_OPT_SYNTAX_ATT}, {CS_OPT_DETAIL, CS_OPT_ON}},
		x86Code32,
		"X86 32 (AT&T syntax)",
	},
	platform{
		CS_ARCH_X86,
		CS_MODE_32,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		x86Code32,
		"X86 32 (Intel syntax)",
	},
	platform{
		CS_ARCH_X86,
		CS_MODE_64,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		x86Code64,
		"X86 64 (Intel syntax)",
	},
}

var ppcTests = platforms{
	platform{
		CS_ARCH_PPC,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		ppcCode,
		"PPC-64",
	},
	platform{
		CS_ARCH_PPC,
		CS_MODE_BIG_ENDIAN + CS_MODE_QPX,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		ppcCode2,
		"PPC-64 + QPX",
	},
}

var sysZTests = platforms{
	platform{
		CS_ARCH_SYSZ,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		sysZCode,
		"SystemZ",
	},
}

var sparcTests = platforms{
	platform{
		CS_ARCH_SPARC,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		sparcCode,
		"Sparc",
	},
	platform{
		CS_ARCH_SPARC,
		CS_MODE_BIG_ENDIAN + CS_MODE_V9,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		sparcV9Code,
		"SparcV9",
	},
}

var xcoreTests = platforms{
	platform{
		CS_ARCH_XCORE,
		CS_MODE_BIG_ENDIAN,
		[]option{{CS_OPT_DETAIL, CS_OPT_ON}},
		xcoreCode,
		"XCore",
	},
}

func dumpHex(code []byte, buf *bytes.Buffer) {
	for _, b := range code {
		// This deliberately leaves a stray space at EOL to match the C tests.
		fmt.Fprintf(buf, "0x%.2x ", b)
	}
	fmt.Fprintf(buf, "\n")
}

func saveFile(t *testing.T, filename, content string) {
	t.Helper()
	fd, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create file %s: %v", filename, err)
	}
	defer fd.Close()

	if _, err := fd.WriteString(content); err != nil {
		t.Fatalf("Failed to write to file %s: %v", filename, err)
	}
}

func assertEqual[T comparable](t *testing.T, format string, a, b T) bool {
	if a != b {
		t.Helper()
		t.Errorf(format, a, b)
	}
	return a == b
}
