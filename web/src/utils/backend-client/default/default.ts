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
  ListRequestsResponseResponse,
  UserListRequestsUpcomingParams,
  UserListRequestsPastParams,
  AccessRuleDetail,
  ErrorResponseResponse,
  ProviderSetupResponseResponse,
  CreateProviderSetupRequestBody,
  LookupAccessRule,
  AccessRuleLookupParams,
  Favorite,
  CreateFavoriteRequestBody,
  FavoriteDetail
} from '.././types'
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
 * display pending requests and approved requests that are currently active or scheduled to begin some time in future.
 * @summary Your GET endpoint
 */
export const userListRequestsUpcoming = (
    params?: UserListRequestsUpcomingParams,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ListRequestsResponseResponse>(
      {url: `/api/v1/requests/upcoming`, method: 'get',
        params
    },
      options);
    }
  

export const getUserListRequestsUpcomingKey = (params?: UserListRequestsUpcomingParams,) => [`/api/v1/requests/upcoming`, ...(params ? [params]: [])];

    
export type UserListRequestsUpcomingQueryResult = NonNullable<Awaited<ReturnType<typeof userListRequestsUpcoming>>>
export type UserListRequestsUpcomingQueryError = ErrorType<unknown>

export const useUserListRequestsUpcoming = <TError = ErrorType<unknown>>(
 params?: UserListRequestsUpcomingParams, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof userListRequestsUpcoming>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getUserListRequestsUpcomingKey(params) : null);
  const swrFn = () => userListRequestsUpcoming(params, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * display show cancelled, expired, and revoked requests.

 * @summary Your GET endpoint
 */
export const userListRequestsPast = (
    params?: UserListRequestsPastParams,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ListRequestsResponseResponse>(
      {url: `/api/v1/requests/past`, method: 'get',
        params
    },
      options);
    }
  

export const getUserListRequestsPastKey = (params?: UserListRequestsPastParams,) => [`/api/v1/requests/past`, ...(params ? [params]: [])];

    
export type UserListRequestsPastQueryResult = NonNullable<Awaited<ReturnType<typeof userListRequestsPast>>>
export type UserListRequestsPastQueryError = ErrorType<unknown>

export const useUserListRequestsPast = <TError = ErrorType<unknown>>(
 params?: UserListRequestsPastParams, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof userListRequestsPast>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getUserListRequestsPastKey(params) : null);
  const swrFn = () => userListRequestsPast(params, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Marks an access rule as archived.
Any pending requests for this access rule will be cancelled.
 * @summary Archive Access Rule
 */
export const adminArchiveAccessRule = (
    ruleId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<AccessRuleDetail>(
      {url: `/api/v1/admin/access-rules/${ruleId}/archive`, method: 'post'
    },
      options);
    }
  

/**
 * Begins the guided setup process for a new Access Provider.
 * @summary Begin the setup process for a new Access Provider
 */
export const createProvidersetup = (
    createProviderSetupRequestBody: CreateProviderSetupRequestBody,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ProviderSetupResponseResponse>(
      {url: `/api/v1/admin/providersetups`, method: 'post',
      headers: {'Content-Type': 'application/json', },
      data: createProviderSetupRequestBody
    },
      options);
    }
  

/**
 * Removes an in-progress provider setup and deletes all data relating to it.

Returns the deleted provider.
 * @summary Delete an in-progress provider setup
 */
export const deleteProvidersetup = (
    providersetupId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ProviderSetupResponseResponse>(
      {url: `/api/v1/admin/providersetups/${providersetupId}`, method: 'delete'
    },
      options);
    }
  

/**
 * endpoint returns an array of relevant access rules (used in combination with granted cli)
 * @summary Lookup an access rule based on the target
 */
export const accessRuleLookup = (
    params?: AccessRuleLookupParams,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<LookupAccessRule[]>(
      {url: `/api/v1/access-rules/lookup`, method: 'get',
        params
    },
      options);
    }
  

export const getAccessRuleLookupKey = (params?: AccessRuleLookupParams,) => [`/api/v1/access-rules/lookup`, ...(params ? [params]: [])];

    
export type AccessRuleLookupQueryResult = NonNullable<Awaited<ReturnType<typeof accessRuleLookup>>>
export type AccessRuleLookupQueryError = ErrorType<ErrorResponseResponse>

export const useAccessRuleLookup = <TError = ErrorType<ErrorResponseResponse>>(
 params?: AccessRuleLookupParams, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof accessRuleLookup>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getAccessRuleLookupKey(params) : null);
  const swrFn = () => accessRuleLookup(params, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * @summary ListFavorites
 */
export const userListFavorites = (
    
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<Favorite[]>(
      {url: `/api/v1/favorites`, method: 'get'
    },
      options);
    }
  

export const getUserListFavoritesKey = () => [`/api/v1/favorites`];

    
export type UserListFavoritesQueryResult = NonNullable<Awaited<ReturnType<typeof userListFavorites>>>
export type UserListFavoritesQueryError = ErrorType<ErrorResponseResponse>

export const useUserListFavorites = <TError = ErrorType<ErrorResponseResponse>>(
  options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof userListFavorites>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getUserListFavoritesKey() : null);
  const swrFn = () => userListFavorites(requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * @summary Create Favorite
 */
export const userCreateFavorite = (
    createFavoriteRequestBody: CreateFavoriteRequestBody,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<Favorite>(
      {url: `/api/v1/favorites`, method: 'post',
      headers: {'Content-Type': 'application/json', },
      data: createFavoriteRequestBody
    },
      options);
    }
  

/**
 * @summary Get Favorite
 */
export const userGetFavorite = (
    id: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<FavoriteDetail>(
      {url: `/api/v1/favorites/${id}`, method: 'get'
    },
      options);
    }
  

export const getUserGetFavoriteKey = (id: string,) => [`/api/v1/favorites/${id}`];

    
export type UserGetFavoriteQueryResult = NonNullable<Awaited<ReturnType<typeof userGetFavorite>>>
export type UserGetFavoriteQueryError = ErrorType<ErrorResponseResponse>

export const useUserGetFavorite = <TError = ErrorType<ErrorResponseResponse>>(
 id: string, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof userGetFavorite>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options ?? {}

  const isEnabled = swrOptions?.enabled !== false && !!(id)
    const swrKey = swrOptions?.swrKey ?? (() => isEnabled ? getUserGetFavoriteKey(id) : null);
  const swrFn = () => userGetFavorite(id, requestOptions);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

