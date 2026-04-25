SSTPA Tool - Software Requirements Specification (SRS) Version: 5.3 Date: April 10, 2026.

Permission is granted to collaborators and contractors working under authorization from Nicholas Triska to use, reproduce, and modify this document for the purpose of developing the SSTPA Tool and its derivatives.

Any distribution or reuse outside this scope requires prior written consent.



# 1.  Introduction

## 1.1  SRS Definitions

Purpose This Software Requirements Specification (SRS) defines the complete functional and non-functional requirements for the Systems Security-Theoretic Process Analysis (SSTPA) Tool version. This document is intended to be the single source of truth for all project stakeholders, including developers, testers, and project managers. It describes the system's features, capabilities, operational environment, and constraints, ensuring a common understanding and guiding the design, implementation, and verification of the software.

The imperative terminology used in this SRS is standard, but defined here for clarity.


-  Statements are line-items or groups of line items (e.g. this list).  Statements are requirements when they contain an imperative listed below.

-  "Shall" used in a statement indicates its implementation is mandatory and its correct behavior must be tested.

-  "Should" used in a statement indicates it is treated as "Shall" unless justification is provided and permission granted to omit or defer this requirement.

-  "Will" used in a statement indicates an expected behavior which occurs as the result of other requirements and therefore needs no special action.  Where a "Will" statement is likely not to occur, notification and explanation must be given.

-  "May" used in a statement indicates the requirement is optional.  If employed, it is treated as a "Shall" but no justification is needed for omitting it.

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
The SSTPA Tool will have the following data models:
Core Data Model: supports the primary purpose of SSSTPA Tool modeling large complex systems
Tool Data Model:  utility data for SSTPA Tool internal use to include User information, Version number, licenses for components etc... .
Reference Data Model:  MITRE ATT&CK data set, EMB3D data set, NIST 800-53 catalog of Controls data set used by users to associate relevant nodes.
Help Data Model:  Used to support Help functions.


The Core Data Model is one of the key innovations built into SSTPA Tool and will be leveraged by both the Frontend and the Backend.  The Core Data Model is a single sub-graph template which will be generated whenever a new (:System) node is created. 

### 1.3.1 Core Data Model Node Definition 
The Core Data Model extends the formal definition of a System as composed of Behaviors and Elements where each Element can be conceived of as a System. SSTPA Tool resolves how to start this recursive model by creating the Capability an Sandbox as starting points.  The Capability is for the primary system under development where the Sandbox is for temporary User use outside the baseline.  Behavior is not in the model but is represented by Interface and Function which describe external and internal behaviors respectively.  The model is further enhanced by certain properties as primary nodes.  These are State and Asset.  The model is extended further with a set of views with a limit of one view of each type per System.  These are Purpose, ControlStructure and Security.  Only Purpose is limited to one per System as it is the repository of all requirements allocated to the System.  Because both ControlStructure and Security have subordinate Nodes to which Requirements are allocated, there may be more than one ControlStructure or Security view.

(:Capability)
Capability is intended to be the anchor point for Tier 1 System Nodes.  From a traditional engineering work flow, Capability represents the both the customer intent for the overall project and all the requirements for the entire project.  SSTPA Tools is intended to take data from a customer and assist a technical team in all facets of systems and system security engineering to satisfy the customer's need. Only an Admin can create a (:Capability) as it represents the start of a project.
(:Sandbox)
Sandbox is a user created node similar to a capability but intended for independent development outside the formal project.  the User may clone elements of the Sandbox system into the main project in the same way as other elements are cloned. Expected use case for this feature is development of common high tier Systems early in a project lifecycle before top-down decomposition gets to them.

(:System)
A System is composed of Behaviors and Elements such that the behaviors could not originate from any individual Element alone and each Element can itself be conceived of as a System.  The core innovation for SSTPA Tools is the focus on single Systems of Interest (SoI) to mitigate the complexity of large systems. The System is the focus of the tool and the GUI.  Systems theory traditionally composes a System with Behaviors and Elements.  SSTPA Tools (and most of the Engineering world) decomposes Behaviors into Interfaces and Functions.  SSTPA Tools also extends the System model with additional views on the system.  In the SSTPA Tool, most of hte "properties" of a typical system are pulled out into the child nodes.

(:Environment)  
Environment is not a component of System, but captures the external context and threat environment the System operates in.  All components of Environment are external to the system excepting Attack which is the hazardous interaction of Environment on the System.  Environment is related to state in that the User should construct a set of unique states for each environment, but this is not enforced.  An example of an Environment is:  Wearhouse Storage, Operation in Intended Environment, Field Software Update, Hardware Installation, Decommissioning, etc...

(:Interface)  
Interface is a behavior of the System which acts on other Systems.  Interfaces communicate with Interfaces on other Systems by sharing in Connections. Interfaces can participate in Functional Flows, but only with Functions within their SoI.  Like Functions, Interfaces are abstractions and must be allocated to Elements to be realized.

(:Function)
Function is a behavior of the System which acts internally to the System.  They perform the primary behaviors of the system as well as derived security behavior, but they are only abstractions which must be allocated to Elements to be realized.  SSTPA Tool will have an add-on tool to perform Functional Flow analysis which will extend into STPA analysis when the functional flow becomes cyclic.

(:Element)
Elements are the tangible components of the System which can parent a child System.  Elements may just be Elements and in their  SoI, they perform that role,  or they can be conceived of as a System which is why they can only parent one system.  Functions and Interfaces are allocated to Elements.  Elements may contain Assets.  Elements are the target of Attacks.

(:Purpose)
Purpose is a system "view" or attribute imposed on a Engineered Systems by human beings.  Purpose is imposed by Constraints and Requirements.  A Realized system has no "purpose" but it can be Validated to assure its realized behavior is consistent with the Purpose of the Engineered System.  SSTPA Tools is focused on assuring System Security purpose is realized.  intended to assist in the development of Engineered Systems and the Analysis of Realized Engineered Systems.  SSTPA Tools enforses the limit of one Purpose per System.

(:State) 
State is a "view" of the behavior of a System in time. In a single Functional State the System exposes a fixed set of behaviors. In a System Security State, all security assurances on a System are constant. In both cases, these conditions change on state Transition, defined as a relationship rather than as a node. These two types of States are not distinguished as a User should assure their alignment. 

A (:State)-[:TRANSITIONS_TO]->(:State) relationship MAY represent: 
1. a functional transition needed for normal System behavior, 
2. a transition required to satisfy a specific (:Countermeasure), or 
3. a transition which is both functional and required to satisfy a specific (:Countermeasure). 

The distinction SHALL be represented by relationship properties on [:TRANSITIONS_TO]. 


(:ControlStructure)
Control Structure is a "view" of the system functional flow allowing STPA analysis.  In this view, Functions and Interfaces are "cast" into STPA roles for analysis, identification of Hazards and development of Countermeasures to their attacks them.

(:Asset)
The Asset is a something valuable in the System which need assurance against Loss. All criticality is derived from one or more Primary Assets in the System.  Assets may be Derived, if they protect the primary Asset but they derive their Asset status from the Primary.  

Asset is separate from Purpose as the Purpose of the System is not to protect its own Assets.  Assets are needed to achieve the Systems Purpose.  If an Asset is compromised, the stakeholders of the system will suffer harm.  Assets may be specific to Criticality Regime or may cross criticality regimes.  A regime is a governance structure which has approval authorities (who will grant Certification or Approval to Operate) in addition to formal or customary rules for development (e.g. DO 178C for development of Flight Critical software) and a community of practitioner who predominantly work within the Regime. 

(:Security)
Security is a "view" of the System containing Controls and Countermeasures.  Security is responsible for providing effective assurance needed by Assets.

(:Connection)
The Connection is the true interconnection between systems and is intended to present as the OSI communications model.  Interfaces participate in Connections.  It only exists to cover the necessity that a Connection between Systems must be owned by one and only one entity.  System Interfaces may joint Connections which the System does not own.

(:Hazard)
Hazard is the threat in the environment.  It is composed of Threat Actors, Tactics, Techniques, and Tooling from the MITRE ATTA&K framework and user created reference nodes in the MITRE format.  These component nodes are NOT part of the Core Data Model.  License forbids modification of data provided by MITRE or assertion of ownership.  To skirt this, Nodes in the Core Data Model have external reference relationships to this data and can have properties which are system specific.  Likewise, Users can create Reference data using that model and externally reference it.  This data may also be exported and sent to MITRE for proposed inclusion in subsequent releases of their framework.

(:Requirement) FIX
Requirement is the a specification statement used to realize Purpose, Constraint, Countermeasure, Interface, Function, Element, Connection, and Capability intent. A Requirement may be parented by one or more higher-level Requirement nodes and may be verified by one or more Verification procedures. Requirement parentage supports decomposition and allocation across Systems of Interest (SoI) and across tiers in the hierarchy. Requirements are analytical objects used to show traceability from customer intent through realization. 

(:Constraint)
Constraint is a rule or condition the System must enforce (e.g. weight, speed limits etc...) violation of a Constraint may cause a Hazard such that it makes the System susceptible to a Threat.

(:Control)
A Control is a security abstraction (e.g. Role Based Access Control).  These are either externally referenced or user created per the intermediate set of properties combining the NIST SP 800-53 controls catalog and the MITRE ATT&CK Mitigation, EMB3D MID sets or User created.

(:Countermeasure)
a Countermeasure is security feature which blocks Attacks on System Elements.  Where Controls are abstractions, Countermeasure implementations satisfy Controls.

(:Validation)
Validation is a procedure (similar to a verification) which assures the System when realized operates as intended.  If a System passes all validation procedures (and they were all written to perfection) the System is fit for purpose.  All Validation procedures operates on the entire holistic system in the intended environments of operation.

(:Verification)
Verification is a procedure for assuring a Requirement is implemented correctly.  If the realized system successfully pass the procedure, the Requirement is met.


(:ControlAlgorithm)
A Control Algorith is a Function or Interface of the system cast in the role of controlling a process.  All processes are controlled by something.

(:ControlAction)
A control Action is a message used by the Control Algorithm to control a Controlled Process.  There may be many Control Actions.  Some may put the System in an unsafe or unsecure position and expose the System to a Hazard.  All control behaviors have Control Actions.

(:ControlledProcess)
A Controlled Process is a Function which actually performs the intended behavior.  If done locally, it is a Function, if done in another System, it is an Interface.  All control behaviors have a Controlled Process.

(:Feedback)
Feedback is a message sent by the Controlled Process with information on how it is performing its intended behavior.  Not all systems use Feedback but most should.  This forms a control loop from what otherwise would be a Functional Flow.

(:ProcessModel)
A Process Model is a Function (never an Interface) which consumes feedback on the control behavior and uses that to tune the Control Algorithm.  For control behaviors (a form of functional flow) the Process Model is not mandatory, but without it there is no controlled feedback loop.  The Process Model is the system feature which determines the nature of the feedback loop (positive, negative, etc...)

(:Attack)
An Attack is action on a System by a threat actor who is a component of a Hazard.  SSTPA Tools follows a novel path in System Security applications by deemphasizing the Threat Actor and instead emphasizing the Attack which may be either a Tactic, Technique or Procedure depending on the level of abstraction. These will be referenced from a data set the User may add to but will be initially populated from the MITRE ATT&CK Tactic and Technique set supplemented by their EM4B TID data set. 

(:Loss)
A Loss is the compromise of a required security attribute on an Asset within a Criticality domain.  Departing just a bit from STPA, SSTPA does not treat Loss as the thing which must be avoided.  Rather, Loss is the Directed Analytical Graph (DAG) of States, Elements, Attacks and Countermeasures which lead to compromise in a specific Environment.  The resulting set of DAGs must be acceptable!  SSTPA Tools centers its analysis on developing Countermeasures and their Requirements to make each Loss DAG acceptably unlikely as to support a favorable certification decision within any criticality domain. 

For each SoI, one Loss will be created by the Frontend for each (:Asset)-Criticality-Assurance-(:Environment).  For example, an Asset which is Flight Critical with the Assurance of Availability in the "In-Flight" environment while also having the Security critical Assurance property of confidentiality in all environments will generate a single (:Loss) node for Flight Critical Availability and a number of (:Loss) nodes for Security Critical confidentiality Assurance for every (:Environment).

The Attack Tree Add-on Tool is the method for the Frontend to develop Loss.  The Attack Tree transforms the seemingly cyclical structures involved with State, Element, Attack and Countermeasure into a Directed Analytical Graph (DAG).  The Attack Tree DAG has as its Root the Loss. The root fans out using Sequential AND (SAND) and Exclusive Or (XOR) logic to form branches which terminate in either an Attack or in a new Derived Asset which forms the root of another Loss.  The branches of the tree terminating in Attack are Residual Vulnerabilities.  If acceptable nothing need be done.  If not acceptable, new Countermeasures must be added to the System which may alter State behavior or Asset Element relationships.  The Loss is calculated again to see if the Residual Vulnerabilities are acceptable.  The process repeats until this is true.  

For branches which terminate in a new Derived Asset, this too must be acceptable.  A derived Asset is something which inherits its Asset status by protecting another Asset.  Cryptographic keys are a primary example.  By themselves, they are not an Asset.  Use a symmetric key to encrypt an Asset needing Confidentiality, Authenticity or Non-Repudiation and the primary Asset may be acceptably protected, but at the cost of generating a derived Asset, the cryptographic key.  Criticality may change for Derived Assets, for example, when protecting the confidentiality of a Data Asset on a computer against a root-kit Attack, the Operating System (Element subject to Attack) may be declared a Derived Asset with the Assurance of Authenticity rather than Confidentiality because compromise of its Authenticity with a root-Kit Attack would expose the data Asset to a breach of Confidentiality.  In this case, a new Loss would be created for the Derived Asset.


 
### 1.3.2 General Modeling Rules
All relationships SHALL use Neo4j syntax: -[:RELATIONSHIP_NAME]->
All relationship names SHALL use UPPERCASE_SNAKE_CASE
All node labels SHALL be singular
Reverse relationships SHALL NOT be explicitly created unless required for performance or semantics
All properties with no value SHALL use "Null"

Where the same source and destination (:State) pair is used both for ordinary System behavior and to satisfy a (:Countermeasure), the preferred representation SHALL be a single [:TRANSITIONS_TO] relationship with TransitionKind = BOTH rather than duplicate parallel relationships. 

#### 1.3.2.1 Graph Integrity and Scale Constraints 
Duplicate Relationships: 
• The Backend SHALL prevent duplicate logical relationships between the same node pair for the same relationship type unless multiplicity is explicitly required and distinguished by properties. 

Recursive Relationship Governance: 
• Recursive relationships SHALL declare whether they are: 
  - acyclic (DAG), or 
  - cyclic-by-design (bounded analytical loop) 

• The Backend SHALL enforce the declared constraint. 

Traversal Safety: 
• All recursive traversals SHALL use bounded depth by default. 
• The system SHALL NOT execute unbounded recursive queries. 

Node Degree Monitoring: 
• The Backend SHOULD monitor node degree. 
• The Backend SHOULD warn or reject operations when configurable thresholds are exceeded. 

Pagination: 
• All list-returning endpoints SHALL support pagination and maximum result limits. 


### 1.3.3 Capability Relationship
The Capability holds place at the top of the hierarchy and serves as the anchor point for Tier 1 Systems.  In practical terms, the Capability captures the intent of the entire project from the customer's perspective.  The expected work flow for STPA Tools is to populate the Capability with the customers intent and associate their requirements set (which may be extensive and may exceed one thousand Requirements).  From these, Tier 1 systems will be created and requirements allocated.  The Capability relationships are:

(:Capability)-[:HAS_SYSTEM]->(:System)
(:Capability)-[:HAS_REQUIREMENT]->(:Requirement)

(:Sandbox)-[:HAS_SYSTEM]->(:System)

#### 1.3.2.1  Capability Constraints
Capability represents the ultimate customer's intent and there will be one Capability per project.  The requirements related to the Capability are ultimate customer requirements generally provided at the start of the project and will be barons.

Constraints:
Sandboxed Systems are intended to be temporary with elements ultimately cloned into a project system then deleted.  (:System) nodes associated with (:Sandbox) shall not be associated with a (:Capability) or in its lineage.  This behavior shall be enforced by the Backend when (:System) nodes are created or related to a (:Sandbox)


### 1.3.3 System Relationships
The system is the conceptual core of SSTPA Tools which focuses primarily on a single System of Interest (SoI) to mitigate complexity. The rich system model presented minimizes complexity by defining primary relationships and primary nodes a system view (similar to the DoDAF concept).  Secondary Nodes and relationships are the components of those views.   Systems are created in tiers in the hierarchy through parent-child relationships held through the (:Element) node.  Systems connect with other systems through the (:Connection) which must be owned by a single system.  Systems use the Connection through Interfaces.  The System relationships are:

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

The Frontend shall on creation of a (:System) create one of each primary node  type excepting (:Connection) with the Name property set to "New".
Users will have the capability to add primary node types the (:System) node through the Frontend.  From these Primary Nodes, the User will be able to enter new Secondary Nodes.


#### 1.3.3.1  System Relationships Constraints
Constraints: 

• (:Connection) SHALL NOT participate in System hierarchy relationships
• A (:System) SHALL have one and only one (:Purpose) Node 

A (:Connection) SHALL be owned by exactly one (:System) through [:HAS_CONNECTION]. 
Interfaces from the owning System and from external Systems MAY participate in that Connection through [:PARTICIPATES_IN] via (:Interface) Nodes. 
Connection ownership SHALL NOT imply that all participating Interfaces belong to the owning System. 


### 1.3.4 Core Data Model Secondary Relationships
Secondary nodes and relationships consist of the contents of the "views" under the System node.  
 
(:Environment)-[:HAS_HAZARD]->(:Hazard)

(:Connection)-[:HAS_REQUIREMENT]->(:Requirement) 

(:Interface)-[:HAS_REQUIREMENT]->(:Requirement)
(:Interface)-[:IMPLEMENTS]->(:ControlAlgorithm)
(:Interface)-[:IMPLEMENTS]->(:ControlledProcess)
(:Interface)-[:PARTICIPATES_IN]->(:Connection) 
(:Interface)-[:CONNECTS]->(:Function)


(:Function)-[:HAS_REQUIREMENT]->(:Requirement)
(:Function)-[:IMPLEMENTS]->(:ControlAlgorithm)
(:Function)-[:IMPLEMENTS]->(:ControlledProcess)
(:Function)-[:IMPLEMENTS]->(:ProcessModel)
(:Function)-[:FLOWS_TO_FUNCTION]->(:Function)
Note: (:Function) may only relate to (:Function) within its own sub-graph
(:Function)-[:FLOWS_TO_INTERFACE]->(:Interface)
Note: (:Function) may only relate to (:Function) within its own sub-graph

(:Element)-[:HAS_REQUIREMENT]->(:Requirement)
(:Element)-[:CONTAINS]->(:Asset)
(:Element)-[:PARENTS]->(:System)

(:Purpose)-[:HAS_CONSTRAINT]->(:Constraint)
(:Purpose)-[:HAS_REQUIREMENT]->(:Requirement)
(:Purpose)-[:HAS_VALIDATION]->(:Validation)

(:State)-[:TRANSITIONS_TO]->(:State)
Note: [:TRANSITIONS_TO] is the canonical representation of a state transition in the Core Data Model. It SHALL remain a relationship rather than a node. The semantic role of the transition SHALL be distinguished by relationship properties, including whether the transition is functional, countermeasure-required, or both. 

(:State)-[:HAS_HAZARD]->(:Hazard)
(:State)-[:CONTAINS]->(:Asset)

(:ControlStructure)-[:HAS_CONTROL_ALGORITHM]->(:ControlAlgorithm)
(:ControlStructure)-[:HAS_PROCESS_MODEL]->(:ProcessModel)
(:ControlStructure)-[:HAS_CONTROLLED_PROCESS]->(:ControlledProcess)
(:ControlStructure)-[:HAS_CONTROL_ACTION]->(:ControlAction)
(:ControlStructure)-[:HAS_FEEDBACK]->(:Feedback)


(:Asset)-[:HAS_LOSS]->(:Loss)

(:Security)-[:HAS_CONTROL]->(:Control)
(:Security)-[:HAS_COUNTERMEASURE]->(:Countermeasure)


(:Control)-[:MITIGATES]->(:Hazard)
(:Countermeasure)-[:SATISFIES]->(:Control)
 

#### 1.3.4.1 Secondary Relationship Constraints
Preferred Modeling Rule: 
Cross-System interaction SHALL be modeled through (:Connection). 

Connection Constraints: 
• Each (:Connection) SHALL relate to two or more (:Interface) nodes. 
• An (:Interface) SHALL NOT participate more than once in the same (:Connection). 
• (:Connection) SHALL NOT parent or child another (:Connection). 
• All (:Connection) Requirements SHALL belong to the same SoI. 

Interface Constraints for (:Interface)-[:IMPLEMENTS]->
• A single (:Interface) node shall have no more than one [:IMPLEMENTS] relationship as the individual Interface can only play one role in an STPA functional flow.

Function Constraints for (:Function)-[:IMPLEMENTS]->
• A single (:Function) node shall have no more than one [:IMPLEMENTS] relationship as the individual function can only play one role in an STPA functional flow.


Function Constraints for (:Function)-[:FLOWS_TO_FUNCTION]->(:Function) 
• Both Functions SHALL belong to the same SoI. 
• Cycles ARE allowed but SHALL NOT imply hierarchy, ownership, or requirement parentage. 
• Duplicate edges SHALL NOT be created. 
• Backend SHALL enforce same SoI constraint. 
• Graph queries for Function flow visualization SHALL NOT perform unbounded recursive expansion.

Constraint on (:Function)-[:FLOWS_TO_INTERFACE]->(:Interface)
• Both nodes SHALL belong to the same SoI. 
• SHALL NOT create cross-SoI Function relationships. 

Constraints on (:Requirement)-[:PARENTS]->(:Requirement) 
• SHALL form a Directed Acyclic Graph (DAG) 
• SHALL NOT create cycles 
• Backend SHALL enforce 
• Duplicate edges SHALL NOT exist 

Constraints on (:Element)-[:PARENTS]->(:System)
• Shall allow no more than one relationship
• SHALL form a Directed Acyclic Graph (DAG) enforced by the Backend such that the child (:System) Tier property is the (:Element) Tier Property +1
• SHALL form a Directed Acyclic Graph (DAG) enforced by having parented (:System) HID Index is calculated by appending the (:Element) HID Sequence_Number to its HID Index (e.g. E_1.2.3_4 parents S_1.2.3.4).
• Backend SHALL enforce constraints


Constraints on (:State)-[:TRANSITIONS_TO]->(:State) 
• [:TRANSITIONS_TO] SHALL be treated as a recursive relationship that is cyclic-by-design. 
• Backend SHALL enforce bounded traversal for all state-transition queries and analyses. 
• Duplicate logical [:TRANSITIONS_TO] relationships between the same source (:State) and destination (:State) SHALL NOT exist unless distinguished by relationship properties. 
• A transition required to satisfy a specific (:Countermeasure) SHALL be represented using properties on the [:TRANSITIONS_TO] relationship rather than by introducing a separate Transition node. 
• Where a [:TRANSITIONS_TO] relationship is required to satisfy a specific (:Countermeasure), the relationship SHALL identify that (:Countermeasure) by HID and/or uuid in relationship properties. 
• A single [:TRANSITIONS_TO] relationship MAY be both functional and countermeasure-required. 


---------------------------------------- 

1.3.5 Core Data Model Tertiary Relationships
Tertiary relationships model the cross-cutting relationships across "views" seeking to map the complexity of within the System.  Complexity involves cyclic structures within the sub-graph  Add-on tools will identify and manage these cyclic relationships.  

(:Constraint)-[:HAS_REQUIREMENT]->(:Requirement)

(:Hazard)-[:VIOLATES]->(:Constraint)
(:Hazard)-[:THREATENS]->(:Asset)

(:Control)-[:ENFORCES]->(:Constraint)
(:Control)-[:MITIGATES]->(:Hazard)


(:Countermeasure)-[:SATISFIES]->(:Control)
(:Countermeasure)-[:HAS_REQUIREMENT]->(:Requirement)
(:Countermeasure)-[:APPLIES_TO_FUNCTION]->(:Function)
(:Countermeasure)-[:APPLIES_TO_INTERFACE]->(:Interface)
(:Countermeasure)-[:APPLIES_TO_ELEMENT]->(:Element)
(:Countermeasure)-[:APPLIES_TO_STATE]->(:State)
Note: When a (:Countermeasure) requires a change in System state, the required transition SHALL be represented by properties on the corresponding (:State)-[:TRANSITIONS_TO]->(:State) relationship. The (:Countermeasure)-[:APPLIES_TO_STATE]->(:State) relationship identifies the affected State nodes but does not replace the transition relationship. 

(:Countermeasure)-[:APPLIES_TO_FEEDBACK]->(:Feedback)

(:Requirement)-[:PARENTS]->(:Requirement)
Note:  Parentage relationships for requirements may be within the System (sub-graph) or across systems (sub-graphs).
(:Requirement)-[:VERIFIED_BY]->(:Verification)

(:ControlAlgorithm)-[:GENERATES]->(:ControlAction)


(:ControlAction)-[:CAUSES]->(:Hazard)
(:ControlAction)-[:COMMANDS]->(:ControlledProcess)
(:ControlledProcess)-[:PRODUCES]->(:Feedback)
(:Feedback)-[:INFORMS]->(:ProcessModel)
(:ProcessModel)-[:TUNES]->(:ControlAlgorithm)

(:Loss)-[:HAS_ENVIRONMENT]->(:Environment)
(:Loss)-[:HAS_ELEMENT]->(:Element)
(:Loss)-[:HAS_STATE]->(:State)
(:Loss)-[:HAS_ATTACK]->(:Attack)
(:Loss)-[:HAS_COUNTERMEASURE]->(:Countermeasure)

(:Hazard)-[:USES_ATTACK]->(:Attack)
(:Attack)-[:DEFEATS]->(:Countermeasure)
(:Countermeasure)-[:BLOCKS]->(:Attack)
(:Attack)-[:EXPLOITS]->(:Element)

#### 1.3.5.1 Core Data Model Tertiary Relationships Constraints
Each (:Loss) SHALL be associated with exactly one (:Environment) because Loss is defined for a specific Asset-Criticality-Assurance-Environment combination.  (:Loss) SHALL be constrained to only one Criticality and one Assurance property.  All other nodes may have multiple Criticalities and Assurances.

#### 1.3.5.2 Validation Rules

If a [:TRANSITIONS_TO] relationship has TransitionKind = COUNTERMEASURE_REQUIRED or BOTH, the referenced (:Countermeasure) identified in relationship properties SHALL exist and SHALL belong to the same SoI as both endpoint (:State) nodes unless explicitly justified as a cross-SoI analytical relationship. 

### 1.3.6 Identity Model (HID + UUID)
Each node SHALL contain:
HID (Hierarchical Identifier)
uuid (Globally unique identifier)

HID Format
{TYPE}_{INDEX}_{SEQUENCE}

Example:
SYS_1.2.3_0
UUID Property
uuid: apoc.create.uuid()


 
#### 1.3.6.1 Node Type Identifier
The Node Identifier uniquely identifies each Node Type.  In STPA analysis it is common to identify nodes with a letter and a number.  Each Node type in the SSTPA Tool shall have a unique one, two or three character identifier as listed below in the format {Node Type} {Node Type Identifier}:

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
SSTPA_Tool SST
UserRegistry URG
AdminRegistry ARG
User USR
Admin ADM


#### 1.3.6.2  Index
The index uniquely identifies the Sub-graph a Node belongs to and is constructed to depict its position in the entire hierarchy.

The Index will be unique for each sub-graph and every node in the sub-graph will have the same Index.
The Index for the Capability shall be null as the data set only contains one capability whose only purpose is to attach tier 1 systems.  

Tool Data nodes do not belong to a System of Interest sub-graph and SHALL use the global empty-index HID convention.  The singleton Tool Data nodes SHALL use SST__1 for (:SSTPA_Tool), URG__1 for (:UserRegistry), and ARG__1 for (:AdminRegistry).  Registered User and Admin nodes SHALL use USR__{SEQUENCE} and ADM__{SEQUENCE} respectively.

When a node is created it shall inherit the Index of the sub-graph it belongs to excepting (:System) nodes.

When a (:System) is created as a child of a capability the index shall be calculated as the next highest integer value of other System children unless there are no other System children then it gets an index of "1".

When a System is created as the child of an Element its Index shall be the index of the Parent Element concatenated with "." concatenated with the (:Element) Node HID Sequence Number property.  For example if an (:Element) has an HID of E_1.2.3_4 than its child (:System) will have an HID of S_1.2.3.4_0.

(:Element) Nodes shall have zero or one child (:System) nodes and this constraint will be enforced by Frontend Software.  The Relationship between an (:Element) node and its single (:System) node is:

(:Element)--[:PARENTS]-->(:System)
Note, The (:System) related to here is in a child sub-graph where the new System HID index is set to the concatenation of the (:Element) Index with the (:Element) Sequence Number.

##### 1.3.6.2.1 Index Strategy

The Backend SHALL create the following indexes:

CREATE INDEX node_hid_index IF NOT EXISTS FOR (n) ON (n.HID);
CREATE INDEX node_uuid_index IF NOT EXISTS FOR (n) ON (n.uuid);
CREATE INDEX node_name_index IF NOT EXISTS FOR (n) ON (n.Name);
CREATE INDEX node_type_index IF NOT EXISTS FOR (n) ON (n.TypeName);


#### 1.3.6.3 Sequence Number
The Sequence Number is intended to distinguish nodes of the same type within the same SoI sub-graph.  
The Sequence Number for a System shall be "0" because there is only one System in the SoI sub-Graph.
The Sequence Number for a Node other than a System Node shall be next highest integer value of other nodes of the same node type in the sub-graph unless there are no others of that node type in the sub-graph, then it is the first and its value is "1". 


1.3.7 Common Property Groups
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

Property Groups are not node properties and only for organizing the display of properties and shall be enforced by the Frontend.
Property types shall be enforced by both the Frontend and the Backend.
Property defaults and ability to edit shall be enforced by the Frontend.
 
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

#### 1.3.7.1 Data Ownership Rules 
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


### 1.3.8 Type Unique Property and Relationship Groups
Each Node type will have, in addition, not common properties and relationship groups unique to its type.
Formatting rules from 1.3.7 apply here.
Headings below are Node Type names to which the unique Property Groups and Properties apply.

For node types authorized in Section 1.3.10.4 to assign imported external references by [:REFERENCES], Section 1.3.8 SHALL define property groups used to capture node-local interpretation, applicability, implementation, evidence, and analysis specific to that node. These properties SHALL NOT duplicate or overwrite authoritative imported reference item properties. Imported reference item content remains read-only and authoritative. Node-local external reference properties apply only to the SSTPA node and may differ between nodes referencing the same imported item. 

#### 1.3.8.1 Capability

Mission:
AST "A Capability To:" String edit default: "Null"
BMO "By Means Of:" String edit default: "Null"
IOTCT "In Order To Contribute To:" String edit default: "Null"

#### 1.3.8.2 System

Mission:
AST "A System To:" String edit default: "Null"
BMO "By Means Of:" String edit default: "Null"
IOTCT "In Order To Contribute To:" String edit default: "Null"

#### 1.3.8.3 Environment

Context:
Context "Context" String edit default: "Null"

#### 1.3.8.4  Connection
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


#### 1.3.8.5  Interface
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


#### 1.3.8.5.1  Interface Outgoing Relationship Properties
Only outgoing relationships with properties are identified here.

[:PARTICIPATES_IN] and [:CONNECTS] shall have hte following properties:
RelationshipNature "Nature:" Enum {PHYSICAL, LOGICAL, BOTH} edit default: "LOGICAL"  
PhysicalType "Physical Type:" String edit default: "Null"  
Example: universal joint, shaft, hydraulic linkage  
LogicalLayer "OSI Layer:" Enum{N/A, Layer 1: Physical, Layer2: Data Link, Layer 3: Network, Layer 4: Transport, Layer 5 Session, Layer 6: Presentation, Layer 7: Application} edit default: "Null"  
Protocol "Protocol:" String edit default: "Null"  
FlowDirectionality "Directionality:" Enum {Unidirectional, Bidirectional, Multicast} edit default: "Unidirectional"  
TimingClass "Timing Class:" String edit default: "Null"  
SecurityClass "Security Classification:" String edit default: "Null"  



#### 1.3.8.6 Function
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

#### 1.3.8.6.1  Function Outgoing Relationship Properties
Only outgoing relationships with properties are identified here.

[:FLOWS_TO_FUNCTION] and [:FLOWS_TO_INTERFACE] shallhave the following properties:
RelationshipNature "Nature:" Enum {PHYSICAL, LOGICAL, BOTH} edit default: "LOGICAL"  
PhysicalType "Physical Type:" String edit default: "Null"  
Example: universal joint, shaft, hydraulic linkage  
LogicalLayer "OSI Layer:" Enum{N/A, Layer 1: Physical, Layer2: Data Link, Layer 3: Network, Layer 4: Transport, Layer 5 Session, Layer 6: Presentation, Layer 7: Application} edit default: "Null"  
Protocol "Protocol:" String edit default: "Null"  
FlowDirectionality "Directionality:" Enum {Unidirectional, Bidirectional, Multicast} edit default: "Unidirectional"  
TimingClass "Timing Class:" String edit default: "Null"  
SecurityClass "Security Classification:" String edit default: "Null"  


#### 1.3.8.7 Element

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

Reference Characterization: 
ReferenceApplicabilityStatement "Applicability Statement:" String edit default: "Null" 
ReferenceExposureDescription "Exposure Description:" String edit default: "Null" 
ReferenceAssumption "Assumption:" String edit default: "Null" 

Threat / Property Context: 
ThreatSurface "Threat Surface:" String edit default: "Null" 
TechnologyType "Technology Type:" String edit default: "Null" 
DeploymentContext "Deployment Context:" String edit default: "Null" 





#### 1.3.8.8 Purpose
None

#### 1.3.8.9 State


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



#### 1.3.8.10 ControlStructure
None

#### 1.3.8.11 Asset
Type:
IsPrimary  "Primary:" Boolean edit default: "Null"

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


#### 1.3.8.12 Constraint
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
Baron  "Baron" Boolean fixed default: "True"


#### 1.3.8.14 Validation
Validation:
VStatement  "Validation Statement: " String edit default: "Null"
VMethod  "Method: " Enum {Inspection, Demonstration, Analysis, Test, Similarity} edit default: "Null"

#### 1.3.8.15 Control

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



#### 1.3.8.16 Countermeasure
None

#### 1.3.8.17 Verification

Verification:
Procedures: "Null"

#### 1.3.8.18 ControlAlgorithm
None

#### 1.3.8.19 ProcessModel
None

#### 1.3.8.20 ControlAction
None

#### 1.3.8.21 Feedback
None

#### 1.3.8.22 ControlledProcess
None

#### 1.3.8.23 Hazard
None

#### 1.3.8.24 Loss

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


Attack Tree:
AttackTreeJSON "Attack Tree Source serialized JSON document fixed default: N/A
AttackTreeFormat "Format" String fixed default: "SSTPA-ATF-1.0"
AttackTreeVersion — "Version: " Integer fixed default: N/A
AttackTreeCreated — "Created: " datetime fixed default: N/A
AttackTreeStatus — Enum {AUTO_GENERATED, ANALYST_REFINED, BASELINED, EXPORTED}


#### 1.3.8.25  Attack
None


----------------------------------------
## 1.4 Tool Data Model
SSTPA Tools is intended to support a large dispersed engineering team. In this initial version the Backend and Frontend will be bundled for use on a single system.  Data in the Core Data Model will have ownership and also record the creator.  Owner and creator contact information will be captured as email address.  SSTPA Tools will have an add-on tool called "Message Center" which allows users to message each other and allows the system to notify Users of changes to properties and relationships on data they own.


Tool Data Model consists of utility information for SSTPA Tools maintained by the backend.  It consists of:
Data on the SSTPA Tool: (:SSTPA_Tool)
Tool Data will parent a single node SSTPA Tools Data which shall have properties typical of a commercial software installation to include the license.
The SSTPA Tool Data root SHALL be represented as a singular (:SSTPA_Tool) node with Name "SSTPA Tools Data".  The nodes called "Users" and "Admins" below are singular registry nodes represented by labels (:UserRegistry) and (:AdminRegistry), with Name values "Users" and "Admins" respectively.  "Users" and "Admins" SHALL NOT be used as plural Neo4j labels.
Data on Users: (:User)
Users can own data and edit properties identified with "edit".
Data on Admins (:Admin)
Admins cannot own data but can edit and commit certain specific properties identified as for Admins Only otherwise their Commit is invalid.  Admin Users can access the backend Webserver to view telemetry information.
Data on the Host system (:Host)
Data on Messages for Users: (:Message) 
Data on Messages for Admins: (:Message) 

### 1.4.1  Onboarding
The SSTPA Tool Installer shall be configured to capture the Installer's Name and e-mail contact information.  It shall establish an Admin account and a User account for the Installer.
When invoking the SSTPA Tool, it shall present a login window asking for User Name and Password and there will be a link below with the label "New User" and a link labeled "Admin".  If the User enters a valid Username and password, the GUI will start with the identified (:User) as the active User for the purposes of the ownership and permissions model.  If  the User clicks the "New User" Link, they will be presented with a screen where they can create a new (:User) node by editing its properties. then they will be redirected to the Login screen.  If the User clicks on "Admin" link than a new Admin login window appears with the same behavior as the User login window.  If a valid Admin properly logis in, then the active User for the purposes of ownership and permissions model is the (:Admin). Under User Name and Pasword is a link with label "I want to be an Admin".  If the User clicks this, a new screen is presented which has the Admin User Name and password fields and a New Admin and New Admin Password field.  A valid existing Admin must "login" to authorize the creation of a new (:Admin) who sets their User Name and Password and other properties.  on success the User is sent to hte Admin login screen to login as a new Admin.


### 1.4.2 Admin Data
On installation of SSTPA Tools it will have , the "Installer" shall be required to provide an Admin email and establish a user account.  Admin cannot own data but can edit data and certain elements of data are fixed for normal users can be edited by Admin.  These properties are explicitly identified.
All users will "login" to the system as either an existing User or a new User.  If new User they will setup a User account.  Per the initial security model, no password or authentication will be used.

Tool Data will parent a single AdminRegistry node called Admins and nodes associated with each registered Admin will be a child of the AdminRegistry node with properties describing the Admin.

### 1.4.3 User Data
Tool Data will parent a single UserRegistry node called Users and nodes associated with each registered user will be a child of the UserRegistry node with properties describing the User.


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


1.
1.3.10 External Reference Framework Data Model 

The SSTPA Tool SHALL support imported reference frameworks in the Backend for read-only use by the GUI and Add-on Tools. These frameworks will include: 
•	NIST SP 800-53r5
•	MITRE ATT&CK
•	MITRE EMB3D

These frameworks will be manually inserted into the database prior to delivery.  Future versions of SSTPA Tools will allow update of this data which shall be fixed (not editable by any user).  This restriction enforces the data license agreement.  Future external reference update capability will update from a provided BLOB image matching the data structure in the database. SSTPA Tools will not scrape the internet for data.  

The purpose of these imported reference frameworks is to allow the User to: 
•	Navigate authoritative external reference data
•	Read the properties of imported reference items
•	Assign an external reference to selected SSTPA node types without cloning the reference item into the SSTPA analytical graph

Imported reference framework data SHALL be stored in Neo4j as a distinct but connected graph structure separate from the SSTPA System of Interest (SoI) sub-graphs. This is consistent with the SRS requirement that the Backend import NIST SP 800-53, MITRE ATT&CK, and MITRE ESTM data into graphical form and that the GUI provide add-on tools for populating selected node properties from those datasets. 


## 1.5 Reference Data Model

### 1.5.1 External Reference Framework Nodes 

The Backend SHALL support the following imported framework node labels: 
•	(:ReferenceFramework)
•	(:ReferenceItem)

The Backend MAY additionally support framework-specific labels, including: 
•	(:NistControl)
•	(:NistControlEnhancement)
•	(:AttackTactic)
•	(:AttackTechnique)
•	(:AttackMitigation)
•	(:Emb3dThreat)
•	(:Emb3dMitigation)
•	(:Emb3dProperty)

All imported reference nodes SHALL also carry the label (:ReferenceItem) except (:ReferenceFramework). 


### 1.5.2 Reference Framework Identity 

Each imported (:ReferenceItem) node SHALL contain: 
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

At minimum, the following relationships SHALL be supported: 
•	(:ReferenceFramework)-[:HAS_ITEM]->(:ReferenceItem)
•	(:ReferenceItem)-[:HAS_CHILD]->(:ReferenceItem)
•	(:ReferenceItem)-[:RELATED_TO]->(:ReferenceItem)

The Backend MAY support framework-specific relationships reflecting the source data semantics. 


### 1.5.4 SSTPA Node External Reference Relationships 

The following SSTPA node types SHALL support assignment of external references: 
•	(:Control)
•	(:Element)
•	(:System)
•	(:Hazard)
•	(:Attack)
•	(:Countermeasure)

The relationship between an SSTPA node and an imported reference item SHALL be: 
•	(:Control)-[:REFERENCES]->(:ReferenceItem)
•	(:Element)-[:REFERENCES]->(:ReferenceItem)
•	(:System)-[:REFERENCES]->(:ReferenceItem)
•	(:Hazard)-[:REFERENCES]->(:ReferenceItem)
•	(:Attack)-[:REFERENCES]->(:ReferenceItem)
•	(:Countermeasure)-[:REFERENCES]->(:ReferenceItem)

Reverse relationships SHALL NOT be explicitly created unless required for performance or semantics. 

###1.5.5 Allowed Reference Assignments 

The Backend SHALL validate that only semantically valid reference assignments are allowed. 

The initial allowed assignments SHALL be: 
•	(:Control) → NIST SP 800-53 Controls, ATT&CK Mitigations, EMB3D Mitigations
•	(:Element) → EMB3D Properties, EMB3D Threats
•	(:System) → EMB3D Properties, NIST SP 800-53 Controls
•	(:Hazard) → EMB3D Threats, ATT&CK Techniques
•	(:Attack) → ATT&CK Tactics, ATT&CK Techniques, EMB3D Threats
•	(:Countermeasure) → EMB3D Mitigations, ATT&CK Mitigations, NIST SP 800-53 Controls

The Backend SHALL reject assignments outside these allowed mappings. 


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

Startup Software shall launch from a desktop icon or Command line.

The Startup Software shall present the user with a dialog with default SSTPA theme and startup animation.

The Startup Software shall allow the user to change Theme and if a new theme is chosen shall change to that theme and play its startup animation. 

The Startup Software shall connect to a Backend or start the Backend on the Local Machine.  In future versions, Startup Software will connect to remote Backends.

Startup Software shall present the User with the Backends' list of Users and Roles and allow the User to select or Add a User.  I know this is bad security practice, but this is a stand in for a more comprehensive security implementation in a future version.   and e-mail contact information.

Startup Software shall launch the Frontend Software after User selects or adds a User ID enters information 

On receiving the Shutdown command from the Frontend software, Startup Software shall assure both frontend and backend are properly shutdown preserving stored data (i.e. don't kill the database while transactions are in process). 



## 2.2 BackEnd

The backend database will include the graph database and support software needed for ACID compliance.  It will use the most current stable NEO4J Community Edition with a defined pathway to the Enterprise Edition on customer desire.  Backend will be divided into docker containers and Docker Compose.  User will connect and interact with a reverse proxy which will connect to the database.  The reverse proxy will collect and present telemetry on backend performance.

The Backend shall be configured to execute CYPHER\_25 scripts.

The back-end shall support multiple concurrent connections.

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

only the reverse proxy shall be internet-facing. 

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
Backend Software shall be written in the most current stable version of the Go language.


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

The backend Database shall be the latest stable Neo4j Community Edition.

Responsibilities: 
persist graph data 
provide ACID transactions 
accept Bolt protocol connections from backend 


Typical internal connection: 

Go backend -> neo4j:7687 
Do not expose Neo4j publicly. 


### 2.2.6 Telemetry

The backend shall use Open Telemetry Collector for backend telemetry

Responsibilities: 

receive telemetry from backend 
batch/process telemetry 
export traces to Tempo 
optionally expose Prometheus-scrapable metrics or forward OTLP metrics 


Typical flow: 

Go backend -> OTel Collector -> Tempo/Jaeger 


### 2.2.7 Metrics

The backend shall use Prometheus for metrics.

Responsibilities: 
scrape /metrics endpoints 
store time-series metrics 
answer PromQL queries from Grafana 


Typical scrape targets: 

backend:8080/metrics 
otel-collector metrics endpoint 

optionally Neo4j exporter if you add one 


### 2.2.8 Traces

The backend shall use Tempo.

Responsibilities: 
store distributed traces 
let Grafana drill into request traces 


Typical flow: 

OTel Collector -> Tempo/Jaeger 
Grafana -> Tempo/Jaeger 


### 2.2.9 Dashboard

The backend shall use Grafana. 

Responsibilities: 
display metrics dashboards 
display traces 
correlate slow requests with backend metrics 


Typical data sources: 
Prometheus 
Tempo

Grafina shall present an accessible dashboard via the reverse proxy

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

The Backend SHALL provide import tools for: 
•	NIST SP 800-53r5
•	MITRE ATT&CK
•	MITRE EMB3D

The import process SHALL: 
•	Convert source data into graph format
•	Preserve framework version information
•	Preserve source identifiers
•	Preserve hierarchy and related-item relationships where supported by the source data
•	Avoid creating duplicate imported reference items for the same framework version and source identifier

The Backend SHALL support re-import or update of framework data. 


2.2.10.10.2 Framework Retrieval Requirements 

The Backend SHALL provide endpoints to retrieve: 
•	List of available reference frameworks
•	Framework version metadata
•	Reference item by ExternalID
•	Reference item by internal uuid
•	Reference items by framework and type
•	Parent-child structure within a framework
•	Related reference items for a selected item

Responses SHALL include: 
•	Framework name
•	Framework version
•	External identifier
•	Type
•	Name
•	ShortDescription
•	LongDescription
•	SourceURI
•	Related items summary

2.2.10.10.3 Framework Search Requirements 

The Backend SHALL support search across imported framework data. 

Search SHALL support: 
•	exact search by ExternalID
•	partial search by Name
•	partial search by ShortDescription
•	filtering by framework
•	filtering by framework version
•	filtering by imported reference item type

Search results SHALL include: 
•	Framework name
•	Framework version
•	ExternalID
•	Name
•	Reference item type
•	ShortDescription

2.2.10.10.4 Assignment Validation Requirements 

The Backend SHALL validate [:REFERENCES] relationships before creation. 

Validation SHALL confirm: 
•	the source SSTPA node exists
•	the reference item exists
•	the selected SSTPA node type is permitted to reference the selected imported item type
•	duplicate assignment does not already exist

The API SHALL return: 
•	Valid / invalid
•	Reason for invalidity

2.2.10.10.5 Assignment Mutation Requirements 

The Backend SHALL provide endpoints to: 
•	create a [:REFERENCES] relationship
•	delete a [:REFERENCES] relationship
•	list all [:REFERENCES] relationships for a selected SSTPA node

All such write operations SHALL be transactional and ACID compliant. 


##### 2.2.10.10.6 Performance Requirements 

The Backend SHALL create indexes for imported reference data sufficient to support efficient: 
•	ExternalID lookup
•	framework filtering
•	item type filtering
•	text search by name and description

The Backend SHALL optimize retrieval for: 
•	framework hierarchy navigation
•	reference item inspection
•	assignment validation
•	Data Drawer relationship display




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

Backend shall allow display and configuration of ports, configs, and volumes
Backend shall send configuration information to Frontend or Startup on connection 


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
The GUI shall execute as a stand alone desktop application.
The GUI shall operate in a single window with the capability to execute add-on tools in pop-up windows
The GUI shall connect to the Backend and commit data only after commit confirmation in a confirmation dialog.

## 3.2  GUI Style
The GUI shall have a style defined by a Style sheet (.css file) with a default style and user selectable alternate styles.
The GUI shall be organized in panels.

### 3.2.1  Default_Style.css file

The Default_Style.css is the single default User Interface (UI) style.  

Default\_Style.css shall implement a "high-tech, dark, liquid glass aesthetic to include presentation of data in a data grid pattern on tilt cards with background animation simulating slow dynamic motion and details for user edit of specific records presented using the right-side drawer interaction pattern.

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
	The top Panel of the SSSTPA Tools GUI Shall show SSTPA Logo on top left with "SSTPA Tools" name and version at center.  right shall contain status information from Backend in smaller font to include connection IP, port and connection status (in Courier font and a contrasting color). In the right it shall also show the current User name, mail Icon titled "Message Center" (an add-on Tool) and a gear Icon for changing GUI parameters such as changing Style and displaying license, and system version information.   

3.3.1 Message Center Add-on Tool 
The Branding Panel SHALL display a mail icon labeled or tooltiped as Message Center. 

The Message Center SHALL display an unread indicator when unread messages exist. 

Selecting the Message Center icon SHALL open a pop-up window. 

The Message Center pop-up SHALL display the current user’s mailbox only. 

The Message Center pop-up SHALL not change the current SoI. 

The Message Center SHALL be a specialized add-on tool consistent with the pop-up model already allowed by the SRS. 



## 3.4  SSTPA Control Panel
The SSTPA Control Panel SHALL be below the SSTPA Branding Panel and contain icons for Add-on Tools starting from the Left with the "Navigator Tool" going in sequence with Add-on Tools ending with the "Reports" drop-down Menue. 

The SSTPA Control Panel shall present an ICON for 'Shutdown" at the far right of the panel as a typical power icon but in red color. 

If User selects a menu item or icon which does not have real functionality attached, it shall present an alert dialog titled "Under Construction" with a construction icon and an "OK" button.  On click of the "OK" button, the alert will close. 

### 3.4.1 The "Navigator Tool" 
The "Navigator Tool" will perform the following core functions:
It is the means by which he SSTPA Tools GUI User Selects a System of Interest (SoI) for the rest of the GUI
It allows the User to search the entire Hierarchy to identify specific nodes and display it graphicly
It allows the User to explore the System Hierarchy while not changing the current SoI
It allows the Use to navigate to and select a node from another System to clone into the current SoI
It allows the User to select a (:Connection) node owned by another System and connect an SoI (:Interface) to it
User can set Navigator Tool to display participants and owner in a specific (:Connection) and graphically visualize selected connections up to all of them.


The tool described here shall be branded at top of window as "Navigator Tool"


The Navigator Tool shall: 
•	place (:Capability) at the top or central anchor position
•	display connected (:System) nodes using force-directed or constrained hierarchical layout behavior
•	support smooth zoom, pan, and animated re-centering
•	preserve visible spatial continuity during user interaction
•	prevent navigation states where all nodes are moved completely out of view

---------------------------------------- 

3.4.1.1 Modes of Operation 
The Hierarchy Search Tool shall support four modes: 

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
This operation shall only be made available when cloning nodes having [:HAS_REQUIREMENT] relationships (if no requirements just use 1))   

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
If no SoI has been selected, the tool shall scale to show the full hierarchy. 
If an SoI has already been selected, the tool shall initially center the graph on the current SoI and visually distinguish it from all other nodes


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

Navigator Pop-up Window shall be composed of the following elements: 
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

•	Avoid unnecessary full-layout recomputation during interaction

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

In the Hierarchy View, the Requirement Tool shall depict the focused requirement and its Heritage and lineage to a User selected depth.  

In the Allocation View, the Requirements Tool shall depict  the Valid Node with a [:HAS_REQUIREMENT] relationship and all allocated requirements

The User may move between views by, in the Hierarchy View selecting to display Allocations and the User selecting one. or in the Allocation View by the User selecting one Requirement and selecting it for display in Hierarchy View.


---------------------------------------- 

3.4.2.3 Supported Node Context 

The tool SHALL support invocation ONLY when the Data Drawer is open for: 
•	(:Capability)
•	(:Purpose)
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
The Reports Dropdown Menu shall list the following reports to create
System Description
System Specification
Requirement‑Traceability Gap Analysis
Controls List

### 3.4.3.1 System Description Report
System Description Report is a text based hierarchical description of the SoI, its primary nodes and relationships followed by its secondary nodes and relationships.
Report shall be in text, markdown, MS Word or PDF format.
  
### 3.4.3.2 System Specification Report
System Specification Report is a text based list of requirements for the SoI organized by the node they are related to within the SoI. It begins with a description of the SoI and its properties. One section per primary element type with [HAS_REQUIREMENT] relationship.  Sub-section for each entity of the type followed by an ordered list of requirements showing uuid and RequirementStatement properties for each.
Report shall be in text, markdown, MS Word or PDF format.

### 3.4.3.3 Requirement‑Traceability Gap Analysis
Requirement‑Traceability Gap Analysis is a text based report identifying problematic requirements in the SoI.  This report is less informational and focused to remediating action. It is organized in the same way as the System Specification Report excepting when referring to Requirement properties it shows the UUID followed by the analytical properties: Baseline, Orphan and Baron.   Orphan and Baron properties are not user editable, but shall be by the generation of this report.  

Baseline is not set by running this report but the (:Requirement) Baseline property is reported (note:  projects deal with baselines in a number of ways and the tool must be flexible with this property to support most use cases).

For a (:Requirement) property "Orphan" shall be true if any of these is true:
1.  It has no parent (:Requirement)

In other words, a Requirement cannot be created without a parent, so an Orphan is likely the result of a node deletion and the Requirement is flagged foe reparenting or removal.

For a (:Requirement) property "Baron" shall be true if any of these is true:
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

Reference Tool shall initialize in Assignment Mode when the Data Drawer for a valid node type is open otherwise Reference Tool opens in Research Mode.
The User shall be able to switch the Research Tool into Research Mode at any time and switch back to the view in Assignment Mode.
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

The User SHALL be able to switch between modes without changing the current SoI. 

---------------------------------------- 

#### 3.4.6.3 Invocation 

The Flow Tool SHALL: 

• Be launched from the SSTPA Control Panel  
• Initialize to Functional Flow Mode if a Data Drawer is open for (:Function) or (:Interface) and by default  
• If a Data Drawer is open for (:Function) or (:Interface), center and focus on that node
• Initialize to STPA Control Flow mode if a Data Drawer is open for (:ControlStructure), (:ControlAlgorithm), (:ProcessModel), (:ControlledProcess), (:ControlAction), or (:Feedback).
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

The GUI shall refresh after the Flow Tool performs a commit: 

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

### 3.4.7  The Loss Tool
#### 3.4.7.1 Purpose 
In the body of evidence needed for certification, nothing is more important than the analysis of Loss and the identification and approval of Residual Vulnerabilities.  All Assets will have Residual Vulnerabilities, the purpose of the Loss Tool is to develop a Loss View which is both comprehensible and acceptable to system stakeholders and in particular certification authorities.

In SSTPA Tools, a (:Loss) Node will be created from the commit of an (:Asset) Node.  One (:Loss) node will be created for every Asset Criticality and Assurance pair.  The User uses the Loss Tool to allocate Loss to (:Environment) as their should be one (:Loss) for every Criticality Assurance and Environment combination but there are many Environments where a Loss is not relevant and this is the User's call to make.

The Attack Tree is a Directed Analytical Graph (DAG) developed from the cyclical relationships derived from (:Loss).  The Attack Tree "un-rolls the relationship structure into linear form and assigns Sequential AND (SAND) and Exclusive OR (XOR) relationships not available in the graph model.  The Loss Tool will auto-generates an Attack Tree diagram based on the conventions of Structured Attack Tree Analysis and allows the User to extend it, modify it and add Countermeasures and Attacks.  The User must terminate each branch with an Attack (which will be a Residual Vulnerability) or a new (:Asset) node with derived criticality and Assurances which will spawn new (:Loss) nodes as children under the new derived Assets.

#### 3.4.7.1 Invocation 

The Loss Tool SHALL: 

• Be launched from the SSTPA Control Panel  
• If a Data Drawer is open for (:Loss) and there is a valid AttackTreeJSON property than the Loss Tool shall display the serialized JSON document in its canvas.
• If a Data Drawer is open for (:Loss) and there is not a valid AttackTreeJSON property than the Loss Tool shall create a serialized JSON document and display it in its canvas.
• If the GUI is in any other state, the Loss Tool shall present the User with (:Loss) nodes from the SoI and their short description and allow the user operate on one which shall be processed as above.




### 3.4.8 Message Center 

#### 3.4.8.1 Purpose 
The Message Center provides the current user access to direct messages and owner-change notification messages. 
The Message Center allows the current user to send direct messages to other users and the Admin.

##### 3.4.8.2 Window behavior 
The Message Center SHALL open in a pop-up window. 
The pop-up SHALL be closable without affecting staged edits in the Data Drawer. 
The pop-up SHALL support refresh. 
The pop-up SHALL preserve the current SoI and current GUI navigation state. 


##### 3.4.8.3 Mailbox list view 
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


##### 3.4.8.4 Message open behavior 

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


##### 3.4.8.5 Direct messaging 
Users SHALL be able to send direct messages to other users. 
Direct messages SHALL be stored in the recipient mailbox. 
Direct messages MAY optionally reference one or more HIDs. 

##### 3.4.8.6 Change notification messages 
The system SHALL generate change notification messages automatically on commit when required by ownership rules. 
Sender SHALL be the current user who committed the change. 
Recipient SHALL be the Owner of the affected node. 
The HID column SHALL show the primary affected HID; where multiple HIDs are affected, the detail view SHALL show the full list. 


##### 3.4.8.7 Delete behavior 
Delete in the current version SHOULD be soft delete. 
Deleting a message from a mailbox SHALL remove it from the current user’s active list only. 
Deleted messages SHOULD remain recoverable for audit unless system retention rules later remove them. 

##### 3.4.8.8 Read state 
Opening a message SHALL mark it read unless the user closes before content load completes. 
Unread count SHALL update after read and delete actions. 



## 3.5 System of Interest Panel
The SSTPA System of Interest Panel shall be below the SSTPA Control Panel.
Data presented in the System of Interest Panel shall be not editable.  The user will be able to edit data in the associated data drawer.  Data in this panel will be updated when new data from the data drawer is committed.

If there is no current System of Interest (SoI) selected, the panel shall display at top center "Select a System of Interest".

The SSTPA Tools System of Interest Panel shall display SoI properties: HID, Name and ShortDescription.

The SSTPA Tools System of Interest Panel shall display an "Edit" icon which on-click shall open all system properties in Data Drawer.  


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

All Node deletions shall follow an alert / confirm pattern.

When a Node is deleted other than a (:System) or (:Element) node parenting a (:System) node, the GUI shall identify all nodes within the SoI which are orphaned by this action and include this notification in the Alert/Confirm Dialog.  If there are any orphaned nodes, the Alert / Confirm dialog shall include a warning:  "WARNING:  Cancel and Re-Associate the Following Nodes or They will be Deleted". 

Node deletion SHALL NOT automatically cascade outside the current SoI.

Deletion of (:System) Nodes SHALL require explicit user confirmation and preview.


## 3.7  The Data Drawer
The right-side Data Drawer Shall be the single edit surface for the GUI and will implicitly validate node associations.  Note; Add-on Tools may also allow edit and commit but this is in a pop-up window outside the GUI.

All node associations SHALL be validated via Backend API prior to commit.


### 3.7.1 General Behavior 
•	The Data Drawer SHALL slide in from the right side of the GUI but shall not obscure or de-activate the Branding Panel or the Control Panel.
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

When Removing relationship is selected, the effected node shall be assessed if relationship removal leaves it an orphan.  If the node is not an orphan, the relationship can be safely removed and the change is made when the Data Drawer is committed.  If the Node will become an orphan when the relationship is removed, it shall be treated as a deleted node subject to the same alert/confirm prior to the "Commit".

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

If owner notification generation fails, the Frontend SHALL display the overall commit failure and no staged changes shall be considered committed. 


### 3.7.6 Navigation 

The Data Drawer MAY support: 
•	Navigation between related nodes
•	Breadcrumb display of relationship context

Navigation in SSTPA Tools is through the SoI Navigator.  Navigation within a Data Drawer to related nodes within the SoI is allowed.  When displaying relationships outside the SoI (for (:Interface) and (:Element) nodes) attempts to edit a node outside the SoI shall be responded with an Alert dialog with the Message "Navigate to:  " n.HID " to edit".  The HID for the node shall be copiable via icon to allow pasting into the SoI Navigator.

### 3.7.7 Consistency Requirement 

All node editing SHALL occur through the Data Drawer. 

Pop-up windows SHALL NOT be used for editing Secondary or Tertiary nodes. 

Pop-ups SHALL be reserved for: 
•	SoI selection
•	Graph visualization tools
•	Specialized add-on tools




# 5  SSTPA Tool Component Copyright

All SSTPA Tool software components shall include a copyright statement:

"2025 Nicholas Triska. All rights reserved.
The SSTPA Tools software and all associated modules, binaries, and source code are proprietary intellectual property of Nicholas Triska.  Unauthorized reproduction, modification, or distribution is strictly prohibited.  Licensed copies may be used under specific contractual terms provided by the author."



All SSTPA Data components shall include a copyright statement:

"2025 Nicholas Triska.
The SSTPA Tools is proprietary software. However, users retain ownership of data and reports generated during legitimate use of the software, except for embedded proprietary schemas and templates."


# 6 SSTPA Tool Constraints

The SSTPA Tools shall operate on an air-gapped Microsoft Windows 11 Enterprise based network with no access to the internet.   

As the SSTPA tool is developed on a Linux based system the SSTPA Tool shall also function on a system with the following characteristics: Operating System: Ubuntu Studio 25.04

KDE Plasma Version: 6.3.4

KDE Frameworks Version: 6.12.0

Qt Version: 6.8.3

Kernel Version: 6.14.0-27-generic (64-bit)

Graphics Platform: Wayland

Processors: 28   Intel  Core  i7-14700K

Memory: 31.1 GiB of RAM

Graphics Processor: Intel  Graphics



# 6.1 SSTPA Tool Architecture

the SSTPA Tool architecture shall be implemented with minimum complexity.  When integrating a capability, the developer shall asses if libraries or functions already existing in the code-base can execute the new capability before introducing a new library or function.



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
The SSTPA Tools shall operate on an air-gapped Microsoft Windows 11 Enterprise based network with no access to the internet.  The Installer shall execute a window needing expected available components in this architecture and present the User with resource requirements needed for successful installation prior to installation.



