# jsoup Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the human-adjudicated architecture state of the jsoup corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `java-jsoup`
- Source revision: `d24b16d952530f3aefe87e91c344b23fe4b8a7fc`
- Commit subject: `Link PR in changelog`
- Language: Java
- Baseline role: external clean/real-world control

The baseline records corpus state only. It does not prescribe detector changes.

## Current Pitlord observation

A repository-wide dependency-pressure scan reports 18 advisory warnings: 13 production files and 5 test files.

Package-scoped scans of the production packages containing those findings currently report:

- `src/main/java/org/jsoup/nodes/` — no findings
- `src/main/java/org/jsoup/parser/` — no findings
- `src/main/java/org/jsoup/select/` — no findings
- `src/main/java/org/jsoup/internal/` — no findings
- `src/main/java/org/jsoup/safety/` — no findings
- `src/main/java/org/jsoup/helper/` — one high finding on `HttpConnection.java`

The difference between repository-wide and package-scoped output is preserved as an observation, not treated here as ground truth.

## Human-adjudicated production state

### HTTP and helper package

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `helper/DataUtil.java` | clean | Cohesive internal data/input utility: stream opening, charset detection, parsing handoff, and related input helpers. |
| `helper/HttpConnection.java` | intentional broad seam; maintenance watch | Implements the public `Connection` session/request/response contract. The file contains the `HttpConnection` implementation plus nested `Base`, `Request`, `Response`, and `KeyVal` implementations. Its breadth is substantial, but the inspected responsibilities remain within HTTP request/session execution and the public connection abstraction. |
| `helper/W3CDom.java` | clean | Cohesive adapter between jsoup's DOM and W3C DOM/XPath/serialization surfaces. |

`HttpConnection` is the only current production candidate that remains a statistical outlier at its natural package boundary. The baseline records it as a large maintenance-watch seam, not as demonstrated unrelated-responsibility sprawl.

### DOM nodes

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `nodes/Document.java` | clean; intentional domain root | Represents the document root and document-specific behavior such as structure, output settings, parser/session context, and document helpers. |
| `nodes/Element.java` | intentional broad public domain API; maintenance watch | Central HTML element abstraction. Its broad dependency surface follows its public responsibilities for DOM traversal, manipulation, selection, attributes, text, and element operations. It is large enough to watch for maintainability pressure, but the inspected breadth is intrinsic to the `Element` abstraction rather than unrelated ownership. |
| `nodes/FormElement.java` | clean | Focused `Element` specialization for associated controls, form data, and form submission preparation. |

### Parser

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `parser/TreeBuilder.java` | clean | Focused abstract tree-building runtime: parser/tokenizer/document/stack lifecycle and common insertion mechanics. |
| `parser/HtmlTreeBuilder.java` | intentional algorithmic complexity | Cohesive HTML tree-builder implementation. Its broad DOM/parser references support one parsing algorithm and its state/context management. |
| `parser/HtmlTreeBuilderState.java` | intentional state-machine representation | Large enum whose states each embody HTML tree-builder processing and transitions. Size and dependency breadth reflect the parsing state machine rather than multiple architectural owners. |
| `parser/XmlTreeBuilder.java` | clean | Focused XML tree builder with XML namespace and insertion behavior. |

### Safety

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `safety/Cleaner.java` | clean | Focused safelist-based HTML sanitization and validation. |

### Selectors

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `select/Evaluator.java` | intentional evaluator taxonomy | Abstract selector evaluator plus cohesive nested evaluator implementations for selector predicates. File size reflects the evaluator family rather than unrelated ownership. |
| `select/QueryParser.java` | clean | Focused CSS selector parser that produces evaluator trees. |

## Non-production peer class

The five repository-wide findings under `src/test/` are test architecture, not production dependency-pressure findings in this baseline. Tests may legitimately span broad production surfaces.

Current observed test findings are:

- `src/test/java/org/jsoup/integration/ConnectTest.java`
- `src/test/java/org/jsoup/nodes/ElementTest.java`
- `src/test/java/org/jsoup/parser/HtmlParserTest.java`
- `src/test/java/org/jsoup/parser/PositionTest.java`
- `src/test/java/org/jsoup/select/SelectorTest.java`

## Ground-truth summary

The inspected jsoup production architecture should not be characterized as broadly unhealthy from the current dependency-pressure warnings.

No inspected warning demonstrated unrelated architectural responsibility sprawl. `HttpConnection` and `Element` are legitimate maintenance-watch seams because they are large central public abstractions. `HtmlTreeBuilder`, `HtmlTreeBuilderState`, and `Evaluator` carry substantial but cohesive algorithmic/taxonomic complexity. The remaining warned production classes are clean within their stated responsibilities.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Dapper](dapper.md)
- [Space Rocks Clean Control v1](space-rocks-clean-v1.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
