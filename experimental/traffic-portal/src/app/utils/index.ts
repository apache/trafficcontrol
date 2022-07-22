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

export * from "./order-by";
export * from "./fuzzy";
export * from "./ip";
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
