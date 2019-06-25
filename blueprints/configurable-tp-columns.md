# *What* is being asked for?
As an SRE, I need the ability to override the default order and visibility of table columns in Traffic Portal and persist the overridden settings for future sessions. In addition, the values of hidden columns should be excluded from the global table search.

# *Why* is this necessary?
When viewing a collection of resources (servers, delivery services, etc.), an SRE may require knowledge of different resource attributes depending on the current task. By allowing an SRE to configure the visibility and order of TP table columns, an SRE can retrieve the necessary information in a timely fashion. In addition, by allowing an SRE to hide columns and exclude the value of those columns from a table search, it provides a more powerful and targeted search mechanism.

# *How* could this be implemented?
 - Add drag and drop capabilities to each table column so the user can reorder the columns to fit their needs. There is a jquery datatables extension (ColReorder) that could accomplish this.
 - Add a drop down with a check box for each table column so the user can toggle the visibility of the column. Column visibility can be controlled using the API of jquery datatables.
 - Exclude a column from the global table search if it is not visible. Column searchability cannot be controlled using the API of jquery datatables so the jquery datatables library may need to be patched to accomplish this.
 - Leverage the browser's local storage to persist the user-defined order and visibility of the table columns.

## *Which* TC components may be impacted?

	- [x] Traffic Portal
 	- [ ] Traffic Ops
		- [ ] API
		- [ ] Database
	- [ ] Traffic Monitor
	- [ ] Traffic Router
	- [ ] Traffic Stats
	- [ ] Traffic Vault
	- [ ] ORT



