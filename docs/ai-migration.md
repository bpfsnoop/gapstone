# AI Migration Playbook (Gapstone ↔ Capstone)

This document defines the **repeatable migration flow** to upgrade Gapstone to future Capstone versions.

## Scope

Use this when any of these change:
- Capstone major/minor version
- Capstone Python constant layout (`*_const.py`)
- Capstone operand/decomposer structs in C headers
- Expected test/spec output formatting

Primary goals:
- Regenerate all `*_constants.go` via `./genconst` (no manual constant drift)
- Regenerate `*.SPEC` via `./genspec`
- Update decomposers/tests to new Capstone APIs
- Keep `make gotest` green

---

## Ground Rules

1. **Generated files are authoritative**
   - Do not hand-edit generated `*_constants.go` unless it is an emergency unblock.
   - Fix generation logic in `genconst` instead.

2. **Specs are test-driven**
   - `genspec` must source `.SPEC` from executable test output against the current Capstone build.

3. **Architecture-by-architecture migration**
   - For each arch: constants → decomposer/tests → spec refresh → focused test.

4. **Fix root causes first**
   - Prefer script/runtime adaptation over one-off constant aliases.

---

## Prerequisites

From repo root:

```bash
git submodule update --init --recursive
cd capstone
cmake -B build \
  -DCMAKE_BUILD_TYPE=Release \
  -DCAPSTONE_USE_ARCH_REGISTRATION=1 \
  -DCAPSTONE_ARCHITECTURE_DEFAULT=1 \
  -DCAPSTONE_BUILD_SHARED_LIBS=1 \
  -DCAPSTONE_BUILD_CSTOOL=0
cmake --build build
cd ..
```

Sanity check:

```bash
make gotest
```

Expect failures during migration; this establishes the baseline.

---

## Standard Migration Flow

### Phase 1: Regenerate constants

```bash
./genconst ./capstone/bindings/python/capstone
```

Expected behavior of `genconst`:
- Accept direct constants, aliases, and expressions from Python files
- Handle compatibility wrappers (e.g. arm64/sysz alias modules)
- Deduplicate repeated names in source files
- Emit Go-safe const expressions and format files

If a generated constant fails to compile, fix `genconst` rules and regenerate.

---

### Phase 2: Adapt shared core constants (if needed)

When Capstone adds cross-arch operand primitives (example: `CS_OP_SPECIAL`, `CS_OP_MEM_REG`, `CS_OP_MEM_IMM`), ensure they are available in `engine_constants.go`.

Then re-run:

```bash
./genconst ./capstone/bindings/python/capstone
```

---

### Phase 3: Regenerate SPEC files

```bash
./genspec ./capstone/build
```

Current contract of `genspec`:
- Runs selected Go tests (e.g. `TestArm`, `TestX86`, etc.)
- Reads `*.SPEC.test` outputs produced on mismatch
- Promotes them to authoritative `*.SPEC`

If `genspec` fails because tests do not compile, continue with Phase 4 first.

---

### Phase 4: Decomposer and test API migration

Compile/test failures usually reveal C API changes in:
- Operand enums (renames/splits)
- Union field layout changes
- Struct members (e.g. `writeback` → `post_index`)
- Arch-specific operand sub-structures

Approach:
1. Read corresponding Capstone header in `capstone/include/capstone/*.h`
2. Update Go decomposer (`*_decomposer.go`) mapping logic
3. Update arch tests (`*_decomposer_test.go`) for enum/name changes
4. Re-run targeted test

Example targeted run:

```bash
CGO_CFLAGS='-O1 -I$(pwd)/capstone/include' \
CGO_LDFLAGS='-O1 -g -L$(pwd)/capstone/build -lcapstone' \
go test -run '^TestArm64$' -count=1 .
```

---

### Phase 5: Architecture loop

Process each architecture independently:

1. Regenerate constants (`genconst`)
2. Fix compile/runtime deltas for this arch
3. Regenerate/validate related `.SPEC`
4. Run focused test

Suggested order (fast signal first):
1. x86
2. arm
3. arm64
4. mips
5. ppc
6. sysz
7. sparc
8. xcore
9. remaining optional/legacy arches

---

### Phase 6: Full validation

```bash
make gotest
```

Acceptance criteria:
- `make gotest` passes
- No `*.SPEC.test` leftovers
- Generated files come only from scripts
- No manual compatibility hacks remain unexplained

---

## Failure Playbook

### A) `could not determine what C.X refers to`

Likely causes:
- Generated constant points to missing C symbol
- Python source uses alias/expression not direct C constant

Action:
- Fix `genconst` transformation logic
- Regenerate constants

---

### B) Missing operand enums in tests (e.g. old names removed)

Action:
- Map old semantic usage to new enum set in tests/decomposer
- Prefer supporting all relevant new equivalents (grouped cases)

---

### C) Decomposer field no longer exists

Action:
- Inspect current header struct
- Update Go struct field extraction and memory mapping
- Re-run arch-focused tests

---

### D) `genspec` fails early

Action:
- Treat as compile/runtime blocker in tests
- Fix failing test/decomposer first
- Re-run `genspec`

---

## Update Checklist for Future Versions

- [ ] Capstone submodule updated
- [ ] Capstone build artifacts regenerated
- [ ] `genconst` executed successfully
- [ ] Wrapper/alias modules still handled (`arm64_const.py`, `sysz_const.py`-style patterns)
- [ ] `engine_constants.go` includes newly required shared constants
- [ ] `genspec` executed successfully
- [ ] Decomposer mappings updated for changed headers
- [ ] Arch-focused tests pass
- [ ] `make gotest` fully green
- [ ] No pending `*.SPEC.test` files

---

## Maintenance Notes

- Keep this document and scripts aligned: when `genconst`/`genspec` behavior changes, update this file in the same PR.
- Prefer adding small migration-safe transformations in scripts over post-generation manual edits.
- If Capstone changes test infrastructure again, preserve the contract: **SPECs must be reproducible from executable tests in this repo.**
