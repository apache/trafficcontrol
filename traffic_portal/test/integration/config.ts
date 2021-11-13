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
import { resolve } from "path";

import { emptyDir } from "fs-extra";
import { Config, browser } from 'protractor';
import { JUnitXmlReporter } from 'jasmine-reporters';
import HtmlReporter from "protractor-beautiful-reporter";

import { API } from './CommonUtils';
import * as conf from "./config.json"
import { prerequisites } from "./Data";
import { isTestingConfig } from "./config.model";

const downloadsPath = resolve('Downloads');
export const randomize = Math.random().toString(36).substring(3, 7);
export const twoNumberRandomize = Math.floor(Math.random() * 101);

export let config: Config = conf;
if (config.capabilities) {
  config.capabilities.chromeOptions.prefs.download.default_directory = downloadsPath;
} else {
  config.capabilities = {chromeOptions: {prefs: {download: {default_directory: downloadsPath}}}};
}

if (!config.params) {
  throw new Error("no testing parameters provided - cannot proceed");
}

try {
  if (!isTestingConfig(config.params)) {
    throw new Error();
  }
} catch (e) {
  const msg = e instanceof Error ? e.message : String(e);
  throw new Error(`invalid testing params: ${msg}`);
}

export const testingConfig = config.params;
export const api = new API(testingConfig);

config.onPrepare = async function () {
    await browser.waitForAngularEnabled(true);
    await browser.driver.manage().window().maximize();
    emptyDir('./Reports/', function (err) {
      console.log(err);
    });

    if (config.params.junitReporter === true) {
        jasmine.getEnv().addReporter(
            new JUnitXmlReporter({
                savePath: '/portaltestresults',
                filePrefix: 'portaltestresults',
                consolidateAll: true
            }));
    }
    else {
        jasmine.getEnv().addReporter(new HtmlReporter({
            baseDirectory: './Reports/',
            clientDefaults: {
                showTotalDurationIn: "header",
                totalDurationFormat: "hms"
            },
            jsonsSubfolder: 'jsons',
            screenshotsSubfolder: 'images',
            takeScreenShotsOnlyForFailedSpecs: true,
            docTitle: 'Traffic Portal Test Cases'
        }).getJasmine2Reporter());
    }

    await api.UseAPI(prerequisites);
}
