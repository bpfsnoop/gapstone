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
	"testing"
)

func TestDetailTest(t *testing.T) {
	final := new(bytes.Buffer)

	spec_file := "test_detail.SPEC"
	var maj, min int
	if ver, err := New(0, 0); err == nil {
		maj, min = ver.Version()
		ver.Close()
	}

	t.Logf("Detailed Test. Capstone Version: %v.%v", maj, min)

	runDetailTest := func(t *testing.T, i int, platform platform) {
		t.Logf("%2d> %s", i, platform.comment)
		if shouldSkipPlatform(platform.comment) {
			t.Skipf("Skipping platform: %s", platform.comment)
			return
		}
		engine, err := New(platform.arch, platform.mode)
		if err != nil {
			t.Errorf("Failed to initialize engine %v", err)
			return
		}
		defer engine.Close()

		for _, opt := range platform.options {
			engine.SetOption(opt.ty, opt.value)
		}

		insns, err := engine.Disasm([]byte(platform.code), address, 0)
		if err == nil {
			fmt.Fprintf(final, "****************\n")
			fmt.Fprintf(final, "Platform: %s\n", platform.comment)
			fmt.Fprintf(final, "Code: ")
			dumpHex([]byte(platform.code), final)
			fmt.Fprintf(final, "Disasm:\n")
			for _, insn := range insns {
				fmt.Fprintf(
					final,
					"0x%x:\t%s\t\t%s // insn-ID: %v, insn-mnem: %s\n",
					insn.Address,
					insn.Mnemonic,
					insn.OpStr,
					insn.Id,
					engine.InsnName(insn.Id),
				)
				if len(insn.RegistersRead) > 0 {
					fmt.Fprint(final, "\tImplicit registers read: ")
					for _, reg := range insn.RegistersRead {
						fmt.Fprintf(final, "%s ", engine.RegName(reg))
					}
					fmt.Fprintf(final, "\n")
				}
				if len(insn.RegistersWritten) > 0 {
					fmt.Fprint(final, "\tImplicit registers modified: ")
					for _, reg := range insn.RegistersWritten {
						fmt.Fprintf(final, "%s ", engine.RegName(reg))
					}
					fmt.Fprintf(final, "\n")
				}
				if len(insn.Groups) > 0 {
					fmt.Fprintf(final, "\tThis instruction belongs to groups: ")
					for _, grp := range insn.Groups {
						fmt.Fprintf(final, "%v ", engine.GroupName(grp))
					}
					fmt.Fprintf(final, "\n")
				}
			}
			fmt.Fprintf(final, "0x%x:\n", insns[len(insns)-1].Address+insns[len(insns)-1].Size)
			fmt.Fprintf(final, "\n")
		} else {
			t.Errorf("Disassembly error: %v\n", err)
		}
	}

	for i, platform := range detailTests {
		t.Run(platform.comment, func(t *testing.T) {
			runDetailTest(t, i, platform)
		})
	}

	spec, err := os.ReadFile(spec_file)
	if err != nil {
		t.Errorf("Cannot read spec file %v: %v", spec_file, err)
	}
	if fs := final.String(); string(spec) != fs {
		saveFile(t, spec_file+".test", fs)
		t.Skip("Output failed to match spec!")
	} else {
		t.Logf("Clean diff with %v.\n", spec_file)
	}
}
