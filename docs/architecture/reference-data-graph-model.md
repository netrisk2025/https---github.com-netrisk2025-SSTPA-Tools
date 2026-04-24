# Reference Data Graph Model

## Purpose

This document defines the graph data model and normalized pipeline artifact shape for the SSTPA reference-data subgraph. It is based on the raw framework deliveries currently stored in:

- `reference-data/raw/mitre-attack/v18.1`
- `reference-data/raw/mitre-emb3d/v2.0.1`
- `reference-data/raw/nist-sp800-53/v5.2.0`

The design supports four goals:

1. Preserve the source frameworks as an immutable, navigable graph.
2. Preserve both intra-framework and cross-framework relationships where the raw data supports them.
3. Distinguish items that are browsable from items that are valid assignment targets for `[:REFERENCES]`.
4. Allow future user-authored reference items in the same semantic space without mutating licensed imported content.

## Raw-Data Findings That Drive The Model

### 1. MITRE ATT&CK is not safely keyable by `ExternalID` alone

The combined ATT&CK Enterprise, Mobile, and ICS corpus contains:

- repeated STIX objects across bundle files
- repeated `external_id` values across different object types

Examples found in `v18.1`:

- the same STIX object appears in multiple bundles, such as malware and data components shared by Enterprise, Mobile, and ICS
- deprecated ATT&CK mitigations can reuse a `Txxxx` identifier that also appears on a technique

Result:

- the canonical import key for ATT&CK must be the source object identifier (`stix.id`)
- `ExternalID` remains a preserved publisher field and search alias, but not the import merge key

### 2. ATT&CK hierarchy is partly explicit and partly embedded

ATT&CK relationships come from several places:

- STIX `relationship` objects:
  - `uses`
  - `mitigates`
  - `detects`
  - `subtechnique-of`
  - `targets`
  - `attributed-to`
  - `revoked-by`
- matrix object `tactic_refs`
- technique `kill_chain_phases`
- detection strategy `x_mitre_analytic_refs`
- analytic `x_mitre_log_source_references`
- data component `x_mitre_data_source_ref`

Result:

- normalization must produce edges from both explicit relationship rows and embedded object fields

### 3. ATT&CK contains valid research items that are not valid assignment targets

Examples present in the raw ATT&CK data:

- campaigns
- malware
- tools
- intrusion sets
- analytics
- detection strategies
- data sources
- data components
- ICS assets
- matrices

Result:

- the model must separate `navigable` from `selectable`
- the Reference Tool must browse all imported items, but only a subset may be assigned to SSTPA core nodes

### 4. NIST SP 800-53 has strong internal graph structure

The OSCAL catalog contains:

- families
- controls
- control enhancements
- link semantics including:
  - `related`
  - `required`
  - `incorporated-into`
  - `moved-to`

Result:

- NIST hierarchy is first-class, not just text
- enhancement and cross-control links should be imported as graph edges

### 5. Cross-framework links exist at multiple confidence levels

Strong cross-framework links found in raw data:

- ATT&CK for ICS mitigations include structured labels such as `NIST SP 800-53 Rev. 5 - AC-3; SC-7`
- 42 ATT&CK ICS mitigations yielded 50 direct edges to current NIST 5.2.0 controls

Derived but still valuable cross-framework links found in raw data:

- EMB3D threat evidence includes ATT&CK technique and software identifiers in markdown text
- ATT&CK objects contain NIST citations that point to NIST publications but not always to specific controls

Result:

- the graph must distinguish authoritative mappings from derived or citation-based links

### 6. User-authored reference content must not mutate imported items

The project requirement is to allow local reference items of the same semantic types as imported content while preserving license restrictions on imported data.

Result:

- imported items and local items must coexist in the same overall reference graph
- imported items remain immutable
- local items must use a separate authority and keyspace

## Relationship Inventory

### Packaging Assessment

The source packaging does not make internal relationships impossible, but it does make them uneven.

- NIST SP 800-53 is the cleanest source. Its OSCAL catalog is nested and explicit, so family, control, enhancement, and most internal links can be extracted directly.
- MITRE ATT&CK is internally rich but structurally mixed. Some relationships are explicit STIX `relationship` objects, while others are embedded in fields such as `tactic_refs`, `kill_chain_phases`, `x_mitre_analytic_refs`, and `x_mitre_data_source_ref`.
- MITRE EMB3D is internally simpler than ATT&CK. Its core threat, mitigation, and property relationships are explicit, but category-style navigation is mostly field-based rather than node-based.
- Cross-framework relationships are the real asymmetry. NIST is not aware of MITRE in the raw files examined. ATT&CK is only partly aware of NIST. EMB3D mentions ATT&CK and NIST in text, but usually not as native graph edges.

Conclusion:

- internal framework relationships should be materialized during normalization
- cross-framework links should be created in a second linking pass after per-framework normalization
- backend import should load explicit graph artifacts and should not be responsible for discovering relationships

### Internal Relationships By Framework

#### NIST SP 800-53 Internal Relationships

The raw OSCAL catalog encodes a strong internal hierarchy plus explicit cross-links between controls and enhancements.

```cypher
(:NistFamily)
    -[:HAS_CHILD]->
(:NistControl)
    -[:HAS_CHILD]->
(:NistControlEnhancement)

(:NistControl)
    -[:RELATED]->
(:NistControl)

(:NistControlEnhancement)
    -[:RELATED]->
(:NistControl)

(:NistControlEnhancement)
    -[:RELATED]->
(:NistControlEnhancement)

(:NistControlEnhancement)
    -[:REQUIRED]->
(:NistControl)

(:NistControl)
    -[:MOVED_TO]->
(:NistControl)
```

Observed source characteristics:

- hierarchy comes from nested OSCAL `groups`, `controls`, and child `controls`
- most semantic cross-links come from OSCAL `links`
- `related` and `required` links resolve cleanly to other catalog items
- some `incorporated-into` and `moved-to` links target OSCAL statement fragments such as `#ac-2_smt.k` rather than full control nodes

Normalization recommendation:

- create explicit item-to-item edges for resolvable family, control, and enhancement relationships
- preserve non-item fragment targets as fragment references during normalization
- only promote fragment anchors to graph nodes if the UI later needs fragment-level navigation

#### MITRE ATT&CK Internal Relationships

The raw ATT&CK bundles encode a broad internal research graph, but the relationships are split between STIX relationship rows and embedded object fields.

```cypher
(:AttackMatrix)
    -[:HAS_CHILD]->
(:AttackTactic)
    -[:HAS_CHILD]->
(:AttackTechnique)

(:AttackTechnique)
    -[:SUBTECHNIQUE_OF]->
(:AttackTechnique)

(:AttackMitigation)
    -[:MITIGATES]->
(:AttackTechnique)

(:AttackGroup)
    -[:USES]->
(:AttackTechnique)

(:AttackCampaign)
    -[:USES]->
(:AttackTechnique)

(:AttackGroup)
    -[:USES]->
(:AttackSoftware)

(:AttackCampaign)
    -[:USES]->
(:AttackSoftware)

(:AttackCampaign)
    -[:ATTRIBUTED_TO]->
(:AttackGroup)

(:AttackDetectionStrategy)
    -[:DETECTS]->
(:AttackTechnique)

(:AttackDetectionStrategy)
    -[:HAS_CHILD]->
(:AttackAnalytic)

(:AttackDataSource)
    -[:HAS_CHILD]->
(:AttackDataComponent)

(:AttackAnalytic)
    -[:USES]->
(:AttackDataComponent)

(:AttackTechnique)
    -[:TARGETS]->
(:AttackAsset)

(:AttackTechnique)
    -[:REVOKED_BY]->
(:AttackTechnique)
```

Observed source characteristics:

- matrices reference tactics through `tactic_refs`
- tactic-to-technique placement is derived from technique `kill_chain_phases`
- sub-technique hierarchy is explicit via `subtechnique-of`
- group, campaign, malware, and tool behavior comes from STIX `uses`
- detection strategy to analytic and analytic to data component links are embedded field references rather than STIX relationship rows
- data source to data component membership is also field-based
- the same STIX object can appear in multiple ATT&CK bundles, so ATT&CK must be deduplicated by STIX `id`, not by raw file membership

Normalization recommendation:

- perform source-specific extraction in `staged/` so embedded ATT&CK references become explicit candidate edges
- create the ATT&CK internal graph in normalization pass 1 after de-duplicating shared STIX objects across Enterprise, Mobile, and ICS bundles
- preserve bundle membership separately from semantic relationships

#### MITRE EMB3D Internal Relationships

The raw EMB3D STIX file is structurally smaller than ATT&CK and its core internal edges are explicit, but some navigation groupings are not first-class nodes in the source.

```cypher
(:Emb3dProperty)
    -[:SUBPROPERTY_OF]->
(:Emb3dProperty)

(:Emb3dProperty)
    -[:RELATES_TO]->
(:Emb3dThreat)

(:Emb3dMitigation)
    -[:MITIGATES]->
(:Emb3dThreat)

(:Emb3dCategory)
    -[:HAS_CHILD]->
(:Emb3dProperty)

(:Emb3dCategory)
    -[:HAS_CHILD]->
(:Emb3dThreat)
```

Observed source characteristics:

- `subproperty-of`, `relates-to`, and `mitigates` are explicit in STIX relationship objects
- property categories are field values such as `Hardware`, `System Software`, `Application Software`, and `Networking`
- threat categories are also field values rather than standalone nodes
- mitigation categories are not first-class in the source and may need to be derived indirectly through related threats or curated navigation collections

Normalization recommendation:

- import explicit property, threat, and mitigation edges directly
- derive category navigation structures during normalization rather than pretending they already exist in raw STIX

### Cross-Framework Relationships By Framework Pair

#### NIST SP 800-53 <-> MITRE ATT&CK

The source relationship is one-way and uneven.

```cypher
(:AttackMitigation)
    -[:MAPS_TO]->
(:NistControl)
```

Observed source characteristics:

- no native NIST-to-ATT&CK relationships were found in the OSCAL catalog
- ATT&CK ICS mitigations contain direct standards labels such as `NIST SP 800-53 Rev. 5 - AC-3; SC-7`
- those labels can be parsed into deterministic control-to-mitigation crosswalks
- outside of those mappings, most ATT&CK references to NIST are citations to NIST publications rather than links to specific SP 800-53 controls

Normalization recommendation:

- create direct cross-framework edges only where the ATT&CK source identifies specific NIST control IDs
- preserve general NIST publication mentions as citations, not as control nodes
- do not invent reverse NIST-to-ATT&CK edges from absence of source data

#### NIST SP 800-53 <-> MITRE EMB3D

No native item-to-item crosswalk was found between these frameworks in the raw files examined.

Observed source characteristics:

- no NIST-to-EMB3D relationships were found in the NIST OSCAL catalog
- EMB3D contains references to NIST publications such as cryptography and firmware guidance
- those are publication references, not explicit SP 800-53 control mappings

Normalization recommendation:

- preserve NIST publication mentions as citations
- do not generate `Emb3dMitigation -> NistControl` or `Emb3dThreat -> NistControl` edges unless a future curated crosswalk is supplied

#### MITRE ATT&CK <-> MITRE EMB3D

This relationship is also one-way in the raw data examined.

```cypher
(:Emb3dThreat)
    -[:REFERENCES_ATTACK_TECHNIQUE]->
(:AttackTechnique)

(:Emb3dThreat)
    -[:REFERENCES_ATTACK_SOFTWARE]->
(:AttackSoftware)
```

Observed source characteristics:

- no native ATT&CK-to-EMB3D relationships were found in ATT&CK raw bundles
- EMB3D threat evidence often includes ATT&CK technique IDs and ATT&CK software IDs inside markdown text
- those references can be deterministically resolved when the cited ATT&CK IDs exist in the imported ATT&CK release

Normalization recommendation:

- create ATT&CK/EMB3D cross-links in a dedicated linker pass after both frameworks have been normalized
- mark these links as derived from EMB3D evidence text rather than source-native STIX edges

### Recommended Pipeline Stage For Relationship Creation

The reference graph should be built in two distinct derivation passes before backend import.

| Relationship class | Examples | Recommended stage | Reason |
| --- | --- | --- | --- |
| source preservation | raw upstream files, licenses, checksums | `raw/` | preserve publisher artifacts unchanged |
| source-specific unpacking | ATT&CK embedded refs, OSCAL nesting, EMB3D category fields | `staged/` | extract framework-specific structures without yet imposing the canonical graph |
| internal framework hierarchy | NIST family/control/enhancement, ATT&CK matrix/tactic/technique, EMB3D subproperty trees | normalization pass 1 | this is where the source becomes a canonical internal graph |
| internal framework semantic edges | ATT&CK `uses`, `mitigates`, `detects`; NIST `related`, `required`; EMB3D `mitigates`, `relates-to` | normalization pass 1 | edges are still framework-local and should be resolved before any cross-framework matching |
| derived navigation structures | ATT&CK bundle collections, EMB3D categories, optional NIST fragment references | normalization pass 1 | these are internal representations of source structure, not cross-framework inferences |
| cross-framework deterministic links | ATT&CK ICS mitigation to NIST control, EMB3D evidence mention to ATT&CK ID | normalization pass 2 linker | linking works best after each framework already has stable normalized keys |
| heuristic or curated crosswalks | future manually maintained links not present in source | normalization pass 2 linker | keeps non-authoritative logic out of import and preserves provenance |
| database loading | Neo4j node and edge creation | backend import | importer should load explicit artifacts, not interpret raw source semantics |

Recommended implementation shape within the current repo layout:

- keep `raw/`, `staged/`, and `normalized/` as the top-level pipeline stages
- treat `normalized/` as two logical passes:
  - pass 1: per-framework internal graph normalization
  - pass 2: cross-framework link resolution
- if explicit separation is helpful, store those outputs under:
  - `reference-data/normalized/internal/<framework>/<version>/`
  - `reference-data/normalized/linked/<release-set>/`

Direct answer to the pipeline question:

- yes, the internal representation should be created before backend import
- no, backend import is not the right place to discover internal or cross-framework relationships
- the right place is a source-aware normalization pass followed by a separate linking pass

## Node Types And Properties

### `(:ReferenceFramework)`

One node per framework family.

Examples:

- `MITRE ATT&CK`
- `MITRE EMB3D`
- `NIST SP 800-53`

Required properties:

- `FrameworkKey`
- `FrameworkName`
- `Publisher`

### `(:ReferenceRelease)`

One node per immutable imported release.

Examples:

- `attack@18.1`
- `emb3d@2.0.1`
- `nist-sp800-53@5.2.0`

Required properties:

- `ReleaseKey`
- `FrameworkKey`
- `FrameworkVersion`
- `LicenseName`
- `LicensePath`
- `SourceManifestPath`
- `ImportedAt`
- `RawArtifactSet`

### `(:ReferenceCollection)`

A navigational grouping inside a release. Collections are read-only and not assignable.

Use cases:

- ATT&CK bundle membership: `enterprise-attack`, `mobile-attack`, `ics-attack`
- derived EMB3D category grouping: `Hardware`, `System Software`, `Application Software`, `Networking`
- future local overlays if needed

Required properties:

- `CollectionKey`
- `ReleaseKey`
- `CollectionType`
- `Name`
- `AuthorityKind`

### `(:ReferenceItem)`

Every browsable framework object. Imported framework-specific labels are attached here.

Examples:

- `(:ReferenceItem:NistFamily)`
- `(:ReferenceItem:NistControl)`
- `(:ReferenceItem:NistControlEnhancement)`
- `(:ReferenceItem:AttackMatrix)`
- `(:ReferenceItem:AttackTactic)`
- `(:ReferenceItem:AttackTechnique)`
- `(:ReferenceItem:AttackMitigation)`
- `(:ReferenceItem:AttackCampaign)`
- `(:ReferenceItem:AttackGroup)`
- `(:ReferenceItem:AttackSoftware)`
- `(:ReferenceItem:AttackDataSource)`
- `(:ReferenceItem:AttackDataComponent)`
- `(:ReferenceItem:AttackDetectionStrategy)`
- `(:ReferenceItem:AttackAnalytic)`
- `(:ReferenceItem:AttackAsset)`
- `(:ReferenceItem:Emb3dProperty)`
- `(:ReferenceItem:Emb3dThreat)`
- `(:ReferenceItem:Emb3dMitigation)`

Common property set for all `:ReferenceItem` nodes:

Required properties:

- `ReferenceKey`
- `FrameworkKey`
- `FrameworkVersion`
- `ReleaseKey`
- `AuthorityKind`
- `Mutability`
- `SourceObjectID`
- `ExternalID`
- `ExternalType`
- `CanonicalType`
- `Name`
- `ShortDescription`
- `LongDescription`
- `SourceURI`
- `Navigable`
- `Selectable`
- `AssignableTo`
- `LifecycleState`
- `Imported`
- `LastUpdated`
- `RawData`

Property rules:

- `ReferenceKey` is the stable import key and must be unique within the entire reference graph
- `AuthorityKind` is one of `imported`, `derived`, or `local`
- `Mutability` is `immutable` for imported and derived content, `editable` for local content
- `Selectable` means valid as a `[:REFERENCES]` endpoint
- `AssignableTo` is the SSTPA node-type allow-list used by the Reference Tool and backend validation
- `LifecycleState` is one of `active`, `deprecated`, `revoked`, `superseded`, or `draft`

Framework-specific `:ReferenceItem` subtypes inherit the common `:ReferenceItem` property set. The primary imported subtypes are:

- `NistFamily`
- `NistControl`
- `NistControlEnhancement`
- `AttackMatrix`
- `AttackTactic`
- `AttackTechnique`
- `AttackMitigation`
- `AttackCampaign`
- `AttackGroup`
- `AttackSoftware`
- `AttackDataSource`
- `AttackDataComponent`
- `AttackDetectionStrategy`
- `AttackAnalytic`
- `AttackAsset`
- `Emb3dProperty`
- `Emb3dThreat`
- `Emb3dMitigation`

Future local reference items may reuse these same semantic subtype labels while remaining distinct by `ReferenceKey`, `AuthorityKind`, and `Mutability`.

### `(:ReferenceCitation)`

Preserves raw bibliographic and external-reference data without pretending every citation is a framework item.

Use cases:

- ATT&CK `external_references`
- EMB3D mitigation reference blocks
- EMB3D threat evidence references
- future non-imported frameworks such as CWE, CVE, IEC 62443, or NIST publications outside SP 800-53

Required properties:

- `CitationKey`
- `SourceName`
- `ExternalID`
- `URL`
- `Description`
- `SourceField`
- `ResolutionState`

## Relationship Classes

### Structural

- `(:ReferenceFramework)-[:HAS_RELEASE]->(:ReferenceRelease)`
- `(:ReferenceRelease)-[:HAS_COLLECTION]->(:ReferenceCollection)`
- `(:ReferenceRelease)-[:HAS_ITEM]->(:ReferenceItem)`
- `(:ReferenceCollection)-[:CONTAINS]->(:ReferenceCollection|:ReferenceItem)`
- `(:ReferenceItem)-[:HAS_CHILD]->(:ReferenceItem)`

Rules:

- use `HAS_CHILD` only for true navigation hierarchy
- use `CONTAINS` for bundle or category membership

### Semantic Item-to-Item

- `(:ReferenceItem)-[:RELATED_TO]->(:ReferenceItem)`

Required edge properties:

- `RelationType`
- `RelationClass`
- `CrossFramework`
- `Derived`
- `Confidence`
- `SourceField`
- `RawValue`

Examples:

- ATT&CK `uses`
- ATT&CK `mitigates`
- ATT&CK `detects`
- ATT&CK `targets`
- ATT&CK `attributed-to`
- ATT&CK `subtechnique-of`
- ATT&CK `revoked-by`
- NIST `related`
- NIST `required`
- NIST `incorporated-into`
- NIST `moved-to`
- EMB3D `mitigates`
- EMB3D `relates-to`
- EMB3D `subproperty-of`
- ATT&CK ICS mitigation to NIST control crosswalk
- EMB3D threat evidence to ATT&CK technique or software link

Recommended property values:

- `RelationClass = structural-derived` for edges inferred from embedded source fields
- `RelationClass = source-relationship` for native relationship rows
- `RelationClass = crosswalk` for direct standards mappings
- `RelationClass = citation-resolution` for links resolved from citation content
- `Confidence = authoritative` for direct source mappings
- `Confidence = derived` for deterministic parsing
- `Confidence = heuristic` for looser future matching

### Citation

- `(:ReferenceItem)-[:CITES]->(:ReferenceCitation)`
- `(:ReferenceCitation)-[:RESOLVES_TO]->(:ReferenceItem|:ReferenceRelease|:ReferenceFramework)`

Rules:

- preserve all raw citations as citations
- only create `RESOLVES_TO` when normalization can match the citation to an imported object or release with deterministic evidence

## Navigable Versus Selectable

All imported items should be navigable. Only the following classes are selectable by default:

| Framework item class | Selectable | Assignable to |
| --- | --- | --- |
| `NistControl` | yes | `System`, `Control`, `Countermeasure` |
| `NistControlEnhancement` | no | none by default |
| `AttackTactic` | yes | `Attack` |
| `AttackTechnique` | yes | `Attack`, `Hazard` |
| `AttackMitigation` | yes | `Control`, `Countermeasure` |
| `Emb3dThreat` | yes | `Element`, `Hazard`, `Attack` |
| `Emb3dProperty` | yes | `Element`, `System` |
| `Emb3dMitigation` | yes | `Control`, `Countermeasure` |
| `AttackCampaign` | no | none |
| `AttackGroup` | no | none |
| `AttackSoftware` | no | none |
| `AttackDataSource` | no | none |
| `AttackDataComponent` | no | none |
| `AttackDetectionStrategy` | no | none |
| `AttackAnalytic` | no | none |
| `AttackAsset` | no | none |
| `NistFamily` | no | none |
| `AttackMatrix` | no | none |

This matches the SRS intent while still allowing broad research-mode navigation.

## Keying And Identity Rules

### Canonical import key

The importer must `MERGE` on `ReferenceKey`, not `ExternalID`.

Recommended construction:

- ATT&CK imported item: `<releaseKey>:stix:<stix-id>`
- EMB3D imported item: `<releaseKey>:stix:<stix-id>`
- NIST imported item: `<releaseKey>:oscal:<normalized-id>`
- local item: `local:<uuid>`

### Preserved publisher ID

`ExternalID` remains mandatory, but it is a display and search field, not the merge key.

### Lookup behavior

The backend should support:

- exact lookup by `ReferenceKey`
- exact lookup by `ExternalID`
- text lookup by `Name`, description, and aliases

Because ATT&CK `ExternalID` collisions exist, exact `ExternalID` lookup may return more than one result and must therefore remain filterable by item type.

## Versioning And Multi-Bundle Handling

### Framework version coexistence

Each imported release remains distinct.

Examples:

- `attack@18.1` can coexist with a future `attack@19.0`
- `nist-sp800-53@5.2.0` can coexist with a future `nist-sp800-53@5.3.0`

No release overwrites another release in-place.

### ATT&CK bundle deduplication

The ATT&CK release is a single logical release with multiple raw bundles.

Normalization rules:

- de-duplicate ATT&CK objects by STIX `id`
- preserve bundle membership through `ReferenceCollection` nodes
- attach the same `ReferenceItem` to more than one collection when present in multiple raw files

## Imported Versus Local Reference Content

To support user-authored reference items without violating licenses:

- imported items use `AuthorityKind = imported` and `Mutability = immutable`
- local items use `AuthorityKind = local` and `Mutability = editable`
- local items may carry the same framework-specific labels and `AssignableTo` semantics as imported items
- local items must never overwrite or reuse imported `ReferenceKey` values

Recommended future pattern:

- attach local items to the same `ReferenceFramework`
- optionally attach them to a local `ReferenceCollection` overlay for organization

## Framework-Specific Normalization Rules

### MITRE ATT&CK

Create items for:

- matrices
- tactics
- techniques
- sub-techniques
- mitigations
- campaigns
- intrusion sets
- malware
- tools
- data sources
- data components
- detection strategies
- analytics
- ICS assets

Create relationships from:

- STIX relationship rows
- matrix `tactic_refs`
- technique `kill_chain_phases`
- detection strategy `x_mitre_analytic_refs`
- analytic `x_mitre_log_source_references`
- data component `x_mitre_data_source_ref`

Cross-framework handling:

- ATT&CK ICS mitigation labels containing `NIST SP 800-53 Rev. 5 - ...` become `RELATED_TO` crosswalk edges to imported NIST controls
- ATT&CK citations to NIST publications become `ReferenceCitation` nodes

### NIST SP 800-53

Create items for:

- families
- controls
- control enhancements

Create hierarchy:

- family `HAS_CHILD` control
- control `HAS_CHILD` enhancement

Create semantic relationships from OSCAL links:

- `related`
- `required`
- `incorporated-into`
- `moved-to`

### MITRE EMB3D

Create items for:

- properties
- threats
- mitigations

Create derived collections for category navigation:

- Hardware
- System Software
- Application Software
- Networking

Create relationships from:

- `subproperty-of`
- `relates-to`
- `mitigates`

Cross-framework handling:

- EMB3D threat evidence that contains ATT&CK technique or software identifiers becomes `ReferenceCitation` plus resolved `RELATED_TO` edges when the identifier exists in the imported ATT&CK release
- EMB3D mitigation IEC 62443 mappings and NIST publication mentions remain citations unless that framework is later imported

## Normalized Pipeline Artifact Format

Each normalized release directory should contain:

- `framework.json`
- `collections.ndjson`
- `items.ndjson`
- `relationships.ndjson`
- `citations.ndjson`
- `manifest.yaml`

### `framework.json`

Contains release metadata:

- framework identity
- version
- license metadata
- raw artifact hashes
- normalization timestamp
- import compatibility version

### `collections.ndjson`

One record per `ReferenceCollection`.

Required fields:

- `collectionKey`
- `releaseKey`
- `collectionType`
- `name`
- `authorityKind`
- `parentCollectionKey`

### `items.ndjson`

One record per `ReferenceItem`.

Required fields:

- `referenceKey`
- `releaseKey`
- `frameworkKey`
- `frameworkVersion`
- `authorityKind`
- `mutability`
- `sourceObjectId`
- `externalId`
- `externalType`
- `canonicalType`
- `name`
- `shortDescription`
- `longDescription`
- `sourceUri`
- `navigable`
- `selectable`
- `assignableTo`
- `lifecycleState`
- `tags`
- `aliases`
- `collectionKeys`
- `rawData`

### `relationships.ndjson`

One record per edge between collections and items or between items and items.

Required fields:

- `fromKey`
- `toKey`
- `relationshipType`
- `relationClass`
- `crossFramework`
- `derived`
- `confidence`
- `sourceField`
- `rawValue`

### `citations.ndjson`

One record per `ReferenceCitation`, plus optional resolution records.

Required fields:

- `citationKey`
- `itemKey`
- `sourceName`
- `externalId`
- `url`
- `description`
- `sourceField`
- `resolutionState`
- `resolvedTargetKey`

## Import Rules

The backend importer should:

1. `MERGE` frameworks by `FrameworkKey`
2. `MERGE` releases by `ReleaseKey`
3. `MERGE` collections by `CollectionKey`
4. `MERGE` items by `ReferenceKey`
5. `MERGE` citations by `CitationKey`
6. create structural and semantic edges idempotently

The importer must not:

- merge on ATT&CK `ExternalID`
- drop deprecated or revoked items
- mutate imported items after load other than idempotent refresh of the same release

## Reference Tool Implications

The Reference Tool should use:

- `Navigable` to decide browse visibility
- `Selectable` and `AssignableTo` to decide assignment availability
- `LifecycleState` to mute deprecated or revoked content
- `ReferenceCitation` and `RELATED_TO` edges to drive related-item and “follow reference” panels

This enables the required behavior:

- users can browse ATT&CK campaigns and other research-only nodes
- users can follow ATT&CK and EMB3D links into NIST or ATT&CK where a deterministic mapping exists
- users can select valid NIST controls from within the same tool flow

## Required Refinements To Existing Assumptions

The current SRS statement that `ExternalID SHALL be unique within a given framework version` is not true for the full ATT&CK raw corpus now being imported.

The implementation should therefore adopt:

- `ReferenceKey` as the unique merge key
- `ExternalID` as a preserved source identifier and search alias

This keeps the SRS intent while making the import correct for the actual publisher data.
