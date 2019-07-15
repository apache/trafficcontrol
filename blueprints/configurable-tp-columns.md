# CONFIGURABLE TABLE COLUMNS

## Problem Description
<!--
*What* is being asked for?
*Why* is this necessary?
*How* will this be used?
-->

As an SRE, I need the ability to override the default order and visibility of table columns in Traffic Portal and persist the overridden settings for future sessions. In addition, the values of hidden columns should be excluded from the global table search.

When viewing a collection of resources (servers, delivery services, etc.), an SRE may require knowledge of different resource attribute values depending on the current task. By allowing an SRE to configure the visibility and order of TP table columns, an SRE can retrieve the necessary information in a timely fashion. In addition, by allowing an SRE to hide columns and exclude the value of those columns from a table search, it provides a more powerful and targeted search mechanism.

## Proposed Change
<!--
*How* will this be implemented (at a high level)?
-->

Add the ability to reorder TP table columns using drag and drop and provide the user with the ability to toggle on/off the visibility of the column. Column visibility will determine if a column is included/excluded from the global search. Any changes made by the user shall be persisted.

### Traffic Portal Impact
<!--
*How* will this impact Traffic Portal?
What new UI changes will be required?
Will entirely new pages/views be necessary?
Will a new field be added to an existing form?
How will the user interact with the new UI changes?
-->

 - Add drag and drop capabilities to each table column so the user can reorder the columns to fit their needs. There is a jquery datatables extension (ColReorder) that could accomplish this.
 - Add a drop down with a check box for each table column so the user can toggle the visibility of the column. Column visibility can be controlled using the API of jquery datatables.
 - Exclude a column from the global table search if it is not visible. Column searchability cannot be controlled using the API of jquery datatables so the jquery datatables library may need to be patched to accomplish this.
 - Leverage the browser's local storage to persist the user-defined order and visibility of the table columns.

### Traffic Ops Impact
<!--
*How* will this impact Traffic Ops (at a high level)?
-->

N/A

#### REST API Impact
<!--
*How* will this impact the Traffic Ops REST API?

What new endpoints will be required?
How will existing endpoints be changed?
What will the requests and responses look like?
What fields are required or optional?
What are the defaults for optional fields?
What are the validation constraints?
-->

N/A

#### Client Impact
<!--
*How* will this impact Traffic Ops REST API clients (Go, Python, Java)?

If new endpoints are required, will corresponding client methods be added?
-->

N/A

#### Data Model Impact
<!--
*How* will this impact the Traffic Ops data model?

What changes to the lib/go-tc structs will be required?
-->

N/A

#### Database Impact
<!--
*How* will this impact the database schema?

What new tables and columns will be required?
How will existing tables and columns be changed?
What are the column data types and modifiers?
What are the FK references and constraints?
-->

N/A

### ORT Impact
<!--
*How* will this impact ORT?
-->

N/A

### Traffic Monitor Impact
<!--
*How* will this impact Traffic Monitor?

Will new profile parameters be required?
-->

N/A

### Traffic Router Impact
<!--
*How* will this impact Traffic Router?

Will new profile parameters be required?
How will the CRConfig be changed?
How will changes in Traffic Ops data be reflected in the CRConfig?
Will Traffic Router remain backwards-compatible with old CRConfigs?
Will old Traffic Routers remain forwards-compatible with new CRConfigs?
-->

N/A

### Traffic Stats Impact
<!--
*How* will this impact Traffic Stats?
-->

N/A

### Traffic Vault Impact
<!--
*How* will this impact Traffic Vault?

Will there be any new data stored in or removed from Riak?
Will there be any changes to the Riak requests and responses?
-->

N/A

### Documentation Impact
<!--
*How* will this impact the documentation?

What new documentation will be required?
What existing documentation will need to be updated?
-->

This new functionality needs to be documented in the "Using Traffic Portal" section of the documentation and any existing screenshots of TP tables need to be modified.

### Testing Impact
<!--
*How* will this impact testing?

What is the high-level test plan?
How should this be tested?
Can this be tested within the existing test frameworks?
How should the existing frameworks be enhanced in order to test this properly?
-->

UI tests can be added to test the proposed behaviour. It is probably overkill to test this functionality on all tables but one or two should suffice.

### Performance Impact
<!--
*How* will this impact performance?

Are the changes expected to improve performance in any way?
Is there anything particularly CPU, network, or storage-intensive to be aware of?
What are the known bottlenecks to be aware of that may need to be addressed?
-->

N/A

### Security Impact
<!--
*How* will this impact overall security?

Are there any security risks to be aware of?
What privilege level is required for these changes?
Do these changes increase the attack surface (e.g. new untrusted input)?
How will untrusted input be validated?
If these changes are used maliciously or improperly, what could go wrong?
Will these changes adhere to multi-tenancy?
Will data be protected in transit (e.g. via HTTPS or TLS)?
Will these changes require sensitive data that should be encrypted at rest?
Will these changes require handling of any secrets?
Will new SQL queries properly use parameter binding?
-->

N/A

### Upgrade Impact
<!--
*How* will this impact the upgrade of an existing system?

Will a database migration be required?
Do the various components need to be upgraded in a specific order?
Will this affect the ability to rollback an upgrade?
Are there any special steps to be followed before an upgrade can be done?
Are there any special steps to be followed during the upgrade?
Are there any special steps to be followed after the upgrade is complete?
-->

This functionality is not dependent on the upgrade of any other TC component. TP can be upgraded or rolled back as needed.

### Operations Impact
<!--
*How* will this impact overall operation of the system?

Will the changes make it harder to operate the system?
Will the changes introduce new configuration that will need to be managed?
Can the changes be easily automated?
Do the changes have known limitations or risks that operators should be made aware of?
Will the changes introduce new steps to be followed for existing operations?
-->

These changes will allow TP users to get to relevant resource data faster. However, some training or a demo of the new functionality is required.

### Developer Impact
<!--
*How* will this impact other developers?

Will it make it easier to set up a development environment?
Will it make the code easier to maintain?
What do other developers need to know about these changes?
Are the changes straightforward, or will new developer instructions be necessary?
-->

N/A

## Alternatives
<!--
What are some of the alternative solutions for this problem?
What are the pros and cons of each approach?
What design trade-offs were made and why?
-->

N/A

## Dependencies
<!--
Are there any significant new dependencies that will be required?
How were the dependencies assessed and chosen?
How will the new dependencies be managed?
Are the dependencies required at build-time, run-time, or both?
-->

The following client-side dependency is required:
- https://datatables.net/extensions/colreorder/

## References
<!--
Include any references to external links here.
-->

- https://datatables.net/extensions/colreorder/
