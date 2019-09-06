<!--
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
-->
# Configurable Table Columns

## Problem Description
As an operations user, I need the ability to override the default order and visibility of table columns in Traffic Portal and persist the overridden settings for future sessions. In addition, the values of hidden columns should be excluded from the global table search.

When viewing a collection of resources (servers, delivery services, etc.), an operations user may require knowledge of different resource attribute values depending on the current task. By allowing an operations user to configure the visibility and order of TP table columns, an operations user can retrieve the necessary information in a timely fashion. In addition, by allowing an operations user to hide columns and exclude the value of those columns from a table search, it provides a more powerful and targeted search mechanism.

## Proposed Change

Add the ability to reorder TP table columns using drag and drop and provide the user with the ability to toggle on/off the visibility of the column. Column visibility will determine if a column is included/excluded from the global search. Any changes made by the user shall be persisted.

### Traffic Portal Impact

 - Add drag and drop capabilities to each table column so the user can reorder the columns to fit their needs. There is a jquery datatables extension (ColReorder) that could accomplish this.
 - Add a drop down with a check box for each table column so the user can toggle the visibility of the column. Column visibility can be controlled using the API of jquery datatables.
 - Exclude a column from the global table search if it is not visible. Column searchability cannot be controlled using the API of jquery datatables so the jquery datatables library may need to be patched to accomplish this.
 - Leverage the browser's local storage to persist the user-defined order and visibility of the table columns.

### Traffic Ops Impact

N/A

#### REST API Impact

N/A

#### Client Impact

N/A

#### Data Model / Database Impact

N/A

### ORT Impact

N/A

### Traffic Monitor Impact

N/A

### Traffic Router Impact

N/A

### Traffic Stats Impact

N/A

### Traffic Vault Impact

N/A

### Documentation Impact

This new functionality needs to be documented in the "Using Traffic Portal" section of the documentation and any existing screenshots of TP tables need to be modified.

### Testing Impact

UI tests can be added to test the proposed behaviour. It is probably overkill to test this functionality on all tables but one or two should suffice.

### Performance Impact

N/A

### Security Impact

N/A

### Upgrade Impact

This functionality is not dependent on the upgrade of any other TC component. TP can be upgraded or rolled back as needed.

### Operations Impact

These changes will allow TP users to get to relevant resource data faster. However, some training or a demo of the new functionality is required.

### Developer Impact

N/A

## Alternatives

N/A

## Dependencies

The following client-side dependency is required:
- https://datatables.net/extensions/colreorder/

## References

- https://datatables.net/extensions/colreorder/
