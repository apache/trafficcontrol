/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

export * from "./fuzzy";
export * from "./ip";
export * from "./logging";
export * from "./order-by";
export * from "./time";

/**
 * These are the values that may be given to the `autocomplete` attribute of an
 * HTML form control.
 */
export const enum AutocompleteValue {
	// Disable/enable with no guidance.
	/**
	 * The browser is not permitted to automatically enter or select a value for
	 * this field. It is possible that the document or application provides its
	 * own autocomplete feature, or that security concerns require that the
	 * field's value not be automatically entered.
	 *
	 * Note that in most modern browsers, setting autocomplete to "off" will not
	 * prevent a password manager from asking the user if they would like to
	 * save username and password information, or from automatically filling in
	 * those values in a site's login form.
	 */
	OFF = "off",
	/**
	 * The browser is allowed to automatically complete the input. No guidance
	 * is provided as to the type of data expected in the field, so the browser
	 * may use its own judgement.
	 */
	ON = "on",

	// Misc.
	/** A preferred language, given as a valid BCP 47 language tag. */
	LANGUAGE = "language",
	/**
	 * The URL of an image representing the person, company, or contact
	 * information given in the other fields in the form.
	 */
	PHOTO = "photo",
	/**
	 * A URL, such as a home page or company web site address as appropriate
	 * given the context of the other fields in the form.
	 */
	URL = "url",

	// Physical addresses
	/**
	 * A street address. This can be multiple lines of text, and should fully
	 * identify the location of the address within its second administrative
	 * level (typically a city or town), but should not include the city name,
	 * ZIP or postal code, or country name.
	 */
	STREET_ADDRESS = "street-address",
	/**
	 * The first line of the street address. This should only be present if the
	 * {@link AutocompleteValue.STREET_ADDRESS} is not present.
	 */
	ADDRESS_LINE_1 = "address-line-1",
	/**
	 * The second line of the street address. This should only be present if the
	 * {@link AutocompleteValue.STREET_ADDRESS} is not present.
	 */
	ADDRESS_LINE_2 = "address-line-2",
	/**
	 * The third line of the street address. This should only be present if the
	 * {@link AutocompleteValue.STREET_ADDRESS} is not present.
	 */
	ADDRESS_LINE_3 = "address-line-3",
	/**
	 * The first administrative level in the address. This is typically the
	 * province in which the address is located. In the United States, this
	 * would be the state. In Switzerland, the canton. In the United Kingdom,
	 * the post town.
	 */
	ADDRESS_LEVEL_1 = "address-level-1",
	/**
	 * The second administrative level, in addresses with at least two of them.
	 * In countries with two administrative levels, this would typically be the
	 * city, town, village, or other locality in which the address is located.
	 */
	ADDRESS_LEVEL_2 = "address-level-2",
	/**
	 * The third administrative level, in addresses with at least three
	 * administrative levels.
	 */
	ADDRESS_LEVEL_3 = "address-level-3",
	/**
	 * The finest-grained administrative level, in addresses which have four
	 * levels.
	 */
	ADDRESS_LEVEL_4 = "address-level-4",
	/** A country or territory code. */
	COUNTRY = "country",
	/** A country or territory name. */
	COUNTRY_NAME = "country-name",
	/** A postal code (in the United States, this is the ZIP code). */
	POSTAL_CODE = "postal code",

	// Credit card information
	/**
	 * The full name as printed on or associated with a payment instrument such
	 * as a credit card. Using a full name field is preferred, typically, over
	 * breaking the name into pieces.
	 */
	CREDIT_CARD_NAME = "cc-name",
	/**
	 * A given (first) name as given on a payment instrument like a credit
	 * card.
	 */
	CREDIT_CARD_GIVEN_NAME = "cc-given-name",
	/** A middle name as given on a payment instrument or credit card. */
	CREDIT_CARD_ADDITIONAL_NAME = "cc-additional-name",
	/** A family name, as given on a credit card. */
	CREDIT_CARD_FAMILY_NAME = "cc-family-name",
	/**
	 * A credit card number or other number identifying a payment method, such
	 * as an account number.
	 */
	CREDIT_CARD_NUMBER = "cc-number",
	/**
	 * A payment method expiration date, typically in the form "MM/YY" or
	 * "MM/YYYY".
	 */
	CREDIT_CARD_EXP = "cc-exp",
	/** The month in which the payment method expires. */
	CREDIT_CARD_EXP_MONTH = "cc-exp-month",
	/** The year in which the payment method expires. */
	CREDIT_CARD_EXP_YEAR = "cc-exp-year",
	/**
	 * The security code for the payment instrument; on credit cards, this is
	 * the 3-digit verification number on the back of the card.
	 */
	CREDIT_CARD_CSC = "cc-csc",
	/** The type of payment instrument (such as "Visa" or "Master Card"). */
	CREDIT_CARD_TYPE = "cc-type",

	// Contact information
	/**
	 * A full telephone number, including the country code. If you need to break
	 * the phone number up into its components, you can use
	 * {@link AutocompleteValue.TELEPHONE_COUNTRY_CODE},
	 * {@link AutocompleteValue.TELEPHONE_NATIONAL},
	 * {@link AutocompleteValue.TELEPHONE_AREA_CODE},
	 * {@link AutocompleteValue.TELEPHONE_LOCAL}, and/or
	 * {@link AutocompleteValue.TELEPHONE_EXTENSION} for those fields.
	 */
	TELEPHONE = "tel",
	/**
	 * The country code, such as "1" for the United States, Canada, and other
	 * areas in North America and parts of the Caribbean.
	 */
	TELEPHONE_COUNTRY_CODE = "tel-country-code",
	/**
	 * The entire phone number without the country code component, including a
	 * country-internal prefix. For the phone number "1-855-555-6502", this
	 * field's value would be "855-555-6502".
	 */
	TELEPHONE_NATIONAL = "tel-national",
	/**
	 * The area code, with any country-internal prefix applied if appropriate.
	 */
	TELEPHONE_AREA_CODE = "tel-area-code",
	/**
	 * The phone number without the country or area code. This can be split
	 * further into two parts, for phone numbers which have an exchange number
	 * and then a number within the exchange. For the phone number "555-6502",
	 * use {@link AutocompleteValue.TELEPHONE_LOCAL_PREFIX} for "555" and
	 * {@link AutocompleteValue.TELEPHONE_LOCAL_SUFFIX} for "6502".
	 */
	TELEPHONE_LOCAL = "tel-local",
	/**
	 * The exchange number part of a telephone number without a country or area
	 * code.
	 */
	TELEPHONE_LOCAL_PREFIX = "tel-local-prefix",
	/**
	 * A number within an exchange for a telephone number without a country or
	 * area code.
	 */
	TELEPHONE_LOCAL_SUFFIX = "tel-local-suffix",
	/**
	 * A telephone extension code within the phone number, such as a room or
	 * suite number in a hotel or an office extension in a company.
	 */
	TELEPHONE_EXTENSION = "tel-extension",
	/** A URL for an instant messaging protocol endpoint, such as "xmpp:username@example.net". */
	IMPP = "impp",
	/** An email address */
	EMAIL = "email",

	// Identification/authentication
	/** A username or account name. */
	USERNAME = "username",
	/**
	 * A new password. When creating a new account or changing passwords, this
	 * should be used for an "Enter your new password" or "Confirm new password"
	 * field, as opposed to a general "Enter your current password" field that
	 * might be present. This may be used by the browser both to avoid
	 * accidentally filling in an existing password and to offer assistance in
	 * creating a secure password
	 */
	NEW_PASSWORD = "new-password",
	/** The user's current password. */
	CURRENT_PASSWORD = "current-password",
	/** A one-time code used for verifying user identity. */
	ONE_TIME_CODE = "one-time-code",
	/**
	 * A job title, or the title a person has within an organization, such as
	 * "Senior Technical Writer", "President", or "Assistant Troop Leader".
	 */
	ORGANIZATION_TITLE = "organization-title",
	/**
	 * A company or organization name, such as "Acme Widget Company" or "Girl
	 * Scouts of America".
	 */
	ORGANIZATION = "organization",
	/** A birth date, as a full date. */
	BIRTHDAY = "bday",
	/** The day of the month of a birth date. */
	BIRTHDAY_DAY = "bday-day",
	/** The month of the year of a birth date. */
	BIRTHDAY_MONTH = "bday-month",
	/** The year of a birth date. */
	BIRTHDAY_YEAR = "bday-year",
	/**
	 * A gender identity (such as "Female", "Fa'afafine", "Male"), as freeform
	 * text without newlines.
	 */
	SEX = "sex",
	/**
	 * A gender identity (such as "Female", "Fa'afafine", "Male"), as freeform
	 * text without newlines.
	 *
	 * Alias of "sex", because the field is actually meant to more generally
	 * express gender than biological sex.
	 */
	GENDER = "sex",
	/**
	 * The field expects the value to be a person's full name. Using "name"
	 * rather than breaking the name down into its components is generally
	 * preferred because it avoids dealing with the wide diversity of human
	 * names and how they are structured; however, you can use the following
	 * autocomplete values if you do need to break the name down into its
	 * components:
	 * - {@link AutocompleteValue.HONORIFIC_PREFIX}
	 * - {@link AutocompleteValue.GIVEN_NAME}
	 * - {@link AutocompleteValue.ADDITIONAL_NAME}
	 * - {@link AutocompleteValue.FAMILY_NAME}
	 * - {@link AutocompleteValue.HONORIFIC_SUFFIX}
	 * - {@link AutocompleteValue.NICKNAME}
	 */
	NAME = "name",

	// Typically, `NAME` should be used instead of any of these, to avoid having
	// to deal with the vast variety of human naming customs. Nevertheless, they
	// are valid `autocomplete` attribute values.
	/**
	 * The prefix or title, such as "Mrs.", "Mr.", "Miss", "Ms.", "Dr.", or
	 * "Mlle.".
	 */
	HONORIFIC_PREFIX = "honorific-prefix",
	/** The given (or "first") name. */
	GIVEN_NAME = "given-name",
	/** The middle name. */
	ADDITIONAL_NAME = "additional-name",
	/** The family (or "last") name. */
	FAMILY_NAME = "family-name",
	/** The suffix, such as "Jr.", "B.Sc.", "PhD.", "MBASW", or "IV". */
	HONORIFIC_SUFFIX = "honorific-suffix",
	/** A nickname or handle. */
	NICKNAME = "nickname",
}

/**
 * Checks if an object is a string. Useful for passing as a type guard into
 * generic functions. In general, you should just use `typeof` instead.
 *
 * @param x The object to check.
 * @returns `true` if `x` is a string, `false` otherwise.
 */
export function isString(x: unknown): x is string {
	return typeof(x) === "string";
}

/**
 * Checks if an object is a number. Useful for passing as a type guard into
 * generic functions. In general, you should just use `typeof` instead.
 *
 * @param x The object to check.
 * @returns `true` if `x` is a number, `false` otherwise.
 */
export function isNumber(x: unknown): x is number {
	return typeof(x) === "number";
}

/**
 * Checks if an object is a boolean. Useful for passing as a type guard into
 * generic functions. In general, you should just use `typeof` instead.
 *
 * @param x The object to check.
 * @returns `true` if `x` is a boolean, `false` otherwise.
 */
export function isBoolean(x: unknown): x is boolean {
	return typeof(x) === "boolean";
}

/**
 * isRecord checks that the passed object is a Record (object).
 *
 * @param x The object to check.
 * @returns `true` if `x` is a Record, `false` otherwise.
 */
export function isRecord(x: unknown): x is Record<PropertyKey, unknown>;
/**
 * isRecord checks that the passed object is a Record with string property
 * values.
 *
 * @param x The object to check.
 * @param type Indicates we are checking for string properties.
 * @returns `true` if `x` is a Record of properties with string values, `false`
 * otherwise.
 */
export function isRecord(x: unknown, type: "string"): x is Record<PropertyKey, string>;
/**
 * isRecord checks that the passed object is a Record with numeric property
 * values.
 *
 * @param x The object to check.
 * @param type Indicates we are checking for number properties.
 * @returns `true` if `x` is a Record of properties with numeric values, `false`
 * otherwise.
 */
export function isRecord(x: unknown, type: "number"): x is Record<PropertyKey, number>;
/**
 * isRecord checks that the passed object is a Record with a certain, homogenous
 * type.
 *
 * @param x The object to check.
 * @param checker A type guard that will be used to ensure each property is of
 * the type for which it checks.
 * @returns `true` if `x` is a Record of properties with values that satisfy
 * `checker`, `false` otherwise.
 */
export function isRecord<T>(x: unknown, checker: (p: unknown) => p is T): x is Record<PropertyKey, T>;
/**
 * isRecord checks that the passed object is a Record, optionally of a certain,
 * homogenous type.
 *
 * @param x The object to check.
 * @param checker Either the name of a primitive type to check, or a custom
 * type checker function to use to check the type of each property value. If not
 * given, the values of properties is not verified.
 * @returns `true` if `x` is a Record of properties with values optionally given
 * by `checker`.
 */
export function isRecord<T>(x: unknown, checker?: "string" | "number" | ((p: unknown) => p is T)): x is Record<PropertyKey, T> {
	if (typeof x !== "object" || x === null) {
		return false;
	}
	if (!checker) {
		return true;
	}
	let chk;
	switch (checker) {
		case "string":
			chk = isString;
			break;
		case "number":
			chk = isNumber;
			break;
		default:
			chk = checker;
	}
	return Object.values(x).every(chk);
}

/**
 * isArray checks if something is an Array. This call signature is provided for
 * generic completeness - it should basically never be used. Instead, use the
 * built-in `Array.isArray` method.
 *
 * @param a The possible array to check.
 * @returns `true` if `a` is any kind of Array, `false` otherwise.
 */
export function isArray(a: unknown): a is Array<unknown>;
/**
 * isArray checks if some object is an array of strings. This is exactly
 * equivalent to using {@link isStringArray}.
 *
 * @param a The possible array to check.
 * @param type Indicates we are checking for an array of strings.
 * @returns `true` if `a` is an array containing only strings, `false`
 * otherwise.
 */
export function isArray(a: unknown, type: "string"): a is Array<string>;
/**
 * isArray checks if some object is an array of numbers. This is exactly
 * equivalent to using {@link isNumberArray}.
 *
 * @param a The possible array to check.
 * @param type Indicates we are checking for an array of numbers.
 * @returns `true` if `a` is an array containing only numbers, `false`
 * otherwise.
 */
export function isArray(a: unknown, type: "number"): a is Array<number>;
/**
 * isArray checks if some object is a homogenous array of some specific type.
 *
 * @param a The possible array to check.
 * @param checker A type guard that will be used to verify that `a` - if it
 * indeed be an array - contains only types of data that satisfy the guard.
 * @returns `true` if `a` is an array containing only values of type `T`,
 * `false` otherwise.
 */
export function isArray<T>(a: unknown, checker: (x: unknown) => x is T): a is Array<T>;
/**
 * isArray checks if something is an Array, optionally with some homogenous
 * value.
 *
 * @param a The possible array to check.
 * @param checker If given, this enforces that `a` is a homogenous array. If
 * this is the name of a primitive, it automatically checks for that primitive.
 * More complex types require that this be a type guard in its own right, to
 * match against each element of `a`.
 * @returns `true` if `a` is an array, satisfying a homogeneity checker if one
 * is provided, `false` otherwise.
 */
export function isArray<T>(a: unknown, checker?: "string" | "number" | ((x: unknown) => x is T)): a is Array<T> {
	if (!Array.isArray(a)) {
		return false;
	}
	if (!checker) {
		return true;
	}

	let chk: (x: unknown) => x is T | string | number;
	switch (checker) {
		case "number":
			chk = isNumber;
			break;
		case "string":
			chk = isString;
			break;
		default:
			chk = checker;
	}

	return a.every(chk);
}

/**
 * isStringArray checks if the passed object is a homogeneous array of strings.
 *
 * @param sa The potential string array.
 * @returns `true` if `sa` is an array and contains only strings, `false`
 * otherwise.
 */
export const isStringArray = (sa: unknown): sa is Array<string> => isArray(sa, "string");
/**
 * isNumberArray checks if the passed object is a homogeneous array of numbers.
 *
 * @param na The potential numeric array.
 * @returns `true` if `na` is an array and contains only numbers, `false`
 * otherwise.
 */
export const isNumberArray = (na: unknown): na is Array<number> => isArray(na, "number");

/**
 * hasProperty checks, generically, whether some variable passed as `o` has the
 * property `k`, where `k` is of some specific type.
 *
 * @note This only works for 'object' sub-types. Other basic types have easier
 * checks.
 *
 * @example
 * hasProperty({}, "id", isNumber); // returns `false`.
 *
 * @example
 * const isNum: (x: unknown) => x is number = x => typeof x === number;
 * // returns `false`
 * hasProperty({wrong: "type"}, "wrong", isNum, "number");
 *
 * @param o The object to check.
 * @param k The key for which to check in the object.
 * @param s This type guard will be used to narrow the *type* of the property,
 * as well as its existence.
 * @returns `true` if o has a property `k` such that the value of `o.k`
 * satisfies `s`, `false` otherwise.
 * @throws {TypeError} when called improperly.
 */
export function hasProperty<T extends object, K extends PropertyKey, S>
(o: T, k: K, s: (x: unknown) => x is S): o is T & Record<K, S>;
/**
 * hasProperty checks, generically, whether some variable passed as `o` has the
 * property `k` such that `o.k` is a string.
 *
 * @note This only works for 'object' sub-types. Other basic types have easier
 * checks.
 *
 * @example
 * hasProperty({}, "id", "string"); // returns `false`.
 *
 * @example
 * hasProperty({wrongType: 5}, "wrongType", "string"); // returns `false`
 *
 * @param o The object to check.
 * @param k The key for which to check in the object.
 * @param s indicates that we are checking for a string value of `o.k`.
 * @returns `true` if o has a property `k` such that the value of `o.k` is a
 * string, `false` otherwise.
 * @throws {TypeError} when called improperly.
 */
export function hasProperty<T extends object, K extends PropertyKey>
(o: T, k: K, s: "string"): o is T & Record<K, string>;
/**
 * hasProperty checks, generically, whether some variable passed as `o` has the
 * property `k` such that `o.k` is a number.
 *
 * @note This only works for 'object' sub-types. Other basic types have easier
 * checks.
 *
 * @example
 * hasProperty({}, "id", "number"); // returns `false`.
 *
 * @example
 * hasProperty({wrong: "type"}, "wrong", "number"); // returns `false`
 *
 * @param o The object to check.
 * @param k The key for which to check in the object.
 * @param s indicates that we are checking for a numeric value of `o.k`.
 * @returns `true` if o has a property `k` such that the value of `o.k` is a
 * number, `false` otherwise.
 * @throws {TypeError} when called improperly.
 */
export function hasProperty<T extends object, K extends PropertyKey>
(o: T, k: K, s: "number"): o is T & Record<K, number>;
/**
 * hasProperty checks, generically, whether some variable passed as `o` has the
 * property `k` such that `o.k` is a boolean.
 *
 * @note This only works for 'object' sub-types. Other basic types have easier
 * checks.
 *
 * @example
 * hasProperty({}, "id", "boolean"); // returns `false`.
 *
 * @example
 * hasProperty({wrong: "type"}, "wrong", "boolean"); // returns `false`
 *
 * @param o The object to check.
 * @param k The key for which to check in the object.
 * @param s indicates that we are checking for a boolean value of `o.k`.
 * @returns `true` if o has a property `k` such that the value of `o.k` is a
 * boolean, `false` otherwise.
 * @throws {TypeError} when called improperly.
 */
export function hasProperty<T extends object, K extends PropertyKey>
(o: T, k: K, s: "boolean"): o is T & Record<K, boolean>;
/**
 * hasProperty checks, generically, whether some variable passed as `o` has the
 * property `k`.
 *
 * @note This only works for 'object' sub-types. Other basic types have easier
 * checks.
 *
 * @example
 * hasProperty({}, "id"); // returns `false`.
 *
 * @param o The object to check.
 * @param k The key for which to check in the object.
 * @returns `true` if o has a property `k`, `false` otherwise.
 * @throws {TypeError} when called improperly.
 */
export function hasProperty<T extends object, K extends PropertyKey>
(o: T, k: K): o is T & Record<K, unknown>;
/**
 * hasProperty checks, generically, whether some variable passed as `o` has the
 * property `k`. It can also optionally narrow the type of `o.k`.
 *
 * @note This only works for 'object' sub-types. Other basic types have easier
 * checks.
 *
 * @example
 * hasProperty({}, "id"); // returns `false`
 *
 * @example
 * const isNum: (x: unknown) => x is number = x => typeof x === number;
 * // returns `false`.
 * hasProperty({wrong: "type"}, "wrong", isNum);
 *
 * @param o The object to check.
 * @param k The key for which to check in the object.
 * @param s If provided, this type guard can be used to narrow the *type* of
 * the property, as well as its existence.
 * @returns `true` if `o` has a property `k` such that any provided `s` is
 * satisfied by `o.k`.
 * @throws {Error} when the type check fails.
 */
export function hasProperty<T extends object, K extends string | number, S = unknown>(
	o: T,
	k: K,
	s?: "string" | "number" | "boolean" | ((x: unknown) => x is S)
): o is T & Record<K, S> {
	if (o === null || !Object.prototype.hasOwnProperty.call(o, k)) {
		return false;
	}
	if (s) {
		const val = (o as Record<K, unknown>)[k];
		switch (s) {
			case "string":
			case "number":
			case "boolean":
				return typeof(val) === s;
		}
		return s(val);
	}
	return true;
}

/**
 * Checks if the input implements the ArrayBufferView interface. NodeJS has a
 * built-in for this, but that won't be available in the browser.
 *
 * @param x The object to check.
 * @returns `true` if `x` is a typed array or a DataView, `false` otherwise.
 */
export function isArrayBufferView(x: unknown): x is ArrayBufferView {
	if (!x || typeof(x) !== "object") {
		return false;
	}
	switch(x.constructor) {
		case Int8Array:
		case Uint8Array:
		case Uint8ClampedArray:
		case Int16Array:
		case Uint16Array:
		case Int32Array:
		case Uint32Array:
		case DataView:
			return true;
	}
	return false;
}
