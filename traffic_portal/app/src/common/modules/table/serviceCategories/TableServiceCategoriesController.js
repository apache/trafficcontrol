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
 *@typedef ServiceCategory
 * @property {string} name
 * @property {string} lastUpdated
 */

/**
 * @param {ServiceCategory} serviceCategories
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 */
var TableServiceCategoriesController = function (
    serviceCategories,
    $scope,
    $state,
    locationUtils
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
            headerName: "Last Updated",
            field: "lastUpdated",
            hide: true,
            filter: "agDateColumnFilter",
        },
    ];

    /** @type {import("../agGrid/CommonGridController").CGC.DropDownOption[]} */
    $scope.dropDownOptions = [
        {
            name: "createServiceCategoryMenuItem",
            href: "#!/service-categories/new",
            text: "Create New Service Category",
            type: 2,
        },
    ];

    /** Reloads all resolved data for the view. */
    $scope.refresh = () => {
        $state.reload();
    };

    /** Options, configuration, data and callbacks for the ag-grid table. */
    /** @type {import("../agGrid/CommonGridController").CGC.GridSettings} */
    $scope.gridOptions = {
        onRowClick: (row) => {
            locationUtils.navigateToPath(
                `/service-categories/edit?name=${encodeURIComponent(
                    row.data.name
                )}`
            );
        },
    };

    $scope.serviceCategories = serviceCategories.map((serviceCategory) => ({
        ...serviceCategory,
        lastUpdated: new Date(
            serviceCategory.lastUpdated.replace(" ", "T").replace("+00", "Z")
        ),
    }));
};

TableServiceCategoriesController.$inject = [
    "serviceCategories",
    "$scope",
    "$state",
    "locationUtils",
];
module.exports = TableServiceCategoriesController;
