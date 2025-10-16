"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.RegisterRoutes = RegisterRoutes;
/* tslint:disable */
/* eslint-disable */
// WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
const runtime_1 = require("@tsoa/runtime");
// WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
const auth_controller_1 = require("./controllers/auth.controller");
// WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
const user_controller_1 = require("./controllers/user.controller");
const auth_middleware_1 = require("./middleware/auth.middleware");
// @ts-ignore - no great way to install types from subpackage
const promiseAny = require('promise.any');
// WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
const models = {
    "LoginResponse": {
        "dataType": "refObject",
        "properties": {
            "user": { "dataType": "nestedObjectLiteral", "nestedProperties": { "created_at": { "dataType": "datetime", "required": true }, "active": { "dataType": "boolean", "required": true }, "email_verified": { "dataType": "boolean", "required": true }, "contact_email": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }], "required": true }, "login_email": { "dataType": "string", "required": true }, "last_name": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }], "required": true }, "first_name": { "dataType": "string", "required": true }, "id": { "dataType": "double", "required": true } }, "required": true },
            "token": { "dataType": "string", "required": true },
            "session_id": { "dataType": "double", "required": true },
            "expires_at": { "dataType": "datetime", "required": true },
            "token_type": { "dataType": "string", "required": true },
        },
        "additionalProperties": false,
    },
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    "LoginRequest": {
        "dataType": "refObject",
        "properties": {
            "login_email": { "dataType": "string", "required": true },
            "password": { "dataType": "string", "required": true },
        },
        "additionalProperties": false,
    },
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    "LogoutResponse": {
        "dataType": "refObject",
        "properties": {
            "message": { "dataType": "string", "required": true },
            "sessions_ended": { "dataType": "double", "required": true },
        },
        "additionalProperties": false,
    },
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    "TokenValidationResponse": {
        "dataType": "refObject",
        "properties": {
            "valid": { "dataType": "boolean", "required": true },
            "user": { "dataType": "nestedObjectLiteral", "nestedProperties": { "active": { "dataType": "boolean", "required": true }, "email_verified": { "dataType": "boolean", "required": true }, "login_email": { "dataType": "string", "required": true }, "last_name": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }], "required": true }, "first_name": { "dataType": "string", "required": true }, "id": { "dataType": "double", "required": true } } },
            "session": { "dataType": "nestedObjectLiteral", "nestedProperties": { "expires_at": { "dataType": "datetime", "required": true }, "created_at": { "dataType": "datetime", "required": true }, "ip_address": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }], "required": true }, "id": { "dataType": "double", "required": true } } },
        },
        "additionalProperties": false,
    },
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    "SessionInfo": {
        "dataType": "refObject",
        "properties": {
            "id": { "dataType": "double", "required": true },
            "ip_address": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }], "required": true },
            "user_agent": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }], "required": true },
            "created_at": { "dataType": "datetime", "required": true },
            "expires_at": { "dataType": "datetime", "required": true },
            "is_current": { "dataType": "boolean" },
        },
        "additionalProperties": false,
    },
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    "UserResponse": {
        "dataType": "refObject",
        "properties": {
            "id": { "dataType": "double", "required": true },
            "first_name": { "dataType": "string", "required": true },
            "last_name": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }], "required": true },
            "login_email": { "dataType": "string", "required": true },
            "contact_email": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }], "required": true },
            "contact_phone_no": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }], "required": true },
            "adk_session_id": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }], "required": true },
            "email_verified": { "dataType": "boolean", "required": true },
            "active": { "dataType": "boolean", "required": true },
            "created_at": { "dataType": "datetime", "required": true },
            "updated_at": { "dataType": "datetime", "required": true },
        },
        "additionalProperties": false,
    },
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    "UserCreateRequest": {
        "dataType": "refObject",
        "properties": {
            "first_name": { "dataType": "string", "required": true },
            "last_name": { "dataType": "string", "required": true },
            "login_email": { "dataType": "string", "required": true },
            "password": { "dataType": "string", "required": true },
            "contact_email": { "dataType": "string" },
            "contact_phone_no": { "dataType": "string" },
            "adk_session_id": { "dataType": "string" },
        },
        "additionalProperties": false,
    },
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    "UserUpdateRequest": {
        "dataType": "refObject",
        "properties": {
            "first_name": { "dataType": "string" },
            "last_name": { "dataType": "string" },
            "login_email": { "dataType": "string" },
            "password": { "dataType": "string" },
            "contact_email": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }] },
            "contact_phone_no": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }] },
            "adk_session_id": { "dataType": "union", "subSchemas": [{ "dataType": "string" }, { "dataType": "enum", "enums": [null] }] },
            "email_verified": { "dataType": "boolean" },
            "active": { "dataType": "boolean" },
        },
        "additionalProperties": false,
    },
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
};
const validationService = new runtime_1.ValidationService(models);
// WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
function RegisterRoutes(app) {
    // ###########################################################################################################
    //  NOTE: If you do not see routes for all of your controllers in this file, then you might not have informed tsoa of where to look
    //      Please look into the "controllerPathGlobs" config option described in the readme: https://github.com/lukeautry/tsoa
    // ###########################################################################################################
    app.post('/auth/login', ...((0, runtime_1.fetchMiddlewares)(auth_controller_1.AuthController)), ...((0, runtime_1.fetchMiddlewares)(auth_controller_1.AuthController.prototype.login)), function AuthController_login(request, response, next) {
        const args = {
            requestBody: { "in": "body", "name": "requestBody", "required": true, "ref": "LoginRequest" },
            userAgent: { "in": "header", "name": "user-agent", "dataType": "string" },
            ipAddress: { "in": "header", "name": "x-forwarded-for", "dataType": "string" },
        };
        // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
        let validatedArgs = [];
        try {
            validatedArgs = getValidatedArgs(args, request, response);
            const controller = new auth_controller_1.AuthController();
            const promise = controller.login.apply(controller, validatedArgs);
            promiseHandler(controller, promise, response, 200, next);
        }
        catch (err) {
            return next(err);
        }
    });
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    app.post('/auth/logout', authenticateMiddleware([{ "BearerAuth": [] }]), ...((0, runtime_1.fetchMiddlewares)(auth_controller_1.AuthController)), ...((0, runtime_1.fetchMiddlewares)(auth_controller_1.AuthController.prototype.logout)), function AuthController_logout(request, response, next) {
        const args = {
            request: { "in": "request", "name": "request", "required": true, "dataType": "object" },
        };
        // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
        let validatedArgs = [];
        try {
            validatedArgs = getValidatedArgs(args, request, response);
            const controller = new auth_controller_1.AuthController();
            const promise = controller.logout.apply(controller, validatedArgs);
            promiseHandler(controller, promise, response, 200, next);
        }
        catch (err) {
            return next(err);
        }
    });
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    app.post('/auth/logout-all', authenticateMiddleware([{ "BearerAuth": [] }]), ...((0, runtime_1.fetchMiddlewares)(auth_controller_1.AuthController)), ...((0, runtime_1.fetchMiddlewares)(auth_controller_1.AuthController.prototype.logoutAll)), function AuthController_logoutAll(request, response, next) {
        const args = {
            request: { "in": "request", "name": "request", "required": true, "dataType": "object" },
        };
        // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
        let validatedArgs = [];
        try {
            validatedArgs = getValidatedArgs(args, request, response);
            const controller = new auth_controller_1.AuthController();
            const promise = controller.logoutAll.apply(controller, validatedArgs);
            promiseHandler(controller, promise, response, 200, next);
        }
        catch (err) {
            return next(err);
        }
    });
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    app.get('/auth/validate', authenticateMiddleware([{ "BearerAuth": [] }]), ...((0, runtime_1.fetchMiddlewares)(auth_controller_1.AuthController)), ...((0, runtime_1.fetchMiddlewares)(auth_controller_1.AuthController.prototype.validateToken)), function AuthController_validateToken(request, response, next) {
        const args = {
            request: { "in": "request", "name": "request", "required": true, "dataType": "object" },
        };
        // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
        let validatedArgs = [];
        try {
            validatedArgs = getValidatedArgs(args, request, response);
            const controller = new auth_controller_1.AuthController();
            const promise = controller.validateToken.apply(controller, validatedArgs);
            promiseHandler(controller, promise, response, 200, next);
        }
        catch (err) {
            return next(err);
        }
    });
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    app.get('/auth/sessions', authenticateMiddleware([{ "BearerAuth": [] }]), ...((0, runtime_1.fetchMiddlewares)(auth_controller_1.AuthController)), ...((0, runtime_1.fetchMiddlewares)(auth_controller_1.AuthController.prototype.getSessions)), function AuthController_getSessions(request, response, next) {
        const args = {
            request: { "in": "request", "name": "request", "required": true, "dataType": "object" },
        };
        // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
        let validatedArgs = [];
        try {
            validatedArgs = getValidatedArgs(args, request, response);
            const controller = new auth_controller_1.AuthController();
            const promise = controller.getSessions.apply(controller, validatedArgs);
            promiseHandler(controller, promise, response, 200, next);
        }
        catch (err) {
            return next(err);
        }
    });
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    app.get('/users', ...((0, runtime_1.fetchMiddlewares)(user_controller_1.UsersController)), ...((0, runtime_1.fetchMiddlewares)(user_controller_1.UsersController.prototype.getAllUsers)), function UsersController_getAllUsers(request, response, next) {
        const args = {};
        // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
        let validatedArgs = [];
        try {
            validatedArgs = getValidatedArgs(args, request, response);
            const controller = new user_controller_1.UsersController();
            const promise = controller.getAllUsers.apply(controller, validatedArgs);
            promiseHandler(controller, promise, response, 200, next);
        }
        catch (err) {
            return next(err);
        }
    });
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    app.get('/users/:id', ...((0, runtime_1.fetchMiddlewares)(user_controller_1.UsersController)), ...((0, runtime_1.fetchMiddlewares)(user_controller_1.UsersController.prototype.getUserById)), function UsersController_getUserById(request, response, next) {
        const args = {
            id: { "in": "path", "name": "id", "required": true, "dataType": "double" },
        };
        // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
        let validatedArgs = [];
        try {
            validatedArgs = getValidatedArgs(args, request, response);
            const controller = new user_controller_1.UsersController();
            const promise = controller.getUserById.apply(controller, validatedArgs);
            promiseHandler(controller, promise, response, 200, next);
        }
        catch (err) {
            return next(err);
        }
    });
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    app.post('/users', ...((0, runtime_1.fetchMiddlewares)(user_controller_1.UsersController)), ...((0, runtime_1.fetchMiddlewares)(user_controller_1.UsersController.prototype.createUser)), function UsersController_createUser(request, response, next) {
        const args = {
            requestBody: { "in": "body", "name": "requestBody", "required": true, "ref": "UserCreateRequest" },
        };
        // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
        let validatedArgs = [];
        try {
            validatedArgs = getValidatedArgs(args, request, response);
            const controller = new user_controller_1.UsersController();
            const promise = controller.createUser.apply(controller, validatedArgs);
            promiseHandler(controller, promise, response, 201, next);
        }
        catch (err) {
            return next(err);
        }
    });
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    app.put('/users/:id', ...((0, runtime_1.fetchMiddlewares)(user_controller_1.UsersController)), ...((0, runtime_1.fetchMiddlewares)(user_controller_1.UsersController.prototype.updateUser)), function UsersController_updateUser(request, response, next) {
        const args = {
            id: { "in": "path", "name": "id", "required": true, "dataType": "double" },
            requestBody: { "in": "body", "name": "requestBody", "required": true, "ref": "UserUpdateRequest" },
        };
        // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
        let validatedArgs = [];
        try {
            validatedArgs = getValidatedArgs(args, request, response);
            const controller = new user_controller_1.UsersController();
            const promise = controller.updateUser.apply(controller, validatedArgs);
            promiseHandler(controller, promise, response, 200, next);
        }
        catch (err) {
            return next(err);
        }
    });
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    app.delete('/users/:id', ...((0, runtime_1.fetchMiddlewares)(user_controller_1.UsersController)), ...((0, runtime_1.fetchMiddlewares)(user_controller_1.UsersController.prototype.deleteUser)), function UsersController_deleteUser(request, response, next) {
        const args = {
            id: { "in": "path", "name": "id", "required": true, "dataType": "double" },
        };
        // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
        let validatedArgs = [];
        try {
            validatedArgs = getValidatedArgs(args, request, response);
            const controller = new user_controller_1.UsersController();
            const promise = controller.deleteUser.apply(controller, validatedArgs);
            promiseHandler(controller, promise, response, 204, next);
        }
        catch (err) {
            return next(err);
        }
    });
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    function authenticateMiddleware(security = []) {
        return async function runAuthenticationMiddleware(request, _response, next) {
            // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
            // keep track of failed auth attempts so we can hand back the most
            // recent one.  This behavior was previously existing so preserving it
            // here
            const failedAttempts = [];
            const pushAndRethrow = (error) => {
                failedAttempts.push(error);
                throw error;
            };
            const secMethodOrPromises = [];
            for (const secMethod of security) {
                if (Object.keys(secMethod).length > 1) {
                    const secMethodAndPromises = [];
                    for (const name in secMethod) {
                        secMethodAndPromises.push((0, auth_middleware_1.expressAuthentication)(request, name, secMethod[name])
                            .catch(pushAndRethrow));
                    }
                    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
                    secMethodOrPromises.push(Promise.all(secMethodAndPromises)
                        .then(users => { return users[0]; }));
                }
                else {
                    for (const name in secMethod) {
                        secMethodOrPromises.push((0, auth_middleware_1.expressAuthentication)(request, name, secMethod[name])
                            .catch(pushAndRethrow));
                    }
                }
            }
            // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
            try {
                request['user'] = await promiseAny.call(Promise, secMethodOrPromises);
                next();
            }
            catch (err) {
                // Show most recent error as response
                const error = failedAttempts.pop();
                error.status = error.status || 401;
                next(error);
            }
            // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
        };
    }
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    function isController(object) {
        return 'getHeaders' in object && 'getStatus' in object && 'setStatus' in object;
    }
    function promiseHandler(controllerObj, promise, response, successStatus, next) {
        return Promise.resolve(promise)
            .then((data) => {
            let statusCode = successStatus;
            let headers;
            if (isController(controllerObj)) {
                headers = controllerObj.getHeaders();
                statusCode = controllerObj.getStatus() || statusCode;
            }
            // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
            returnHandler(response, statusCode, data, headers);
        })
            .catch((error) => next(error));
    }
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    function returnHandler(response, statusCode, data, headers = {}) {
        if (response.headersSent) {
            return;
        }
        Object.keys(headers).forEach((name) => {
            response.set(name, headers[name]);
        });
        if (data && typeof data.pipe === 'function' && data.readable && typeof data._read === 'function') {
            response.status(statusCode || 200);
            data.pipe(response);
        }
        else if (data !== null && data !== undefined) {
            response.status(statusCode || 200).json(data);
        }
        else {
            response.status(statusCode || 204).end();
        }
    }
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    function responder(response) {
        return function (status, data, headers) {
            returnHandler(response, status, data, headers);
        };
    }
    ;
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
    function getValidatedArgs(args, request, response) {
        const fieldErrors = {};
        const values = Object.keys(args).map((key) => {
            const name = args[key].name;
            switch (args[key].in) {
                case 'request':
                    return request;
                case 'query':
                    return validationService.ValidateParam(args[key], request.query[name], name, fieldErrors, undefined, { "noImplicitAdditionalProperties": "throw-on-extras" });
                case 'queries':
                    return validationService.ValidateParam(args[key], request.query, name, fieldErrors, undefined, { "noImplicitAdditionalProperties": "throw-on-extras" });
                case 'path':
                    return validationService.ValidateParam(args[key], request.params[name], name, fieldErrors, undefined, { "noImplicitAdditionalProperties": "throw-on-extras" });
                case 'header':
                    return validationService.ValidateParam(args[key], request.header(name), name, fieldErrors, undefined, { "noImplicitAdditionalProperties": "throw-on-extras" });
                case 'body':
                    return validationService.ValidateParam(args[key], request.body, name, fieldErrors, undefined, { "noImplicitAdditionalProperties": "throw-on-extras" });
                case 'body-prop':
                    return validationService.ValidateParam(args[key], request.body[name], name, fieldErrors, 'body.', { "noImplicitAdditionalProperties": "throw-on-extras" });
                case 'formData':
                    if (args[key].dataType === 'file') {
                        return validationService.ValidateParam(args[key], request.file, name, fieldErrors, undefined, { "noImplicitAdditionalProperties": "throw-on-extras" });
                    }
                    else if (args[key].dataType === 'array' && args[key].array.dataType === 'file') {
                        return validationService.ValidateParam(args[key], request.files, name, fieldErrors, undefined, { "noImplicitAdditionalProperties": "throw-on-extras" });
                    }
                    else {
                        return validationService.ValidateParam(args[key], request.body[name], name, fieldErrors, undefined, { "noImplicitAdditionalProperties": "throw-on-extras" });
                    }
                case 'res':
                    return responder(response);
            }
        });
        if (Object.keys(fieldErrors).length > 0) {
            throw new runtime_1.ValidateError(fieldErrors, '');
        }
        return values;
    }
    // WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
}
// WARNING: This file was auto-generated with tsoa. Please do not modify it. Re-run tsoa to re-generate this file: https://github.com/lukeautry/tsoa
