# Volt MX LotusScript Toolkit Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the Volt MX LotusScript toolkit corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `volt-mx-ls-toolkit`
- Source revision: `48d0c5b14cef215dabeefa3c182caa11036c0ad8`
- Commit subject: `Merge pull request #4 from HCL-TECH-SOFTWARE/feature/pswithers`
- Language: LotusScript
- Baseline role: external integration-toolkit control

The baseline records corpus state only. It does not prescribe detector changes.

## Current Pitlord observation

The repository-wide dependency-pressure scan reports zero findings.

## Manually adjudicated production state

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `notes/Code/ScriptLibraries/VoltMXHttpHelper.lss` | clean; intentional integration seam | Focused Notes/Volt MX request helper: resolves target documents, applies payload updates and builds response fields. |
| `notes/Code/Agents/(RunAgentAsync).lsa` | clean; intentional submission seam | HTTP/request-to-asynchronous-agent submission boundary. |
| `notes/Code/Agents/ProcessAsyncRequests.lsa` | maintenance watch | Small scheduled worker that owns queue traversal, database/agent routing, execution and request-status transition. The responsibilities are coherent for a dispatcher, but the implicit document-field workflow and failure boundary make it change-sensitive. |
| `node-red/UpdateContactsAndThreadsBE.lsa` | maintenance watch | Cohesive endpoint/integration unit with broad request validation, document/view mutation and response responsibilities. |
| `node-red/EchoDoc.lsa` | clean | Small HTTP endpoint adapter with narrow ownership. |

The submission and dispatcher agents share an implicit document-based request contract. That is frozen as maintenance risk, not as demonstrated cross-unit dependency pathology.

## Ground-truth summary

The toolkit is mostly a collection of intentionally broad integration seams. No broad architectural pathology was established. The async dispatcher and its implicit workflow contract are maintenance watches that a cross-file dependency-pressure detector may legitimately miss despite the repository-wide zero result.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
