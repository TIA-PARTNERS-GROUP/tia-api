# TIA Partner API
API documentation for the TIA Partner platform.

## Version: 1.0

### Security
**BearerAuth**  

| apiKey | *API Key* |
| ------ | --------- |
| Description | Type "Bearer" followed by a space and your JWT token. |
| Name | Authorization |
| In | header |

**Schemes:** http

---
### /auth/login

#### POST
##### Summary

User Login

##### Description

Authenticates a user with email and password, creating a new session and returning a JWT token.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| login | body | Login Credentials | Yes | [ports.LoginInput](#portslogininput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successful login, returns user data and token | [ports.LoginResponse](#portsloginresponse) |
| 400 | Invalid request body or validation error | object |
| 401 | Invalid email/password or account deactivated | object |
| 500 | Internal server error | object |

### /auth/logout

#### POST
##### Summary

User Logout

##### Description

Invalidates the current user session (token).

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Successfully logged out (No Content) |  |
| 401 | Unauthorized or missing token | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /auth/me

#### GET
##### Summary

Get Current User

##### Description

Retrieves the profile of the currently authenticated user based on the provided JWT token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User profile retrieved successfully | [ports.UserResponse](#portsuserresponse) |
| 401 | Unauthorized or token missing/invalid | object |
| 500 | Internal server error or invalid context | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users

#### GET
##### Summary

List All Users

##### Description

Retrieves a list of all user profiles. Requires authentication.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of users | [ [ports.UserResponse](#portsuserresponse) ] |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Register New User

##### Description

Registers a new user account. Does not require prior authentication.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| user | body | User registration details (Name, Email, Password) | Yes | [ports.UserCreationSchema](#portsusercreationschema) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | User created successfully | [ports.UserResponse](#portsuserresponse) |
| 400 | Invalid request body or validation failed | object |
| 409 | ErrUserAlreadyExists | object |
| 500 | Internal server error | object |

---
### /businesses

#### GET
##### Summary

List All Businesses

##### Description

Retrieves a list of all business profiles with optional filtering.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| operator_user_id | query | Filter by Operator User ID | No | integer |
| business_type | query | Filter by business type | No | string |
| business_category | query | Filter by business category | No | string |
| business_phase | query | Filter by business phase | No | string |
| search | query | Search by name or description | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of businesses | [ [ports.BusinessResponse](#portsbusinessresponse) ] |
| 400 | Invalid query parameters | object |
| 500 | Internal server error | object |

#### POST
##### Summary

Create New Business

##### Description

Creates a new business profile, restricted to the authenticated user (OperatorUserID must match Auth UserID).

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| business | body | Business creation details | Yes | [ports.CreateBusinessInput](#portscreatebusinessinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Business created successfully | [ports.BusinessResponse](#portsbusinessresponse) |
| 400 | Invalid request body, validation failed | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: OperatorUserID mismatch | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /businesses/{id}

#### GET
##### Summary

Get Business by ID

##### Description

Retrieves a business profile by its unique ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Business ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Business retrieved successfully | [ports.BusinessResponse](#portsbusinessresponse) |
| 400 | Invalid business ID | object |
| 404 | Business not found | object |
| 500 | Internal server error | object |

#### PUT
##### Summary

Update Business Profile

##### Description

Updates an existing business profile. Only the designated OperatorUserID can perform this action.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Business ID | Yes | integer |
| business | body | Fields to update | Yes | [ports.UpdateBusinessInput](#portsupdatebusinessinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Business updated successfully | [ports.BusinessResponse](#portsbusinessresponse) |
| 400 | Invalid request body, validation failed | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Not the operator user | object |
| 404 | Business not found | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Delete Business Profile

##### Description

Deletes a business profile. Only the designated OperatorUserID can perform this action.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Business ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Business deleted successfully (No Content) |  |
| 400 | Invalid business ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Not the operator user or business is in use | object |
| 404 | Business not found | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /businesses/{id}/tags

#### GET
##### Summary

Get All Tags for a Business

##### Description

Retrieves all tags associated with a specific business.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Business ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of business tags | [ports.BusinessTagsResponse](#portsbusinesstagsresponse) |
| 400 | Invalid business ID | object |
| 500 | Internal server error | object |

#### POST
##### Summary

Add Tag to Business

##### Description

Creates a new tag (e.g., 'client', 'service') and associates it with a specific business.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Business ID to associate the tag with | Yes | integer |
| tag | body | Tag details (TagType, Description) | Yes | [ports.CreateBusinessTagInput](#portscreatebusinesstaginput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Tag created and associated successfully | [ports.BusinessTagResponse](#portsbusinesstagresponse) |
| 400 | Invalid business ID, request body, or validation failed | object |
| 401 | Unauthorized | object |
| 409 | ErrBusinessTagAlreadyExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /tags/{id}

#### DELETE
##### Summary

Delete Business Tag

##### Description

Deletes a specific business tag entry by its unique Tag ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Unique Business Tag ID (NOT the Business ID) | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Tag deleted successfully (No Content) |  |
| 400 | Invalid tag ID | object |
| 401 | Unauthorized | object |
| 404 | ErrBusinessTagNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /businesses/{id}/connections

#### GET
##### Summary

List Connections for a Business

##### Description

Retrieves a list of all connections (initiated and received) associated with a specific business ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Business ID | Yes | integer |
| status | query | Filter by connection status (pending, active, rejected, inactive) | No | string |
| type | query | Filter by connection type (Partnership, Client, Supplier, etc.) | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of connections | [ports.BusinessConnectionsResponse](#portsbusinessconnectionsresponse) |
| 400 | Invalid business ID or query parameters | object |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /connections

#### POST
##### Summary

Initiate Business Connection Request

##### Description

Creates a new connection request between two businesses. The initiating user is taken from the auth context.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| connection | body | Connection request details (InitiatingBusinessID, ReceivingBusinessID, ConnectionType) | Yes | [ports.CreateBusinessConnectionInput](#portscreatebusinessconnectioninput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Connection request created successfully | [ports.BusinessConnectionResponse](#portsbusinessconnectionresponse) |
| 400 | Invalid request body, validation failed, or ErrCannotConnectToSelf | object |
| 401 | Unauthorized | object |
| 409 | ErrBusinessConnectionAlreadyExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /connections/{id}

#### GET
##### Summary

Get Business Connection by ID

##### Description

Retrieves a specific connection record by its unique ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Connection ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Connection retrieved successfully | [ports.BusinessConnectionResponse](#portsbusinessconnectionresponse) |
| 400 | Invalid connection ID | object |
| 401 | Unauthorized | object |
| 404 | ErrBusinessConnectionNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### PUT
##### Summary

Update Business Connection Details

##### Description

Updates modifiable fields of an existing connection (e.g., Notes, Type). This is typically restricted to the initiating user.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Connection ID | Yes | integer |
| connection | body | Fields to update (e.g., ConnectionType, Notes) | Yes | [ports.UpdateBusinessConnectionInput](#portsupdatebusinessconnectioninput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Connection updated successfully | [ports.BusinessConnectionResponse](#portsbusinessconnectionresponse) |
| 400 | Invalid request body or connection ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not authorized to update) | object |
| 404 | ErrBusinessConnectionNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Delete Business Connection

##### Description

Deletes a specific connection record. Restricted to the initiating user or the receiving business's operator.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Connection ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Connection deleted successfully (No Content) |  |
| 400 | Invalid connection ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not authorized to delete) | object |
| 404 | ErrBusinessConnectionNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /connections/{id}/accept

#### PATCH
##### Summary

Accept Pending Connection

##### Description

Updates the status of a specific connection request to 'active'. Restricted to the receiving business's operator.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Connection ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Connection successfully accepted and set to active | [ports.BusinessConnectionResponse](#portsbusinessconnectionresponse) |
| 400 | Invalid connection ID or ErrConnectionNotPending | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the receiving business's operator) | object |
| 404 | ErrBusinessConnectionNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /connections/{id}/reject

#### PATCH
##### Summary

Reject Pending Connection

##### Description

Updates the status of a specific connection request to 'rejected'. Restricted to the receiving business's operator.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Connection ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Connection successfully rejected | [ports.BusinessConnectionResponse](#portsbusinessconnectionresponse) |
| 400 | Invalid connection ID or ErrConnectionNotPending | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the receiving business's operator) | object |
| 404 | ErrBusinessConnectionNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /businesses/{id}/tags

#### GET
##### Summary

Get All Tags for a Business

##### Description

Retrieves all tags associated with a specific business.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Business ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of business tags | [ports.BusinessTagsResponse](#portsbusinesstagsresponse) |
| 400 | Invalid business ID | object |
| 500 | Internal server error | object |

#### POST
##### Summary

Add Tag to Business

##### Description

Creates a new tag (e.g., 'client', 'service') and associates it with a specific business.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Business ID to associate the tag with | Yes | integer |
| tag | body | Tag details (TagType, Description) | Yes | [ports.CreateBusinessTagInput](#portscreatebusinesstaginput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Tag created and associated successfully | [ports.BusinessTagResponse](#portsbusinesstagresponse) |
| 400 | Invalid business ID, request body, or validation failed | object |
| 401 | Unauthorized | object |
| 409 | ErrBusinessTagAlreadyExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /tags/{id}

#### DELETE
##### Summary

Delete Business Tag

##### Description

Deletes a specific business tag entry by its unique Tag ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Unique Business Tag ID (NOT the Business ID) | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Tag deleted successfully (No Content) |  |
| 400 | Invalid tag ID | object |
| 401 | Unauthorized | object |
| 404 | ErrBusinessTagNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /daily-activities

#### GET
##### Summary

Get All Daily Activities

##### Description

Retrieves a list of all defined daily activities.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of daily activities | [ [models.DailyActivity](#modelsdailyactivity) ] |
| 500 | Internal server error | object |

#### POST
##### Summary

Create Daily Activity

##### Description

Creates a new daily activity definition (e.g., "30-minute meditation").

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| activity | body | Activity details (Name, Description) | Yes | [ports.CreateDailyActivityInput](#portscreatedailyactivityinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Activity created successfully | [models.DailyActivity](#modelsdailyactivity) |
| 400 | Invalid request body or validation failed | object |
| 401 | Unauthorized | object |
| 409 | ErrActivityNameExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /daily-activities/{id}

#### GET
##### Summary

Get Daily Activity by ID

##### Description

Retrieves a specific daily activity definition by its unique ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Activity ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Activity retrieved successfully | [models.DailyActivity](#modelsdailyactivity) |
| 400 | Invalid activity ID | object |
| 404 | ErrDailyActivityNotFound | object |
| 500 | Internal server error | object |

### /daily-activities/{id}/enrolments

#### GET
##### Summary

Get All Enrolments for Activity

##### Description

Retrieves a list of all users currently enrolled in a specified daily activity.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Daily Activity ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of user enrolments | [ [ports.ActivityEnrolmentResponse](#portsactivityenrolmentresponse) ] |
| 400 | Invalid activity ID | object |
| 500 | Internal server error | object |

#### POST
##### Summary

Enrol User in Daily Activity

##### Description

Enrols the authenticated user in a specified daily activity.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Daily Activity ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Enrolment successful | [ports.ActivityEnrolmentResponse](#portsactivityenrolmentresponse) |
| 400 | Invalid activity ID | object |
| 401 | Unauthorized | object |
| 404 | Activity or user not found | object |
| 409 | ErrAlreadyEnrolled | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Withdraw User from Daily Activity

##### Description

Withdraws the authenticated user from a specified daily activity.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Daily Activity ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Withdrawal successful (No Content) |  |
| 400 | Invalid activity ID | object |
| 401 | Unauthorized | object |
| 404 | ErrEnrolmentNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/enrolments

#### GET
##### Summary

Get All Enrolments for User

##### Description

Retrieves a list of all daily activities a specific user is currently enrolled in.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of activity enrolments for the user | [ [ports.UserEnrolmentResponse](#portsuserenrolmentresponse) ] |
| 400 | Invalid user ID | object |
| 500 | Internal server error | object |

---
### /daily-activities/{id}/enrolments

#### GET
##### Summary

Get All Enrolments for Activity

##### Description

Retrieves a list of all users currently enrolled in a specified daily activity.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Daily Activity ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of user enrolments | [ [ports.ActivityEnrolmentResponse](#portsactivityenrolmentresponse) ] |
| 400 | Invalid activity ID | object |
| 500 | Internal server error | object |

#### POST
##### Summary

Enrol User in Daily Activity

##### Description

Enrols the authenticated user in a specified daily activity.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Daily Activity ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Enrolment successful | [ports.ActivityEnrolmentResponse](#portsactivityenrolmentresponse) |
| 400 | Invalid activity ID | object |
| 401 | Unauthorized | object |
| 404 | Activity or user not found | object |
| 409 | ErrAlreadyEnrolled | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Withdraw User from Daily Activity

##### Description

Withdraws the authenticated user from a specified daily activity.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Daily Activity ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Withdrawal successful (No Content) |  |
| 400 | Invalid activity ID | object |
| 401 | Unauthorized | object |
| 404 | ErrEnrolmentNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/enrolments

#### GET
##### Summary

Get All Enrolments for User

##### Description

Retrieves a list of all daily activities a specific user is currently enrolled in.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of activity enrolments for the user | [ [ports.UserEnrolmentResponse](#portsuserenrolmentresponse) ] |
| 400 | Invalid user ID | object |
| 500 | Internal server error | object |

---
### /events

#### GET
##### Summary

Get All Events

##### Description

Retrieves a list of all events with optional filtering.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| user_id | query | Filter by User ID | No | integer |
| event_type | query | Filter by event type | No | string |
| start_date | query | Filter events created after this date (YYYY-MM-DD) | No | string |
| end_date | query | Filter events created before this date (YYYY-MM-DD) | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of events | [ [ports.EventResponse](#portseventresponse) ] |
| 400 | Invalid query parameters | object |
| 500 | Failed to retrieve events | object |

#### POST
##### Summary

Create New Event Record

##### Description

Creates a new internal system event record. The UserID is automatically injected from the authenticated context.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| event | body | Event details (EventType, Payload) | Yes | [ports.CreateEventInput](#portscreateeventinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Event created successfully | [ports.EventResponse](#portseventresponse) |
| 400 | Invalid input data or validation error | object |
| 401 | Unauthorized or missing authentication context | object |
| 500 | Failed to create event | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /events/{id}

#### GET
##### Summary

Get Event by ID

##### Description

Retrieves a specific event record by its unique ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Event ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Event retrieved successfully | [ports.EventResponse](#portseventresponse) |
| 400 | Invalid event ID | object |
| 404 | Event not found | object |
| 500 | Failed to retrieve event | object |

---
### /feedback

#### GET
##### Summary

Get All Feedback

##### Description

Retrieves a list of all submitted feedback. Requires authentication.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of feedback entries | [ [ports.FeedbackResponse](#portsfeedbackresponse) ] |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /feedback/{id}

#### GET
##### Summary

Get Feedback by ID

##### Description

Retrieves a specific feedback entry by its unique ID. Requires authentication.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Feedback ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Feedback entry retrieved successfully | [ports.FeedbackResponse](#portsfeedbackresponse) |
| 400 | Invalid feedback ID | object |
| 401 | Unauthorized | object |
| 404 | ErrFeedbackNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Delete Feedback

##### Description

Deletes a specific feedback entry by its unique ID. Requires authentication.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Feedback ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Feedback deleted successfully (No Content) |  |
| 400 | Invalid feedback ID | object |
| 401 | Unauthorized | object |
| 404 | ErrFeedbackNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /inferred-connections

#### POST
##### Summary

Create Inferred Connection Record

##### Description

Creates a new record for a potential connection inferred by a model. Intended for internal/system use.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| connection | body | Inferred connection details | Yes | [ports.CreateInferredConnectionInput](#portscreateinferredconnectioninput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Connection record created successfully | [ports.InferredConnectionResponse](#portsinferredconnectionresponse) |
| 400 | Invalid request body or validation error | object |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /inferred-connections/source/{entityType}/{entityID}

#### GET
##### Summary

Get Inferred Connections by Source Entity

##### Description

Retrieves all potential connections inferred from a specific source entity (e.g., a Project or Business).

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| entityType | path | Type of the source entity (e.g., business, project) | Yes | string |
| entityID | path | ID of the source entity | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of inferred connections | [ [ports.InferredConnectionResponse](#portsinferredconnectionresponse) ] |
| 400 | Invalid entity ID | object |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /l2e-responses

#### POST
##### Summary

Submit New L2E Response

##### Description

Submits a user's response or data payload for a Learn-to-Earn (L2E) module. The UserID is taken from the auth context.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| response | body | L2E response data payload | Yes | [ports.CreateL2EResponseInput](#portscreatel2eresponseinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Response recorded successfully | [ports.L2EResponseResponse](#portsl2eresponseresponse) |
| 400 | Invalid request body or validation error | object |
| 401 | Unauthorized or invalid context | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/l2e-responses

#### GET
##### Summary

Get All L2E Responses for User

##### Description

Retrieves all L2E responses submitted by a specific user.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of L2E responses | [ [ports.L2EResponseResponse](#portsl2eresponseresponse) ] |
| 400 | Invalid user ID | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /notifications

#### POST
##### Summary

Create Notification (Internal/System Use)

##### Description

Creates a new notification record. Requires authentication, typically used by system services.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| notification | body | Notification details (ReceiverUserID, Title, Message, Type) | Yes | [ports.CreateNotificationInput](#portscreatenotificationinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Notification created successfully | [ports.NotificationResponse](#portsnotificationresponse) |
| 400 | Invalid request body or validation error | object |
| 401 | Unauthorized | object |
| 404 | ErrReceiverNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/notifications

#### GET
##### Summary

Get Notifications for User

##### Description

Retrieves all notifications for the specified user. Requires authentication and self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID | Yes | integer |
| read | query | Filter by read status (true/false) | No | boolean |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of notifications | [ [ports.NotificationResponse](#portsnotificationresponse) ] |
| 400 | Invalid user ID or query parameter | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/notifications/read-all

#### PATCH
##### Summary

Mark All Notifications as Read

##### Description

Marks all unread notifications for the specified user as read. Requires authentication and self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Count of notifications marked as read | object |
| 400 | Invalid user ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/notifications/{notificationID}

#### DELETE
##### Summary

Delete Single Notification

##### Description

Deletes a specific notification for the user. Requires authentication and self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID | Yes | integer |
| notificationID | path | Notification ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Notification deleted successfully (No Content) |  |
| 400 | Invalid user/notification ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 404 | Notification not found for this user | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/notifications/{notificationID}/read

#### PATCH
##### Summary

Mark Single Notification as Read

##### Description

Marks a specific notification as read. Requires authentication and self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID | Yes | integer |
| notificationID | path | Notification ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Notification marked as read | [ports.NotificationResponse](#portsnotificationresponse) |
| 400 | Invalid user/notification ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 404 | Notification not found for this user | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /projects

#### GET
##### Summary

Get All Projects

##### Description

Retrieves a list of all project records.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of projects | [ [ports.ProjectResponse](#portsprojectresponse) ] |
| 401 | Unauthorized | object |
| 500 | Failed to retrieve projects | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Create New Project

##### Description

Creates a new project record. Requires authentication.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| project | body | Project creation details (Name, ManagedByUserID, ProjectStatus) | Yes | [ports.CreateProjectInput](#portscreateprojectinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Project created successfully | [ports.ProjectResponse](#portsprojectresponse) |
| 400 | Invalid request body or validation failed | object |
| 401 | Unauthorized | object |
| 404 | ErrManagerNotFound | object |
| 409 | ErrProjectNameExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}

#### GET
##### Summary

Get Project by ID

##### Description

Retrieves a specific project record by its unique ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Project retrieved successfully | [ports.ProjectResponse](#portsprojectresponse) |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 404 | ErrProjectNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### PUT
##### Summary

Update Project Details

##### Description

Updates an existing project record. Only the Project Manager can perform this action.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| project | body | Fields to update | Yes | [ports.UpdateProjectInput](#portsupdateprojectinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Project updated successfully | [ports.ProjectResponse](#portsprojectresponse) |
| 400 | Invalid project ID or request body | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Delete Project

##### Description

Deletes a project record and all related data (members, regions, skills). Only the Project Manager can perform this action.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Project deleted successfully (No Content) |  |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}/applicants

#### GET
##### Summary

Get Applicants for Project

##### Description

Retrieves a list of all users who have applied to a specific project. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of applicants | [ [ports.ProjectApplicantResponse](#portsprojectapplicantresponse) ] |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}/apply

#### POST
##### Summary

Apply to Project

##### Description

Submits an application for the authenticated user to join a project. The UserID is taken from the auth context.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Application submitted successfully (Created) |  |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 404 | ErrProjectNotFound or user not found | object |
| 409 | ErrAlreadyApplied | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Withdraw Application

##### Description

Withdraws the authenticated user's application from a project.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Application withdrawn successfully (No Content) |  |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 404 | ErrApplicationNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}/members

#### GET
##### Summary

Get All Project Members

##### Description

Retrieves a list of all members associated with a project. Accessible by any authenticated user.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of project members | [ports.ProjectMembersResponse](#portsprojectmembersresponse) |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Add Project Member

##### Description

Adds a user to a project with a specified role. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| member | body | Member details (UserID, Role) | Yes | [ports.AddProjectMemberInput](#portsaddprojectmemberinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Member added successfully | [ports.ProjectMemberResponse](#portsprojectmemberresponse) |
| 400 | Invalid project ID, request body, or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectNotFound or ErrUserNotFound | object |
| 409 | ErrProjectMemberAlreadyExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}/members/{userID}

#### GET
##### Summary

Get Specific Project Member

##### Description

Retrieves a specific project member record by Project ID and User ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| userID | path | User ID of the member | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Member retrieved successfully | [ports.ProjectMemberResponse](#portsprojectmemberresponse) |
| 400 | Invalid ID | object |
| 401 | Unauthorized | object |
| 404 | ErrProjectMemberNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### PUT
##### Summary

Update Project Member Role

##### Description

Updates the role of an existing project member. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| userID | path | User ID of the member | Yes | integer |
| role | body | New role for the member | Yes | [ports.UpdateProjectMemberRoleInput](#portsupdateprojectmemberroleinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Member role updated successfully | [ports.ProjectMemberResponse](#portsprojectmemberresponse) |
| 400 | Invalid ID, request body, or ErrInvalidRole | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectMemberNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Remove Project Member

##### Description

Removes a user from a project. Allowed for the **Project Manager** (to remove anyone) or the **User** (to remove self).

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| userID | path | User ID of the member to remove | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Member removed successfully (No Content) |  |
| 400 | Invalid ID or ErrCannotRemoveManager | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the manager and not the user) | object |
| 404 | ErrProjectMemberNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}/regions

#### GET
##### Summary

Get Regions for Project

##### Description

Retrieves all geographical regions associated with a specific project.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of regions | [ [ports.ProjectRegionResponse](#portsprojectregionresponse) ] |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Add Region to Project

##### Description

Associates a geographical region (identified by its short code/ID) with a project. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| region | body | Region details (RegionID) | Yes | [ports.AddProjectRegionInput](#portsaddprojectregioninput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Region associated successfully | [ports.ProjectRegionResponse](#portsprojectregionresponse) |
| 400 | Invalid project ID, request body, or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 409 | ErrRegionAlreadyAdded | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}/regions/{regionID}

#### DELETE
##### Summary

Remove Region from Project

##### Description

Dissociates a geographical region from a project. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| regionID | path | Region ID (e.g., USA, AUS) | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Region removed successfully (No Content) |  |
| 400 | Invalid project ID or missing regionID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectRegionNotFound or ErrProjectNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}/skills

#### GET
##### Summary

Get Skills Required by Project

##### Description

Retrieves all skill requirements for a specific project.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of required skills | [ports.ProjectSkillsResponse](#portsprojectskillsresponse) |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Add Skill Requirement to Project

##### Description

Associates a specific skill (by Skill ID) with a project and sets its importance level. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| skill | body | Skill details (SkillID, Importance) | Yes | [ports.CreateProjectSkillInput](#portscreateprojectskillinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Skill requirement added successfully | [ports.ProjectSkillResponse](#portsprojectskillresponse) |
| 400 | Invalid project ID, skill ID, request body, or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectNotFound or ErrSkillNotFound | object |
| 409 | ErrProjectSkillAlreadyExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}/skills/{skillID}

#### PUT
##### Summary

Update Project Skill Importance

##### Description

Updates the importance level (required, preferred, optional) for an existing project skill. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| skillID | path | Skill ID | Yes | integer |
| update | body | New importance level | Yes | [ports.UpdateProjectSkillInput](#portsupdateprojectskillinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Skill importance updated successfully | [ports.ProjectSkillResponse](#portsprojectskillresponse) |
| 400 | Invalid ID, request body, or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectSkillNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Remove Skill Requirement from Project

##### Description

Removes a skill requirement association from a project. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| skillID | path | Skill ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Skill requirement removed successfully (No Content) |  |
| 400 | Invalid ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectSkillNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /projects/{id}/applicants

#### GET
##### Summary

Get Applicants for Project

##### Description

Retrieves a list of all users who have applied to a specific project. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of applicants | [ [ports.ProjectApplicantResponse](#portsprojectapplicantresponse) ] |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}/apply

#### POST
##### Summary

Apply to Project

##### Description

Submits an application for the authenticated user to join a project. The UserID is taken from the auth context.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Application submitted successfully (Created) |  |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 404 | ErrProjectNotFound or user not found | object |
| 409 | ErrAlreadyApplied | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Withdraw Application

##### Description

Withdraws the authenticated user's application from a project.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Application withdrawn successfully (No Content) |  |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 404 | ErrApplicationNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/applications

#### GET
##### Summary

Get Applications for User

##### Description

Retrieves a list of all projects the specified user has applied to. Requires authentication and self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of user applications | [ [ports.UserApplicationResponse](#portsuserapplicationresponse) ] |
| 400 | Invalid user ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /projects/{id}/members

#### GET
##### Summary

Get All Project Members

##### Description

Retrieves a list of all members associated with a project. Accessible by any authenticated user.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of project members | [ports.ProjectMembersResponse](#portsprojectmembersresponse) |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Add Project Member

##### Description

Adds a user to a project with a specified role. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| member | body | Member details (UserID, Role) | Yes | [ports.AddProjectMemberInput](#portsaddprojectmemberinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Member added successfully | [ports.ProjectMemberResponse](#portsprojectmemberresponse) |
| 400 | Invalid project ID, request body, or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectNotFound or ErrUserNotFound | object |
| 409 | ErrProjectMemberAlreadyExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}/members/{userID}

#### GET
##### Summary

Get Specific Project Member

##### Description

Retrieves a specific project member record by Project ID and User ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| userID | path | User ID of the member | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Member retrieved successfully | [ports.ProjectMemberResponse](#portsprojectmemberresponse) |
| 400 | Invalid ID | object |
| 401 | Unauthorized | object |
| 404 | ErrProjectMemberNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### PUT
##### Summary

Update Project Member Role

##### Description

Updates the role of an existing project member. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| userID | path | User ID of the member | Yes | integer |
| role | body | New role for the member | Yes | [ports.UpdateProjectMemberRoleInput](#portsupdateprojectmemberroleinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Member role updated successfully | [ports.ProjectMemberResponse](#portsprojectmemberresponse) |
| 400 | Invalid ID, request body, or ErrInvalidRole | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectMemberNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Remove Project Member

##### Description

Removes a user from a project. Allowed for the **Project Manager** (to remove anyone) or the **User** (to remove self).

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| userID | path | User ID of the member to remove | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Member removed successfully (No Content) |  |
| 400 | Invalid ID or ErrCannotRemoveManager | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the manager and not the user) | object |
| 404 | ErrProjectMemberNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/project-memberships

#### GET
##### Summary

Get Projects by User

##### Description

Retrieves a list of all projects the specified user is a member of. Requires self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID | Yes | integer |
| role | query | Filter by member role (manager, contributor, reviewer) | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of project memberships | [ports.ProjectMembersResponse](#portsprojectmembersresponse) |
| 400 | Invalid user ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /projects/{id}/regions

#### GET
##### Summary

Get Regions for Project

##### Description

Retrieves all geographical regions associated with a specific project.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of regions | [ [ports.ProjectRegionResponse](#portsprojectregionresponse) ] |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Add Region to Project

##### Description

Associates a geographical region (identified by its short code/ID) with a project. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| region | body | Region details (RegionID) | Yes | [ports.AddProjectRegionInput](#portsaddprojectregioninput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Region associated successfully | [ports.ProjectRegionResponse](#portsprojectregionresponse) |
| 400 | Invalid project ID, request body, or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 409 | ErrRegionAlreadyAdded | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}/regions/{regionID}

#### DELETE
##### Summary

Remove Region from Project

##### Description

Dissociates a geographical region from a project. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| regionID | path | Region ID (e.g., USA, AUS) | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Region removed successfully (No Content) |  |
| 400 | Invalid project ID or missing regionID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectRegionNotFound or ErrProjectNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /projects/{id}/skills

#### GET
##### Summary

Get Skills Required by Project

##### Description

Retrieves all skill requirements for a specific project.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of required skills | [ports.ProjectSkillsResponse](#portsprojectskillsresponse) |
| 400 | Invalid project ID | object |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Add Skill Requirement to Project

##### Description

Associates a specific skill (by Skill ID) with a project and sets its importance level. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| skill | body | Skill details (SkillID, Importance) | Yes | [ports.CreateProjectSkillInput](#portscreateprojectskillinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Skill requirement added successfully | [ports.ProjectSkillResponse](#portsprojectskillresponse) |
| 400 | Invalid project ID, skill ID, request body, or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectNotFound or ErrSkillNotFound | object |
| 409 | ErrProjectSkillAlreadyExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /projects/{id}/skills/{skillID}

#### PUT
##### Summary

Update Project Skill Importance

##### Description

Updates the importance level (required, preferred, optional) for an existing project skill. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| skillID | path | Skill ID | Yes | integer |
| update | body | New importance level | Yes | [ports.UpdateProjectSkillInput](#portsupdateprojectskillinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Skill importance updated successfully | [ports.ProjectSkillResponse](#portsprojectskillresponse) |
| 400 | Invalid ID, request body, or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectSkillNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Remove Skill Requirement from Project

##### Description

Removes a skill requirement association from a project. Only accessible by the **Project Manager**.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Project ID | Yes | integer |
| skillID | path | Skill ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Skill requirement removed successfully (No Content) |  |
| 400 | Invalid ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the project manager) | object |
| 404 | ErrProjectSkillNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /skills

#### GET
##### Summary

Get All Skills with Filters

##### Description

Retrieves a list of all skills, with options to filter by category, activity status, or search term.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| category | query | Filter by skill category | No | string |
| active | query | Filter by active status (true/false) | No | boolean |
| search | query | Search by name, category, or description | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of skills | [ [ports.SkillResponse](#portsskillresponse) ] |
| 400 | Invalid query parameters | object |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Create New Skill

##### Description

Creates a new global skill record.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| skill | body | Skill creation details (Name, Category) | Yes | [ports.CreateSkillInput](#portscreateskillinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Skill created successfully | [ports.SkillResponse](#portsskillresponse) |
| 400 | Invalid request body or validation failed | object |
| 401 | Unauthorized | object |
| 409 | ErrSkillNameExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /skills/{id}

#### GET
##### Summary

Get Skill by ID

##### Description

Retrieves a specific skill record by its unique ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Skill ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Skill retrieved successfully | [ports.SkillResponse](#portsskillresponse) |
| 400 | Invalid skill ID | object |
| 401 | Unauthorized | object |
| 404 | ErrSkillNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### PUT
##### Summary

Update Skill

##### Description

Updates the details of an existing skill (e.g., category, name).

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Skill ID | Yes | integer |
| update | body | Fields to update (Category, Name, Active) | Yes | [ports.UpdateSkillInput](#portsupdateskillinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Skill updated successfully | [ports.SkillResponse](#portsskillresponse) |
| 400 | Invalid skill ID or request body | object |
| 401 | Unauthorized | object |
| 404 | ErrSkillNotFound | object |
| 409 | ErrSkillNameExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Delete Skill

##### Description

Deletes a specific skill record. Fails if the skill is currently in use by a user or project.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Skill ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Skill deleted successfully (No Content) |  |
| 400 | Invalid skill ID | object |
| 401 | Unauthorized | object |
| 404 | ErrSkillNotFound | object |
| 409 | ErrSkillInUse | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /skills/{id}/toggle-status

#### PATCH
##### Summary

Toggle Skill Status

##### Description

Toggles the active status of a skill (Active -> Inactive, or vice versa).

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Skill ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Skill status toggled successfully | [ports.SkillResponse](#portsskillresponse) |
| 400 | Invalid skill ID | object |
| 401 | Unauthorized | object |
| 404 | ErrSkillNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/skills

#### GET
##### Summary

Get User Skills

##### Description

Retrieves all skills and proficiency levels associated with the user. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of user skills | [ports.UserSkillsResponse](#portsuserskillsresponse) |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Add User Skill

##### Description

Adds a new skill and its proficiency level to the authenticated user's profile. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |
| skill | body | Skill details (SkillID, ProficiencyLevel) | Yes | [ports.CreateUserSkillInput](#portscreateuserskillinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Skill added successfully | [ports.UserSkillResponse](#portsuserskillresponse) |
| 400 | Invalid request body or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 404 | ErrUserNotFound or ErrSkillNotFound | object |
| 409 | ErrUserSkillAlreadyExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/skills/{skillID}

#### PUT
##### Summary

Update User Skill Proficiency

##### Description

Updates the proficiency level for an existing skill associated with the user. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |
| skillID | path | Skill ID | Yes | integer |
| update | body | New proficiency level | Yes | [ports.UpdateUserSkillInput](#portsupdateuserskillinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Proficiency updated successfully | [ports.UserSkillResponse](#portsuserskillresponse) |
| 400 | Invalid ID, request body, or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 404 | ErrUserSkillNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Remove User Skill

##### Description

Removes a skill association from the user's profile. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |
| skillID | path | Skill ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Skill removed successfully (No Content) |  |
| 400 | Invalid ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 404 | ErrUserSkillNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /publications

#### GET
##### Summary

Get All Publications

##### Description

Retrieves a list of all publication records.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of publications | [ [ports.PublicationResponse](#portspublicationresponse) ] |
| 401 | Unauthorized | object |
| 500 | Failed to retrieve publications | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Create New Publication

##### Description

Creates a new publication (post, article, case study, etc.). The UserID in the body must match the authenticated user.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| publication | body | Publication creation details (Title, UserID, Content, Type) | Yes | [ports.CreatePublicationInput](#portscreatepublicationinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Publication created successfully | [ports.PublicationResponse](#portspublicationresponse) |
| 400 | Invalid request body or validation failed | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Cannot create publication for another user | object |
| 409 | ErrPublicationSlugExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /publications/id/{id}

#### GET
##### Summary

Get Publication by ID

##### Description

Retrieves a specific publication record by its unique ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Publication ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Publication retrieved successfully | [ports.PublicationResponse](#portspublicationresponse) |
| 400 | Invalid publication ID | object |
| 401 | Unauthorized | object |
| 404 | ErrPublicationNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /publications/slug/{slug}

#### GET
##### Summary

Get Publication by Slug

##### Description

Retrieves a specific publication record by its unique URL slug.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| slug | path | Publication URL Slug | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Publication retrieved successfully | [ports.PublicationResponse](#portspublicationresponse) |
| 400 | Missing slug | object |
| 401 | Unauthorized | object |
| 404 | ErrPublicationNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /publications/{id}

#### PUT
##### Summary

Update Publication

##### Description

Updates an existing publication record. Only the Author (UserID) can perform this action.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Publication ID | Yes | integer |
| update | body | Fields to update (Title, Content, Published, etc.) | Yes | [ports.UpdatePublicationInput](#portsupdatepublicationinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Publication updated successfully | [ports.PublicationResponse](#portspublicationresponse) |
| 400 | Invalid publication ID or request body | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the author) | object |
| 404 | ErrPublicationNotFound | object |
| 409 | ErrPublicationSlugExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Delete Publication

##### Description

Deletes a publication record. Only the Author (UserID) can perform this action.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Publication ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Publication deleted successfully (No Content) |  |
| 400 | Invalid publication ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the author) | object |
| 404 | ErrPublicationNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /subscriptions

#### POST
##### Summary

Create New Subscription Plan

##### Description

Creates a new recurring subscription plan definition. Requires authentication (implies admin/privileged access).

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| plan | body | Subscription plan details (Name, Price, ValidDays/ValidMonths) | Yes | [ports.CreateSubscriptionInput](#portscreatesubscriptioninput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Subscription plan created successfully | [ports.SubscriptionResponse](#portssubscriptionresponse) |
| 400 | Invalid request body or validation failed | object |
| 401 | Unauthorized | object |
| 409 | ErrSubscriptionNameExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /subscriptions/subscribe

#### POST
##### Summary

Subscribe User to a Plan

##### Description

Creates a new UserSubscription record for the authenticated user, starting their access to a plan. Enforces self-subscription.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| subscription | body | Subscription details (UserID and SubscriptionID) | Yes | [ports.UserSubscribeInput](#portsusersubscribeinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | User subscribed successfully | [ports.UserSubscriptionResponse](#portsusersubscriptionresponse) |
| 400 | Invalid request body or validation failed | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Cannot subscribe for another user | object |
| 404 | ErrSubscriptionNotFound or ErrUserNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /subscriptions/{id}

#### GET
##### Summary

Get Subscription Plan by ID

##### Description

Retrieves a specific subscription plan definition by its unique ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Subscription Plan ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Subscription plan retrieved successfully | [ports.SubscriptionResponse](#portssubscriptionresponse) |
| 400 | Invalid subscription ID | object |
| 401 | Unauthorized | object |
| 404 | ErrSubscriptionNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/subscriptions

#### GET
##### Summary

Get Active Subscriptions for User

##### Description

Retrieves all currently active user subscription records (those not yet expired). Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID (must match authenticated user) | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of active user subscriptions | [ [ports.UserSubscriptionResponse](#portsusersubscriptionresponse) ] |
| 400 | Invalid user ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Cannot view another user's subscriptions | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/subscriptions/{userSubscriptionID}

#### DELETE
##### Summary

Cancel User Subscription

##### Description

Cancels a specific user subscription record by deleting it from the database. Enforces self-management and ownership of the record.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID (must match authenticated user) | Yes | integer |
| userSubscriptionID | path | User Subscription Record ID to cancel | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Subscription cancelled successfully (No Content) |  |
| 400 | Invalid user or subscription ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: You are not the owner of this record | object |
| 404 | ErrUserSubscriptionNotFound | object |
| 500 | Internal error during cancellation | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /users

#### GET
##### Summary

List All Users

##### Description

Retrieves a list of all user profiles. Requires authentication.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of users | [ [ports.UserResponse](#portsuserresponse) ] |
| 401 | Unauthorized | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Register New User

##### Description

Registers a new user account. Does not require prior authentication.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| user | body | User registration details (Name, Email, Password) | Yes | [ports.UserCreationSchema](#portsusercreationschema) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | User created successfully | [ports.UserResponse](#portsuserresponse) |
| 400 | Invalid request body or validation failed | object |
| 409 | ErrUserAlreadyExists | object |
| 500 | Internal server error | object |

### /users/{id}

#### GET
##### Summary

Get User by ID

##### Description

Retrieves a user's profile by their unique ID. Requires authentication.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User retrieved successfully | [ports.UserResponse](#portsuserresponse) |
| 400 | Invalid user ID | object |
| 401 | Unauthorized | object |
| 404 | ErrUserNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### PUT
##### Summary

Update User Profile

##### Description

Updates the authenticated user's profile information. Requires self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID (must match authenticated user) | Yes | integer |
| update | body | Fields to update (e.g., FirstName, ContactEmail) | Yes | [ports.UserUpdateSchema](#portsuserupdateschema) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Profile updated successfully | [ports.UserResponse](#portsuserresponse) |
| 400 | Invalid user ID or request body | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Cannot update another user's profile | object |
| 404 | ErrUserNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Delete User Account

##### Description

Deletes the authenticated user's account. Requires self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID (must match authenticated user) | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Account deleted successfully (No Content) |  |
| 400 | Invalid user ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Cannot delete another user's profile | object |
| 404 | ErrUserNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/applications

#### GET
##### Summary

Get Applications for User

##### Description

Retrieves a list of all projects the specified user has applied to. Requires authentication and self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of user applications | [ [ports.UserApplicationResponse](#portsuserapplicationresponse) ] |
| 400 | Invalid user ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/config

#### PUT
##### Summary

Set or Update User Configuration

##### Description

Creates a new configuration entry for a user, or updates an existing one for the given config_type. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |
| config | body | Configuration Data | Yes | [ports.SetUserConfigInput](#portssetuserconfiginput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Configuration successfully saved/updated | [ports.UserConfigResponse](#portsuserconfigresponse) |
| 400 | Invalid request body or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Cannot modify another user's config | object |
| 500 | Database error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/config/{configType}

#### GET
##### Summary

Get User Configuration by Type

##### Description

Retrieves a specific configuration entry for a user. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |
| configType | path | Configuration Type (e.g., user_preferences) | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Configuration successfully retrieved | [ports.UserConfigResponse](#portsuserconfigresponse) |
| 400 | Invalid user ID or missing configType | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Cannot view another user's config | object |
| 404 | ErrUserConfigNotFound | object |
| 500 | Database error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Delete User Configuration by Type

##### Description

Deletes a specific configuration entry for a user. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |
| configType | path | Configuration Type (e.g., notification_settings) | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Configuration successfully deleted (No Content) |  |
| 400 | Invalid user ID or missing configType | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Cannot delete another user's config | object |
| 404 | ErrUserConfigNotFound | object |
| 500 | Database error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/notifications

#### GET
##### Summary

Get Notifications for User

##### Description

Retrieves all notifications for the specified user. Requires authentication and self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID | Yes | integer |
| read | query | Filter by read status (true/false) | No | boolean |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of notifications | [ [ports.NotificationResponse](#portsnotificationresponse) ] |
| 400 | Invalid user ID or query parameter | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/notifications/read-all

#### PATCH
##### Summary

Mark All Notifications as Read

##### Description

Marks all unread notifications for the specified user as read. Requires authentication and self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Count of notifications marked as read | object |
| 400 | Invalid user ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/notifications/{notificationID}

#### DELETE
##### Summary

Delete Single Notification

##### Description

Deletes a specific notification for the user. Requires authentication and self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID | Yes | integer |
| notificationID | path | Notification ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Notification deleted successfully (No Content) |  |
| 400 | Invalid user/notification ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 404 | Notification not found for this user | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/notifications/{notificationID}/read

#### PATCH
##### Summary

Mark Single Notification as Read

##### Description

Marks a specific notification as read. Requires authentication and self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID | Yes | integer |
| notificationID | path | Notification ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Notification marked as read | [ports.NotificationResponse](#portsnotificationresponse) |
| 400 | Invalid user/notification ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 404 | Notification not found for this user | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/project-memberships

#### GET
##### Summary

Get Projects by User

##### Description

Retrieves a list of all projects the specified user is a member of. Requires self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID | Yes | integer |
| role | query | Filter by member role (manager, contributor, reviewer) | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of project memberships | [ports.ProjectMembersResponse](#portsprojectmembersresponse) |
| 400 | Invalid user ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/skills

#### GET
##### Summary

Get User Skills

##### Description

Retrieves all skills and proficiency levels associated with the user. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of user skills | [ports.UserSkillsResponse](#portsuserskillsresponse) |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Add User Skill

##### Description

Adds a new skill and its proficiency level to the authenticated user's profile. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |
| skill | body | Skill details (SkillID, ProficiencyLevel) | Yes | [ports.CreateUserSkillInput](#portscreateuserskillinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Skill added successfully | [ports.UserSkillResponse](#portsuserskillresponse) |
| 400 | Invalid request body or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 404 | ErrUserNotFound or ErrSkillNotFound | object |
| 409 | ErrUserSkillAlreadyExists | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/skills/{skillID}

#### PUT
##### Summary

Update User Skill Proficiency

##### Description

Updates the proficiency level for an existing skill associated with the user. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |
| skillID | path | Skill ID | Yes | integer |
| update | body | New proficiency level | Yes | [ports.UpdateUserSkillInput](#portsupdateuserskillinput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Proficiency updated successfully | [ports.UserSkillResponse](#portsuserskillresponse) |
| 400 | Invalid ID, request body, or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 404 | ErrUserSkillNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Remove User Skill

##### Description

Removes a skill association from the user's profile. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |
| skillID | path | Skill ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Skill removed successfully (No Content) |  |
| 400 | Invalid ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden (Not the target user) | object |
| 404 | ErrUserSkillNotFound | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/subscriptions

#### GET
##### Summary

Get Active Subscriptions for User

##### Description

Retrieves all currently active user subscription records (those not yet expired). Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID (must match authenticated user) | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of active user subscriptions | [ [ports.UserSubscriptionResponse](#portsusersubscriptionresponse) ] |
| 400 | Invalid user ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Cannot view another user's subscriptions | object |
| 500 | Internal server error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/subscriptions/{userSubscriptionID}

#### DELETE
##### Summary

Cancel User Subscription

##### Description

Cancels a specific user subscription record by deleting it from the database. Enforces self-management and ownership of the record.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | Target User ID (must match authenticated user) | Yes | integer |
| userSubscriptionID | path | User Subscription Record ID to cancel | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Subscription cancelled successfully (No Content) |  |
| 400 | Invalid user or subscription ID | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: You are not the owner of this record | object |
| 404 | ErrUserSubscriptionNotFound | object |
| 500 | Internal error during cancellation | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /users/{id}/config

#### PUT
##### Summary

Set or Update User Configuration

##### Description

Creates a new configuration entry for a user, or updates an existing one for the given config_type. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |
| config | body | Configuration Data | Yes | [ports.SetUserConfigInput](#portssetuserconfiginput) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Configuration successfully saved/updated | [ports.UserConfigResponse](#portsuserconfigresponse) |
| 400 | Invalid request body or validation error | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Cannot modify another user's config | object |
| 500 | Database error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /users/{id}/config/{configType}

#### GET
##### Summary

Get User Configuration by Type

##### Description

Retrieves a specific configuration entry for a user. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |
| configType | path | Configuration Type (e.g., user_preferences) | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Configuration successfully retrieved | [ports.UserConfigResponse](#portsuserconfigresponse) |
| 400 | Invalid user ID or missing configType | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Cannot view another user's config | object |
| 404 | ErrUserConfigNotFound | object |
| 500 | Database error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### DELETE
##### Summary

Delete User Configuration by Type

##### Description

Deletes a specific configuration entry for a user. Enforces self-management.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | User ID | Yes | integer |
| configType | path | Configuration Type (e.g., notification_settings) | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | Configuration successfully deleted (No Content) |  |
| 400 | Invalid user ID or missing configType | object |
| 401 | Unauthorized | object |
| 403 | Forbidden: Cannot delete another user's config | object |
| 404 | ErrUserConfigNotFound | object |
| 500 | Database error | object |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### Models

#### models.Business

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean |  | No |
| address | string |  | No |
| businessCategory | [models.BusinessCategory](#modelsbusinesscategory) |  | No |
| businessPhase | [models.BusinessPhase](#modelsbusinessphase) |  | No |
| businessTags | [ [models.BusinessTag](#modelsbusinesstag) ] |  | No |
| businessType | [models.BusinessType](#modelsbusinesstype) |  | No |
| city | string |  | No |
| contactEmail | string |  | No |
| contactName | string |  | No |
| contactPhoneNo | string |  | No |
| country | string |  | No |
| createdAt | string |  | No |
| description | string |  | No |
| id | integer |  | No |
| initiatingConnections | [ [models.BusinessConnection](#modelsbusinessconnection) ] |  | No |
| name | string |  | No |
| operatorUser | [models.User](#modelsuser) |  | No |
| operatorUserID | integer |  | No |
| postalCode | string |  | No |
| projects | [ [models.Project](#modelsproject) ] |  | No |
| publications | [ [models.Publication](#modelspublication) ] |  | No |
| receivingConnections | [ [models.BusinessConnection](#modelsbusinessconnection) ] |  | No |
| state | string |  | No |
| tagline | string |  | No |
| updatedAt | string |  | No |
| value | number |  | No |
| website | string |  | No |

#### models.BusinessCategory

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.BusinessCategory | string |  |  |

#### models.BusinessConnection

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| connectionType | [models.BusinessConnectionType](#modelsbusinessconnectiontype) |  | No |
| createdAt | string |  | No |
| id | integer |  | No |
| initiatedByUser | [models.User](#modelsuser) |  | No |
| initiatedByUserID | integer |  | No |
| initiatingBusiness | [models.Business](#modelsbusiness) |  | No |
| initiatingBusinessID | integer |  | No |
| notes | string |  | No |
| receivingBusiness | [models.Business](#modelsbusiness) |  | No |
| receivingBusinessID | integer |  | No |
| status | [models.BusinessConnectionStatus](#modelsbusinessconnectionstatus) |  | No |
| updatedAt | string |  | No |

#### models.BusinessConnectionStatus

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.BusinessConnectionStatus | string |  |  |

#### models.BusinessConnectionType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.BusinessConnectionType | string |  |  |

#### models.BusinessPhase

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.BusinessPhase | string |  |  |

#### models.BusinessTag

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| business | [models.Business](#modelsbusiness) |  | No |
| businessID | integer |  | No |
| createdAt | string |  | No |
| description | string |  | No |
| id | integer |  | No |
| tagType | [models.BusinessTagType](#modelsbusinesstagtype) |  | No |

#### models.BusinessTagType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.BusinessTagType | string |  |  |

#### models.BusinessType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.BusinessType | string |  |  |

#### models.DailyActivity

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string |  | No |
| enrolments | [ [models.DailyActivityEnrolment](#modelsdailyactivityenrolment) ] |  | No |
| id | integer |  | No |
| name | string |  | No |

#### models.DailyActivityEnrolment

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| dailyActivity | [models.DailyActivity](#modelsdailyactivity) |  | No |
| dailyActivityID | integer |  | No |
| user | [models.User](#modelsuser) |  | No |
| userID | integer |  | No |

#### models.DailyActivityProgressStatus

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.DailyActivityProgressStatus | string |  |  |

#### models.L2EResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| dateAdded | string |  | No |
| id | integer |  | No |
| response | object | FIX: Apply swaggertype:"object" | No |
| user | [models.User](#modelsuser) |  | No |
| userID | integer |  | No |

#### models.Notification

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| actionURL | string |  | No |
| createdAt | string |  | No |
| id | integer |  | No |
| message | string |  | No |
| notificationType | string |  | No |
| read | boolean |  | No |
| receiverUser | [models.User](#modelsuser) |  | No |
| receiverUserID | integer |  | No |
| relatedEntityID | integer |  | No |
| relatedEntityType | string |  | No |
| senderUser | [models.User](#modelsuser) |  | No |
| senderUserID | integer |  | No |
| title | string |  | No |

#### models.Project

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| actualEndDate | string |  | No |
| business | [models.Business](#modelsbusiness) |  | No |
| businessID | integer |  | No |
| createdAt | string |  | No |
| description | string |  | No |
| id | integer |  | No |
| managedByUserID | integer |  | No |
| managingUser | [models.User](#modelsuser) |  | No |
| name | string |  | No |
| projectMembers | [ [models.ProjectMember](#modelsprojectmember) ] |  | No |
| projectRegions | [ [models.ProjectRegion](#modelsprojectregion) ] |  | No |
| projectSkills | [ [models.ProjectSkill](#modelsprojectskill) ] |  | No |
| projectStatus | [models.ProjectStatus](#modelsprojectstatus) |  | No |
| startDate | string |  | No |
| targetEndDate | string |  | No |
| updatedAt | string |  | No |

#### models.ProjectApplicant

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| project | [models.Project](#modelsproject) |  | No |
| projectID | integer |  | No |
| user | [models.User](#modelsuser) |  | No |
| userID | integer |  | No |

#### models.ProjectMember

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| joinedAt | string |  | No |
| project | [models.Project](#modelsproject) |  | No |
| projectID | integer |  | No |
| role | [models.ProjectMemberRole](#modelsprojectmemberrole) |  | No |
| user | [models.User](#modelsuser) |  | No |
| userID | integer |  | No |

#### models.ProjectMemberRole

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.ProjectMemberRole | string |  |  |

#### models.ProjectRegion

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| project | [models.Project](#modelsproject) |  | No |
| projectID | integer |  | No |
| region | [models.Region](#modelsregion) |  | No |
| regionID | string |  | No |

#### models.ProjectSkill

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| importance | [models.ProjectSkillImportance](#modelsprojectskillimportance) |  | No |
| project | [models.Project](#modelsproject) |  | No |
| projectID | integer |  | No |
| skill | [models.Skill](#modelsskill) |  | No |
| skillID | integer |  | No |

#### models.ProjectSkillImportance

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.ProjectSkillImportance | string |  |  |

#### models.ProjectStatus

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.ProjectStatus | string |  |  |

#### models.Publication

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| business | [models.Business](#modelsbusiness) |  | No |
| businessID | integer |  | No |
| content | string |  | No |
| createdAt | string |  | No |
| excerpt | string |  | No |
| id | integer |  | No |
| publicationType | [models.PublicationType](#modelspublicationtype) |  | No |
| published | boolean |  | No |
| publishedAt | string |  | No |
| slug | string |  | No |
| thumbnail | string |  | No |
| title | string |  | No |
| updatedAt | string |  | No |
| user | [models.User](#modelsuser) |  | No |
| userID | integer |  | No |
| videoURL | string |  | No |

#### models.PublicationType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.PublicationType | string |  |  |

#### models.Region

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string |  | No |
| name | string |  | No |

#### models.Skill

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean |  | No |
| category | string |  | No |
| createdAt | string |  | No |
| description | string |  | No |
| id | integer |  | No |
| name | string |  | No |
| projectSkills | [ [models.ProjectSkill](#modelsprojectskill) ] |  | No |
| userSkills | [ [models.UserSkill](#modelsuserskill) ] |  | No |

#### models.Subscription

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | integer |  | No |
| name | string |  | No |
| price | number |  | No |
| validDays | integer |  | No |
| validMonths | integer |  | No |

#### models.User

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean |  | No |
| adkSessionID | string |  | No |
| businesses | [ [models.Business](#modelsbusiness) ] |  | No |
| contactEmail | string |  | No |
| contactPhoneNo | string |  | No |
| createdAt | string |  | No |
| dailyActivityEnrolments | [ [models.DailyActivityEnrolment](#modelsdailyactivityenrolment) ] |  | No |
| dailyActivityProgress | [ [models.UserDailyActivityProgress](#modelsuserdailyactivityprogress) ] |  | No |
| emailVerified | boolean |  | No |
| firstName | string |  | No |
| id | integer |  | No |
| initiatedConnections | [ [models.BusinessConnection](#modelsbusinessconnection) ] |  | No |
| l2EResponses | [ [models.L2EResponse](#modelsl2eresponse) ] |  | No |
| lastName | string |  | No |
| loginEmail | string |  | No |
| managedProjects | [ [models.Project](#modelsproject) ] |  | No |
| passwordHash | string |  | No |
| passwordResetRequestedAt | string |  | No |
| passwordResetToken | [ integer ] |  | No |
| projectApplicants | [ [models.ProjectApplicant](#modelsprojectapplicant) ] |  | No |
| projectMemberships | [ [models.ProjectMember](#modelsprojectmember) ] |  | No |
| publications | [ [models.Publication](#modelspublication) ] |  | No |
| receivedNotifications | [ [models.Notification](#modelsnotification) ] |  | No |
| sentNotifications | [ [models.Notification](#modelsnotification) ] |  | No |
| updatedAt | string |  | No |
| userConfigs | [ [models.UserConfig](#modelsuserconfig) ] |  | No |
| userSessions | [ [models.UserSession](#modelsusersession) ] |  | No |
| userSkills | [ [models.UserSkill](#modelsuserskill) ] |  | No |
| userSubscriptions | [ [models.UserSubscription](#modelsusersubscription) ] |  | No |

#### models.UserConfig

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| config | object | FIX: Apply swaggertype:"object" (already done, confirming for completeness) | No |
| configType | string |  | No |
| id | integer |  | No |
| user | [models.User](#modelsuser) |  | No |
| userID | integer |  | No |

#### models.UserDailyActivityProgress

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| dailyActivity | [models.DailyActivity](#modelsdailyactivity) |  | No |
| dailyActivityID | integer |  | No |
| date | string | FIX: Apply swaggertype:"string" | No |
| progress | integer |  | No |
| status | [models.DailyActivityProgressStatus](#modelsdailyactivityprogressstatus) |  | No |
| user | [models.User](#modelsuser) |  | No |
| userID | integer |  | No |

#### models.UserSession

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| createdAt | string |  | No |
| expiresAt | string |  | No |
| id | integer |  | No |
| ipaddress | string |  | No |
| revokedAt | string |  | No |
| tokenHash | string |  | No |
| user | [models.User](#modelsuser) |  | No |
| userAgent | string |  | No |
| userID | integer |  | No |

#### models.UserSkill

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| createdAt | string |  | No |
| proficiencyLevel | [models.UserSkillProficiency](#modelsuserskillproficiency) |  | No |
| skill | [models.Skill](#modelsskill) |  | No |
| skillID | integer |  | No |
| user | [models.User](#modelsuser) |  | No |
| userID | integer |  | No |

#### models.UserSkillProficiency

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.UserSkillProficiency | string |  |  |

#### models.UserSubscription

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| dateFrom | string |  | No |
| dateTo | string |  | No |
| id | integer |  | No |
| isTrial | boolean |  | No |
| subscription | [models.Subscription](#modelssubscription) |  | No |
| subscriptionID | integer |  | No |
| user | [models.User](#modelsuser) |  | No |
| userID | integer |  | No |

#### ports.ActivityEnrolmentResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| daily_activity_id | integer |  | No |
| user | [ports.UserResponse](#portsuserresponse) |  | No |

#### ports.AddProjectMemberInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| project_id | integer |  | Yes |
| role | [models.ProjectMemberRole](#modelsprojectmemberrole) | *Enum:* `"manager"`, `"contributor"`, `"reviewer"` | Yes |
| user_id | integer |  | Yes |

#### ports.AddProjectRegionInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| project_id | integer |  | Yes |
| region_id | string |  | Yes |

#### ports.BusinessConnectionResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| connection_type | [models.BusinessConnectionType](#modelsbusinessconnectiontype) |  | No |
| created_at | string |  | No |
| id | integer |  | No |
| initiated_by_user | [ports.UserResponse](#portsuserresponse) |  | No |
| initiated_by_user_id | integer |  | No |
| initiating_business | [ports.BusinessResponse](#portsbusinessresponse) |  | No |
| initiating_business_id | integer |  | No |
| notes | string |  | No |
| receiving_business | [ports.BusinessResponse](#portsbusinessresponse) |  | No |
| receiving_business_id | integer |  | No |
| status | [models.BusinessConnectionStatus](#modelsbusinessconnectionstatus) |  | No |
| updated_at | string |  | No |

#### ports.BusinessConnectionsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| connections | [ [ports.BusinessConnectionResponse](#portsbusinessconnectionresponse) ] |  | No |
| count | integer |  | No |

#### ports.BusinessResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean |  | No |
| address | string |  | No |
| business_category | [models.BusinessCategory](#modelsbusinesscategory) |  | No |
| business_phase | [models.BusinessPhase](#modelsbusinessphase) |  | No |
| business_type | [models.BusinessType](#modelsbusinesstype) |  | No |
| city | string |  | No |
| contact_email | string |  | No |
| contact_name | string |  | No |
| contact_phone_no | string |  | No |
| country | string |  | No |
| created_at | string |  | No |
| description | string |  | No |
| id | integer |  | No |
| name | string |  | No |
| operator_user | [ports.UserResponse](#portsuserresponse) |  | No |
| operator_user_id | integer |  | No |
| postal_code | string |  | No |
| state | string |  | No |
| tagline | string |  | No |
| updated_at | string |  | No |
| value | number |  | No |
| website | string |  | No |

#### ports.BusinessTagResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| business | [ports.BusinessResponse](#portsbusinessresponse) |  | No |
| business_id | integer |  | No |
| created_at | string |  | No |
| description | string |  | No |
| id | integer |  | No |
| tag_type | [models.BusinessTagType](#modelsbusinesstagtype) |  | No |

#### ports.BusinessTagsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| count | integer |  | No |
| tags | [ [ports.BusinessTagResponse](#portsbusinesstagresponse) ] |  | No |

#### ports.CreateBusinessConnectionInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| connection_type | [models.BusinessConnectionType](#modelsbusinessconnectiontype) | *Enum:* `"Partnership"`, `"Supplier"`, `"Client"`, `"Referral"`, `"Collaboration"` | Yes |
| initiated_by_user_id | integer |  | Yes |
| initiating_business_id | integer |  | Yes |
| notes | string |  | No |
| receiving_business_id | integer |  | Yes |

#### ports.CreateBusinessInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string |  | No |
| business_category | [models.BusinessCategory](#modelsbusinesscategory) |  | Yes |
| business_phase | [models.BusinessPhase](#modelsbusinessphase) |  | Yes |
| business_type | [models.BusinessType](#modelsbusinesstype) |  | Yes |
| city | string |  | No |
| contact_email | string |  | No |
| contact_name | string |  | No |
| contact_phone_no | string |  | No |
| country | string |  | No |
| description | string |  | No |
| name | string |  | Yes |
| operator_user_id | integer |  | Yes |
| postal_code | string |  | No |
| state | string |  | No |
| tagline | string |  | No |
| value | number |  | No |
| website | string |  | No |

#### ports.CreateBusinessTagInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| business_id | integer |  | Yes |
| description | string |  | Yes |
| tag_type | [models.BusinessTagType](#modelsbusinesstagtype) | *Enum:* `"client"`, `"service"`, `"specialty"` | Yes |

#### ports.CreateDailyActivityInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string |  | Yes |
| name | string |  | Yes |

#### ports.CreateEventInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| event_type | string |  | Yes |
| payload | object |  | Yes |
| user_id | integer |  | No |

#### ports.CreateFeedbackInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| content | string |  | Yes |
| email | string |  | Yes |
| name | string |  | Yes |

#### ports.CreateInferredConnectionInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| confidence_score | number |  | Yes |
| connection_type | string |  | Yes |
| model_version | string |  | No |
| source_entity_id | integer |  | Yes |
| source_entity_type | string |  | Yes |
| target_entity_id | integer |  | Yes |
| target_entity_type | string |  | Yes |

#### ports.CreateL2EResponseInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| response | object |  | Yes |

#### ports.CreateNotificationInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| action_url | string |  | No |
| message | string |  | Yes |
| notification_type | string |  | Yes |
| receiver_user_id | integer |  | Yes |
| related_entity_id | integer |  | No |
| related_entity_type | string |  | No |
| sender_user_id | integer |  | No |
| title | string |  | Yes |

#### ports.CreateProjectInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| business_id | integer |  | No |
| description | string |  | No |
| managed_by_user_id | integer |  | Yes |
| name | string |  | Yes |
| project_status | [models.ProjectStatus](#modelsprojectstatus) |  | Yes |
| region_ids | [ string ] |  | No |
| start_date | string |  | No |
| target_end_date | string |  | No |

#### ports.CreateProjectSkillInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| importance | [models.ProjectSkillImportance](#modelsprojectskillimportance) | *Enum:* `"required"`, `"preferred"`, `"optional"` | Yes |
| project_id | integer |  | Yes |
| skill_id | integer |  | Yes |

#### ports.CreatePublicationInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| business_id | integer |  | No |
| content | string |  | Yes |
| excerpt | string |  | No |
| publication_type | [models.PublicationType](#modelspublicationtype) |  | Yes |
| published | boolean |  | No |
| thumbnail | string |  | No |
| title | string |  | Yes |
| user_id | integer |  | Yes |
| video_url | string |  | No |

#### ports.CreateSkillInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean |  | No |
| category | string |  | Yes |
| description | string |  | No |
| name | string |  | Yes |

#### ports.CreateSubscriptionInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | Yes |
| price | number |  | Yes |
| valid_days | integer |  | No |
| valid_months | integer |  | No |

#### ports.CreateUserSkillInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| proficiency_level | [models.UserSkillProficiency](#modelsuserskillproficiency) | *Enum:* `"beginner"`, `"intermediate"`, `"advanced"`, `"expert"` | Yes |
| skill_id | integer |  | Yes |
| user_id | integer |  | Yes |

#### ports.DailyActivityResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string |  | No |
| id | integer |  | No |
| name | string |  | No |

#### ports.EventResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| event_type | string |  | No |
| id | integer |  | No |
| payload | object |  | No |
| timestamp | string |  | No |
| user | [ports.UserResponse](#portsuserresponse) |  | No |
| user_id | integer |  | No |

#### ports.FeedbackResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| content | string |  | No |
| date_submitted | string |  | No |
| email | string |  | No |
| id | integer |  | No |
| name | string |  | No |

#### ports.InferredConnectionResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| confidence_score | number |  | No |
| connection_type | string |  | No |
| created_at | string |  | No |
| id | integer |  | No |
| model_version | string |  | No |
| source_entity_id | integer |  | No |
| source_entity_type | string |  | No |
| target_entity_id | integer |  | No |
| target_entity_type | string |  | No |

#### ports.L2EResponseResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| date_added | string |  | No |
| id | integer |  | No |
| response | object |  | No |
| user_id | integer |  | No |

#### ports.LoginInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| login_email | string |  | Yes |
| password | string |  | Yes |

#### ports.LoginResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| expires_at | string |  | No |
| session_id | integer |  | No |
| token | string |  | No |
| token_type | string |  | No |
| user | [ports.UserResponse](#portsuserresponse) |  | No |

#### ports.NotificationResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| action_url | string |  | No |
| created_at | string |  | No |
| id | integer |  | No |
| message | string |  | No |
| notification_type | string |  | No |
| read | boolean |  | No |
| receiver | [ports.UserResponse](#portsuserresponse) |  | No |
| related_entity_id | integer |  | No |
| related_entity_type | string |  | No |
| sender | [ports.UserResponse](#portsuserresponse) |  | No |
| title | string |  | No |

#### ports.ProjectApplicantResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| project_id | integer |  | No |
| user | [ports.UserResponse](#portsuserresponse) |  | No |
| user_id | integer |  | No |

#### ports.ProjectMemberResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| joined_at | string |  | No |
| project | [ports.ProjectResponse](#portsprojectresponse) |  | No |
| project_id | integer |  | No |
| role | [models.ProjectMemberRole](#modelsprojectmemberrole) |  | No |
| user | [ports.UserResponse](#portsuserresponse) |  | No |
| user_id | integer |  | No |

#### ports.ProjectMembersResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| count | integer |  | No |
| members | [ [ports.ProjectMemberResponse](#portsprojectmemberresponse) ] |  | No |

#### ports.ProjectRegionResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| project_id | integer |  | No |
| region | [ports.RegionResponse](#portsregionresponse) |  | No |
| region_id | string |  | No |

#### ports.ProjectResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| actual_end_date | string |  | No |
| business | [ports.BusinessResponse](#portsbusinessresponse) |  | No |
| created_at | string |  | No |
| description | string |  | No |
| id | integer |  | No |
| manager | [ports.UserResponse](#portsuserresponse) |  | No |
| members | [ [ports.ProjectMemberResponse](#portsprojectmemberresponse) ] |  | No |
| name | string |  | No |
| project_status | [models.ProjectStatus](#modelsprojectstatus) |  | No |
| regions | [ [ports.RegionResponse](#portsregionresponse) ] |  | No |
| start_date | string |  | No |
| target_end_date | string |  | No |
| updated_at | string |  | No |

#### ports.ProjectSkillResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| importance | [models.ProjectSkillImportance](#modelsprojectskillimportance) |  | No |
| project | [ports.ProjectResponse](#portsprojectresponse) |  | No |
| project_id | integer |  | No |
| skill | [ports.SkillResponse](#portsskillresponse) |  | No |
| skill_id | integer |  | No |

#### ports.ProjectSkillsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| count | integer |  | No |
| skills | [ [ports.ProjectSkillResponse](#portsprojectskillresponse) ] |  | No |

#### ports.PublicationResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| author | [ports.UserResponse](#portsuserresponse) |  | No |
| business | [ports.BusinessResponse](#portsbusinessresponse) |  | No |
| content | string |  | No |
| created_at | string |  | No |
| excerpt | string |  | No |
| id | integer |  | No |
| publication_type | [models.PublicationType](#modelspublicationtype) |  | No |
| published | boolean |  | No |
| published_at | string |  | No |
| slug | string |  | No |
| thumbnail | string |  | No |
| title | string |  | No |
| updated_at | string |  | No |
| video_url | string |  | No |

#### ports.RegionResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string |  | No |
| name | string |  | No |

#### ports.SetUserConfigInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| config | object |  | Yes |
| config_type | string |  | Yes |

#### ports.SkillResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean |  | No |
| category | string |  | No |
| created_at | string |  | No |
| description | string |  | No |
| id | integer |  | No |
| name | string |  | No |

#### ports.SubscriptionResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | integer |  | No |
| name | string |  | No |
| price | number |  | No |
| valid_days | integer |  | No |
| valid_months | integer |  | No |

#### ports.UpdateBusinessConnectionInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| connection_type | [models.BusinessConnectionType](#modelsbusinessconnectiontype) | *Enum:* `"Partnership"`, `"Supplier"`, `"Client"`, `"Referral"`, `"Collaboration"` | No |
| notes | string |  | No |
| status | [models.BusinessConnectionStatus](#modelsbusinessconnectionstatus) | *Enum:* `"pending"`, `"active"`, `"rejected"`, `"inactive"` | No |

#### ports.UpdateBusinessInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean |  | No |
| address | string |  | No |
| business_category | [models.BusinessCategory](#modelsbusinesscategory) |  | No |
| business_phase | [models.BusinessPhase](#modelsbusinessphase) |  | No |
| business_type | [models.BusinessType](#modelsbusinesstype) |  | No |
| city | string |  | No |
| contact_email | string |  | No |
| contact_name | string |  | No |
| contact_phone_no | string |  | No |
| country | string |  | No |
| description | string |  | No |
| name | string |  | No |
| postal_code | string |  | No |
| state | string |  | No |
| tagline | string |  | No |
| value | number |  | No |
| website | string |  | No |

#### ports.UpdateProjectInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| actual_end_date | string |  | No |
| business_id | integer |  | No |
| description | string |  | No |
| managed_by_user_id | integer |  | No |
| name | string |  | No |
| project_status | [models.ProjectStatus](#modelsprojectstatus) |  | No |
| region_ids | [ string ] |  | No |
| start_date | string |  | No |
| target_end_date | string |  | No |

#### ports.UpdateProjectMemberRoleInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| role | [models.ProjectMemberRole](#modelsprojectmemberrole) | *Enum:* `"manager"`, `"contributor"`, `"reviewer"` | Yes |

#### ports.UpdateProjectSkillInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| importance | [models.ProjectSkillImportance](#modelsprojectskillimportance) | *Enum:* `"required"`, `"preferred"`, `"optional"` | No |

#### ports.UpdatePublicationInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| content | string |  | No |
| excerpt | string |  | No |
| published | boolean |  | No |
| thumbnail | string |  | No |
| title | string |  | No |
| video_url | string |  | No |

#### ports.UpdateSkillInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean |  | No |
| category | string |  | No |
| description | string |  | No |
| name | string |  | No |

#### ports.UpdateUserSkillInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| proficiency_level | [models.UserSkillProficiency](#modelsuserskillproficiency) | *Enum:* `"beginner"`, `"intermediate"`, `"advanced"`, `"expert"` | No |

#### ports.UserApplicationResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| project | [ports.ProjectResponse](#portsprojectresponse) |  | No |
| project_id | integer |  | No |
| user_id | integer |  | No |

#### ports.UserConfigResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| config | object |  | No |
| config_type | string |  | No |
| user_id | integer |  | No |

#### ports.UserCreationSchema

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| adk_session_id | string |  | No |
| contact_email | string |  | No |
| contact_phone_no | string |  | No |
| first_name | string |  | Yes |
| last_name | string |  | No |
| login_email | string |  | Yes |
| password | string |  | Yes |

#### ports.UserEnrolmentResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| daily_activity | [ports.DailyActivityResponse](#portsdailyactivityresponse) |  | No |
| user_id | integer |  | No |

#### ports.UserResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean |  | No |
| contact_email | string |  | No |
| contact_phone_no | string |  | No |
| created_at | string |  | No |
| email_verified | boolean |  | No |
| first_name | string |  | No |
| id | integer |  | No |
| last_name | string |  | No |
| login_email | string |  | No |
| updated_at | string |  | No |

#### ports.UserSkillResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created_at | string |  | No |
| proficiency_level | [models.UserSkillProficiency](#modelsuserskillproficiency) |  | No |
| skill | [ports.SkillResponse](#portsskillresponse) |  | No |
| skill_id | integer |  | No |
| user | [ports.UserResponse](#portsuserresponse) |  | No |
| user_id | integer |  | No |

#### ports.UserSkillsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| count | integer |  | No |
| skills | [ [ports.UserSkillResponse](#portsuserskillresponse) ] |  | No |

#### ports.UserSubscribeInput

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| subscription_id | integer |  | Yes |
| user_id | integer |  | Yes |

#### ports.UserSubscriptionResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| date_from | string |  | No |
| date_to | string |  | No |
| id | integer |  | No |
| is_trial | boolean |  | No |
| subscription | [ports.SubscriptionResponse](#portssubscriptionresponse) |  | No |
| user | [ports.UserResponse](#portsuserresponse) |  | No |

#### ports.UserUpdateSchema

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean |  | No |
| adk_session_id | string |  | No |
| contact_email | string |  | No |
| contact_phone_no | string |  | No |
| email_verified | boolean |  | No |
| first_name | string |  | No |
| last_name | string |  | No |
| login_email | string |  | No |
| password | string |  | No |
