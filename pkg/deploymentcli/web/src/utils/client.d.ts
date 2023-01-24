import type {
  OpenAPIClient,
  Parameters,
  UnknownParamsObject,
  OperationResponse,
  AxiosRequestConfig,
} from 'openapi-client-axios'; 

declare namespace Components {
    namespace Responses {
        export interface ErrorResponse {
            error?: string;
        }
        export interface GrantResponse {
            grant?: /**
             * Grant
             * A temporary assignment of a user to a principal.
             */
            Schemas.Grant;
        }
        export interface HealthResponse {
            health?: /* ProviderHealth */ Schemas.ProviderHealth;
        }
    }
    namespace Schemas {
        /**
         * AccessRule
         */
        export interface AccessRule {
            id?: string;
            provider?: /* Provider */ Provider;
            /**
             * This is the target value for the provider, for AWS it will be a Role ARN, for Okta - a group name etc.
             * example:
             * arn:${Partition}:iam::${Account}:role/${RoleNameWithPath}
             */
            provider_target?: string;
            /**
             * The start time of the grant in Unix nanoseconds.
             * example:
             * 1257894000000000000
             */
            from_hours?: null | number;
            /**
             * The end time of the grant in Unix nanoseconds.
             * example:
             * 1257897600000000000
             */
            to_hours?: null | number;
            /**
             * A max duration for the access in Unix nanoseconds.
             * example:
             * 1257897600000000000
             */
            max_duration?: number;
            time_constraints?: string | {
                [key: string]: any;
            };
            start_date?: number;
            end_date?: number;
            /**
             * Array of group ids that the access rule applies to
             */
            groups?: /* Group */ Group[];
            approval_required?: boolean;
            approvers?: /* User */ User[];
            integrations?: {
                integration_type?: string;
                integration_fields?: {
                    input_name?: string;
                    input_placeholder?: string;
                    input_type?: string;
                }[];
            }[];
            ""?: string;
        }
        /**
         * CreateGrant
         * A grant to be created.
         */
        export interface CreateGrant {
            /**
             * The email address of the user to grant access to.
             */
            subject: string; // email
            /**
             * The ID of the provider to grant access to.
             * example:
             * okta
             */
            provider: string;
            /**
             * Provider-specific grant data. Must match the provider's schema.
             */
            with: {
                [key: string]: any;
            };
            /**
             * The start time of the grant in Unix milliseconds.
             * example:
             * 1257894000000000000
             */
            start: number;
            /**
             * The end time of the grant in Unix milliseconds.
             * example:
             * 1257897600000000000
             */
            end: number;
        }
        /**
         * Grant
         * A temporary assignment of a user to a principal.
         */
        export interface Grant {
            /**
             * example:
             * aba0dcba-0a8c-4393-ad92-69510326b29a
             */
            id: string;
            /**
             * The current state of the grant.
             */
            status: "pending" | "active" | "deactivated" | "error";
            /**
             * The email address of the user to grant access to.
             */
            subject: string; // email
            /**
             * The ID of the provider to grant access to.
             * example:
             * okta
             */
            provider: string;
            /**
             * Provider-specific grant data. Must match the provider's schema.
             */
            with: {
                [key: string]: any;
            };
            /**
             * The start time of the grant in Unix nanoseconds.
             * example:
             * 1257894000000000000
             */
            start: number;
            /**
             * The end time of the grant in Unix nanoseconds.
             * example:
             * 1257897600000000000
             */
            end: number;
        }
        /**
         * Group
         */
        export interface Group {
            id?: string;
            name?: string;
            group_members?: /* User */ User[];
            provider_type?: "aws" | "okta";
        }
        /**
         * Provider
         */
        export interface Provider {
            /**
             * References the provider's unique ID
             */
            id?: string;
            name?: string;
            type?: "okta" | "aws";
        }
        /**
         * ProviderHealth
         */
        export interface ProviderHealth {
            /**
             * The provider ID.
             * example:
             * okta
             */
            id: string;
            /**
             * Whether the provider is healthy.
             */
            healthy: boolean;
            /**
             * A descriptive error message, if the provider isn't healthy.
             * example:
             * API_TOKEN secret has not been provided
             */
            error?: string | null;
        }
        /**
         * Request
         *
         */
        export interface Request {
            id?: string;
            requestor_id: /* User */ User;
            approver_id?: {
                approver?: /* User */ User;
                decided_at?: number;
            }[];
            /**
             * Contains the decision/status of the request
             *
             */
            status?: string;
            reason?: string;
            duration?: number;
            requested_at?: number;
        }
        /**
         * User
         */
        export interface User {
            id?: string;
            oid_sub?: string;
            email?: string;
            name?: string;
            picture?: string;
            is_admin?: boolean;
            created_at?: number;
            updated_at?: number;
        }
    }
}
declare namespace Paths {
    namespace ApiV1Grants$GrantIdRevoke {
        namespace Parameters {
            export type GrantId = string;
        }
        export interface PathParameters {
            grantId: Parameters.GrantId;
        }
    }
    namespace ApiV1Groups$GroupId {
        namespace Parameters {
            export type GroupId = string;
        }
        export interface PathParameters {
            groupId: Parameters.GroupId;
        }
    }
    namespace ApiV1Requests$RequestId {
        namespace Parameters {
            export type RequestId = string;
        }
        export interface PathParameters {
            requestId: Parameters.RequestId;
        }
    }
    namespace ApiV1Rules$RuleId {
        namespace Parameters {
            export type RuleId = string;
        }
        export interface PathParameters {
            ruleId: Parameters.RuleId;
        }
    }
    namespace ApiV1Users$UserId {
        namespace Parameters {
            export type UserId = string;
        }
        export interface PathParameters {
            userId: Parameters.UserId;
        }
    }
    namespace ApiV1Users$UserIdRequests {
        namespace Parameters {
            export type UserId = string;
        }
        export interface PathParameters {
            userId: Parameters.UserId;
        }
    }
    namespace GetApiV1GroupsGroupId {
        namespace Responses {
            export type $200 = /* Group */ Components.Schemas.Group;
        }
    }
    namespace GetApiV1Rules {
        namespace Responses {
            export type $200 = /* AccessRule */ Components.Schemas.AccessRule[];
        }
    }
    namespace GetApiV1UsersUserIdRequests {
        namespace Responses {
            export type $200 = /**
             * Request
             *
             */
            Components.Schemas.Request[];
        }
    }
    namespace GetGrants {
        namespace Responses {
            export interface $200 {
                grants?: /**
                 * Grant
                 * A temporary assignment of a user to a principal.
                 */
                Components.Schemas.Grant[];
            }
        }
    }
    namespace GetHealth {
        namespace Responses {
            export type $200 = Components.Responses.HealthResponse;
            export type $500 = Components.Responses.HealthResponse;
        }
    }
    namespace PostApiV1Requests {
        namespace Responses {
            export interface $200 {
            }
        }
    }
    namespace PostApiV1Rules {
        export type RequestBody = /* AccessRule */ Components.Schemas.AccessRule;
        namespace Responses {
            export interface $200 {
            }
        }
    }
    namespace PostGrants {
        export type RequestBody = /**
         * CreateGrant
         * A grant to be created.
         */
        Components.Schemas.CreateGrant;
        namespace Responses {
            export type $201 = Components.Responses.GrantResponse;
            export type $400 = Components.Responses.ErrorResponse;
        }
    }
    namespace PostGrantsRevoke {
        namespace Responses {
            export interface $200 {
            }
        }
    }
}

export interface OperationMethods {
  /**
   * get-grants - List Grants
   * 
   * List grants.
   */
  'get-grants'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: any,
    config?: AxiosRequestConfig  
  ): OperationResponse<Paths.GetGrants.Responses.$200>
  /**
   * post-grants - Create Grant
   * 
   * Create a grant.
   * 
   * The returned grant ID will depend on the Access Handler's runtime. When running on AWS Lambda with Step Functions, this ID is the invocation ID of the Step Functions workflow run.
   */
  'post-grants'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: Paths.PostGrants.RequestBody,
    config?: AxiosRequestConfig  
  ): OperationResponse<Paths.PostGrants.Responses.$201>
  /**
   * post-grants-revoke - Revoke grant
   * 
   * Revoke an active grant.
   */
  'post-grants-revoke'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: any,
    config?: AxiosRequestConfig  
  ): OperationResponse<Paths.PostGrantsRevoke.Responses.$200>
  /**
   * get-health - Healthcheck
   * 
   * Returns information on the health of the runtime and providers. If any healthchecks fail the response code will be 500 (Internal Server Error).
   */
  'get-health'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: any,
    config?: AxiosRequestConfig  
  ): OperationResponse<Paths.GetHealth.Responses.$200>
  /**
   * get-api-v1-users-userId - Your GET endpoint
   * 
   * For individual user actions
   */
  'get-api-v1-users-userId'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: any,
    config?: AxiosRequestConfig  
  ): OperationResponse<any>
  /**
   * get-api-v1-users - Your GET endpoint
   * 
   * Fetch a list of users
   */
  'get-api-v1-users'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: any,
    config?: AxiosRequestConfig  
  ): OperationResponse<any>
  /**
   * get-api-v1-requests - Your GET endpoint
   * 
   * Get all requests
   */
  'get-api-v1-requests'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: any,
    config?: AxiosRequestConfig  
  ): OperationResponse<any>
  /**
   * post-api-v1-requests - Create a request
   */
  'post-api-v1-requests'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: any,
    config?: AxiosRequestConfig  
  ): OperationResponse<Paths.PostApiV1Requests.Responses.$200>
  /**
   * get-api-v1-requests-$-requestId - Your GET endpoint
   */
  'get-api-v1-requests-$-requestId'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: any,
    config?: AxiosRequestConfig  
  ): OperationResponse<any>
  /**
   * get-api-v1-rules - Your GET endpoint
   * 
   * Get all access rules
   */
  'get-api-v1-rules'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: any,
    config?: AxiosRequestConfig  
  ): OperationResponse<Paths.GetApiV1Rules.Responses.$200>
  /**
   * post-api-v1-rules - Create an access rule
   */
  'post-api-v1-rules'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: Paths.PostApiV1Rules.RequestBody,
    config?: AxiosRequestConfig  
  ): OperationResponse<Paths.PostApiV1Rules.Responses.$200>
  /**
   * get-api-v1-rule - Your GET endpoint
   */
  'get-api-v1-rule'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: any,
    config?: AxiosRequestConfig  
  ): OperationResponse<any>
  /**
   * get-api-v1-groups-groupId - Your GET endpoint
   * 
   * Gets detailed Group by group id
   */
  'get-api-v1-groups-groupId'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: any,
    config?: AxiosRequestConfig  
  ): OperationResponse<Paths.GetApiV1GroupsGroupId.Responses.$200>
  /**
   * get-api-v1-users-userId-requests - Your GET endpoint
   * 
   * Gets an array of requests for a given userId
   */
  'get-api-v1-users-userId-requests'(
    parameters?: Parameters<UnknownParamsObject> | null,
    data?: any,
    config?: AxiosRequestConfig  
  ): OperationResponse<Paths.GetApiV1UsersUserIdRequests.Responses.$200>
}

export interface PathsDictionary {
  ['/api/v1/grants']: {
    /**
     * get-grants - List Grants
     * 
     * List grants.
     */
    'get'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: any,
      config?: AxiosRequestConfig  
    ): OperationResponse<Paths.GetGrants.Responses.$200>
    /**
     * post-grants - Create Grant
     * 
     * Create a grant.
     * 
     * The returned grant ID will depend on the Access Handler's runtime. When running on AWS Lambda with Step Functions, this ID is the invocation ID of the Step Functions workflow run.
     */
    'post'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: Paths.PostGrants.RequestBody,
      config?: AxiosRequestConfig  
    ): OperationResponse<Paths.PostGrants.Responses.$201>
  }
  ['/api/v1/grants/{grantId}/revoke']: {
    /**
     * post-grants-revoke - Revoke grant
     * 
     * Revoke an active grant.
     */
    'post'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: any,
      config?: AxiosRequestConfig  
    ): OperationResponse<Paths.PostGrantsRevoke.Responses.$200>
  }
  ['/api/v1/health']: {
    /**
     * get-health - Healthcheck
     * 
     * Returns information on the health of the runtime and providers. If any healthchecks fail the response code will be 500 (Internal Server Error).
     */
    'get'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: any,
      config?: AxiosRequestConfig  
    ): OperationResponse<Paths.GetHealth.Responses.$200>
  }
  ['/api/v1/users/{userId}']: {
    /**
     * get-api-v1-users-userId - Your GET endpoint
     * 
     * For individual user actions
     */
    'get'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: any,
      config?: AxiosRequestConfig  
    ): OperationResponse<any>
  }
  ['/api/v1/users/']: {
    /**
     * get-api-v1-users - Your GET endpoint
     * 
     * Fetch a list of users
     */
    'get'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: any,
      config?: AxiosRequestConfig  
    ): OperationResponse<any>
  }
  ['/api/v1/requests/']: {
    /**
     * get-api-v1-requests - Your GET endpoint
     * 
     * Get all requests
     */
    'get'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: any,
      config?: AxiosRequestConfig  
    ): OperationResponse<any>
    /**
     * post-api-v1-requests - Create a request
     */
    'post'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: any,
      config?: AxiosRequestConfig  
    ): OperationResponse<Paths.PostApiV1Requests.Responses.$200>
  }
  ['/api/v1/requests/{requestId}']: {
    /**
     * get-api-v1-requests-$-requestId - Your GET endpoint
     */
    'get'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: any,
      config?: AxiosRequestConfig  
    ): OperationResponse<any>
  }
  ['/api/v1/rules/']: {
    /**
     * get-api-v1-rules - Your GET endpoint
     * 
     * Get all access rules
     */
    'get'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: any,
      config?: AxiosRequestConfig  
    ): OperationResponse<Paths.GetApiV1Rules.Responses.$200>
    /**
     * post-api-v1-rules - Create an access rule
     */
    'post'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: Paths.PostApiV1Rules.RequestBody,
      config?: AxiosRequestConfig  
    ): OperationResponse<Paths.PostApiV1Rules.Responses.$200>
  }
  ['/api/v1/rules/{ruleId}']: {
    /**
     * get-api-v1-rule - Your GET endpoint
     */
    'get'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: any,
      config?: AxiosRequestConfig  
    ): OperationResponse<any>
  }
  ['/api/v1/groups/{groupId}']: {
    /**
     * get-api-v1-groups-groupId - Your GET endpoint
     * 
     * Gets detailed Group by group id
     */
    'get'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: any,
      config?: AxiosRequestConfig  
    ): OperationResponse<Paths.GetApiV1GroupsGroupId.Responses.$200>
  }
  ['/api/v1/users/{userId}/requests']: {
    /**
     * get-api-v1-users-userId-requests - Your GET endpoint
     * 
     * Gets an array of requests for a given userId
     */
    'get'(
      parameters?: Parameters<UnknownParamsObject> | null,
      data?: any,
      config?: AxiosRequestConfig  
    ): OperationResponse<Paths.GetApiV1UsersUserIdRequests.Responses.$200>
  }
}

export type Client = OpenAPIClient<OperationMethods, PathsDictionary>
