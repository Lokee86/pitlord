# Dependency-Knot Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the first manually adjudicated real-corpus reference state for Pitlord's `dependency-knots` detector.

## Overview

The audit uses the same 16 pinned corpora as dependency-pressure calibration. No knot threshold or relation policy was changed while the corpus output was adjudicated.

Raw scan state at the initial detector commit (`e00e87e`) was 29 findings across eight corpora; the other eight corpora were quiet. Scoring that unchanged detector against the frozen projections yields 7 true positives, 31 true negatives, 8 false positives, 0 false negatives, 3 severity mismatches, and 0 unlabelled findings. Labelled precision is 46.7%; labelled recall is 100%.

| Corpus | Raw knots | Frozen interpretation |
| --- | ---: | --- |
| Dapper | 1 | One bounded public/internal factory cycle; allowed warning. |
| Polly | 0 | Quiet negative control. |
| Spectre.Console | 5 | One directory-granularity false positive plus four intentional catalog/domain cycles. |
| Gson | 3 | Two real cross-package knots plus one bounded public/internal helper cycle. |
| HikariCP | 2 | One intentional public/runtime cycle plus one real config/utility package cycle. |
| JMH | 0 | Quiet negative control; known generator pressure is non-cyclic. |
| Maven | 5 | Three real package knots plus two intentional API/implementation cycles. |
| jsoup | 1 | One large real package knot spanning the core DOM/select/parser/helper architecture. |
| kotlinx.coroutines | 5 | Five bounded public/internal or algorithmic runtime cycles; allowed warnings. |
| Now in Android | 0 | Quiet negative control. |
| detekt | 0 | Quiet negative control. |
| Lexicanter | 0 | Quiet negative control; known layout debt is non-cyclic. |
| DUnit LotusScript | 0 | Quiet negative control. |
| JSONParser LotusScript | 0 | Quiet; single-unit pressure cannot form a file dependency knot. |
| Volt MX LotusScript toolkit | 0 | Quiet negative control. |
| Space Rocks | 7 | Seven false architectural cycles from vendored tooling or callback/interface dispatch projection. |

## Frozen adjudication

### Required knot pressure

Seven current components are real architectural dependency knots and should remain detectable as advisory findings:

- Gson `FormattingStyle.java` component: 16 files across public JSON tree, stream, internal and binding packages.
- Gson `Gson.java` component: 19 files across public facade, annotations, internal binding and SQL adapter packages.
- HikariCP `HikariConfig.java` ↔ `PropertyElf.java`: the configuration layer uses the property utility while the utility special-cases `HikariConfig`.
- Maven compatibility artifact component: 12 files across eight legacy artifact/repository/resolver packages.
- Maven lifecycle execution-plan component: public `MavenExecutionPlan` directly owns `internal.ExecutionPlanItem`, while lifecycle internals depend back on the public plan type.
- Maven project/artifact factory component: `MavenProject` directly calls an internal lifecycle factory which in turn consumes `MavenProject`.
- jsoup core component: 59 files across seven packages. Bidirectional package dependencies are explicit; for example `nodes` imports selector types while `select` imports node types.

The large Gson, Maven and jsoup components are structural maintenance pressure even where the responsibilities are cohesive. This detector is advisory; `required` means the structural condition belongs to this detector, not that the corpus should be rewritten.

### Allowed cyclic seams

Fourteen current findings are real file/type cycles but are intentionally bounded or intrinsic to a cohesive API/runtime shape. They may remain visible, but the reference does not require them:

- Dapper: `BulkCopy` public convenience factory ↔ internal `DynamicBulkCopy` implementation.
- Spectre.Console: BoxBorder, TableBorder and TreeGuide known-instance catalogs ↔ their concrete implementations; `IHasTreeNodes` ↔ recursive `TreeNode` domain contract.
- Gson: public `ReflectionAccessFilter` constants ↔ internal helper.
- HikariCP: public `HikariDataSource` ↔ pool runtime; the SCC's reverse public-package edge includes a documentation-oriented `HikariDataSource` import in `HikariPool`.
- Maven: public Session/service request API cycle and Injector/API ↔ DI implementation cycle.
- kotlinx.coroutines: dispatcher, scope, Flow/context, SharedFlow/StateFlow and JVM thread-context public/internal implementation cycles.

Allowed components are warning-level at most in the frozen projection.

### False architectural cycles

Eight current findings must be absent from the generalized knot result:

- Spectre.Console ANSI component: all types share the same `Spectre.Console` namespace; the extra `Utilities/` directory is organizational rather than an architectural boundary.
- Space Rocks GUT component: vendored/addon test tooling, not product architecture.
- Space Rocks game-server inbound component: interface calls are projected back to adapter implementations even though the `inbound` package does not depend on the parent `networking` package.
- Space Rocks gameplay composition and realtime transport components: callback/Callable dispatch produces reverse call edges without reverse source ownership.
- Space Rocks profile HTTP and local-store components: interface/factory callbacks project reverse edges into the game-server composition root; the player-data packages do not import the game-server command package.
- Space Rocks packet-router component: function-field callbacks project calls back into the parent router file; `inbound/router.go` does not import the parent networking package.

These cases demonstrate that mutual reachability over runtime call targets is not automatically a source dependency cycle.

## Calibration implications

The frozen state identified the tuning questions without changing the reference set:

- distinguish static dependency direction from callback/interface dispatch targets;
- treat vendored/addon tooling as a separate source role;
- avoid treating directory splits as language architecture when the language exposes a stronger package/namespace boundary;
- keep intentional factory/catalog/public-internal cycles warning-level or suppressible without erasing large cross-package knots.

## Tuned detector result

The implementation was then tuned against those unchanged references. The current projection scores:

- 7 true positives;
- 39 true negatives;
- 0 false positives;
- 0 false negatives;
- 0 severity mismatches; and
- 0 unlabelled findings.

Labelled precision and recall are both 100% on the frozen required/absent expectations.

The tuned semantics are deliberately narrower than the initial detector:

- runtime `calls` do not define source dependency knots;
- imports, inheritance/implementation, trait use, overrides, includes, `depends-on`, and conversions remain source-dependency evidence;
- Arcana language namespaces and multi-file modules define architectural regions when available, with filesystem directories as fallback only;
- GUT addon code is non-production test tooling;
- tiny two-file parent/`impl`, parent/`internal`, and parent/`implementation` cycles are deferred as bounded implementation seams; and
- density alone no longer promotes a two-region knot to high severity.

The machine-readable projections are the unchanged `*.dependency-knots.json` files in [Calibration references](references/INDEX.md).

## Related docs

- [Calibration evaluation harness](HARNESS.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)
- [Current limitations](../../limits/current-limitations.md)

## Notes

This is reference corpus state, not a claim that every source cycle must be removed. Re-adjudicate only when a pinned corpus revision changes or concrete source evidence disproves a frozen judgment.
