# Spectre.Console Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the Spectre.Console corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `cs-spectre`
- Source revision: `555a00ae317abc9b7a49100cda3251584b0a8861`
- Commit subject: `docs: add Turkish README translation`
- Language: C#/.NET
- Baseline role: external real-world control

The baseline records corpus state only. It does not prescribe detector changes.

## Current Pitlord observation

A repository-wide dependency-pressure scan reports 20 advisory findings: 8 high and 12 warning. Nineteen are production files and one is a test.

## Manually adjudicated production state

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `src/Spectre.Console/AnsiConsole.cs` | clean; intentional broad public facade | Small fragment of the partial `AnsiConsole` type. The public facade is deliberately distributed across `AnsiConsole.*.cs` files for prompts, live rendering, markup, progress, recording, screen, state and writing. The high file fan-out does not represent one 97-dependency implementation unit. |
| `src/Spectre.Console/BoxBorder.Known.cs` | clean; intentional catalog | Registry of built-in box-border implementations. |
| `src/Spectre.Console/TableBorder.Known.cs` | clean; intentional catalog | Registry of built-in table-border implementations. |
| `src/Spectre.Console/Prompts/MultiSelectionPrompt.cs` | maintenance watch | Cohesive but stateful interactive prompt controller covering navigation, input, paging and rendering. |
| `src/Spectre.Console/Prompts/SelectionPrompt.cs` | maintenance watch | Cohesive single-selection counterpart with similar interaction state. |
| `src/Spectre.Console/Widgets/Panel.cs` | clean | Cohesive panel widget and configuration surface. |
| `src/Spectre.Console/Widgets/Table/TableRenderer.cs` | clean; intentional renderer seam | Focused table rendering engine. |
| `src/Extensions/Spectre.Console.Json/JsonParser.cs` | clean | Cohesive JSON tokenizer/parser pipeline for the extension. |
| `src/Spectre.Console/Enrichment/ProfileEnricher.cs` | maintenance watch | Purposeful profile/environment enrichment seam whose default/global policy makes it change-sensitive. |
| `src/Spectre.Console/Live/LiveRenderable.cs` | clean | Focused live-rendering state/delegation abstraction. |
| `src/Spectre.Console/Live/Progress/Renderers/DefaultProgressRenderer.cs` | maintenance watch | Cohesive but operationally dense progress/live renderer. |
| `src/Spectre.Console/Widgets/Charts/BarChart.cs` | clean | Cohesive public chart abstraction. |
| `src/Spectre.Console/Widgets/Charts/BreakdownChart.cs` | clean | Cohesive public chart abstraction. |
| `src/Spectre.Console/Widgets/Exceptions/ExceptionRenderableBuilder.cs` | maintenance watch | Dense but cohesive exception-formatting pipeline. |
| `src/Spectre.Console/Widgets/Layout/Layout.cs` | maintenance watch | Cohesive layout tree/region allocation engine with central state and constraint logic. |
| `src/Spectre.Console/Widgets/Paragraph.cs` | clean | Cohesive paragraph measurement/rendering abstraction. |
| `src/Spectre.Console/Widgets/Rule.cs` | clean | Narrow separator/rule widget. |
| `src/Spectre.Console/Widgets/Table/Table.cs` | maintenance watch | Intentional public table model/facade with a broad mutation/configuration surface. |
| `src/Spectre.Console/Widgets/Tree.cs` | clean | Cohesive tree widget. |

No inspected production finding established unrelated architectural responsibility sprawl. The interactive/rendering controllers above are maintenance watches because of state and API breadth, not because their dependencies cross unrelated domains.

## Non-production peer class

`src/Spectre.Console.Tests/Unit/AlternateScreenTests.cs` is test architecture and is not a production dependency-pressure finding in this baseline.

## Ground-truth summary

Spectre.Console is a broad UI/rendering library with many intentionally central public widgets, renderers, catalogs and facades. `AnsiConsole.cs` in particular must be interpreted as part of a partial public type rather than as a standalone 97-dependency unit. The audit identifies several stateful rendering and prompt components as maintenance watches, but no current production finding is frozen as demonstrated architectural pathology.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
