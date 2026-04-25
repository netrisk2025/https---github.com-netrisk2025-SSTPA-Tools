SSTPA Tool - Software Requirements Specification (SRS) Version: 0.5.7 Date: April 23, 2026.

Permission is granted to collaborators and contractors working under authorization from Nicholas Triska to use, reproduce, and modify this document for the purpose of developing the SSTPA Tool and its derivatives.

Any distribution or reuse outside this scope requires prior written consent.



# 1.  Introduction

## 1.1  SRS Definitions

Purpose This Software Requirements Specification (SRS) defines the complete functional and non-functional requirements for the Systems Security-Theoretic Process Analysis (SSTPA) Tool version. This document is intended to be the single source of truth for all project stakeholders, including developers, testers, and project managers. It describes the system's features, capabilities, operational environment, and constraints, ensuring a common understanding and guiding the design, implementation, and verification of the software.

The imperative terminology used in this SRS is standard, but defined here for clarity.


-  Statements are line-items or groups of line items (e.g. this list).  Statements are requirements when they contain an imperative listed below.

-  "SHALL" used in a statement indicates its implementation is mandatory and its correct behavior must be tested.

-  "Should" used in a statement indicates it is treated as "SHALL" unless justification is provided and permission granted to omit or defer this requirement.

-  "Will" used in a statement indicates an expected behavior which occurs as the result of other requirements and therefore needs no special action.  Where a "Will" statement is likely not to occur, notification and explanation must be given.

-  "May" used in a statement indicates the requirement is optional.  If employed, it is treated as a "SHALL" but no justification is needed for omitting it.

-  Statements without an imperative are either background information informing design, definitions or headings used for organization.

This SRS tries to describe components as they are to be implemented, specifically as individual instances which forces description of components to be singular.  In most cases, the correct interpretation is the plural.  For example the statement:  "**Requirement** specifies **Countermeasure**" is to be interpreted as one or many **Requirement** components specifies one or many **Countermeasure** components.



## 1.2 Overview of SSTPA

Systems Security Theoretic Process Analysis (SSTPA) is a systems security engineering methodology derived from Systems Theoretic Process Analysis (STPA) by Nicholas Triska in 2025.  STPA (System-Theoretic Process Analysis) is a relatively new hazard analysis technique based on an extended model of accident causation. In addition to component failures, STPA assumes that accidents can also be caused by unsafe interactions of system components (elements), none of which may have failed. STPA was developed by Nancy Leveson at MIT where it continues to evolve.


Simplistically, SSTPA extends STPA by introducing the concept of *^Asset**, Criticality and the Security Attributes (Confidentiality, Integrity, Availability, Authenticity, Non-Repudiation, Durability, and Trustworthy) to make theoretic process analysis useful for System Security.  STPA focuses on Safety where the core idea is that a **Loss* may be caused by a **Hazard**.  STPA implies one Criticality, Safety.  SSTPA is **Asset** centric; **Loss** occurs when the Security Assurance on an **Asset** is compromised and conditions where this may occur are **Hazard**. **Hazard** conditions should be isolated by **StateBehavior** and mitigated by **Countermeasure** realized through **Requirement**. 

Criticality is a regime.  Pragmatically, a regime is an environment with decision makers, shared values, and accepted rules/processes.  The Safety Critical regime is similar to the Flight Critical regime but very different from the system security regimes which include Mission, Cyber-Security, Anti-Tamper, Software/Hardware Assurance, Export Control, Privacy, Surety, and others.  Each regime is mediated by experts who interact with regime decision makers to navigate regime specific certification / approval.   The SSTPA methodology is intended for use by those regime experts to realize systems that satisfy regime goals, show due diligence and speed the certification or approval process.     
 
SSTPA extends the theoretic process analysis with rigorous systems engineering methodology.  STPA develops **Control** which "must" be integrated for the system to be acceptable.  SSTPA maps **Countermeasure** which satisfies the **Control** and is specified by **Requirement** verified by **Verification**.  The **Requirement** how organizations realize engineered systems.  Further, SSTPA associates **Validation** to **System** assuring the realized system faithfully executes the **System** **Purpose**.

The SSTPA methodology mitigates the complexity around system security analysis and design in large real-world engineered systems.  These are developed top-down in a stepwise formal gated process rather than "all at once".  The SSTPA Methodology accompanies the traditional engineering process and combats complexity by focusing analysis and development of requirements at the level of a System of Interest (SoI).  This approach enhances the System Specification tree by providing (as output) system security requirements for each leaf and branch on the specification tree.    


## 1.2.1 Overview of SSTPA Tools

SSTPA Tools is intended to scale the SSTPA Methodology to the largest, most complex real-world systems.  

It will provide an ACID graph database backend capable of meeting the system security needs of a large and dispersed engineering teams in the development of a large hierarchical system.  

Backend will contain a database for project information, System information such as User names, and email addresses which will map to Owner and Creator properties to data objects and Reference data.

The backend will be developed to contain reference frameworks from MITRE (ATT&CK and EMB3D and the NIST SP800-53 Controls Catalog). Add-on Tools will allow users to search, review and associate this data via reference to valid nodes. 
 
It will have backend tools to support telemetry to include a separately accessible dashboard for telemetry display using Grafana.  Through this interface an Admin User will be able to select and save System Data.
 
It will provide a desktop Graphic User Interface (GUI) application which connects to a Backend and navigates the System nodes in the dataset, displays data from sub-graphs on selected system (the SoI), Displays sub-graph node properties in data drawers, edits nodes and their properties and commits changes to the backend.
 
The GUI will have Add-on Tools to operate on specific Nodes to:
Navigate the System Hierarch and select a System of Interest (SoI) or clone specific node properties.
Review and associate Reference Data to Valid nodes in the SoI.
Display and manage the Hierarchy of Requirements in SysML 2 compliant.
Display and manage State Transition diagrams in SysML 2 compliant using existing (:State) nodes and [:TRANSITIONS_TO] relationships. 
Display and Manage Functional Flow diagrams in SysML 2 compliant.
Develop and display the STPA Control Flow of an SoI by assigning Functions and Interfaces to STPA Roles.
Develop, Display and analyze Asset Loss as an Attack Tree.


The GUI will produce reports which include System Specifications and System descriptions based on Backend data.

### 1.2.2 SSTPA Tool Theory of Execution
**System** fail when they either stop producing value or produce harm (negative value).  This can happen with failure of function (not modeled in SSTPA Tools and not its purpose) or when **Security Assurances** are removed from **Asset** by **Attack** or flaws in the System's exposed in **Environment** subject to **Hazard**.    System Security is responsible for:
Identifying **Assets**
Assigning criticality regimes to each **Asset** as a property
Assigning Security Assurances to each **Asset** as a property
Identifying **Hazard** to each **Asset**
Developing **Control** and relating them to **Asset**
Developing **Countermeasure** to satisfy **Control** assurance
Developing **Requirement** to realize the **Countermeasure**.


Key to managing complexity is the focus on the "system" as the directed analytical graph of "Systems of Interest (SoI)". Each SoI forms a sub-graph of identical form which are analyzed by the tool independently.  The tool allows expert  decomposition **System** into subordinate **System** in a hierarchical (tree) structure using a rich system model consisting of **Environment**, **Function**, **Interface**, **Element**, **Purpose**, **State**, **ControlStructure**, and **Asset** nodes each a large set of default and user configurable properties.  The hierarchy is created by declaring an **Element** as a **System** which causes a new sub-graph to be created using the full rich system model.  The **System** in this sub-graph is parented to that **Element** from the superior sub-graph. In this way, the most complex systems are modeled. 

SSTPA Tools contains a main GUI for creating new Nodes, assigning property values and associating relationships. A key aspect of SSTPA Tools is all data is "owned" by a User.  Any user may change data butwhen a user changes data or relationships they do not own, the owner is notified via a internal message center tool.  SSTPA Tools also contains a set of "Add-on Tools" which perform analysis and assist in system development and analysis using SysML 2.0 MBSE compliant visualizations where applicable.  

The core innovation of the SSTPA Tool is the treatment of The **System** as the summation of  **System** components addressed individually as Systems of Interest (SoI).  The Graphic User Interface (GUI) will allow navigation through the hierarchy and allow selection of a single SOI.  The components of that SOI will be organized by type and presented to the User/Analyst for edit and display.  

#### 1.2.2.1  SSTPA Tools Work Flow
SSTPA Tools is intended to be a human centered set of tools allowing experts to design, specify and realize intrinsically safe and secure systems.  The intent is to be as flexible as possible to allow the human ingenuity of the users to identify and resolve the complexities of system security and create a verifiably safe and secure system which produces evidence to satisfy external certification authorities in multiple criticality domains (safety, flight, surety, mission and security).  SSTPA Tools primarily support a top-down system decomposition with bottom-up system realization. Future versions of SSTPA Tools will have a Sandbox capability allowing users to develop an isolated System independently for later integration with the full model.  

The expected typical work-flow is:
1.  Organization is contracted to realize a "Capability" valued by a customer or client.
2.  Customer or client provides the organization with a description of the "Capability", a Specification containing capability requirements, and a Statement of Work (SoW) identifying criticalities and certifications needed the organization needs to meet.
3.  The Organization has or purchases SSTPA Tools (good choice!) to use in parallel with DOORS and MBSE tools such as No-Magic CAMEO.
4.  Installer uses the SSTPA Tools Installer to install SSTPA Tools on a single computer (MVP version) and becomes the Admin and SSTPA Tools first User.
5.  Members of the Engineering team register with SSTPA Tools as Users via the Admin.
6.  A User creates the Capability from customer provided description, Specification and Statement of Work and populates Capability Requirement Nodes.
7.  Organization creates Tier 1 architecture (Systems Elements, Functions and Interfaces) identifies Assets, criticalities and assurances outside the SSTPA Tool
8.  User, uses SSTPA Tool to capture Tier-1 architecture and uses Add-on Requirements Tool to allocate Capability Requirement to Tier-1 Systems ( allocated to System-->Purpose).
9.  User Uses SSTPA Tool derive and allocate requirements from Purpose to Interfaces Functions and Elements within each SoI.
10. Users Use SSTPA Tools to identify Connections to other Systems and assigned Interfaces to participate.
11. User Specifies Tier-1 Connections within SSTPA Tool, develops System Environment, Purpose Constraints, Hazards and States for each SoI.
12. Users use SSTPA Tools to develop Controls and Countermeasures sufficient to protect Assets from Loss
13. Users develop validation criteria which if met assure the System is functional, safe and secure
14. Users derive child Systems from Elements which instantiates Core System data model for the new SoI and copies all requirements allocated to the Element into the Purpose of the Child System.
15. User uses SSTPA Tool to repeat steps 9-14 until the capability is decomposed to the point where it can be realized.
16. Users develop Verification procedures for the lowest tier systems to satisfy requirements and begin the process of realization
17. Organization realizes the highest tier system and performs verification of requirements (requirements implemented correctly) then validates the System (system meets intended purpose)
18. User uses SSTPA Tools to generate body of evidence to assure external certification authority the System is Safe and Secure and acceptable for the criticality domains it must work in. 
19. The organization integrates the lower tier systems with their peers to create the next lower tier System
20. User uses SSTPA Tools to again create Verification procedures for Requirements of that lower tier System.
21. Organization performs verification of requirements (requirements implemented correctly) then validates the System (system meets intended purpose)
22. User again uses SSTPA Tools to generate body of evidence to assure external certification authority the System is Safe and Secure and acceptable for the criticality domains it must work in. 
23. Cycle repeats until the entire Capability is realized, verified, validated and certified.
24  Organization uses SSTPA Tool to support capability sustainment throughout lifecycle to include continuous certification and accreditation.


### 1.2.3 SSTPA Tool in the Systems Engineering Space

DOORS is a powerful tool for managing requirements, and can associate cyber-security properties (a task it is never challenged to do) its relational database structure allows for only storage and recall.  SSTPA Tool with its graph-based structure allows experts to surface insights down and across the hierarchy.

CAMEO is a powerful Model Based Systems Engineering (MBSE) tool but its focus is on individual system design and cannot be easily adapted for cyber-security (it has been tried).  SSTPA Tool replicates CAMEOs requirements diagraming capability for System Security.

Both these tools and others like them (e.g. the entire IBM Rational series of Tools) are focused on systems engineering of the primary behaviors of a System.  Primary behaviors are those the customer wants and is willing to pay for.  These tools do not do a good job developing System Security behaviors needed to protect Assets needed to perform the system's primary behaviors.  These tools can configuration manage System Security requirements and model System Security functions which align with SysML but that is not their purpose.  SSTPA Tools intends to be orthogonal to these tools by focusing on System Security.  It is expected the users of SSTPA Tools will  likely NOT enter primary functional requirements into SSTPA Tool or use it to manage the primary behavior of the Project.  SSTPA Tools may need to model primary behavior functions when they have associated criticality.  The focus of SSTPA Tools is not to design the primary behavior, but develop controls, countermeasures and requirements needed to assure it meets its criticality  assurances and that these attributes can be effectively documented to achieve Certification and/or Approval to Operate.  

SSTPA is derived from STPA and is also intended to support Safety analysis, Flight Criticality analysis, Mission Criticality analysis, Surety analysis as well as Cyber-Security.  The term 'Surety" is used to include special purpose certification needs to include: Nuclear surety, Rust abatement, pharmaceutical purity, Cyber-Safe etc... Use in these domains may need additional add-on tools and reports, but hte structure of SSTPA Tools should be sound enough to address these needs.


## 1.3 SSTPA Tool Core System Data Model
SSTPA Tools  will have the following data models:
Core Data Model: supports the primary purpose of SSSTPA Tool modeling large complex systems
Tool Data Model:  utility data for SSTPA Tool internal use to include User information, Version number, licenses for components etc... .
Reference Data Model:  MITRE ATT&CK data set, EMB3D data set, NIST 800-53 catalog of Controls data set and user created data in these formats and additional formats to persist across the entire hierarchy. Help  Help Data Model:  Used to support Help  and tutorial functions.


The Core Data Model is the authoritative graph model for representing Systems of Interest (SoI), their decomposition, security analysis, requirements traceability, control structure, asset assurance, loss analysis, and assurance-case evidence. 

The Core Data Model SHALL be implemented as a Neo4j graph. The Backend SHALL validate all node labels, relationship types, relationship direction, cardinality, SoI membership, and recursive traversal constraints before committing graph mutations. 

The Core Data Model SHALL be the single authoritative schema used by the Backend, Frontend, Add-on Tools, reports, and validation logic.

### 1.3.1 Canonical Modeling Concepts 

#### 1.3.1.1 System of Interest 
A System of Interest (SoI) is the analytical sub-graph rooted at exactly one (:System) node. 
All nodes created as part of that SoI SHALL share the same HID Index as the root (:System), except child (:System) nodes created through (:Element)-[:PARENTS]->(:System). 
The SoI boundary SHALL be determined by HID Index. 
The Backend SHALL treat HID Index as the canonical SoI membership indicator. 
Cross-SoI relationships SHALL be prohibited unless explicitly allowed by this SRS. 

#### 1.3.1.2 Behavior 
System behavior SHALL be represented only by: 
- (:Function), for behavior internal to the SoI 
- (:Interface), for behavior exposed to, or interacting with, other Systems 

The Core Data Model SHALL NOT define a separate (:Behavior) node. 

#### 1.3.1.3 Purpose, Assets, and Security Assurance 
(:Purpose) represents human-imposed intent for the engineered System. 

(:Asset) represents something valuable having Criticality that requires Assurance. 

(:Security) represents the security view containing Controls and Countermeasures used to protect Assets. 

Purpose is realized by Requirements and validated by Validation procedures. 

Assets are protected by Controls and Countermeasures. 

Security assurance is represented through the relationship chain: 

(:Asset) 
  <-[:THREATENS]-(:Hazard) 
  <-[:MITIGATES]-(:Control) 
  <-[:SATISFIES]-(:Countermeasure) 
  -[:HAS_REQUIREMENT]->(:Requirement) 
  -[:VERIFIED_BY]->(:Verification) 

#### 1.3.1.4 Hazard and Attack 
A (:Hazard) is a system condition or environmental condition that can make compromise of an Asset possible to include the presence of a Threat Actor or conditions within the system such as a Control Action. 

An (:Attack) is an action, technique, tactic, procedure, or exploit path that acts on an Element, Function, Interface or defeats a Countermeasure. 
A Hazard SHALL NOT be treated as the same concept as an Attack.  (:Attack) is a projection of a (:Hazard) into the (:System) through action on an (:Element). (:Interface), or (:Function). 
Hazards MAY reference external threat framework items. 
Attacks MAY reference external attack framework items. 
Hazards and Attacks SHALL be related only when the Attack is a concrete means by which the Hazard may be realized. 

#### 1.3.1.5 Loss 
A (:Loss) represents a specific unacceptable compromise case for: 

- one (:Asset) 
- one Criticality 
- one Assurance 
- one (:Environment) 

Loss SHALL be modeled as an analytical root object. 
The attack tree associated with a Loss SHALL be represented as data on the (:Loss) node or through explicit Loss analysis relationships. 
The Loss node SHALL NOT itself be defined as the DAG. The DAG is the analytical representation of how the Loss may occur. 
Each (:Loss) SHALL have exactly one true Criticality property and exactly one true Assurance property. 

#### 1.3.1.6 Control, Countermeasure, Requirement, Verification 
A (:Control) is an abstract security or assurance objective. 
A (:Countermeasure) is a concrete feature, design element, procedure, or mechanism that satisfies one or more Controls. 
A (:Requirement) is a specification statement that realizes Purpose, Constraint, Countermeasure, Interface, Function, Element, Connection, Capability, or other authorized intent. 
A (:Verification) is a procedure confirming that a Requirement is implemented correctly. 

The canonical traceability direction SHALL be: 

(:Control)-[:ENFORCES]->(:Constraint) 
(:Control)-[:MITIGATES]->(:Hazard) 
(:Countermeasure)-[:SATISFIES]->(:Control) 
(:Countermeasure)-[:HAS_REQUIREMENT]->(:Requirement) 
(:Requirement)-[:VERIFIED_BY]->(:Verification) 

Reverse semantic interpretation SHALL NOT be used to define traceability. 

#### 1.3.1.7 Validation 
(:Validation) is a procedure confirming that the realized System satisfies its intended Purpose in its intended Environment. 
Validation SHALL be related to (:Purpose), not to individual implementation Requirements unless explicitly added in a future version. 

#### 1.3.1.8 Assurance Case / GSN 
(:Goal), (:Strategy), (:Context), (:Justification), (:Assumption), and (:Solution) SHALL represent Goal Structured Notation (GSN) assurance-case artifacts. 
GSN nodes SHALL be used to structure evidence and argumentation about Assets, Loss, Verification, Validation, and other evidence-bearing objects. 
GSN relationships SHALL use Neo4j relationship syntax and UPPERCASE_SNAKE_CASE. 

--- 

### 1.3.2 Canonical Modeling Rules 
All node labels SHALL be singular. 
All relationship names SHALL use UPPERCASE_SNAKE_CASE. 
All relationships SHALL use Neo4j directed relationship syntax: 
(:Source)-[:RELATIONSHIP_NAME]->(:Target) 
Reverse relationships SHALL NOT be explicitly created unless required for performance or explicitly authorized by this SRS. 
All properties with no value SHALL use Null. 
Relationship names SHALL be unique in meaning across the Core Data Model. 
Duplicate logical relationships between the same source node, target node, and relationship type SHALL be prohibited unless multiplicity is explicitly allowed and distinguished by relationship properties. 

Recursive relationships SHALL declare whether they are: 
- acyclic, or 
- cyclic-by-design with bounded traversal 

The Backend SHALL enforce recursive relationship constraints. 
The system SHALL NOT execute unbounded recursive graph traversals. 
All list-returning Backend endpoints SHALL support pagination and maximum result limits. 

--- 

### 1.3.3 Canonical Node Labels 
The Core Data Model SHALL include the following node labels: 

Project / hierarchy: 
- (:Capability) 
- (:Sandbox) 
- (:System) 

System primary nodes: 
- (:Environment) 
- (:Connection) 
- (:Interface) 
- (:Function) 
- (:Element) 
- (:Purpose) 
- (:State) 
- (:ControlStructure) 
- (:Asset) 
- (:Security) 
- (:FunctionalFlow) 

Security / assurance nodes: 
- (:Regime) 
- (:Hazard) 
- (:Attack) 
- (:Control) 
- (:Countermeasure) 
- (:Loss) 

Requirements / verification nodes: 
- (:Constraint) 
- (:Requirement) 
- (:Validation) 
- (:Verification) 

STPA / control-flow role nodes: 
- (:ControlAlgorithm) 
- (:ControlAction) 
- (:ControlledProcess) 
- (:Feedback) 
- (:ProcessModel) 

GSN assurance-case nodes: 
- (:Goal) 
- (:Strategy) 
- (:Context) 
- (:Justification) 
- (:Assumption) 
- (:Solution) 

--- 

### 1.3.4 Canonical Relationship Model 

#### 1.3.4.1 Project and System Hierarchy 

(:Capability)-[:HAS_SYSTEM]->(:System) 
(:Capability)-[:HAS_REQUIREMENT]->(:Requirement) 

(:Sandbox)-[:HAS_SYSTEM]->(:System) 

(:Element)-[:PARENTS]->(:System) 

Constraints: 

- A (:Capability) SHALL be the root of the project hierarchy. 
- A (:Sandbox) SHALL be outside the Capability baseline. 
- A (:System) SHALL NOT be related to both (:Capability) lineage and (:Sandbox) lineage. 
- A (:Element) SHALL parent zero or one child (:System). 
- (:Element)-[:PARENTS]->(:System) SHALL form a Directed Acyclic Graph. 
- Child (:System) HID Index SHALL be derived from the parent (:Element) HID Index and Sequence Number. 

#### 1.3.4.2 System Composition 

(:System)-[:ACTS_IN]->(:Environment) 
(:System)-[:HAS_CONNECTION]->(:Connection) 
(:System)-[:HAS_INTERFACE]->(:Interface) 
(:System)-[:HAS_FUNCTION]->(:Function) 
(:System)-[:HAS_ELEMENT]->(:Element) 
(:System)-[:REALIZES]->(:Purpose) 
(:System)-[:EXHIBITS]->(:State) 
(:System)-[:EXECUTES]->(:ControlStructure) 
(:System)-[:HAS_ASSET]->(:Asset) 
(:System)-[:HAS_SECURITY]->(:Security) 
(:System)-[:HAS_FUNCTIONAL_FLOW]->(:FunctionalFlow) 

Constraints: 

- A (:System) SHALL have exactly one (:Purpose). 
- A (:System) SHALL have at least one (:Environment). 
- A (:System) SHALL have at least one (:State). 
- A (:System) SHALL have exactly one (:Security) unless explicitly extended in a future version. 
- (:Connection) SHALL NOT participate in hierarchy relationships. 
- A (:Connection) SHALL be owned by exactly one (:System) through [:HAS_CONNECTION]. 

#### 1.3.4.3 Environment, State, and Hazard 

(:Environment)-[:HAS_HAZARD]->(:Hazard) 

(:State)-[:TRANSITIONS_TO]->(:State) 
(:State)-[:HAS_HAZARD]->(:Hazard) 
(:State)-[:CONTAINS]->(:Asset) 

Constraints: 

- [:TRANSITIONS_TO] SHALL be the canonical representation of state transition. 
- The Core Data Model SHALL NOT contain a (:Transition) node. 
- [:TRANSITIONS_TO] SHALL be cyclic-by-design. 
- All [:TRANSITIONS_TO] traversal SHALL be bounded. 
- Duplicate logical transitions between the same source and destination (:State) SHALL NOT exist unless distinguished by relationship properties. 
- TransitionKind SHALL be one of: FUNCTIONAL, COUNTERMEASURE_REQUIRED, BOTH. 
- If TransitionKind is COUNTERMEASURE_REQUIRED or BOTH, RequiredByCountermeasureHID and/or RequiredByCountermeasureUUID SHALL identify an existing (:Countermeasure). 
- The referenced (:Countermeasure) SHALL belong to the same SoI unless an explicitly justified cross-SoI analytical exception is recorded. 

#### 1.3.4.4 Connections and Cross-System Interaction 

(:Interface)-[:PARTICIPATES_IN]->(:Connection) 
(:Interface)-[:CONNECTS]->(:Function) 

Constraints: 

- Cross-System interaction SHALL be modeled through (:Connection). 
- Each (:Connection) SHALL relate to two or more (:Interface) nodes. 
- An (:Interface) SHALL NOT participate more than once in the same (:Connection). 
- Connection ownership SHALL NOT imply that all participating Interfaces belong to the owning System. 
- (:Connection) Requirements SHALL belong to the owning SoI. 

#### 1.3.4.5 Functional Flow 

(:Function)-[:FLOWS_TO_FUNCTION]->(:Function) 
(:Function)-[:FLOWS_TO_INTERFACE]->(:Interface) 

(:FunctionalFlow)-[:CONTAINS]->(:Function) 
(:FunctionalFlow)-[:CONTAINS]->(:Interface) 
(:FunctionalFlow)-[:CONTAINS]->(:Connection) 
(:FunctionalFlow)-[:CONTAINS]->(:Element) 
(:FunctionalFlow)-[:CONTAINS]->(:Asset) 

Constraints: 

- [:FLOWS_TO_FUNCTION] SHALL only relate Functions within the same SoI. 
- [:FLOWS_TO_INTERFACE] SHALL only relate a Function and Interface within the same SoI. 
- Functional flow cycles ARE allowed. 
- Functional flow cycles SHALL NOT imply ownership, hierarchy, or requirement parentage. 
- Functional flow traversal SHALL be bounded. 

#### 1.3.4.6 Element, Function, Interface, and Asset Allocation 

(:Function)-[:ALLOCATED_TO]->(:Element) 
(:Interface)-[:ALLOCATED_TO]->(:Element) 

(:Function)-[:CONTAINS]->(:Asset) 
(:Interface)-[:CONTAINS]->(:Asset) 
(:Element)-[:CONTAINS]->(:Asset) 

Constraints: 

- Functions and Interfaces are abstractions. 
- Functions and Interfaces SHALL be allocated to Elements to be realized. 
- Assets MAY be contained by Elements, Functions, Interfaces, or States. 
- Asset containment SHALL NOT imply ownership transfer across SoI boundaries. 

#### 1.3.4.7 Purpose, Constraint, Requirement, Validation 

(:Purpose)-[:HAS_CONSTRAINT]->(:Constraint) 
(:Purpose)-[:HAS_REQUIREMENT]->(:Requirement) 
(:Purpose)-[:HAS_VALIDATION]->(:Validation) 

(:Constraint)-[:HAS_REQUIREMENT]->(:Requirement) 

(:Requirement)-[:PARENTS]->(:Requirement) 
(:Requirement)-[:VERIFIED_BY]->(:Verification) 

Constraints: 

- (:Requirement)-[:PARENTS]->(:Requirement) SHALL form a Directed Acyclic Graph. 
- Requirement parentage MAY cross SoI boundaries only where allowed by Requirements Tool rules. 
- Duplicate parentage edges SHALL NOT exist. 
- Requirements allocated only to (:Purpose) SHALL be treated as unallocated for gap-analysis purposes unless explicitly exempted. 
- The analytical property currently named Barren SHOULD be renamed Barren. 

#### 1.3.4.8 Requirement-Bearing Nodes 

The following nodes SHALL be authorized to own direct [:HAS_REQUIREMENT] relationships: 

(:Capability)-[:HAS_REQUIREMENT]->(:Requirement) 
(:Purpose)-[:HAS_REQUIREMENT]->(:Requirement) 
(:Connection)-[:HAS_REQUIREMENT]->(:Requirement) 
(:Interface)-[:HAS_REQUIREMENT]->(:Requirement) 
(:Function)-[:HAS_REQUIREMENT]->(:Requirement) 
(:Element)-[:HAS_REQUIREMENT]->(:Requirement) 
(:Constraint)-[:HAS_REQUIREMENT]->(:Requirement) 
(:Countermeasure)-[:HAS_REQUIREMENT]->(:Requirement) 

No other node type SHALL create [:HAS_REQUIREMENT] relationships unless explicitly authorized by a later version of this SRS. 

#### 1.3.4.9 Security, Controls, Countermeasures 

(:System)-[:HAS_SECURITY]->(:Security) 

(:Security)-[:HAS_CONTROL]->(:Control) 
(:Security)-[:HAS_COUNTERMEASURE]->(:Countermeasure) 

(:Control)-[:ENFORCES]->(:Constraint) 
(:Control)-[:MITIGATES]->(:Hazard) 

(:Countermeasure)-[:SATISFIES]->(:Control) 
(:Countermeasure)-[:HAS_REQUIREMENT]->(:Requirement) 
(:Countermeasure)-[:APPLIES_TO_FUNCTION]->(:Function) 
(:Countermeasure)-[:APPLIES_TO_INTERFACE]->(:Interface) 
(:Countermeasure)-[:APPLIES_TO_ELEMENT]->(:Element) 
(:Countermeasure)-[:APPLIES_TO_STATE]->(:State) 
(:Countermeasure)-[:APPLIES_TO_FEEDBACK]->(:Feedback) 
(:Countermeasure)-[:BLOCKS]->(:Attack) 

Constraints: 

- (:Control) SHALL represent abstract assurance intent. 
- (:Countermeasure) SHALL represent concrete implementation or design response. 
- Requirements SHALL realize Countermeasures. 
- Verification SHALL verify Requirements. 
- Countermeasure-driven state behavior SHALL be represented by properties on [:TRANSITIONS_TO], not by creating a transition node. 
- [:APPLIES_TO_STATE] identifies affected State nodes but SHALL NOT replace [:TRANSITIONS_TO]. 

#### 1.3.4.10 Hazard and Attack 

(:Hazard)-[:VIOLATES]->(:Constraint) 
(:Hazard)-[:THREATENS]->(:Asset) 
(:Hazard)-[:USES_ATTACK]->(:Attack) 

(:Attack)-[:EXPLOITS]->(:Element) 
(:Attack)-[:DEFEATS]->(:Countermeasure) 

Constraints: 

- Hazard SHALL represent a threatening condition. 
- Attack SHALL represent an action or exploit path. 
- Attack MAY terminate a Loss attack-tree branch as a residual vulnerability. 
- Hazard and Attack SHALL remain separate node types. 

#### 1.3.4.11 Asset, Regime, Loss, and GSN Goal 

(:System)-[:HAS_ASSET]->(:Asset) 

(:Asset)-[:HAS_REGIME]->(:Regime) 
(:Asset)-[:HAS_LOSS]->(:Loss) 
(:Asset)-[:HAS_GOAL]->(:Goal) 

(:Loss)-[:HAS_ENVIRONMENT]->(:Environment) 
(:Loss)-[:HAS_ELEMENT]->(:Element) 
(:Loss)-[:HAS_STATE]->(:State) 
(:Loss)-[:HAS_ATTACK]->(:Attack) 
(:Loss)-[:HAS_COUNTERMEASURE]->(:Countermeasure) 

Constraints: 

- (:Asset)-[:HAS_REGIME]->(:Regime) SHALL replace any use of [:HAS] or [:Has] for Regime. 
- (:Asset)-[:HAS_GOAL]->(:Goal) SHALL replace generic [:HAS] relationships to GSN Goal nodes. 
- Each (:Loss) SHALL be associated with exactly one (:Environment). 
- Each (:Loss) SHALL be associated with exactly one (:Asset) through the inverse path of (:Asset)-[:HAS_LOSS]->(:Loss). 
- Each (:Loss) SHALL have exactly one true Criticality property. 
- Each (:Loss) SHALL have exactly one true Assurance property. 
- Loss analysis DAGs SHALL be represented by Loss relationships and/or AttackTreeJSON, not by redefining Loss as the DAG itself. 

#### 1.3.4.12 STPA Control Structure 

(:ControlStructure)-[:HAS_CONTROL_ALGORITHM]->(:ControlAlgorithm) 
(:ControlStructure)-[:HAS_PROCESS_MODEL]->(:ProcessModel) 
(:ControlStructure)-[:HAS_CONTROLLED_PROCESS]->(:ControlledProcess) 
(:ControlStructure)-[:HAS_CONTROL_ACTION]->(:ControlAction) 
(:ControlStructure)-[:HAS_FEEDBACK]->(:Feedback) 

(:Interface)-[:IMPLEMENTS]->(:ControlAlgorithm) 
(:Interface)-[:IMPLEMENTS]->(:ControlledProcess) 

(:Function)-[:IMPLEMENTS]->(:ControlAlgorithm) 
(:Function)-[:IMPLEMENTS]->(:ControlledProcess) 
(:Function)-[:IMPLEMENTS]->(:ProcessModel) 

(:ControlAlgorithm)-[:GENERATES]->(:ControlAction) 
(:ControlAction)-[:COMMANDS]->(:ControlledProcess) 
(:ControlAction)-[:CAUSES]->(:Hazard) 
(:ControlledProcess)-[:PRODUCES]->(:Feedback) 
(:Feedback)-[:INFORMS]->(:ProcessModel) 
(:ProcessModel)-[:TUNES]->(:ControlAlgorithm) 

Constraints: 

- A single (:Function) SHALL have no more than one [:IMPLEMENTS] relationship into an STPA role node. 
- A single (:Interface) SHALL have no more than one [:IMPLEMENTS] relationship into an STPA role node. 
- Control-loop relationships MAY form bounded analytical cycles. 
- Control-loop traversal SHALL be bounded. 

#### 1.3.4.13 GSN Assurance Case 

(:Goal)-[:SUPPORTED_BY]->(:Goal) 
(:Goal)-[:SUPPORTED_BY]->(:Strategy) 
(:Goal)-[:SUPPORTED_BY]->(:Solution) 

(:Goal)-[:IN_CONTEXT_OF]->(:Context) 
(:Goal)-[:IN_CONTEXT_OF]->(:Justification) 
(:Goal)-[:IN_CONTEXT_OF]->(:Assumption) 

(:Strategy)-[:IN_CONTEXT_OF]->(:Context) 
(:Strategy)-[:IN_CONTEXT_OF]->(:Justification) 
(:Strategy)-[:IN_CONTEXT_OF]->(:Assumption) 

(:Context)-[:HAS_ENVIRONMENT]->(:Environment) 

(:Solution)-[:HAS_VALIDATION]->(:Validation) 
(:Solution)-[:HAS_VERIFICATION]->(:Verification) 
(:Solution)-[:HAS_LOSS]->(:Loss) 

Constraints: 

- GSN relationships SHALL use Neo4j relationship syntax. 
- Generic [:HAS] relationships SHALL NOT be used where a semantic relationship exists. 
- [:SUPPORTED_BY] relationships SHALL form a DAG unless a future version explicitly authorizes cyclic assurance-case structures. 

--- 

### 1.3.5 Canonical Cross-SoI Relationship Rules 

The following relationships MAY cross SoI boundaries when validated by the Backend: 

- (:Interface)-[:PARTICIPATES_IN]->(:Connection) 
- (:Requirement)-[:PARENTS]->(:Requirement) 

The following relationships SHALL NOT cross SoI boundaries unless explicitly justified and recorded as an analytical exception: 

- (:State)-[:TRANSITIONS_TO]->(:State) 
- (:Function)-[:FLOWS_TO_FUNCTION]->(:Function) 
- (:Function)-[:FLOWS_TO_INTERFACE]->(:Interface) 
- (:Countermeasure)-[:APPLIES_TO_STATE]->(:State) 
- (:Countermeasure)-[:APPLIES_TO_FUNCTION]->(:Function) 
- (:Countermeasure)-[:APPLIES_TO_INTERFACE]->(:Interface) 
- (:Countermeasure)-[:APPLIES_TO_ELEMENT]->(:Element) 

The Backend SHALL reject unauthorized cross-SoI relationships. 

--- 

### 1.3.6 Canonical Recursive Relationship Governance 

The following relationships SHALL be acyclic: 

- (:Element)-[:PARENTS]->(:System) 
- (:Requirement)-[:PARENTS]->(:Requirement) 
- (:Goal)-[:SUPPORTED_BY]->(:Goal), unless explicitly extended 

The following relationships SHALL be cyclic-by-design and bounded: 

- (:State)-[:TRANSITIONS_TO]->(:State) 
- (:Function)-[:FLOWS_TO_FUNCTION]->(:Function) 
- Control-loop relationships involving ControlAlgorithm, ControlAction, ControlledProcess, Feedback, and ProcessModel 

The Backend SHALL enforce the declared recursive behavior. 
All recursive traversal SHALL require a maximum depth parameter. 
The Backend SHALL provide safe default maximum depths for all recursive traversals. 


### 1.3.7 System Creation Behavior 

When an (:System) is created from an (:Element) through the relationship: (:Element)-[:PARENTS]->(:System) the following behaviors SHALL occur:
- (:System) is created with a new HID
- One of each  (:Purpose), (:Environment) and (:State) node with default properties are created and related to the new (:System)
- All (:Requirement) nodes related to the parent (:Element) or related to a (:Function) or an (:Interface) related to the parent (:Element) are copied to a new (:Requirement) under the new (:System) (:Purpose) node with the same properties excepting HID and uuid which is modified to reflect the new (:System).
- All (:Asset) nodes related to the parent (:Element) or related to a (:Function) or an (:Interface) related to the parent (:Element) are copied to new (:Asset) nodes with the same properties excepting HID and uuid which is modified to reflect the new (:System)
- New (:Loss) and (:Goal) nodes are created based on the new (:Asset) nodes

------------------


### 1.3.8 Identity Model (HID + UUID)
Each node SHALL contain:
HID (Hierarchical Identifier)
uuid (Globally unique identifier)

HID Format
{TYPE}_{INDEX}_{SEQUENCE}

Example:
SYS_1.2.3_0
UUID Property
uuid: apoc.create.uuid()


 
#### 1.3.8.1 Node Type Identifier
The Node Identifier uniquely identifies each Node Type.  In STPA analysis it is common to identify nodes with a letter and a number.  Each Node type in the SSTPA Tool SHALL have a unique one, two or three character identifier as listed below in the format {Node Type} {Node Type Identifier}:

Capability CAP
Sandbox SB
System SYS
Environment ENV
Connection CNN
Interface INT
Function FUN
Element EL
Purpose PUR
State ST
ControlStructure CS
Asset AST
Security SEC
Constraint CONSTR
Requirement REQ
Validation VAL
Control CTRL
Countermeasure CM
Verification VER
ControlAlgorithm CAL
ProcessModel PM
ControlAction ACT
Feedback FB
ControlledProcess CP
Hazard HAZ
Loss LOS
Attack ATK
Regime REG
Goal G
Strategy SGY
Context CX
Assumption ASS
Justification JUS
Solution SOL


#### 1.3.8.2  Index
The index uniquely identifies the Sub-graph a Node belongs to and is constructed to depict its position in the entire hierarchy.

The Index will be unique for each sub-graph and every node in the sub-graph will have the same Index.
The Index for the Capability SHALL be null as the data set only contains one capability whose only purpose is to attach tier 1 systems.  

When a node is created it SHALL inherit the Index of the sub-graph it belongs to excepting (:System) nodes.

When a (:System) is created as a child of a capability the index SHALL be calculated as the next highest integer value of other System children unless there are no other System children then it gets an index of "1".

When a System is created as the child of an Element its Index SHALL be the index of the Parent Element concatenated with "." concatenated with the (:Element) Node HID Sequence Number property.  For example if an (:Element) has an HID of E_1.2.3_4 than its child (:System) will have an HID of S_1.2.3.4_0.

(:Element) Nodes SHALL have zero or one child (:System) nodes and this constraint will be enforced by Frontend Software.  The Relationship between an (:Element) node and its single (:System) node is:

(:Element)--[:PARENTS]-->(:System)
Note, The (:System) related to here is in a child sub-graph where the new System HID index is set to the concatenation of the (:Element) Index with the (:Element) Sequence Number.

##### 1.3.8.2.1 Index Strategy

The Backend SHALL create the following indexes:

CREATE INDEX node_hid_index IF NOT EXISTS FOR (n) ON (n.HID);
CREATE INDEX node_uuid_index IF NOT EXISTS FOR (n) ON (n.uuid);
CREATE INDEX node_name_index IF NOT EXISTS FOR (n) ON (n.Name);
CREATE INDEX node_type_index IF NOT EXISTS FOR (n) ON (n.TypeName);


#### 1.3.8.3 Sequence Number
The Sequence Number is intended to distinguish nodes of the same type within the same SoI sub-graph.  
The Sequence Number for a System SHALL be "0" because there is only one System in the SoI sub-Graph.
The Sequence Number for a Node other than a System Node SHALL be next highest integer value of other nodes of the same node type in the sub-graph unless there are no others of that node type in the sub-graph, then it is the first and its value is "1". 


1.3.9 Common Property Groups
This section serves to both define the core data model and define how those properties should be displayed.  This is done for consistency across representations.  The progressive disclosure pattern will be:

Node Type-->Node-->Property Groups-->Properties

When Node Type is displayed, under it will be the integer number of nodes of that type.  If other than (0) than when user toggles it will progressively disclose vertically the specific nodes showing HID and Name properties.  below each Node will be the word "Property Groups"  When toggled by the user this progressively discloses all property groups associated with the node starting with the two common property groups.  When a Property Group name is toggled by the user, properties within that group are displayed and can, if allowed be edited. 

For clarity, properties are organized in "Property Groups". in the GUI following the progressive disclosure pattern, when a specific node is progressively revealed there will be a single carrot below it with the word "Properties".  when toggled by the user the GUI will reveal the Property Groups associated with the specific node type.  First displayed will be Property Groups common to all Nodes and Relationships  These common Property Groups are "ID:" and "Description"  These are described below with their specific properties.

The format used to specify Property Groups and Properties will be:
Property Group Name:
Property "Display Name for Property" Data_Type, editability, default: ""

Property Group listings may be followed by a statement addressing constraints on those properties.

	The "Property Group Name" is what is progressively disclosed when the user toggles "Property Groups" and when it is toggled, will progressively disclose its properties.
	"Property" is the property name maintained by the backend on that specific node
	"Display Name" is what the GUI or add-on tool displays to the user as the property name.
	Data_Type is the data type used by the backend and frontend to represent the property.   Note, where possible, Boolean types will be represented by check boxes rather than "True" / "False"
	Editability is direction to the front end to allow the non-privileged user to edit the property.  "edit"=yes, "fixed"=no.  some properties are fixed on creation while others are "fixed" when an analytical report is run.  Some "fixed" properties may be editable only if the current user is "Admin".
	":Default: """ indicates what is between the "" is to be used as the default property value on node creation cast to the correct property type. when default is "Null" value is null, when value is "N/A" the Frontend must enforce a specific property value at time of creation and there is no default value (e.g. property "uuid" has default: "N/A" because a unique identifier is assigned at creation). 

Property Groups are not node properties and only for organizing the display of properties and SHALL be enforced by the Frontend.
Property types SHALL be enforced by both the Frontend and the Backend.
Property defaults and ability to edit SHALL be enforced by the Frontend.
 
Common Property Groups:

ID:
Name "Name:" String edit default: "New"
HID "Hierarchical Identifier: " Structure fixed default: "N/A"
uuid "UUID: " String fixed default: "N/A"
TypeName "Node Type Name " String fixed default: "N/A"
Owner "Data Owner: " String fixed default: "N/A"
OwnerEmail "Owner Email " String fixed default: "N/A"
Creator "Creator: " String fixed default: "N/A"
CreatorEmail "Creator Email: " String fixed default: "N/A"
Created "Created: " datetime fixed default: "N/A" 
LastTouch "Last Touch: " datetime() fixed default: "N/A"
VersionID "Data Schema Version:  " String fixed default: "N/A"

Description:
ShortDescription "Short Description: " String fixed default: "Null"
LongDescription "Full Description:" String fixed default: "Null"

#### 1.3.9.1 Data Ownership Rules 
Every node SHALL have exactly one Owner and one Creator. 
On creation of any node, Creator, CreatorEmail, Owner, and OwnerEmail SHALL be assigned to the current authenticated user. 
Created SHALL be set to the current timestamp on creation. 
LastTouch SHALL be set to the current timestamp on creation and on every committed modification to that node. 
Creator and CreatorEmail SHALL be immutable after node creation except when current user is Admin. 
Owner and OwnerEmail SHALL be editable such that the current user can assume ownership only.  If current user is Admin ownership assignment may be to any registered user. 
Owner and OwnerEmail if changed are always changed as a pair from backend (:User) node properties.

Ownership change SHALL be treated as a node modification for notification purposes. 
Relationship changes involving a node SHALL be treated as changes to that node for notification purposes. 
For a relationship between two existing nodes with different owners, the change SHALL be considered to affect both endpoint nodes. 
If the current user commits a change to a node or relationship and the current user is not the Owner of the affected node, the system SHALL generate a message to the Owner’s mailbox describing the change. 
Message generation SHALL occur within the same transaction as the committed data change. 
Failure to create required ownership-notification messages SHALL cause the overall commit transaction to fail. 


### 1.3.10 Type Unique Property and Relationship Groups
Each Node type will have, in addition, not common properties and relationship groups unique to its type.
Formatting rules from 1.3.7 apply here.
Headings below are Node Type names to which the unique Property Groups and Properties apply.

For node types authorized in Section 1.3.10.4 to assign imported external references by [:REFERENCES], Section 1.3.8 SHALL define property groups used to capture node-local interpretation, applicability, implementation, evidence, and analysis specific to that node. These properties SHALL NOT duplicate or overwrite authoritative imported reference item properties. Imported reference item content remains read-only and authoritative. Node-local external reference properties apply only to the SSTPA node and may differ between nodes referencing the same imported item. 

#### 1.3.10.1 Capability

Mission:
AST "A Capability To:" String edit default: "Null"
BMO "By Means Of:" String edit default: "Null"
IOTCT "In Order To Contribute To:" String edit default: "Null"

#### 1.3.10.2 System

Mission:
AST "A System To:" String edit default: "Null"
BMO "By Means Of:" String edit default: "Null"
IOTCT "In Order To Contribute To:" String edit default: "Null"

#### 1.3.10.3 Environment

Context:
Context "Context" String edit default: "Null"

#### 1.3.10.4  Connection
Reason:
Connection_Description "Rational:" String edit default: "Null"

Properties: 
ConnectionType "Connection Type: " String edit default: "Null"
OSILayer "OSI Layer: " integer edit default: "Null"  
Protocol: "Protocol: "String edit default: "Null" 
Directionality "Directionality: " Enum {Unidirectional, Bidirectional, Multicast, Null} edit default: "Null" 
TimingClass "Timing Class: " String (default "Null") edit default: "Null" 
SecurityClass "Security Classification: " String (default "Null") edit default: "Null" 
PayloadDescription "Payload Description" String (default "Null") 

Criticality:
SafetyCritical "Safety:" Boolean edit default: "False"
SafetyLevel "Level" Integer edit default: "Null"
SafetyDescription: "Description:  " string edit default "Null"

MissionCritical: "Mission:" Boolean edit default: "False"
MissionLevel: "Level" Integer edit default: "Null"
MissionDescription "Description:  " string edit default "Null"

FlightCritical "Flight: " Boolean edit default: "False"
FlightLevel "Level" Integer edit default: "Null"
FlightDescription "Description:  " string edit default "Null"

SecurityCritical "Security:" Boolean edit default: "False"
SecurityLevel "Level" Integer edit default: "Null"
SecurityDescription "Description:  " string edit default "Null"

Assurances:
Confidentiality "Confidentiality" Boolean edit default: "False"
Availability  "Availability" Boolean edit default: "False"
Authenticity  "Authenticity" Boolean edit default: "False"
NonRepudiation "Non-Repudiation" Boolean edit default: "False"
Durability  "Durability" Boolean edit default: "False"
Privacy "Privacy" Boolean edit default: "False"
Trustworthy "Trust" Boolean edit default: "False"


#### 1.3.10.5  Interface
Criticality:
SafetyCritical "Safety:" Boolean edit default: "False"
SafetyLevel "Level" Integer edit default: "Null"
SafetyDescription: "Description:  "string edit default "Null"

MissionCritical: "Mission:" Boolean edit default: "False"
MissionLevel: "Level" Integer edit default: "Null"
MissionDescription "Description:  "string edit default "Null"

FlightCritical "Flight: " Boolean edit default: "False"
FlightLevel "Level" Integer edit default: "Null"
FlightDescription "Description:  " string edit default "Null"

SecurityCritical "Security:" Boolean edit default: "False"
SecurityLevel "Level" Integer edit default: "Null"
SecurityDescription "Description:  "string edit default "Null"

Assurances:
Confidentiality "Confidentiality" Boolean edit default: "False"
Availability  "Availability" Boolean edit default: "False"
Authenticity  "Authenticity" Boolean edit default: "False"
NonRepudiation "Non-Repudiation" Boolean edit default: "False"
Durability  "Durability" Boolean edit default: "False"
Privacy "Privacy" Boolean edit default: "False"
Trustworthy "Trust" Boolean edit default: "False"


#### 1.3.10.5.1  Interface Outgoing Relationship Properties
Only outgoing relationships with properties are identified here.

[:PARTICIPATES_IN] and [:CONNECTS] SHALL have hte following properties:
RelationshipNature "Nature:" Enum {PHYSICAL, LOGICAL, BOTH} edit default: "LOGICAL"  
PhysicalType "Physical Type:" String edit default: "Null"  
Example: universal joint, shaft, hydraulic linkage  
LogicalLayer "OSI Layer:" Enum{N/A, Layer 1: Physical, Layer2: Data Link, Layer 3: Network, Layer 4: Transport, Layer 5 Session, Layer 6: Presentation, Layer 7: Application} edit default: "Null"  
Protocol "Protocol:" String edit default: "Null"  
FlowDirectionality "Directionality:" Enum {Unidirectional, Bidirectional, Multicast} edit default: "Unidirectional"  
TimingClass "Timing Class:" String edit default: "Null"  
SecurityClass "Security Classification:" String edit default: "Null"  



#### 1.3.10.6 Function
Criticality:
SafetyCritical "Safety:" Boolean edit default: "False"
SafetyLevel "Level" Integer edit default: "Null"
SafetyDescription: "Description:  " edit default "Null"

MissionCritical: "Mission:" Boolean edit default: "False"
MissionLevel: "Level" Integer edit default: "Null"
MissionDescription "Description:  " edit default "Null"

FlightCritical "Flight: " Boolean edit default: "False"
FlightLevel "Level" Integer edit default: "Null"
FlightDescription "Description:  " edit default "Null"

SecurityCritical "Security:" Boolean edit default: "False"
SecurityLevel "Level" Integer edit default: "Null"
SecurityDescription "Description:  " edit default "Null"

Assurances:
Confidentiality "Confidentiality" Boolean edit default: "False"
Availability  "Availability" Boolean edit default: "False"
Authenticity  "Authenticity" Boolean edit default: "False"
NonRepudiation "Non-Repudiation" Boolean edit default: "False"
Durability  "Durability" Boolean edit default: "False"
Privacy "Privacy" Boolean edit default: "False"
Trustworthy "Trust" Boolean edit default: "False"

#### 1.3.10.6.1  Function Outgoing Relationship Properties
Only outgoing relationships with properties are identified here.

[:FLOWS_TO_FUNCTION] and [:FLOWS_TO_INTERFACE] SHALLhave the following properties:
RelationshipNature "Nature:" Enum {PHYSICAL, LOGICAL, BOTH} edit default: "LOGICAL"  
PhysicalType "Physical Type:" String edit default: "Null"  
Example: universal joint, shaft, hydraulic linkage  
LogicalLayer "OSI Layer:" Enum{N/A, Layer 1: Physical, Layer2: Data Link, Layer 3: Network, Layer 4: Transport, Layer 5 Session, Layer 6: Presentation, Layer 7: Application} edit default: "Null"  
Protocol "Protocol:" String edit default: "Null"  
FlowDirectionality "Directionality:" Enum {Unidirectional, Bidirectional, Multicast} edit default: "Unidirectional"  
TimingClass "Timing Class:" String edit default: "Null"  
SecurityClass "Security Classification:" String edit default: "Null"  


#### 1.3.10.7 Element

Criticality:
SafetyCritical "Safety:" Boolean edit default: "False"
SafetyLevel "Level" Integer edit default: "Null"
SafetyDescription: "Description:  " string edit default "Null"

MissionCritical: "Mission:" Boolean edit default: "False"
MissionLevel: "Level" Integer edit default: "Null"
MissionDescription "Description:  " string edit default "Null"

FlightCritical "Flight: " Boolean edit default: "False"
FlightLevel "Level" Integer edit default: "Null"
FlightDescription "Description:  " string edit default "Null"

SecurityCritical "Security:" Boolean edit default: "False"
SecurityLevel "Level" Integer edit default: "Null"
SecurityDescription "Description:  " string edit default "Null"

Assurances:
Confidentiality "Confidentiality" Boolean edit default: "False"
Availability  "Availability" Boolean edit default: "False"
Authenticity  "Authenticity" Boolean edit default: "False"
NonRepudiation "Non-Repudiation" Boolean edit default: "False"
Durability  "Durability" Boolean edit default: "False"
Privacy "Privacy" Boolean edit default: "False"
Trustworthy "Trust" Boolean edit default: "False"

Reference Characterization: 
ReferenceApplicabilityStatement "Applicability Statement:" String edit default: "Null" 
ReferenceExposureDescription "Exposure Description:" String edit default: "Null" 
ReferenceAssumption "Assumption:" String edit default: "Null" 

Threat / Property Context: 
ThreatSurface "Threat Surface:" String edit default: "Null" 
TechnologyType "Technology Type:" String edit default: "Null" 
DeploymentContext "Deployment Context:" String edit default: "Null" 



#### 1.3.10.8 Purpose
None

#### 1.3.10.9 State
none

Transitions:
The [:TRANSITIONS_TO] relationship SHALL support the following properties: 
TransitionKind "Transition Kind:" Enum {FUNCTIONAL, COUNTERMEASURE_REQUIRED, BOTH} edit default: "FUNCTIONAL" 
Trigger "Trigger:" String edit default: "Null" 
GuardCondition "Guard Condition:" String edit default: "Null" 
Rationale "Rationale:" String edit default: "Null" 
Where TransitionKind = COUNTERMEASURE_REQUIRED or BOTH, RequiredByCountermeasureHID and/or RequiredByCountermeasureUUID SHALL identify the governing (:Countermeasure). 

Countermeasure Traceability: 
RequiredByCountermeasureHID "Required By Countermeasure HID:" String fixed default: "Null" 
RequiredByCountermeasureUUID "Required By Countermeasure UUID:" String fixed default: "Null" 

Analysis: 
Priority "Priority:" Integer edit default: "Null" 
ResidualRiskNote "Residual Risk Note:" String edit default: "Null" 



#### 1.3.10.10 ControlStructure

Control Structures:
ControlStructureJSON "Diagram Source: "  serialized JSON document fixed default: N/A



#### 1.3.10.11 Asset
Type:
IsPrimary  "Primary:" Boolean edit default: "Null"

Criticality:
SafetyCritical "Safety:" Boolean edit default: "False"
SafetyLevel "Level" Integer edit default: "Null"
SafetyDescription: "Description:  " string edit default "Null"

MissionCritical: "Mission:" Boolean edit default: "False"
MissionLevel: "Level" Integer edit default: "Null"
MissionDescription "Description:  " string edit default "Null"

FlightCritical "Flight: " Boolean edit default: "False"
FlightLevel "Level" Integer edit default: "Null"
FlightDescription "Description:  " string edit default "Null"

SecurityCritical "Security:" Boolean edit default: "False"
SecurityLevel "Level" Integer edit default: "Null"
SecurityDescription "Description:  " string edit default "Null"

Assurances:
Confidentiality "Confidentiality" Boolean edit default: "False"
Availability  "Availability" Boolean edit default: "False"
Authenticity  "Authenticity" Boolean edit default: "False"
NonRepudiation "Non-Repudiation" Boolean edit default: "False"
Durability  "Durability" Boolean edit default: "False"
Privacy "Privacy" Boolean edit default: "False"
Trustworthy "Trust" Boolean edit default: "False"


#### 1.3.10.12 Constraint
Constraint:
CStatement "Constraint Statement:"  string edit default: None

#### 1.3.8.13 Requirement
Requirement:
RStatement: "Text: " String edit default: "Null"
VMethod  "Method: " Enum {Inspection, Demonstration, Analysis, Test, Similarity} edit default: "Null"
VStatement: "Verification Statement: " String edit default: "Null"

Analytical State:
Baseline "Baseline:  " String fixed default: "None"
Orphan   "Orphan" Boolean fixed default: "True"
Barren  "Barren" Boolean fixed default: "True"


#### 1.3.10.14 Validation
Validation:
VStatement  "Validation Statement: " String edit default: "Null"
VMethod  "Method: " Enum {Inspection, Demonstration, Analysis, Test, Similarity} edit default: "Null"

#### 1.3.10.15 Control

Control:
ControlStatement: "Control Statement: " String edit default: "Null"
SatisfactionStatement "Satisfaction Statement: " String edit default: "Null"


NIST SP 800-53r5
ControlID:  "Null"
ControlName:  "Null"
BaseControl:  "Null"
Discussion:  "Null"
RelatedControls:  "Null"
ControlEnhancement:  "Null"
EnhancementDiscussion:  "Null"
EnhancementRelatedControls:  "Null"
References:  "Null"
EvidenceofImplementation:  "Null"


#### 1.3.10.16 Countermeasure
None

#### 1.3.10.17 Verification

Verification:
Procedures: "Null"

#### 1.3.10.18 ControlAlgorithm
cloned from related node

#### 1.3.10.19 ProcessModel
cloned from related node

#### 1.3.10.20 ControlAction
User defined

#### 1.3.10.21 Feedback
User defined

#### 1.3.10.22 ControlledProcess
cloned from related node

#### 1.3.10.23 Hazard
User defined or cloned from reference data

#### 1.3.10.24 Loss

Criticality:
SafetyCritical "Safety:" Boolean edit default: "False"
SafetyLevel "Level" Integer edit default: "Null"
SafetyDescription: "Description:  " string edit default "Null"

MissionCritical: "Mission:" Boolean edit default: "False"
MissionLevel: "Level" Integer edit default: "Null"
MissionDescription "Description:  " edit default "Null"

FlightCritical "Flight: " Boolean edit default: "False"
FlightLevel "Level" Integer edit default: "Null"
FlightDescription "Description:  " string edit default "Null"

SecurityCritical "Security:" Boolean edit default: "False"
SecurityLevel "Level" Integer edit default: "Null"
SecurityDescription "Description:  " string edit default "Null"

Constraint on Criticality properties
The Frontend SHALL enforce the rule that a (:Loss) has a single true criticality.  All others must be "False".

Attack Tree:
AttackTreeFormat — "Format" String, default "SSTPA-ATF-1.0"
AttackTreeVersion — "Version: "Integer fixed default: "-"
AttackTreeCreated — "Created: " datetime fixed default: "Null"
AttackTreeLastModified — "Last Touch: "datetime fixed default: "Null"
AttackTreeCreatedBy — "By: " String fixed default: "Null"
AttackTreeCreatedByEmail "Contact: " String fixed default: "Null"
AttackTreeStatus — Enum {AUTO_GENERATED, ANALYST_REFINED, BASELINED, EXPORTED} edit default: "Null"

Attack Tree JSON:
AttackTreeJSON — "Attack Tree JSON: " Serialized JSON document fixed default: "Null"

NOTE:  Example serialized JSON document for the AttackTreeJSON is in section 3.4.7

Assurances:
Confidentiality "Confidentiality" Boolean edit default: "False"
Availability  "Availability" Boolean edit default: "False"
Authenticity  "Authenticity" Boolean edit default: "False"
NonRepudiation "Non-Repudiation" Boolean edit default: "False"
Durability  "Durability" Boolean edit default: "False"
Privacy "Privacy" Boolean edit default: "False"
Trustworthy "Trust" Boolean edit default: "False"

Constraint on Assurance properties
The Frontend SHALL enforce the rule that a (:Loss) has a single true Assurance.  All others must be "False".


#### 1.3.10.25  Attack
User defined or cloned from reference data


#### 1.3.10.26 FunctionalFlow
Flow Diagram:
FunctionalFlowJSON "Functional Flow Source: " serialized JSON document fixed default: N/A

#### 1.3.10.27 Goal
GSN Info:
GoalID "GSN ID: Integer fixed default: N/A
GoalStatement "Goal Statement" String edit default: "Null"

Diagram Source:
GoalStructure "Goal Structure Source:  serialized JSON document fixed default: "Null"

#### 1.3.10.28 Context
GSN Info:
ContextID "GSN ID: Integer fixed default: N/A
ContextStatement "Context: " String edit default: "Null"

#### 1.3.10.29 Assumption
GSN Info:
AssumptionID "GSN ID: Integer fixed default: N/A
AssumptionStatement "Assumption: " String edit default: "Null"

#### 1.3.10.30 Justification
GSN Info:
JustificationID "GSN ID: Integer fixed default: N/A
JustificationStatement "Justification: " String edit default: "Null"

#### 1.3.10.31 Strategy
GSN Info:
StrategyID "GSN ID: Integer fixed default: N/A
StrategyStatement "Strategy: " String edit default: "Null"

#### 1.3.10.32 Solution
GSN Info:
SolutionID "GSN ID: Integer fixed default: N/A
SolutionStatement "Solution: " String edit default: "Null"

#### 1.3.10.32 Regime
Criticality:
SafetyCritical "Safety:" Boolean edit default: "False"
SafetyLevel "Level" Integer edit default: "Null"
SafetyDescription: "Description:  " edit default "Null"

MissionCritical: "Mission:" Boolean edit default: "False"
MissionLevel: "Level" Integer edit default: "Null"
MissionDescription "Description:  " edit default "Null"

FlightCritical "Flight: " Boolean edit default: "False"
FlightLevel "Level" Integer edit default: "Null"
FlightDescription "Description: " string edit default "Null"

SecurityCritical "Security:" Boolean edit default: "False"
SecurityLevel "Level" Integer edit default: "Null"
SecurityDescription "Description:  " string edit default "Null"

Contact Info:
AuthorityName "Name: " string edit default "Null"
AuthorityTitle "Title: " string edit default "Null"
AuthorityOrg "Organization: " string edit default "Null"
AuthorityEmail  "Email: " string edit default "Null"
SupInfo "Supplumental: " string edit default "Null"

Authoritative Documentation:
DomainGuidance "Guidance: " string edit default "Null"


----------------------------------------
## 1.4 Tool Data Model
SSTPA Tools is intended to support a large dispersed engineering team. In this initial version the Backend and Frontend will be bundled for use on a single system.  Data in the Core Data Model will have ownership and also record the creator.  Owner and creator contact information will be captured as email address.  SSTPA Tools will have an add-on tool called "Message Center" which allows users to message each other and allows the system to notify Users of changes to properties and relationships on data they own.


Tool Data Model consists of utility information for SSTPA Tools maintained by the backend.  It consists of:
Data on the SSTPA Tool: (:SSTPA_Tool)
Tool Data will parent a single node SSTPA Tools Data which SHALL have properties typical of a commercial software installation to include the license.
Data on Users: (:User)
Users can own data and edit properties identified with "edit".
Data on Admins (:Admin)
Admins cannot own data but can edit and commit certain specific properties identified as for Admins Only otherwise their Commit is invalid.  Admin Users can access the backend Webserver to view telemetry information.
Data on the Host system (:Host)
Data on Messages for Users: (:Message) 
Data on Messages for Admins: (:Message) 

### 1.4.1  Onboarding
The SSTPA Tool Installer SHALL be configured to capture the Installer's Name and e-mail contact information.  It SHALL establish an Admin account and a User account for the Installer.
When invoking the SSTPA Tool, it SHALL present a login window asking for User Name and Password and there will be a link below with the label "New User" and a link labeled "Admin".  If the User enters a valid Username and password, the GUI will start with the identified (:User) as the active User for the purposes of the ownership and permissions model.  If  the User clicks the "New User" Link, they will be presented with a screen where they can create a new (:User) node by editing its properties. then they will be redirected to the Login screen.  If the User clicks on "Admin" link than a new Admin login window appears with the same behavior as the User login window.  If a valid Admin properly logs in, then the active User for the purposes of ownership and permissions model is the (:Admin). Under User Name and Password is a link with label "I want to be an Admin".  If the User clicks this, a new screen is presented which has the Admin User Name and password fields and a New Admin and New Admin Password field.  A valid existing Admin must "login" to authorize the creation of a new (:Admin) who sets their User Name and Password and other properties.  on success the User is sent to hte Admin login screen to login as a new Admin.


### 1.4.2 Admin Data
On installation of SSTPA Tools it will have , the "Installer" SHALL be required to provide an Admin email and establish a user account.  Admin cannot own data but can edit data and certain elements of data are fixed for normal users can be edited by Admin.  These properties are explicitly identified.
All users will "login" to the system as either an existing User or a new User.  If new User they will setup a User account.  Per the initial security model, no password or authentication will be used.

Tool Data will parent a single node called Admins and nodes associated with each registered Admin will be a child of the Admins node with properties describing the Admin

### 1.4.3 User Data
Tool Data will parent a single node called Users and nodes associated with each registered user will be a child of the Users node with properties describing the User


### 1.4.4 Messaging Data Model 
Nodes 
(:User) 
UserName 
UserEmail 
UserHash or equivalent identifier 

(:Mailbox) 
MailboxID 
Owner 
OwnerEmail 
UnreadCount 
Created 
LastTouch 

(:Message) 
MessageID / uuid 
Subject 
Body 
MessageType enum {DIRECT, CHANGE_NOTIFICATION, SYSTEM} 
SentAt 
ReadAt 
DeletedAt 
Sender 
SenderEmail 
Recipient 
RecipientEmail 
RelatedNodeHIDs 
RelatedRelationshipTypes 
CommitID 
IsRead 
IsDeleted 
RequiresApproval (default False) 
ApprovalStatus enum {NOT_APPLICABLE, PENDING, APPROVED, REJECTED} 


Relationships 

(:User)-[:OWNS_MAILBOX]->(:Mailbox) 

(:Mailbox)-[:HAS_MESSAGE]->(:Message) 

(:Message)-[:RELATES_TO]->(:System|:Environment|:Interface|:Function|:Element|:Purpose|:State|:ControlStructure|:Asset|:Control|:Constraint|:Requirement|:Validation|:Countermeasure|:Verification|:ControlAlgorithm|:ProcessModel|:ControlAction|:Feedback|:ControlledProcess|:Hazard|:Loss|:Attack) 

(:Message)-[:REPLY_TO]->(:Message) for replies 



### 1.3.10 External Reference Framework Data Model 

The SSTPA Tool SHALL support imported reference frameworks in the Backend for read-only use by the GUI and Add-on Tools. These frameworks will include: 
•	NIST SP 800-53r5
•	MITRE ATT&CK
•	MITRE EMB3D

These frameworks will be manually inserted into the database prior to delivery.  Future versions of SSTPA Tools will allow update of this data which SHALL be fixed (not editable by any user).  This restriction enforces the data license agreement.  Future external reference update capability will update from a provided BLOB image matching the data structure in the database. SSTPA Tools will not scrape the internet for data.  

The purpose of these imported reference frameworks is to allow the User to: 
•	Navigate authoritative external reference data
•	Read the properties of imported reference items
•	clone external reference properties and tag the Core data Node to the Reference Data node

Imported reference framework data SHALL be stored in Neo4j as a distinct but connected graph structure separate from the SSTPA System of Interest (SoI) sub-graphs. This is consistent with the SRS requirement that the Backend import NIST SP 800-53, MITRE ATT&CK, and MITRE ESTM data into graphical form and that the GUI provide add-on tools for populating selected node properties from those datasets. 


## 1.5 Reference Data Model

### 1.5.1 External Reference Framework Nodes 

### 1.5.1.1 MITRE ATT&CK Framework Nodes
The SSTPA Tools Development System SHALL create a data pipeline to transform MITRE ATT&CK data from a raw preserved state through an intermediate format into a graphical representation loaded into the Backend.
Note:
MITRE ATT&CK Framework version 19 releases on April 28, 2026.
MITRE ATT&CK framework uses the stix2 data format
cti-python-stix2 is a  repository that provides Python APIs for serializing and de-serializing STIX2 JSON content, along with higher-level APIs for common tasks, including data markings, versioning, and for resolving STIX IDs across multiple data sources.
https://github.com/oasis-open/cti-python-stix2.git

ATT&CK STIX data URL:  https://github.com/mitre-attack/attack-stix-data.git
Note:  ATT&CK data is in four sub-directories; Enterprise-attack, ics-attack, and mobile-attack.  Within those directories, the files are organized by version number. version should be 19.


Backend Representation:
(:ATT&CK_19)-[HAS]->
(:CyberCampaign)
(:ThreatActors)
(:Tactics)
(:Techniques)
(:Tooling)
(:SoftwareElements)
(:DataElements)
(:DetectionStrategies)
(:Analytics)
(:Mitigation)


Users SHALL be able to create Reference nodes of the above type and of type (not included in MITRE ATT&CK) (:AK_Procedure) 
(:AK_Procedures) are very system specific so MITRE does not support this, but Users are working on specific systems so will develop (:Attack) nodes and using tools will persist them across the heirarchy as (:AK_Procedures).

Backend Internal Relationships:

(:CyberCampaign)-[:USES]->(:ThreatActors)-[:USES]->(:Tactics)-[:ACHIEVED_THROUGH]->(:Techniques)-[:USES]->(:AK_Procedures)
(:ThreatActors)-[:USES]->(:Tooling)
(:Tooling)-[:CONTAINS]->(:SoftwareElements)-[CREATES]->(:DataElements)
(:Procedures)-[:USE]->(:Tooling)
(:Tooling)-[:FINGERPRINTS]->(:ThreatActors)
(:SoftwareElements)-[:FINGERPRINTS]->(:ThreatActors)
(:DataElements)-[:FINGERPRINTS]->(:ThreatActors)
(:Detection_Strategies)-[:INFORM]->(:Analytics)-[:RECOMMEND]->(:Mitigation)


### 1.5.1.2 MITRE EMB3D Framework Nodes
The SSTPA Tools Development System SHALL create a data pipeline to transform MITRE EMB3D data from a raw preserved state through an intermediate format into a graphical representation loaded into the Backend.

EMB3D Data URL is: https://github.com/mitre/emb3d/blob/main/assets/emb3d-stix-2.0.1.json
The EMB3D data is generated using the OASIS CTI TC’s Python STIX2 library in UTF-8 text encoding

Note:  EMB3D uses type properties in its embd-stix-2.0.1.json file.  Reference data Model transforms the stix-2.0.1 file into graphical format by creating new Node types based on the "type' value name.  The mapping below shows the Node definition followed by the canonical stix-2.0.1 key value pair.

(:Vulnerability) "type": "vulnerability"
(:Course_Of_Action) "type": "course-of-action"
(:Device) "type": "x-mitre-emb3d-property" 
[:MITIGATES] "type": "relationship" AND  "relationship_type": "mitigates"
[:RELATES_TO] "type": "relationship" AND "relationship_type": "relates-to"

Backend Representation:

(:EMB3D_21)-[HAS]->
(:Vulnerability)
(:Course_Of_Action)
(:Device)

Backend Internal Relationships:

(:Course_Of_Action)-[:MITIGATES]->(:Vulnerability)
(:Vulnerability)-[:RELATES_TO]->(:Device)



### 1.5.1.3 NIST SP800-53 Catalog of Controls Nodes
NIST is maintaining OSCAL content for multiple revisions of the NIST Special Publication (SP) 800-53. The XML, JSON, and YAML versions of SP800-53 given here are derived from the NIST publications. The OSCAL XML, JSON, and YAML variants are all equivalent in their information content and are provided to support tooling on different format-specific implementation stacks. These OSCAL files are intended to faithfully represent the control-related content from the published documents in machine-readable formats.
Data URL:  nist.gov/SP800-53/rev5/json  


(:NIST800-53R5)-[HAS]->
(:NIST_Control)


### 1.5.2 Reference Framework Identity 

Each imported Reference Framework node ((:NIST800-53R5),(:EMB3D_21),(:ATT&CK_19)) SHALL contain the following properties: 
•	FrameworkName
•	FrameworkVersion
•	ExternalID
•	ExternalType
•	Name
•	ShortDescription
•	LongDescription
•	SourceURI
•	Imported
•	LastUpdated
•	RawData

ExternalID SHALL be unique within a given framework version. 

The Backend SHALL preserve the source framework identifier for each imported item. 


### 1.5.3 Framework Structure Relationships 

The Backend SHALL support framework hierarchy and navigation relationships for imported reference data. 


### 1.5.4 SSTPA Core Data Node External Reference Relationships 
Core Data Nodes SHALL be able to clone Reference Node properties directly which will be filtered of reference specific properties (e.g. "spec_version": "2.1").


(:Countermeasure) Nodes in the Core Data Set SHALL be able to clone the properties of the following nodes:
(:DetectionStrategies)
(:Analytics)
(:Mitigation)
(:Course_Of_Action)


(:Attack) Nodes in the Core Data Set SHALL be able to clone the properties of the following nodes:
(:Tactics)
(:Techniques)
(:Vulnerability)

(:Element) Nodes 
in the Core Data Set SHALL be able to clone the properties of the following nodes:
(:SoftwareElements)
(:DataElements)

(:Hazard) Nodes in the Core Data Set SHALL be able to clone the properties of the following nodes:
(:Vulnerability)
(:ThreatActors)

(:Control) Nodes in the Core Data Set SHALL be able to clone the properties of the following nodes:
(:NIST_Control)


### 1.5.6 Read-Only Constraint 

Imported reference framework nodes SHALL be read-only to the User. 
The User SHALL NOT edit imported framework nodes from the GUI. 
The only permitted mutation involving imported reference framework data from the GUI SHALL be creation or removal of a [:REFERENCES] relationship between an SSTPA node and a valid imported reference item. 


## 1.6 Help Data Model
The Help Data Model needs Help

The Help Data Model consists of:
Help information on GUI fields and input boxes
Tutorial Information
Definitions and description of SSTPA Termonology


# 2 Structure

Primary elements of the SSTPA Tool will be:

1. The Startup Software
2. The back-end graph database and supporting application (called the "BackEnd")
3. The User facing components to include GUI interface and Add-on Tools (called the (Front-End)
4. The Installer


## 2.1  Startup Software

The Startup Software for this version of SSTPA Tools is a stand in for a security application which authorizes users prior to allowing connection to the Frontend or the Backend.  In this version, all users are authorized and the primary use case is connection from frontend to backend on the same physical machine.  

The User / Analyst launches SSTPA Tools from Startup Software, which collects user information, activates the database and launches the GUI.  When the user exits, the Startup Software will assure the Database connections are properly closed.

Startup Software SHALL launch from a desktop icon or Command line.

The Startup Software SHALL present the user with a dialog with default SSTPA theme and startup animation.

The Startup Software SHALL allow the user to change Theme and if a new theme is chosen SHALL change to that theme and play its startup animation. 

The Startup Software SHALL connect to a Backend or start the Backend on the Local Machine.  In future versions, Startup Software will connect to remote Backends.

Startup Software SHALL present the User with the Backends' list of Users and Roles and allow the User to select or Add a User.  I know this is bad security practice, but this is a stand in for a more comprehensive security implementation in a future version.   and e-mail contact information.

Startup Software SHALL launch the Frontend Software after User selects or adds a User ID enters information 

On receiving the Shutdown command from the Frontend software, Startup Software SHALL assure both frontend and backend are properly shutdown preserving stored data (i.e. don't kill the database while transactions are in process). 

## 2.2 BackEnd

The backend database will include the graph database and support software needed for ACID compliance.  It will use the most current stable NEO4J Community Edition with a defined pathway to the Enterprise Edition on customer desire.  Backend will be divided into docker containers and Docker Compose.  User will connect and interact with a reverse proxy which will connect to the database.  The reverse proxy will collect and present telemetry on backend performance.

The Backend SHALL be configured to execute CYPHER\_25 scripts.

The back-end SHALL support multiple concurrent connections.

The Backend layout is shown below.


Internet / Remote Clients 
                              | 
                              | HTTPS :443 
                              v 
                    +----------------------+ 
                    |   Caddy   | 
                    | TLS + reverse proxy | 
                    +----------+----------+ 
                               | 
                               | HTTP :8080 
                               | (internal only) 
                               v 
                    +----------------------+ 
                    |      Go Backend      | 
                    | chi + Neo4j driver   | 
                    | Prom/OTel instr.     | 
                    +----+-----------+-----+ 
                         |           | 
         Bolt :7687      |           | OTLP :4317 / :4318 
      (internal only)    |           | 
                         v           v 
                +----------------+  +----------------------+ 
                |     Neo4j      |  | OTel Collector       | 
                | Community Ed.  |  | traces/metrics pipe  | 
                +----------------+  +----------+-----------+ 
                                                | 
                                                | scrape/export 
                                                | 
                                  +-------------+-------------+ 
                                  |                           | 
                                  v                           v 
                         +----------------+          +----------------+ 
                         |   Prometheus   |          | Tempo   | 
                         | metrics store  |          | trace store    | 
                         +-------+--------+          +--------+-------+ 
                                 |                            | 
                                 +-------------+--------------+ 
                                               | 
                                               v 
                                       +---------------+ 
                                       |    Grafana    | 
                                       | dashboards    | 
                                       +---------------+ 

The core idea is to have two network zones: 

1. Public edge network 

Only the reverse proxy is exposed here.  caddy accepts HTTPS from remote clients forwards requests to the Go backend.  Grafana is exposed for remote dashboard access during development 


2. Private backend network 

Everything else talks here. 

Go backend 
Neo4j 
OpenTelemetry Collector 
Prometheus 
Grafana 
Tempo

only the reverse proxy SHALL be internet-facing. 

### 2.2.1 User Facing Container
Backend should put Reverse Proxy and the Grafana in the same container with sufficient software / configuration to present telemetry to external user.

### 2.2.2  Reverse Proxy

Reverse proxy: Caddy

Responsibilities: 

terminate TLS 
expose port 443 
redirect 80 -> 443 
proxy /api/* to the Go backend 
expose Grafana during development 


Typical traffic: 

Client -> Caddy -> Go backend 


--- 
## 2.2.3 Database Container
Backend Should put non-user interacting applications and the database into the a single container.

### 2.2.4 Backend Software
Backend Software SHALL be written in the most current stable version of the Go language.


Go Software Responsibilities: 

expose REST API 
handle auth, validation, routing 
start Neo4j transactions 
expose /metrics for Prometheus 
emit OpenTelemetry traces/metrics/logs 
Typical internal connections: 
to Neo4j on neo4j:7687 
to OTel collector on otel-collector:4317 or 4318 

### 2.2.5 Backend Database

The backend Database SHALL be the latest stable Neo4j Community Edition.

Responsibilities: 
persist graph data 
provide ACID transactions 
accept Bolt protocol connections from backend 


Typical internal connection: 

Go backend -> neo4j:7687 
Do not expose Neo4j publicly. 


### 2.2.6 Telemetry

The backend SHALL use Open Telemetry Collector for backend telemetry

Responsibilities: 

receive telemetry from backend 
batch/process telemetry 
export traces to Tempo 
optionally expose Prometheus-scrapable metrics or forward OTLP metrics 


Typical flow: 

Go backend -> OTel Collector -> Tempo/Jaeger 


### 2.2.7 Metrics

The backend SHALL use Prometheus for metrics.

Responsibilities: 
scrape /metrics endpoints 
store time-series metrics 
answer PromQL queries from Grafana 


Typical scrape targets: 

backend:8080/metrics 
otel-collector metrics endpoint 

optionally Neo4j exporter if you add one 


### 2.2.8 Traces

The backend SHALL use Tempo.

Responsibilities: 
store distributed traces 
let Grafana drill into request traces 


Typical flow: 

OTel Collector -> Tempo/Jaeger 
Grafana -> Tempo/Jaeger 


### 2.2.9 Dashboard

The backend SHALL use Grafana. 

Responsibilities: 
display metrics dashboards 
display traces 
correlate slow requests with backend metrics 


Typical data sources: 
Prometheus 
Tempo

Grafina SHALL present an accessible dashboard via the reverse proxy

### 2.2.10 Backend API Requirements 

The Backend SHALL expose a REST API to support all Frontend and tool interactions. 

---------------------------------------- 

#### 2.2.10.1 General Requirements 
•	The API SHALL use HTTPS
•	The API SHALL return JSON
•	All endpoints SHALL support concurrent access
•	All write operations SHALL be transactional
---------------------------------------- 

#### 2.2.10.2 Node Retrieval 

The Backend SHALL provide endpoints for node lookup: 
•	Retrieve node by HID
•	Retrieve node by uuid
•	Retrieve node by type

Responses SHALL include: 
•	All node properties
•	HID
•	uuid
•	TypeName
•	Containing SoI
---------------------------------------- 

#### 2.2.10.3 Hierarchy Retrieval 

The Backend SHALL provide: 
•	Full system hierarchy (Capability → Systems)
•	Parent-child relationships between Systems

This endpoint SHALL: 
•	Support efficient graph rendering
•	Minimize payload size
---------------------------------------- 

#### 2.2.10.4 Search 

The Backend SHALL support search queries across nodes. 

Search SHALL support: 
•	HID (exact)
•	uuid (exact)
•	Name (partial)
•	ShortDescription (partial)
•	Node Type filtering

Search results SHALL include: 
•	Node metadata
•	Containing SoI
•	Node type
---------------------------------------- 

#### 2.2.10.5 Relationship Validation 

The Backend SHALL validate relationships before creation. 

Validation SHALL: 
•	Confirm allowed node types
•	Enforce relationship rules
•	Prevent invalid associations

The API SHALL return: 
•	Valid / invalid
•	Reason for invalidity
---------------------------------------- 

#### 2.2.10.6 Context Retrieval 

The Backend SHALL provide context for any node: 
•	Containing (:System)
•	Path within hierarchy
•	Parent relationships
---------------------------------------- 

#### 2.2.10.7 Performance Requirements 

The Backend SHALL: 
•	Use indexes on:

•	HID
•	uuid
•	Name
•	TypeName

•	Provide optimized queries for:

•	hierarchy traversal
•	search operations
---------------------------------------- 

#### 2.2.10.8 Transaction Requirements 
•	All mutations SHALL be ACID compliant
•	Relationship creation SHALL be atomic
•	Validation SHALL occur prior to commit
---------------------------------------- 

##### 2.2.10.8.1 Ownership and Change Notification Requirements 

The Backend SHALL determine the set of affected nodes for each commit. 
Affected nodes SHALL include: 
any node whose property values changed 
both endpoint nodes of any created relationship 
both endpoint nodes of any removed relationship 

For each affected node, the Backend SHALL compare the current user with the node Owner. 
If the current user is not the Owner, the Backend SHALL create a CHANGE_NOTIFICATION message addressed to that Owner. 
A single commit MAY generate multiple messages. 
The Backend SHOULD aggregate multiple changes for the same Owner within one commit into a single message. 
Each change-notification message SHALL include: 
Subject 
Sent timestamp 
Sender 
Recipient 
affected HID or HIDs 
change type summary 
old owner and current owner where ownership changed 
commit identifier 

Message creation SHALL occur in the same ACID transaction as the graph mutation. 

If the data mutation succeeds, required messages SHALL also succeed. 

If required messages fail, the entire transaction SHALL roll back. 

The current version SHALL notify only through internal mailbox messaging. 

Future attachment to organization email exchange SHALL be supported through an integration boundary and SHALL NOT change the internal mailbox requirement. 

-----------------------------

##### 2.2.10.8.2 Ownership Change Rules 

Any user MAY change ownership in the current version. 

Ownership change SHALL itself generate a notification to the prior Owner when performed by a different user. 

Ownership change SHALL update Owner, OwnerEmail, and LastTouch. 

Ownership change SHALL NOT modify Creator or CreatorEmail. 

-----------------------------


#### 2.2.10.9 Security (Placeholder) 

The Backend SHALL support: 
•	User identification
•	Role-based access (future implementation)

2.2.10.10 External Reference Framework API Requirements 

The Backend SHALL expose REST API endpoints to support import, retrieval, search, navigation, inspection, and assignment of external reference framework data. 

All write operations SHALL be transactional, consistent with existing Backend API requirements. 

2.2.10.10.1 Messaging API Requirements 

Include endpoints such as: 

GET /api/messages 
GET /api/messages/{messageId} 
POST /api/messages 
POST /api/messages/{messageId}/reply 
POST /api/messages/{messageId}/read 
DELETE /api/messages/{messageId} 
GET /api/messages/unread-count 

Response requirements: 

list view returns subject, datetime, HID summary, sender, message type, read/unread 
detail view returns full body plus related HIDs and reply chain 
list query SHALL support sort by subject, datetime, HID, and sender 
list query SHALL support ascending and descending order 
------------------------------------------


2.2.10.10.1 Framework Import Requirements 

The SSTPA Development System SHALL provide import tools for: 
•	NIST SP 800-53r5
•	MITRE ATT&CK
•	MITRE EMB3D

The import process SHALL: 
•	Convert source data into graph format
•	Preserve framework version information
•	Preserve source identifiers
•	Preserve hierarchy and related-item relationships where supported by the source data
•	Avoid creating duplicate imported reference items for the same framework version and source identifier

The Backend SHALL support converted data set. 


## 2.3 Docker networks

Backend should use at least two Docker networks, an edge network and a backend network. 

### 2.3.1 Edge Network

For public-facing traffic. 

Members: 

caddy
grafana (proxy publishes dashboards) 


### 2.3.2 backend Network

For internal service-to-service traffic. 

Members: 

caddy
backend 
neo4j 
otel-collector 
prometheus 
grafana 
tempo


So the proxy sits on both networks: 
on edge to accept external traffic on backend to forward internally Everything else sits only on backend. 


### 2.3.3 Security model

The backend architecture should give the following defaults: 

Publicly exposed 
	443 on reverse proxy 
	maybe 80 for redirect to 443 


Internal only 
	backend 8080 
	Neo4j 7687 
	Prometheus 9090 
	Grafana 3000 unless intentionally proxied 
	Tempo/Jaeger ports 
	OTel Collector ports 


###  2.3.4 Docker Compose Topology
Docker Compose should have a Docker Compose-style topology aligned to below:

services: 
  caddy: 
    image: caddy:latest 
    ports: 
      - "80:80" 
      - "443:443" 
    networks: 
      - edge 
      - backend 
    depends_on: 
      - backend 

  backend: 
    image: my-backend:latest 
    networks: 
      - backend 
    depends_on: 
      - neo4j 
      - otel-collector 

  neo4j: 
    image: neo4j:community 
    networks: 
      - backend 
    volumes: 
      - neo4j-data:/data 

  otel-collector: 
    image: otel/opentelemetry-collector:latest 
    networks: 
      - backend 

  prometheus: 
    image: prom/prometheus:latest 
    networks: 
      - backend 

  tempo: 
    image: grafana/tempo:latest 
    networks: 
      - backend 

  grafana: 
    image: grafana/grafana:latest 
    networks: 
      - backend 

networks: 
  edge: 
  backend: 

volumes: 
  neo4j-data: 

Backend SHALL allow display and configuration of ports, configs, and volumes
Backend SHALL send configuration information to Frontend or Startup on connection 


###  2.3.5   Backend Telemetry
Request flow 

1. Client sends HTTPS request to api.example.com 
2. Caddy terminates TLS 
3. Proxy forwards to backend:8080 
4. Backend authenticates request 
5. Backend starts trace/span 
6. Backend opens Neo4j transaction 
7. Neo4j commits/rolls back 
8. Backend returns JSON response 
9. Proxy returns HTTPS response 

Telemetry flow 

1. Backend records request counter + latency histogram 
2. Backend creates OpenTelemetry spans 
3. Prometheus scrapes /metrics 
4. OTel Collector receives spans 
5. Collector exports traces to Tempo 
6. Grafana queries Prometheus and Tempo 
7. Operator sees metrics + traces together 


# 3 Frontend
The Front end includes the Graphic User Interface (GUI) and Add-on Tools.  It also connects to the Backend which serves as its datastore.

The GUI SHALL utilize Backend API endpoints defined in Section 2.2.10.

## 3.1  GUI Overview

The core of the Frontend is the "GUI".  As the principal user facing component it will be branded as "SSTPA Tools". 
The GUI SHALL execute as a stand alone desktop application.
The GUI SHALL operate in a single window with the capability to execute add-on tools in pop-up windows
The GUI SHALL connect to the Backend and commit data only after commit confirmation in a confirmation dialog.

## 3.2  GUI Style
The GUI SHALL have a style defined by a Style sheet (.css file) with a default style and user selectable alternate styles.
The GUI SHALL be organized in panels.

### 3.2.1  Default_Style.css file

The Default_Style.css is the single default User Interface (UI) style.  

Default\_Style.css SHALL implement a "high-tech, dark, liquid glass aesthetic to include presentation of data in a data grid pattern on tilt cards with background animation simulating slow dynamic motion and details for user edit of specific records presented using the right-side drawer interaction pattern.

The visual style foundation should use the Liquid Glass (or "frosted translucent") pattern implemented through JavaScript libraries for data grid, animations and UI components and HTML5/CSS to enable the effect.


#### 3.2.1.1 Liquid Glass Core Techniques

Default\_Style.css should use the "liquid glass" or "frosted translucent" aesthetic using a CSS backdrop-filter to provide blur, contrast, brightness effects on elements behind a translucent surface.

Example code:

.glass {
&nbsp; background: rgba(30, 30, 30, 0.35);
&nbsp; backdrop-filter: blur(12px) saturate(150%) contrast(120%);
&nbsp; border: 1px solid rgba(255, 255, 255, 0.1);
}


Default\_Style.css should use a CSS filter applied directly to elements to provide subtle glow effects.

Example code:
filter: drop-shadow(0 0 6px rgba(0,255,255,0.5));
Default\_Style.css should use gradients and subtle noise textures to avoid the flat plastic look.
Example code:
background: linear-gradient(135deg, rgba(50,50,60,0.4), rgba(20,20,25,0.4));

Default\_Style.css should allow dynamic theme tweaking.

Example code:

:root {
&nbsp; --glow-color: #00ffff;
&nbsp; --glass-alpha: 0.35;
}

Default\_Style.css should use CSS masks and clipping for curved, liquid-like shapes.


## 3.3  SSTPA Tools branding Panel
	The top Panel of the SSSTPA Tools GUI SHALL show SSTPA Logo on top left with "SSTPA Tools" name and version at center.  right SHALL contain status information from Backend in smaller font to include connection IP, port and connection status (in Courier font and a contrasting color). In the right it SHALL also show the current User name, mail Icon titled "Message Center" (an add-on Tool) and a gear Icon for changing GUI parameters such as changing Style and displaying license, and system version information.   

3.3.1 Message Center Add-on Tool 
The Branding Panel SHALL display a mail icon labeled or tooltiped as Message Center. 

The Message Center SHALL display an unread indicator when unread messages exist. 

Selecting the Message Center icon SHALL open a pop-up window. 

The Message Center pop-up SHALL display the current user’s mailbox only. 

The Message Center pop-up SHALL not change the current SoI. 

The Message Center SHALL be a specialized add-on tool consistent with the pop-up model already allowed by the SRS. 



## 3.4  SSTPA Control Panel
The SSTPA Control Panel SHALL be below the SSTPA Branding Panel and contain icons for Add-on Tools starting from the Left with the "Navigator Tool" going in sequence with Add-on Tools ending with the "Reports" drop-down Menue. 

The SSTPA Control Panel SHALL present an ICON for 'Shutdown" at the far right of the panel as a typical power icon but in red color. 

If User selects a menu item or icon which does not have real functionality attached, it SHALL present an alert dialog titled "Under Construction" with a construction icon and an "OK" button.  On click of the "OK" button, the alert will close. 

### 3.4.1 The "Navigator Tool" 
The "Navigator Tool" will perform the following core functions:
It is the means by which he SSTPA Tools GUI User Selects a System of Interest (SoI) for the rest of the GUI
It allows the User to search the entire Hierarchy to identify specific nodes and display it graphicly
It allows the User to explore the System Hierarchy while not changing the current SoI
It allows the Use to navigate to and select a node from another System to clone into the current SoI
It allows the User to select a (:Connection) node owned by another System and connect an SoI (:Interface) to it
User can set Navigator Tool to display participants and owner in a specific (:Connection) and graphically visualize selected connections up to all of them.


The tool described here SHALL be branded at top of window as "Navigator Tool"


The Navigator Tool SHALL: 
•	place (:Capability) at the top or central anchor position
•	display connected (:System) nodes using force-directed or constrained hierarchical layout behavior
•	support smooth zoom, pan, and animated re-centering
•	preserve visible spatial continuity during user interaction
•	prevent navigation states where all nodes are moved completely out of view

---------------------------------------- 

3.4.1.1 Modes of Operation 
The Hierarchy Search Tool SHALL support four modes: 

a. SoI Selection Mode 
•	Allows selection of a (:System) node as the current SoI
•	On confirmation, the selected node SHALL become the current SoI
•	All GUI panels SHALL update to reflect the new SoI

b. Association Selection Mode 
•	Allows selection of nodes outside the current SoI for association
•	SHALL NOT change the current SoI
•	SHALL return the selected node to the calling context

Use cases include: 
•	(:Interface)-[:PARTICIPATES_IN]->(:Connection) association where (:Connection) belongs to another SoI.
•	(:Requirement)-[:PARENTS]->(:Requirement) association
•	Other cross-SoI associations as defined by relationship rules

c. Search / Locate Mode 
•	Allows users to locate nodes across the hierarchy

•	SHALL support search by:
•	HID
•	uuid
•	Name
•	ShortDescription
•	Node Type

•	SHALL allow selection, centering, and optional action based on context

d.  Clone Node 
The Navigator Tool SHALL provide a Clone Mode enabling duplication of either: 
1.	A single node with no retained relationships ("Clone Node")
2.	A single node with all [:HAS_REQUIREMENT] relationships only ("Clone Node With Requirements")



1). Clone Node (Properties Only, No Relationships) 

This operation SHALL duplicate a single node with no retained relationships. 

---------------------------------------- 

Behavior 
•	The selected node SHALL be cloned with:
•	All properties copied
•	No inbound relationships retained
•	No outbound relationships retained
---------------------------------------- 

Insertion Behavior 
•	The cloned node SHALL be attached to a user-selected parent node using a valid relationship defined in the Core Data Model
---------------------------------------- 

Validation Rules 
•	The selected parent MUST support the relationship type per Core Data Model rules
•	Invalid parent nodes SHALL be visually disabled
•	Clone execution SHALL be blocked until a valid parent is selected
---------------------------------------- 

Identity Rules 

The cloned node SHALL receive: 
•	New uuid
•	New HID Index of the destination SoI
•	New sequence number assigned per node-type rules
---------------------------------------- 

Relationship Rules 
•	Exactly one relationship SHALL be created:
•	Between the selected parent node and the cloned node
•	No other relationships SHALL exist on the cloned node

---------------------------------------- 

General Clone Mode Requirements 
•	The user SHALL:
•	Select a source System or node
•	Select clone type
•	Select a valid destination System or parent node
•	The tool SHALL:
•	Visually distinguish valid and invalid targets
•	Prevent execution until all validation rules pass

•	Clone operations SHALL NOT:
•	Modify source nodes
•	Create invalid or incomplete graph structures


---------------------------------------- 
2). Clone Node with Requirements (Properties Only, Only [:HAS_REQUIREMENT] relationships) 

This operation SHALL duplicate a single node with all related requirements with other relationships on the node and its requirements stripped.
This operation SHALL only be made available when cloning nodes having [:HAS_REQUIREMENT] relationships (if no requirements just use 1))   

---------------------------------------- 

Behavior 
•	The selected node SHALL be cloned with:
•	All properties copied
•	No inbound relationships retained
•	No outbound relationships retained except [HAS_REQUIREMENT]

---------------------------------------- 

Insertion Behavior 
•	The cloned node SHALL be attached to a user-selected parent node using a valid relationship defined in the Core Data Model
•	(:Requirement) nodes cloned SHALL have their Orphan property set to "True"



---------------------------------------- 

Validation Rules 
•	The selected parent MUST support the relationship type per Core Data Model rules
•	Invalid parent nodes SHALL be visually disabled
•	Clone execution SHALL be blocked until a valid parent is selected
---------------------------------------- 

Identity Rules 

The cloned node and cloned related requirements nodes SHALL receive: 
•	New uuid
•	New HID Index of the destination SoI
•	New sequence number assigned per node-type rules
---------------------------------------- 

Relationship Rules 
•	Exactly one relationship SHALL be created:
•	Between the selected parent node and the cloned node
•	No other relationships SHALL exist on the cloned node except [:HAS_REQUIREMENT]
•	(:Requirement) nodes cloned SHALL be related to the parent SoI (:Purpose) node such that (:Purpose)-[HAS_REQUIREMENT]->(:Requirement)


---------------------------------------- 

General Clone Mode Requirements 
•	The user SHALL:
•	Select a source System or node
•	Select clone type
•	Select a valid destination System or parent node
•	The tool SHALL:
•	Visually distinguish valid and invalid targets
•	Prevent execution until all validation rules pass

•	Clone operations SHALL NOT:
•	Modify source nodes
•	Create invalid or incomplete graph structures


---------------------------------------- 



3.4.1.2 Hierarchy Visualization 

The tool SHALL display a graph-based visualization of the system hierarchy. 

Requirements: 
•	SHALL display (:Capability) and (:System) nodes by default
•	SHALL represent parent-child (:System) relationships directly
•	SHALL NOT display intermediate (:Element) nodes in default view
•	SHALL allow user to roll down primary or primary and  secondary nodes where primary and secondary relationships are depicted as a line
•	SHALL visually resemble:
•	Obsidian Graph View
•	Neo4j Browser graph

The graph SHALL: 
•	Support zoom, pan, and animated transitions
•	Maintain visual continuity during navigation
•	Prevent out-of-bounds navigation
---------------------------------------- 

3.4.1.3 SoI and Selection Behavior 
If no SoI has been selected, the tool SHALL scale to show the full hierarchy. 
If an SoI has already been selected, the tool SHALL initially center the graph on the current SoI and visually distinguish it from all other nodes


•	The current SoI SHALL be visually distinct
•	The tool SHALL support a temporary selection independent of the SoI
•	Selecting a node SHALL NOT change the SoI unless explicitly confirmed

Visual distinction SHALL exist for: 
•	Current SoI
•	Temporary selection
•	Search results
•	Hover state
•	Valid/invalid selection targets
---------------------------------------- 

3.4.1.4 Search and Filtering 

The tool SHALL provide a search interface. 

Capabilities: 
•	exact search by HID
•	exact search by uuid
•	partial search by Name
•	partial search by ShortDescription
•	filtering by node type
•	filtering by current mode-valid node types
•	case-insensitive matching for textual fields
•	optional exact-match toggle
•	optional incremental search while typing


Search results SHALL: 
•	Highlight nodes in the graph
•	Be listed in a synchronized results panel
•	Allow selection and graph centering

Exact HID/uuid matches SHALL: 
•	Automatically center the graph
•	Automatically select the node

Pattern: 
•	global search box sends debounced query
•	backend returns candidate nodes
•	graph recenters and highlights

---------------------------------------- 

3.4.1.5 Cross-SoI Association Support 

When invoked for association: 
•	The current SoI SHALL remain unchanged
•	The tool SHALL display:
•	Source node HID
•	Source node type
•	Intended relationship
•	Only valid target node types SHALL be selectable
•	Invalid nodes SHALL be visually muted

On confirmation: 
•	The selected node SHALL be returned to the caller
•	No navigation state SHALL change
---------------------------------------- 

3.4.1.6 Node Scope and Expansion 

Default scope: 
•	(:Capability)
•	(:System)

Optional expanded scope MAY include: 
•	(:Interface)
•	(:Requirement)
•	Other node types

When enabled: 
•	Non-System nodes SHALL be visually attached to their containing (:System)
•	SHALL be toggleable on/off
---------------------------------------- 

3.4.1.7 Information Display Controls 

Navigator Pop-up Window SHALL be composed of the following elements: 
•	HierarchySearchDialog
•	HierarchyGraphPane
•	HierarchySearchPanel
•	SelectionContextHeader
•	NodeInspectorMiniPanel
•	ModeActionBar


The tool SHALL provide toggles for: 
•	All HID
•	All Name
•	Selected HID
•	Selected Name
•	Node Type
•	Search highlights

The selected node SHALL always display: 
•	HID
•	Name
•	Type
---------------------------------------- 

3.4.1.8 Selection Actions 

Actions SHALL be mode-dependent: 


SoI Selection Mode 
•	Select SoI
•	Cancel
•	Close

Association Selection Mode 
•	Associate Selected Node
•	Cancel
•	Close

Search Mode 
•	Center on Selected
•	Use as SoI
•	Return Selected Node
•	Close

Clone Mode
•	Center on Selected
•	Expand selected (primary or primary and secondary) SoI
•	Select Node to Clone and clone type
•	Select valid parent (dim invalid)
•	Perform clone operation (copy properties, change HID, etc...)
•	Close

Only valid actions SHALL be enabled. 

---------------------------------------- 

3.4.1.9 Interaction Requirements 

The graph SHALL support: 
•	Zoom (wheel and controls)
•	Pan (drag)
•	Click selection
•	Hover highlight
•	Animated centering
•	Keyboard navigation
•	Escape to close
---------------------------------------- 

3.4.1.10 Performance Requirements 

The tool SHALL: 
•	Load hierarchy efficiently
•	Render only required nodes initially
•	Support progressive loading for large graphs
•	Maintain UI responsiveness

Exact HID/uuid lookup SHALL be faster than general search. 

---------------------------------------- 

3.4.1.11 Data Integration Requirements 

The tool SHALL retrieve data from the Backend. 

Required capabilities: 
•	System hierarchy retrieval
•	Node lookup by HID
•	Node lookup by uuid
•	Text search
•	SoI context lookup
•	Relationship validation

The tool SHALL edit nodes only when performing clone functions. 

The Navigator Tool SHALL execute all clone operations through Backend API interactions as transactional graph mutations. 

---------------------------------------- 

Required Backend Capabilities 

The Backend SHALL support: 
•	Retrieval of:

•	Nodes by HID and uuid
•	SoI membership via HID Index
•	All relationships within a given SoI

•	Validation of:

•	Allowed relationship types
•	Parent-child compatibility
•	Node type constraints
---------------------------------------- 

Clone Execution Requirements 

All clone operations SHALL: 
•	Execute as a single ACID-compliant transaction
•	Perform all validation prior to mutation
•	Fully roll back on any failure
---------------------------------------- 

Node Clone Processing Rules 

The Backend SHALL: 
1.	Clone only the selected node
2.	Copy all properties
3.	Remove all relationships
4.	Assign new identity values
5.	Create exactly one valid parent relationship
---------------------------------------- 

Identity and HID Rules 

The Backend SHALL: 
•	Generate new uuid values for all cloned nodes

•	Recompute HID values using:

•	Destination SoI Index
•	Node Type Identifier
•	Correct sequence numbering rules

•	Ensure:

•	No HID collisions
•	No sequence conflicts within the destination SoI
---------------------------------------- 

Relationship Integrity Rules 
•	No cloned node SHALL retain relationships not explicitly allowed by clone type
•	No relationships SHALL be created to nodes outside the destination SoI
•	The resulting graph SHALL conform fully to Core Data Model constraints
---------------------------------------- 

Cross-SoI Constraints 
•	Clone operations MAY create nodes in a different SoI than the current SoI

•	This behavior SHALL be explicitly allowed as an exception to the Out-of-SoI Editing Constraint

•	Clone operations SHALL NOT:

•	Modify existing nodes in other SoIs
•	Create relationships to existing external nodes
---------------------------------------- 

Error Handling 

If validation fails, the Backend SHALL return: 
•	Failure status
•	Specific reason (e.g., invalid parent, relationship violation, identity conflict)

No partial clone SHALL be committed. 

---------------------------------------- 

Performance Requirements 
•	System clone SHALL efficiently traverse nodes using HID Index
•	Clone operations SHALL scale to large SoIs
•	UI responsiveness SHALL be maintained during execution



---------------------------------------- 

3.4.1.12 Out-of-SoI Editing Constraint 

Selection of nodes outside the current SoI SHALL NOT allow editing. 

If edit is attempted: 
•	The GUI SHALL display:
•	"Navigate to: [HID] to edit"


3.4.1.13 Visual Encoding Requirements (No Icons Policy) 

The Navigator Tool SHALL use shape and color exclusively to distinguish node types and SHALL NOT use icons within graph nodes or within the legend for node-type identification. 

Node type representation SHALL be consistent across all modes and sessions. 

The Navigator Tool SHALL assign each node type a unique combination of: 
•	Shape
•	Fill color
•	Border (stroke) style

Minimum required node representations SHALL include: 
•	(:Capability): unique shape and color
•	(:System): unique shape and color
•	(:Interface): unique shape and color
•	(:Requirement): unique shape and color
•	All other node types: visually distinct but lower emphasis

Relationship types SHALL be represented using: 
•	Line style (solid, dashed, dotted)
•	Stroke thickness
•	Color

Icons SHALL NOT be used to represent: 
•	Node type
•	Relationship type
•	Node state
---------------------------------------- 

3.4.1.14 Node State Visualization Requirements 

The Navigator Tool SHALL visually distinguish node states using non-icon methods such as: 
•	Stroke thickness
•	Glow effects
•	Opacity
•	Color variation
•	Animation (subtle, non-distracting)

The following states SHALL be visually distinct: 
•	Current System of Interest (SoI)
•	Temporary selection
•	Hover state
•	Search result match
•	Valid selection target
•	Invalid selection target
•	Non-selectable node

Invalid or non-selectable nodes SHALL remain visible but visually muted. 

Valid targets SHALL be clearly distinguishable from invalid targets in all modes. 

---------------------------------------- 

3.4.1.15 Labeling and Text Display 

The Navigator Tool SHALL provide independent toggles for: 
•	HID labels
•	Name labels
•	Node Type labels

The Navigator Tool SHALL support display modes: 
•	No labels
•	Selected node labels only
•	All visible node labels

The selected node SHALL always display: 
•	HID
•	Name
•	Type

The tool SHALL implement label collision reduction strategies including: 
•	Zoom-based label visibility
•	Truncation
•	Overlap avoidance where feasible
---------------------------------------- 

3.4.1.16 Search Results Panel Behavior 

The Navigator Tool SHALL include a synchronized Search Results Panel. 

Search results SHALL: 
•	Display HID, Name, Type, and containing SoI
•	Be sorted by relevance
•	Prioritize exact HID and uuid matches

Interaction requirements: 
•	Selecting a result SHALL center and select the node in the graph
•	Selecting a node in the graph SHALL highlight the corresponding result (if present)
•	Results SHALL support incremental loading for large datasets
•	The panel SHALL support a “show more results” mechanism
---------------------------------------- 

3.4.1.17 Selected Node Detail Panel 

The Navigator Tool SHALL include a Node Detail Panel displaying: 
•	HID
•	Name
•	Type
•	uuid
•	Containing SoI
•	ShortDescription
•	Path to root (hierarchy)

The Node Detail Panel SHALL: 
•	Update immediately upon selection
•	Provide copy controls for HID and uuid
•	Display mode-relevant actions

The Path-to-Root display SHALL: 
•	Show the full hierarchy chain from Capability to selected node
•	Allow user interaction to center any node in the path
---------------------------------------- 

3.4.1.18 Minimap and Viewport Controls 

The Navigator Tool SHALL include a minimap. 

The minimap SHALL: 
•	Display the extent of the currently loaded graph
•	Indicate the current viewport
•	Allow click or drag navigation

The Navigator Tool SHALL provide: 
•	Zoom controls (in addition to mouse wheel)
•	“Center on selected node” action
•	“Fit graph to view” action
---------------------------------------- 

3.4.1.19 Graph Expansion and Scope Control 

The Navigator Tool SHALL support controlled expansion beyond default hierarchy view. 

Supported view scopes SHALL include: 
1.	Hierarchy only (Capability + Systems)
2.	Hierarchy + primary nodes
3.	Hierarchy + primary and secondary nodes

Expansion behavior SHALL: 
•	Be user-controlled and reversible
•	Preserve current selection and viewport when feasible
•	Maintain visual attachment of non-System nodes to their containing System
---------------------------------------- 

3.4.1.20 Legend Requirements 

The Navigator Tool SHALL display a legend. 

The legend SHALL describe: 
•	Node shapes and colors (by type)
•	Relationship line styles
•	Node state visual treatments

The legend SHALL: 
•	Use shape and color only (no icons)
•	Update dynamically if visual encoding changes based on mode or scope
---------------------------------------- 

3.4.1.21 Graph Export Requirements 

The Navigator Tool SHALL support export of graph visualizations. 

Export capabilities SHALL include: 
•	PNG format
•	SVG format

The export SHALL support: 
•	Current viewport export
•	Full visible graph export

The exported output SHALL preserve: 
•	Node shapes
•	Node colors
•	Labels (based on current visibility settings)
•	Relationship styles
---------------------------------------- 

3.4.1.22 Layout Stability Requirements 

The Navigator Tool SHALL maintain layout stability. 

The graph layout engine SHALL: 
•	Preserve relative node positioning where feasible

•	Avoid unnecessary full-layout re-computation during interaction

•	Maintain spatial continuity during:

•	Selection
•	Search
•	Expansion
•	Mode switching

The Navigator Tool SHALL support: 
•	Force-directed layout mode
•	Hierarchical layout mode

Switching layouts SHALL NOT reset user context unless explicitly requested. 

---------------------------------------- 

3.4.1.23 Interaction Feedback Requirements 

All user actions SHALL provide immediate visual feedback. 

The Navigator Tool SHALL: 
•	Animate centering transitions
•	Highlight selected nodes clearly
•	Indicate active mode prominently
•	Disable invalid actions visually
---------------------------------------- 

3.4.1.24 Reporting and Capture Use Case Support 

The Navigator Tool SHALL support use in reporting workflows. 

The tool SHALL: 
•	Allow users to visually frame portions of the graph
•	Preserve layout and styling in exports
•	Support repeated capture of consistent visual states


### 3.4.2  The Requirements Tool
#### 3.4.2.1 Purpose 

The Requirement Tool is a graphical add-on tool used to: 
•	Visualize (:Requirement) nodes in a hierarchical Tier structure with user selectable depth above and below the selected requirement.
•	Enable structured navigation of requirement parentage and lineage outside the System Hierarchy
•	Display (:Requirement) allocations to valid node types within the SoI
•	Allow user to develop and Allocate requirements to valid nodes.
•	Provide controlled editing, creation, and deletion of (:Requirement) nodes
•	Generate SysML 2.0-compliant requirement diagrams for analysis and reporting

The tool SHALL be visually and interactively consistent with the Navigator Tool. 

3.4.2.2 Invocation Requirements 

The Requirement Tool SHALL: 
•	Be launched only from the SSTPA Control Panel
•	Initialize to the Requirements Hierarchy view if a Data Drawer is Active for a (:Requirement) is opened and hte Tool will center on that (:Requirement) 
•	Initialize to the Requirements Allocation View if a Data Drawer is active for an entity with relationships to (:Requirement) nodes  
•	Display Requirements and other entities in a SysML 2.0 visualization.  

In the Hierarchy View, the Requirement Tool SHALL depict the focused requirement and its Heritage and lineage to a User selected depth.  

In the Allocation View, the Requirements Tool SHALL depict  the Valid Node with a [:HAS_REQUIREMENT] relationship and all allocated requirements

The User may move between views by, in the Hierarchy View selecting to display Allocations and the User selecting one. or in the Allocation View by the User selecting one Requirement and selecting it for display in Hierarchy View.


---------------------------------------- 

3.4.2.3 Supported Node Context 

The tool SHALL support invocation ONLY when the Data Drawer is open for: 
•	(:Capability)
•	(:Requirement)
•	(:Connection)
•	(:Element)
•	(:Interface)
•	(:Function)
•	(:Constraint)
•	(:Countermeasure)


The tool SHALL: 
•	Load all (:Requirement) nodes associated with the invoking node
•	Display both direct and inherited requirement relationships
•	Shift context on user action
•	Revert back to original context on user action (back arrow)
---------------------------------------- 

3.4.2.4 Requirement Hierarchy Model (Tier System) 

The Requirement Management Tool SHALL represent requirements in hierarchical tiers.  All (:Requirement) Nodes belonging to the same sub-graph (:System) as defined by HID Index will have the same tier.  Tier number is the distance between the sub-graph and the (:Capability) such that:
•	Tier 0 → (:Capability)-level requirements
•	Tier 1 → (:Requirement) in a sub-graph whose (:System) is a direct child of (:Capability)
•	Tier 2 → (:Requirement) in a sub-graph whose (:System) is a direct child of an (:Element) in a sub-graph whose (:System) is a direct child of (:Capability)
•	Tier N → (:Requirement) where N is the number of sub-graph relationships to :(Capability)

 
The tool SHALL: 
•	Determine tier level based on parent-child relationship count to :(Capability)
•	Support cross-SoI parentage relationships (note:  parent-child relationships between requirements may cross tier boundaries e.g. a Tier 5 requirement could parent a Tier 4 requirement if that other requirement were in another sub-graph.  This behavior is bad practice, but does happen in large complex systems and will be allowed so long as the two SOI's common ancestor is of a lower tier then both of them (prevents one parenting an ancestor).  
•	Display tiers visually in structured layout
---------------------------------------- 

3.4.2.5 Visualization Requirements (SysML 2.0 Compliance) 

The tool SHALL render requirement diagrams consistent with SysML 2.0 requirement diagram conventions. 

The diagram SHALL include as user toggleable display properties:  
•	Requirement nodes displayed as structured blocks containing:
•	HID
•	Name
•	Requirement Statement (RStatement)
•	Parent-child relationships rendered as directed links
•	Association relationships to:
•	(:Purpose)
•	(:Connection)
•	(:Element)
•	(:Interface)
•	(:Function)
•	(:Constraint)
•	(:Countermeasure)

Relationship representation SHALL use: 
•	Directed edges for parentage
•	Labeled edges for association types
•	No icons; shape and color only
---------------------------------------- 

3.4.2.6 Visual Encoding Requirements 

The Requirement Management Tool SHALL follow the same visual encoding rules as the Navigator Tool for interfaces and controls otherwise SysML 2.0 is authoritative in the diagram itself.  
•	Node types SHALL be distinguished by shape and color only
•	Requirement nodes SHALL have a unique shape and color
•	Parent-child relationships SHALL be visually distinct from association relationships
•	Node states SHALL include:
•	Selected
•	Hover
•	Editable
•	Invalid (when applicable)
---------------------------------------- 

3.4.2.7 Parentage and Lineage Controls 

The tool SHALL provide user controls to: 

Parent Traversal 
•	Display parent requirements up to a user-selected tier depth
•	Allow user input (e.g., Tier N range limit)
•	Dynamically update the diagram

Child Traversal 
•	Display child requirements down to a user-selected tier depth
•	Support expansion and collapse

Combined View 
•	Allow simultaneous display of parent and child relationships
---------------------------------------- 

3.4.2.8 Interaction Requirements 

The diagram SHALL support: 
•	Zoom (mouse wheel and controls)
•	Pan (drag)
•	Node selection
•	Hover highlighting
•	Animated centering
•	Expand/collapse of requirement branches

Selecting a node SHALL: 
•	Highlight it
•	Display its properties in a Requirement Detail Panel
---------------------------------------- 

3.4.2.9 Requirement Editing Model 

The Requirement Management Tool SHALL support editing of: 
•	RStatement
•	VMethod
•	VStatement
•	Name
•	ShortDescription

Editing SHALL: 
•	Use a right-side Data Drawer consistent with GUI standards
•	Follow the same staged editing model
•	Require Commit confirmation

All edits SHALL: 
•	Be validated via Backend API
•	Be persisted only after successful Commit
---------------------------------------- 

3.4.2.10 Requirement Creation 

The tool SHALL allow creation of new (:Requirement) nodes. 

Creation SHALL require: 
•	Selection of a valid parent requirement and associated node (node with [:HAS_REQUIREMENT] relationship)
•	Assignment of required properties

New nodes SHALL: 
•	Receive valid HID and uuid values
•	Be inserted into the correct tier

---------------------------------------- 

3.4.2.11 Requirement Deletion 

Deletion SHALL: 
•	Follow the Alert / Confirm pattern
•	Identify and display dependent relationships
•	Warn of orphaned nodes

Deletion SHALL NOT: 
•	Cascade outside the current SoI without explicit confirmation
---------------------------------------- 

3.4.2.12 Association Management 

The tool SHALL allow association of requirements to: 
•	(:Capability)
•	(:System)
•	(:Element)
•	(:Interface)
•	(:Function)
•	(:Countermeasure)

The tool SHALL: 
•	Validate all associations via Backend API
•	Prevent invalid associations
•	Visually distinguish valid vs invalid targets
---------------------------------------- 

3.4.2.13 Data Synchronization 

Upon Commit: 
•	The Backend SHALL persist all changes transactionally
•	The Main Panel SHALL refresh automatically
•	The Data Drawer SHALL update to reflect committed state


3.4.2.14 Export Requirements 

The tool SHALL support export of requirement diagrams. 

Supported formats SHALL include: 
•	PNG
•	SVG

Exports SHALL: 
•	Preserve SysML 2.0 layout
•	Preserve node shapes and colors
•	Preserve labels and relationships

The user SHALL be able to export: 
•	Current viewport
•	Full diagram
---------------------------------------- 

3.4.2.15 Performance Requirements 

The tool SHALL: 
•	Efficiently load requirement hierarchies
•	Support progressive loading for large requirement sets
•	Maintain UI responsiveness during expansion
---------------------------------------- 

3.4.2.16 Backend Integration Requirements 

The Backend SHALL support: 
•	Retrieval of requirement hierarchies
•	Retrieval of parent/child relationships
•	Retrieval of associated system elements
•	Validation of requirement associations
•	Transactional creation, update, and deletion

All operations SHALL be ACID-compliant. 

---------------------------------------- 
3.4.2.17 Layout Requirements 

The tool SHALL support: 
•	Hierarchical layout (default for tiers)
•	Stable layout during interaction
•	Optional re-layout on demand

The layout SHALL: 
•	Visually group requirements by tier
•	Minimize edge crossings where feasible
---------------------------------------- 
3.4.2.18 Reporting Integration 

The Requirement Management Tool SHALL support: 
•	Generation of diagrams suitable for reports
•	Consistent visual formatting across exports
•	Reproducible diagram states
---------------------------------------- 


### 3.4.3 Reports Dropdown Menu
The Reports Dropdown Menu SHALL list the following reports to create
System Description
System Specification
Requirement‑Traceability Gap Analysis
Controls List

### 3.4.3.1 System Description Report
System Description Report is a text based hierarchical description of the SoI, its primary nodes and relationships followed by its secondary nodes and relationships.
Report SHALL be in text, markdown, MS Word or PDF format.
  
### 3.4.3.2 System Specification Report
System Specification Report is a text based list of requirements for the SoI organized by the node they are related to within the SoI. It begins with a description of the SoI and its properties. One section per primary element type with [HAS_REQUIREMENT] relationship.  Sub-section for each entity of the type followed by an ordered list of requirements showing uuid and RequirementStatement properties for each.
Report SHALL be in text, markdown, MS Word or PDF format.

### 3.4.3.3 Requirement‑Traceability Gap Analysis
Requirement‑Traceability Gap Analysis is a text based report identifying problematic requirements in the SoI.  This report is less informational and focused to remediating action. It is organized in the same way as the System Specification Report excepting when referring to Requirement properties it shows the UUID followed by the analytical properties: Baseline, Orphan and Barren.   Orphan and Barren properties are not user editable, but SHALL be by the generation of this report.  

Baseline is not set by running this report but the (:Requirement) Baseline property is reported (note:  projects deal with baselines in a number of ways and the tool must be flexible with this property to support most use cases).

For a (:Requirement) property "Orphan" SHALL be true if any of these is true:
1.  It has no parent (:Requirement)

In other words, a Requirement cannot be created without a parent, so an Orphan is likely the result of a node deletion and the Requirement is flagged foe reparenting or removal.

For a (:Requirement) property "Barren" SHALL be true if any of these is true:
1.  It has no child (:Requirement)
2.  it has no [HAS_REQUIREMENT] relationship other than with (:Purpose).  

In other words, every requirement should be allocated to an Interface, Function or Element though there are valid exceptions.

For this analysis, valid incoming [:HAS_REQUIREMENT] sources are: 
(:Connection), (:Interface), (:Function), (:Element), (:Purpose), (:Constraint), and (:Countermeasure). Allocated or Derived Requirements must be assigned to nodes other than Purpose.


3.4.4 The Reference Tool 

The Reference Tool is an Add-on Tool used to navigate imported external reference frameworks, inspect authoritative reference item properties, and assign valid external references to selected valid SSTPA nodes. 

The tool described here SHALL be branded at top of window as “Reference Tool”. 


3.4.4.1 Purpose 

The Reference Tool SHALL allow the User to: 
1.	Relate a valid Reference data to a valid node type in the active Data Drawer in Assignment Mode
2.	Navigate imported reference frameworks and display its node properties and relationships without changing the current SoI in Research Mode
3.	Search for a specific external reference item

Reference Tool SHALL initialize in Assignment Mode when the Data Drawer for a valid node type is open otherwise Reference Tool opens in Research Mode.
The User SHALL be able to switch the Research Tool into Research Mode at any time and switch back to the view in Assignment Mode.
The User will not be able to switch to Assignment Mode without a valid node type in the active Data Drawer.


3.4.4.2 Invocation 

The Reference Tool SHALL initialize into the Assignment Mode for the following supported node types. 

•	(:Control)
•	(:Element)
•	(:System)
•	(:Hazard)
•	(:Attack)
•	(:Countermeasure)

On launch, the tool SHALL display: 
•	Source node HID
•	Source node Name
•	Source node Type
•	Current SoI
•	Allowed framework filters for that source node type

Launching the tool SHALL NOT change the current SoI. 


3.4.4.3 Modes of Operation 

The Reference Catalog Tool SHALL support two modes: 


a. Research Mode 
•	Allows the User to navigate imported framework data
•	SHALL NOT modify the source SSTPA node
•	SHALL support hierarchical navigation where available
•	SHALL display the selected reference item in a read-only inspector

b. Assignment Mode 
•	Allows the User to select a valid imported reference item and assign it to the source SSTPA node by switching to Assignment Mode
•	Allows User to follow internal references to a valid node for assignment 
•	SHALL validate the proposed assignment through the Backend prior to commit
•	SHALL return the selected reference item to the calling Data Drawer context

3.4.4.4 Layout 

The Reference Catalog Tool pop-up window SHALL include: 
•	ReferenceCatalogDialog
•	ReferenceFrameworkSelector
•	ReferenceTypeFilterBar
•	ReferenceSearchPanel
•	ReferenceHierarchyPane
•	ReferenceResultsGrid
•	ReferenceInspectorPanel
•	ReferenceActionBar

3.4.4.5 Framework Selection and Filtering 

The tool SHALL allow the User to filter by: 
•	Framework
•	Framework version
•	Imported reference item type

The tool SHALL restrict displayed item types based on the source SSTPA node type. 

Invalid framework item types for the current source SSTPA node SHALL be visually muted or hidden. 


3.4.4.6 Search and Locate 

The tool SHALL provide a search interface. 

Search SHALL support: 
•	exact search by ExternalID
•	partial search by Name
•	partial search by ShortDescription
•	filtering by framework
•	filtering by item type
•	optional incremental search while typing

Search results SHALL: 
•	be listed in a synchronized results panel
•	allow selection of a result
•	update the read-only inspector on selection

3.4.4.7 Hierarchy Navigation 

Where the imported framework supports hierarchy, the tool SHALL allow navigation by parent-child structure. 

The tool SHALL support at minimum: 
•	family/control/enhancement style navigation for NIST SP 800-53
•	tactic/technique style navigation for ATT&CK
•	category/property-threat-mitigation style navigation for EMB3D

The tool MAY additionally display related imported reference items in a secondary related-items panel. 


3.4.4.8 Read-Only Reference Inspector 

The selected imported reference item SHALL be displayed in a read-only inspector. 

At minimum, the inspector SHALL display: 
•	Framework name
•	Framework version
•	ExternalID
•	Reference item type
•	Name
•	ShortDescription
•	LongDescription
•	SourceURI

The inspector MAY additionally display: 
•	parent item
•	child items
•	related items
•	framework-specific fields

The User SHALL NOT edit imported reference item content. 


3.4.4.9 Selection Actions 

Actions SHALL be mode-dependent. 


Browse / Inspect Mode 
•	Expand Selected
•	Collapse Selected
•	Center on Selected
•	Close

Assign Reference Mode 
•	Assign Selected Reference
•	Cancel
•	Close

Only valid actions SHALL be enabled. 


3.4.4.10 Data Drawer Integration 

On successful assignment, the Reference Catalog Tool SHALL return the selected imported reference item to the calling Data Drawer. 

The Data Drawer SHALL display assigned external references in a relationship group using: 
•	ExternalID
•	Name
•	Framework
•	Reference item type

The Data Drawer SHALL allow: 
•	launching the Reference Catalog Tool to add a reference
•	removing an existing [:REFERENCES] relationship
•	opening an assigned reference item in read-only inspection mode

3.4.4.11 Out-of-SoI Editing Constraint 

The Reference Catalog Tool SHALL NOT be treated as editing of the imported reference item or navigation to another SoI. 

Assignment of a [:REFERENCES] relationship to an imported reference item SHALL be allowed even though the imported reference item is not part of the current SoI. 

The tool SHALL NOT allow editing of any node outside the current SoI. 


3.4.4.12 Performance Requirements 

The tool SHALL: 
•	load framework metadata efficiently
•	render only required results initially
•	support progressive loading for large framework datasets
•	maintain UI responsiveness during search and navigation

Exact ExternalID lookup SHALL be faster than general text search. 

---------------------------------------- 

3.4.4.13 Data Integration Requirements 

The Reference Catalog Tool SHALL retrieve data from the Backend. 

Required capabilities: 
•	framework list retrieval
•	framework version retrieval
•	reference item lookup by ExternalID
•	reference item lookup by uuid
•	framework text search
•	framework hierarchy retrieval
•	related reference item retrieval
•	assignment validation
•	reference relationship creation and removal

The tool SHALL edit SSTPA nodes only by creating or removing [:REFERENCES] relationships through the Backend. 



3.4.4.14 Test and Verification Requirements 

The Reference Catalog Tool SHALL be verified through test and analysis. 

The system SHALL verify that: 
•	imported framework items are retrievable by ExternalID
•	framework hierarchy is navigable where source data supports hierarchy
•	imported reference item properties are displayed read-only
•	invalid source-node-to-reference-item assignments are rejected
•	valid assignments create exactly one [:REFERENCES] relationship
•	removal of an assigned reference deletes only the [:REFERENCES] relationship
•	imported reference items are not modified by any GUI action
•	all assignment mutations are transactional and roll back on failure


--------------------------------------

### 3.4.5 The State Tool  
The State Tool is an Add-on Tool used to visualize, create, edit, and analyze SysML 2 aligned State Transition diagrams for the current System of Interest (SoI) using existing (:State) nodes and (:State)-[:TRANSITIONS_TO]->(:State) relationships. 
 
The tool described here SHALL be branded at top of window as "State Tool". 

#### 3.4.5.1 Purpose 
The State Tool SHALL allow the User to: 
1. Display (:State) nodes and [:TRANSITIONS_TO] relationships in a SysML 2 aligned state-transition visualization 
2. Create new (:State) nodes within the active SoI 
3. Create, edit, and remove [:TRANSITIONS_TO] relationships between valid (:State) nodes in the active SoI 
4. View and edit transition relationship properties defined in Section 1.3.8.9 
5. Associate and/or create related (:Requirement), (:Countermeasure), and (:Hazard) nodes in the active SoI 
6. Display state transition criteria and related node relationships in a graph-like analytical view 
7. Support analysis of how Hazards, Countermeasures, and Requirements relate to state behavior without changing the canonical Core Data Model representation of transitions as relationships rather than nodes 

The State Tool SHALL be visually and interactively consistent with the Navigator Tool and Requirement Tool. 


#### 3.4.5.2 Invocation 
 The State Tool SHALL: 
• Be launched from the SSTPA Control Panel "State Tool" button 
• Initialize to the State Diagram View if a Data Drawer is active for a (:State) node and center on that (:State) 
• Initialize to the State Context View if a Data Drawer is active for a (:Countermeasure), (:Hazard), or (:Requirement) related to one or more (:State) nodes 
• Initialize to the full active-SoI State Diagram View when no specific (:State) context is active 

 
#### 3.4.5.3 Supported Node Context 
 The tool SHALL support invocation when Data Drawer is open for: 
• (:State) 
• (:Countermeasure) 
• (:Hazard) 
• (:Requirement) 
• (:System) 
 
The tool SHALL: 
• Load all (:State) nodes in the current SoI needed for the selected view 
• Load all [:TRANSITIONS_TO] relationships among displayed (:State) nodes 
• Load related (:Hazard), (:Countermeasure), and (:Requirement) nodes needed for the selected context 
• Allow the User to shift focus without changing the current SoI 
• Provide a back action to return to the invoking context 
• Allow user to move and arrange objects on the canvas



#### 3.4.5.4 State Model Compliance  
The State Tool SHALL use the Core Data Model representation of state behavior already defined by the SRS. 

The State Tool SHALL: 
• Treat (:State)-[:TRANSITIONS_TO]->(:State) as the canonical representation of a transition 
• SHALL NOT introduce a Transition node into the Core Data Model 
• Distinguish the semantic role of transitions using relationship properties, including TransitionKind 
• Support transitions whose TransitionKind is FUNCTIONAL, COUNTERMEASURE_REQUIRED, or BOTH 
• Support transition traceability to a governing (:Countermeasure) by RequiredByCountermeasureHID and/or RequiredByCountermeasureUUID when applicable 

The tool SHALL preserve the preferred modeling rule that where the same source and destination (:State) pair is used both for ordinary behavior and to satisfy a (:Countermeasure), the preferred representation is a single [:TRANSITIONS_TO] relationship with TransitionKind = BOTH. 

#### 3.4.5.5 Modes of Operation 
The State Tool SHALL support three modes: 

a. Diagram View 
• Displays the current SoI state-transition diagram 
• Allows selection of (:State) nodes and [:TRANSITIONS_TO] relationships 
• Allows creation of new (:State) nodes 
• Allows creation of new [:TRANSITIONS_TO] relationships 
• Allows editing of displayed node and relationship properties through standard SSTPA edit patterns 

b. Context View 
• Displays a selected (:State), (:Countermeasure), (:Hazard), or (:Requirement) and its related state-transition context 
• SHALL highlight related transitions and related nodes 
• SHALL support filtering to the selected analytical context 
• SHALL NOT change the current SoI 

c. Criteria / Relationship View 
• Displays a graph-like analytical view centered on selected (:State) nodes and [:TRANSITIONS_TO] relationships 
• SHALL show transition criteria such as Trigger, GuardCondition, Rationale, Priority, and ResidualRiskNote 
• SHALL show related (:Hazard), (:Countermeasure), and (:Requirement) nodes and their relationships to state behavior 
• SHALL support expansion and collapse of related-node groupings 


#### 3.4.5.6 Visualization Requirements 
The State Tool SHALL render diagrams consistent with SysML 2 state-transition diagram conventions to the maximum extent practical within the SSTPA Tool visual style. 

The diagram SHALL include user-toggleable display of: 
• (:State) nodes displayed as SysML 2 aligned state blocks 
• HID 
• Name 
• ShortDescription 
• [:TRANSITIONS_TO] relationships rendered as directed transitions 
• Transition labels derived from relationship properties, including Trigger and GuardCondition where present 
• Visual distinction for TransitionKind values FUNCTIONAL, COUNTERMEASURE_REQUIRED, and BOTH 
• Optional display of related (:Hazard), (:Countermeasure), and (:Requirement) nodes as analytical overlays or side-panel-linked objects 

The diagram SHALL use: 
• Directed edges for transitions 
• Shape and color only for node-type distinction 
• No icons within the diagram for node or relationship type identification 

#### 3.4.5.7 Visual Encoding Requirements 

The State Tool SHALL follow the same visual encoding rules established for the Navigator Tool unless SysML 2 convention is authoritative within the diagram itself. 

The tool SHALL visually distinguish: 
• (:State) 
• (:Hazard) 
• (:Countermeasure) 
• (:Requirement) 

The tool SHALL visually distinguish transition semantics using non-icon methods such as: 
• Line style 
• Stroke thickness 
• Color 
• Label treatment 
• Glow or highlight state 

The following transition states SHALL be visually distinct: 
• Selected transition 
• Hover state 
• Editable transition 
• Invalid transition proposal 
• TransitionKind = FUNCTIONAL 
• TransitionKind = COUNTERMEASURE_REQUIRED 
• TransitionKind = BOTH 

#### 3.4.5.8 Interaction Requirements 

The diagram SHALL support: 
• Zoom (mouse wheel and controls) 
• Pan (drag) 
• Node selection 
• Relationship selection 
• Hover highlighting 
• Animated centering 
• Expand/collapse of related analytical overlays 
• Keyboard navigation 
• Escape to close 

Selecting a (:State) node SHALL: 
• Highlight the node 
• Display its properties and related relationships in a State Detail Panel 

Selecting a [:TRANSITIONS_TO] relationship SHALL: 
• Highlight the relationship 
• Display its relationship properties in a Transition Detail Panel 
• Display related (:Countermeasure), (:Hazard), and (:Requirement) nodes where present 

#### 3.4.5.9 State Creation 

The State Tool SHALL allow creation of new (:State) nodes within the current SoI. 

Creation SHALL: 
• Use the standard SSTPA staged editing and Commit confirmation model 
• Assign valid HID and uuid values 
• Assign the new node to the active SoI 
• Open the created (:State) in the standard Data Drawer or State Detail Panel for further editing 

New (:State) nodes SHALL receive: 
• Common properties per Section 1.3.7 
• Type-specific defaults per Section 1.3.8.9 
• Correct Owner, Creator, and LastTouch behavior per Section 1.3.7.1 

#### 3.4.5.10 Transition Creation and Editing 

The State Tool SHALL allow the User to create a [:TRANSITIONS_TO] relationship between two valid (:State) nodes in the active SoI. 

Transition creation SHALL: 
• Require selection of a source (:State) 
• Require selection of a destination (:State) 
• Stage relationship properties prior to Commit 
• Validate duplicate-logical-relationship constraints before Commit 
• Validate countermeasure traceability fields where TransitionKind requires them 

The tool SHALL allow editing of the following transition relationship properties: 
• TransitionKind 
• Trigger 
• GuardCondition 
• Rationale 
• RequiredByCountermeasureHID 
• RequiredByCountermeasureUUID 
• Priority 
• ResidualRiskNote 

If TransitionKind = COUNTERMEASURE_REQUIRED or BOTH, the tool SHALL require RequiredByCountermeasureHID and/or RequiredByCountermeasureUUID to identify the governing (:Countermeasure) before Commit is allowed. 

#### 3.4.5.11 Related Node Association and Creation 

The State Tool SHALL allow association and/or creation of the following related node types within the current SoI: 
• (:Requirement) 
• (:Countermeasure) 
• (:Hazard) 

The tool SHALL support: 
• Associating an existing (:Hazard) to a (:State) via [:HAS_HAZARD] 
• Associating an existing (:Countermeasure) to a (:State) via [:APPLIES_TO_STATE] 
• Associating an existing (:Requirement) to a (:Countermeasure) via [:HAS_REQUIREMENT] 
• Creating new related (:Hazard), (:Countermeasure), and (:Requirement) nodes using standard SSTPA staged editing behavior 

The State Tool SHALL NOT create invalid direct relationships not defined in the Core Data Model. 

For (:Requirement), the State Tool SHALL support its creation and association through the valid node that owns the requirement relationship, typically (:Countermeasure), (:Purpose), or another valid requirement-bearing node per the Core Data Model. 

#### 3.4.5.12 Criteria / Relationship View Requirements 

The Criteria / Relationship View SHALL provide a graph-like analytical display of: 
• Selected (:State) nodes 
• Their outgoing and incoming [:TRANSITIONS_TO] relationships 
• Transition criteria and analysis properties 
• Related (:Hazard) nodes 
• Related (:Countermeasure) nodes 
• Related (:Requirement) nodes through valid intermediate nodes 

This view SHALL allow the User to: 
• Filter by TransitionKind 
• Filter by selected (:Countermeasure) 
• Filter by selected (:Hazard) 
• Filter by selected (:Requirement) 
• Toggle display of transition criteria labels 
• Toggle display of related-node overlays 
• Center on a selected node or transition 
• Export the current analytical view 

#### 3.4.5.13 Validation Requirements 

The State Tool SHALL validate all proposed mutations through the Backend API prior to Commit. 

Validation SHALL confirm: 
• Both transition endpoints are valid (:State) nodes 
• Both endpoint (:State) nodes belong to the same SoI unless explicitly allowed by future analytical extension 
• Duplicate logical [:TRANSITIONS_TO] relationships do not already exist unless distinguished by valid relationship properties 
• TransitionKind values are valid 
• RequiredByCountermeasureHID and/or RequiredByCountermeasureUUID identify an existing (:Countermeasure) when TransitionKind = COUNTERMEASURE_REQUIRED or BOTH 
• Any referenced governing (:Countermeasure) belongs to the same SoI as both endpoint (:State) nodes unless explicitly justified as a cross-SoI analytical relationship 
• All other proposed relationships conform to the Core Data Model 

The API SHALL return: 
• Valid / invalid 
• Reason for invalidity 

#### 3.4.5.14 Data Drawer Integration 

On successful selection from the State Tool, the calling Data Drawer SHALL be able to display: 
• Related (:State) nodes 
• [:TRANSITIONS_TO] relationships and their properties 
• Related (:Hazard) nodes 
• Related (:Countermeasure) nodes 
• Related (:Requirement) nodes as reachable through valid requirement-bearing nodes 

The Data Drawer SHALL allow: 
• Launching the State Tool from a valid node context 
• Removing a valid relationship subject to orphan and deletion rules already defined by the SRS 
• Opening selected related nodes for edit within the SoI 

#### 3.4.5.15 Export Requirements 

The State Tool SHALL support export of state diagrams and analytical views. 

Supported formats SHALL include: 
• PNG 
• SVG 

Exports SHALL preserve: 
• Node shapes 
• Node colors 
• Labels based on current visibility settings 
• Relationship directionality 
• Relationship style distinctions 
• Visible transition criteria labels 

The User SHALL be able to export: 
• Current viewport 
• Full visible diagram 

#### 3.4.5.16 Performance Requirements 

The State Tool SHALL: 
• Efficiently load state-transition diagrams for the active SoI 
• Support progressive loading for large state-transition graphs 
• Maintain UI responsiveness during expansion, filtering, and selection 
• Use bounded traversal for recursive transition analysis 
• Avoid unbounded recursive expansion of state-transition relationships 

Exact HID and uuid lookup for (:State) nodes SHALL be faster than general text search. 

#### 3.4.5.17 Data Integration Requirements 

The State Tool SHALL retrieve data from the Backend. 

Required capabilities: 
• Retrieval of all (:State) nodes within the current SoI 
• Retrieval of [:TRANSITIONS_TO] relationships and their properties 
• Retrieval of related (:Hazard), (:Countermeasure), and (:Requirement) context 
• Validation of transition creation and edit operations 
• Transactional creation, update, and deletion of permitted nodes and relationships 

The State Tool SHALL execute all mutations through Backend API interactions as transactional graph mutations. 

All write operations SHALL be ACID compliant. 

#### 3.4.5.18 Layout Stability Requirements 

The State Tool SHALL maintain layout stability during: 
• Selection 
• Search 
• Filtering 
• Expansion 
• Mode switching 

The layout engine SHALL: 
• Preserve relative node positioning where feasible 
• Avoid unnecessary full-layout recomputation during interaction 
• Support a state-diagram-oriented layout mode 
• Support a force-directed analytical layout mode for Criteria / Relationship View 

Switching layouts SHALL NOT reset user context unless explicitly requested. 

#### 3.4.5.19 Reporting and Capture Use Case Support 

The State Tool SHALL support use in reporting workflows. 

The tool SHALL: 
• Allow users to visually frame portions of the state diagram 
• Preserve layout and styling in exports 
• Support repeated capture of consistent visual states 
• Support generation of figures suitable for insertion into System Description, System Specification, and future analytical reports 

#### 3.4.5.20 Test and Verification Requirements 

The State Tool SHALL be verified through test and analysis. 

The system SHALL verify that: 
• (:State) nodes in the active SoI are retrievable and displayable 
• valid [:TRANSITIONS_TO] relationships are retrievable and displayable 
• new (:State) nodes receive correct HID and uuid values 
• valid transitions can be created and edited 
• invalid transitions are rejected 
• TransitionKind semantics are correctly represented 
• required countermeasure traceability is enforced when TransitionKind = COUNTERMEASURE_REQUIRED or BOTH 
• the tool does not create Transition nodes outside the Core Data Model 
• all permitted mutations are transactional and roll back on failure 
• exported diagrams preserve visible relationship direction and labeling 

---------------------------------------- 

### 3.4.6 The Flow Tool 

The Flow Tool is an Add-on Tool used to visualize, create, edit, and analyze Functional Flow and STPA Control Flow diagrams for the current System of Interest (SoI) using (:Function), (:Interface), (:Connection), and ControlStructure-related nodes. 

The tool described here SHALL be branded at top of window as "Flow Tool". 

#### 3.4.6.1 Purpose 

The Flow Tool SHALL allow the User to: 

• Visualize and analyze Functional Flow between (:Function) and (:Function) nodes and (:Function) and (:Interface) nodes  
• Visualize and analyze STPA Control Flow using (:ControlStructure) roles
• Store and retrieve visualizations from a structured JSON file as property of  (:ControlStructure) or (:FunctionalFlow) 
• Create, edit, and remove flow relationships between (:Function) and (:Interface) nodes  
• Create and relate (:Function), (:Interface), (:Requirement), and (:Countermeasure) nodes  
• Associate (:Interface) nodes to (:Connection) nodes (including cross-SoI ownership cases)  
• Define the nature of flow relationships including physical and logical (OSI-based) characteristics  
• Create and manage Feedback relationships in flows  
• Display and filter flows associated with (:Countermeasure) nodes  
• Assign (:Function) and (:Interface) nodes to STPA roles in (:ControlStructure)  
• Create and assign (:ControlAction) and (:Feedback) nodes 
• Commit validated changes to the Backend  

The Flow Tool SHALL be visually and interactively consistent with the Navigator Tool and Requirements Tool. 

---------------------------------------- 

#### 3.4.6.2 Modes of Operation 

The Flow Tool SHALL support two modes: 

a. "Functional Flow" Mode  
b. "STPA Control Flow" Mode  

For both modes the visualization SHALL come from a  structured JSON file as property of either  (:ControlStructure) or (:FunctionalFlow) which will hold the elements and their position on the canvas.
The Tool SHALL operate on the Structured JSON file used for visualization to capture User changes.


The User SHALL be able to switch between modes without changing the current SoI. 

---------------------------------------- 

#### 3.4.6.3 Invocation 

The Flow Tool SHALL: 

• Be launched from the SSTPA Control Panel  
• Initialize to Functional Flow Mode if a Data Drawer is open for (:Function) or (:Interface) and by default  
• If a Data Drawer is open for (:Function) or (:Interface), center and focus on that node
• Initialize to STPA Control Flow mode if a Data Drawer is open for (:ControlStructure), (:ControlAlgorithm), (:ProcessModel), (:ControlledProcess), (:ControlAction), or (:Feedback).
• If neither Drawer is open, the default SHALL be the Functional Flow Mode.
• If there is more than one (:FunctionalFlow) node, the Tool SHALL present the name properties for all and allow the User to select which to operate on. 
• If there is more than one (:ControlStructure) node, the Tool SHALL present the name properties for all and allow the User to select which to operate on.
• Load all relevant flow relationships within the active SoI  
• NOT change the current SoI  

---------------------------------------- 

#### 3.4.6.4 Scope Constraints 

The Flow Tool SHALL: 

• Restrict node creation, relationship creation, and editing to the current SoI  
• Allow association of (:Interface) nodes in the SoI to (:Connection) nodes owned by another SoI  
• NOT allow editing of nodes outside the current SoI  
• Enforce all Core Data Model constraints  

---------------------------------------- 

#### 3.4.6.5 Functional Flow Mode 

In Functional Flow Mode, the tool SHALL: 

• Display (:Function) and (:Interface) nodes  
• Display relationships: 
  • (:Function)-[:FLOWS_TO_FUNCTION]->(:Function) 
  • (:Function)-[:FLOWS_TO_INTERFACE]->(:Interface) 
  • (:Interface)-[:CONNECTS]->(:Function) 
  • (:Interface)-[:PARTICIPATES_IN]->(:Connection) 

• Allow creation and editing of these relationships  
• Allow creation of new (:Function) and (:Interface) nodes  
• Allow assignment of (:Requirement) and (:Countermeasure) nodes  
• Allow creation and display of Feedback relationships 
• Allow user to move and arrange objects on the canvas while maintaining relationships. 

---------------------------------------- 

#### 3.4.6.6 Countermeasure Overlay 

The Flow Tool SHALL allow: 

• Display of nodes and relationships associated with (:Countermeasure)  
• Filtering of flows impacted by Countermeasures  
• Visualization of how Countermeasures alter flow behavior  

---------------------------------------- 

#### 3.4.6.7 Feedback Relationships 

The Flow Tool SHALL support Feedback relationships: 

• (:ControlledProcess)-[:PRODUCES]->(:Feedback)  
• (:Feedback)-[:INFORMS]->(:ProcessModel)  

The tool SHALL allow: 

• Creation of (:Feedback) nodes  
• Assignment and editing of properties and their values 
• Visualization within both modes  

---------------------------------------- 

#### 3.4.6.8 STPA Control Flow Mode 

In STPA Control Flow Mode, the tool SHALL: 

• Display (:ControlStructure) and its child nodes: 
  • (:ControlAlgorithm) 
  • (:ControlledProcess) 
  • (:ProcessModel) 
  • (:ControlAction) 
  • (:Feedback) 

• Allow casting of (:Function) and (:Interface) nodes into STPA roles  
• Display Control Flow relationships: 
  • (:ControlAlgorithm)-[:GENERATES]->(:ControlAction) 
  • (:ControlAction)-[:COMMANDS]->(:ControlledProcess) 
  • (:ControlledProcess)-[:PRODUCES]->(:Feedback) 
  • (:Feedback)-[:INFORMS]->(:ProcessModel) 
  • (:ProcessModel)-[:TUNES]->(:ControlAlgorithm) 

---------------------------------------- 

#### 3.4.6.9 STPA Role Assignment Rules 

The Flow Tool SHALL enforce: 

• (:Function) MAY be assigned to: 
  • (:ControlAlgorithm) 
  • (:ControlledProcess) 
  • (:ProcessModel) 

• (:Interface) MAY be assigned to: 
  • (:ControlAlgorithm) 
  • (:ControlledProcess) 

• (:Interface) SHALL NOT be assigned to (:ProcessModel) 

• Validation SHALL reject invalid assignments  

---------------------------------------- 

#### 3.4.6.10 ControlAction and Feedback Nodes 

The Flow Tool SHALL allow: 

• Creation of (:ControlAction) nodes  
• Creation of (:Feedback) nodes
• Assignment and editing of properties and their values to include relating (:Hazard) nodes to (:ControlAction) and relating (:Countermeasure) nodes to (:Feedback).   
• Assignment into STPA Control Flow  

---------------------------------------- 

#### 3.4.6.11 Visualization Requirements 

The Flow Tool SHALL: 

• Render diagrams in a canvas consistent with SysML 2 functional and control flow conventions  
• Use directed edges for flow  
• Distinguish: 
  • Function vs Interface 
  • Physical vs Logical flow 
  • Control (STPA) vs linear vs Feedback flow 
• Use shape and color only (no icons)  
• Use a right side panel in the pop-up window for all editing arranged in a manner similar to its GUI data drawer representation.  
• Display relationship properties  

---------------------------------------- 

#### 3.4.6.12 Interaction Requirements 

The tool SHALL support: 

• Zoom, pan, selection, hover  
• Node and relationship selection  
• Editing via Data Drawer  
• Animated centering  
• Expand/collapse  
• User interaction to position objects on the canvas

Selecting a node SHALL: 

• Highlight the node  
• Display its properties  

Selecting a relationship SHALL: 

• Display its properties  

---------------------------------------- 

#### 3.4.6.13 Node and Relationship Creation 

The tool SHALL allow creation of: 

• (:Function) 
• (:Interface) 
• (:Requirement) 
• (:Countermeasure) 

The tool SHALL: 

• Assign valid HID and uuid  
• Enforce SoI constraints  
• Use staged edit + Commit model  

---------------------------------------- 

#### 3.4.6.14 Validation Requirements 

The Backend SHALL validate: 

• Same-SoI constraints for flow relationships  
• Valid node types for relationships  
• STPA role assignment rules  
• Duplicate relationship prevention  
• Valid Connection participation rules  

---------------------------------------- 

#### 3.4.6.15 GUI and Data Drawer Integration 

The GUI SHALL refresh after the Flow Tool performs a commit: 

---------------------------------------- 

#### 3.4.6.16 Export Requirements 

The tool SHALL support export: 

• PNG  
• SVG  

Exports SHALL preserve: 

• Node shapes and colors  
• Relationship styles  
• Labels  

---------------------------------------- 

#### 3.4.6.17 Performance Requirements 

The Flow Tool SHALL: 

• Use bounded traversal  
• Support progressive loading  
• Maintain responsiveness  
• Prevent unbounded recursive queries  

---------------------------------------- 

#### 3.4.6.18 Backend Integration 

The Flow Tool SHALL: 

• Use Backend API for all operations  
• Execute mutations as ACID transactions  
• Fully rollback on failure  

---------------------------------------- 

#### 3.4.6.19 Test and Verification 

The system SHALL verify: 

• Functional flow relationships are correctly created and validated  
• STPA role assignments enforce constraints  
• Invalid assignments are rejected  
• ControlAction and Feedback nodes behave correctly  
• Flow properties persist correctly  
• All operations are transactional  
---------------------------------------- 

### 3.4.7  The Asset Manager Tool
#### 3.4.7.1 Purpose 
The Asset Manager Tool is an Add-on Tool used to create, inspect, edit, organize, and analyze Assets within the current System of Interest (SoI). 

The Asset Manager Tool SHALL provide a structured, table-oriented, and methodologically guided interface for managing: 

• (:Asset) nodes  
• (:Regime) nodes  
• (:Loss) nodes  
• Root (:Goal) nodes  
• Asset relationships to Elements, Functions, Interfaces, States, and Environments  

The tool SHALL guide the User in defining Assets and their associated certification structures while preserving User flexibility and avoiding rigid prescriptive workflows. 

The Asset Manager Tool SHALL support both: 

1. efficient expert workflows; and  
2. progressive disclosure for less experienced Users  

The tool SHALL be branded at the top of the pop-up window as “Asset Manager Tool”. 

---------------------------------------- 

#### 3.4.7.2 Core Concepts 

##### Asset Types 

Assets SHALL be classified into two types: 

• PRIMARY  
• DERIVED  

PRIMARY Assets: 

• represent intrinsically valuable entities  
• SHALL define their own Criticality and Assurance needs  

DERIVED Assets: 

• derive their value from enabling compromise of a PRIMARY Asset  
• SHALL reference one or more PRIMARY Assets  
• SHALL inherit Criticality from the referenced PRIMARY Asset(s)  
• MAY define additional Assurance requirements  

Example: 

A cryptographic key used to protect a data Asset is a DERIVED Asset. 

##### Asset Relationships 

Assets MAY be related to: 

• (:Element)  
• (:Function)  
• (:Interface)  
• (:State)  
• (:Environment)  

These relationships SHALL define where the Asset exists, is processed, or is exposed. 

---------------------------------------- 

#### 3.4.7.3 Invocation 

The Asset Manager Tool SHALL be launched from the SSTPA Control Panel. 

If a Data Drawer is open for: 

• (:Asset) → the tool SHALL open focused on that Asset  
• (:System) → the tool SHALL display all Assets in the SoI  
• (:Element), (:Function), (:Interface), (:State), (:Environment) → the tool SHALL filter Assets associated with that node  

If no context exists, the tool SHALL display all Assets in the current SoI. 

Opening the Asset Manager Tool SHALL NOT change the current SoI. 

---------------------------------------- 

#### 3.4.7.4 Asset Table View 

The primary interface SHALL be a table displaying all Assets in the current SoI. 

Each row SHALL represent one (:Asset). 

Columns SHALL include: 

• HID  
• Name  
• Asset Type (PRIMARY / DERIVED)  
• Criticality (multi-value)  
• Assurance (multi-value)  
• Associated Regimes  
• Associated Elements  
• Associated Functions  
• Associated Interfaces  
• Associated States  
• Associated Environments  
• Derived From (if DERIVED)  
• Number of Loss nodes  
• Goal Structure status  
• Validation status  

The table SHALL support: 

• sorting  
• filtering  
• column visibility control  
• multi-select  
• inline editing (where permitted)  
• search (HID, Name, description)  

---------------------------------------- 

#### 3.4.7.5 Progressive Disclosure UX 

The Asset Manager Tool SHALL use progressive disclosure to manage complexity. 

The UI SHALL support expandable panels per Asset: 

Level 1 (collapsed): 

• Summary row (table view) 

Level 2 (expanded row): 

• Core properties  
• Criticality and Assurance selection  
• Regime selection  
• Asset relationships (high-level) 

Level 3 (detail panel or modal): 

• Full property editing  
• Relationship editing  
• Loss configuration  
• Goal access  

The tool SHALL NOT require Users to complete all fields before creating an Asset. 

The tool SHALL guide but SHALL NOT enforce strict sequencing. 

---------------------------------------- 

#### 3.4.7.6 Asset Creation 

The tool SHALL allow creation of new (:Asset) nodes. 

On creation: 

• HID and uuid SHALL be generated  
• Asset Type SHALL be selected (PRIMARY or DERIVED)  
• default properties SHALL be initialized  
• Owner, Creator, Created, and LastTouch SHALL be assigned  

The tool SHALL prompt the User to: 

1. define Criticality values  
2. define Assurance values  
3. assign or create Regimes  
4. associate Environments  

---------------------------------------- 

#### 3.4.7.7 Automatic Node Generation 

For each Asset, the tool SHALL automatically create: 

For each combination of: 

• Criticality  
• Assurance  
• Environment  

The tool SHALL create: 

• one (:Loss) node  
• one Root (:Goal) node  

The tool SHALL associate: 

(:Asset)-[:HAS_LOSS]->(:Loss)  
(:Asset)-[:HAS_GOAL]->(:Goal)  

Each (:Loss) SHALL: 

• reference the Asset  
• reference one Environment  
• define one Criticality  
• define one Assurance  

Each Root Goal SHALL: 

• be associated with the Asset-Loss pair  
• be initialized with a default GoalStatement  

---------------------------------------- 

#### 3.4.7.8 Regime Management 

##### Regime Concept 

A (:Regime) represents a certification authority or governing standard. 

Each Asset SHALL have one or more Regimes per Criticality. 

Regimes MAY differ between: 

• PRIMARY Assets  
• DERIVED Assets  
• different Criticalities  

##### Master Regime Node 

The system SHALL support a reusable master Regime node: 

(:MasterRegime) 

The Master Regime SHALL serve as a template. 

Properties SHALL include: 

• Name  
• Authority  
• Standard  
• Description  
• Certification Scope  
• Metadata fields  

##### Regime Cloning 

The Asset Manager Tool SHALL allow Users to: 

• select a Master Regime  
• clone it into an Asset-specific (:Regime) node  

Cloning SHALL: 

• copy all properties  
• generate new HID and uuid  
• associate the new Regime with the Asset  

Relationship: 

(:Asset)-[:HAS_REGIME]->(:Regime) 

The tool SHALL support: 

• editing cloned Regimes  
• creating new Master Regimes  
• reusing Master Regimes across the SoI  

##### UX Requirements 

The tool SHALL provide: 

• Regime selection dropdown  
• searchable Master Regime list  
• “Clone Regime” action  
• “Create New Regime” action  
• inline Regime editing  

---------------------------------------- 

#### 3.4.7.9 Asset Relationship Allocation 

The tool SHALL allow allocation of Assets to: 

• Elements  
• Functions  
• Interfaces  
• States  
• Environments  

The tool SHALL provide: 

• multi-select pickers  
• graph-assisted selection (optional)  
• filtered lists based on SoI  

The tool SHALL validate: 

• same SoI membership  
• valid relationship types  

The tool SHALL allow batch allocation. 

---------------------------------------- 

#### 3.4.7.10 Derived Asset Handling 

For DERIVED Assets: 

The tool SHALL require association to at least one PRIMARY Asset. 

Relationship: 

(:Asset)-[:DERIVED_FROM]->(:Asset) 

Constraints: 

• target SHALL be PRIMARY  
• DERIVED Asset SHALL inherit Criticality  
• tool SHALL visually indicate inherited Criticality  

The tool SHALL allow additional Assurance values. 

---------------------------------------- 

#### 3.4.7.11 Loss Editing Integration 

The tool SHALL allow editing of Loss nodes. 

The tool SHALL display: 

• Loss HID  
• Criticality  
• Assurance  
• Environment  
• associated Goal  

The tool SHALL allow: 

• opening Loss in Loss Tool  
• inline editing of allowed properties  
• regeneration of missing Loss nodes  

---------------------------------------- 

#### 3.4.7.12 Goal Integration 

The tool SHALL allow access to Goal Structures. 

The tool SHALL display: 

• Root Goal status  
• completeness indicator  
• evidence indicator  

The tool SHALL allow: 

• opening Goal Keeper Tool  
• navigating to Root Goal  

---------------------------------------- 

#### 3.4.7.13 Modes of Operation 

The tool SHALL support: 

a. Table View  
• asset overview  
• sorting/filtering  

b. Detail View  
• full Asset editing  
• relationship editing  

c. Regime View  
• Master Regime management  
• cloning and editing  

d. Validation View  
• missing relationships  
• incomplete definitions  
• inconsistent Criticality/Assurance  

---------------------------------------- 

#### 3.4.7.14 Validation Requirements 

The tool SHALL validate: 

• DERIVED Assets reference PRIMARY Assets  
• Loss nodes exist for each Criticality/Assurance/Environment  
• Root Goal exists for each Loss  
• Regime exists for each Criticality  
• relationships are valid and within SoI  

The tool SHALL provide: 

• warnings (non-blocking)  
• errors (blocking)  

---------------------------------------- 

#### 3.4.7.15 Interaction Requirements 

The tool SHALL support: 

• inline editing  
• batch operations  
• undo before commit  
• search  
• filtering  
• hover highlighting  
• keyboard navigation  
• commit confirmation  

---------------------------------------- 

#### 3.4.7.16 Performance Requirements 

The tool SHALL: 

• support large Asset sets  
• use pagination  
• support lazy loading  
• maintain UI responsiveness  

---------------------------------------- 

#### 3.4.7.17 Backend Integration 

The tool SHALL: 

• retrieve Assets for SoI  
• retrieve Regimes and Master Regimes  
• create and update Assets  
• create Loss and Goal nodes automatically  
• validate all changes before commit  
• persist all changes transactionally  

---------------------------------------- 

#### 3.4.7.18 UX Design Principles 

The Asset Manager Tool SHALL: 

• guide without enforcing rigid workflows  
• minimize cognitive load  
• expose complexity progressively  
• support expert speed workflows  
• maintain consistency with other Add-on Tools  

The tool SHALL prioritize: 

• clarity of relationships  
• ease of navigation  
• minimal data duplication  
• fast editing cycles  

---------------------------------------- 

#### 3.4.7.19 Test and Verification Requirements 

The system SHALL verify: 

• Asset creation works for PRIMARY and DERIVED  
• Loss and Goal nodes are automatically created  
• Regimes can be cloned from Master Regime  
• DERIVED Assets inherit Criticality  
• relationships are valid  
• UI interactions persist correctly  
• validation rules trigger correctly  
• transactions commit atomically  


### 3.4.8  The Environment Manager Tool
#### 3.4.8.1 Purpose 
The Environment Manager Tool is intended to be the primary method where Environments and their relationships are entered.  It identifies and characterize all Environments in the SoI in one place and allow Users to edit their properties and relationships.  Information will be presented in a mix of table and graphic format.  When data is presented in a table Asset manager will maintain the underlying graphical structure. 

Environment Manager SHALL:
Display summary of all Environments currently associated with the SoI with Environment Name and short description.  
It will Identify associated (:State) and (:Hazard) relationships using the progressive disclosure pattern for editing their properties in a table
On selection of Environment, display graphic relationships of Environment to SoI (:State), and (:Hazard) and allow user to add or remove Relationships
Allow User to create New (:Hazard) nodes from Reference Data by allowing User to review and select (clone and own pattern) (:ThreatActor) nodes from the Reference Data Set.
Allow User to create New Environment and auto generates, on commit the associated (:Loss) and (:Goal) nodes (Note:  (:Asset) SHALL have one (:Loss) and one(:Goal) node for each Criticality, Assurance, Environment)


### 3.4.9  The State-Trace Tool
#### 3.4.9.1 Purpose 
The State Trace Tool is intended to be the primary method where (:Assets) are assigned to (:State).  In a table based format each row will be a (:State) and each column will be an (:Element), (:Function) or (:Interface) related to a specific (:Asset).  If there are more than one (:Asset) in the SoI, the tool will allow the User to select the (:Asset) and generate the appropriate Table.  Cells in teh Tble will be clickable and on-click, the Asset will be related to the state.  Further, on commit the structured JSON file associated with each associated (:Loss) will be modified to record the (:Asset) resides in the { (:Element), (:Function), or (:Interface) } in that (:State). 

State-Trace Tool SHALL:
Display relationship between (:State) and (:Element), (:Function, (:Interface) for each (:Asset) in a table format allowing user so to make assignments
Update each associated (:Loss) node's Structured JSON property record the the status to be displayed in teh Loss Tool.  



### 3.4.10  The Loss Tool
#### 3.4.10.1 Purpose 
In the body of evidence needed for certification, nothing is more important than the analysis of Loss and the identification and approval of Residual Vulnerabilities.  All Assets will have Residual Vulnerabilities, the purpose of the Loss Tool is to develop a Loss View which is both comprehensible and acceptable to system stakeholders and in particular certification authorities.

In SSTPA Tools, a (:Loss) Node will be created from the commit of an (:Asset) Node.  One (:Loss) node will be created for every Asset Criticality and Assurance pair.  The User uses the Loss Tool to allocate Loss to (:Environment) as their should be one (:Loss) for every Criticality Assurance and Environment combination but there are many Environments where a Loss is not relevant and this is the User's call to make.

The Attack Tree is a Directed Analytical Graph (DAG) developed from the cyclical relationships derived from (:Loss).  The Attack Tree "un-rolls the relationship structure into linear form and assigns Sequential AND (SAND) and Exclusive OR (XOR) relationships not available in the graph model.  The Loss Tool will auto-generates an Attack Tree diagram based on the conventions of Structured Attack Tree Analysis and allows the User to extend it, modify it and add Countermeasures and Attacks.  The User must terminate each branch with an Attack (which will be a Residual Vulnerability) or a new (:Asset) node with derived criticality and Assurances which will spawn new (:Loss) nodes as children under the new derived Assets.


Loss Tool SHALL:
Display the Directed Analytical Graph of the one specific (:Loss) as an Attack Tree following the rules of Structured Attack Tree Analysis from a structured JSON Property of (:Loss)
It will allow the User to relate (:Attack) nodes to (:Element), (:Function), and (:Interface) nodes in the attack tree diagram
Allow User to select or create New (:Attack) nodes from Reference Data by allowing User to review and select (clone and own pattern) (:Tactics) or (:Techniques) from the MITRE ATT&CK reference data set or a (:TID) from the MITRE EMB3D Reference Data Set. 
Allow Users to relate (:Countermeasure) nodes to (:Attack) nodes to [:BLOCK] them
Allow User to select or create New (:Countermeasure) nodes from Reference Data by allowing User to review and select (clone and own pattern) relevent nodes from (:Mitigation) the MITRE ATT&CK reference data set or (:MID) from the MITRE EMB3D Reference Data Set.
Allow User to create new (:Asset) nodes with the property of "Derived) and relate them to (:Countermeasure) as a means of terminating the tree.  On commit, introductio nof a new (:Asset) will auto-generate new (:Loss) and (:Goal) nodes.

#### 3.4.10.1 Invocation 

The Loss Tool SHALL: 

• Be launched from the SSTPA Control Panel  
• If a Data Drawer is open for (:Loss) and there is a valid AttackTreeJSON property than the Loss Tool SHALL display the serialized JSON document in its canvas.
• If a Data Drawer is open for (:Loss) and there is not a valid AttackTreeJSON property than the Loss Tool SHALL create a serialized JSON document and display it in its canvas.
• If the GUI is in any other state, the Loss Tool SHALL present the User with (:Loss) nodes from the SoI and their short description and allow the user operate on one which SHALL be processed as above.


------------------------
### 3.4.11  The Goal Keeper Tool
#### 3.4.11.1 Purpose 

#### 3.4.11.1 Purpose 

The Goal Keeper Tool is an Add-on Tool used to create, display, edit, validate, persist, and export a formal certification argument for an Asset Loss using Goal Structuring Notation (GSN) Community Standard Version 3. 
The Goal Keeper Tool SHALL represent the certification argument as a Directed Acyclic Graph (DAG) rooted at a single (:Goal) node associated with an (:Asset) and its corresponding (:Loss). 
The Root Goal SHALL be automatically created when the Frontend creates the associated (:Asset) and (:Loss) nodes. 
The Goal Keeper Tool SHALL allow the User to construct an assurance case by creating, associating, arranging, editing, and saving GSN nodes and relationships according to GSN rules. 
The GSN DAG SHALL terminate in (:Solution) nodes. (:Solution) nodes SHALL reference evidence-bearing SSTPA nodes, including: 

• (:Validation) 
• (:Verification) 
• (:Loss) 

The Goal Keeper Tool SHALL preserve both: 
1. the authoritative semantic graph stored in the Backend as nodes and relationships; and 
2. the user-manipulated visual layout stored as structured diagram JSON. 

The tool described here SHALL be branded at the top of the pop-up window as “Goal Keeper Tool”. 

---------------------------------------- 

#### 3.4.11.2 Invocation 

The Goal Keeper Tool SHALL be launched from the SSTPA Control Panel. 
If a Data Drawer is open for a (:Goal), the Goal Keeper Tool SHALL open the Goal Structure rooted at that (:Goal). 
If a Data Drawer is open for an (:Asset), the Goal Keeper Tool SHALL display the Goal Structures associated with that Asset and allow the User to open one. 
If a Data Drawer is open for a (:Loss), the Goal Keeper Tool SHALL open the Goal Structure associated with that Loss, if one exists. 
If no valid Data Drawer context exists, the Goal Keeper Tool SHALL present the User with a list of available (:Asset), (:Loss), and Root (:Goal) combinations in the current SoI. 

The selection list SHALL display: 

• Asset HID 
• Asset Name 
• Loss HID 
• Loss Name 
• Root Goal HID 
• Root Goal Name 
• Criticality 
• Assurance 
• Environment 
• Goal Structure status, if available 

Selecting a Goal Structure SHALL open the associated GSN DAG. 
Opening the Goal Keeper Tool SHALL NOT change the current SoI. 

---------------------------------------- 

#### 3.4.11.3 Supported Node Context 

The Goal Keeper Tool SHALL support invocation when the Data Drawer is open for: 

• (:Asset) 
• (:Loss) 
• (:Goal) 
• (:Strategy) 
• (:Context) 
• (:Justification) 
• (:Assumption) 
• (:Solution) 
• (:Validation) 
• (:Verification) 

The tool SHALL load all GSN nodes reachable from the selected Root Goal using valid GSN relationships. 
The tool SHALL load referenced evidence nodes needed to display terminal Solution evidence. 
The tool SHALL allow the User to shift focus within the Goal Structure without changing the current SoI. 
The tool SHALL provide a back action to return to the invoking context. 

---------------------------------------- 

#### 3.4.11.4 GSN Node Types 

The Goal Keeper Tool SHALL support the following GSN node types: 

• (:Goal) 
• (:Strategy) 
• (:Context) 
• (:Justification) 
• (:Assumption) 
• (:Solution) 

The Root Goal SHALL be a (:Goal) node. 
The Root Goal SHALL represent the top-level certification claim for a specific Asset-Loss case. 
(:Goal) nodes SHALL represent claims. 
(:Strategy) nodes SHALL represent reasoning or inference patterns used to decompose or support Goals. 
(:Context) nodes SHALL represent contextual information needed to interpret a Goal or Strategy. 
(:Justification) nodes SHALL represent rationale supporting a Goal or Strategy. 
(:Assumption) nodes SHALL represent assumptions relied upon by a Goal or Strategy. 
(:Solution) nodes SHALL represent references to evidence supporting a Goal. 

---------------------------------------- 

#### 3.4.11.5 GSN Relationship Types 

The Goal Keeper Tool SHALL support the following relationships: 

(:Goal)-[:SUPPORTED_BY]->(:Goal) 
(:Goal)-[:SUPPORTED_BY]->(:Strategy) 
(:Goal)-[:SUPPORTED_BY]->(:Solution) 

(:Goal)-[:IN_CONTEXT_OF]->(:Context) 
(:Goal)-[:IN_CONTEXT_OF]->(:Justification) 
(:Goal)-[:IN_CONTEXT_OF]->(:Assumption) 

(:Strategy)-[:SUPPORTED_BY]->(:Goal) 
(:Strategy)-[:SUPPORTED_BY]->(:Solution) 

(:Strategy)-[:IN_CONTEXT_OF]->(:Context) 
(:Strategy)-[:IN_CONTEXT_OF]->(:Justification) 
(:Strategy)-[:IN_CONTEXT_OF]->(:Assumption) 

(:Context)-[:HAS_ENVIRONMENT]->(:Environment) 

(:Solution)-[:HAS_VALIDATION]->(:Validation) 
(:Solution)-[:HAS_VERIFICATION]->(:Verification) 
(:Solution)-[:HAS_LOSS]->(:Loss) 

The Backend SHALL validate all Goal Keeper relationships before commit. 
The Goal Keeper DAG SHALL NOT contain cycles. 
The Backend SHALL reject any relationship that would create a cycle in the Goal Structure. 
The Backend SHALL prevent duplicate logical GSN relationships between the same source node, target node, and relationship type. 

---------------------------------------- 

#### 3.4.11.6 GSN Structure Rules 

A Goal Structure SHALL have exactly one Root Goal. 
The Root Goal SHALL have no incoming [:SUPPORTED_BY] relationship from another GSN node within the same Goal Structure. 
Every non-root GSN node SHALL be reachable from the Root Goal. 
A (:Goal) node MAY be supported by one or more (:Goal), (:Strategy), or (:Solution) nodes. 
A (:Strategy) node SHALL be used to explain how a Goal is decomposed into supporting Goals or Solutions. 
A (:Strategy) node SHOULD have at least one outgoing [:SUPPORTED_BY] relationship to a (:Goal) or (:Solution). 
A (:Solution) node SHALL be terminal with respect to [:SUPPORTED_BY] relationships. 
A (:Solution) node SHALL NOT have outgoing [:SUPPORTED_BY] relationships. 
A (:Solution) node SHALL reference at least one evidence-bearing node before the Goal Structure can be marked complete. 

Valid evidence-bearing references for (:Solution) are: 

• (:Validation) 
• (:Verification) 
• (:Loss) 

(:Context), (:Justification), and (:Assumption) nodes SHALL NOT support Goals directly through [:SUPPORTED_BY]. 
(:Context), (:Justification), and (:Assumption) nodes SHALL only be related through [:IN_CONTEXT_OF] or explicitly authorized context relationships. 

The Goal Keeper Tool SHALL visually identify incomplete or invalid GSN structures. 

---------------------------------------- 

#### 3.4.11.7 Modes of Operation 

The Goal Keeper Tool SHALL support the following modes: 

a. Goal Structure View 

• Displays the full GSN DAG rooted at the selected Root Goal 
• Allows selection of GSN nodes and relationships 
• Allows creation and association of valid GSN nodes 
• Allows editing of GSN node properties 
• Allows association of evidence nodes to Solution nodes 
• Allows manual diagram layout and visual organization 

b. Evidence View 

• Displays terminal (:Solution) nodes and their referenced evidence 
• Allows inspection of linked (:Validation), (:Verification), and (:Loss) nodes 
• Highlights unsupported Goals and Solutions without evidence 
• Allows the User to navigate from evidence to the supporting Goal path 

c. Validation View 

• Displays structural completeness and rule violations 
• Identifies unsupported Goals 
• Identifies terminal Goals without Solutions 
• Identifies Solutions without evidence 
• Identifies unreachable GSN nodes 
• Identifies invalid relationship types 
• Identifies cycle attempts or detected cycles 
• Provides actionable remediation messages 

d. Presentation / Export View 

• Provides a clean report-oriented rendering of the Goal Structure 
• Preserves GSN shapes, labels, relationships, and evidence references 
• Allows export of the current viewport or full Goal Structure 

---------------------------------------- 

#### 3.4.11.8 Visualization Requirements 

The Goal Keeper Tool SHALL render diagrams consistent with GSN Community Standard Version 3 to the maximum extent practical within the SSTPA Tool visual style. 

The diagram SHALL use the following visual conventions: 

• (:Goal) nodes SHALL be displayed as rectangles. 
• (:Strategy) nodes SHALL be displayed as parallelograms. 
• (:Solution) nodes SHALL be displayed as circles. 
• (:Context) nodes SHALL be displayed as rectangles with rounded sides. 
• (:Justification) nodes SHALL be displayed as ovals. 
• (:Assumption) nodes SHALL be displayed as ovals. 

All GSN nodes SHALL display: 

• GSN ID 
• Node type label 
• HID 
• Name 
• Statement property 

Statement properties SHALL include: 

• GoalStatement for (:Goal) 
• StrategyStatement for (:Strategy) 
• ContextStatement for (:Context) 
• JustificationStatement for (:Justification) 
• AssumptionStatement for (:Assumption) 
• SolutionStatement for (:Solution) 

[:SUPPORTED_BY] relationships SHALL be displayed as directed arrows with solid arrowheads. 
[:IN_CONTEXT_OF] relationships SHALL be displayed as directed arrows with hollow arrowheads. 

Evidence relationships from (:Solution) nodes SHALL be displayed as evidence references attached to the Solution node or shown in an adjacent evidence list box. 
Evidence nodes SHALL NOT be rendered inside the Solution circle. 
The tool SHALL visually distinguish: 

• Root Goal 
• Selected node 
• Hover state 
• Invalid node 
• Incomplete node 
• Terminal Solution 
• Evidence-linked Solution 
• Evidence-missing Solution 

---------------------------------------- 

#### 3.4.11.9 Diagram Persistence Requirements 

The Goal Keeper Tool SHALL persist the visual state of the Goal Structure as structured JSON. 
The persisted diagram JSON SHALL store presentation and reconstruction data only. 
The persisted diagram JSON SHALL NOT be authoritative for semantic GSN relationships. 
The authoritative Goal Structure SHALL be the Neo4j graph of GSN nodes and relationships. 
The persisted diagram JSON SHALL include: 

• schema version 
• Root Goal HID 
• Root Goal uuid 
• tool type 
• viewport 
• zoom level 
• node positions 
• edge routing 
• collapsed/expanded state 
• display toggles 
• layout mode 
• selected visual options 
• evidence panel display state 

The Goal Keeper Tool SHOULD use the common Diagram View Persistence Model defined for graphical Add-on Tools. 
If the SRS does not yet define a common (:DiagramView) node, the Goal Keeper Tool MAY persist diagram JSON in the (:Goal).GoalStructure property for MVP implementation. 
Future versions SHOULD migrate (:Goal).GoalStructure into a common (:DiagramView) node related to the Root Goal. 
On opening an existing Goal Structure, the Frontend SHALL: 

1. retrieve the authoritative GSN graph from the Backend; 
2. retrieve the persisted diagram JSON, if present; 
3. reconcile saved diagram references against existing graph nodes and relationships; 
4. restore valid node positions, edge routing, viewport, and display options; 
5. ignore stale references to deleted nodes or relationships; 
6. visually notify the User if stale references were ignored. 

On Commit, semantic graph changes and diagram JSON changes SHALL be committed transactionally. 
If semantic graph persistence succeeds but diagram JSON persistence fails, the entire transaction SHALL roll back unless the User explicitly elects to save semantic changes without layout persistence. 

---------------------------------------- 

#### 3.4.11.10 Node Creation and Editing 

The Goal Keeper Tool SHALL allow creation of the following nodes: 

• (:Goal) 
• (:Strategy) 
• (:Context) 
• (:Justification) 
• (:Assumption) 
• (:Solution) 

Creation SHALL: 

• use the standard SSTPA staged editing and Commit confirmation model; 
• assign valid HID and uuid values; 
• assign the node to the active SoI; 
• assign valid GSN ID values; 
• assign Owner, Creator, Created, and LastTouch properties according to SRS ownership rules; 
• open the created node in the Goal Keeper detail panel or Data Drawer for editing. 

The Goal Keeper Tool SHALL allow editing of: 

• Name 
• ShortDescription 
• LongDescription 
• GoalStatement 
• StrategyStatement 
• ContextStatement 
• JustificationStatement 
• AssumptionStatement 
• SolutionStatement 

The Goal Keeper Tool SHALL NOT allow editing of fixed identity properties except as allowed by Admin rules. 

---------------------------------------- 

#### 3.4.11.11 Relationship Creation and Editing 

The Goal Keeper Tool SHALL allow the User to create valid GSN relationships by dragging, selecting, or using explicit Add Relationship actions. 
The tool SHALL visually distinguish valid and invalid relationship targets. 
Invalid relationship targets SHALL be muted or disabled. 
The tool SHALL prevent Commit until all staged relationships pass Backend validation. 

Relationship creation SHALL validate: 

• source node type 
• target node type 
• relationship type 
• same SoI constraint 
• Root Goal reachability 
• DAG acyclicity 
• duplicate relationship prevention 
• GSN rule compliance 

The tool SHALL allow deletion of GSN relationships subject to orphan and structural validation rules. 

If deleting a relationship causes one or more GSN nodes to become unreachable from the Root Goal, the tool SHALL warn the User before Commit. 

---------------------------------------- 

#### 3.4.11.12 Solution Evidence Association 

The Goal Keeper Tool SHALL allow (:Solution) nodes to reference evidence-bearing SSTPA nodes. 

The tool SHALL allow association of: 

• (:Solution)-[:HAS_VALIDATION]->(:Validation) 
• (:Solution)-[:HAS_VERIFICATION]->(:Verification) 
• (:Solution)-[:HAS_LOSS]->(:Loss) 

The tool SHALL allow multiple evidence references per (:Solution). 
The tool SHALL display evidence references in a list box attached to or visually associated with the Solution node. 

The evidence list SHALL display: 

• Evidence node type 
• HID 
• Name 
• ShortDescription 
• Verification or Validation method, where applicable 
• Loss Criticality and Assurance, where applicable 

The tool SHALL allow opening an evidence node in read-only inspection mode when the evidence node is outside the current editable context. 
The tool SHALL allow opening an evidence node for edit when it belongs to the current SoI and the User has edit authority. 
A (:Solution) node without at least one evidence reference SHALL be visually marked incomplete. 

---------------------------------------- 

#### 3.4.11.13 Asset, Loss, and Root Goal Integration 

When a new (:Asset) and associated (:Loss) are created, the Frontend SHALL automatically create the corresponding Root (:Goal). 
There SHALL be one Root (:Goal) for each certification argument associated with an Asset-Loss case. 
The Root Goal SHOULD be initialized with a default GoalStatement derived from: 

• Asset Name 
• Asset HID 
• Loss Name 
• Loss HID 
• Criticality 
• Assurance 
• Environment 

Example default Root Goal statement: 

“The evidence supports certification that [Asset] maintains [Assurance] for [Criticality] in [Environment] such that [Loss] is acceptably mitigated.” 
The User SHALL be able to edit the Root Goal statement. 
The Goal Keeper Tool SHALL preserve the association among: 

• (:Asset) 
• (:Loss) 
• Root (:Goal) 

The tool SHALL NOT allow deletion of the Root Goal without explicit warning that the certification argument for the Asset-Loss case will be removed. 

---------------------------------------- 

#### 3.4.11.14 Validation Requirements 

The Goal Keeper Tool SHALL validate all proposed mutations through the Backend API before Commit. 

Validation SHALL confirm: 

• exactly one Root Goal exists for the opened Goal Structure; 
• the Goal Structure is a DAG; 
• all non-root GSN nodes are reachable from the Root Goal; 
• all relationship types are valid for their source and target node types; 
• no duplicate logical relationships exist; 
• all Solution nodes terminate the supported-by structure; 
• Solution nodes do not have outgoing [:SUPPORTED_BY] relationships; 
• Solution nodes marked complete have at least one evidence reference; 
• GSN IDs are unique within the Goal Structure; 
• all nodes belong to the same SoI unless explicitly allowed by evidence-reference rules; 
• diagram JSON is well formed and compatible with the supported schema version. 

The API SHALL return: 

• valid / invalid 
• reason for invalidity 
• affected node HID or relationship 
• recommended corrective action, where practical 

---------------------------------------- 

#### 3.4.11.15 Interaction Requirements 

The Goal Keeper diagram SHALL support: 

• zoom 
• pan 
• node selection 
• relationship selection 
• hover highlighting 
• drag-to-position nodes 
• drag-to-create relationships where practical 
• animated centering 
• expand/collapse of branches 
• expand/collapse of evidence lists 
• keyboard navigation 
• Escape to close 
• undo of staged, uncommitted diagram operations 
• Commit confirmation 

Selecting a node SHALL: 

• highlight the node; 
• display its properties; 
• display incoming and outgoing relationships; 
• display validation state; 
• display available actions for that node type. 

Selecting a relationship SHALL: 

• highlight the relationship; 
• display source node, target node, and relationship type; 
• display validation state; 
• allow deletion if permitted. 

---------------------------------------- 

#### 3.4.11.16 Layout Requirements 

The Goal Keeper Tool SHALL support: 

• manual layout 
• hierarchical top-down layout 
• hierarchical left-to-right layout 
• automatic layout on demand 

The layout engine SHALL preserve relative node positioning where feasible. 

Switching layout modes SHALL NOT discard User-created manual positioning unless the User confirms re-layout. 

The tool SHALL maintain layout stability during: 

• selection 
• editing 
• validation 
• evidence association 
• branch expansion 
• mode switching 
• save and reopen 

---------------------------------------- 

#### 3.4.11.17 Search and Navigation 

The Goal Keeper Tool SHALL provide search within the open Goal Structure. 

Search SHALL support: 

• HID 
• uuid 
• GSN ID 
• Name 
• node type 
• statement text 
• evidence HID 
• evidence node type 

Search results SHALL: 

• be listed in a synchronized results panel; 
• highlight matching nodes in the diagram; 
• allow centering on a selected result; 
• prioritize exact HID, uuid, and GSN ID matches. 

The tool SHALL provide a path-to-root display for the selected node. 

The path-to-root display SHALL show the chain from Root Goal to selected node. 

---------------------------------------- 

#### 3.4.11.18 Data Drawer Integration 

On successful selection from the Goal Keeper Tool, the calling Data Drawer SHALL be able to display: 

• related GSN nodes 
• GSN relationships 
• evidence references 
• validation status 
• Goal Structure status 

The Data Drawer SHALL allow: 

• launching the Goal Keeper Tool from valid node context; 
• opening related GSN nodes; 
• opening referenced evidence nodes; 
• removing valid relationships subject to validation and deletion rules; 
• refreshing after commit. 

The Main Panel SHALL refresh after the Goal Keeper Tool commits changes. 

---------------------------------------- 

#### 3.4.11.19 Export Requirements 

The Goal Keeper Tool SHALL support export of Goal Structures. 

Supported export formats SHALL include: 

• PNG 
• SVG 
• JSON 
• Markdown summary 

The PNG and SVG exports SHALL preserve: 

• GSN node shapes 
• node labels 
• statement text 
• relationship direction 
• relationship style 
• evidence list boxes 
• visible validation markers, if enabled 

The JSON export SHALL include: 

• diagram schema version 
• Root Goal HID and uuid 
• GSN nodes 
• GSN relationships 
• evidence references 
• layout information 
• viewport and display settings 

The User SHALL be able to export: 

• current viewport 
• full visible Goal Structure 
• evidence summary 
• validation findings 

Exports SHALL be suitable for insertion into System Description, System Specification, certification package, and body-of-evidence reports. 

---------------------------------------- 

#### 3.4.11.20 Performance Requirements 

The Goal Keeper Tool SHALL: 

• efficiently load Goal Structures for the active SoI; 
• support progressive loading for large Goal Structures; 
• maintain UI responsiveness during layout, editing, validation, and export; 
• use bounded traversal for all recursive GSN queries; 
• prevent unbounded recursive graph expansion; 
• support lazy loading of evidence-node detail. 

Exact HID, uuid, and GSN ID lookup SHALL be faster than general text search. 

---------------------------------------- 

#### 3.4.11.21 Backend Integration Requirements 

The Goal Keeper Tool SHALL retrieve and mutate data through the Backend API. 

Required Backend capabilities include: 

• retrieval of Root Goal by Asset and Loss; 
• retrieval of complete GSN DAG by Root Goal; 
• retrieval of GSN node properties; 
• retrieval of GSN relationships; 
• retrieval of Solution evidence references; 
• validation of GSN relationship creation; 
• validation of DAG acyclicity; 
• validation of Solution evidence completeness; 
• transactional creation, update, and deletion of permitted GSN nodes; 
• transactional creation and deletion of permitted GSN relationships; 
• transactional persistence of diagram JSON; 
• export data retrieval. 

All Goal Keeper write operations SHALL be ACID compliant. 
All semantic graph mutations and diagram persistence updates SHALL commit as a single transaction unless explicitly separated by User confirmation. 

---------------------------------------- 

#### 3.4.11.22 Error Handling 

If validation fails, the Goal Keeper Tool SHALL display: 

• failure status; 
• specific rule violated; 
• affected node or relationship; 
• recommended corrective action, where practical. 

No partial Goal Structure mutation SHALL be committed. 

If diagram JSON cannot be reconciled with the current semantic graph, the tool SHALL: 

• open the valid portion of the Goal Structure; 
• ignore stale layout references; 
• notify the User; 
• allow the User to save the repaired layout. 

If a referenced evidence node has been deleted, the tool SHALL mark the associated Solution as incomplete. 

---------------------------------------- 

#### 3.4.11.23 Reporting and Certification Package Support 

The Goal Keeper Tool SHALL support certification-package workflows. 

The tool SHALL allow Users to: 

• frame and export portions of the Goal Structure; 
• generate a complete evidence summary; 
• identify unsupported claims; 
• identify claims supported only by assumptions; 
• identify Solutions lacking evidence; 
• identify referenced Loss, Verification, and Validation nodes; 
• preserve reproducible diagram states for reports. 

The Goal Keeper Tool SHALL support generation of figures and summaries suitable for future certification argument reports. 

---------------------------------------- 

#### 3.4.11.24 Test and Verification Requirements 

The Goal Keeper Tool SHALL be verified through test and analysis. 

The system SHALL verify that: 

• Root Goals are automatically created with Asset-Loss creation; 
• a Goal Structure opens from Asset, Loss, or Goal context; 
• GSN nodes are created with valid HID, uuid, and GSN ID values; 
• valid GSN relationships are accepted; 
• invalid GSN relationships are rejected; 
• cycle creation is rejected; 
• duplicate logical relationships are rejected; 
• non-root unreachable nodes are detected; 
• Solution nodes can reference Validation, Verification, and Loss nodes; 
• Solution nodes without evidence are marked incomplete; 
• diagram JSON persists and reloads correctly; 
• stale diagram references are handled gracefully; 
• exports preserve visible GSN structure; 
• all write operations are transactional and roll back on failure. 


### 3.4.12 Message Center 

#### 3.4.12.1 Purpose 
The Message Center provides the current user access to direct messages and owner-change notification messages. 
The Message Center allows the current user to send direct messages to other users and the Admin.

##### 3.4.12.2 Window behavior 
The Message Center SHALL open in a pop-up window. 
The pop-up SHALL be closable without affecting staged edits in the Data Drawer. 
The pop-up SHALL support refresh. 
The pop-up SHALL preserve the current SoI and current GUI navigation state. 


##### 3.4.12.3 Mailbox list view 
The message list SHALL display columns: 
Subject 
DateTime 
HID 
Sender 


Additional recommended columns: 
Message Type 
Read/Unread 


The list SHALL support: 
sort by clicking column headers 
reverse sort by repeated click 
row selection 
keyboard navigation 
unread highlighting 


##### 3.4.12.4 Message open behavior 

Clicking the message row or row icon SHALL open the message detail view. 
The detail view SHALL display: 
subject 
sender 
sent datetime 
related HID or HIDs 
full message body 

The detail view SHALL support: 
Reply 
Delete 
Close 


##### 3.4.12.5 Direct messaging 
Users SHALL be able to send direct messages to other users. 
Direct messages SHALL be stored in the recipient mailbox. 
Direct messages MAY optionally reference one or more HIDs. 

##### 3.4.12.6 Change notification messages 
The system SHALL generate change notification messages automatically on commit when required by ownership rules. 
Sender SHALL be the current user who committed the change. 
Recipient SHALL be the Owner of the affected node. 
The HID column SHALL show the primary affected HID; where multiple HIDs are affected, the detail view SHALL show the full list. 


##### 3.4.12.7 Delete behavior 
Delete in the current version SHOULD be soft delete. 
Deleting a message from a mailbox SHALL remove it from the current user’s active list only. 
Deleted messages SHOULD remain recoverable for audit unless system retention rules later remove them. 

##### 3.4.12.8 Read state 
Opening a message SHALL mark it read unless the user closes before content load completes. 
Unread count SHALL update after read and delete actions. 



## 3.5 System of Interest Panel
The SSTPA System of Interest Panel SHALL be below the SSTPA Control Panel.
Data presented in the System of Interest Panel SHALL be not editable.  The user will be able to edit data in the associated data drawer.  Data in this panel will be updated when new data from the data drawer is committed.

If there is no current System of Interest (SoI) selected, the panel SHALL display at top center "Select a System of Interest".

The SSTPA Tools System of Interest Panel SHALL display SoI properties: HID, Name and ShortDescription.

The SSTPA Tools System of Interest Panel SHALL display an "Edit" icon which on-click SHALL open all system properties in Data Drawer.  


## 3.6  Main Panel
The intent for the Main Panel is to present enough information on Nodes in the SoI to allow the user to open a Data Drawer to edit them. It will be organized by Primary Node Types each specific Node within the Type grouping can expand to show its Secondary Node relationships with the capability to edit those nodes.  This extends though the ability of each secondary node to expand to expose tertiary nodes.   

The Main Panel SHALL present all data for the currently selected System of Interest (SoI) using a hierarchical, single-window card interface with progressive disclosure.

### 3.6.1 Structure 

The Main Panel SHALL be organized into collapsible Node Type Sections (Primary Types): 
• Environment 
• Connection 
• Interface 
• Function 
• Element 
• Purpose 
• State 
• ControlStructure 
• Asset 
• Security 

When the Security section is expanded, it SHALL disclose related (:Control) and (:Countermeasure) nodes through their relationship groups

Each Node Type Section SHALL: 
•	Display the Node Type name and count
•	Include an Add button for creating a new node of that type
•	Support expand/collapse behavior

### 3.6.2 Primary Entity Cards 

Within each Node Type Section, individual entities SHALL be displayed as cards. 

Each card header SHALL display: 
•	HID
•	Name
•	ShortDescription
•	Node Type badge

Each card SHALL include actions: 
•	Expand / Collapse
•	Edit (opens Data Drawer)
•	Delete (with confirmation)

### 3.6.3 Relationship Groups (Nested) 

When a card is expanded, it SHALL display Relationship Groups corresponding to its outgoing relationships. 

Each Relationship Group SHALL: 
•	Display the relationship name and count
•	Be collapsible
•	Include buttons:
•	Add (create new related node)
•	Associate (link existing node)

For State cards, the TRANSITIONS_TO relationship group SHALL distinguish transitions by TransitionKind so the user can visually identify functional transitions, countermeasure-required transitions, and transitions serving both roles. 



### 3.6.4 Secondary and Tertiary Entity Cards 

Within each Relationship Group, related entities SHALL be displayed as cards. 
•	Secondary entity cards MAY expand to show their own relationship groups (tertiary level)
•	This pattern SHALL recursively apply to all supported hierarchy levels

3.6.5 Repeated Entity Representation 

An entity MAY appear multiple times in the Main Panel when related to multiple parent entities. 

Requirements: 
•	All instances SHALL reference the same underlying node (via HID and uuid)
•	Editing any instance SHALL update the single underlying node
•	Each instance SHALL display its HID to preserve identity clarity

3.6.6 Visual Density Controls 

To manage complexity, the UI SHALL support: 
•	Expand/Collapse at all levels
•	Lazy loading or virtualization for large lists
•	Display of counts for collapsed sections


#### 3.6.7  Node Deletion
Node deletion is complex as the consequence of a user mistake are grave and the potential for unintended cascade is high.  Therefore the following Node Deletion rules are established.

All Node deletions SHALL follow an alert / confirm pattern.

When a Node is deleted other than a (:System) or (:Element) node parenting a (:System) node, the GUI SHALL identify all nodes within the SoI which are orphaned by this action and include this notification in the Alert/Confirm Dialog.  If there are any orphaned nodes, the Alert / Confirm dialog SHALL include a warning:  "WARNING:  Cancel and Re-Associate the Following Nodes or They will be Deleted". 

Node deletion SHALL NOT automatically cascade outside the current SoI.

Deletion of (:System) Nodes SHALL require explicit user confirmation and preview.


## 3.7  The Data Drawer
The right-side Data Drawer SHALL be the single edit surface for the GUI and will implicitly validate node associations.  Note; Add-on Tools may also allow edit and commit but this is in a pop-up window outside the GUI.

All node associations SHALL be validated via Backend API prior to commit.


### 3.7.1 General Behavior 
•	The Data Drawer SHALL slide in from the right side of the GUI but SHALL not obscure or de-activate the Branding Panel or the Control Panel.
•	Only one Data Drawer SHALL be open at a time
•	It should not be possible to open another data drawer while one is already open
•	A Drawer may be exited through "Commit" of its contents or "Close" or "Cancel" without save
•	Opening a new node SHALL replace the current drawer content

### 3.7.2 Header 

The top of the Data Drawer SHALL display: 
•	Node Type
•	Node Name
•	HID

It SHALL include: 
•	Commit button
•	Cancel button
•	Close (X) icon

### 3.7.3 Property Groups 

Properties SHALL be grouped and displayed vertically: 
•	Common Properties
•	Type-specific Properties

Each group SHALL: 
•	Be collapsible (roll up/down)
•	Allow editing of editable fields only

All empty values SHALL be displayed as "Null". 


3.7.4 Relationship Groups in Drawer 

The Data Drawer SHALL also display relationship groups for the selected node. These will be below the Type-specific properties.

Each Relationship Group SHALL: 
•	Show related nodes (HID, Name, ShortDescription)
•	Include actions:
•	Add (create new related node)
•	Associate (link existing node)
•	Remove relationship

For (:State)-[:TRANSITIONS_TO]->(:State) relationships, the Data Drawer SHALL display relationship properties including TransitionKind, Trigger, GuardCondition, Rationale, and any associated Countermeasure HID/uuid traceability fields.



Selecting "Add" a related or "Create" a related node node SHALL open that node in the Data Drawer only after displaying the "Commit" dialog allowing the user to commit the current data in the drawer before closing it.  "Cancel" will have the effect of canceling the create related or add related node action.

When Removing relationship is selected, the effected node SHALL be assessed if relationship removal leaves it an orphan.  If the node is not an orphan, the relationship can be safely removed and the change is made when the Data Drawer is committed.  If the Node will become an orphan when the relationship is removed, it SHALL be treated as a deleted node subject to the same alert/confirm prior to the "Commit".

### 3.7.5 Editing Model 
•	All edits SHALL be staged in the Data Drawer
•	Changes SHALL NOT be persisted until Commit is confirmed
•	Commit SHALL trigger validation and backend update

Data Drawer / commit insertion 

6) Section 3.7 The Data Drawer 

This is the critical functional insertion point because the SRS already says all edits are staged there and commit triggers backend update. 


#### 3.7.5.1 Commit Notification Behavior 
On Commit confirmation, the Frontend SHALL submit the staged delta to the Backend. 

The Frontend SHALL not independently determine final ownership notification recipients. 

The Backend SHALL determine affected owners and generate required messages. 

If the commit will modify nodes not owned by the current user, the confirmation dialog SHOULD display a notice stating that owner notification messages will be generated. 

The commit response SHOULD include a summary: 
nodes changed 
relationships changed 
number of messages generated 
recipients notified 

If owner notification generation fails, the Frontend SHALL display the overall commit failure and no staged changes SHALL be considered committed. 


### 3.7.6 Navigation 

The Data Drawer MAY support: 
•	Navigation between related nodes
•	Breadcrumb display of relationship context

Navigation in SSTPA Tools is through the SoI Navigator.  Navigation within a Data Drawer to related nodes within the SoI is allowed.  When displaying relationships outside the SoI (for (:Interface) and (:Element) nodes) attempts to edit a node outside the SoI SHALL be responded with an Alert dialog with the Message "Navigate to:  " n.HID " to edit".  The HID for the node SHALL be copiable via icon to allow pasting into the SoI Navigator.

### 3.7.7 Consistency Requirement 

All node editing SHALL occur through the Data Drawer. 

Pop-up windows SHALL NOT be used for editing Secondary or Tertiary nodes. 

Pop-ups SHALL be reserved for: 
•	SoI selection
•	Graph visualization tools
•	Specialized add-on tools




# 5  SSTPA Tool Component Copyright

All SSTPA Tool software components SHALL include a copyright statement:

"2025 Nicholas Triska. All rights reserved.
The SSTPA Tools software and all associated modules, binaries, and source code are proprietary intellectual property of Nicholas Triska.  Unauthorized reproduction, modification, or distribution is strictly prohibited.  Licensed copies may be used under specific contractual terms provided by the author."



All SSTPA Data components SHALL include a copyright statement:

"2025 Nicholas Triska.
The SSTPA Tools is proprietary software. However, users retain ownership of data and reports generated during legitimate use of the software, except for embedded proprietary schemas and templates."


# 6 SSTPA Tool Constraints

The SSTPA Tools SHALL operate on an air-gapped Microsoft Windows 11 Enterprise based network with no access to the internet.   

As the SSTPA tool is developed on a Linux based system the SSTPA Tool SHALL also function on a system with the following characteristics: Operating System: Ubuntu Studio 25.04

KDE Plasma Version: 6.3.4

KDE Frameworks Version: 6.12.0

Qt Version: 6.8.3

Kernel Version: 6.14.0-27-generic (64-bit)

Graphics Platform: Wayland

Processors: 28   Intel  Core  i7-14700K

Memory: 31.1 GiB of RAM

Graphics Processor: Intel  Graphics



# 6.1 SSTPA Tool Architecture

the SSTPA Tool architecture SHALL be implemented with minimum complexity.  When integrating a capability, the developer SHALL asses if libraries or functions already existing in the code-base can execute the new capability before introducing a new library or function.



### 6.2  UI Tech Stack

Core stack 
•	Tauri
•	React
•	TypeScript
•	Vite
•	Tailwind CSS

UI behavior 
•	Headless UI or Radix primitives style component approach
•	Framer Motion for expand/collapse and drawer animations

Data/state 
•	Zustand for UI state
•	TanStack Query for backend fetch/mutate/cache

Large-list performance 
•	react-virtual for long nested card sections

Graph visualization 
•	Cytoscape.js for SoI graph selection popup
•	Do not use the graph library for the main editor surface

Recommended Cytoscape ecosystem: 
	•	cytoscape
	•	react-cytoscapejs
	•	layout plugin such as:
	•	cytoscape-fcose for force-directed layouts
	•	optionally dagre-style layout if you want a clearer top-down hierarchy mode
	•	cytoscape-cose-bilkent is also worth considering for stable medium-sized graphs


Optional effects 
•	plain CSS for glass effect

AG Grid only for: 
•	Requirements Traceability Matrix
•	search results
•	report tabulations


# Installer  
The SSTPA Tools SHALL operate on an air-gapped Microsoft Windows 11 Enterprise based network with no access to the internet.  The Installer SHALL execute a window needing expected available components in this architecture and present the User with resource requirements needed for successful installation prior to installation.



