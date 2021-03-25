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
import { API } from './CommonUtils/API';
import { Config, browser } from 'protractor'
import * as conf from "./config.json"

let path = require('path');
let downloadsPath = path.resolve('Downloads');
let randomize = Math.random().toString(36).substring(3, 7);
let HtmlReporter = require('protractor-beautiful-reporter');
let twoNumberRandomize = Math.floor(Math.random() * 101);
exports.twoNumberRandomize = twoNumberRandomize;
exports.randomize = randomize;

export let config: Config = conf;
config.capabilities.chromeOptions.prefs.download.default_directory = downloadsPath;
config.onPrepare = async function () {
    await browser.waitForAngularEnabled(true);

    var fs = require('fs-extra');

    fs.emptyDir('./Reports/', function (err) {
      console.log(err);
    });

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

    try {
      let api = new API();
      let setupFile = 'Data/Prerequisites/user.setup.json';
      let setupData = JSON.parse(fs.readFileSync(setupFile));
      let output = await api.UseAPI(setupData);
      if (output != null){
        throw new Error(output)
      }
    } catch (error) {
      throw error
    }
}
