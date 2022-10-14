/**
 * Generated by orval v6.9.6 🍺
 * Do not edit manually.
 * Approvals
 * Granted Approvals API
 * OpenAPI spec version: 1.0
 */
import useSwr from 'swr'
import type {
  SWRConfiguration,
  Key
} from 'swr'
import type {
  ListAccessRulesDetailResponseResponse,
  AdminListAccessRulesParams,
  AccessRuleDetail,
  ErrorResponseResponse,
  CreateAccessRuleRequestBody,
  ListRequestsResponseResponse,
  AdminListRequestsParams,
  User,
  UpdateUserBody,
  ListUserResponseResponse,
  GetUsersParams,
  CreateUserRequestBody,
  ListGroupsResponseResponse,
  GetGroupsParams,
  Group,
  CreateGroupRequestBody,
  Provider,
  ListProviderArgOptionsParams,
  ListProviderSetupsResponseResponse,
  ProviderSetupResponseResponse,
  ProviderSetupInstructions,
  CompleteProviderSetupResponseResponse,
  ProviderSetupStepCompleteRequestBody,
  IdentityConfigurationResponseResponse
} from '.././types'
import type {
  ArgSchemaResponseResponse,
  ArgOptionsResponseResponse
} from '.././types/accesshandler-openapi.yml'
import { customInstance } from '../../custom-instance'
import type { ErrorType } from '../../custom-instance'


  
  // eslint-disable-next-line
  type SecondParameter<T extends (...args: any) => any> = T extends (
  config: any,
  args: infer P,
) => any
  ? P
  : never;

/**
 * List all access rules
 * @summary List Access Rules
 */
export const adminListAccessRules = (
    params?: AdminListAccessRulesParams,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ListAccessRulesDetailResponseResponse>(
      {url: `/api/v1/admin/access-rules`, method: 'get',
        params
    },
      options);
    }
  

export const getAdminListAccessRulesKey = (params?: AdminListAccessRulesParams,) => [`/api/v1/admin/access-rules`, ...(params ? [params]: [])];

    
export type AdminListAccessRulesQueryResult = NonNullable<Awaited<ReturnType<typeof adminListAccessRules>>>
export type AdminListAccessRulesQueryError = ErrorType<unknown>

export const useAdminListAccessRules = <TError = ErrorType<unknown>>(
 params?: AdminListAccessRulesParams, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof adminListAccessRules>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getAdminListAccessRulesKey(params) : null);
  const swrFn = () => adminListAccessRules(params, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Create an access rule
 * @summary Create Access Rule
 */
export const adminCreateAccessRule = (
    createAccessRuleRequestBody: CreateAccessRuleRequestBody,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<AccessRuleDetail>(
      {url: `/api/v1/admin/access-rules`, method: 'post',
      headers: {'Content-Type': 'application/json', },
      data: createAccessRuleRequestBody
    },
      options);
    }
  

/**
 * Get an Access Rule.
 * @summary Get Access Rule
 */
export const adminGetAccessRule = (
    ruleId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<AccessRuleDetail>(
      {url: `/api/v1/admin/access-rules/${ruleId}`, method: 'get'
    },
      options);
    }
  

export const getAdminGetAccessRuleKey = (ruleId: string,) => [`/api/v1/admin/access-rules/${ruleId}`];

    
export type AdminGetAccessRuleQueryResult = NonNullable<Awaited<ReturnType<typeof adminGetAccessRule>>>
export type AdminGetAccessRuleQueryError = ErrorType<ErrorResponseResponse>

export const useAdminGetAccessRule = <TError = ErrorType<ErrorResponseResponse>>(
 ruleId: string, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof adminGetAccessRule>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false && !!(ruleId)
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getAdminGetAccessRuleKey(ruleId) : null);
  const swrFn = () => adminGetAccessRule(ruleId, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Updates an Access Rule. Updating a rule creates a new version.
 * @summary Update Access Rule
 */
export const adminUpdateAccessRule = (
    ruleId: string,
    createAccessRuleRequestBody: CreateAccessRuleRequestBody,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<AccessRuleDetail>(
      {url: `/api/v1/admin/access-rules/${ruleId}`, method: 'put',
      headers: {'Content-Type': 'application/json', },
      data: createAccessRuleRequestBody
    },
      options);
    }
  

/**
 * Returns a version history for a particular Access Rule.
 * @summary Get Access Rule version history
 */
export const adminGetAccessRuleVersions = (
    ruleId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ListAccessRulesDetailResponseResponse>(
      {url: `/api/v1/admin/access-rules/${ruleId}/versions`, method: 'get'
    },
      options);
    }
  

export const getAdminGetAccessRuleVersionsKey = (ruleId: string,) => [`/api/v1/admin/access-rules/${ruleId}/versions`];

    
export type AdminGetAccessRuleVersionsQueryResult = NonNullable<Awaited<ReturnType<typeof adminGetAccessRuleVersions>>>
export type AdminGetAccessRuleVersionsQueryError = ErrorType<void>

export const useAdminGetAccessRuleVersions = <TError = ErrorType<void>>(
 ruleId: string, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof adminGetAccessRuleVersions>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false && !!(ruleId)
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getAdminGetAccessRuleVersionsKey(ruleId) : null);
  const swrFn = () => adminGetAccessRuleVersions(ruleId, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Returns a specific version for an Access Rule.
 * @summary Get Access Rule Version
 */
export const adminGetAccessRuleVersion = (
    ruleId: string,
    version: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<AccessRuleDetail>(
      {url: `/api/v1/admin/access-rules/${ruleId}/versions/${version}`, method: 'get'
    },
      options);
    }
  

export const getAdminGetAccessRuleVersionKey = (ruleId: string,
    version: string,) => [`/api/v1/admin/access-rules/${ruleId}/versions/${version}`];

    
export type AdminGetAccessRuleVersionQueryResult = NonNullable<Awaited<ReturnType<typeof adminGetAccessRuleVersion>>>
export type AdminGetAccessRuleVersionQueryError = ErrorType<ErrorResponseResponse>

export const useAdminGetAccessRuleVersion = <TError = ErrorType<ErrorResponseResponse>>(
 ruleId: string,
    version: string, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof adminGetAccessRuleVersion>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false && !!(ruleId && version)
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getAdminGetAccessRuleVersionKey(ruleId,version) : null);
  const swrFn = () => adminGetAccessRuleVersion(ruleId,version, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Return a list of all requests
 * @summary Your GET endpoint
 */
export const adminListRequests = (
    params?: AdminListRequestsParams,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ListRequestsResponseResponse>(
      {url: `/api/v1/admin/requests`, method: 'get',
        params
    },
      options);
    }
  

export const getAdminListRequestsKey = (params?: AdminListRequestsParams,) => [`/api/v1/admin/requests`, ...(params ? [params]: [])];

    
export type AdminListRequestsQueryResult = NonNullable<Awaited<ReturnType<typeof adminListRequests>>>
export type AdminListRequestsQueryError = ErrorType<unknown>

export const useAdminListRequests = <TError = ErrorType<unknown>>(
 params?: AdminListRequestsParams, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof adminListRequests>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getAdminListRequestsKey(params) : null);
  const swrFn = () => adminListRequests(params, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Update a user including group membership
 * @summary Update User
 */
export const updateUser = (
    userId: string,
    updateUserBody: UpdateUserBody,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<User>(
      {url: `/api/v1/admin/users/${userId}`, method: 'post',
      headers: {'Content-Type': 'application/json', },
      data: updateUserBody
    },
      options);
    }
  

/**
 * Fetch a list of users
 * @summary Returns a list of users
 */
export const getUsers = (
    params?: GetUsersParams,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ListUserResponseResponse>(
      {url: `/api/v1/admin/users`, method: 'get',
        params
    },
      options);
    }
  

export const getGetUsersKey = (params?: GetUsersParams,) => [`/api/v1/admin/users`, ...(params ? [params]: [])];

    
export type GetUsersQueryResult = NonNullable<Awaited<ReturnType<typeof getUsers>>>
export type GetUsersQueryError = ErrorType<unknown>

export const useGetUsers = <TError = ErrorType<unknown>>(
 params?: GetUsersParams, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof getUsers>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getGetUsersKey(params) : null);
  const swrFn = () => getUsers(params, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Create new user in the Cognito user pool if it is enabled.
 * @summary Create User
 */
export const createUser = (
    createUserRequestBody: CreateUserRequestBody,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<User>(
      {url: `/api/v1/admin/users`, method: 'post',
      headers: {'Content-Type': 'application/json', },
      data: createUserRequestBody
    },
      options);
    }
  

/**
 * Gets all groups
 * @summary List groups
 */
export const getGroups = (
    params?: GetGroupsParams,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ListGroupsResponseResponse>(
      {url: `/api/v1/admin/groups`, method: 'get',
        params
    },
      options);
    }
  

export const getGetGroupsKey = (params?: GetGroupsParams,) => [`/api/v1/admin/groups`, ...(params ? [params]: [])];

    
export type GetGroupsQueryResult = NonNullable<Awaited<ReturnType<typeof getGroups>>>
export type GetGroupsQueryError = ErrorType<unknown>

export const useGetGroups = <TError = ErrorType<unknown>>(
 params?: GetGroupsParams, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof getGroups>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getGetGroupsKey(params) : null);
  const swrFn = () => getGroups(params, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Create new group in the Cognito user pool if it is enabled.
 * @summary Create Group
 */
export const createGroup = (
    createGroupRequestBody: CreateGroupRequestBody,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<Group>(
      {url: `/api/v1/admin/groups`, method: 'post',
      headers: {'Content-Type': 'application/json', },
      data: createGroupRequestBody
    },
      options);
    }
  

/**
 * Returns information for a group.
 * @summary Get Group Details
 */
export const getGroup = (
    groupId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<Group>(
      {url: `/api/v1/admin/groups/${groupId}`, method: 'get'
    },
      options);
    }
  

export const getGetGroupKey = (groupId: string,) => [`/api/v1/admin/groups/${groupId}`];

    
export type GetGroupQueryResult = NonNullable<Awaited<ReturnType<typeof getGroup>>>
export type GetGroupQueryError = ErrorType<unknown>

export const useGetGroup = <TError = ErrorType<unknown>>(
 groupId: string, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof getGroup>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false && !!(groupId)
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getGetGroupKey(groupId) : null);
  const swrFn = () => getGroup(groupId, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * List providers
 * @summary List providers
 */
export const listProviders = (
    
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<Provider[]>(
      {url: `/api/v1/admin/providers`, method: 'get'
    },
      options);
    }
  

export const getListProvidersKey = () => [`/api/v1/admin/providers`];

    
export type ListProvidersQueryResult = NonNullable<Awaited<ReturnType<typeof listProviders>>>
export type ListProvidersQueryError = ErrorType<ErrorResponseResponse>

export const useListProviders = <TError = ErrorType<ErrorResponseResponse>>(
  options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof listProviders>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getListProvidersKey() : null);
  const swrFn = () => listProviders(requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Get provider by id
 * @summary List providers
 */
export const getProvider = (
    providerId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<Provider>(
      {url: `/api/v1/admin/providers/${providerId}`, method: 'get'
    },
      options);
    }
  

export const getGetProviderKey = (providerId: string,) => [`/api/v1/admin/providers/${providerId}`];

    
export type GetProviderQueryResult = NonNullable<Awaited<ReturnType<typeof getProvider>>>
export type GetProviderQueryError = ErrorType<ErrorResponseResponse>

export const useGetProvider = <TError = ErrorType<ErrorResponseResponse>>(
 providerId: string, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof getProvider>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false && !!(providerId)
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getGetProviderKey(providerId) : null);
  const swrFn = () => getProvider(providerId, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * gets the jsonschema describing the args for this provider
 * @summary Get provider arg schema
 */
export const getProviderArgs = (
    providerId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ArgSchemaResponseResponse>(
      {url: `/api/v1/admin/providers/${providerId}/args`, method: 'get'
    },
      options);
    }
  

export const getGetProviderArgsKey = (providerId: string,) => [`/api/v1/admin/providers/${providerId}/args`];

    
export type GetProviderArgsQueryResult = NonNullable<Awaited<ReturnType<typeof getProviderArgs>>>
export type GetProviderArgsQueryError = ErrorType<ErrorResponseResponse>

export const useGetProviderArgs = <TError = ErrorType<ErrorResponseResponse>>(
 providerId: string, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof getProviderArgs>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false && !!(providerId)
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getGetProviderArgsKey(providerId) : null);
  const swrFn = () => getProviderArgs(providerId, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Returns the options for a particular Access Provider argument. The options may be cached. To refresh the cache, pass the `refresh` query parameter.
 * @summary List provider arg options
 */
export const listProviderArgOptions = (
    providerId: string,
    argId: string,
    params?: ListProviderArgOptionsParams,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ArgOptionsResponseResponse>(
      {url: `/api/v1/admin/providers/${providerId}/args/${argId}/options`, method: 'get',
        params
    },
      options);
    }
  

export const getListProviderArgOptionsKey = (providerId: string,
    argId: string,
    params?: ListProviderArgOptionsParams,) => [`/api/v1/admin/providers/${providerId}/args/${argId}/options`, ...(params ? [params]: [])];

    
export type ListProviderArgOptionsQueryResult = NonNullable<Awaited<ReturnType<typeof listProviderArgOptions>>>
export type ListProviderArgOptionsQueryError = ErrorType<ErrorResponseResponse>

export const useListProviderArgOptions = <TError = ErrorType<ErrorResponseResponse>>(
 providerId: string,
    argId: string,
    params?: ListProviderArgOptionsParams, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof listProviderArgOptions>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false && !!(providerId && argId)
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getListProviderArgOptionsKey(providerId,argId,params) : null);
  const swrFn = () => listProviderArgOptions(providerId,argId,params, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * List providers which are still in the process of being set up.
 * @summary List the provider setups in progress
 */
export const listProvidersetups = (
    
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ListProviderSetupsResponseResponse>(
      {url: `/api/v1/admin/providersetups`, method: 'get'
    },
      options);
    }
  

export const getListProvidersetupsKey = () => [`/api/v1/admin/providersetups`];

    
export type ListProvidersetupsQueryResult = NonNullable<Awaited<ReturnType<typeof listProvidersetups>>>
export type ListProvidersetupsQueryError = ErrorType<unknown>

export const useListProvidersetups = <TError = ErrorType<unknown>>(
  options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof listProvidersetups>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getListProvidersetupsKey() : null);
  const swrFn = () => listProvidersetups(requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Get the setup instructions for an Access Provider.
 * @summary Get an in-progress provider setup
 */
export const getProvidersetup = (
    providersetupId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ProviderSetupResponseResponse>(
      {url: `/api/v1/admin/providersetups/${providersetupId}`, method: 'get'
    },
      options);
    }
  

export const getGetProvidersetupKey = (providersetupId: string,) => [`/api/v1/admin/providersetups/${providersetupId}`];

    
export type GetProvidersetupQueryResult = NonNullable<Awaited<ReturnType<typeof getProvidersetup>>>
export type GetProvidersetupQueryError = ErrorType<unknown>

export const useGetProvidersetup = <TError = ErrorType<unknown>>(
 providersetupId: string, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof getProvidersetup>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false && !!(providersetupId)
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getGetProvidersetupKey(providersetupId) : null);
  const swrFn = () => getProvidersetup(providersetupId, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Get the setup instructions for an Access Provider.
 * @summary Get the setup instructions for an Access Provider
 */
export const getProvidersetupInstructions = (
    providersetupId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ProviderSetupInstructions>(
      {url: `/api/v1/admin/providersetups/${providersetupId}/instructions`, method: 'get'
    },
      options);
    }
  

export const getGetProvidersetupInstructionsKey = (providersetupId: string,) => [`/api/v1/admin/providersetups/${providersetupId}/instructions`];

    
export type GetProvidersetupInstructionsQueryResult = NonNullable<Awaited<ReturnType<typeof getProvidersetupInstructions>>>
export type GetProvidersetupInstructionsQueryError = ErrorType<unknown>

export const useGetProvidersetupInstructions = <TError = ErrorType<unknown>>(
 providersetupId: string, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof getProvidersetupInstructions>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false && !!(providersetupId)
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getGetProvidersetupInstructionsKey(providersetupId) : null);
  const swrFn = () => getProvidersetupInstructions(providersetupId, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Validates the configuration values for an access provider being setup.

Will return a HTTP200 OK response even if there are validation errors. The errors can be found by inspecting the validation diagnostics in the `configValidation` field.

Will return a HTTP400 response if the provider cannot be validated (for example, the config values for the provider are incomplete).
 * @summary Validate the configuration for a Provider Setup
 */
export const validateProvidersetup = (
    providersetupId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ProviderSetupResponseResponse>(
      {url: `/api/v1/admin/providersetups/${providersetupId}/validate`, method: 'post'
    },
      options);
    }
  

/**
 * If Runtime Configuration is enabled, this will write the Access Provider to the configuration storage and activate it. If Runtime Configuration is disabled, this endpoint does nothing.
 * @summary Complete a ProviderSetup
 */
export const completeProvidersetup = (
    providersetupId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<CompleteProviderSetupResponseResponse>(
      {url: `/api/v1/admin/providersetups/${providersetupId}/complete`, method: 'post'
    },
      options);
    }
  

/**
 * The updated provider setup.
 * @summary Update the completion status for a Provider setup step
 */
export const submitProvidersetupStep = (
    providersetupId: string,
    stepIndex: number,
    providerSetupStepCompleteRequestBody: ProviderSetupStepCompleteRequestBody,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ProviderSetupResponseResponse>(
      {url: `/api/v1/admin/providersetups/${providersetupId}/steps/${stepIndex}/complete`, method: 'put',
      headers: {'Content-Type': 'application/json', },
      data: providerSetupStepCompleteRequestBody
    },
      options);
    }
  

/**
 * Run the identity sync operation on demand
 * @summary Sync Identity
 */
export const identitySync = (
    
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<void>(
      {url: `/api/v1/admin/identity/sync`, method: 'post'
    },
      options);
    }
  

/**
 * Get information about the identity configuration
 * @summary Get identity configuration
 */
export const identityConfiguration = (
    
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<IdentityConfigurationResponseResponse>(
      {url: `/api/v1/admin/identity`, method: 'get'
    },
      options);
    }
  

export const getIdentityConfigurationKey = () => [`/api/v1/admin/identity`];

    
export type IdentityConfigurationQueryResult = NonNullable<Awaited<ReturnType<typeof identityConfiguration>>>
export type IdentityConfigurationQueryError = ErrorType<ErrorResponseResponse>

export const useIdentityConfiguration = <TError = ErrorType<ErrorResponseResponse>>(
  options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof identityConfiguration>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getIdentityConfigurationKey() : null);
  const swrFn = () => identityConfiguration(requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

