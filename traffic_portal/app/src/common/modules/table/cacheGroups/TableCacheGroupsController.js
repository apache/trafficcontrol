/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/**
 * @typedef CacheGroup
 * @property {number} id
 * @property {string} name
 * @property {number} shortName
 * @property {number} latitude
 * @property {number} longitude
 * @property {string} parentCachegroupName
 * @property {string} secondaryParentCachegroupName
 * @property {string} typeName
 * @property {string} lastUpdated
 */

/**
 * @param {CacheGroup} cacheGroup
 * @returns  {string}
 */
const getHref = (cacheGroup) => `#!/cache-groups/${cacheGroup.id}`;

/**
 * @param {CacheGroup[]} cacheGroups
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/CacheGroupService")} cacheGroupService
 * @param {import("../../../models/MessageModel")} messageModel
 */
var TableCacheGroupsController = function (
    cacheGroups,
    $scope,
    $state,
    $uibModal,
    locationUtils,
    cacheGroupService,
    messageModel
) {
    /**** Constants, scope data, etc. ****/

    /** The columns of the ag-grid table */
    $scope.columns = [
        {
            headerName: "Name",
            field: "name",
            hide: false,
        },
        {
            headerName: "Short Name",
            field: "shortName",
            hide: false,
        },
        {
            headerName: "Type",
            field: "typeName",
            hide: false,
        },
        {
            headerName: "1st Parent",
            field: "parentCachegroupName",
            hide: false,
        },
        {
            headerName: "2nd Parent",
            field: "secondaryParentCachegroupName",
            hide: false,
        },
        {
            headerName: "Latitude",
            field: "latitude",
            hide: false,
        },
        {
            headerName: "Longitude",
            field: "longitude",
            hide: false,
        },
        {
            headerName: "ID",
            field: "id",
            filter: "agNumberColumnFilter",
            hide: true,
        },
        {
            headerName: "Last Updated",
            field: "lastUpdated",
            hide: true,
            filter: "agDateColumnFilter",
        },
    ];

    /** @type {import("../agGrid/CommonGridController").CGC.DropDownOption[]} */
    $scope.dropDownOptions = [
        {
            name: "createCacheGroupMenuItem",
            href: "#!/cache-groups/new",
            text: "Create New Cache Group",
            type: 2,
        },
    ];

    /** Reloads all resolved data for the view. */
    $scope.refresh = () => {
        $state.reload();
    };

    /**
     * Deletes a Cache Group if confirmation is given.
     * @param {CacheGroup} cacheGroup
     */
    function confirmDelete(cacheGroup) {
        const params = {
            title: `Delete Cache Group: ${cacheGroup.name}`,
            key: cacheGroup.name,
        };
        const modalInstance = $uibModal.open({
            templateUrl: "common/modules/dialog/delete/dialog.delete.tpl.html",
            controller: "DialogDeleteController",
            size: "md",
            resolve: { params },
        });
        modalInstance.result
            .then(() => {
                cacheGroupService
                    .deleteCacheGroup(cacheGroup.id)
                    .then((result) => {
                        messageModel.setMessages(result.alerts, false);
                        $scope.refresh();
                    });
            })
            .catch((e) => console.error("failed to delete Cache Group:", e));
    }

    /**
     * Queues servers updates on a Cache Group if CDN is selected
     * @param {CacheGroup} cacheGroup
     */
    function confirmQueueServerUpdates(cacheGroup) {
        const params = {
            title: `Queue Server Updates: ${cacheGroup.name}`,
            message: "Please select a CDN",
        };
        const modalInstance = $uibModal.open({
            templateUrl: "common/modules/dialog/select/dialog.select.tpl.html",
            controller: "DialogSelectController",
            size: "md",
            resolve: {
                params,
                collection: (cdnService) => cdnService.getCDNs(),
            },
        });
        modalInstance.result.then((cdn) =>
            cacheGroupService.queueServerUpdates(cacheGroup.id, cdn.id)
        );
    }

    /**
     * Clears servers updates on a Cache Group if confirmation is given.
     * @param {CacheGroup} cacheGroup
     */
    function confirmClearServerUpdates(cacheGroup) {
        const params = {
            title: `Clear Server Updates: ${cacheGroup.name}`,
            message: "Please select a CDN",
        };
        const modalInstance = $uibModal.open({
            templateUrl: "common/modules/dialog/select/dialog.select.tpl.html",
            controller: "DialogSelectController",
            size: "md",
            resolve: {
                params,
                collection: (cdnService) => cdnService.getCDNs(),
            },
        });
        modalInstance.result.then((cdn) => {
            cacheGroupService.clearServerUpdates(cacheGroup.id, cdn.id);
        });
    }

    /** @type {import("../agGrid/CommonGridController").CGC.ContextMenuOption[]} */
    $scope.contextMenuOptions = [
        {
            getHref,
            getText: (cacheGroup) => `Open ${cacheGroup.name} in a new tab`,
            newTab: true,
            type: 2,
        },
        { type: 0 },
        {
            getHref,
            text: "Edit",
            type: 2,
        },
        {
            onClick: (cacheGroup) => confirmDelete(cacheGroup),
            text: "Delete",
            type: 1,
        },
        { type: 0 },
        {
            onClick: (cacheGroup) => confirmQueueServerUpdates(cacheGroup),
            text: "Queue Server Updates",
            type: 1,
        },
        {
            onClick: (cacheGroup) => confirmClearServerUpdates(cacheGroup),
            text: "Clear Server Updates",
            type: 1,
        },
        { type: 0 },
        {
            getHref: (cacheGroup) => `#!/cache-groups/${cacheGroup.id}/asns`,
            text: "Manage ASNs",
            type: 2,
        },
        {
            getHref: (cacheGroup) => `#!/cache-groups/${cacheGroup.id}/servers`,
            text: "Manage Servers",
            type: 2,
        },
    ];

    /** Options, configuration, data and callbacks for the ag-grid table. */
    /** @type {import("../agGrid/CommonGridController").CGC.GridSettings} */
    $scope.gridOptions = {
        onRowClick: function (row) {
            locationUtils.navigateToPath(`/cache-groups/${row.data.id}`);
        },
    };

    $scope.cacheGroups = cacheGroups.map((cacheGroup) => ({
        ...cacheGroup,
        lastUpdated: new Date(
            cacheGroup.lastUpdated.replace(" ", "T").replace("+00", "Z")
        ),
    }));
};

TableCacheGroupsController.$inject = [
    "cacheGroups",
    "$scope",
    "$state",
    "$uibModal",
    "locationUtils",
    "cacheGroupService",
    "messageModel",
];
module.exports = TableCacheGroupsController;
