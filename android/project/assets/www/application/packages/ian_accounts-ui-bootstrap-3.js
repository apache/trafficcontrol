//////////////////////////////////////////////////////////////////////////
//                                                                      //
// This is a generated file. You can view the original                  //
// source in your browser if your browser supports source maps.         //
// Source maps are supported by all recent versions of Chrome, Safari,  //
// and Firefox, and by Internet Explorer 11.                            //
//                                                                      //
//////////////////////////////////////////////////////////////////////////


(function () {

/* Imports */
var Meteor = Package.meteor.Meteor;
var global = Package.meteor.global;
var meteorEnv = Package.meteor.meteorEnv;
var Session = Package.session.Session;
var Spacebars = Package.spacebars.Spacebars;
var Accounts = Package['accounts-base'].Accounts;
var _ = Package.underscore._;
var Template = Package.templating.Template;
var i18n = Package['anti:i18n'].i18n;
var Blaze = Package.blaze.Blaze;
var UI = Package.blaze.UI;
var Handlebars = Package.blaze.Handlebars;
var HTML = Package.htmljs.HTML;

/* Package-scope variables */
var ptPT, ptBR, zhCN, zhTW, srCyrl, srLatn, accountsUIBootstrap3, i, $modal;

(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/accounts_ui.js                                                                //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
if (!Accounts.ui){                                                                                                    // 1
	Accounts.ui = {};                                                                                                    // 2
}                                                                                                                     // 3
                                                                                                                      // 4
if (!Accounts.ui._options) {                                                                                          // 5
	Accounts.ui._options = {                                                                                             // 6
		extraSignupFields: [],                                                                                              // 7
		requestPermissions: {},                                                                                             // 8
		requestOfflineToken: {},                                                                                            // 9
		forceApprovalPrompt: {},                                                                                            // 10
		forceEmailLowercase: false,                                                                                         // 11
		forceUsernameLowercase: false,                                                                                      // 12
		forcePasswordLowercase: false                                                                                       // 13
	};                                                                                                                   // 14
}                                                                                                                     // 15
                                                                                                                      // 16
Accounts.ui.navigate = function (route, hash) {                                                                       // 17
	// if router is iron-router                                                                                          // 18
	if (window.Router && _.isFunction(Router.go)) {                                                                      // 19
		Router.go(route, hash);                                                                                             // 20
	}                                                                                                                    // 21
}                                                                                                                     // 22
                                                                                                                      // 23
Accounts.ui.config = function(options) {                                                                              // 24
	// validate options keys                                                                                             // 25
	var VALID_KEYS = ['onCreate', 'passwordSignupFields', 'extraSignupFields', 'forceEmailLowercase', 'forceUsernameLowercase','forcePasswordLowercase',
	'requestPermissions', 'requestOfflineToken', 'forceApprovalPrompt'];                                                 // 27
                                                                                                                      // 28
	_.each(_.keys(options), function(key) {                                                                              // 29
		if (!_.contains(VALID_KEYS, key)){                                                                                  // 30
			throw new Error("Accounts.ui.config: Invalid key: " + key);                                                        // 31
		}                                                                                                                   // 32
	});                                                                                                                  // 33
                                                                                                                      // 34
	if (options.onCreate && typeof options.onCreate === 'function') {                                                    // 35
		Accounts.ui._options.onCreate = options.onCreate;                                                                   // 36
	} else if (! options.onCreate ) {                                                                                    // 37
		//ignore and skip                                                                                                   // 38
	} else {                                                                                                             // 39
		throw new Error("Accounts.ui.config: Value for 'onCreate' must be a" +                                              // 40
				" function");                                                                                                     // 41
	}                                                                                                                    // 42
                                                                                                                      // 43
	options.extraSignupFields = options.extraSignupFields || [];                                                         // 44
                                                                                                                      // 45
	// deal with `passwordSignupFields`                                                                                  // 46
	if (options.passwordSignupFields) {                                                                                  // 47
		if (_.contains([                                                                                                    // 48
			"USERNAME_AND_EMAIL_CONFIRM",                                                                                      // 49
			"USERNAME_AND_EMAIL",                                                                                              // 50
			"USERNAME_AND_OPTIONAL_EMAIL",                                                                                     // 51
			"USERNAME_ONLY",                                                                                                   // 52
			"EMAIL_ONLY"                                                                                                       // 53
		], options.passwordSignupFields)) {                                                                                 // 54
			if (Accounts.ui._options.passwordSignupFields){                                                                    // 55
				throw new Error("Accounts.ui.config: Can't set `passwordSignupFields` more than once");                           // 56
			} else {                                                                                                           // 57
				Accounts.ui._options.passwordSignupFields = options.passwordSignupFields;                                         // 58
			}                                                                                                                  // 59
		} else {                                                                                                            // 60
			throw new Error("Accounts.ui.config: Invalid option for `passwordSignupFields`: " + options.passwordSignupFields);
		}                                                                                                                   // 62
	}                                                                                                                    // 63
                                                                                                                      // 64
	Accounts.ui._options.forceEmailLowercase = options.forceEmailLowercase;                                              // 65
	Accounts.ui._options.forceUsernameLowercase = options.forceUsernameLowercase;                                        // 66
	Accounts.ui._options.forcePasswordLowercase = options.forcePasswordLowercase;                                        // 67
                                                                                                                      // 68
	// deal with `requestPermissions`                                                                                    // 69
	if (options.requestPermissions) {                                                                                    // 70
		_.each(options.requestPermissions, function(scope, service) {                                                       // 71
			if (Accounts.ui._options.requestPermissions[service]) {                                                            // 72
				throw new Error("Accounts.ui.config: Can't set `requestPermissions` more than once for " + service);              // 73
			} else if (!(scope instanceof Array)) {                                                                            // 74
				throw new Error("Accounts.ui.config: Value for `requestPermissions` must be an array");                           // 75
			} else {                                                                                                           // 76
				Accounts.ui._options.requestPermissions[service] = scope;                                                         // 77
			}                                                                                                                  // 78
		});                                                                                                                 // 79
	}                                                                                                                    // 80
	if (typeof options.extraSignupFields !== 'object' || !options.extraSignupFields instanceof Array) {                  // 81
		throw new Error("Accounts.ui.config: `extraSignupFields` must be an array.");                                       // 82
	} else {                                                                                                             // 83
		if (options.extraSignupFields) {                                                                                    // 84
			_.each(options.extraSignupFields, function(field, index) {                                                         // 85
				if (!field.fieldName || !field.fieldLabel){                                                                       // 86
					throw new Error("Accounts.ui.config: `extraSignupFields` objects must have `fieldName` and `fieldLabel` attributes.");
				}                                                                                                                 // 88
				if (typeof field.visible === 'undefined'){                                                                        // 89
					field.visible = true;                                                                                            // 90
				}                                                                                                                 // 91
				Accounts.ui._options.extraSignupFields[index] = field;                                                            // 92
			});                                                                                                                // 93
		}                                                                                                                   // 94
	}                                                                                                                    // 95
                                                                                                                      // 96
	// deal with `requestOfflineToken`                                                                                   // 97
	if (options.requestOfflineToken) {                                                                                   // 98
		_.each(options.requestOfflineToken, function (value, service) {                                                     // 99
			if (service !== 'google'){                                                                                         // 100
				throw new Error("Accounts.ui.config: `requestOfflineToken` only supported for Google login at the moment.");      // 101
			}                                                                                                                  // 102
			if (Accounts.ui._options.requestOfflineToken[service]) {                                                           // 103
				throw new Error("Accounts.ui.config: Can't set `requestOfflineToken` more than once for " + service);             // 104
			} else {                                                                                                           // 105
				Accounts.ui._options.requestOfflineToken[service] = value;                                                        // 106
			}                                                                                                                  // 107
		});                                                                                                                 // 108
	}                                                                                                                    // 109
                                                                                                                      // 110
	// deal with `forceApprovalPrompt`                                                                                   // 111
	if (options.forceApprovalPrompt) {                                                                                   // 112
		_.each(options.forceApprovalPrompt, function (value, service) {                                                     // 113
			if (service !== 'google'){                                                                                         // 114
				throw new Error("Accounts.ui.config: `forceApprovalPrompt` only supported for Google login at the moment.");      // 115
			}                                                                                                                  // 116
			if (Accounts.ui._options.forceApprovalPrompt[service]) {                                                           // 117
				throw new Error("Accounts.ui.config: Can't set `forceApprovalPrompt` more than once for " + service);             // 118
			} else {                                                                                                           // 119
				Accounts.ui._options.forceApprovalPrompt[service] = value;                                                        // 120
			}                                                                                                                  // 121
		});                                                                                                                 // 122
	}                                                                                                                    // 123
};                                                                                                                    // 124
                                                                                                                      // 125
Accounts.ui._passwordSignupFields = function() {                                                                      // 126
	return Accounts.ui._options.passwordSignupFields || "EMAIL_ONLY";                                                    // 127
};                                                                                                                    // 128
                                                                                                                      // 129
                                                                                                                      // 130
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/en.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("en", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Reset your password",                                                                                       // 3
		newPassword: "New password",                                                                                        // 4
		newPasswordAgain: "New Password (again)",                                                                           // 5
		cancel: "Cancel",                                                                                                   // 6
		submit: "Set password"                                                                                              // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Choose a password",                                                                                         // 10
		newPassword: "New password",                                                                                        // 11
		newPasswordAgain: "New Password (again)",                                                                           // 12
		cancel: "Close",                                                                                                    // 13
		submit: "Set password"                                                                                              // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Email address verified",                                                                                 // 17
		dismiss: "Dismiss"                                                                                                  // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Dismiss",                                                                                                 // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Change password",                                                                                        // 24
		signOut: "Sign out"                                                                                                 // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Sign in",                                                                                                  // 28
		up: "Join"                                                                                                          // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "or"                                                                                                            // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Create",                                                                                                   // 35
		signIn: "Sign in",                                                                                                  // 36
		forgot: "Forgot password?",                                                                                         // 37
		createAcc: "Create account"                                                                                         // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Email",                                                                                                     // 41
		reset: "Reset password",                                                                                            // 42
		invalidEmail: "Invalid email"                                                                                       // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Cancel"                                                                                                      // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Change password",                                                                                          // 49
		cancel: "Cancel"                                                                                                    // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Sign in with",                                                                                         // 53
		configure: "Configure",                                                                                             // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Sign out"                                                                                                 // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "No login services configured"                                                                     // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Username or Email",                                                                               // 63
		username: "Username",                                                                                               // 64
		email: "Email",                                                                                                     // 65
		password: "Password"                                                                                                // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Username",                                                                                               // 69
		email: "Email",                                                                                                     // 70
		emailOpt: "Email (optional)",                                                                                       // 71
		password: "Password",                                                                                               // 72
		passwordAgain: "Password (again)"                                                                                   // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Current Password",                                                                                // 76
		newPassword: "New Password",                                                                                        // 77
		newPasswordAgain: "New Password (again)"                                                                            // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		emailSent: "Email sent",                                                                                            // 81
		passwordChanged: "Password changed"                                                                                 // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "There was an error",                                                                                 // 85
		userNotFound: "User not found",                                                                                     // 86
		invalidEmail: "Invalid email",                                                                                      // 87
		incorrectPassword: "Incorrect password",                                                                            // 88
		usernameTooShort: "Username must be at least 3 characters long",                                                    // 89
		passwordTooShort: "Password must be at least 6 characters long",                                                    // 90
		passwordsDontMatch: "Passwords don't match",                                                                        // 91
		newPasswordSameAsOld: "New and old passwords must be different",                                                    // 92
		signupsForbidden: "Signups forbidden"                                                                               // 93
	}                                                                                                                    // 94
});                                                                                                                   // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/es.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("es", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Restablece tu contraseña",                                                                                  // 3
		newPassword: "Nueva contraseña",                                                                                    // 4
		newPasswordAgain: "Nueva contraseña (otra vez)",                                                                    // 5
		cancel: "Cancelar",                                                                                                 // 6
		submit: "Guardar"                                                                                                   // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Escribe una contraseña",                                                                                    // 10
		newPassword: "Nueva contraseña",                                                                                    // 11
		newPasswordAgain: "Nueva contraseña (otra vez)",                                                                    // 12
		cancel: "Cerrar",                                                                                                   // 13
		submit: "Guardar contraseña"                                                                                        // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Correo electrónico verificado",                                                                          // 17
		dismiss: "Cerrar"                                                                                                   // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Cerrar",                                                                                                  // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Cambiar contraseña",                                                                                     // 24
		signOut: "Cerrar sesión"                                                                                            // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Iniciar sesión",                                                                                           // 28
		up: "registrarse"                                                                                                   // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "o"                                                                                                             // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Crear",                                                                                                    // 35
		signIn: "Iniciar sesión",                                                                                           // 36
		forgot: "¿Ha olvidado su contraseña?",                                                                              // 37
		createAcc: "Registrarse"                                                                                            // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Correo electrónico",                                                                                        // 41
		reset: "Restablecer contraseña",                                                                                    // 42
		invalidEmail: "Correo electrónico inválido"                                                                         // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Cancelar"                                                                                                    // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Cambiar contraseña",                                                                                       // 49
		cancel: "Cancelar"                                                                                                  // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Inicia sesión con",                                                                                    // 53
		configure: "Configurar",                                                                                            // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Cerrar sesión"                                                                                            // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "No hay ningún servicio configurado"                                                               // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Usuario o correo electrónico",                                                                    // 63
		username: "Usuario",                                                                                                // 64
		email: "Correo electrónico",                                                                                        // 65
		password: "Contraseña"                                                                                              // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Usuario",                                                                                                // 69
		email: "Correo electrónico",                                                                                        // 70
		emailOpt: "Correo elect. (opcional)",                                                                               // 71
		password: "Contraseña",                                                                                             // 72
		passwordAgain: "Contraseña (otra vez)"                                                                              // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Contraseña Actual",                                                                               // 76
		newPassword: "Nueva Contraseña",                                                                                    // 77
		newPasswordAgain: "Nueva Contraseña (otra vez)"                                                                     // 78
	},                                                                                                                   // 79
	infoMessages: {                                                                                                      // 80
		emailSent: "Email enviado",                                                                                         // 81
		passwordChanged: "Contraseña modificada"                                                                            // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Ha ocurrido un error",                                                                               // 85
		userNotFound: "El usuario no existe",                                                                               // 86
		invalidEmail: "Correo electrónico inválido",                                                                        // 87
		incorrectPassword: "Contraseña incorrecta",                                                                         // 88
		usernameTooShort: "El nombre de usuario tiene que tener 3 caracteres como mínimo",                                  // 89
		passwordTooShort: "La contraseña tiene que tener 6 caracteres como mínimo",                                         // 90
		passwordsDontMatch: "Las contraseñas no son iguales",                                                               // 91
		newPasswordSameAsOld: "La contraseña nueva y la actual no pueden ser iguales"                                       // 92
	}                                                                                                                    // 93
});                                                                                                                   // 94
                                                                                                                      // 95
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/ca.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("ca", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Restablir la contrasenya",                                                                                  // 3
		newPassword: "Nova contrasenya",                                                                                    // 4
		newPasswordAgain: "Nova contrasenya (un altre cop)",                                                                // 5
		cancel: "Cancel·lar",                                                                                               // 6
		submit: "Guardar"                                                                                                   // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Escriu una contrasenya",                                                                                    // 10
		newPassword: "Nova contrasenya",                                                                                    // 11
		newPasswordAgain: "Nova contrasenya (un altre cop)",                                                                // 12
		cancel: "Tancar",                                                                                                   // 13
		submit: "Guardar contrasenya"                                                                                       // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Adreça electrònica verificada",                                                                          // 17
		dismiss: "Tancar"                                                                                                   // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Tancar",                                                                                                  // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Canviar contrasenya",                                                                                    // 24
		signOut: "Tancar sessió"                                                                                            // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Iniciar sessió",                                                                                           // 28
		up: "Registrar-se"                                                                                                  // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "o bé"                                                                                                          // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Crear",                                                                                                    // 35
		signIn: "Iniciar sessió",                                                                                           // 36
		forgot: "Ha oblidat la contrasenya?",                                                                               // 37
		createAcc: "Registrar-se"                                                                                           // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Adreça electrònica",                                                                                        // 41
		reset: "Restablir contrasenya",                                                                                     // 42
		invalidEmail: "Adreça invàlida"                                                                                     // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Cancel·lar"                                                                                                  // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Canviar contrasenya",                                                                                      // 49
		cancel: "Cancel·lar"                                                                                                // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Inicia sessió amb",                                                                                    // 53
		configure: "Configurar"                                                                                             // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Tancar sessió"                                                                                            // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "No hi ha cap servei configurat"                                                                   // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Usuari o correu electrònic",                                                                      // 63
		username: "Usuari",                                                                                                 // 64
		email: "Adreça electrònica",                                                                                        // 65
		password: "Contrasenya"                                                                                             // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Usuari",                                                                                                 // 69
		email: "Adreça electrònica",                                                                                        // 70
		emailOpt: "Adreça elect. (opcional)",                                                                               // 71
		password: "Contrasenya",                                                                                            // 72
		passwordAgain: "Contrasenya (un altre cop)"                                                                         // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Contrasenya Actual",                                                                              // 76
		newPassword: "Nova Contrasenya",                                                                                    // 77
		newPasswordAgain: "Nova Contrasenya (un altre cop)"                                                                 // 78
	},                                                                                                                   // 79
	infoMessages: {                                                                                                      // 80
		emailSent: "Email enviat",                                                                                          // 81
		passwordChanged: "Contrasenya canviada"                                                                             // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Hi ha hagut un error",                                                                               // 85
		userNotFound: "L'usuari no existeix",                                                                               // 86
		invalidEmail: "Adreça invàlida",                                                                                    // 87
		incorrectPassword: "Contrasenya incorrecta",                                                                        // 88
		usernameTooShort: "El nom d'usuari ha de tenir 3 caracters com a mínim",                                            // 89
		passwordTooShort: "La contrasenya ha de tenir 6 caracters como a mínim",                                            // 90
		passwordsDontMatch: "Les contrasenyes no són iguals",                                                               // 91
		newPasswordSameAsOld: "La contrasenya nova i l'actual no poden ser iguals"                                          // 92
	}                                                                                                                    // 93
});                                                                                                                   // 94
                                                                                                                      // 95
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/fr.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("fr", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Réinitialiser mon mot de passe",                                                                            // 3
		newPassword: "Nouveau mot de passe",                                                                                // 4
		newPasswordAgain: "Nouveau mot de passe (confirmation)",                                                            // 5
		cancel: "Annuler",                                                                                                  // 6
		submit: "Définir le mot de passe"                                                                                   // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Choisir un mot de passe",                                                                                   // 10
		newPassword: "Nouveau mot de passe",                                                                                // 11
		newPasswordAgain: "Nouveau mot de passe (confirmation)",                                                            // 12
		cancel: "Fermer",                                                                                                   // 13
		submit: "Définir le mot de passe"                                                                                   // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "L'adresse email a été vérifiée",                                                                         // 17
		dismiss: "Fermer"                                                                                                   // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Fermer",                                                                                                  // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Changer le mot de passe",                                                                                // 24
		signOut: "Déconnexion"                                                                                              // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Connexion",                                                                                                // 28
		up: "Inscription"                                                                                                   // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "ou"                                                                                                            // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Créer",                                                                                                    // 35
		signIn: "Connexion",                                                                                                // 36
		forgot: "Mot de passe oublié ?",                                                                                    // 37
		createAcc: "Inscription"                                                                                            // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Email",                                                                                                     // 41
		reset: "Réinitialiser le mot de passe",                                                                             // 42
		invalidEmail: "L'adresse email est invalide"                                                                        // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Annuler"                                                                                                     // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Changer le mot de passe",                                                                                  // 49
		cancel: "Annuler"                                                                                                   // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Se connecter avec",                                                                                    // 53
		configure: "Configurer",                                                                                            // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Déconnexion"                                                                                              // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Aucun service d'authentification n'est configuré"                                                 // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Nom d'utilisateur ou email",                                                                      // 63
		username: "Nom d'utilisateur",                                                                                      // 64
		email: "Email",                                                                                                     // 65
		password: "Mot de passe"                                                                                            // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Nom d'utilisateur",                                                                                      // 69
		email: "Email",                                                                                                     // 70
		emailOpt: "Email (optionnel)",                                                                                      // 71
		password: "Mot de passe",                                                                                           // 72
		passwordAgain: "Mot de passe (confirmation)"                                                                        // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Mot de passe actuel",                                                                             // 76
		newPassword: "Nouveau mot de passe",                                                                                // 77
		newPasswordAgain: "Nouveau mot de passe (confirmation)"                                                             // 78
	},                                                                                                                   // 79
	infoMessages: {                                                                                                      // 80
		emailSent: "Email envoyé",                                                                                          // 81
		passwordChanged: "Mot de passe modifié"                                                                             // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Il y avait une erreur",                                                                              // 85
		userNotFound: "Utilisateur non trouvé",                                                                             // 86
		invalidEmail: "L'adresse email est invalide",                                                                       // 87
		incorrectPassword: "Mot de passe incorrect",                                                                        // 88
		usernameTooShort: "Le nom d'utilisateur doit comporter au moins 3 caractères",                                      // 89
		passwordTooShort: "Le mot de passe doit comporter au moins 6 caractères",                                           // 90
		passwordsDontMatch: "Les mots de passe ne sont pas identiques",                                                     // 91
		newPasswordSameAsOld: "Le nouveau et le vieux mot de passe doivent être différent"                                  // 92
	}                                                                                                                    // 93
});                                                                                                                   // 94
                                                                                                                      // 95
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/de.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("de", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Passwort zurücksetzen",                                                                                     // 3
		newPassword: "Neues Passwort",                                                                                      // 4
		newPasswordAgain: "Neues Passwort (wiederholen)",                                                                   // 5
		cancel: "Abbrechen",                                                                                                // 6
		submit: "Passwort ändern"                                                                                           // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Passwort wählen",                                                                                           // 10
		newPassword: "Neues Passwort",                                                                                      // 11
		newPasswordAgain: "Neues Passwort (wiederholen)",                                                                   // 12
		cancel: "Schließen",                                                                                                // 13
		submit: "Passwort ändern"                                                                                           // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Email Adresse verifiziert",                                                                              // 17
		dismiss: "Schließen"                                                                                                // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Schließen"                                                                                                // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Passwort ändern",                                                                                        // 24
		signOut: "Abmelden"                                                                                                 // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Anmelden",                                                                                                 // 28
		up: "Registrieren"                                                                                                  // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "oder"                                                                                                          // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Erstellen",                                                                                                // 35
		signIn: "Anmelden",                                                                                                 // 36
		forgot: "Passwort vergessen?",                                                                                      // 37
		createAcc: "Account erstellen"                                                                                      // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Email",                                                                                                     // 41
		reset: "Passwort zurücksetzen",                                                                                     // 42
		invalidEmail: "Ungültige Email Adresse"                                                                             // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Abbrechen"                                                                                                   // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Passwort ändern",                                                                                          // 49
		cancel: "Abbrechen"                                                                                                 // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Anmelden mit",                                                                                         // 53
		configure: "Konfigurieren",                                                                                         // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Abmelden"                                                                                                 // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Keine Anmelde Services konfiguriert"                                                              // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Benutzername oder Email",                                                                         // 63
		username: "Benutzername",                                                                                           // 64
		email: "Email",                                                                                                     // 65
		password: "Passwort"                                                                                                // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Benutzername",                                                                                           // 69
		email: "Email",                                                                                                     // 70
		emailOpt: "Email (freiwillig)",                                                                                     // 71
		password: "Passwort",                                                                                               // 72
		passwordAgain: "Passwort (wiederholen)"                                                                             // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Aktuelles Passwort",                                                                              // 76
		newPassword: "Neues Passwort",                                                                                      // 77
		newPasswordAgain: "Neues Passwort (wiederholen)"                                                                    // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		sent: "Email gesendet",                                                                                             // 81
		passwordChanged: "Passwort geändert"                                                                                // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Es gab einen Fehler",                                                                                // 85
		userNotFound: "Benutzer nicht gefunden",                                                                            // 86
		invalidEmail: "Ungültige Email Adresse",                                                                            // 87
		incorrectPassword: "Falsches Passwort",                                                                             // 88
		usernameTooShort: "Der Benutzername muss mindestens 3 Buchstaben lang sein",                                        // 89
		passwordTooShort: "Passwort muss mindestens 6 Zeichen lang sein",                                                   // 90
		passwordsDontMatch: "Die Passwörter stimmen nicht überein",                                                         // 91
		newPasswordSameAsOld: "Neue und aktuelle Passwörter müssen unterschiedlich sein"                                    // 92
	}                                                                                                                    // 93
});                                                                                                                   // 94
                                                                                                                      // 95
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/it.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("it", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Reimposta la password",                                                                                     // 3
		newPassword: "Nuova password",                                                                                      // 4
		newPasswordAgain: "Nuova password (di nuovo)",                                                                      // 5
		cancel: "Annulla",                                                                                                  // 6
		submit: "Imposta password"                                                                                          // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Scegli una password",                                                                                       // 10
		newPassword: "Nuova password",                                                                                      // 11
		newPasswordAgain: "Nuova password (di nuovo)",                                                                      // 12
		cancel: "Chiudi",                                                                                                   // 13
		submit: "Imposta password"                                                                                          // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Indirizzo email verificato",                                                                             // 17
		dismiss: "Chiudi"                                                                                                   // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Chiudi",                                                                                                  // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Cambia password",                                                                                        // 24
		signOut: "Esci"                                                                                                     // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Accedi",                                                                                                   // 28
		up: "Registrati"                                                                                                    // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "oppure"                                                                                                        // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Crea",                                                                                                     // 35
		signIn: "Accedi",                                                                                                   // 36
		forgot: "Password dimenticata?",                                                                                    // 37
		createAcc: "Crea un account"                                                                                        // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Email",                                                                                                     // 41
		reset: "Reimposta la password",                                                                                     // 42
		invalidEmail: "Email non valida"                                                                                    // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Cancella"                                                                                                    // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Cambia password",                                                                                          // 49
		cancel: "Annulla"                                                                                                   // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Accedi con",                                                                                           // 53
		configure: "Configura",                                                                                             // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Esci"                                                                                                     // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Nessun servizio di accesso configurato"                                                           // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Username o Email",                                                                                // 63
		username: "Username",                                                                                               // 64
		email: "Email",                                                                                                     // 65
		password: "Password"                                                                                                // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Username",                                                                                               // 69
		email: "Email",                                                                                                     // 70
		emailOpt: "Email (opzionale)",                                                                                      // 71
		password: "Password",                                                                                               // 72
		passwordAgain: "Password (di nuovo)"                                                                                // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Password corrente",                                                                               // 76
		newPassword: "Nuova password",                                                                                      // 77
		newPasswordAgain: "Nuova password (di nuovo)"                                                                       // 78
	},                                                                                                                   // 79
	infoMessages: {                                                                                                      // 80
		emailSent: "Email inviata",                                                                                         // 81
		passwordChanged: "Password changed"                                                                                 // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "C'era un errore",                                                                                    // 85
		userNotFound: "Username non trovato",                                                                               // 86
		invalidEmail: "Email non valida",                                                                                   // 87
		incorrectPassword: "Password errata",                                                                               // 88
		usernameTooShort: "La Username deve essere almeno di 3 caratteri",                                                  // 89
		passwordTooShort: "La Password deve essere almeno di 6 caratteri",                                                  // 90
		passwordsDontMatch: "Le password non corrispondono",                                                                // 91
		newPasswordSameAsOld: "Nuove e vecchie password devono essere diversi"                                              // 92
	}                                                                                                                    // 93
});                                                                                                                   // 94
                                                                                                                      // 95
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/pt-PT.i18n.js                                                            //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
ptPT = {                                                                                                              // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Esqueci-me da palavra-passe",                                                                               // 3
		newPassword: "Nova palavra-passe",                                                                                  // 4
		newPasswordAgain: "Nova palavra-passe (confirmacao)",                                                               // 5
		cancel: "Cancelar",                                                                                                 // 6
		submit: "Alterar palavra-passe"                                                                                     // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Introduza a nova palavra-passe",                                                                            // 10
		newPassword: "Nova palavra-passe",                                                                                  // 11
		newPasswordAgain: "Nova palavra-passe (confirmacao)",                                                               // 12
		cancel: "Fechar",                                                                                                   // 13
		submit: "Alterar palavra-passe"                                                                                     // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "E-mail verificado!",                                                                                     // 17
		dismiss: "Ignorar"                                                                                                  // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Ignorar"                                                                                                  // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Mudar palavra-passe",                                                                                    // 24
		signOut: "Sair"                                                                                                     // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Entrar",                                                                                                   // 28
		up: "Registar"                                                                                                      // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "ou"                                                                                                            // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Criar",                                                                                                    // 35
		signIn: "Entrar",                                                                                                   // 36
		forgot: "Esqueci-me da palavra-passe",                                                                              // 37
		createAcc: "Registar"                                                                                               // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "E-mail",                                                                                                    // 41
		reset: "Alterar palavra-passe",                                                                                     // 42
		sent: "E-mail enviado",                                                                                             // 43
		invalidEmail: "E-mail inválido"                                                                                     // 44
	},                                                                                                                   // 45
	loginButtonsBackToLoginLink: {                                                                                       // 46
		back: "Cancelar"                                                                                                    // 47
	},                                                                                                                   // 48
	loginButtonsChangePassword: {                                                                                        // 49
		submit: "Mudar palavra-passe",                                                                                      // 50
		cancel: "Cancelar"                                                                                                  // 51
	},                                                                                                                   // 52
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 53
		signInWith: "Entrar com",                                                                                           // 54
		configure: "Configurar"                                                                                             // 55
	},                                                                                                                   // 56
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 57
		signOut: "Sair"                                                                                                     // 58
	},                                                                                                                   // 59
	loginButtonsLoggedOut: {                                                                                             // 60
		noLoginServices: "Nenhum servico de login configurado"                                                              // 61
	},                                                                                                                   // 62
	loginFields: {                                                                                                       // 63
		usernameOrEmail: "Utilizador ou E-mail",                                                                            // 64
		username: "Utilizador",                                                                                             // 65
		email: "E-mail",                                                                                                    // 66
		password: "Palavra-passe"                                                                                           // 67
	},                                                                                                                   // 68
	signupFields: {                                                                                                      // 69
		username: "Utilizador",                                                                                             // 70
		email: "E-mail",                                                                                                    // 71
		emailOpt: "E-mail (opcional)",                                                                                      // 72
		password: "Palavra-passe",                                                                                          // 73
		passwordAgain: "Palavra-passe (confirmacão)"                                                                        // 74
	},                                                                                                                   // 75
	changePasswordFields: {                                                                                              // 76
		currentPassword: "Palavra-passe atual",                                                                             // 77
		newPassword: "Nova palavra-passe",                                                                                  // 78
		newPasswordAgain: "Nova palavra-passe (confirmacao)"                                                                // 79
	},                                                                                                                   // 80
	infoMessages: {                                                                                                      // 81
		emailSent: "E-mail enviado",                                                                                        // 82
		passwordChanged: "Palavra-passe alterada"                                                                           // 83
	},                                                                                                                   // 84
	errorMessages: {                                                                                                     // 85
		genericTitle: "Houve um erro",                                                                                      // 86
		usernameTooShort: "Utilizador precisa de ter mais de 3 caracteres",                                                 // 87
		invalidEmail: "E-mail inválido",                                                                                    // 88
		passwordTooShort: "Palavra-passe precisa ter mais de 6 caracteres",                                                 // 89
		passwordsDontMatch: "As Palavras-passe estão diferentes",                                                           // 90
		userNotFound: "Utilizador não encontrado",                                                                          // 91
		incorrectPassword: "Palavra-passe incorreta",                                                                       // 92
		newPasswordSameAsOld: "A nova palavra-passe tem de ser diferente da antiga"                                         // 93
	}                                                                                                                    // 94
};                                                                                                                    // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/pt-BR.i18n.js                                                            //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
ptBR = {                                                                                                              // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Esqueceu sua senha?",                                                                                       // 3
		newPassword: "Nova senha",                                                                                          // 4
		newPasswordAgain: "Nova senha (confirmacao)",                                                                       // 5
		cancel: "Cancelar",                                                                                                 // 6
		submit: "Alterar senha"                                                                                             // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Digite a nova senha",                                                                                       // 10
		newPassword: "Nova senha",                                                                                          // 11
		newPasswordAgain: "Nova senha (confirmacao)",                                                                       // 12
		cancel: "Fechar",                                                                                                   // 13
		submit: "Alterar senha"                                                                                             // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "E-mail verificado!",                                                                                     // 17
		dismiss: "Ignorar"                                                                                                  // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Ignorar"                                                                                                  // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Mudar senha",                                                                                            // 24
		signOut: "Sair"                                                                                                     // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Entrar",                                                                                                   // 28
		up: "Cadastrar"                                                                                                     // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "ou"                                                                                                            // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Criar",                                                                                                    // 35
		signIn: "Login",                                                                                                    // 36
		forgot: "Esqueceu sua senha?",                                                                                      // 37
		createAcc: "Cadastrar"                                                                                              // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "E-mail",                                                                                                    // 41
		reset: "Alterar senha",                                                                                             // 42
		invalidEmail: "E-mail inválido"                                                                                     // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Cancelar"                                                                                                    // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Mudar senha",                                                                                              // 49
		cancel: "Cancelar"                                                                                                  // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Logar com",                                                                                            // 53
		configure: "Configurar",                                                                                            // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Sair"                                                                                                     // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Nenhum servico de login configurado"                                                              // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Usuário ou E-mail",                                                                               // 63
		username: "Usuário",                                                                                                // 64
		email: "E-mail",                                                                                                    // 65
		password: "Senha"                                                                                                   // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Usuário",                                                                                                // 69
		email: "E-mail",                                                                                                    // 70
		emailOpt: "E-mail (opcional)",                                                                                      // 71
		password: "Senha",                                                                                                  // 72
		passwordAgain: "Senha (confirmacão)"                                                                                // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Senha atual",                                                                                     // 76
		newPassword: "Nova Senha",                                                                                          // 77
		newPasswordAgain: "Nova Senha (confirmacao)"                                                                        // 78
	},                                                                                                                   // 79
	infoMessages: {                                                                                                      // 80
		emailSent: "E-mail enviado",                                                                                        // 81
		passwordChanged: "Senha alterada"                                                                                   // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Houve um erro",                                                                                      // 85
		userNotFound: "Usuário não encontrado",                                                                             // 86
		invalidEmail: "E-mail inválido",                                                                                    // 87
		incorrectPassword: "Senha incorreta",                                                                               // 88
		usernameTooShort: "Usuário precisa ter mais de 3 caracteres",                                                       // 89
		passwordTooShort: "Senha precisa ter mais de 6 caracteres",                                                         // 90
		passwordsDontMatch: "Senhas estão diferentes",                                                                      // 91
		newPasswordSameAsOld: "A nova senha tem de ser diferente da antiga"                                                 // 92
	}                                                                                                                    // 93
};                                                                                                                    // 94
                                                                                                                      // 95
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/pt.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map('pt', ptPT);                                                                                                 // 1
i18n.map('pt-PT', ptPT);                                                                                              // 2
i18n.map('pt-BR', ptBR);                                                                                              // 3
                                                                                                                      // 4
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/ru.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("ru", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Сбросить пароль",                                                                                           // 3
		newPassword: "Новый пароль",                                                                                        // 4
		newPasswordAgain: "Новый пароль (еще раз)",                                                                         // 5
		cancel: "Отмена",                                                                                                   // 6
		submit: "Сохранить пароль"                                                                                          // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Выбрать пароль",                                                                                            // 10
		newPassword: "Новый пароль",                                                                                        // 11
		newPasswordAgain: "Новый пароль (еще раз)",                                                                         // 12
		cancel: "Отмена",                                                                                                   // 13
		submit: "Сохранить пароль"                                                                                          // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Email подтвержден",                                                                                      // 17
		dismiss: "Закрыть"                                                                                                  // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
			dismiss: "Закрыть"                                                                                                 // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Изменить пароль",                                                                                        // 24
		signOut: "Выйти"                                                                                                    // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Войти",                                                                                                    // 28
		up: "Зарегистрироваться"                                                                                            // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "или"                                                                                                           // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Создать",                                                                                                  // 35
		signIn: "Войти",                                                                                                    // 36
		forgot: "Забыли пароль?",                                                                                           // 37
		createAcc: "Создать аккаунт"                                                                                        // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Email",                                                                                                     // 41
		reset: "Сбросить пароль",                                                                                           // 42
		invalidEmail: "Некорректный email"                                                                                  // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Отмена"                                                                                                      // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Изменить пароль",                                                                                          // 49
		cancel: "Отмена"                                                                                                    // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Войти через",                                                                                          // 53
		configure: "Настроить вход через",                                                                                  // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Выйти"                                                                                                    // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Сервис для входа не настроен"                                                                     // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Имя пользователя или email",                                                                      // 63
		username: "Имя пользователя",                                                                                       // 64
		email: "Email",                                                                                                     // 65
		password: "Пароль"                                                                                                  // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Имя пользователя",                                                                                       // 69
		email: "Email",                                                                                                     // 70
		emailOpt: "Email (необязательный)",                                                                                 // 71
		password: "Пароль",                                                                                                 // 72
		passwordAgain: "Пароль (еще раз)"                                                                                   // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Текущий пароль",                                                                                  // 76
		newPassword: "Новый пароль",                                                                                        // 77
		newPasswordAgain: "Новый пароль (еще раз)"                                                                          // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		sent: "Вам отправлено письмо",                                                                                      // 81
		passwordChanged: "Пароль изменён"                                                                                   // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Там была ошибка",                                                                                    // 85
		userNotFound: "Пользователь не найден",                                                                             // 86
		invalidEmail: "Некорректный email",                                                                                 // 87
		incorrectPassword: "Неправильный пароль",                                                                           // 88
		usernameTooShort: "Имя пользователя должно быть длиной не менее 3-х символов",                                      // 89
		passwordTooShort: "Пароль должен быть длиной не менее 6-ти символов",                                               // 90
		passwordsDontMatch: "Пароли не совпадают",                                                                          // 91
		newPasswordSameAsOld: "Новый и старый пароли должны быть разными"                                                   // 92
	}                                                                                                                    // 93
});                                                                                                                   // 94
                                                                                                                      // 95
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/el.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("el", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Ακύρωση κωδικού",                                                                                           // 3
		newPassword: "Νέος κωδικός",                                                                                        // 4
		newPasswordAgain: "Νέος Κωδικός (ξανά)",                                                                            // 5
		cancel: "Ακύρωση",                                                                                                  // 6
		submit: "Ορισμός κωδικού"                                                                                           // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Επιλέξτε κωδικό",                                                                                           // 10
		newPassword: "Νέος κωδικός",                                                                                        // 11
		newPasswordAgain: "Νέος Κωδικός (ξανά)",                                                                            // 12
		cancel: "Ακύρωση",                                                                                                  // 13
		submit: "Ορισμός κωδικού"                                                                                           // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Ο λογαριασμός ηλεκτρονικού ταχυδρομείου έχει επιβεβαιωθεί",                                              // 17
		dismiss: "Κλείσιμο"                                                                                                 // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Κλείσιμο",                                                                                                // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Αλλαγή κωδικού",                                                                                         // 24
		signOut: "Αποσύνδεση"                                                                                               // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Είσοδος",                                                                                                  // 28
		up: "Εγγραφή"                                                                                                       // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "ή"                                                                                                             // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Δημιουργία",                                                                                               // 35
		signIn: "Είσοδος",                                                                                                  // 36
		forgot: "Ξεχάσατε τον κωδικό σας;",                                                                                 // 37
		createAcc: "Δημιουργία λογαριασμού"                                                                                 // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Ηλεκτρονικό ταχυδρομείο (email)",                                                                           // 41
		reset: "Ακύρωση κωδικού",                                                                                           // 42
		invalidEmail: "Μη έγκυρος λογαριασμός ηλεκτρονικού ταχυδρομείου (email)"                                            // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Επιστροφή"                                                                                                   // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Αλλαγή κωδικού",                                                                                           // 49
		cancel: "Ακύρωση"                                                                                                   // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Είσοδος με",                                                                                           // 53
		configure: "Διαμόρφωση",                                                                                            // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Αποσύνδεση"                                                                                               // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Δεν έχουν διαμορφωθεί υπηρεσίες εισόδου"                                                          // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Όνομα χρήστη ή Λογαριασμός Ηλεκτρονικού Ταχυδρομείου",                                            // 63
		username: "Όνομα χρήστη",                                                                                           // 64
		email: "Ηλεκτρονικό ταχυδρομείο (email)",                                                                           // 65
		password: "Κωδικός"                                                                                                 // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Όνομα χρήστη",                                                                                           // 69
		email: "Ηλεκτρονικό ταχυδρομείο (email)",                                                                           // 70
		emailOpt: "Ηλεκτρονικό ταχυδρομείο (προαιρετικό)",                                                                  // 71
		password: "Κωδικός",                                                                                                // 72
		passwordAgain: "Κωδικός (ξανά)"                                                                                     // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Ισχύων Κωδικός",                                                                                  // 76
		newPassword: "Νέος Κωδικός",                                                                                        // 77
		newPasswordAgain: "Νέος Κωδικός (ξανά)"                                                                             // 78
	},                                                                                                                   // 79
	infoMessages: {                                                                                                      // 80
		emailSent: "Το email έχει αποσταλεί",                                                                               // 81
		passwordChanged: "Password changed"                                                                                 // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Υπήρξε ένα σφάλμα",                                                                                  // 85
		userNotFound: "Ο χρήστης δεν βρέθηκε",                                                                              // 86
		invalidEmail: "Μη έγκυρος λογαριασμός ηλεκτρονικού ταχυδρομείου (email)",                                           // 87
		incorrectPassword: "Λάθος κωδικός",                                                                                 // 88
		usernameTooShort: "Το όνομα χρήστη πρέπει να είναι τουλάχιστον 3 χαρακτήρες",                                       // 89
		passwordTooShort: "Ο κωδικός πρέπει να είναι τουλάχιστον 6 χαρακτήρες",                                             // 90
		passwordsDontMatch: "Οι κωδικοί δεν ταιριάζουν",                                                                    // 91
		newPasswordSameAsOld: "Νέα και παλιά κωδικούς πρόσβασης πρέπει να είναι διαφορετική"                                // 92
	}                                                                                                                    // 93
});                                                                                                                   // 94
                                                                                                                      // 95
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/ko.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("ko", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "비밀번호 초기화하기",                                                                                                // 3
		newPassword: "새로운 비밀번호",                                                                                            // 4
		newPasswordAgain: "새로운 비밀번호 (확인)",                                                                                  // 5
		cancel: "취소",                                                                                                       // 6
		submit: "변경"                                                                                                        // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "비밀번호를 입력해주세요",                                                                                              // 10
		newPassword: "새로운 비밀번호",                                                                                            // 11
		newPasswordAgain: "새로운 비밀번호 (확인)",                                                                                  // 12
		cancel: "닫기",                                                                                                       // 13
		submit: "변경"                                                                                                        // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "이메일 주소가 인증되었습니다",                                                                                        // 17
		dismiss: "취소"                                                                                                       // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "취소",                                                                                                      // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "비밀번호 변경하기",                                                                                              // 24
		signOut: "로그아웃"                                                                                                     // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "로그인",                                                                                                      // 28
		up: "계정 만들기"                                                                                                        // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "또는"                                                                                                            // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "만들기",                                                                                                      // 35
		signIn: "로그인",                                                                                                      // 36
		forgot: "비밀번호를 잊어버리셨나요?",                                                                                           // 37
		createAcc: "계정 만들기"                                                                                                 // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "이메일 주소",                                                                                                    // 41
		reset: "비밀번호 초기화하기",                                                                                                // 42
		invalidEmail: "올바르지 않은 이메일 주소입니다"                                                                                   // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "취소"                                                                                                          // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "비밀번호 변경하기",                                                                                                // 49
		cancel: "취소"                                                                                                        // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "다음으로 로그인하기:",                                                                                          // 53
		configure: "설정",                                                                                                    // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "로그아웃"                                                                                                     // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "사용 가능한 로그인 서비스가 없습니다"                                                                             // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "사용자 이름 또는 이메일 주소",                                                                                // 63
		username: "사용자 이름",                                                                                                 // 64
		email: "이메일 주소",                                                                                                    // 65
		password: "비밀번호"                                                                                                    // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "사용자 이름",                                                                                                 // 69
		email: "이메일 주소",                                                                                                    // 70
		emailOpt: "이메일 주소 (선택)",                                                                                            // 71
		password: "비밀번호",                                                                                                   // 72
		passwordAgain: "비밀번호 (확인)"                                                                                          // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "현재 비밀번호",                                                                                         // 76
		newPassword: "새로운 비밀번호",                                                                                            // 77
		newPasswordAgain: "새로운 비밀번호 (확인)"                                                                                   // 78
	},                                                                                                                   // 79
	infoMessages: {                                                                                                      // 80
		emailSent: "이메일이 발송되었습니다",                                                                                          // 81
		passwordChanged: "비밀번호가 변경되었습니다"                                                                                    // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "오류가 발생했습니다",                                                                                         // 85
		userNotFound: "찾을 수 없는 회원입니다",                                                                                      // 86
		invalidEmail: "잘못된 이메일 주소",                                                                                         // 87
		incorrectPassword: "비밀번호가 틀렸습니다",                                                                                   // 88
		usernameTooShort: "사용자 이름은 최소 3글자 이상이어야 합니다",                                                                       // 89
		passwordTooShort: "비밀번호는 최소 6글자 이상이어야 합니다",                                                                         // 90
		passwordsDontMatch: "비밀번호가 같지 않습니다",                                                                                // 91
		newPasswordSameAsOld: "새 비밀번호와 기존 비밀번호는 달라야합니다"                                                                     // 92
	}                                                                                                                    // 93
});                                                                                                                   // 94
                                                                                                                      // 95
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/ar.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("ar", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "استرجع كلمة المرور",                                                                                        // 3
		newPassword: "كلمة المرور الجديدة",                                                                                 // 4
		newPasswordAgain: "أعد كتابة كلمة السر الجديدة",                                                                    // 5
		cancel: "إلغ",                                                                                                      // 6
		submit: "تم"                                                                                                        // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "اختر كلمة سر",                                                                                              // 10
		newPassword: "كلمة السر",                                                                                           // 11
		newPasswordAgain: "أعد كتابة كلمة السر الجديدة",                                                                    // 12
		cancel: "أغلق",                                                                                                     // 13
		submit: "تم"                                                                                                        // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "تم تأكيد البريد",                                                                                        // 17
		dismiss:  "حسنًا"                                                                                                   // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "حسنًا"                                                                                                    // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "غير كلمة السر",                                                                                          // 24
		signOut: "تسجيل الخروج"                                                                                             // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "دخول",                                                                                                     // 28
		up: "إنشاء حساب"                                                                                                    // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "أو"                                                                                                            // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "أنشئ",                                                                                                     // 35
		signIn: "دخول",                                                                                                     // 36
		forgot: "نسيت كلمة السر؟",                                                                                          // 37
		createAcc: "أنشئ حسابا"                                                                                             // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "البريد",                                                                                                    // 41
		reset: "إعادة تعين كلمة السر",                                                                                      // 42
		invalidEmail: "البريد خاطئ"                                                                                         // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "إلغ"                                                                                                         // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "غير كلمة السر",                                                                                            // 49
		cancel: "إلغ"                                                                                                       // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "سجل الدخول عبر",                                                                                       // 53
		configure: "تعيين"                                                                                                  // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "اخرج"                                                                                                     // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "لا يوجد خدمة دخول مفعله"                                                                          // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "اسم المستخدم او عنوان البريد",                                                                    // 63
		username: "اسم المستخدم",                                                                                           // 64
		email: "البريد",                                                                                                    // 65
		password: "كلمة السر"                                                                                               // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "اسم المستخدم",                                                                                           // 69
		email: "البريد",                                                                                                    // 70
		emailOpt: "-اختياري- البريد",                                                                                       // 71
		password: "كلمة السر",                                                                                              // 72
		passwordAgain: "أعد كتابة كلمة السر"                                                                                // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "كلمة السر الحالية",                                                                               // 76
		newPassword: "كلمة السر الجديدة",                                                                                   // 77
		newPasswordAgain: "أعد كتابة كلمة السر الجديدة"                                                                     // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		emailSent: "تم الارسال",                                                                                            // 81
		passwordChanged: "تمت إعادة تعيين كلمة السر"                                                                        // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "كان هناك خطأ",                                                                                       // 85
		userNotFound: "المستخدم غير موجود",                                                                                 // 86
		invalidEmail: "بريد خاطئ",                                                                                          // 87
		incorrectPassword: "كلمة السر خطأ",                                                                                 // 88
		usernameTooShort: "اسم المستخدم لابد ان يكون علي الاقل ٣ حروف",                                                     // 89
		passwordTooShort: "كلمة السر لابد ان تكون علي الاقل ٦ احرف",                                                        // 90
		passwordsDontMatch: "كلمة السر غير متطابقة",                                                                        // 91
		newPasswordSameAsOld: "لابد من اختيار كلمة سر مختلفة عن السابقة",                                                   // 92
		signupsForbidden: "التسجيل مغلق"                                                                                    // 93
	}                                                                                                                    // 94
});                                                                                                                   // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/pl.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("pl", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Resetuj hasło",                                                                                             // 3
		newPassword: "Nowe hasło",                                                                                          // 4
		newPasswordAgain: "Nowe hasło (powtórz)",                                                                           // 5
		cancel: "Anuluj",                                                                                                   // 6
		submit: "Ustaw hasło"                                                                                               // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Wybierz hasło",                                                                                             // 10
		newPassword: "Nowe hasło",                                                                                          // 11
		newPasswordAgain: "Nowe hasło (powtórz)",                                                                           // 12
		cancel: "Zamknij",                                                                                                  // 13
		submit: "Ustaw hasło"                                                                                               // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Adres email został zweryfikowany",                                                                       // 17
		dismiss: "Odrzuć"                                                                                                   // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Odrzuć"                                                                                                   // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Zmień hasło",                                                                                            // 24
		signOut: "Wyloguj się"                                                                                              // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Zaloguj się",                                                                                              // 28
		up: "Zarejestruj się"                                                                                               // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "lub"                                                                                                           // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Stwórz",                                                                                                   // 35
		signIn: "Zaloguj się ",                                                                                             // 36
		forgot: "Nie pamiętasz hasła?",                                                                                     // 37
		createAcc: "Utwórz konto"                                                                                           // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Email",                                                                                                     // 41
		reset: "Resetuj hasło",                                                                                             // 42
		invalidEmail: "Niepoprawny email"                                                                                   // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Anuluj"                                                                                                      // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Zmień hasło",                                                                                              // 49
		cancel: "Anuluj"                                                                                                    // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Zaloguj się przez",                                                                                    // 53
		configure: "Configure"                                                                                              // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Wyloguj się"                                                                                              // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Nie skonfigurowano usługi logowania"                                                              // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Nazwa użytkownika lub email",                                                                     // 63
		username: "Nazwa użytkownika",                                                                                      // 64
		email: "Email",                                                                                                     // 65
		password: "Hasło"                                                                                                   // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Nazwa użytkownika",                                                                                      // 69
		email: "Email",                                                                                                     // 70
		emailOpt: "Email (opcjonalnie)",                                                                                    // 71
		password: "Hasło",                                                                                                  // 72
		passwordAgain: "Hasło (powtórz)"                                                                                    // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Obecne hasło",                                                                                    // 76
		newPassword: "Nowe hasło",                                                                                          // 77
		newPasswordAgain: "Nowe hasło (powtórz)"                                                                            // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		emailSent: "Wysłano email",                                                                                         // 81
		passwordChanged: "Hasło zostało zmienione"                                                                          // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Wystąpił błąd",                                                                                      // 85
		userNotFound: "Nie znaleziono użytkownika",                                                                         // 86
		invalidEmail: "Niepoprawny email",                                                                                  // 87
		incorrectPassword: "Niepoprawne hasło",                                                                             // 88
		usernameTooShort: "Nazwa użytkowika powinna mieć przynajmniej 3 znaki",                                             // 89
		passwordTooShort: "Hasło powinno składać się przynajmnej z 6 znaków",                                               // 90
		passwordsDontMatch: "Hasło są różne",                                                                               // 91
		newPasswordSameAsOld: "Nowe hasło musi się różnić od starego",                                                      // 92
		signupsForbidden: "Rejstracja wyłączona"                                                                            // 93
	}                                                                                                                    // 94
});                                                                                                                   // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/zh-CN.i18n.js                                                            //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
zhCN = {                                                                                                              // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "重置密码",                                                                                                      // 3
		newPassword: "新密码",                                                                                                 // 4
		newPasswordAgain: "重复新密码",                                                                                          // 5
		cancel: "取消",                                                                                                       // 6
		submit: "确定"                                                                                                        // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "选择一个密码",                                                                                                    // 10
		newPassword: "新密码",                                                                                                 // 11
		newPasswordAgain: "重复新密码",                                                                                          // 12
		cancel: "取消",                                                                                                       // 13
		submit: "确定"                                                                                                        // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Email地址验证",                                                                                              // 17
		dismiss: "失败"                                                                                                       // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "失败"                                                                                                       // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "更改密码",                                                                                                   // 24
		signOut: "退出"                                                                                                       // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "登录",                                                                                                       // 28
		up: "注册"                                                                                                            // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "或"                                                                                                             // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "新建",                                                                                                       // 35
		signIn: "登陆",                                                                                                       // 36
		forgot: "忘记密码?",                                                                                                    // 37
		createAcc: "新建用户"                                                                                                   // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Email",                                                                                                     // 41
		reset: "重置密码",                                                                                                      // 42
		invalidEmail: "email格式不正确"                                                                                          // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "取消"                                                                                                          // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "更改密码",                                                                                                     // 49
		cancel: "取消"                                                                                                        // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "登陆以",                                                                                                  // 53
		configure: "配置"                                                                                                     // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "退出"                                                                                                       // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "没有配置登录服务"                                                                                         // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "用户名或者Email",                                                                                      // 63
		username: "用户名",                                                                                                    // 64
		email: "Email",                                                                                                     // 65
		password: "密码"                                                                                                      // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "用户名",                                                                                                    // 69
		email: "Email",                                                                                                     // 70
		emailOpt: "Email (可选)",                                                                                             // 71
		password: "密码",                                                                                                     // 72
		passwordAgain: "重复密码"                                                                                               // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "当前密码",                                                                                            // 76
		newPassword: "新密码",                                                                                                 // 77
		newPasswordAgain: "重复新密码"                                                                                           // 78
	},                                                                                                                   // 79
	infoMessages: {                                                                                                      // 80
		emailSent: "发送Email",                                                                                               // 81
		passwordChanged: "密码更改成功"                                                                                           // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "出現了一個錯誤",                                                                                            // 85
		userNotFound: "用户不存在",                                                                                              // 86
		invalidEmail: "Email格式不正确",                                                                                         // 87
		incorrectPassword: "密码错误",                                                                                          // 88
		usernameTooShort: "用户名必需大于3位",                                                                                      // 89
		passwordTooShort: "密码必需大于6位",                                                                                       // 90
		passwordsDontMatch: "密码输入不一致",                                                                                      // 91
		newPasswordSameAsOld: "新密码和旧的不能一样",                                                                                 // 92
		signupsForbidden: "没有权限"                                                                                            // 93
	}                                                                                                                    // 94
};                                                                                                                    // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/zh-TW.i18n.js                                                            //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
zhTW = {                                                                                                              // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "重設密碼",                                                                                                      // 3
		newPassword: "新密碼",                                                                                                 // 4
		newPasswordAgain: "重複新密碼",                                                                                          // 5
		cancel: "取消",                                                                                                       // 6
		submit: "確定"                                                                                                        // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "選擇一個密碼",                                                                                                    // 10
		newPassword: "新密碼",                                                                                                 // 11
		newPasswordAgain: "重複新密碼",                                                                                          // 12
		cancel: "取消",                                                                                                       // 13
		submit: "確定"                                                                                                        // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Email驗證",                                                                                                // 17
		dismiss: "失敗"                                                                                                       // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "失敗"                                                                                                       // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "更改密碼",                                                                                                   // 24
		signOut: "登出"                                                                                                       // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "登入",                                                                                                       // 28
		up: "註冊"                                                                                                            // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "或"                                                                                                             // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "新建",                                                                                                       // 35
		signIn: "登入",                                                                                                       // 36
		forgot: "忘记密碼?",                                                                                                    // 37
		createAcc: "新建用户"                                                                                                   // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Email",                                                                                                     // 41
		reset: "重設密碼",                                                                                                      // 42
		invalidEmail: "email格式不正確"                                                                                          // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "取消"                                                                                                          // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "更改密碼",                                                                                                     // 49
		cancel: "取消"                                                                                                        // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "登入以",                                                                                                  // 53
		configure: "設定"                                                                                                     // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "退出"                                                                                                       // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "沒有設定登入服务"                                                                                         // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "用户名或者Email",                                                                                      // 63
		username: "用户名",                                                                                                    // 64
		email: "Email",                                                                                                     // 65
		password: "密碼"                                                                                                      // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "用户名",                                                                                                    // 69
		email: "Email",                                                                                                     // 70
		emailOpt: "Email (可選)",                                                                                             // 71
		password: "密碼",                                                                                                     // 72
		passwordAgain: "重複密碼"                                                                                               // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "目前密碼",                                                                                            // 76
		newPassword: "新密碼",                                                                                                 // 77
		newPasswordAgain: "重複新密碼"                                                                                           // 78
	},                                                                                                                   // 79
	infoMessages: {                                                                                                      // 80
		emailSent: "發送Email",                                                                                               // 81
		passwordChanged: "密碼更改成功"                                                                                           // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "出現了一個錯誤",                                                                                            // 85
		userNotFound: "用户不存在",                                                                                              // 86
		invalidEmail: "Email格式不正確",                                                                                         // 87
		incorrectPassword: "密碼错误",                                                                                          // 88
		usernameTooShort: "用户名必需大于3位",                                                                                      // 89
		passwordTooShort: "密碼必需大于6位",                                                                                       // 90
		passwordsDontMatch: "密碼输入不一致",                                                                                      // 91
		newPasswordSameAsOld: "新密碼和舊的不能一樣",                                                                                 // 92
		signupsForbidden: "沒有權限"                                                                                            // 93
	}                                                                                                                    // 94
};                                                                                                                    // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/zh.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("zh", zhCN);                                                                                                 // 1
i18n.map("zh-CN", zhCN);                                                                                              // 2
i18n.map("zh-TW", zhTW);                                                                                              // 3
                                                                                                                      // 4
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/nl.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("nl", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Wachtwoord resetten",                                                                                       // 3
		newPassword: "Nieuw wachtwoord",                                                                                    // 4
		newPasswordAgain: "Nieuw wachtwoord (opnieuw)",                                                                     // 5
		cancel: "Annuleren",                                                                                                // 6
		submit: "Wachtwoord instellen"                                                                                      // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Stel een wachtwoord in",                                                                                    // 10
		newPassword: "Nieuw wachtwoord",                                                                                    // 11
		newPasswordAgain: "Nieuw wachtwoord (opnieuw)",                                                                     // 12
		cancel: "Sluiten",                                                                                                  // 13
		submit: "Wachtwoord instellen"                                                                                      // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "E-mailadres geverifieerd",                                                                               // 17
		dismiss: "Sluiten"                                                                                                  // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Sluiten",                                                                                                 // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Wachtwoord veranderen",                                                                                  // 24
		signOut: "Afmelden"                                                                                                 // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Aanmelden",                                                                                                // 28
		up: "Registreren"                                                                                                   // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "of"                                                                                                            // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Aanmaken",                                                                                                 // 35
		signIn: "Aanmelden",                                                                                                // 36
		forgot: "Wachtwoord vergeten?",                                                                                     // 37
		createAcc: "Account aanmaken"                                                                                       // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "E-mailadres",                                                                                               // 41
		reset: "Wachtwoord opnieuw instellen",                                                                              // 42
		invalidEmail: "Ongeldig e-mailadres"                                                                                // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Annuleren"                                                                                                   // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Wachtwoord veranderen",                                                                                    // 49
		cancel: "Annuleren"                                                                                                 // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Aanmelden via",                                                                                        // 53
		configure: "Instellen",                                                                                             // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Afmelden"                                                                                                 // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Geen aanmelddienst ingesteld"                                                                     // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Gebruikersnaam of e-mailadres",                                                                   // 63
		username: "Gebruikersnaam",                                                                                         // 64
		email: "E-mailadres",                                                                                               // 65
		password: "Wachtwoord"                                                                                              // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Gebruikersnaam",                                                                                         // 69
		email: "E-mailadres",                                                                                               // 70
		emailOpt: "E-mailadres (niet verplicht)",                                                                           // 71
		password: "Wachtwoord",                                                                                             // 72
		passwordAgain: "Wachtwoord (opnieuw)"                                                                               // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Huidig wachtwoord",                                                                               // 76
		newPassword: "Nieuw wachtwoord",                                                                                    // 77
		newPasswordAgain: "Nieuw wachtwoord (opnieuw)"                                                                      // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		emailSent: "E-mail verstuurd",                                                                                      // 81
		passwordChanged: "Wachtwoord veranderd"                                                                             // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Er is een fout opgetreden",                                                                          // 85
		userNotFound: "Gebruiker niet gevonden",                                                                            // 86
		invalidEmail: "Ongeldig e-mailadres",                                                                               // 87
		incorrectPassword: "Onjuist wachtwoord",                                                                            // 88
		usernameTooShort: "De gebruikersnaam moet minimaal uit 3 tekens bestaan",                                           // 89
		passwordTooShort: "Het wachtwoord moet minimaal uit 6 tekens bestaan",                                              // 90
		passwordsDontMatch: "De wachtwoorden komen niet overeen",                                                           // 91
		newPasswordSameAsOld: "Het oude en het nieuwe wachtwoord mogen niet hetzelfde zijn",                                // 92
		signupsForbidden: "Aanmeldingen niet toegestaan"                                                                    // 93
	}                                                                                                                    // 94
});                                                                                                                   // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/ja.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("ja", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "パスワードを再設定する",                                                                                               // 3
		newPassword: "新しいパスワード",                                                                                            // 4
		newPasswordAgain: "新しいパスワード (確認)",                                                                                  // 5
		cancel: "キャンセル",                                                                                                    // 6
		submit: "パスワードを設定"                                                                                                  // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "パスワードを決める",                                                                                                 // 10
		newPassword: "新しいパスワード",                                                                                            // 11
		newPasswordAgain: "新しいパスワード (確認)",                                                                                  // 12
		cancel: "閉じる",                                                                                                      // 13
		submit: "パスワードを設定"                                                                                                  // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "メールアドレス菅確認されました",                                                                                        // 17
		dismiss: "閉じる"                                                                                                      // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "閉じる",                                                                                                     // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "パスワードを変える",                                                                                              // 24
		signOut: "サインアウト"                                                                                                   // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "サインイン",                                                                                                    // 28
		up: "参加"                                                                                                            // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "または"                                                                                                           // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "作成",                                                                                                       // 35
		signIn: "サインイン",                                                                                                    // 36
		forgot: "パスワードを忘れましたか?",                                                                                            // 37
		createAcc: "アカウントを作成"                                                                                               // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "メール",                                                                                                       // 41
		reset: "パスワードを再設定する",                                                                                               // 42
		invalidEmail: "無効なメール"                                                                                              // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "キャンセル"                                                                                                       // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "パスワードを変える",                                                                                                // 49
		cancel: "キャンセル"                                                                                                     // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "サインインする",                                                                                              // 53
		configure: "設定する",                                                                                                  // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "サインアウト"                                                                                                   // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "ログインサービスが設定されていません"                                                                               // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "ユーザー名またはメール",                                                                                     // 63
		username: "ユーザー名",                                                                                                  // 64
		email: "メール",                                                                                                       // 65
		password: "パスワード"                                                                                                   // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "ユーザー名",                                                                                                  // 69
		email: "メール",                                                                                                       // 70
		emailOpt: "メール (オプション)",                                                                                            // 71
		password: "パスワード",                                                                                                  // 72
		passwordAgain: "パスワード (確認)"                                                                                         // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "現在のパスワード",                                                                                        // 76
		newPassword: "新しいパスワード",                                                                                            // 77
		newPasswordAgain: "新しいパスワード (確認)"                                                                                   // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		emailSent: "メールを送りました",                                                                                             // 81
		passwordChanged: "パスワードが変わりました"                                                                                     // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "エラーが発生しました",                                                                                         // 85
		userNotFound: "ユーザーが見つかりません",                                                                                       // 86
		invalidEmail: "無効なメール",                                                                                             // 87
		incorrectPassword: "無効なパスワード",                                                                                      // 88
		usernameTooShort: "ユーザー名は3文字以上でなければいけません",                                                                          // 89
		passwordTooShort: "パスワードは6文字以上でなければいけません",                                                                          // 90
		passwordsDontMatch: "パスワードが一致しません",                                                                                 // 91
		newPasswordSameAsOld: "新しいパスワードは古いパスワードと違っていなければいけません",                                                             // 92
		signupsForbidden: "サインアップが許可されませんでした"                                                                               // 93
	}                                                                                                                    // 94
});                                                                                                                   // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/he.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("he", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "איפוס סיסמא",                                                                                               // 3
		newPassword: "סיסמא חדשה",                                                                                          // 4
		newPasswordAgain: "סיסמא חדשה (שוב)",                                                                               // 5
		cancel: "ביטול",                                                                                                    // 6
		submit: "קביעת סיסמא"                                                                                               // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "בחירת סיסמא",                                                                                               // 10
		newPassword: "סיסמא חדשה",                                                                                          // 11
		newPasswordAgain: "סיסמא חדשה (שוב)",                                                                               // 12
		cancel: "סגירה",                                                                                                    // 13
		submit: "קביעת סיסמא"                                                                                               // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "כתובת דואר אלקטרוני אומתה",                                                                              // 17
		dismiss: "סגירה"                                                                                                    // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "סגירה",                                                                                                   // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "שינוי סיסמא",                                                                                            // 24
		signOut: "יציאה"                                                                                                    // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "כניסה",                                                                                                    // 28
		up: "רישום"                                                                                                         // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "או"                                                                                                            // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "יצירה",                                                                                                    // 35
		signIn: "התחברות",                                                                                                  // 36
		forgot: "סיסמא נשכחה?",                                                                                             // 37
		createAcc: "יצירת חשבון"                                                                                            // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "דואר אלקטרוני (אימייל)",                                                                                    // 41
		reset: "איפוס סיסמא",                                                                                               // 42
		invalidEmail: "כתובת דואר אלקטרוני לא חוקית"                                                                        // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "ביטול"                                                                                                       // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "שינוי סיסמא",                                                                                              // 49
		cancel: "ביטול"                                                                                                     // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "התחברות עםh",                                                                                          // 53
		configure: "הגדרות",                                                                                                // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "התנתקות"                                                                                                  // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "שירותים התחברות נוספים לא הוגדרו"                                                                 // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "שם משתמש או כתובת דואר אלקטרוני",                                                                 // 63
		username: "שם משתמש",                                                                                               // 64
		email: "דואר אלקטרוני",                                                                                             // 65
		password: "סיסמא"                                                                                                   // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "שם משתמש",                                                                                               // 69
		email: "דואר אלקטרוני",                                                                                             // 70
		emailOpt: "דואר אלקטרוני (לא חובה)",                                                                                // 71
		password: "סיסמא",                                                                                                  // 72
		passwordAgain: "סיסמא (שוב)"                                                                                        // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "סיסמא נוכחית",                                                                                    // 76
		newPassword: "סיסמא חדשה",                                                                                          // 77
		newPasswordAgain: "סיסמא חדשה (שוב)"                                                                                // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		emailSent: "דואר אלקטרוני נשלח",                                                                                    // 81
		passwordChanged: "סיסמא שונתה"                                                                                      // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "תרעה שגיאה",                                                                                         // 85
		userNotFound: "משתמש/ת לא קיים/מת",                                                                                 // 86
		invalidEmail: "כתובת דואר אלקטרוני לא חוקי",                                                                        // 87
		incorrectPassword: "סיסמא שגויה",                                                                                   // 88
		usernameTooShort: "שם משתמש חייב להיות בן 3 תוים לפחות",                                                            // 89
		passwordTooShort: "סיסמא חייבת להיות בת 6 תווים לפחות",                                                             // 90
		passwordsDontMatch: "סיסמאות לא תואמות",                                                                            // 91
		newPasswordSameAsOld: "סיסמא חדשה וישנה חייבות להיות שונות",                                                        // 92
		signupsForbidden: "אין אפשרות לרישום"                                                                               // 93
	}                                                                                                                    // 94
});                                                                                                                   // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/sv.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("sv", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Återställ ditt lösenord",                                                                                   // 3
		newPassword: "Nytt lösenord",                                                                                       // 4
		cancel: "Avbryt",                                                                                                   // 5
		submit: "Spara lösenord"                                                                                            // 6
	},                                                                                                                   // 7
	enrollAccountDialog: {                                                                                               // 8
		title: "Välj ett lösenord",                                                                                         // 9
		newPassword: "Nytt lösenord",                                                                                       // 10
		cancel: "Stäng",                                                                                                    // 11
		submit: "Spara lösenord"                                                                                            // 12
	},                                                                                                                   // 13
	justVerifiedEmailDialog: {                                                                                           // 14
		verified: "Epostadressen verifierades",                                                                             // 15
		dismiss: "Avfärda"                                                                                                  // 16
	},                                                                                                                   // 17
	loginButtonsMessagesDialog: {                                                                                        // 18
		dismiss: "Avfärda",                                                                                                 // 19
	},                                                                                                                   // 20
	loginButtonsLoggedInDropdownActions: {                                                                               // 21
		password: "Byt lösenord",                                                                                           // 22
		signOut: "Logga ut"                                                                                                 // 23
	},                                                                                                                   // 24
	loginButtonsLoggedOutDropdown: {                                                                                     // 25
		signIn: "Logga in",                                                                                                 // 26
		up: "Skapa konto"                                                                                                   // 27
	},                                                                                                                   // 28
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 29
		or: "eller"                                                                                                         // 30
	},                                                                                                                   // 31
	loginButtonsLoggedOutPasswordService: {                                                                              // 32
		create: "Skapa",                                                                                                    // 33
		signIn: "Logga in",                                                                                                 // 34
		forgot: "Glömt ditt lösenord?",                                                                                     // 35
		createAcc: "Skapa konto"                                                                                            // 36
	},                                                                                                                   // 37
	forgotPasswordForm: {                                                                                                // 38
		email: "E-postadress",                                                                                              // 39
		reset: "Återställ lösenord",                                                                                        // 40
		sent: "E-post skickat",                                                                                             // 41
		invalidEmail: "Ogiltig e-postadress"                                                                                // 42
	},                                                                                                                   // 43
	loginButtonsBackToLoginLink: {                                                                                       // 44
		back: "Avbryt"                                                                                                      // 45
	},                                                                                                                   // 46
	loginButtonsChangePassword: {                                                                                        // 47
		submit: "Byt lösenord",                                                                                             // 48
		cancel: "Avbryt"                                                                                                    // 49
	},                                                                                                                   // 50
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 51
		signInWith: "Logga in med",                                                                                         // 52
		configure: "Konfigurera",                                                                                           // 53
	},                                                                                                                   // 54
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 55
		signOut: "Logga ut"                                                                                                 // 56
	},                                                                                                                   // 57
	loginButtonsLoggedOut: {                                                                                             // 58
		noLoginServices: "Inga inloggningstjänster har konfigurerats"                                                       // 59
	},                                                                                                                   // 60
	loginFields: {                                                                                                       // 61
		usernameOrEmail: "Användarnamn eller e-postadress",                                                                 // 62
		username: "Användarnamn",                                                                                           // 63
		email: "E-postadress",                                                                                              // 64
		password: "Lösenord"                                                                                                // 65
	},                                                                                                                   // 66
	signupFields: {                                                                                                      // 67
		username: "Användarnamn",                                                                                           // 68
		email: "E-postadress",                                                                                              // 69
		emailOpt: "E-postadress (valfritt)",                                                                                // 70
		password: "Lösenord",                                                                                               // 71
		passwordAgain: "Upprepa lösenord"                                                                                   // 72
	},                                                                                                                   // 73
	changePasswordFields: {                                                                                              // 74
		currentPassword: "Nuvarande lösenord",                                                                              // 75
		newPassword: "Nytt lösenord",                                                                                       // 76
		newPasswordAgain: "Upprepa nytt lösenord"                                                                           // 77
	},                                                                                                                   // 78
	infoMessages : {                                                                                                     // 79
		emailSent: "E-post skickat",                                                                                        // 80
		passwordChanged: "Lösenord ändrat"                                                                                  // 81
	},                                                                                                                   // 82
	errorMessages: {                                                                                                     // 83
		genericTitle: "Ett fel har uppstått",                                                                               // 84
		userNotFound: "Ingen användare hittades",                                                                           // 85
		invalidEmail: "Ogiltig e-postadress",                                                                               // 86
		incorrectPassword: "Fel lösenord",                                                                                  // 87
		usernameTooShort: "Användarnamnet måste vara minst 3 tecken långt",                                                 // 88
		passwordTooShort: "Lösenordet bör vara längre än 6 tecken",                                                         // 89
		passwordsDontMatch: "Lösenorden matchar inte",                                                                      // 90
		newPasswordSameAsOld: "Den nya och gamla lösenordet bör inte vara samma",                                           // 91
		signupsForbidden: "Sign up är inte tillåtet"                                                                        // 92
	}                                                                                                                    // 93
});                                                                                                                   // 94
                                                                                                                      // 95
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/uk.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map('uk', {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Скинути пароль",                                                                                            // 3
		newPassword: "Новий пароль",                                                                                        // 4
		newPasswordAgain: "Новий пароль (ще раз)",                                                                          // 5
		cancel: "Скасувати",                                                                                                // 6
		submit: "Зберегти пароль"                                                                                           // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Обрати пароль",                                                                                             // 10
		newPassword: "Новий пароль",                                                                                        // 11
		newPasswordAgain: "Новий пароль (ще раз)",                                                                          // 12
		cancel: "Скасувати",                                                                                                // 13
		submit: "Зберегти пароль"                                                                                           // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Email підтверджено",                                                                                     // 17
			dismiss: "Закрити"                                                                                                 // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
			dismiss: "Закрити"                                                                                                 // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Змінити пароль",                                                                                         // 24
		signOut: "Вийти"                                                                                                    // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Ввійти",                                                                                                   // 28
		up: "Зареєструватись"                                                                                               // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "або"                                                                                                           // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Створити",                                                                                                 // 35
		signIn: "Ввійти",                                                                                                   // 36
		forgot: "Забули пароль?",                                                                                           // 37
		createAcc: "Створити аккаунт"                                                                                       // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Email",                                                                                                     // 41
		reset: "Скинути пароль",                                                                                            // 42
		invalidEmail: "Некорректный email"                                                                                  // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Скасувати"                                                                                                   // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Змінити пароль",                                                                                           // 49
		cancel: "Скасувати"                                                                                                 // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Ввійти через",                                                                                         // 53
		configure: "Налаштувати вхід через",                                                                                // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Вийти"                                                                                                    // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Сервіс для входу не налаштований"                                                                 // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Им’я користувача або email",                                                                      // 63
		username: "Им’я користувача",                                                                                       // 64
		email: "Email",                                                                                                     // 65
		password: "Пароль"                                                                                                  // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Им’я користувача",                                                                                       // 69
		email: "Email",                                                                                                     // 70
		emailOpt: "Email (необов’язковий)",                                                                                 // 71
		password: "Пароль",                                                                                                 // 72
		passwordAgain: "Пароль (ще раз)"                                                                                    // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Поточний пароль",                                                                                 // 76
		newPassword: "Новий пароль",                                                                                        // 77
		newPasswordAgain: "Новий пароль (ще раз)"                                                                           // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		sent: "Вам відправлено лист",                                                                                       // 81
		passwordChanged: "Пароль змінено"                                                                                   // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Там була помилка",                                                                                   // 85
		userNotFound: "Користувача не знайдено",                                                                            // 86
		invalidEmail: "Некорректний email",                                                                                 // 87
		incorrectPassword: "Невірний пароль",                                                                               // 88
		usernameTooShort: "Им’я користувача повинно бути довжиною не менше 3-ьох символів",                                 // 89
		passwordTooShort: "Пароль повинен бути довжиною не менше 6-ти символів",                                            // 90
		passwordsDontMatch: "Паролі не співпадають",                                                                        // 91
		newPasswordSameAsOld: "Новий та старий паролі повинні бути різними"                                                 // 92
	}                                                                                                                    // 93
});                                                                                                                   // 94
                                                                                                                      // 95
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/fi.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("fi", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Palauta salasana",                                                                                          // 3
		newPassword: "Uusi salasana",                                                                                       // 4
		newPasswordAgain: "Uusi salasana (uudestaan)",                                                                      // 5
		cancel: "Peruuta",                                                                                                  // 6
		submit: "Aseta salasana"                                                                                            // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Valitse salasana",                                                                                          // 10
		newPassword: "Uusi salasana",                                                                                       // 11
		newPasswordAgain: "Uusi salasana (uudestaan)",                                                                      // 12
		cancel: "Sulje",                                                                                                    // 13
		submit: "Aseta salasana"                                                                                            // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Sähköpostiosoite vahvistettu",                                                                           // 17
		dismiss: "Sulje"                                                                                                    // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Sulje",                                                                                                   // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Vaihda salasana",                                                                                        // 24
		signOut: "Kirjaudu ulos"                                                                                            // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Kirjaudu",                                                                                                 // 28
		up: "Rekisteröidy"                                                                                                  // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "tai"                                                                                                           // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Luo",                                                                                                      // 35
		signIn: "Kirjaudu",                                                                                                 // 36
		forgot: "Unohditko salasanasi?",                                                                                    // 37
		createAcc: "Luo tunnus"                                                                                             // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Sähköpostiosoite",                                                                                          // 41
		reset: "Palauta salasana",                                                                                          // 42
		invalidEmail: "Virheellinen sähköpostiosoite"                                                                       // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Peruuta"                                                                                                     // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Vaihda salasana",                                                                                          // 49
		cancel: "Peruuta"                                                                                                   // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Kirjaudu käyttäen",                                                                                    // 53
		configure: "Konfiguroi",                                                                                            // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Kirjaudu ulos"                                                                                            // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Kirjautumispalveluita ei konfiguroituna"                                                          // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Käyttäjätunnus tai sähköpostiosoite",                                                             // 63
		username: "Käyttäjätunnus",                                                                                         // 64
		email: "Sähköpostiosoite",                                                                                          // 65
		password: "Salasana"                                                                                                // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Käyttäjätunnus",                                                                                         // 69
		email: "Sähköpostiosoite",                                                                                          // 70
		emailOpt: "Sähköposti (vapaaehtoinen)",                                                                             // 71
		password: "Salasana",                                                                                               // 72
		passwordAgain: "Salasana (uudestaan)"                                                                               // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Nykyinen salasana",                                                                               // 76
		newPassword: "Uusi salasana",                                                                                       // 77
		newPasswordAgain: "Uusi salasana (uudestaan)"                                                                       // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		emailSent: "Sähköposti lähetetty",                                                                                  // 81
		passwordChanged: "Salasana vaihdettu"                                                                               // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Tapahtui virhe",                                                                                     // 85
		userNotFound: "Käyttäjää ei löytynyt",                                                                              // 86
		invalidEmail: "Virheellinen sähköpostiosoite",                                                                      // 87
		incorrectPassword: "Virheellinen salasana",                                                                         // 88
		usernameTooShort: "Käyttäjätunnuksen on oltava vähintään 3 merkkiä pitkä",                                          // 89
		passwordTooShort: "Salasanan on oltava vähintään 6 merkkiä pitkä",                                                  // 90
		passwordsDontMatch: "Salasanat eivät täsmää",                                                                       // 91
		newPasswordSameAsOld: "Uuden ja vanhan salasanan on oltava eri",                                                    // 92
		signupsForbidden: "Rekisteröityminen estetty"                                                                       // 93
	}                                                                                                                    // 94
});                                                                                                                   // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/vi.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("vi", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Đặt lại mật khẩu",                                                                                          // 3
		newPassword: "Mật khẩu mới",                                                                                        // 4
		newPasswordAgain: "Xác nhận mật khẩu mới",                                                                          // 5
		cancel: "Thoát",                                                                                                    // 6
		submit: "Lưu lại"                                                                                                   // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Cài đặt mật khẩu",                                                                                          // 10
		newPassword: "Mật khẩu mới",                                                                                        // 11
		newPasswordAgain: "Xác nhận mật khẩu mới",                                                                          // 12
		cancel: "Thoát",                                                                                                    // 13
		submit: "Lưu lại"                                                                                                   // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Địa chỉ Email đã được xác nhận.",                                                                        // 17
		dismiss: "Đóng"                                                                                                     // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Đóng",                                                                                                    // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Thay đổi mật khẩu",                                                                                      // 24
		signOut: "Đăng xuất"                                                                                                // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Đăng nhập",                                                                                                // 28
		up: "Đăng ký"                                                                                                       // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "hoặc"                                                                                                          // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Tạo mới",                                                                                                  // 35
		signIn: "Đăng nhập",                                                                                                // 36
		forgot: "Quên mật khẩu?",                                                                                           // 37
		createAcc: "Khởi tạo tài khoản"                                                                                     // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Email",                                                                                                     // 41
		reset: "Cài lại mật khẩu",                                                                                          // 42
		invalidEmail: "Email không hợp lệ"                                                                                  // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Thoát"                                                                                                       // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Lưu lại",                                                                                                  // 49
		cancel: "Thoát"                                                                                                     // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Đăng nhập bằng",                                                                                       // 53
		configure: "Cấu hình",                                                                                              // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Đăng xuất"                                                                                                // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Không có dịch vụ nào được cấu hình"                                                               // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Tên đăng nhập hoặc Email",                                                                        // 63
		username: "Tên đăg nhập",                                                                                           // 64
		email: "Email",                                                                                                     // 65
		password: "Mật khẩu"                                                                                                // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Tên đăng nhập",                                                                                          // 69
		email: "Email",                                                                                                     // 70
		emailOpt: "Email (Không bắt buộc)",                                                                                 // 71
		password: "Mật khẩu",                                                                                               // 72
		passwordAgain: "Xác nhận mật khẩu"                                                                                  // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Mật khẩu hiện tại",                                                                               // 76
		newPassword: "Mật khẩu mới",                                                                                        // 77
		newPasswordAgain: "Xác nhận mật khẩu mới"                                                                           // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		emailSent: "Gửi Email",                                                                                             // 81
		passwordChanged: "Mật khẩu đã được thay đổi"                                                                        // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Có lỗi xảy ra",                                                                                      // 85
		userNotFound: "Người dùng không tồn tại",                                                                           // 86
		invalidEmail: "Email không hợp lệ",                                                                                 // 87
		incorrectPassword: "Sai mật khẩu",                                                                                  // 88
		usernameTooShort: "Tên đăng nhập phải có ít nhất 3 ký tự",                                                          // 89
		passwordTooShort: "Mật khẩu phải có ít nhất 6 ký tự",                                                               // 90
		passwordsDontMatch: "Xác nhận mật khẩu không khớp",                                                                 // 91
		newPasswordSameAsOld: "Mật khẩu mới và cũ phải khác nhau",                                                          // 92
		signupsForbidden: "Tạm khoá đăng ký"                                                                                // 93
	}                                                                                                                    // 94
});                                                                                                                   // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/sk.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map('sk', {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: 'Obnovenie hesla',                                                                                           // 3
		newPassword: 'Nové heslo',                                                                                          // 4
		newPasswordAgain: 'Nové heslo (opakujte)',                                                                          // 5
		cancel: 'Zrušiť',                                                                                                   // 6
		submit: 'Zmeniť heslo'                                                                                              // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: 'Nastaviť heslo',                                                                                            // 10
		newPassword: 'Nové heslo',                                                                                          // 11
		newPasswordAgain: 'Nové heslo (opakujte)',                                                                          // 12
		cancel: 'Zatvoriť',                                                                                                 // 13
		submit: 'Nastaviť heslo'                                                                                            // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: 'Emailová adresa overená',                                                                                // 17
		dismiss: 'Zavrieť'                                                                                                  // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: 'Zrušiť'                                                                                                   // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: 'Zmeniť heslo',                                                                                           // 24
		signOut: 'Odhlásiť'                                                                                                 // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: 'Prihlásenie',                                                                                              // 28
		up: 'Registrovať'                                                                                                   // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: 'alebo'                                                                                                         // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: 'Vytvoriť',                                                                                                 // 35
		signIn: 'Prihlásiť',                                                                                                // 36
		forgot: 'Zabudli ste heslo?',                                                                                       // 37
		createAcc: 'Vytvoriť účet'                                                                                          // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: 'Email',                                                                                                     // 41
		reset: 'Obnoviť heslo',                                                                                             // 42
		invalidEmail: 'Nesprávný email'                                                                                     // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: 'Späť'                                                                                                        // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: 'Zmeniť heslo',                                                                                             // 49
		cancel: 'Zrušiť'                                                                                                    // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: 'Prihlasiť s',                                                                                          // 53
		configure: 'Nastaviť'                                                                                               // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: 'Odhlásiť'                                                                                                 // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: 'Žiadne prihlasovacie služby'                                                                      // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: 'Užívateľské meno alebo email',                                                                    // 63
		username: 'Užívateľské meno',                                                                                       // 64
		email: 'Email',                                                                                                     // 65
		password: 'Heslo'                                                                                                   // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: 'Užívateľské meno',                                                                                       // 69
		email: 'Email',                                                                                                     // 70
		emailOpt: 'Email (voliteľné)',                                                                                      // 71
		password: 'Heslo',                                                                                                  // 72
		passwordAgain: 'Heslo (opakujte)'                                                                                   // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: 'Vaše heslo',                                                                                      // 76
		newPassword: 'Nové heslo',                                                                                          // 77
		newPasswordAgain: 'Nové heslo (opakujte)'                                                                           // 78
	},                                                                                                                   // 79
	infoMessages: {                                                                                                      // 80
		emailSent: 'Email odoslaný',                                                                                        // 81
		passwordChanged: 'Heslo zmenené'                                                                                    // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: 'Nastala chyba',                                                                                      // 85
		userNotFound: 'Užívateľ sa nenašiel',                                                                               // 86
		invalidEmail: 'Nesprávný email',                                                                                    // 87
		incorrectPassword: 'Nesprávne heslo',                                                                               // 88
		usernameTooShort: 'Užívateľské meno musi obsahovať minimálne 3 znaky',                                              // 89
		passwordTooShort: 'Heslo musi obsahovať minimálne 6 znakov',                                                        // 90
		passwordsDontMatch: 'Hesla sa nezhodujú',                                                                           // 91
		newPasswordSameAsOld: 'Staré a nové hesla sa nezhodujú',                                                            // 92
		signupsForbidden: 'Prihlasovanie je zakázané'                                                                       // 93
	}                                                                                                                    // 94
});                                                                                                                   // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/be.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map('be', {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Скасаваць пароль",                                                                                          // 3
		newPassword: "Новы пароль",                                                                                         // 4
		newPasswordAgain: "Новы пароль (яшче раз)",                                                                         // 5
		cancel: "Скасаваць",                                                                                                // 6
		submit: "Захаваць пароль"                                                                                           // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Выбраць пароль",                                                                                            // 10
		newPassword: "Новы пароль",                                                                                         // 11
		newPasswordAgain: "Новы пароль (яшче раз)",                                                                         // 12
		cancel: "Скасаваць",                                                                                                // 13
		submit: "Захаваць пароль"                                                                                           // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Email пацьверджаны",                                                                                     // 17
		dismiss: "Закрыць"                                                                                                  // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Закрыць"                                                                                                  // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Зьмяніць пароль",                                                                                        // 24
		signOut: "Выйсьці"                                                                                                  // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Увайсьці",                                                                                                 // 28
		up: "Зарэгістравацца"                                                                                               // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "або"                                                                                                           // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Стварыць",                                                                                                 // 35
		signIn: "Увайсьці",                                                                                                 // 36
		forgot: "Забылі пароль?",                                                                                           // 37
		createAcc: "Стварыць акаўнт"                                                                                        // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Email",                                                                                                     // 41
		reset: "Ськінуць пароль",                                                                                           // 42
		invalidEmail: "Некарэктны email"                                                                                    // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Скасаваць"                                                                                                   // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Зьмяніць пароль",                                                                                          // 49
		cancel: "Скасаваць"                                                                                                 // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Увайсьці праз",                                                                                        // 53
		configure: "Наладзіць уваход праз",                                                                                 // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Выйсьці"                                                                                                  // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Сервіс для ўваходу не наладжаны"                                                                  // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Імя карыстальніка або email",                                                                     // 63
		username: "Імя карыстальніка",                                                                                      // 64
		email: "Email",                                                                                                     // 65
		password: "Пароль"                                                                                                  // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Імя карыстальніка",                                                                                      // 69
		email: "Email",                                                                                                     // 70
		emailOpt: "Email (неабавязковы)",                                                                                   // 71
		password: "Пароль",                                                                                                 // 72
		passwordAgain: "Пароль (яшче раз)"                                                                                  // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Цяперашні пароль",                                                                                // 76
		newPassword: "Новы пароль",                                                                                         // 77
		newPasswordAgain: "Новы пароль (яшче раз)"                                                                          // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		sent: "Вам высланы ліст",                                                                                           // 81
		passwordChanged: "Пароль зьменены"                                                                                  // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Там была памылка",                                                                                   // 85
		userNotFound: "Карыстальнік не знойдзены",                                                                          // 86
		invalidEmail: "Некарэктны email",                                                                                   // 87
		incorrectPassword: "Няверны пароль",                                                                                // 88
		usernameTooShort: "Імя карыстальніка павінна быць даўжынёю не меней 3-ох літараў",                                  // 89
		passwordTooShort: "Пароль павінна быць даўжынёю не меней 6-ці літараў",                                             // 90
		passwordsDontMatch: "Паролі не аднолькавыя",                                                                        // 91
		newPasswordSameAsOld: "Новы і стары паролі павінны быць рознымі"                                                    // 92
	}                                                                                                                    // 93
});                                                                                                                   // 94
                                                                                                                      // 95
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/fa.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("fa", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "رمز عبور خود را تنظیم مجدد",                                                                                // 3
		newPassword: "رمز عبور جدید",                                                                                       // 4
		newPasswordAgain: "رمز عبور (دوباره)",                                                                              // 5
		cancel: "لغو",                                                                                                      // 6
		submit: "تنظیم رمز عبور"                                                                                            // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "یک رمز عبور انتخاب کنید",                                                                                   // 10
		newPassword: "رمز عبور جدید",                                                                                       // 11
		newPasswordAgain: "رمز عبور (دوباره)",                                                                              // 12
		cancel: "نزدیک",                                                                                                    // 13
		submit: "تنظیم رمز عبور"                                                                                            // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "آدرس ایمیل تایید",                                                                                       // 17
		dismiss: "پنهان کن"                                                                                                 // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "پنهان کن",                                                                                                // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "تغییر رمز عبور",                                                                                         // 24
		signOut: "خروج"                                                                                                     // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "ورود",                                                                                                     // 28
		up: "بپیوندید"                                                                                                      // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "یا"                                                                                                            // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "ايجاد كردن",                                                                                               // 35
		signIn: "ورود",                                                                                                     // 36
		forgot: "رمز عبور را فراموش کرده اید؟",                                                                             // 37
		createAcc: "ایجاد حساب کاربری"                                                                                      // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "ایمیل",                                                                                                     // 41
		reset: "تنظیم مجدد رمز ورود",                                                                                       // 42
		invalidEmail: "ایمیل نامعتبر"                                                                                       // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "لغو"                                                                                                         // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "تغییر رمز عبور",                                                                                           // 49
		cancel: "لغو"                                                                                                       // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "ورود به سیستم با",                                                                                     // 53
		configure: "پیکربندی",                                                                                              // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "خروج"                                                                                                     // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "بدون خدمات پیکربندی ورود"                                                                         // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "نام کاربری یا پست الکترونیک",                                                                     // 63
		username: "نام کاربری",                                                                                             // 64
		email: "ایمیل",                                                                                                     // 65
		password: "رمز عبور"                                                                                                // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "نام کاربری",                                                                                             // 69
		email: "ایمیل",                                                                                                     // 70
		emailOpt: "ایمیل اختیاری ",                                                                                         // 71
		password: "رمز عبور",                                                                                               // 72
			passwordAgain: "رمز عبور دوباره "                                                                                  // 73
	},                                                                                                                   // 74
		changePasswordFields: {                                                                                             // 75
		currentPassword: "رمز عبور فعلی",                                                                                   // 76
		newPassword: "رمز عبور جدید",                                                                                       // 77
		newPasswordAgain:  "رمز عبر (دوباره"                                                                                // 78
	},                                                                                                                   // 79
		infoMessages : {                                                                                                    // 80
		emailSent: "ایمیل ارسال",                                                                                           // 81
		passwordChanged: "رمز عبور تغییر کرد"                                                                               // 82
	},                                                                                                                   // 83
		errorMessages: {                                                                                                    // 84
		genericTitle: "یک خطای وجود دارد",                                                                                  // 85
	userNotFound: "کاربر پیدا نشد",                                                                                      // 86
	invalidEmail: "ایمیل نامعتبر",                                                                                       // 87
	incorrectPassword: "رمز عبور اشتباه",                                                                                // 88
	usernameTooShort: "نام کاربری حداقل باید 3 کاراکتر باشد",                                                            // 89
	passwordTooShort: "رمز عبور باید حداقل 6 کاراکتر باشد",                                                              // 90
	passwordsDontMatch: "کلمه عبور هماهنگ نیست",                                                                         // 91
	newPasswordSameAsOld: "کلمه عبور جدید و قدیمی باید متفاوت باشد",                                                     // 92
	signupsForbidden: "ثبت نام ممنوع"                                                                                    // 93
	}                                                                                                                    // 94
                                                                                                                      // 95
});                                                                                                                   // 96
                                                                                                                      // 97
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/sr-Cyrl.i18n.js                                                          //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
srCyrl = {                                                                                                            // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Ресетуј своју лозинку",                                                                                     // 3
		newPassword: "Нова лозинка",                                                                                        // 4
		newPasswordAgain: "Нова Лозинка (поново)",                                                                          // 5
		cancel: "Откажи",                                                                                                   // 6
		submit: "Постави лозинку"                                                                                           // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Одабери лозинку",                                                                                           // 10
		newPassword: "Нова лозинка",                                                                                        // 11
		newPasswordAgain: "Нова Лозинка (поново)",                                                                          // 12
		cancel: "Затвори",                                                                                                  // 13
		submit: "Постави лозинку"                                                                                           // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Адреса еПоште је проверена",                                                                             // 17
		dismiss: "Одбаци"                                                                                                   // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Одбаци",                                                                                                  // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Промени лозинку",                                                                                        // 24
		signOut: "Одјави се"                                                                                                // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Пријави се",                                                                                               // 28
		up: "Придружи се"                                                                                                   // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "или"                                                                                                           // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Креирај",                                                                                                  // 35
		signIn: "Пријави се",                                                                                               // 36
		forgot: "Заборавили сте лозинку?",                                                                                  // 37
		createAcc: "Креирај налог"                                                                                          // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "еАдреса",                                                                                                   // 41
		reset: "Ресетуј лозинку",                                                                                           // 42
		invalidEmail: "Неправилна еАдреса"                                                                                  // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Откажи"                                                                                                      // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Промени лозинку",                                                                                          // 49
		cancel: "Откажи"                                                                                                    // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Пријави се са",                                                                                        // 53
		configure: "Подеси",                                                                                                // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Одјави се"                                                                                                // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Пријавни сервиси нису подешени"                                                                   // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Корисничко име или еАдреса",                                                                      // 63
		username: "Корисничко име",                                                                                         // 64
		email: "еАдреса",                                                                                                   // 65
		password: "Лозинка"                                                                                                 // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Лозинка",                                                                                                // 69
		email: "еАдреса",                                                                                                   // 70
		emailOpt: "еАдреса (опционо)",                                                                                      // 71
		password: "Лозинка",                                                                                                // 72
		passwordAgain: "Лозинка (поново)"                                                                                   // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Тренутна Лозинка",                                                                                // 76
		newPassword: "Нова Лозинка",                                                                                        // 77
		newPasswordAgain: "Нова Лозинка (поново)"                                                                           // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		emailSent: "еПорука је послата",                                                                                    // 81
		passwordChanged: "Лозинка је промењена"                                                                             // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Појавила се једна грешка",                                                                           // 85
		userNotFound: "Корисник није пронађен",                                                                             // 86
		invalidEmail: "Неправилна еАдреса",                                                                                 // 87
		incorrectPassword: "Нетачна лозинка",                                                                               // 88
		usernameTooShort: "Корисничко име мора бити најмање 3 знака дуго",                                                  // 89
		passwordTooShort: "Лозинка мора бити најмање 6 знакова дуга",                                                       // 90
		passwordsDontMatch: "Лозинке се не поклапају",                                                                      // 91
		newPasswordSameAsOld: "Нова и стара лозинка морају бити различите",                                                 // 92
		signupsForbidden: "Пријаве забрањене"                                                                               // 93
	}                                                                                                                    // 94
};                                                                                                                    // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/sr-Latn.i18n.js                                                          //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
srLatn = {                                                                                                            // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Resetuj svoju lozinku",                                                                                     // 3
		newPassword: "Nova lozinka",                                                                                        // 4
		newPasswordAgain: "Nova Lozinka (ponovo)",                                                                          // 5
		cancel: "Otkaži",                                                                                                   // 6
		submit: "Postavi lozinku"                                                                                           // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Odaberi lozinku",                                                                                           // 10
		newPassword: "Nova lozinka",                                                                                        // 11
		newPasswordAgain: "Nova Lozinka (ponovo)",                                                                          // 12
		cancel: "Zatvori",                                                                                                  // 13
		submit: "Postavi loziknu"                                                                                           // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Adresa ePošte je proverena",                                                                             // 17
		dismiss: "Odbaci"                                                                                                   // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Odbaci",                                                                                                  // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Promeni lozinku",                                                                                        // 24
		signOut: "Odjavi se"                                                                                                // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Prijavi se",                                                                                               // 28
		up: "Pridruži se"                                                                                                   // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "ili"                                                                                                           // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Kreiraj",                                                                                                  // 35
		signIn: "Prijavi se",                                                                                               // 36
		forgot: "Zaboravili ste lozinku?",                                                                                  // 37
		createAcc: "Kreiraj nalog"                                                                                          // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "eAdresa",                                                                                                   // 41
		reset: "Resetuj lozinku",                                                                                           // 42
		invalidEmail: "Nepravilna eAdresa"                                                                                  // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Otkaži"                                                                                                      // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Promeni lozinku",                                                                                          // 49
		cancel: "Otkaži"                                                                                                    // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Prijavi se sa",                                                                                        // 53
		configure: "Podesi",                                                                                                // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Odjavi se"                                                                                                // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Prijavni servisi nisu podešeni"                                                                   // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Korisničko ime ili eAdresa",                                                                      // 63
		username: "Korisničko ime",                                                                                         // 64
		email: "eAdresa",                                                                                                   // 65
		password: "Lozinka"                                                                                                 // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Lozinka",                                                                                                // 69
		email: "eAdresa",                                                                                                   // 70
		emailOpt: "eAdresa (opciono)",                                                                                      // 71
		password: "Lozinka",                                                                                                // 72
		passwordAgain: "Lozinka (ponovo)"                                                                                   // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Trenutna Lozinkа",                                                                                // 76
		newPassword: "Nova Lozinka",                                                                                        // 77
		newPasswordAgain: "Nova Lozinka (ponovo)"                                                                           // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		emailSent: "ePoruka je poslata",                                                                                    // 81
		passwordChanged: "Lozinka je promenjena"                                                                            // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Pojavila se jedna greška",                                                                           // 85
		userNotFound: "Korisnik nije pronađen",                                                                             // 86
		invalidEmail: "Nepravilna eAdresa",                                                                                 // 87
		incorrectPassword: "Netačna lozinka",                                                                               // 88
		usernameTooShort: "Korisničko ime mora biti najmanje 3 znaka dugo",                                                 // 89
		passwordTooShort: "Lozinka mora biti najmanje 6 znakova duga",                                                      // 90
		passwordsDontMatch: "Lozinke se ne poklapaju",                                                                      // 91
		newPasswordSameAsOld: "Nova i stara lozinka moraju biti različite",                                                 // 92
		signupsForbidden: "Prijave zabranjene"                                                                              // 93
	}                                                                                                                    // 94
};                                                                                                                    // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/sr.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("sr", srCyrl);                                                                                               // 1
i18n.map("sr-Cyrl", srCyrl);                                                                                          // 2
i18n.map("sr-Latn", srLatn);                                                                                          // 3
                                                                                                                      // 4
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n/hu.i18n.js                                                               //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.map("hu", {                                                                                                      // 1
	resetPasswordDialog: {                                                                                               // 2
		title: "Új jelszó választása",                                                                                      // 3
		newPassword: "Új jelszó",                                                                                           // 4
		newPasswordAgain: "Új jelszó (újra)",                                                                               // 5
		cancel: "Mégsem",                                                                                                   // 6
		submit: "Tovább"                                                                                                    // 7
	},                                                                                                                   // 8
	enrollAccountDialog: {                                                                                               // 9
		title: "Jelszó választása",                                                                                         // 10
		newPassword: "Új jelszó",                                                                                           // 11
		newPasswordAgain: "Új jelszó (újra)",                                                                               // 12
		cancel: "Mégsem",                                                                                                   // 13
		submit: "Tovább"                                                                                                    // 14
	},                                                                                                                   // 15
	justVerifiedEmailDialog: {                                                                                           // 16
		verified: "Email cím megerősítve",                                                                                  // 17
		dismiss: "Ok"                                                                                                       // 18
	},                                                                                                                   // 19
	loginButtonsMessagesDialog: {                                                                                        // 20
		dismiss: "Ok",                                                                                                      // 21
	},                                                                                                                   // 22
	loginButtonsLoggedInDropdownActions: {                                                                               // 23
		password: "Jelszó módosítása",                                                                                      // 24
		signOut: "Kijelentkezés"                                                                                            // 25
	},                                                                                                                   // 26
	loginButtonsLoggedOutDropdown: {                                                                                     // 27
		signIn: "Bejelentkezés",                                                                                            // 28
		up: "Regisztráció"                                                                                                  // 29
	},                                                                                                                   // 30
	loginButtonsLoggedOutPasswordServiceSeparator: {                                                                     // 31
		or: "vagy"                                                                                                          // 32
	},                                                                                                                   // 33
	loginButtonsLoggedOutPasswordService: {                                                                              // 34
		create: "Regisztráció",                                                                                             // 35
		signIn: "Bejelentkezés",                                                                                            // 36
		forgot: "Elfelejtetted a jelszavadat?",                                                                             // 37
		createAcc: "Regisztráció"                                                                                           // 38
	},                                                                                                                   // 39
	forgotPasswordForm: {                                                                                                // 40
		email: "Email cím",                                                                                                 // 41
		reset: "Jelszó visszaállítása",                                                                                     // 42
		invalidEmail: "Érvénytelen email cím"                                                                               // 43
	},                                                                                                                   // 44
	loginButtonsBackToLoginLink: {                                                                                       // 45
		back: "Mégsem"                                                                                                      // 46
	},                                                                                                                   // 47
	loginButtonsChangePassword: {                                                                                        // 48
		submit: "Módosít",                                                                                                  // 49
		cancel: "Mégsem"                                                                                                    // 50
	},                                                                                                                   // 51
	loginButtonsLoggedOutSingleLoginButton: {                                                                            // 52
		signInWith: "Bejelentkezés: ",                                                                                      // 53
		configure: "Beállítás",                                                                                             // 54
	},                                                                                                                   // 55
	loginButtonsLoggedInSingleLogoutButton: {                                                                            // 56
		signOut: "Kijelentkezés"                                                                                            // 57
	},                                                                                                                   // 58
	loginButtonsLoggedOut: {                                                                                             // 59
		noLoginServices: "Nincs bejelentkezési szolgáltatás beállítva"                                                      // 60
	},                                                                                                                   // 61
	loginFields: {                                                                                                       // 62
		usernameOrEmail: "Felhasználónév vagy Email cím",                                                                   // 63
		username: "Felhasználónév",                                                                                         // 64
		email: "Email cím",                                                                                                 // 65
		password: "Jelszó"                                                                                                  // 66
	},                                                                                                                   // 67
	signupFields: {                                                                                                      // 68
		username: "Felhasználónév",                                                                                         // 69
		email: "Email cím",                                                                                                 // 70
		emailOpt: "Email cím (nem kötelező)",                                                                               // 71
		password: "Jelszó",                                                                                                 // 72
		passwordAgain: "Jelszó (újra)"                                                                                      // 73
	},                                                                                                                   // 74
	changePasswordFields: {                                                                                              // 75
		currentPassword: "Jelenlegi jelszó",                                                                                // 76
		newPassword: "Új jelszó",                                                                                           // 77
		newPasswordAgain: "Új jelszó (újra)"                                                                                // 78
	},                                                                                                                   // 79
	infoMessages : {                                                                                                     // 80
		emailSent: "Email elküldve",                                                                                        // 81
		passwordChanged: "Jelszó megváltoztatva"                                                                            // 82
	},                                                                                                                   // 83
	errorMessages: {                                                                                                     // 84
		genericTitle: "Hiba történt",                                                                                       // 85
		userNotFound: "Nem létező felhasználó",                                                                             // 86
		invalidEmail: "Érvénytelen email cím",                                                                              // 87
		incorrectPassword: "Hibás jelszó",                                                                                  // 88
		usernameTooShort: "A felhasználónévnek legalább 3 karakter hosszúnak kell lennie",                                  // 89
		passwordTooShort: "A jelszónak legalább 6 karakter hosszúnak kell lennie",                                          // 90
		passwordsDontMatch: "A jelszavak nem egyeznek",                                                                     // 91
		newPasswordSameAsOld: "Az új jelszónak el kell térnie a régi jelszótól",                                            // 92
		signupsForbidden: "A regisztráció le van tiltva"                                                                    // 93
	}                                                                                                                    // 94
});                                                                                                                   // 95
                                                                                                                      // 96
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/i18n.js                                                                       //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
i18n.setDefaultLanguage('en')                                                                                         // 1
                                                                                                                      // 2
accountsUIBootstrap3 = {                                                                                              // 3
	setLanguage: function (lang) {                                                                                       // 4
		return i18n.setLanguage(lang)                                                                                       // 5
	},                                                                                                                   // 6
	getLanguage: function () {                                                                                           // 7
		return i18n.getLanguage()                                                                                           // 8
	},                                                                                                                   // 9
	map: function (lang, obj) {                                                                                          // 10
		return i18n.map(lang, obj)                                                                                          // 11
	}                                                                                                                    // 12
}                                                                                                                     // 13
                                                                                                                      // 14
                                                                                                                      // 15
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/template.login_buttons.js                                                     //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
                                                                                                                      // 1
Template.__checkName("_loginButtons");                                                                                // 2
Template["_loginButtons"] = new Template("Template._loginButtons", (function() {                                      // 3
  var view = this;                                                                                                    // 4
  return Blaze.If(function() {                                                                                        // 5
    return Spacebars.call(view.lookup("currentUser"));                                                                // 6
  }, function() {                                                                                                     // 7
    return [ "\n		", Blaze.Unless(function() {                                                                        // 8
      return Spacebars.call(view.lookup("loggingIn"));                                                                // 9
    }, function() {                                                                                                   // 10
      return [ "\n			", Spacebars.include(view.lookupTemplate("_loginButtonsLoggedIn")), "\n		" ];                    // 11
    }), "\n	" ];                                                                                                      // 12
  }, function() {                                                                                                     // 13
    return [ "\n		", Spacebars.include(view.lookupTemplate("_loginButtonsLoggedOut")), "\n	" ];                       // 14
  });                                                                                                                 // 15
}));                                                                                                                  // 16
                                                                                                                      // 17
Template.__checkName("_loginButtonsLoggedIn");                                                                        // 18
Template["_loginButtonsLoggedIn"] = new Template("Template._loginButtonsLoggedIn", (function() {                      // 19
  var view = this;                                                                                                    // 20
  return Spacebars.include(view.lookupTemplate("_loginButtonsLoggedInDropdown"));                                     // 21
}));                                                                                                                  // 22
                                                                                                                      // 23
Template.__checkName("_loginButtonsLoggedOut");                                                                       // 24
Template["_loginButtonsLoggedOut"] = new Template("Template._loginButtonsLoggedOut", (function() {                    // 25
  var view = this;                                                                                                    // 26
  return Blaze.If(function() {                                                                                        // 27
    return Spacebars.call(view.lookup("services"));                                                                   // 28
  }, function() {                                                                                                     // 29
    return [ " \n		", Blaze.If(function() {                                                                           // 30
      return Spacebars.call(view.lookup("configurationLoaded"));                                                      // 31
    }, function() {                                                                                                   // 32
      return [ "\n			", Blaze.If(function() {                                                                         // 33
        return Spacebars.call(view.lookup("dropdown"));                                                               // 34
      }, function() {                                                                                                 // 35
        return [ " \n				", Spacebars.include(view.lookupTemplate("_loginButtonsLoggedOutDropdown")), "\n			" ];      // 36
      }, function() {                                                                                                 // 37
        return [ "\n				", Spacebars.With(function() {                                                                // 38
          return Spacebars.call(view.lookup("singleService"));                                                        // 39
        }, function() {                                                                                               // 40
          return [ " \n					", Blaze.Unless(function() {                                                              // 41
            return Spacebars.call(view.lookup("logginIn"));                                                           // 42
          }, function() {                                                                                             // 43
            return [ "\n						", HTML.DIV({                                                                           // 44
              "class": "navbar-form"                                                                                  // 45
            }, "\n							", Spacebars.include(view.lookupTemplate("_loginButtonsLoggedOutSingleLoginButton")), "\n						"), "\n					" ];
          }), "\n				" ];                                                                                             // 47
        }), "\n			" ];                                                                                                // 48
      }), "\n		" ];                                                                                                   // 49
    }), "\n	" ];                                                                                                      // 50
  }, function() {                                                                                                     // 51
    return [ "\n		", HTML.DIV({                                                                                       // 52
      "class": "no-services"                                                                                          // 53
    }, Blaze.View("lookup:i18n", function() {                                                                         // 54
      return Spacebars.mustache(view.lookup("i18n"), "loginButtonsLoggedOut.noLoginServices");                        // 55
    })), "\n	" ];                                                                                                     // 56
  });                                                                                                                 // 57
}));                                                                                                                  // 58
                                                                                                                      // 59
Template.__checkName("_loginButtonsMessages");                                                                        // 60
Template["_loginButtonsMessages"] = new Template("Template._loginButtonsMessages", (function() {                      // 61
  var view = this;                                                                                                    // 62
  return [ Blaze.If(function() {                                                                                      // 63
    return Spacebars.call(view.lookup("errorMessage"));                                                               // 64
  }, function() {                                                                                                     // 65
    return [ "\n		", HTML.DIV({                                                                                       // 66
      "class": "alert alert-danger"                                                                                   // 67
    }, Blaze.View("lookup:errorMessage", function() {                                                                 // 68
      return Spacebars.mustache(view.lookup("errorMessage"));                                                         // 69
    })), "\n	" ];                                                                                                     // 70
  }), "\n	", Blaze.If(function() {                                                                                    // 71
    return Spacebars.call(view.lookup("infoMessage"));                                                                // 72
  }, function() {                                                                                                     // 73
    return [ "\n		", HTML.DIV({                                                                                       // 74
      "class": "alert alert-success no-margin"                                                                        // 75
    }, Blaze.View("lookup:infoMessage", function() {                                                                  // 76
      return Spacebars.mustache(view.lookup("infoMessage"));                                                          // 77
    })), "\n	" ];                                                                                                     // 78
  }) ];                                                                                                               // 79
}));                                                                                                                  // 80
                                                                                                                      // 81
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/template.login_buttons_single.js                                              //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
                                                                                                                      // 1
Template.__checkName("_loginButtonsLoggedOutSingleLoginButton");                                                      // 2
Template["_loginButtonsLoggedOutSingleLoginButton"] = new Template("Template._loginButtonsLoggedOutSingleLoginButton", (function() {
  var view = this;                                                                                                    // 4
  return Blaze.If(function() {                                                                                        // 5
    return Spacebars.call(view.lookup("configured"));                                                                 // 6
  }, function() {                                                                                                     // 7
    return [ "\n		", HTML.BUTTON({                                                                                    // 8
      "class": function() {                                                                                           // 9
        return [ "login-button btn btn-block btn-", Spacebars.mustache(view.lookup("capitalizedName")) ];             // 10
      }                                                                                                               // 11
    }, Blaze.View("lookup:i18n", function() {                                                                         // 12
      return Spacebars.mustache(view.lookup("i18n"), "loginButtonsLoggedOutSingleLoginButton.signInWith");            // 13
    }), " ", Blaze.View("lookup:capitalizedName", function() {                                                        // 14
      return Spacebars.mustache(view.lookup("capitalizedName"));                                                      // 15
    })), "\n	" ];                                                                                                     // 16
  }, function() {                                                                                                     // 17
    return [ "\n		", HTML.BUTTON({                                                                                    // 18
      "class": "login-button btn btn-block configure-button btn-danger"                                               // 19
    }, Blaze.View("lookup:i18n", function() {                                                                         // 20
      return Spacebars.mustache(view.lookup("i18n"), "loginButtonsLoggedOutSingleLoginButton.configure");             // 21
    }), " ", Blaze.View("lookup:capitalizedName", function() {                                                        // 22
      return Spacebars.mustache(view.lookup("capitalizedName"));                                                      // 23
    })), "\n	" ];                                                                                                     // 24
  });                                                                                                                 // 25
}));                                                                                                                  // 26
                                                                                                                      // 27
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/template.login_buttons_dropdown.js                                            //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
                                                                                                                      // 1
Template.__checkName("_loginButtonsLoggedInDropdown");                                                                // 2
Template["_loginButtonsLoggedInDropdown"] = new Template("Template._loginButtonsLoggedInDropdown", (function() {      // 3
  var view = this;                                                                                                    // 4
  return HTML.LI({                                                                                                    // 5
    id: "login-dropdown-list",                                                                                        // 6
    "class": "dropdown"                                                                                               // 7
  }, "\n		", HTML.A({                                                                                                 // 8
    "class": "dropdown-toggle",                                                                                       // 9
    "data-toggle": "dropdown"                                                                                         // 10
  }, "\n			", Blaze.View("lookup:displayName", function() {                                                           // 11
    return Spacebars.mustache(view.lookup("displayName"));                                                            // 12
  }), "\n			", Spacebars.With(function() {                                                                            // 13
    return Spacebars.call(view.lookup("user_profile_picture"));                                                       // 14
  }, function() {                                                                                                     // 15
    return [ "\n				", HTML.IMG({                                                                                     // 16
      src: function() {                                                                                               // 17
        return Spacebars.mustache(view.lookup("."));                                                                  // 18
      },                                                                                                              // 19
      width: "30px",                                                                                                  // 20
      height: "30px",                                                                                                 // 21
      "class": "img-circle",                                                                                          // 22
      alt: "#"                                                                                                        // 23
    }), "\n			" ];                                                                                                    // 24
  }), "\n			", HTML.Raw('<b class="caret"></b>'), "\n		"), "\n		", HTML.DIV({                                         // 25
    "class": "dropdown-menu"                                                                                          // 26
  }, "\n			", Blaze.If(function() {                                                                                   // 27
    return Spacebars.call(view.lookup("inMessageOnlyFlow"));                                                          // 28
  }, function() {                                                                                                     // 29
    return [ "\n				", Spacebars.include(view.lookupTemplate("_loginButtonsMessages")), "\n			" ];                    // 30
  }, function() {                                                                                                     // 31
    return [ "\n				", Blaze.If(function() {                                                                          // 32
      return Spacebars.call(view.lookup("inChangePasswordFlow"));                                                     // 33
    }, function() {                                                                                                   // 34
      return [ "\n					", Spacebars.include(view.lookupTemplate("_loginButtonsChangePassword")), "\n				" ];          // 35
    }, function() {                                                                                                   // 36
      return [ "\n					", Spacebars.include(view.lookupTemplate("_loginButtonsLoggedInDropdownActions")), "\n				" ];
    }), "\n			" ];                                                                                                    // 38
  }), "\n		"), "\n	");                                                                                                // 39
}));                                                                                                                  // 40
                                                                                                                      // 41
Template.__checkName("_loginButtonsLoggedInDropdownActions");                                                         // 42
Template["_loginButtonsLoggedInDropdownActions"] = new Template("Template._loginButtonsLoggedInDropdownActions", (function() {
  var view = this;                                                                                                    // 44
  return [ Blaze.If(function() {                                                                                      // 45
    return Spacebars.call(view.lookup("additionalLoggedInDropdownActions"));                                          // 46
  }, function() {                                                                                                     // 47
    return [ "\n		", Spacebars.include(view.lookupTemplate("_loginButtonsAdditionalLoggedInDropdownActions")), "\n	" ];
  }), "\n\n	", Blaze.If(function() {                                                                                  // 49
    return Spacebars.call(view.lookup("allowChangingPassword"));                                                      // 50
  }, function() {                                                                                                     // 51
    return [ "\n		", HTML.BUTTON({                                                                                    // 52
      "class": "btn btn-default btn-block",                                                                           // 53
      id: "login-buttons-open-change-password"                                                                        // 54
    }, Blaze.View("lookup:i18n", function() {                                                                         // 55
      return Spacebars.mustache(view.lookup("i18n"), "loginButtonsLoggedInDropdownActions.password");                 // 56
    })), "\n	" ];                                                                                                     // 57
  }), "\n\n	", HTML.BUTTON({                                                                                          // 58
    "class": "btn btn-block btn-primary",                                                                             // 59
    id: "login-buttons-logout"                                                                                        // 60
  }, Blaze.View("lookup:i18n", function() {                                                                           // 61
    return Spacebars.mustache(view.lookup("i18n"), "loginButtonsLoggedInDropdownActions.signOut");                    // 62
  })) ];                                                                                                              // 63
}));                                                                                                                  // 64
                                                                                                                      // 65
Template.__checkName("_loginButtonsLoggedOutDropdown");                                                               // 66
Template["_loginButtonsLoggedOutDropdown"] = new Template("Template._loginButtonsLoggedOutDropdown", (function() {    // 67
  var view = this;                                                                                                    // 68
  return HTML.LI({                                                                                                    // 69
    id: "login-dropdown-list",                                                                                        // 70
    "class": "dropdown"                                                                                               // 71
  }, "\n		", HTML.A({                                                                                                 // 72
    "class": "dropdown-toggle",                                                                                       // 73
    "data-toggle": "dropdown"                                                                                         // 74
  }, Blaze.View("lookup:i18n", function() {                                                                           // 75
    return Spacebars.mustache(view.lookup("i18n"), "loginButtonsLoggedOutDropdown.signIn");                           // 76
  }), Blaze.Unless(function() {                                                                                       // 77
    return Spacebars.call(view.lookup("forbidClientAccountCreation"));                                                // 78
  }, function() {                                                                                                     // 79
    return [ " / ", Blaze.View("lookup:i18n", function() {                                                            // 80
      return Spacebars.mustache(view.lookup("i18n"), "loginButtonsLoggedOutDropdown.up");                             // 81
    }) ];                                                                                                             // 82
  }), " ", HTML.Raw('<b class="caret"></b>')), "\n		", HTML.DIV({                                                     // 83
    "class": "dropdown-menu"                                                                                          // 84
  }, "\n			", Spacebars.include(view.lookupTemplate("_loginButtonsLoggedOutAllServices")), "\n		"), "\n	");           // 85
}));                                                                                                                  // 86
                                                                                                                      // 87
Template.__checkName("_loginButtonsLoggedOutAllServices");                                                            // 88
Template["_loginButtonsLoggedOutAllServices"] = new Template("Template._loginButtonsLoggedOutAllServices", (function() {
  var view = this;                                                                                                    // 90
  return Blaze.Each(function() {                                                                                      // 91
    return Spacebars.call(view.lookup("services"));                                                                   // 92
  }, function() {                                                                                                     // 93
    return [ "\n	", Blaze.Unless(function() {                                                                         // 94
      return Spacebars.call(view.lookup("hasPasswordService"));                                                       // 95
    }, function() {                                                                                                   // 96
      return [ "\n		", Spacebars.include(view.lookupTemplate("_loginButtonsMessages")), "\n	" ];                      // 97
    }), "\n		", Blaze.If(function() {                                                                                 // 98
      return Spacebars.call(view.lookup("isPasswordService"));                                                        // 99
    }, function() {                                                                                                   // 100
      return [ "\n			", Blaze.If(function() {                                                                         // 101
        return Spacebars.call(view.lookup("hasOtherServices"));                                                       // 102
      }, function() {                                                                                                 // 103
        return [ " \n				", Spacebars.include(view.lookupTemplate("_loginButtonsLoggedOutPasswordServiceSeparator")), "\n			" ];
      }), "\n			", Spacebars.include(view.lookupTemplate("_loginButtonsLoggedOutPasswordService")), "\n		" ];         // 105
    }, function() {                                                                                                   // 106
      return [ "\n			", Spacebars.include(view.lookupTemplate("_loginButtonsLoggedOutSingleLoginButton")), "\n		" ];  // 107
    }), "\n	" ];                                                                                                      // 108
  });                                                                                                                 // 109
}));                                                                                                                  // 110
                                                                                                                      // 111
Template.__checkName("_loginButtonsLoggedOutPasswordServiceSeparator");                                               // 112
Template["_loginButtonsLoggedOutPasswordServiceSeparator"] = new Template("Template._loginButtonsLoggedOutPasswordServiceSeparator", (function() {
  var view = this;                                                                                                    // 114
  return HTML.DIV({                                                                                                   // 115
    "class": "or"                                                                                                     // 116
  }, HTML.Raw('\n		<span class="hline">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>\n		'), HTML.SPAN({
    "class": "or-text"                                                                                                // 118
  }, Blaze.View("lookup:i18n", function() {                                                                           // 119
    return Spacebars.mustache(view.lookup("i18n"), "loginButtonsLoggedOutPasswordServiceSeparator.or");               // 120
  })), HTML.Raw('\n		<span class="hline">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>\n	'));   // 121
}));                                                                                                                  // 122
                                                                                                                      // 123
Template.__checkName("_loginButtonsLoggedOutPasswordService");                                                        // 124
Template["_loginButtonsLoggedOutPasswordService"] = new Template("Template._loginButtonsLoggedOutPasswordService", (function() {
  var view = this;                                                                                                    // 126
  return Blaze.If(function() {                                                                                        // 127
    return Spacebars.call(view.lookup("inForgotPasswordFlow"));                                                       // 128
  }, function() {                                                                                                     // 129
    return [ "\n		", Spacebars.include(view.lookupTemplate("_forgotPasswordForm")), "\n	" ];                          // 130
  }, function() {                                                                                                     // 131
    return [ "\n		", Spacebars.include(view.lookupTemplate("_loginButtonsMessages")), "\n		", Blaze.Each(function() {
      return Spacebars.call(view.lookup("fields"));                                                                   // 133
    }, function() {                                                                                                   // 134
      return [ "\n			", Spacebars.include(view.lookupTemplate("_loginButtonsFormField")), "\n		" ];                   // 135
    }), "\n		", HTML.BUTTON({                                                                                         // 136
      "class": "btn btn-primary col-xs-12 col-sm-12",                                                                 // 137
      id: "login-buttons-password",                                                                                   // 138
      type: "button"                                                                                                  // 139
    }, "\n			", Blaze.If(function() {                                                                                 // 140
      return Spacebars.call(view.lookup("inSignupFlow"));                                                             // 141
    }, function() {                                                                                                   // 142
      return [ "\n				", Blaze.View("lookup:i18n", function() {                                                       // 143
        return Spacebars.mustache(view.lookup("i18n"), "loginButtonsLoggedOutPasswordService.create");                // 144
      }), "\n			" ];                                                                                                  // 145
    }, function() {                                                                                                   // 146
      return [ "\n				", Blaze.View("lookup:i18n", function() {                                                       // 147
        return Spacebars.mustache(view.lookup("i18n"), "loginButtonsLoggedOutPasswordService.signIn");                // 148
      }), "\n			" ];                                                                                                  // 149
    }), "\n		"), "\n		", Blaze.If(function() {                                                                        // 150
      return Spacebars.call(view.lookup("inLoginFlow"));                                                              // 151
    }, function() {                                                                                                   // 152
      return [ "\n			", HTML.DIV({                                                                                    // 153
        id: "login-other-options"                                                                                     // 154
      }, "\n			", Blaze.If(function() {                                                                               // 155
        return Spacebars.call(view.lookup("showForgotPasswordLink"));                                                 // 156
      }, function() {                                                                                                 // 157
        return [ "\n				", HTML.A({                                                                                   // 158
          id: "forgot-password-link",                                                                                 // 159
          "class": "pull-left"                                                                                        // 160
        }, Blaze.View("lookup:i18n", function() {                                                                     // 161
          return Spacebars.mustache(view.lookup("i18n"), "loginButtonsLoggedOutPasswordService.forgot");              // 162
        })), "\n			" ];                                                                                               // 163
      }), "\n			", Blaze.If(function() {                                                                              // 164
        return Spacebars.call(view.lookup("showCreateAccountLink"));                                                  // 165
      }, function() {                                                                                                 // 166
        return [ "\n				", HTML.A({                                                                                   // 167
          id: "signup-link",                                                                                          // 168
          "class": "pull-right"                                                                                       // 169
        }, Blaze.View("lookup:i18n", function() {                                                                     // 170
          return Spacebars.mustache(view.lookup("i18n"), "loginButtonsLoggedOutPasswordService.createAcc");           // 171
        })), "\n			" ];                                                                                               // 172
      }), "\n			"), "\n		" ];                                                                                         // 173
    }), "\n		", Blaze.If(function() {                                                                                 // 174
      return Spacebars.call(view.lookup("inSignupFlow"));                                                             // 175
    }, function() {                                                                                                   // 176
      return [ "\n			", Spacebars.include(view.lookupTemplate("_loginButtonsBackToLoginLink")), "\n		" ];             // 177
    }), "\n	" ];                                                                                                      // 178
  });                                                                                                                 // 179
}));                                                                                                                  // 180
                                                                                                                      // 181
Template.__checkName("_forgotPasswordForm");                                                                          // 182
Template["_forgotPasswordForm"] = new Template("Template._forgotPasswordForm", (function() {                          // 183
  var view = this;                                                                                                    // 184
  return HTML.DIV({                                                                                                   // 185
    "class": "login-form"                                                                                             // 186
  }, "\n		", Spacebars.include(view.lookupTemplate("_loginButtonsMessages")), "\n		", HTML.DIV({                      // 187
    id: "forgot-password-email-label-and-input"                                                                       // 188
  }, " \n			", HTML.INPUT({                                                                                           // 189
    id: "forgot-password-email",                                                                                      // 190
    type: "email",                                                                                                    // 191
    placeholder: function() {                                                                                         // 192
      return Spacebars.mustache(view.lookup("i18n"), "forgotPasswordForm.email");                                     // 193
    },                                                                                                                // 194
    "class": "form-control"                                                                                           // 195
  }), "\n		"), "\n		", HTML.BUTTON({                                                                                  // 196
    "class": "btn btn-primary login-button-form-submit col-xs-12 col-sm-12",                                          // 197
    id: "login-buttons-forgot-password"                                                                               // 198
  }, Blaze.View("lookup:i18n", function() {                                                                           // 199
    return Spacebars.mustache(view.lookup("i18n"), "forgotPasswordForm.reset");                                       // 200
  })), "\n		", Spacebars.include(view.lookupTemplate("_loginButtonsBackToLoginLink")), "\n	");                        // 201
}));                                                                                                                  // 202
                                                                                                                      // 203
Template.__checkName("_loginButtonsBackToLoginLink");                                                                 // 204
Template["_loginButtonsBackToLoginLink"] = new Template("Template._loginButtonsBackToLoginLink", (function() {        // 205
  var view = this;                                                                                                    // 206
  return HTML.BUTTON({                                                                                                // 207
    id: "back-to-login-link",                                                                                         // 208
    "class": "btn btn-default col-xs-12 col-sm-12"                                                                    // 209
  }, Blaze.View("lookup:i18n", function() {                                                                           // 210
    return Spacebars.mustache(view.lookup("i18n"), "loginButtonsBackToLoginLink.back");                               // 211
  }));                                                                                                                // 212
}));                                                                                                                  // 213
                                                                                                                      // 214
Template.__checkName("_loginButtonsFormField");                                                                       // 215
Template["_loginButtonsFormField"] = new Template("Template._loginButtonsFormField", (function() {                    // 216
  var view = this;                                                                                                    // 217
  return Blaze.If(function() {                                                                                        // 218
    return Spacebars.call(view.lookup("visible"));                                                                    // 219
  }, function() {                                                                                                     // 220
    return [ "\n		", HTML.Comment(" TODO: Implement more input types "), "\n		", Blaze.If(function() {                // 221
      return Spacebars.dataMustache(view.lookup("equals"), view.lookup("inputType"), "checkbox");                     // 222
    }, function() {                                                                                                   // 223
      return [ "\n			", HTML.DIV({                                                                                    // 224
        "class": "checkbox"                                                                                           // 225
      }, "\n				", HTML.LABEL(HTML.INPUT({                                                                            // 226
        type: "checkbox",                                                                                             // 227
        id: function() {                                                                                              // 228
          return [ "login-", Spacebars.mustache(view.lookup("fieldName")) ];                                          // 229
        },                                                                                                            // 230
        name: function() {                                                                                            // 231
          return [ "login-", Spacebars.mustache(view.lookup("fieldName")) ];                                          // 232
        },                                                                                                            // 233
        value: "true"                                                                                                 // 234
      }), "\n				", Blaze.View("lookup:fieldLabel", function() {                                                      // 235
        return Spacebars.makeRaw(Spacebars.mustache(view.lookup("fieldLabel")));                                      // 236
      })), "\n			"), "\n		" ];                                                                                        // 237
    }), "\n\n		", Blaze.If(function() {                                                                               // 238
      return Spacebars.dataMustache(view.lookup("equals"), view.lookup("inputType"), "select");                       // 239
    }, function() {                                                                                                   // 240
      return [ "\n			", HTML.DIV({                                                                                    // 241
        "class": "select-dropdown"                                                                                    // 242
      }, "\n			", Blaze.If(function() {                                                                               // 243
        return Spacebars.call(view.lookup("showFieldLabel"));                                                         // 244
      }, function() {                                                                                                 // 245
        return [ "\n				", HTML.LABEL(Blaze.View("lookup:fieldLabel", function() {                                    // 246
          return Spacebars.mustache(view.lookup("fieldLabel"));                                                       // 247
        })), HTML.BR(), "\n			" ];                                                                                    // 248
      }), "\n			", HTML.SELECT({                                                                                      // 249
        id: function() {                                                                                              // 250
          return [ "login-", Spacebars.mustache(view.lookup("fieldName")) ];                                          // 251
        }                                                                                                             // 252
      }, "\n				", Blaze.If(function() {                                                                              // 253
        return Spacebars.call(view.lookup("empty"));                                                                  // 254
      }, function() {                                                                                                 // 255
        return [ "\n					", HTML.OPTION({                                                                             // 256
          value: ""                                                                                                   // 257
        }, Blaze.View("lookup:empty", function() {                                                                    // 258
          return Spacebars.mustache(view.lookup("empty"));                                                            // 259
        })), "\n				" ];                                                                                              // 260
      }), "\n				", Blaze.Each(function() {                                                                           // 261
        return Spacebars.call(view.lookup("data"));                                                                   // 262
      }, function() {                                                                                                 // 263
        return [ "\n					", HTML.OPTION({                                                                             // 264
          value: function() {                                                                                         // 265
            return Spacebars.mustache(view.lookup("value"));                                                          // 266
          }                                                                                                           // 267
        }, Blaze.View("lookup:label", function() {                                                                    // 268
          return Spacebars.mustache(view.lookup("label"));                                                            // 269
        })), "\n				" ];                                                                                              // 270
      }), "\n			"), "\n			"), "\n		" ];                                                                               // 271
    }), "\n\n		", Blaze.If(function() {                                                                               // 272
      return Spacebars.dataMustache(view.lookup("equals"), view.lookup("inputType"), "radio");                        // 273
    }, function() {                                                                                                   // 274
      return [ "\n			", HTML.DIV({                                                                                    // 275
        "class": "radio"                                                                                              // 276
      }, "\n				", Blaze.If(function() {                                                                              // 277
        return Spacebars.call(view.lookup("showFieldLabel"));                                                         // 278
      }, function() {                                                                                                 // 279
        return [ "\n				", HTML.LABEL(Blaze.View("lookup:fieldLabel", function() {                                    // 280
          return Spacebars.mustache(view.lookup("fieldLabel"));                                                       // 281
        })), HTML.BR(), "\n				" ];                                                                                   // 282
      }), "\n				", Blaze.Each(function() {                                                                           // 283
        return Spacebars.call(view.lookup("data"));                                                                   // 284
      }, function() {                                                                                                 // 285
        return [ "\n					", HTML.LABEL(HTML.INPUT(HTML.Attrs({                                                        // 286
          type: "radio",                                                                                              // 287
          id: function() {                                                                                            // 288
            return [ "login-", Spacebars.mustache(Spacebars.dot(view.lookup(".."), "fieldName")), "-", Spacebars.mustache(view.lookup("id")) ];
          },                                                                                                          // 290
          name: function() {                                                                                          // 291
            return [ "login-", Spacebars.mustache(Spacebars.dot(view.lookup(".."), "fieldName")) ];                   // 292
          },                                                                                                          // 293
          value: function() {                                                                                         // 294
            return Spacebars.mustache(view.lookup("value"));                                                          // 295
          }                                                                                                           // 296
        }, function() {                                                                                               // 297
          return Spacebars.attrMustache(view.lookup("checked"));                                                      // 298
        })), " ", Blaze.View("lookup:label", function() {                                                             // 299
          return Spacebars.mustache(view.lookup("label"));                                                            // 300
        })), "\n					", Blaze.If(function() {                                                                         // 301
          return Spacebars.dataMustache(view.lookup("equals"), Spacebars.dot(view.lookup(".."), "radioLayout"), "vertical");
        }, function() {                                                                                               // 303
          return [ "\n						", HTML.BR(), "\n					" ];                                                                // 304
        }), "\n				" ];                                                                                               // 305
      }), "\n			"), "\n		" ];                                                                                         // 306
    }), "\n\n		", Blaze.If(function() {                                                                               // 307
      return Spacebars.call(view.lookup("inputTextual"));                                                             // 308
    }, function() {                                                                                                   // 309
      return [ "\n			", HTML.INPUT({                                                                                  // 310
        id: function() {                                                                                              // 311
          return [ "login-", Spacebars.mustache(view.lookup("fieldName")) ];                                          // 312
        },                                                                                                            // 313
        type: function() {                                                                                            // 314
          return Spacebars.mustache(view.lookup("inputType"));                                                        // 315
        },                                                                                                            // 316
        placeholder: function() {                                                                                     // 317
          return Spacebars.mustache(view.lookup("fieldLabel"));                                                       // 318
        },                                                                                                            // 319
        "class": "form-control"                                                                                       // 320
      }), "\n		" ];                                                                                                   // 321
    }), "\n	" ];                                                                                                      // 322
  });                                                                                                                 // 323
}));                                                                                                                  // 324
                                                                                                                      // 325
Template.__checkName("_loginButtonsChangePassword");                                                                  // 326
Template["_loginButtonsChangePassword"] = new Template("Template._loginButtonsChangePassword", (function() {          // 327
  var view = this;                                                                                                    // 328
  return [ Spacebars.include(view.lookupTemplate("_loginButtonsMessages")), "\n	", Blaze.Each(function() {            // 329
    return Spacebars.call(view.lookup("fields"));                                                                     // 330
  }, function() {                                                                                                     // 331
    return [ "\n		", Spacebars.include(view.lookupTemplate("_loginButtonsFormField")), "\n	" ];                       // 332
  }), "\n	", HTML.BUTTON({                                                                                            // 333
    "class": "btn btn-block btn-primary",                                                                             // 334
    id: "login-buttons-do-change-password"                                                                            // 335
  }, Blaze.View("lookup:i18n", function() {                                                                           // 336
    return Spacebars.mustache(view.lookup("i18n"), "loginButtonsChangePassword.submit");                              // 337
  })), "\n	", HTML.BUTTON({                                                                                           // 338
    "class": "btn btn-block btn-default",                                                                             // 339
    id: "login-buttons-cancel-change-password"                                                                        // 340
  }, Blaze.View("lookup:i18n", function() {                                                                           // 341
    return Spacebars.mustache(view.lookup("i18n"), "loginButtonsChangePassword.cancel");                              // 342
  })) ];                                                                                                              // 343
}));                                                                                                                  // 344
                                                                                                                      // 345
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/template.login_buttons_dialogs.js                                             //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
                                                                                                                      // 1
Template.body.addContent((function() {                                                                                // 2
  var view = this;                                                                                                    // 3
  return [ Spacebars.include(view.lookupTemplate("_resetPasswordDialog")), "\n	", Spacebars.include(view.lookupTemplate("_enrollAccountDialog")), "\n	", Spacebars.include(view.lookupTemplate("_justVerifiedEmailDialog")), "\n	", Spacebars.include(view.lookupTemplate("_configureLoginServiceDialog")), "\n	", Spacebars.include(view.lookupTemplate("_loginButtonsMessagesDialog")) ];
}));                                                                                                                  // 5
Meteor.startup(Template.body.renderToDocument);                                                                       // 6
                                                                                                                      // 7
Template.__checkName("_resetPasswordDialog");                                                                         // 8
Template["_resetPasswordDialog"] = new Template("Template._resetPasswordDialog", (function() {                        // 9
  var view = this;                                                                                                    // 10
  return [ Blaze.If(function() {                                                                                      // 11
    return Spacebars.call(view.lookup("inResetPasswordFlow"));                                                        // 12
  }, function() {                                                                                                     // 13
    return [ "\n		", HTML.DIV({                                                                                       // 14
      "class": "modal",                                                                                               // 15
      id: "login-buttons-reset-password-modal"                                                                        // 16
    }, "\n			", HTML.DIV({                                                                                            // 17
      "class": "modal-dialog"                                                                                         // 18
    }, "\n				", HTML.DIV({                                                                                           // 19
      "class": "modal-content"                                                                                        // 20
    }, "\n					", HTML.DIV({                                                                                          // 21
      "class": "modal-header"                                                                                         // 22
    }, "\n						", HTML.BUTTON({                                                                                      // 23
      type: "button",                                                                                                 // 24
      "class": "close",                                                                                               // 25
      "data-dismiss": "modal",                                                                                        // 26
      "aria-hidden": "true"                                                                                           // 27
    }, HTML.CharRef({                                                                                                 // 28
      html: "&times;",                                                                                                // 29
      str: "×"                                                                                                        // 30
    })), "\n						", HTML.H4({                                                                                        // 31
      "class": "modal-title"                                                                                          // 32
    }, Blaze.View("lookup:i18n", function() {                                                                         // 33
      return Spacebars.mustache(view.lookup("i18n"), "resetPasswordDialog.title");                                    // 34
    })), "\n					"), "\n					", HTML.DIV({                                                                            // 35
      "class": "modal-body"                                                                                           // 36
    }, "\n						", Spacebars.include(view.lookupTemplate("_loginButtonsMessages")), "\n						", HTML.INPUT({          // 37
      id: "reset-password-new-password",                                                                              // 38
      "class": "form-control",                                                                                        // 39
      type: "password",                                                                                               // 40
      placeholder: function() {                                                                                       // 41
        return Spacebars.mustache(view.lookup("i18n"), "resetPasswordDialog.newPassword");                            // 42
      }                                                                                                               // 43
    }), HTML.BR(), "\n						", HTML.INPUT({                                                                           // 44
      id: "reset-password-new-password-again",                                                                        // 45
      "class": "form-control",                                                                                        // 46
      type: "password",                                                                                               // 47
      placeholder: function() {                                                                                       // 48
        return Spacebars.mustache(view.lookup("i18n"), "resetPasswordDialog.newPasswordAgain");                       // 49
      }                                                                                                               // 50
    }), HTML.BR(), "\n					"), "\n					", HTML.DIV({                                                                  // 51
      "class": "modal-footer"                                                                                         // 52
    }, "\n						", HTML.A({                                                                                           // 53
      "class": "btn btn-default",                                                                                     // 54
      id: "login-buttons-cancel-reset-password"                                                                       // 55
    }, Blaze.View("lookup:i18n", function() {                                                                         // 56
      return Spacebars.mustache(view.lookup("i18n"), "resetPasswordDialog.cancel");                                   // 57
    })), "\n						", HTML.BUTTON({                                                                                    // 58
      "class": "btn btn-primary",                                                                                     // 59
      id: "login-buttons-reset-password-button"                                                                       // 60
    }, "\n							", Blaze.View("lookup:i18n", function() {                                                            // 61
      return Spacebars.mustache(view.lookup("i18n"), "resetPasswordDialog.submit");                                   // 62
    }), "\n						"), "\n					"), "\n				"), HTML.Comment(" /.modal-content "), "\n			"), HTML.Comment(" /.modal-dalog "), "\n		"), HTML.Comment(" /.modal "), "\n	" ];
  }), "\n\n	", HTML.DIV({                                                                                             // 64
    "class": "modal",                                                                                                 // 65
    id: "login-buttons-reset-password-modal-success"                                                                  // 66
  }, "\n		", HTML.DIV({                                                                                               // 67
    "class": "modal-dialog"                                                                                           // 68
  }, "\n			", HTML.DIV({                                                                                              // 69
    "class": "modal-content"                                                                                          // 70
  }, "\n				", HTML.DIV({                                                                                             // 71
    "class": "modal-header"                                                                                           // 72
  }, "\n					", HTML.Raw('<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>'), "\n					", HTML.H4({
    "class": "modal-title"                                                                                            // 74
  }, Blaze.View("lookup:i18n", function() {                                                                           // 75
    return Spacebars.mustache(view.lookup("i18n"), "resetPasswordDialog.title");                                      // 76
  })), "\n				"), "\n				", HTML.DIV({                                                                                // 77
    "class": "modal-body"                                                                                             // 78
  }, "\n					", Spacebars.include(view.lookupTemplate("_loginButtonsMessages")), "\n				"), "\n				", HTML.DIV({      // 79
    "class": "modal-footer"                                                                                           // 80
  }, "\n					", HTML.A({                                                                                              // 81
    "class": "btn btn-default",                                                                                       // 82
    id: "login-buttons-dismiss-reset-password-success"                                                                // 83
  }, Blaze.View("lookup:i18n", function() {                                                                           // 84
    return Spacebars.mustache(view.lookup("i18n"), "loginButtonsMessagesDialog.dismiss");                             // 85
  })), "\n				"), "\n			"), HTML.Raw("<!-- /.modal-content -->"), "\n		"), HTML.Raw("<!-- /.modal-dalog -->"), "\n	"), HTML.Raw("<!-- /.modal -->") ];
}));                                                                                                                  // 87
                                                                                                                      // 88
Template.__checkName("_enrollAccountDialog");                                                                         // 89
Template["_enrollAccountDialog"] = new Template("Template._enrollAccountDialog", (function() {                        // 90
  var view = this;                                                                                                    // 91
  return Blaze.If(function() {                                                                                        // 92
    return Spacebars.call(view.lookup("inEnrollAccountFlow"));                                                        // 93
  }, function() {                                                                                                     // 94
    return [ "\n		", HTML.DIV({                                                                                       // 95
      "class": "modal",                                                                                               // 96
      id: "login-buttons-enroll-account-modal"                                                                        // 97
    }, "\n			", HTML.DIV({                                                                                            // 98
      "class": "modal-dialog"                                                                                         // 99
    }, "\n				", HTML.DIV({                                                                                           // 100
      "class": "modal-content"                                                                                        // 101
    }, "\n					", HTML.DIV({                                                                                          // 102
      "class": "modal-header"                                                                                         // 103
    }, "\n						", HTML.BUTTON({                                                                                      // 104
      type: "button",                                                                                                 // 105
      "class": "close",                                                                                               // 106
      "data-dismiss": "modal",                                                                                        // 107
      "aria-hidden": "true"                                                                                           // 108
    }, HTML.CharRef({                                                                                                 // 109
      html: "&times;",                                                                                                // 110
      str: "×"                                                                                                        // 111
    })), "\n						", HTML.H4({                                                                                        // 112
      "class": "modal-title"                                                                                          // 113
    }, Blaze.View("lookup:i18n", function() {                                                                         // 114
      return Spacebars.mustache(view.lookup("i18n"), "enrollAccountDialog.title");                                    // 115
    })), "\n					"), "\n					", HTML.DIV({                                                                            // 116
      "class": "modal-body"                                                                                           // 117
    }, "\n						", HTML.INPUT({                                                                                       // 118
      id: "enroll-account-password",                                                                                  // 119
      "class": "form-control",                                                                                        // 120
      type: "password",                                                                                               // 121
      placeholder: function() {                                                                                       // 122
        return Spacebars.mustache(view.lookup("i18n"), "enrollAccountDialog.newPassword");                            // 123
      }                                                                                                               // 124
    }), HTML.BR(), "\n												", HTML.INPUT({                                                                     // 125
      id: "enroll-account-password-again",                                                                            // 126
      "class": "form-control",                                                                                        // 127
      type: "password",                                                                                               // 128
      placeholder: function() {                                                                                       // 129
        return Spacebars.mustache(view.lookup("i18n"), "enrollAccountDialog.newPasswordAgain");                       // 130
      }                                                                                                               // 131
    }), HTML.BR(), "\n						", Spacebars.include(view.lookupTemplate("_loginButtonsMessages")), "\n					"), "\n					", HTML.DIV({
      "class": "modal-footer"                                                                                         // 133
    }, "\n						", HTML.A({                                                                                           // 134
      "class": "btn btn-default",                                                                                     // 135
      id: "login-buttons-cancel-enroll-account-button"                                                                // 136
    }, Blaze.View("lookup:i18n", function() {                                                                         // 137
      return Spacebars.mustache(view.lookup("i18n"), "enrollAccountDialog.cancel");                                   // 138
    })), "\n						", HTML.BUTTON({                                                                                    // 139
      "class": "btn btn-primary",                                                                                     // 140
      id: "login-buttons-enroll-account-button"                                                                       // 141
    }, "\n							", Blaze.View("lookup:i18n", function() {                                                            // 142
      return Spacebars.mustache(view.lookup("i18n"), "enrollAccountDialog.submit");                                   // 143
    }), "\n						"), "\n					"), "\n				"), HTML.Comment(" /.modal-content "), "\n			"), HTML.Comment(" /.modal-dalog "), "\n		"), HTML.Comment(" /.modal "), "\n	" ];
  });                                                                                                                 // 145
}));                                                                                                                  // 146
                                                                                                                      // 147
Template.__checkName("_justVerifiedEmailDialog");                                                                     // 148
Template["_justVerifiedEmailDialog"] = new Template("Template._justVerifiedEmailDialog", (function() {                // 149
  var view = this;                                                                                                    // 150
  return Blaze.If(function() {                                                                                        // 151
    return Spacebars.call(view.lookup("visible"));                                                                    // 152
  }, function() {                                                                                                     // 153
    return [ "\n		", HTML.DIV({                                                                                       // 154
      "class": "modal",                                                                                               // 155
      id: "login-buttons-email-address-verified-modal"                                                                // 156
    }, "\n			", HTML.DIV({                                                                                            // 157
      "class": "modal-dialog"                                                                                         // 158
    }, "\n				", HTML.DIV({                                                                                           // 159
      "class": "modal-content"                                                                                        // 160
    }, "\n					", HTML.DIV({                                                                                          // 161
      "class": "modal-body"                                                                                           // 162
    }, "\n						", HTML.H4(HTML.B(Blaze.View("lookup:i18n", function() {                                              // 163
      return Spacebars.mustache(view.lookup("i18n"), "justVerifiedEmailDialog.verified");                             // 164
    }))), "\n					"), "\n					", HTML.DIV({                                                                           // 165
      "class": "modal-footer"                                                                                         // 166
    }, "\n						", HTML.BUTTON({                                                                                      // 167
      "class": "btn btn-primary login-button",                                                                        // 168
      id: "just-verified-dismiss-button",                                                                             // 169
      "data-dismiss": "modal"                                                                                         // 170
    }, Blaze.View("lookup:i18n", function() {                                                                         // 171
      return Spacebars.mustache(view.lookup("i18n"), "justVerifiedEmailDialog.dismiss");                              // 172
    })), "\n					"), "\n				"), "\n			"), "\n		"), "\n	" ];                                                           // 173
  });                                                                                                                 // 174
}));                                                                                                                  // 175
                                                                                                                      // 176
Template.__checkName("_configureLoginServiceDialog");                                                                 // 177
Template["_configureLoginServiceDialog"] = new Template("Template._configureLoginServiceDialog", (function() {        // 178
  var view = this;                                                                                                    // 179
  return Blaze.If(function() {                                                                                        // 180
    return Spacebars.call(view.lookup("visible"));                                                                    // 181
  }, function() {                                                                                                     // 182
    return [ "\n		", HTML.DIV({                                                                                       // 183
      "class": "modal",                                                                                               // 184
      id: "configure-login-service-dialog-modal"                                                                      // 185
    }, "\n			", HTML.DIV({                                                                                            // 186
      "class": "modal-dialog"                                                                                         // 187
    }, "\n				", HTML.DIV({                                                                                           // 188
      "class": "modal-content"                                                                                        // 189
    }, "\n					", HTML.DIV({                                                                                          // 190
      "class": "modal-header"                                                                                         // 191
    }, "\n						", HTML.BUTTON({                                                                                      // 192
      type: "button",                                                                                                 // 193
      "class": "close",                                                                                               // 194
      "data-dismiss": "modal",                                                                                        // 195
      "aria-label": "Close"                                                                                           // 196
    }, HTML.SPAN({                                                                                                    // 197
      "aria-hidden": "true"                                                                                           // 198
    }, HTML.CharRef({                                                                                                 // 199
      html: "&times;",                                                                                                // 200
      str: "×"                                                                                                        // 201
    }))), "\n						", HTML.H4({                                                                                       // 202
      "class": "modal-title"                                                                                          // 203
    }, "Configure Service"), "\n					"), "\n					", HTML.DIV({                                                        // 204
      "class": "modal-body"                                                                                           // 205
    }, "\n						", HTML.DIV({                                                                                         // 206
      id: "configure-login-service-dialog",                                                                           // 207
      "class": "accounts-dialog accounts-centered-dialog"                                                             // 208
    }, "\n								", Spacebars.include(view.lookupTemplate("configurationSteps")), "\n								", HTML.P("\n								Now, copy over some details.\n								"), "\n								", HTML.P("\n								", HTML.TABLE("\n									", HTML.COLGROUP("\n										", HTML.COL({
      span: "1",                                                                                                      // 210
      "class": "configuration_labels"                                                                                 // 211
    }), "\n										", HTML.COL({                                                                                    // 212
      span: "1",                                                                                                      // 213
      "class": "configuration_inputs"                                                                                 // 214
    }), "\n									"), "\n									", Blaze.Each(function() {                                                        // 215
      return Spacebars.call(view.lookup("configurationFields"));                                                      // 216
    }, function() {                                                                                                   // 217
      return [ "\n										", HTML.TR("\n											", HTML.TD("\n												", HTML.LABEL({                        // 218
        "for": function() {                                                                                           // 219
          return [ "configure-login-service-dialog-", Spacebars.mustache(view.lookup("property")) ];                  // 220
        }                                                                                                             // 221
      }, Blaze.View("lookup:label", function() {                                                                      // 222
        return Spacebars.mustache(view.lookup("label"));                                                              // 223
      })), "\n											"), "\n											", HTML.TD("\n												", HTML.INPUT({                                  // 224
        id: function() {                                                                                              // 225
          return [ "configure-login-service-dialog-", Spacebars.mustache(view.lookup("property")) ];                  // 226
        },                                                                                                            // 227
        type: "text"                                                                                                  // 228
      }), "\n											"), "\n										"), "\n									" ];                                                         // 229
    }), "\n								"), "\n								"), "\n								", HTML.P({                                                          // 230
      "class": "new-section"                                                                                          // 231
    }, "\n									Choose the login style:\n								"), "\n								", HTML.P("\n									", HTML.CharRef({            // 232
      html: "&emsp;",                                                                                                 // 233
      str: " "                                                                                                        // 234
    }), HTML.INPUT({                                                                                                  // 235
      id: "configure-login-service-dialog-popupBasedLogin",                                                           // 236
      type: "radio",                                                                                                  // 237
      checked: "checked",                                                                                             // 238
      name: "loginStyle",                                                                                             // 239
      value: "popup"                                                                                                  // 240
    }), "\n									", HTML.LABEL({                                                                                   // 241
      "for": "configure-login-service-dialog-popupBasedLogin"                                                         // 242
    }, "Popup-based login (recommended for most applications)"), "\n\n									", HTML.BR(), HTML.CharRef({           // 243
      html: "&emsp;",                                                                                                 // 244
      str: " "                                                                                                        // 245
    }), HTML.INPUT({                                                                                                  // 246
      id: "configure-login-service-dialog-redirectBasedLogin",                                                        // 247
      type: "radio",                                                                                                  // 248
      name: "loginStyle",                                                                                             // 249
      value: "redirect"                                                                                               // 250
    }), "\n									", HTML.LABEL({                                                                                   // 251
      "for": "configure-login-service-dialog-redirectBasedLogin"                                                      // 252
    }, "\n									Redirect-based login (special cases explained\n									", HTML.A({                                // 253
      href: "https://github.com/meteor/meteor/wiki/OAuth-for-mobile-Meteor-clients#popup-versus-redirect-flow",       // 254
      target: "_blank"                                                                                                // 255
    }, "here"), ")\n									"), "\n								"), "\n						"), "\n					"), "\n					", HTML.DIV({                        // 256
      "class": "modal-footer new-section"                                                                             // 257
    }, "\n						", HTML.DIV({                                                                                         // 258
      "class": "login-button btn btn-danger configure-login-service-dismiss-button"                                   // 259
    }, "\n							I'll do this later\n						"), "\n						", HTML.DIV({                                                 // 260
      "class": function() {                                                                                           // 261
        return [ "login-button login-button-configure btn btn-success ", Blaze.If(function() {                        // 262
          return Spacebars.call(view.lookup("saveDisabled"));                                                         // 263
        }, function() {                                                                                               // 264
          return "login-button-disabled";                                                                             // 265
        }) ];                                                                                                         // 266
      },                                                                                                              // 267
      id: "configure-login-service-dialog-save-configuration"                                                         // 268
    }, "\n							Save Configuration\n						"), "\n					"), "\n				"), "\n			"), "\n		"), "\n	" ];                     // 269
  });                                                                                                                 // 270
}));                                                                                                                  // 271
                                                                                                                      // 272
Template.__checkName("_loginButtonsMessagesDialog");                                                                  // 273
Template["_loginButtonsMessagesDialog"] = new Template("Template._loginButtonsMessagesDialog", (function() {          // 274
  var view = this;                                                                                                    // 275
  return HTML.DIV({                                                                                                   // 276
    "class": "modal",                                                                                                 // 277
    id: "login-buttons-message-dialog"                                                                                // 278
  }, "\n		", Blaze.If(function() {                                                                                    // 279
    return Spacebars.call(view.lookup("visible"));                                                                    // 280
  }, function() {                                                                                                     // 281
    return [ "\n		", HTML.DIV({                                                                                       // 282
      "class": "modal-dialog"                                                                                         // 283
    }, "\n			", HTML.DIV({                                                                                            // 284
      "class": "modal-content"                                                                                        // 285
    }, "\n				", HTML.DIV({                                                                                           // 286
      "class": "modal-header"                                                                                         // 287
    }, "\n					", HTML.BUTTON({                                                                                       // 288
      type: "button",                                                                                                 // 289
      "class": "close",                                                                                               // 290
      "data-dismiss": "modal",                                                                                        // 291
      "aria-label": "Close"                                                                                           // 292
    }, HTML.SPAN({                                                                                                    // 293
      "aria-hidden": "true"                                                                                           // 294
    }, HTML.CharRef({                                                                                                 // 295
      html: "&times;",                                                                                                // 296
      str: "×"                                                                                                        // 297
    }))), "\n					", HTML.H4({                                                                                        // 298
      "class": "modal-title"                                                                                          // 299
    }, Blaze.View("lookup:i18n", function() {                                                                         // 300
      return Spacebars.mustache(view.lookup("i18n"), "errorMessages.genericTitle");                                   // 301
    })), "\n				"), "\n				", HTML.DIV({                                                                              // 302
      "class": "modal-body"                                                                                           // 303
    }, "\n					", Spacebars.include(view.lookupTemplate("_loginButtonsMessages")), "\n				"), "\n				", HTML.DIV({    // 304
      "class": "modal-footer"                                                                                         // 305
    }, "\n					", HTML.BUTTON({                                                                                       // 306
      "class": "btn btn-primary login-button",                                                                        // 307
      id: "messages-dialog-dismiss-button",                                                                           // 308
      "data-dismiss": "modal"                                                                                         // 309
    }, Blaze.View("lookup:i18n", function() {                                                                         // 310
      return Spacebars.mustache(view.lookup("i18n"), "loginButtonsMessagesDialog.dismiss");                           // 311
    })), "\n				"), "\n			"), "\n		"), "\n		" ];                                                                      // 312
  }), "\n	");                                                                                                         // 313
}));                                                                                                                  // 314
                                                                                                                      // 315
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/login_buttons_session.js                                                      //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
(function () {                                                                                                        // 1
	var VALID_KEYS = [                                                                                                   // 2
		'dropdownVisible',                                                                                                  // 3
                                                                                                                      // 4
		// XXX consider replacing these with one key that has an enum for values.                                           // 5
		'inSignupFlow',                                                                                                     // 6
		'inForgotPasswordFlow',                                                                                             // 7
		'inChangePasswordFlow',                                                                                             // 8
		'inMessageOnlyFlow',                                                                                                // 9
                                                                                                                      // 10
		'errorMessage',                                                                                                     // 11
		'infoMessage',                                                                                                      // 12
                                                                                                                      // 13
		// dialogs with messages (info and error)                                                                           // 14
		'resetPasswordToken',                                                                                               // 15
		'enrollAccountToken',                                                                                               // 16
		'justVerifiedEmail',                                                                                                // 17
                                                                                                                      // 18
		'configureLoginServiceDialogVisible',                                                                               // 19
		'configureLoginServiceDialogServiceName',                                                                           // 20
		'configureLoginServiceDialogSaveDisabled'                                                                           // 21
	];                                                                                                                   // 22
                                                                                                                      // 23
	var validateKey = function (key) {                                                                                   // 24
		if (!_.contains(VALID_KEYS, key)){                                                                                  // 25
			throw new Error("Invalid key in loginButtonsSession: " + key);                                                     // 26
		}                                                                                                                   // 27
	};                                                                                                                   // 28
                                                                                                                      // 29
	var KEY_PREFIX = "Meteor.loginButtons.";                                                                             // 30
                                                                                                                      // 31
	// XXX we should have a better pattern for code private to a package like this one                                   // 32
	Accounts._loginButtonsSession = {                                                                                    // 33
		set: function(key, value) {                                                                                         // 34
			validateKey(key);                                                                                                  // 35
			if (_.contains(['errorMessage', 'infoMessage'], key)){                                                             // 36
				throw new Error("Don't set errorMessage or infoMessage directly. Instead, use errorMessage() or infoMessage().");
			}                                                                                                                  // 38
                                                                                                                      // 39
			this._set(key, value);                                                                                             // 40
		},                                                                                                                  // 41
                                                                                                                      // 42
		_set: function(key, value) {                                                                                        // 43
			Session.set(KEY_PREFIX + key, value);                                                                              // 44
		},                                                                                                                  // 45
                                                                                                                      // 46
		get: function(key) {                                                                                                // 47
			validateKey(key);                                                                                                  // 48
			return Session.get(KEY_PREFIX + key);                                                                              // 49
		},                                                                                                                  // 50
                                                                                                                      // 51
		closeDropdown: function () {                                                                                        // 52
			this.set('inSignupFlow', false);                                                                                   // 53
			this.set('inForgotPasswordFlow', false);                                                                           // 54
			this.set('inChangePasswordFlow', false);                                                                           // 55
			this.set('inMessageOnlyFlow', false);                                                                              // 56
			this.set('dropdownVisible', false);                                                                                // 57
			this.resetMessages();                                                                                              // 58
		},                                                                                                                  // 59
                                                                                                                      // 60
		infoMessage: function(message) {                                                                                    // 61
			this._set("errorMessage", null);                                                                                   // 62
			this._set("infoMessage", message);                                                                                 // 63
			this.ensureMessageVisible();                                                                                       // 64
		},                                                                                                                  // 65
                                                                                                                      // 66
		errorMessage: function(message) {                                                                                   // 67
			this._set("errorMessage", message);                                                                                // 68
			this._set("infoMessage", null);                                                                                    // 69
			this.ensureMessageVisible();                                                                                       // 70
		},                                                                                                                  // 71
                                                                                                                      // 72
		// is there a visible dialog that shows messages (info and error)                                                   // 73
		isMessageDialogVisible: function () {                                                                               // 74
			return this.get('resetPasswordToken') ||                                                                           // 75
				this.get('enrollAccountToken') ||                                                                                 // 76
				this.get('justVerifiedEmail');                                                                                    // 77
		},                                                                                                                  // 78
                                                                                                                      // 79
		// ensure that somethings displaying a message (info or error) is                                                   // 80
		// visible.  if a dialog with messages is open, do nothing;                                                         // 81
		// otherwise open the dropdown.                                                                                     // 82
		//                                                                                                                  // 83
		// notably this doesn't matter when only displaying a single login                                                  // 84
		// button since then we have an explicit message dialog                                                             // 85
		// (_loginButtonsMessageDialog), and dropdownVisible is ignored in                                                  // 86
		// this case.                                                                                                       // 87
		ensureMessageVisible: function () {                                                                                 // 88
			if (!this.isMessageDialogVisible()){                                                                               // 89
				this.set("dropdownVisible", true);                                                                                // 90
			}                                                                                                                  // 91
		},                                                                                                                  // 92
                                                                                                                      // 93
		resetMessages: function () {                                                                                        // 94
			this._set("errorMessage", null);                                                                                   // 95
			this._set("infoMessage", null);                                                                                    // 96
		},                                                                                                                  // 97
                                                                                                                      // 98
		configureService: function (name) {                                                                                 // 99
			this.set('configureLoginServiceDialogVisible', true);                                                              // 100
			this.set('configureLoginServiceDialogServiceName', name);                                                          // 101
			this.set('configureLoginServiceDialogSaveDisabled', true);                                                         // 102
			setTimeout(function(){                                                                                             // 103
				$('#configure-login-service-dialog-modal').modal();                                                               // 104
			}, 500)                                                                                                            // 105
		}                                                                                                                   // 106
	};                                                                                                                   // 107
}) ();                                                                                                                // 108
                                                                                                                      // 109
                                                                                                                      // 110
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/login_buttons.js                                                              //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
(function() {                                                                                                         // 1
	if (!Accounts._loginButtons){                                                                                        // 2
		Accounts._loginButtons = {};                                                                                        // 3
	}                                                                                                                    // 4
                                                                                                                      // 5
	// for convenience                                                                                                   // 6
	var loginButtonsSession = Accounts._loginButtonsSession;                                                             // 7
                                                                                                                      // 8
	UI.registerHelper("loginButtons", function() {                                                                       // 9
		return Template._loginButtons;                                                                                      // 10
	});                                                                                                                  // 11
                                                                                                                      // 12
	// shared between dropdown and single mode                                                                           // 13
	Template._loginButtons.events({                                                                                      // 14
		'click #login-buttons-logout': function() {                                                                         // 15
			Meteor.logout(function(error) {                                                                                    // 16
				loginButtonsSession.closeDropdown();                                                                              // 17
				if (typeof accountsUIBootstrap3.logoutCallback === 'function') {                                                  // 18
					accountsUIBootstrap3.logoutCallback(error);                                                                      // 19
				}                                                                                                                 // 20
			});                                                                                                                // 21
		}                                                                                                                   // 22
	});                                                                                                                  // 23
                                                                                                                      // 24
	//                                                                                                                   // 25
	// loginButtonLoggedOut template                                                                                     // 26
	//                                                                                                                   // 27
	Template._loginButtonsLoggedOut.helpers({                                                                            // 28
		dropdown: function() {                                                                                              // 29
			return Accounts._loginButtons.dropdown();                                                                          // 30
		},                                                                                                                  // 31
		services: function() {                                                                                              // 32
			return Accounts._loginButtons.getLoginServices();                                                                  // 33
		},                                                                                                                  // 34
		singleService: function() {                                                                                         // 35
			var services = Accounts._loginButtons.getLoginServices();                                                          // 36
			if (services.length !== 1){                                                                                        // 37
				throw new Error(                                                                                                  // 38
					"Shouldn't be rendering this template with more than one configured service");                                   // 39
			}                                                                                                                  // 40
			return services[0];                                                                                                // 41
		},                                                                                                                  // 42
		configurationLoaded: function() {                                                                                   // 43
			return Accounts.loginServicesConfigured();                                                                         // 44
		}                                                                                                                   // 45
	});                                                                                                                  // 46
                                                                                                                      // 47
                                                                                                                      // 48
                                                                                                                      // 49
	//                                                                                                                   // 50
	// loginButtonsLoggedIn template                                                                                     // 51
	//                                                                                                                   // 52
                                                                                                                      // 53
	// decide whether we should show a dropdown rather than a row of                                                     // 54
	// buttons                                                                                                           // 55
	Template._loginButtonsLoggedIn.helpers({                                                                             // 56
		dropdown: function() {                                                                                              // 57
			return Accounts._loginButtons.dropdown();                                                                          // 58
		},                                                                                                                  // 59
		displayName: function() {                                                                                           // 60
			return Accounts._loginButtons.displayName();                                                                       // 61
		}                                                                                                                   // 62
	})                                                                                                                   // 63
                                                                                                                      // 64
                                                                                                                      // 65
                                                                                                                      // 66
	//                                                                                                                   // 67
	// loginButtonsMessage template                                                                                      // 68
	//                                                                                                                   // 69
                                                                                                                      // 70
	Template._loginButtonsMessages.helpers({                                                                             // 71
		errorMessage: function() {                                                                                          // 72
			return loginButtonsSession.get('errorMessage');                                                                    // 73
		},                                                                                                                  // 74
		infoMessage: function() {                                                                                           // 75
			return loginButtonsSession.get('infoMessage');                                                                     // 76
		}                                                                                                                   // 77
	});                                                                                                                  // 78
                                                                                                                      // 79
                                                                                                                      // 80
                                                                                                                      // 81
	//                                                                                                                   // 82
	// helpers                                                                                                           // 83
	//                                                                                                                   // 84
                                                                                                                      // 85
	Accounts._loginButtons.displayName = function() {                                                                    // 86
		var user = Meteor.user();                                                                                           // 87
		if (!user){                                                                                                         // 88
			return '';                                                                                                         // 89
		}                                                                                                                   // 90
                                                                                                                      // 91
		if (user.profile && user.profile.name){                                                                             // 92
			return user.profile.name;                                                                                          // 93
		}                                                                                                                   // 94
		if (user.username){                                                                                                 // 95
			return user.username;                                                                                              // 96
		}                                                                                                                   // 97
		if (user.emails && user.emails[0] && user.emails[0].address){                                                       // 98
			return user.emails[0].address;                                                                                     // 99
		}                                                                                                                   // 100
                                                                                                                      // 101
		return '';                                                                                                          // 102
	};                                                                                                                   // 103
                                                                                                                      // 104
	Accounts._loginButtons.getLoginServices = function() {                                                               // 105
		// First look for OAuth services.                                                                                   // 106
		var services = Package['accounts-oauth'] ? Accounts.oauth.serviceNames() : [];                                      // 107
                                                                                                                      // 108
		// Be equally kind to all login services. This also preserves                                                       // 109
		// backwards-compatibility. (But maybe order should be                                                              // 110
		// configurable?)                                                                                                   // 111
		services.sort();                                                                                                    // 112
                                                                                                                      // 113
		// Add password, if it's there; it must come last.                                                                  // 114
		if (this.hasPasswordService()){                                                                                     // 115
			services.push('password');                                                                                         // 116
		}                                                                                                                   // 117
                                                                                                                      // 118
		return _.map(services, function(name) {                                                                             // 119
			return {                                                                                                           // 120
				name: name                                                                                                        // 121
			};                                                                                                                 // 122
		});                                                                                                                 // 123
	};                                                                                                                   // 124
                                                                                                                      // 125
	Accounts._loginButtons.hasPasswordService = function() {                                                             // 126
		return !!Package['accounts-password'];                                                                              // 127
	};                                                                                                                   // 128
                                                                                                                      // 129
	Accounts._loginButtons.dropdown = function() {                                                                       // 130
		return this.hasPasswordService() || Accounts._loginButtons.getLoginServices().length > 1;                           // 131
	};                                                                                                                   // 132
                                                                                                                      // 133
	// XXX improve these. should this be in accounts-password instead?                                                   // 134
	//                                                                                                                   // 135
	// XXX these will become configurable, and will be validated on                                                      // 136
	// the server as well.                                                                                               // 137
	Accounts._loginButtons.validateUsername = function(username) {                                                       // 138
		if (username.length >= 3) {                                                                                         // 139
			return true;                                                                                                       // 140
		} else {                                                                                                            // 141
			loginButtonsSession.errorMessage(i18n('errorMessages.usernameTooShort'));                                          // 142
			return false;                                                                                                      // 143
		}                                                                                                                   // 144
	};                                                                                                                   // 145
	Accounts._loginButtons.validateEmail = function(email) {                                                             // 146
		if (Accounts.ui._passwordSignupFields() === "USERNAME_AND_OPTIONAL_EMAIL" && email === ''){                         // 147
			return true;                                                                                                       // 148
		}                                                                                                                   // 149
                                                                                                                      // 150
		var re = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
                                                                                                                      // 152
		if (re.test(email)) {                                                                                               // 153
			return true;                                                                                                       // 154
		} else {                                                                                                            // 155
			loginButtonsSession.errorMessage(i18n('errorMessages.invalidEmail'));                                              // 156
			return false;                                                                                                      // 157
		}                                                                                                                   // 158
	};                                                                                                                   // 159
	Accounts._loginButtons.validatePassword = function(password, passwordAgain) {                                        // 160
		if (password.length >= 6) {                                                                                         // 161
			if (typeof passwordAgain !== "undefined" && passwordAgain !== null && password != passwordAgain) {                 // 162
				loginButtonsSession.errorMessage(i18n('errorMessages.passwordsDontMatch'));                                       // 163
				return false;                                                                                                     // 164
			}                                                                                                                  // 165
			return true;                                                                                                       // 166
		} else {                                                                                                            // 167
			loginButtonsSession.errorMessage(i18n('errorMessages.passwordTooShort'));                                          // 168
			return false;                                                                                                      // 169
		}                                                                                                                   // 170
	};                                                                                                                   // 171
                                                                                                                      // 172
	Accounts._loginButtons.rendered = function() {                                                                       // 173
		debugger;                                                                                                           // 174
	};                                                                                                                   // 175
                                                                                                                      // 176
})();                                                                                                                 // 177
                                                                                                                      // 178
                                                                                                                      // 179
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/login_buttons_single.js                                                       //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
(function() {                                                                                                         // 1
	// for convenience                                                                                                   // 2
	var loginButtonsSession = Accounts._loginButtonsSession;                                                             // 3
                                                                                                                      // 4
	Template._loginButtonsLoggedOutSingleLoginButton.events({                                                            // 5
		'click .login-button': function() {                                                                                 // 6
			var serviceName = this.name;                                                                                       // 7
			loginButtonsSession.resetMessages();                                                                               // 8
			var callback = function(err) {                                                                                     // 9
				if (!err) {                                                                                                       // 10
					loginButtonsSession.closeDropdown();                                                                             // 11
				} else if (err instanceof Accounts.LoginCancelledError) {                                                         // 12
					// do nothing                                                                                                    // 13
				} else if (err instanceof Accounts.ConfigError) {                                                                 // 14
					loginButtonsSession.configureService(serviceName);                                                               // 15
				} else {                                                                                                          // 16
					loginButtonsSession.errorMessage(err.reason || "Unknown error");                                                 // 17
				}                                                                                                                 // 18
			};                                                                                                                 // 19
                                                                                                                      // 20
			// XXX Service providers should be able to specify their                                                           // 21
			// `Meteor.loginWithX` method name.                                                                                // 22
			var loginWithService = Meteor["loginWith" + (serviceName === 'meteor-developer' ?  'MeteorDeveloperAccount' :  capitalize(serviceName))];
                                                                                                                      // 24
			var options = {}; // use default scope unless specified                                                            // 25
			if (Accounts.ui._options.requestPermissions[serviceName])                                                          // 26
				options.requestPermissions = Accounts.ui._options.requestPermissions[serviceName];                                // 27
			if (Accounts.ui._options.requestOfflineToken[serviceName])                                                         // 28
				options.requestOfflineToken = Accounts.ui._options.requestOfflineToken[serviceName];                              // 29
			if (Accounts.ui._options.forceApprovalPrompt[serviceName])                                                         // 30
				options.forceApprovalPrompt = Accounts.ui._options.forceApprovalPrompt[serviceName];                              // 31
                                                                                                                      // 32
			loginWithService(options, callback);                                                                               // 33
		}                                                                                                                   // 34
	});                                                                                                                  // 35
                                                                                                                      // 36
	Template._loginButtonsLoggedOutSingleLoginButton.helpers({                                                           // 37
		configured: function() {                                                                                            // 38
			return !!Accounts.loginServiceConfiguration.findOne({                                                              // 39
				service: this.name                                                                                                // 40
			});                                                                                                                // 41
		},                                                                                                                  // 42
		capitalizedName: function() {                                                                                       // 43
			if (this.name === 'github'){                                                                                       // 44
			// XXX we should allow service packages to set their capitalized name                                              // 45
				return 'GitHub';                                                                                                  // 46
			} else {                                                                                                           // 47
				return capitalize(this.name);                                                                                     // 48
			}                                                                                                                  // 49
		}                                                                                                                   // 50
	});                                                                                                                  // 51
                                                                                                                      // 52
                                                                                                                      // 53
	// XXX from http://epeli.github.com/underscore.string/lib/underscore.string.js                                       // 54
	var capitalize = function(str) {                                                                                     // 55
		str = (str == null) ? '' : String(str);                                                                             // 56
		return str.charAt(0).toUpperCase() + str.slice(1);                                                                  // 57
	};                                                                                                                   // 58
})();                                                                                                                 // 59
                                                                                                                      // 60
                                                                                                                      // 61
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/login_buttons_dropdown.js                                                     //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
(function() {                                                                                                         // 1
                                                                                                                      // 2
	// for convenience                                                                                                   // 3
	var loginButtonsSession = Accounts._loginButtonsSession;                                                             // 4
                                                                                                                      // 5
	// events shared between loginButtonsLoggedOutDropdown and                                                           // 6
	// loginButtonsLoggedInDropdown                                                                                      // 7
	Template._loginButtons.events({                                                                                      // 8
		'click input, click .radio, click .checkbox, click option, click select': function(event) {                         // 9
			event.stopPropagation();                                                                                           // 10
		},                                                                                                                  // 11
		'click #login-name-link, click #login-sign-in-link': function(event) {                                              // 12
			event.stopPropagation();                                                                                           // 13
			loginButtonsSession.set('dropdownVisible', true);                                                                  // 14
			Meteor.flush();                                                                                                    // 15
		},                                                                                                                  // 16
		'click .login-close': function() {                                                                                  // 17
			loginButtonsSession.closeDropdown();                                                                               // 18
		}                                                                                                                   // 19
	});                                                                                                                  // 20
                                                                                                                      // 21
	Template._loginButtons.toggleDropdown = function() {                                                                 // 22
		toggleDropdown();                                                                                                   // 23
		focusInput();                                                                                                       // 24
	};                                                                                                                   // 25
                                                                                                                      // 26
	//                                                                                                                   // 27
	// loginButtonsLoggedInDropdown template and related                                                                 // 28
	//                                                                                                                   // 29
                                                                                                                      // 30
	Template._loginButtonsLoggedInDropdown.events({                                                                      // 31
		'click #login-buttons-open-change-password': function(event) {                                                      // 32
			event.stopPropagation();                                                                                           // 33
			loginButtonsSession.resetMessages();                                                                               // 34
			loginButtonsSession.set('inChangePasswordFlow', true);                                                             // 35
			Meteor.flush();                                                                                                    // 36
		}                                                                                                                   // 37
	});                                                                                                                  // 38
                                                                                                                      // 39
	Template._loginButtonsLoggedInDropdown.helpers({                                                                     // 40
		displayName: function() {                                                                                           // 41
			return Accounts._loginButtons.displayName();                                                                       // 42
		},                                                                                                                  // 43
                                                                                                                      // 44
		inChangePasswordFlow: function() {                                                                                  // 45
			return loginButtonsSession.get('inChangePasswordFlow');                                                            // 46
		},                                                                                                                  // 47
                                                                                                                      // 48
		inMessageOnlyFlow: function() {                                                                                     // 49
			return loginButtonsSession.get('inMessageOnlyFlow');                                                               // 50
		},                                                                                                                  // 51
                                                                                                                      // 52
		dropdownVisible: function() {                                                                                       // 53
			return loginButtonsSession.get('dropdownVisible');                                                                 // 54
		},                                                                                                                  // 55
                                                                                                                      // 56
		user_profile_picture: function() {                                                                                  // 57
			var user = Meteor.user();                                                                                          // 58
			if (user && user.profile && user.profile.display_picture) {                                                        // 59
				return user.profile.display_picture;                                                                              // 60
			}                                                                                                                  // 61
			return "";                                                                                                         // 62
		}                                                                                                                   // 63
	});                                                                                                                  // 64
                                                                                                                      // 65
                                                                                                                      // 66
	Template._loginButtonsLoggedInDropdownActions.helpers({                                                              // 67
		allowChangingPassword: function() {                                                                                 // 68
			// it would be more correct to check whether the user has a password set,                                          // 69
			// but in order to do that we'd have to send more data down to the client,                                         // 70
			// and it'd be preferable not to send down the entire service.password document.                                   // 71
			//                                                                                                                 // 72
			// instead we use the heuristic: if the user has a username or email set.                                          // 73
			var user = Meteor.user();                                                                                          // 74
			return user.username || (user.emails && user.emails[0] && user.emails[0].address);                                 // 75
		},                                                                                                                  // 76
		additionalLoggedInDropdownActions: function() {                                                                     // 77
			return Template._loginButtonsAdditionalLoggedInDropdownActions !== undefined;                                      // 78
		}                                                                                                                   // 79
	});                                                                                                                  // 80
                                                                                                                      // 81
                                                                                                                      // 82
	//                                                                                                                   // 83
	// loginButtonsLoggedOutDropdown template and related                                                                // 84
	//                                                                                                                   // 85
                                                                                                                      // 86
	Template._loginButtonsLoggedOutAllServices.events({                                                                  // 87
		'click #login-buttons-password': function(event) {                                                                  // 88
			event.stopPropagation();                                                                                           // 89
			loginOrSignup();                                                                                                   // 90
		},                                                                                                                  // 91
                                                                                                                      // 92
		'keypress #forgot-password-email': function(event) {                                                                // 93
			event.stopPropagation();                                                                                           // 94
			if (event.keyCode === 13){                                                                                         // 95
				forgotPassword();                                                                                                 // 96
			}                                                                                                                  // 97
		},                                                                                                                  // 98
                                                                                                                      // 99
		'click #login-buttons-forgot-password': function(event) {                                                           // 100
			event.stopPropagation();                                                                                           // 101
			forgotPassword();                                                                                                  // 102
		},                                                                                                                  // 103
                                                                                                                      // 104
		'click #signup-link': function(event) {                                                                             // 105
			event.stopPropagation();                                                                                           // 106
			loginButtonsSession.resetMessages();                                                                               // 107
                                                                                                                      // 108
			//check to see if onCreate is populated with a function. If it is, call it                                         // 109
			var onCreateFn = Accounts.ui._options.onCreate;                                                                    // 110
			if (onCreateFn){                                                                                                   // 111
				loginButtonsSession.closeDropdown();                                                                              // 112
				onCreateFn.apply();                                                                                               // 113
                                                                                                                      // 114
			} else {                                                                                                           // 115
				// store values of fields before swtiching to the signup form                                                     // 116
				var username = trimmedElementValueById('login-username');                                                         // 117
				var email = trimmedElementValueById('login-email');                                                               // 118
				var usernameOrEmail = trimmedElementValueById('login-username-or-email');                                         // 119
				// notably not trimmed. a password could (?) start or end with a space                                            // 120
				var password = elementValueById('login-password');                                                                // 121
                                                                                                                      // 122
				loginButtonsSession.set('inSignupFlow', true);                                                                    // 123
				loginButtonsSession.set('inForgotPasswordFlow', false);                                                           // 124
                                                                                                                      // 125
				// force the ui to update so that we have the approprate fields to fill in                                        // 126
				Meteor.flush();                                                                                                   // 127
                                                                                                                      // 128
				// update new fields with appropriate defaults                                                                    // 129
				if (username !== null) {                                                                                          // 130
					document.getElementById('login-username').value = username;                                                      // 131
				} else if (email !== null) {                                                                                      // 132
					document.getElementById('login-email').value = email;                                                            // 133
				} else if (usernameOrEmail !== null) {                                                                            // 134
					if (usernameOrEmail.indexOf('@') === -1) {                                                                       // 135
						document.getElementById('login-username').value = usernameOrEmail;                                              // 136
					} else {                                                                                                         // 137
						document.getElementById('login-email').value = usernameOrEmail;                                                 // 138
					}                                                                                                                // 139
				}                                                                                                                 // 140
			}                                                                                                                  // 141
		},                                                                                                                  // 142
		'click #forgot-password-link': function(event) {                                                                    // 143
			event.stopPropagation();                                                                                           // 144
			loginButtonsSession.resetMessages();                                                                               // 145
                                                                                                                      // 146
			// store values of fields before swtiching to the signup form                                                      // 147
			var email = trimmedElementValueById('login-email');                                                                // 148
			var usernameOrEmail = trimmedElementValueById('login-username-or-email');                                          // 149
                                                                                                                      // 150
			loginButtonsSession.set('inSignupFlow', false);                                                                    // 151
			loginButtonsSession.set('inForgotPasswordFlow', true);                                                             // 152
                                                                                                                      // 153
			// force the ui to update so that we have the approprate fields to fill in                                         // 154
			Meteor.flush();                                                                                                    // 155
			//toggleDropdown();                                                                                                // 156
                                                                                                                      // 157
			// update new fields with appropriate defaults                                                                     // 158
			if (email !== null){                                                                                               // 159
				document.getElementById('forgot-password-email').value = email;                                                   // 160
			} else if (usernameOrEmail !== null){                                                                              // 161
				if (usernameOrEmail.indexOf('@') !== -1){                                                                         // 162
					document.getElementById('forgot-password-email').value = usernameOrEmail;                                        // 163
				}                                                                                                                 // 164
			}                                                                                                                  // 165
		},                                                                                                                  // 166
		'click #back-to-login-link': function(event) {                                                                      // 167
			event.stopPropagation();                                                                                           // 168
			loginButtonsSession.resetMessages();                                                                               // 169
                                                                                                                      // 170
			var username = trimmedElementValueById('login-username');                                                          // 171
			var email = trimmedElementValueById('login-email') || trimmedElementValueById('forgot-password-email'); // Ughh. Standardize on names?
                                                                                                                      // 173
			loginButtonsSession.set('inSignupFlow', false);                                                                    // 174
			loginButtonsSession.set('inForgotPasswordFlow', false);                                                            // 175
                                                                                                                      // 176
			// force the ui to update so that we have the approprate fields to fill in                                         // 177
			Meteor.flush();                                                                                                    // 178
                                                                                                                      // 179
			if (document.getElementById('login-username')){                                                                    // 180
				document.getElementById('login-username').value = username;                                                       // 181
			}                                                                                                                  // 182
			if (document.getElementById('login-email')){                                                                       // 183
				document.getElementById('login-email').value = email;                                                             // 184
			}                                                                                                                  // 185
			// "login-password" is preserved thanks to the preserve-inputs package                                             // 186
			if (document.getElementById('login-username-or-email')){                                                           // 187
				document.getElementById('login-username-or-email').value = email || username;                                     // 188
			}                                                                                                                  // 189
		},                                                                                                                  // 190
		'keypress #login-username, keypress #login-email, keypress #login-username-or-email, keypress #login-password, keypress #login-password-again': function(event) {
			if (event.keyCode === 13){                                                                                         // 192
				loginOrSignup();                                                                                                  // 193
			}                                                                                                                  // 194
		}                                                                                                                   // 195
	});                                                                                                                  // 196
                                                                                                                      // 197
	Template._loginButtonsLoggedOutDropdown.helpers({                                                                    // 198
		forbidClientAccountCreation: function() {                                                                           // 199
			return Accounts._options.forbidClientAccountCreation;                                                              // 200
		}                                                                                                                   // 201
	});                                                                                                                  // 202
                                                                                                                      // 203
	Template._loginButtonsLoggedOutAllServices.helpers({                                                                 // 204
		// additional classes that can be helpful in styling the dropdown                                                   // 205
		additionalClasses: function() {                                                                                     // 206
			if (!Accounts.password) {                                                                                          // 207
				return false;                                                                                                     // 208
			} else {                                                                                                           // 209
				if (loginButtonsSession.get('inSignupFlow')) {                                                                    // 210
					return 'login-form-create-account';                                                                              // 211
				} else if (loginButtonsSession.get('inForgotPasswordFlow')) {                                                     // 212
					return 'login-form-forgot-password';                                                                             // 213
				} else {                                                                                                          // 214
					return 'login-form-sign-in';                                                                                     // 215
				}                                                                                                                 // 216
			}                                                                                                                  // 217
		},                                                                                                                  // 218
                                                                                                                      // 219
		dropdownVisible: function() {                                                                                       // 220
			return loginButtonsSession.get('dropdownVisible');                                                                 // 221
		},                                                                                                                  // 222
                                                                                                                      // 223
		services: function() {                                                                                              // 224
			return Accounts._loginButtons.getLoginServices();                                                                  // 225
		},                                                                                                                  // 226
                                                                                                                      // 227
		isPasswordService: function() {                                                                                     // 228
			return this.name === 'password';                                                                                   // 229
		},                                                                                                                  // 230
                                                                                                                      // 231
		hasOtherServices: function() {                                                                                      // 232
			return Accounts._loginButtons.getLoginServices().length > 1;                                                       // 233
		},                                                                                                                  // 234
                                                                                                                      // 235
		hasPasswordService: function() {                                                                                    // 236
			return Accounts._loginButtons.hasPasswordService();                                                                // 237
		}                                                                                                                   // 238
	});                                                                                                                  // 239
                                                                                                                      // 240
                                                                                                                      // 241
	Template._loginButtonsLoggedOutPasswordService.helpers({                                                             // 242
		fields: function() {                                                                                                // 243
			var loginFields = [{                                                                                               // 244
				fieldName: 'username-or-email',                                                                                   // 245
				fieldLabel: i18n('loginFields.usernameOrEmail'),                                                                  // 246
				visible: function() {                                                                                             // 247
					return _.contains(                                                                                               // 248
						["USERNAME_AND_EMAIL_CONFIRM", "USERNAME_AND_EMAIL", "USERNAME_AND_OPTIONAL_EMAIL"],                            // 249
						Accounts.ui._passwordSignupFields());                                                                           // 250
				}                                                                                                                 // 251
			}, {                                                                                                               // 252
				fieldName: 'username',                                                                                            // 253
				fieldLabel: i18n('loginFields.username'),                                                                         // 254
				visible: function() {                                                                                             // 255
					return Accounts.ui._passwordSignupFields() === "USERNAME_ONLY";                                                  // 256
				}                                                                                                                 // 257
			}, {                                                                                                               // 258
				fieldName: 'email',                                                                                               // 259
				fieldLabel: i18n('loginFields.email'),                                                                            // 260
				inputType: 'email',                                                                                               // 261
				visible: function() {                                                                                             // 262
					return Accounts.ui._passwordSignupFields() === "EMAIL_ONLY";                                                     // 263
				}                                                                                                                 // 264
			}, {                                                                                                               // 265
				fieldName: 'password',                                                                                            // 266
				fieldLabel: i18n('loginFields.password'),                                                                         // 267
				inputType: 'password',                                                                                            // 268
				visible: function() {                                                                                             // 269
					return true;                                                                                                     // 270
				}                                                                                                                 // 271
			}];                                                                                                                // 272
                                                                                                                      // 273
			var signupFields = [{                                                                                              // 274
				fieldName: 'username',                                                                                            // 275
				fieldLabel: i18n('signupFields.username'),                                                                        // 276
				visible: function() {                                                                                             // 277
					return _.contains(                                                                                               // 278
						["USERNAME_AND_EMAIL_CONFIRM", "USERNAME_AND_EMAIL", "USERNAME_AND_OPTIONAL_EMAIL", "USERNAME_ONLY"],           // 279
						Accounts.ui._passwordSignupFields());                                                                           // 280
				}                                                                                                                 // 281
			}, {                                                                                                               // 282
				fieldName: 'email',                                                                                               // 283
				fieldLabel: i18n('signupFields.email'),                                                                           // 284
				inputType: 'email',                                                                                               // 285
				visible: function() {                                                                                             // 286
					return _.contains(                                                                                               // 287
						["USERNAME_AND_EMAIL_CONFIRM", "USERNAME_AND_EMAIL", "EMAIL_ONLY"],                                             // 288
						Accounts.ui._passwordSignupFields());                                                                           // 289
				}                                                                                                                 // 290
			}, {                                                                                                               // 291
				fieldName: 'email',                                                                                               // 292
				fieldLabel: i18n('signupFields.emailOpt'),                                                                        // 293
				inputType: 'email',                                                                                               // 294
				visible: function() {                                                                                             // 295
					return Accounts.ui._passwordSignupFields() === "USERNAME_AND_OPTIONAL_EMAIL";                                    // 296
				}                                                                                                                 // 297
			}, {                                                                                                               // 298
				fieldName: 'password',                                                                                            // 299
				fieldLabel: i18n('signupFields.password'),                                                                        // 300
				inputType: 'password',                                                                                            // 301
				visible: function() {                                                                                             // 302
					return true;                                                                                                     // 303
				}                                                                                                                 // 304
			}, {                                                                                                               // 305
				fieldName: 'password-again',                                                                                      // 306
				fieldLabel: i18n('signupFields.passwordAgain'),                                                                   // 307
				inputType: 'password',                                                                                            // 308
				visible: function() {                                                                                             // 309
					// No need to make users double-enter their password if                                                          // 310
					// they'll necessarily have an email set, since they can use                                                     // 311
					// the "forgot password" flow.                                                                                   // 312
					return _.contains(                                                                                               // 313
						["USERNAME_AND_EMAIL_CONFIRM", "USERNAME_AND_OPTIONAL_EMAIL", "USERNAME_ONLY"],                                 // 314
						Accounts.ui._passwordSignupFields());                                                                           // 315
				}                                                                                                                 // 316
			}];                                                                                                                // 317
                                                                                                                      // 318
			signupFields = signupFields.concat(Accounts.ui._options.extraSignupFields);                                        // 319
                                                                                                                      // 320
			return loginButtonsSession.get('inSignupFlow') ? signupFields : loginFields;                                       // 321
		},                                                                                                                  // 322
                                                                                                                      // 323
		inForgotPasswordFlow: function() {                                                                                  // 324
			return loginButtonsSession.get('inForgotPasswordFlow');                                                            // 325
		},                                                                                                                  // 326
                                                                                                                      // 327
		inLoginFlow: function() {                                                                                           // 328
			return !loginButtonsSession.get('inSignupFlow') && !loginButtonsSession.get('inForgotPasswordFlow');               // 329
		},                                                                                                                  // 330
                                                                                                                      // 331
		inSignupFlow: function() {                                                                                          // 332
			return loginButtonsSession.get('inSignupFlow');                                                                    // 333
		},                                                                                                                  // 334
                                                                                                                      // 335
		showForgotPasswordLink: function() {                                                                                // 336
			return _.contains(                                                                                                 // 337
				["USERNAME_AND_EMAIL_CONFIRM", "USERNAME_AND_EMAIL", "USERNAME_AND_OPTIONAL_EMAIL", "EMAIL_ONLY"],                // 338
				Accounts.ui._passwordSignupFields());                                                                             // 339
		},                                                                                                                  // 340
                                                                                                                      // 341
		showCreateAccountLink: function() {                                                                                 // 342
			return !Accounts._options.forbidClientAccountCreation;                                                             // 343
		}                                                                                                                   // 344
	});                                                                                                                  // 345
                                                                                                                      // 346
	Template._loginButtonsFormField.helpers({                                                                            // 347
		equals: function(a, b) {                                                                                            // 348
			return (a === b);                                                                                                  // 349
		},                                                                                                                  // 350
		inputType: function() {                                                                                             // 351
			return this.inputType || "text";                                                                                   // 352
		},                                                                                                                  // 353
		inputTextual: function() {                                                                                          // 354
			return !_.contains(["radio", "checkbox", "select"], this.inputType);                                               // 355
		}                                                                                                                   // 356
	});                                                                                                                  // 357
                                                                                                                      // 358
	//                                                                                                                   // 359
	// loginButtonsChangePassword template                                                                               // 360
	//                                                                                                                   // 361
	Template._loginButtonsChangePassword.events({                                                                        // 362
		'keypress #login-old-password, keypress #login-password, keypress #login-password-again': function(event) {         // 363
			if (event.keyCode === 13){                                                                                         // 364
				changePassword();                                                                                                 // 365
			}                                                                                                                  // 366
		},                                                                                                                  // 367
		'click #login-buttons-do-change-password': function(event) {                                                        // 368
			event.stopPropagation();                                                                                           // 369
			changePassword();                                                                                                  // 370
		},                                                                                                                  // 371
		'click #login-buttons-cancel-change-password': function(event) {                                                    // 372
			event.stopPropagation();                                                                                           // 373
			loginButtonsSession.resetMessages();                                                                               // 374
			Accounts._loginButtonsSession.set('inChangePasswordFlow', false);                                                  // 375
			Meteor.flush();                                                                                                    // 376
		}                                                                                                                   // 377
	});                                                                                                                  // 378
                                                                                                                      // 379
	Template._loginButtonsChangePassword.helpers({                                                                       // 380
		fields: function() {                                                                                                // 381
			return [{                                                                                                          // 382
				fieldName: 'old-password',                                                                                        // 383
				fieldLabel: i18n('changePasswordFields.currentPassword'),                                                         // 384
				inputType: 'password',                                                                                            // 385
				visible: function() {                                                                                             // 386
					return true;                                                                                                     // 387
				}                                                                                                                 // 388
			}, {                                                                                                               // 389
				fieldName: 'password',                                                                                            // 390
				fieldLabel: i18n('changePasswordFields.newPassword'),                                                             // 391
				inputType: 'password',                                                                                            // 392
				visible: function() {                                                                                             // 393
					return true;                                                                                                     // 394
				}                                                                                                                 // 395
			}, {                                                                                                               // 396
				fieldName: 'password-again',                                                                                      // 397
				fieldLabel: i18n('changePasswordFields.newPasswordAgain'),                                                        // 398
				inputType: 'password',                                                                                            // 399
				visible: function() {                                                                                             // 400
					// No need to make users double-enter their password if                                                          // 401
					// they'll necessarily have an email set, since they can use                                                     // 402
					// the "forgot password" flow.                                                                                   // 403
					return _.contains(                                                                                               // 404
						["USERNAME_AND_OPTIONAL_EMAIL", "USERNAME_ONLY"],                                                               // 405
						Accounts.ui._passwordSignupFields());                                                                           // 406
				}                                                                                                                 // 407
			}];                                                                                                                // 408
		}                                                                                                                   // 409
	});                                                                                                                  // 410
                                                                                                                      // 411
	//                                                                                                                   // 412
	// helpers                                                                                                           // 413
	//                                                                                                                   // 414
                                                                                                                      // 415
	var elementValueById = function(id) {                                                                                // 416
		var element = document.getElementById(id);                                                                          // 417
		if (!element){                                                                                                      // 418
			return null;                                                                                                       // 419
		} else {                                                                                                            // 420
			return element.value;                                                                                              // 421
		}                                                                                                                   // 422
	};                                                                                                                   // 423
                                                                                                                      // 424
	var elementValueByIdForRadio = function(fieldIdPrefix, radioOptions) {                                               // 425
		var value = null;                                                                                                   // 426
		for (i in radioOptions) {                                                                                           // 427
			var element = document.getElementById(fieldIdPrefix + '-' + radioOptions[i].id);                                   // 428
			if (element && element.checked){                                                                                   // 429
				value =  element.value;                                                                                           // 430
			}                                                                                                                  // 431
		}                                                                                                                   // 432
		return value;                                                                                                       // 433
	};                                                                                                                   // 434
                                                                                                                      // 435
	var elementValueByIdForCheckbox = function(id) {                                                                     // 436
		var element = document.getElementById(id);                                                                          // 437
		return element.checked;                                                                                             // 438
	};                                                                                                                   // 439
                                                                                                                      // 440
	var trimmedElementValueById = function(id) {                                                                         // 441
		var element = document.getElementById(id);                                                                          // 442
		if (!element){                                                                                                      // 443
			return null;                                                                                                       // 444
		} else {                                                                                                            // 445
			return element.value.replace(/^\s*|\s*$/g, ""); // trim;                                                           // 446
		}                                                                                                                   // 447
	};                                                                                                                   // 448
                                                                                                                      // 449
	var loginOrSignup = function() {                                                                                     // 450
		if (loginButtonsSession.get('inSignupFlow')){                                                                       // 451
			signup();                                                                                                          // 452
		} else {                                                                                                            // 453
			login();                                                                                                           // 454
		}                                                                                                                   // 455
	};                                                                                                                   // 456
                                                                                                                      // 457
	var login = function() {                                                                                             // 458
		loginButtonsSession.resetMessages();                                                                                // 459
                                                                                                                      // 460
		var username = trimmedElementValueById('login-username');                                                           // 461
		if (username && Accounts.ui._options.forceUsernameLowercase) {                                                      // 462
			username = username.toLowerCase();                                                                                 // 463
		}                                                                                                                   // 464
		var email = trimmedElementValueById('login-email');                                                                 // 465
		if (email && Accounts.ui._options.forceEmailLowercase) {                                                            // 466
			email = email.toLowerCase();                                                                                       // 467
		}                                                                                                                   // 468
		var usernameOrEmail = trimmedElementValueById('login-username-or-email');                                           // 469
		if (usernameOrEmail && Accounts.ui._options.forceEmailLowercase && Accounts.ui._options.forceUsernameLowercase) {   // 470
			usernameOrEmail = usernameOrEmail.toLowerCase();                                                                   // 471
		}                                                                                                                   // 472
                                                                                                                      // 473
		// notably not trimmed. a password could (?) start or end with a space                                              // 474
		var password = elementValueById('login-password');                                                                  // 475
		if (password && Accounts.ui._options.forcePasswordLowercase) {                                                      // 476
			password = password.toLowerCase();                                                                                 // 477
		}                                                                                                                   // 478
                                                                                                                      // 479
		var loginSelector;                                                                                                  // 480
		if (username !== null) {                                                                                            // 481
			if (!Accounts._loginButtons.validateUsername(username)){                                                           // 482
				return;                                                                                                           // 483
			} else {                                                                                                           // 484
				loginSelector = {                                                                                                 // 485
					username: username                                                                                               // 486
				};                                                                                                                // 487
			}                                                                                                                  // 488
		} else if (email !== null) {                                                                                        // 489
			if (!Accounts._loginButtons.validateEmail(email)){                                                                 // 490
				return;                                                                                                           // 491
			} else {                                                                                                           // 492
				loginSelector = {                                                                                                 // 493
					email: email                                                                                                     // 494
				};                                                                                                                // 495
			}                                                                                                                  // 496
		} else if (usernameOrEmail !== null) {                                                                              // 497
			// XXX not sure how we should validate this. but this seems good enough (for now),                                 // 498
			// since an email must have at least 3 characters anyways                                                          // 499
			if (!Accounts._loginButtons.validateUsername(usernameOrEmail)){                                                    // 500
				return;                                                                                                           // 501
			} else {                                                                                                           // 502
				loginSelector = usernameOrEmail;                                                                                  // 503
			}                                                                                                                  // 504
		} else {                                                                                                            // 505
			throw new Error("Unexpected -- no element to use as a login user selector");                                       // 506
		}                                                                                                                   // 507
                                                                                                                      // 508
		Meteor.loginWithPassword(loginSelector, password, function(error, result) {                                         // 509
			if (error) {                                                                                                       // 510
				if (error.reason == 'User not found'){                                                                            // 511
					loginButtonsSession.errorMessage(i18n('errorMessages.userNotFound'))                                             // 512
				} else if (error.reason == 'Incorrect password'){                                                                 // 513
					loginButtonsSession.errorMessage(i18n('errorMessages.incorrectPassword'))                                        // 514
				} else {                                                                                                          // 515
					loginButtonsSession.errorMessage(error.reason || "Unknown error");                                               // 516
				}                                                                                                                 // 517
			} else {                                                                                                           // 518
				loginButtonsSession.closeDropdown();                                                                              // 519
			}                                                                                                                  // 520
		});                                                                                                                 // 521
	};                                                                                                                   // 522
                                                                                                                      // 523
	var toggleDropdown = function() {                                                                                    // 524
		$("#login-dropdown-list").toggleClass("open");                                                                      // 525
	}                                                                                                                    // 526
                                                                                                                      // 527
	var focusInput = function() {                                                                                        // 528
		setTimeout(function() {                                                                                             // 529
			$("#login-dropdown-list input").first().focus();                                                                   // 530
		}, 0);                                                                                                              // 531
	};                                                                                                                   // 532
                                                                                                                      // 533
	var signup = function() {                                                                                            // 534
		loginButtonsSession.resetMessages();                                                                                // 535
                                                                                                                      // 536
		// to be passed to Accounts.createUser                                                                              // 537
		var options = {};                                                                                                   // 538
		if(typeof accountsUIBootstrap3.setCustomSignupOptions === 'function') {                                             // 539
			options = accountsUIBootstrap3.setCustomSignupOptions();                                                           // 540
			if (!(options instanceof Object)){ options = {}; }                                                                 // 541
		}                                                                                                                   // 542
                                                                                                                      // 543
		var username = trimmedElementValueById('login-username');                                                           // 544
		if (username && Accounts.ui._options.forceUsernameLowercase) {                                                      // 545
			username = username.toLowerCase();                                                                                 // 546
		}                                                                                                                   // 547
		if (username !== null) {                                                                                            // 548
			if (!Accounts._loginButtons.validateUsername(username)){                                                           // 549
				return;                                                                                                           // 550
			} else {                                                                                                           // 551
				options.username = username;                                                                                      // 552
			}                                                                                                                  // 553
		}                                                                                                                   // 554
                                                                                                                      // 555
		var email = trimmedElementValueById('login-email');                                                                 // 556
		if (email && Accounts.ui._options.forceEmailLowercase) {                                                            // 557
			email = email.toLowerCase();                                                                                       // 558
		}                                                                                                                   // 559
		if (email !== null) {                                                                                               // 560
			if (!Accounts._loginButtons.validateEmail(email)){                                                                 // 561
				return;                                                                                                           // 562
			} else {                                                                                                           // 563
				options.email = email;                                                                                            // 564
			}                                                                                                                  // 565
		}                                                                                                                   // 566
                                                                                                                      // 567
		// notably not trimmed. a password could (?) start or end with a space                                              // 568
		var password = elementValueById('login-password');                                                                  // 569
		if (password && Accounts.ui._options.forcePasswordLowercase) {                                                      // 570
			password = password.toLowerCase();                                                                                 // 571
		}                                                                                                                   // 572
		if (!Accounts._loginButtons.validatePassword(password)){                                                            // 573
			return;                                                                                                            // 574
		} else {                                                                                                            // 575
			options.password = password;                                                                                       // 576
		}                                                                                                                   // 577
                                                                                                                      // 578
		if (!matchPasswordAgainIfPresent()){                                                                                // 579
			return;                                                                                                            // 580
		}                                                                                                                   // 581
                                                                                                                      // 582
		// prepare the profile object                                                                                       // 583
		// it could have already been set through setCustomSignupOptions                                                    // 584
		if (!(options.profile instanceof Object)){                                                                          // 585
			options.profile = {};                                                                                              // 586
		}                                                                                                                   // 587
                                                                                                                      // 588
		// define a proxy function to allow extraSignupFields set error messages                                            // 589
		var errorFunction = function(errorMessage) {                                                                        // 590
			Accounts._loginButtonsSession.errorMessage(errorMessage);                                                          // 591
		};                                                                                                                  // 592
                                                                                                                      // 593
		var invalidExtraSignupFields = false;                                                                               // 594
		// parse extraSignupFields to populate account's profile data                                                       // 595
		_.each(Accounts.ui._options.extraSignupFields, function(field, index) {                                             // 596
						var value = null;                                                                                               // 597
						var elementIdPrefix = 'login-';                                                                                 // 598
                                                                                                                      // 599
						if (field.inputType === 'radio') {                                                                              // 600
							value = elementValueByIdForRadio(elementIdPrefix + field.fieldName, field.data);                               // 601
						} else if (field.inputType === 'checkbox') {                                                                    // 602
							value = elementValueByIdForCheckbox(elementIdPrefix + field.fieldName);                                        // 603
						} else {                                                                                                        // 604
							value = elementValueById(elementIdPrefix + field.fieldName);                                                   // 605
						}                                                                                                               // 606
                                                                                                                      // 607
			if (typeof field.validate === 'function') {                                                                        // 608
				if (field.validate(value, errorFunction)) {                                                                       // 609
					if (typeof field.saveToProfile !== 'undefined' && !field.saveToProfile){                                         // 610
						options[field.fieldName] = value;                                                                               // 611
					} else {                                                                                                         // 612
						options.profile[field.fieldName] = value;                                                                       // 613
					}                                                                                                                // 614
				} else {                                                                                                          // 615
					invalidExtraSignupFields = true;                                                                                 // 616
				}                                                                                                                 // 617
			} else {                                                                                                           // 618
				options.profile[field.fieldName] = value;                                                                         // 619
			}                                                                                                                  // 620
		});                                                                                                                 // 621
                                                                                                                      // 622
		if (invalidExtraSignupFields){                                                                                      // 623
			return;                                                                                                            // 624
		}                                                                                                                   // 625
                                                                                                                      // 626
		Accounts.createUser(options, function(error) {                                                                      // 627
			if (error) {                                                                                                       // 628
				if (error.reason == 'Signups forbidden'){                                                                         // 629
					loginButtonsSession.errorMessage(i18n('errorMessages.signupsForbidden'))                                         // 630
				} else {                                                                                                          // 631
					loginButtonsSession.errorMessage(error.reason || "Unknown error");                                               // 632
				}                                                                                                                 // 633
			} else {                                                                                                           // 634
				loginButtonsSession.closeDropdown();                                                                              // 635
			}                                                                                                                  // 636
		});                                                                                                                 // 637
	};                                                                                                                   // 638
                                                                                                                      // 639
	var forgotPassword = function() {                                                                                    // 640
		loginButtonsSession.resetMessages();                                                                                // 641
                                                                                                                      // 642
		var email = trimmedElementValueById("forgot-password-email");                                                       // 643
		if (email.indexOf('@') !== -1) {                                                                                    // 644
			Accounts.forgotPassword({                                                                                          // 645
				email: email                                                                                                      // 646
			}, function(error) {                                                                                               // 647
				if (error) {                                                                                                      // 648
					if (error.reason == 'User not found'){                                                                           // 649
						loginButtonsSession.errorMessage(i18n('errorMessages.userNotFound'))                                            // 650
					} else {                                                                                                         // 651
						loginButtonsSession.errorMessage(error.reason || "Unknown error");                                              // 652
					}                                                                                                                // 653
				} else {                                                                                                          // 654
					loginButtonsSession.infoMessage(i18n('infoMessages.emailSent'));                                                 // 655
				}                                                                                                                 // 656
			});                                                                                                                // 657
		} else {                                                                                                            // 658
			loginButtonsSession.errorMessage(i18n('forgotPasswordForm.invalidEmail'));                                         // 659
		}                                                                                                                   // 660
	};                                                                                                                   // 661
	var changePassword = function() {                                                                                    // 662
		loginButtonsSession.resetMessages();                                                                                // 663
		// notably not trimmed. a password could (?) start or end with a space                                              // 664
		var oldPassword = elementValueById('login-old-password');                                                           // 665
		// notably not trimmed. a password could (?) start or end with a space                                              // 666
		var password = elementValueById('login-password');                                                                  // 667
                                                                                                                      // 668
		if (password == oldPassword) {                                                                                      // 669
			loginButtonsSession.errorMessage(i18n('errorMessages.newPasswordSameAsOld'));                                      // 670
			return;                                                                                                            // 671
		}                                                                                                                   // 672
                                                                                                                      // 673
		if (!Accounts._loginButtons.validatePassword(password)){                                                            // 674
			return;                                                                                                            // 675
		}                                                                                                                   // 676
                                                                                                                      // 677
		if (!matchPasswordAgainIfPresent()){                                                                                // 678
			return;                                                                                                            // 679
		}                                                                                                                   // 680
                                                                                                                      // 681
		Accounts.changePassword(oldPassword, password, function(error) {                                                    // 682
			if (error) {                                                                                                       // 683
				if (error.reason == 'Incorrect password'){                                                                        // 684
					loginButtonsSession.errorMessage(i18n('errorMessages.incorrectPassword'))                                        // 685
				} else {                                                                                                          // 686
					loginButtonsSession.errorMessage(error.reason || "Unknown error");                                               // 687
				}                                                                                                                 // 688
			} else {                                                                                                           // 689
				loginButtonsSession.infoMessage(i18n('infoMessages.passwordChanged'));                                            // 690
                                                                                                                      // 691
				// wait 3 seconds, then expire the msg                                                                            // 692
				Meteor.setTimeout(function() {                                                                                    // 693
					loginButtonsSession.resetMessages();                                                                             // 694
				}, 3000);                                                                                                         // 695
			}                                                                                                                  // 696
		});                                                                                                                 // 697
	};                                                                                                                   // 698
                                                                                                                      // 699
	var matchPasswordAgainIfPresent = function() {                                                                       // 700
		// notably not trimmed. a password could (?) start or end with a space                                              // 701
		var passwordAgain = elementValueById('login-password-again');                                                       // 702
		if (passwordAgain !== null) {                                                                                       // 703
			// notably not trimmed. a password could (?) start or end with a space                                             // 704
			var password = elementValueById('login-password');                                                                 // 705
			if (password !== passwordAgain) {                                                                                  // 706
				loginButtonsSession.errorMessage(i18n('errorMessages.passwordsDontMatch'));                                       // 707
				return false;                                                                                                     // 708
			}                                                                                                                  // 709
		}                                                                                                                   // 710
		return true;                                                                                                        // 711
	};                                                                                                                   // 712
})();                                                                                                                 // 713
                                                                                                                      // 714
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/ian_accounts-ui-bootstrap-3/login_buttons_dialogs.js                                                      //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
(function() {                                                                                                         // 1
	// for convenience                                                                                                   // 2
	var loginButtonsSession = Accounts._loginButtonsSession;                                                             // 3
                                                                                                                      // 4
                                                                                                                      // 5
	//                                                                                                                   // 6
	// populate the session so that the appropriate dialogs are                                                          // 7
	// displayed by reading variables set by accounts-urls, which parses                                                 // 8
	// special URLs. since accounts-ui depends on accounts-urls, we are                                                  // 9
	// guaranteed to have these set at this point.                                                                       // 10
	//                                                                                                                   // 11
                                                                                                                      // 12
	if (Accounts._resetPasswordToken) {                                                                                  // 13
		loginButtonsSession.set('resetPasswordToken', Accounts._resetPasswordToken);                                        // 14
	}                                                                                                                    // 15
                                                                                                                      // 16
	if (Accounts._enrollAccountToken) {                                                                                  // 17
		loginButtonsSession.set('enrollAccountToken', Accounts._enrollAccountToken);                                        // 18
	}                                                                                                                    // 19
                                                                                                                      // 20
	// Needs to be in Meteor.startup because of a package loading order                                                  // 21
	// issue. We can't be sure that accounts-password is loaded earlier                                                  // 22
	// than accounts-ui so Accounts.verifyEmail might not be defined.                                                    // 23
	Meteor.startup(function() {                                                                                          // 24
		if (Accounts._verifyEmailToken) {                                                                                   // 25
			Accounts.verifyEmail(Accounts._verifyEmailToken, function(error) {                                                 // 26
				Accounts._enableAutoLogin();                                                                                      // 27
				if (!error){                                                                                                      // 28
					loginButtonsSession.set('justVerifiedEmail', true);                                                              // 29
				}                                                                                                                 // 30
				// XXX show something if there was an error.                                                                      // 31
			});                                                                                                                // 32
		}                                                                                                                   // 33
	});                                                                                                                  // 34
                                                                                                                      // 35
	//                                                                                                                   // 36
	// resetPasswordDialog template                                                                                      // 37
	//                                                                                                                   // 38
                                                                                                                      // 39
	Template._resetPasswordDialog.events({                                                                               // 40
		'click #login-buttons-reset-password-button': function(event) {                                                     // 41
			event.stopPropagation();                                                                                           // 42
			resetPassword();                                                                                                   // 43
		},                                                                                                                  // 44
		'keypress #reset-password-new-password': function(event) {                                                          // 45
			if (event.keyCode === 13){                                                                                         // 46
				resetPassword();                                                                                                  // 47
			}                                                                                                                  // 48
		},                                                                                                                  // 49
		'click #login-buttons-cancel-reset-password': function(event) {                                                     // 50
			event.stopPropagation();                                                                                           // 51
			loginButtonsSession.set('resetPasswordToken', null);                                                               // 52
			Accounts._enableAutoLogin();                                                                                       // 53
			$('#login-buttons-reset-password-modal').modal("hide");                                                            // 54
		},                                                                                                                  // 55
		'click #login-buttons-dismiss-reset-password-success': function(event) {                                            // 56
			event.stopPropagation();                                                                                           // 57
			$('#login-buttons-reset-password-modal-success').modal("hide");                                                    // 58
		}                                                                                                                   // 59
	});                                                                                                                  // 60
                                                                                                                      // 61
	var resetPassword = function() {                                                                                     // 62
		loginButtonsSession.resetMessages();                                                                                // 63
		var newPassword = document.getElementById('reset-password-new-password').value;                                     // 64
		var passwordAgain= document.getElementById('reset-password-new-password-again').value;                              // 65
		if (!Accounts._loginButtons.validatePassword(newPassword,passwordAgain)){                                           // 66
			return;                                                                                                            // 67
		}                                                                                                                   // 68
                                                                                                                      // 69
		Accounts.resetPassword(                                                                                             // 70
			loginButtonsSession.get('resetPasswordToken'), newPassword,                                                        // 71
			function(error) {                                                                                                  // 72
				if (error) {                                                                                                      // 73
					loginButtonsSession.errorMessage(error.reason || "Unknown error");                                               // 74
				} else {                                                                                                          // 75
					$('#login-buttons-reset-password-modal').modal("hide");                                                          // 76
					$('#login-buttons-reset-password-modal-success').modal();                                                        // 77
					loginButtonsSession.infoMessage(i18n('infoMessages.passwordChanged'));                                           // 78
					loginButtonsSession.set('resetPasswordToken', null);                                                             // 79
					Accounts._enableAutoLogin();                                                                                     // 80
				}                                                                                                                 // 81
			});                                                                                                                // 82
	};                                                                                                                   // 83
                                                                                                                      // 84
	Template._resetPasswordDialog.helpers({                                                                              // 85
		inResetPasswordFlow: function() {                                                                                   // 86
			return loginButtonsSession.get('resetPasswordToken');                                                              // 87
		}                                                                                                                   // 88
	});                                                                                                                  // 89
                                                                                                                      // 90
	Template._resetPasswordDialog.onRendered(function() {                                                                // 91
		var $modal = $(this.find('#login-buttons-reset-password-modal'));                                                   // 92
		if (!_.isFunction($modal.modal)) {                                                                                  // 93
			console.error("You have to add a Bootstrap package, i.e. meteor add twbs:bootstrap");                              // 94
		} else {                                                                                                            // 95
			$modal.modal();                                                                                                    // 96
		}                                                                                                                   // 97
	});                                                                                                                  // 98
                                                                                                                      // 99
	//                                                                                                                   // 100
	// enrollAccountDialog template                                                                                      // 101
	//                                                                                                                   // 102
                                                                                                                      // 103
	Template._enrollAccountDialog.events({                                                                               // 104
		'click #login-buttons-enroll-account-button': function() {                                                          // 105
			enrollAccount();                                                                                                   // 106
		},                                                                                                                  // 107
		'keypress #enroll-account-password': function(event) {                                                              // 108
			if (event.keyCode === 13){                                                                                         // 109
				enrollAccount();                                                                                                  // 110
			}                                                                                                                  // 111
		},                                                                                                                  // 112
		'click #login-buttons-cancel-enroll-account-button': function() {                                                   // 113
			loginButtonsSession.set('enrollAccountToken', null);                                                               // 114
			Accounts._enableAutoLogin();                                                                                       // 115
			$modal.modal("hide");                                                                                              // 116
		}                                                                                                                   // 117
	});                                                                                                                  // 118
                                                                                                                      // 119
	var enrollAccount = function() {                                                                                     // 120
		loginButtonsSession.resetMessages();                                                                                // 121
		var password = document.getElementById('enroll-account-password').value;                                            // 122
		var passwordAgain= document.getElementById('enroll-account-password-again').value;                                  // 123
		if (!Accounts._loginButtons.validatePassword(password,passwordAgain)){                                              // 124
			return;                                                                                                            // 125
		}                                                                                                                   // 126
                                                                                                                      // 127
		Accounts.resetPassword(                                                                                             // 128
			loginButtonsSession.get('enrollAccountToken'), password,                                                           // 129
			function(error) {                                                                                                  // 130
				if (error) {                                                                                                      // 131
					loginButtonsSession.errorMessage(error.reason || "Unknown error");                                               // 132
				} else {                                                                                                          // 133
					loginButtonsSession.set('enrollAccountToken', null);                                                             // 134
					Accounts._enableAutoLogin();                                                                                     // 135
					$modal.modal("hide");                                                                                            // 136
				}                                                                                                                 // 137
			});                                                                                                                // 138
	};                                                                                                                   // 139
                                                                                                                      // 140
	Template._enrollAccountDialog.helpers({                                                                              // 141
		inEnrollAccountFlow: function() {                                                                                   // 142
			return loginButtonsSession.get('enrollAccountToken');                                                              // 143
		}                                                                                                                   // 144
	});                                                                                                                  // 145
                                                                                                                      // 146
	Template._enrollAccountDialog.onRendered(function() {                                                                // 147
		$modal = $(this.find('#login-buttons-enroll-account-modal'));                                                       // 148
		if (!_.isFunction($modal.modal)) {                                                                                  // 149
			console.error("You have to add a Bootstrap package, i.e. meteor add twbs:bootstrap");                              // 150
		} else {                                                                                                            // 151
			$modal.modal();                                                                                                    // 152
		}                                                                                                                   // 153
	});                                                                                                                  // 154
                                                                                                                      // 155
	//                                                                                                                   // 156
	// justVerifiedEmailDialog template                                                                                  // 157
	//                                                                                                                   // 158
                                                                                                                      // 159
	Template._justVerifiedEmailDialog.events({                                                                           // 160
		'click #just-verified-dismiss-button': function() {                                                                 // 161
			loginButtonsSession.set('justVerifiedEmail', false);                                                               // 162
		}                                                                                                                   // 163
	});                                                                                                                  // 164
                                                                                                                      // 165
	Template._justVerifiedEmailDialog.helpers({                                                                          // 166
		visible: function() {                                                                                               // 167
			if (loginButtonsSession.get('justVerifiedEmail')) {                                                                // 168
				setTimeout(function() {                                                                                           // 169
					$('#login-buttons-email-address-verified-modal').modal()                                                         // 170
				}, 500)                                                                                                           // 171
			}                                                                                                                  // 172
			return loginButtonsSession.get('justVerifiedEmail');                                                               // 173
		}                                                                                                                   // 174
	});                                                                                                                  // 175
                                                                                                                      // 176
                                                                                                                      // 177
	//                                                                                                                   // 178
	// loginButtonsMessagesDialog template                                                                               // 179
	//                                                                                                                   // 180
                                                                                                                      // 181
	var messagesDialogVisible = function() {                                                                             // 182
		var hasMessage = loginButtonsSession.get('infoMessage') || loginButtonsSession.get('errorMessage');                 // 183
		return !Accounts._loginButtons.dropdown() && hasMessage;                                                            // 184
	}                                                                                                                    // 185
                                                                                                                      // 186
                                                                                                                      // 187
	Template._loginButtonsMessagesDialog.onRendered(function() {                                                         // 188
		var self = this;                                                                                                    // 189
                                                                                                                      // 190
		self.autorun(function() {                                                                                           // 191
			if (messagesDialogVisible()) {                                                                                     // 192
				var $modal = $(self.find('#login-buttons-message-dialog'));                                                       // 193
				if (!_.isFunction($modal.modal)) {                                                                                // 194
					console.error("You have to add a Bootstrap package, i.e. meteor add twbs:bootstrap");                            // 195
				} else {                                                                                                          // 196
					$modal.modal();                                                                                                  // 197
				}                                                                                                                 // 198
			}                                                                                                                  // 199
		});                                                                                                                 // 200
	});                                                                                                                  // 201
                                                                                                                      // 202
	Template._loginButtonsMessagesDialog.events({                                                                        // 203
		'click #messages-dialog-dismiss-button': function() {                                                               // 204
			loginButtonsSession.resetMessages();                                                                               // 205
		}                                                                                                                   // 206
	});                                                                                                                  // 207
                                                                                                                      // 208
	Template._loginButtonsMessagesDialog.helpers({                                                                       // 209
		visible: function() { return messagesDialogVisible(); }                                                             // 210
	});                                                                                                                  // 211
                                                                                                                      // 212
                                                                                                                      // 213
	//                                                                                                                   // 214
	// configureLoginServiceDialog template                                                                              // 215
	//                                                                                                                   // 216
                                                                                                                      // 217
	Template._configureLoginServiceDialog.events({                                                                       // 218
		'click .configure-login-service-dismiss-button': function(event) {                                                  // 219
			event.stopPropagation();                                                                                           // 220
			loginButtonsSession.set('configureLoginServiceDialogVisible', false);                                              // 221
			$('#configure-login-service-dialog-modal').modal('hide');                                                          // 222
		},                                                                                                                  // 223
		'click #configure-login-service-dialog-save-configuration': function() {                                            // 224
			if (loginButtonsSession.get('configureLoginServiceDialogVisible') &&                                               // 225
				!loginButtonsSession.get('configureLoginServiceDialogSaveDisabled')) {                                            // 226
				// Prepare the configuration document for this login service                                                      // 227
				var serviceName = loginButtonsSession.get('configureLoginServiceDialogServiceName');                              // 228
				var configuration = {                                                                                             // 229
					service: serviceName                                                                                             // 230
				};                                                                                                                // 231
				_.each(configurationFields(), function(field) {                                                                   // 232
					configuration[field.property] = document.getElementById(                                                         // 233
						'configure-login-service-dialog-' + field.property).value                                                       // 234
						.replace(/^\s*|\s*$/g, ""); // trim;                                                                            // 235
				});                                                                                                               // 236
                                                                                                                      // 237
				configuration.loginStyle =                                                                                        // 238
				$('#configure-login-service-dialog input[name="loginStyle"]:checked')                                             // 239
				.val();                                                                                                           // 240
                                                                                                                      // 241
				// Configure this login service                                                                                   // 242
				Meteor.call("configureLoginService", configuration, function(error, result) {                                     // 243
					if (error){                                                                                                      // 244
						Meteor._debug("Error configuring login service " + serviceName, error);                                         // 245
					} else {                                                                                                         // 246
						loginButtonsSession.set('configureLoginServiceDialogVisible', false);                                           // 247
					}                                                                                                                // 248
					$('#configure-login-service-dialog-modal').modal('hide');                                                        // 249
				});                                                                                                               // 250
			}                                                                                                                  // 251
		},                                                                                                                  // 252
		// IE8 doesn't support the 'input' event, so we'll run this on the keyup as                                         // 253
		// well. (Keeping the 'input' event means that this also fires when you use                                         // 254
		// the mouse to change the contents of the field, eg 'Cut' menu item.)                                              // 255
		'input, keyup input': function(event) {                                                                             // 256
			// if the event fired on one of the configuration input fields,                                                    // 257
			// check whether we should enable the 'save configuration' button                                                  // 258
			if (event.target.id.indexOf('configure-login-service-dialog') === 0){                                              // 259
				updateSaveDisabled();                                                                                             // 260
			}                                                                                                                  // 261
		}                                                                                                                   // 262
	});                                                                                                                  // 263
                                                                                                                      // 264
	// check whether the 'save configuration' button should be enabled.                                                  // 265
	// this is a really strange way to implement this and a Forms                                                        // 266
	// Abstraction would make all of this reactive, and simpler.                                                         // 267
	var updateSaveDisabled = function() {                                                                                // 268
		var anyFieldEmpty = _.any(configurationFields(), function(field) {                                                  // 269
			return document.getElementById(                                                                                    // 270
				'configure-login-service-dialog-' + field.property).value === '';                                                 // 271
		});                                                                                                                 // 272
                                                                                                                      // 273
		loginButtonsSession.set('configureLoginServiceDialogSaveDisabled', anyFieldEmpty);                                  // 274
	};                                                                                                                   // 275
                                                                                                                      // 276
	// Returns the appropriate template for this login service.  This                                                    // 277
	// template should be defined in the service's package                                                               // 278
	var configureLoginServiceDialogTemplateForService = function() {                                                     // 279
		var serviceName = loginButtonsSession.get('configureLoginServiceDialogServiceName');                                // 280
		return Template['configureLoginServiceDialogFor' + capitalize(serviceName)];                                        // 281
	};                                                                                                                   // 282
                                                                                                                      // 283
	var configurationFields = function() {                                                                               // 284
		var template = configureLoginServiceDialogTemplateForService();                                                     // 285
		return template.fields();                                                                                           // 286
	};                                                                                                                   // 287
                                                                                                                      // 288
	Template._configureLoginServiceDialog.helpers({                                                                      // 289
		configurationFields: function() {                                                                                   // 290
			return configurationFields();                                                                                      // 291
		},                                                                                                                  // 292
                                                                                                                      // 293
		visible: function() {                                                                                               // 294
			return loginButtonsSession.get('configureLoginServiceDialogVisible');                                              // 295
		},                                                                                                                  // 296
                                                                                                                      // 297
		configurationSteps: function() {                                                                                    // 298
			// renders the appropriate template                                                                                // 299
			return configureLoginServiceDialogTemplateForService();                                                            // 300
		},                                                                                                                  // 301
                                                                                                                      // 302
		saveDisabled: function() {                                                                                          // 303
			return loginButtonsSession.get('configureLoginServiceDialogSaveDisabled');                                         // 304
		}                                                                                                                   // 305
	});                                                                                                                  // 306
                                                                                                                      // 307
                                                                                                                      // 308
	;                                                                                                                    // 309
                                                                                                                      // 310
                                                                                                                      // 311
                                                                                                                      // 312
	// XXX from http://epeli.github.com/underscore.string/lib/underscore.string.js                                       // 313
	var capitalize = function(str) {                                                                                     // 314
		str = str == null ? '' : String(str);                                                                               // 315
		return str.charAt(0).toUpperCase() + str.slice(1);                                                                  // 316
	};                                                                                                                   // 317
                                                                                                                      // 318
})();                                                                                                                 // 319
                                                                                                                      // 320
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);


/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package['ian:accounts-ui-bootstrap-3'] = {}, {
  accountsUIBootstrap3: accountsUIBootstrap3
});

})();
