# Archived Status Reports (docs/status/archived/)

Closed history. Every file here was fully annotated inline (docs-health canon:
`~~item~~ done at <hash>` / `Won't implement` / `NOT-DO` / `superseded — <successor>`)
and passed `check-rows --strict` + `bare-rows` BEFORE the move. Open items were routed
to `TODO_LIST.md` / `ROADMAP.md` at sweep time — read nothing here for current work.

**Sweep:** 2026-11-07 docs-health AUDIT (18 files annotated, 18 archived).

| File | Classification | Deciding reason |
| ---- | -------------- | --------------- |
| 2026-01-14_01-38_modernize-fix-compilation-and-test-status.md | ARCHIVE | All items resolved: test failures aligned (`b767071`), integration suite recreated (`8edf0b1`), modernize superseded by Go 1.27 toolchain |
| 2026-03-20_18-07_COMPREHENSIVE_STATUS.md | ARCHIVE | Earliest of the same-evening trio; every issue fixed same evening (`b767071`, `7e85f9e`, `9bae006`) or by the v2.17 waves (`e85bb4f`, `9f4bd71`) |
| 2026-03-20_21-34_COMPREHENSIVE_STATUS_REPORT.md | ARCHIVE | Near-duplicate of the trio; all open work resolved/routed (homebrew → TODO_LIST M18, signs → ROADMAP) |
| 2026-03-20_21-36_COMPREHENSIVE_STATUS.md | ARCHIVE | Most precise of the trio; all findings verified resolved (`b767071`, `a7680ac`, `8edf0b1`) |
| 2026-03-20_22-19_TUI-Implementation-Complete.md | ARCHIVE | TUI shipped (`7e85f9e`, huh v2 `52a5a96`); wishlist routed to ROADMAP |
| 2026-03-21_02-07_HOW_TO_GOLANG-Migration-Status.md | ARCHIVE | koanf migration completed same day (`50c0feb`, `a45d663`); backlog routed |
| 2026-03-21_03-23_COMPREHENSIVE_STATUS_REPORT.md | ARCHIVE | TOP-25 resolved by `9bae006`, `9703270`, `9f4bd71`, `8edf0b1`; rest routed |
| 2026-03-21_08-47_INIT_WIZARD_IMPROVEMENTS_STATUS.md | ARCHIVE | Session work committed (`50c0feb`); Windows-defaults question answered by domain validation (gotcha 15) |
| 2026-03-21_10-33_TEST_FIXES_AND_ARCHITECTURE_ANALYSIS.md | ARCHIVE | Test fixes landed (`a7680ac`); template split-brain killed (`9f4bd71`); coverage routed |
| 2026-03-21_10-37_COMPREHENSIVE_STATUS_REPORT.md | ARCHIVE | Superset of 10-33; every priority resolved/routed |
| 2026-03-25_21-58_CHARMBRACELET_LOG_V2_MIGRATION_COMPLETE.md | ARCHIVE | Migration done (`c99a279`); huh v1/v2 question answered next day (`52a5a96`) |
| 2026-04-30_10-16_GO_BRANDED_ID_MIGRATION_COMPLETE.md | ARCHIVE | Migration done (`cc8a753`, `81d8844`); Top-25 resolved/routed |
| 2026-05-28_06-05_goreleaser-session3-release-verification.md | ARCHIVE | Blockers cleared: arch-lint config (`c6d0ec1`), release proof (`969259e`), E2E (`8edf0b1`); coverage routed |
| 2026-06-30_23-01_BUILD_RECOVERY_AND_ARCHITECTURE_FIX.html | ARCHIVE | Root causes fixed (`41aaf86`); resolution banner added inline (HTML: `<del>` + `~~` canon block); Top-25 resolved/routed |
| 2026-07-28_13-01_deduplication-session-review.md | ARCHIVE | All hygiene items landed (`3cb91da`, `9703270`); remainder routed (M17/M26) |
| 2026-08-17_00-04_WAVE1-GORELEASER-MODERNIZATION-STATUS.md | ARCHIVE | Waves completed (`e85bb4f`…`969259e`); Wave-6 remainder routed to TODO_LIST (M17-M28) |
| 2026-08-17_01-15_WAVES-2-4-EXECUTION-STATUS.md | ARCHIVE | Waves 2-4 + M13-M16 verified landed; open items routed |

**Kept out of the archive on purpose:**
- `docs/planning/2026-08-16_18-45_SUPERB-WIZARD-MODERNIZATION.md` — fully annotated (all 178 tasks struck/routed) but kept in `docs/planning/` as the active execution record; TODO_LIST.md cites its M-numbers.
