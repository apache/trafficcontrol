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
import { prerequisites, LoadCacheGroup, LoadProfile } from "./prerequisites";
import type { CacheGroup, CDN, Profile, Role, Tenant, Type, User } from "./model";
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

const cdns = new Map<string, CDN>();
const tenants = new Map<string, Tenant>();
const createdTenants = new Array<Tenant>();
const types = new Map<string, Type>();
const roles = new Map<string, Role>();
const cacheGroups = new Array<CacheGroup>();
const profiles = new Array<Profile>();
const users = new Array<User>();


/** This is used to prevent trying forever to load Tenants with broken parentage. */
let lastUnloaded = new Set<string>();

async function loadTenants(ts: {name: string; parent: string;}[], loaded: Map<string, Tenant>): Promise<void> {
	if (ts.length === 0) {
		return;
	}

	const unloaded = [];
	const promises = new Array<Promise<void>>();
	for (const t of ts) {
		const parent = loaded.get(t.parent);
		if (parent) {
			ts.filter(ten=>ten.parent === t.name).map(
				ten => {
					ten.parent += randomize;
				}
			);
			t.name += randomize
			const payload = {
				active: true,
				name: t.name,
				parentId: parent.id
			};
			promises.push(api.post<Tenant>("tenants", payload).then(
				createdTenant => {
					loaded.set(createdTenant.name, createdTenant);
					createdTenants.push(createdTenant);
				}
			));
		} else {
			unloaded.push(t);
		}
	}

	await Promise.all(promises);

	if (unloaded.length === lastUnloaded.size && unloaded.every(t=>lastUnloaded.has(t.name))) {
		const msg = unloaded.map(t=>`'${t.name}'`).join(", ");
		throw new Error(`the following Tenants cannot be loaded because their parents don't already exist and aren't in the testing data: ${msg}`);
	}

	lastUnloaded = new Set(unloaded.map(t=>t.name));
	return loadTenants(unloaded, loaded);
}


async function loadCDN(cdn: {domainName: string; name: string}): Promise<void> {
	cdn.name += randomize;
	cdn.domainName += randomize;
	return api.post<CDN>("cdns", {...cdn, dnssecEnabled: false}).then(
		c => {
			cdns.set(c.name, c);
		}
	);
}

async function loadCacheGroup(cg: LoadCacheGroup): Promise<void> {
	const type = types.get(cg.type);
	if (!type) {
		throw new Error(`Cache Group '${cg.name}' cannot be created because its Type '${cg.type}' doesn't exist`);
	}
	if (type.useInTable !== "cachegroup") {
		throw new Error(`Cache Group '${cg.name}' cannot be created because its Type '${cg.type}' is used in the wrong table: ${type.useInTable}`);
	}
	cg.name += randomize;
	cg.shortName += randomize;
	return api.post<CacheGroup>("cachegroups", {...cg, typeId: type.id}).then(
		created => {
			cacheGroups.push(created);
		}
	);
}

async function loadProfile(prof: LoadProfile): Promise<void> {
	const cdn = cdns.get(prof.cdn+randomize);
	if (!cdn) {
		throw new Error(`Profile '${prof.name} could not be created because the CDN it claims to be in, '${prof.cdn}' doesn't exist and isn't in the testing data`);
	}
	const payload = {
		cdn: cdn.id,
		description: prof.description,
		name: prof.name+randomize,
		routingDisabled: prof.routingDisabled,
		type: prof.type
	}
	return api.post<Profile>("profiles", payload).then(
		p => {
			profiles.push(p);
		}
	);
}

async function loadUser(user: {role: string; tenant: string; username: string}): Promise<void> {
	const role = roles.get(user.role);
	if (!role) {
		throw new Error("User '${user.username}' could not be created because its Role '${user.role}' doesn't exist");
	}
	const tenantName = user.tenant;
	let tenant = tenants.get(tenantName);
	if (!tenant) {
		tenant = tenants.get(tenantName+randomize);
		user.tenant += randomize;
	}
	if (!tenant) {
		throw new Error("User '${user.username}' could not be created because its Tenant '${tenantName}' doesn't exist and isn't in the testing data");
	}

	user.username += randomize;
	const payload = {
		...user,
		confirmLocalPasswd: "pa$$word",
		fullName: user.username,
		email: `${user.username}@tp-tests.test`,
		localPasswd: "pa$$word",
		role: role.id,
		tenantId: tenant.id
	};
	return api.post<User>("users", payload).then(
		u => {
			users.push(u);
		}
	);
}

async function setupAPI(): Promise<void> {
	console.debug("loading testing data");
	console.time("testing dataset loaded");
	if (!api.loggedIn) {
		throw new Error("cannot setup before logging in");
	}


	let promises: Array<Promise<void>> = prerequisites.CDNs.map(loadCDN);
	promises.push(api.get<Tenant[]>("tenants").then(
		ts => {
			let foundRoot = false;
			for (const t of ts) {
				tenants.set(t.name, t);
				if (!foundRoot && t.name === "root") {
					foundRoot = true;
				}
			}
			if (!foundRoot) {
				throw new Error("'root' tenant not found in TO - this means the configured login is not for a root user, and so the tests can't load");
			}
		}
	));
	promises.push(api.get<Type[]>("types").then(
		ts => {
			for (const t of ts) {
				types.set(t.name, t);
			}
		}
	));
	promises.push(api.get<Role[]>("roles").then(
		rs => {
			for (const r of rs) {
				roles.set(r.name, r);
			}
		}
	));

	await Promise.all(promises);

	promises = [loadTenants(prerequisites.tenants, tenants)];
	for (const role of prerequisites.roles) {
		const r = roles.get(role.name);
		if (!r) {
			const payload = {
				...role,
				description: `the '${role.name}' Role used in TP tests`
			};
			promises.push(api.post<Role>("roles", payload).then(
				created => {
					roles.set(created.name, created);
				}
			));
		} else if (r.privLevel < role.privLevel) {
			throw new Error(`the tests need Role '${role.name}' to have a privLevel of at least ${role.privLevel} to run, but it only has ${r.privLevel}`);
		}
	}
	promises = promises.concat(promises, prerequisites.cacheGroups.map(loadCacheGroup));
	promises = promises.concat(promises, prerequisites.profiles.map(loadProfile));

	await Promise.all(promises);
	await Promise.all(prerequisites.users.map(loadUser));
	console.timeEnd("testing dataset loaded");
}

config.beforeLaunch = async function () {
	await api.Login();
	await setupAPI();
};

config.onPrepare = async function () {
    await browser.waitForAngularEnabled(true);
    await browser.driver.manage().window().maximize();
    emptyDir('./Reports/', function (err) {
      console.log(err);
    });

	await browser.waitForAngularEnabled(true);
	await browser.driver.manage().window().maximize();
	emptyDir('./Reports/', function (err) {
		console.error(err);
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
};

function leafTenants(ts: Array<Tenant & {parentId: number}>): Array<Tenant & {parentId: number}> {
	const parents = new Set(ts.map(t=>t.parentName));
	return ts.filter(t=>!parents.has(t.name));
}

function handleErr(thing: string, id: number): (e: unknown) => void {
	const tmpl = `failed to clean up ${thing} #${id}:`;
	return (e: unknown) => {
		const msg = e instanceof Error ? e.message : e;
		console.error(tmpl, msg);
	};
}

async function cleanUpTenants(ts: Array<Tenant & {parentId: number}>): Promise<void> {
	if (ts.length === 0) {
		return;
	}

	const leaves = leafTenants(ts);
	if (leaves.length === 0) {
		const msg = leaves.map(t=>`'${t.name.replace(new RegExp(`${randomize}$`), "")}'`).join(", ");
		throw new Error(`The following Tenants could not be cleaned up because of bad parentage: ${msg}`);
	}

	const leafIDs = new Set<number>();
	await Promise.all(leaves.map(
		t => {
			leafIDs.add(t.id);
			return api.delete(`tenants/${t.id}`).catch(handleErr("Tenant", t.id));
		}
	));

	const unDeleted = ts.filter(t=>!leafIDs.has(t.id));
	return cleanUpTenants(unDeleted);
}

const teardownTimingLabel = "testing data torn down";

async function teardownAPI(): Promise<void> {
	console.debug("tearing down testing data");
	// TODO: delete when users can be deleted.
	console.warn("users cannot be cleaned up, so they will be left in the testing environment");
	console.time(teardownTimingLabel);

	let promises = profiles.map(p=>api.delete(`profiles/${p.id}`).catch(handleErr("Profile", p.id)));
	promises = promises.concat(promises, cacheGroups.map(
		cg=>api.delete(`cachegroups/${cg.id}`).catch(handleErr("Cache Group", cg.id))
	));
	// TODO: uncomment once users can be deleted.
	// promises = promises.concat(promises, users.map(
	// 	u=>api.delete(`users/${u.id}`).catch(handleErr("User", u.id))
	// ));
	await Promise.all(promises);

	promises = Array.from(cdns.values()).map(
		c=>api.delete(`cdns/${c.id}`).catch(handleErr("CDN", c.id))
	);
	const noRoots = createdTenants.filter(
		(t): t is Tenant & {parentId: number} => typeof(t.parentId) === "number"
	);
	if (noRoots.length < createdTenants.length) {
		console.warn("refusing to delete the 'root' Tenant");
	}
	promises.push(cleanUpTenants(noRoots).catch(
		err => {
			console.error("Failed to clean up Tenants:", err instanceof Error ? err.message : err);
		}
	));
	await Promise.all(promises);

	console.timeEnd(teardownTimingLabel);
}


config.afterLaunch = teardownAPI;
