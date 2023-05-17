/**
 * Generated by orval v6.7.1 🍺
 * Do not edit manually.
 * Common Fate
 * Common Fate API
 * OpenAPI spec version: 1.0
 */
import useSwr,{
  SWRConfiguration,
  Key
} from 'swr'
import type {
  User,
  AuthUserResponseResponse,
  ReviewResponseResponse,
  ReviewRequestBody,
  ErrorResponseResponse,
  ListRequestEventsResponseResponse
} from '.././types'
import { customInstance, ErrorType } from '../../custom-instance'

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type AsyncReturnType<
T extends (...args: any) => Promise<any>
> = T extends (...args: any) => Promise<infer R> ? R : any;


// eslint-disable-next-line @typescript-eslint/no-explicit-any
  type SecondParameter<T extends (...args: any) => any> = T extends (
  config: any,
  args: infer P,
) => any
  ? P
  : never;

/**
 * Returns a Common Fate user.
 * @summary Get a user
 */
export const userGetUser = (
    userId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<User>(
      {url: `/api/v1/users/${userId}`, method: 'get'
    },
      options);
    }
  

export const getUserGetUserKey = (userId: string,) => [`/api/v1/users/${userId}`];

    
export type UserGetUserQueryResult = NonNullable<AsyncReturnType<typeof userGetUser>>
export type UserGetUserQueryError = ErrorType<void>

export const useUserGetUser = <TError = ErrorType<void>>(
 userId: string, options?: { swr?:SWRConfiguration<AsyncReturnType<typeof userGetUser>, TError> & {swrKey: Key}, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options || {}

  const isEnable = !!(userId)
  const swrKey = swrOptions?.swrKey ?? (() => isEnable ? getUserGetUserKey(userId) : null);
  const swrFn = () => userGetUser(userId, requestOptions);

  const query = useSwr<AsyncReturnType<typeof swrFn>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Returns information about the currently logged in user.
 * @summary Get details for the current user
 */
export const userGetMe = (
    
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<AuthUserResponseResponse>(
      {url: `/api/v1/users/me`, method: 'get'
    },
      options);
    }
  

export const getUserGetMeKey = () => [`/api/v1/users/me`];

    
export type UserGetMeQueryResult = NonNullable<AsyncReturnType<typeof userGetMe>>
export type UserGetMeQueryError = ErrorType<void>

export const useUserGetMe = <TError = ErrorType<void>>(
  options?: { swr?:SWRConfiguration<AsyncReturnType<typeof userGetMe>, TError> & {swrKey: Key}, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options || {}

  const swrKey = swrOptions?.swrKey ?? (() => getUserGetMeKey())
  const swrFn = () => userGetMe(requestOptions);

  const query = useSwr<AsyncReturnType<typeof swrFn>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

/**
 * Review an access request made by a user. The reviewing user must be an approver for a request. Users cannot review their own requests, even if they are an approver for the Access Rule.
 * @summary Review a request
 */
export const userReviewRequest = (
    requestId: string,
    groupId: string,
    reviewRequestBody: ReviewRequestBody,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ReviewResponseResponse>(
      {url: `/api/v1/requests/${requestId}/review/${groupId}`, method: 'post',
      headers: {'Content-Type': 'application/json'},
      data: reviewRequestBody
    },
      options);
    }
  

/**
 * Admins and approvers can revoke access previously approved. Effective immediately 
 * @summary Revoke an active request
 */
export const userRevokeRequest = (
    requestid: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<void>(
      {url: `/api/v1/requests/${requestid}/revoke`, method: 'post'
    },
      options);
    }
  

/**
 * Admins and approvers can cancel access before provisioned. Effective immediately 
 * @summary Revoke an active request
 */
export const userCancelRequest = (
    requestid: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<void>(
      {url: `/api/v1/requests/${requestid}/cancel`, method: 'post'
    },
      options);
    }
  

/**
 * Returns a HTTP401 response if the user is not the requestor or a reviewer.

 * @summary List request events
 */
export const userListRequestEvents = (
    requestId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<ListRequestEventsResponseResponse>(
      {url: `/api/v1/requests/${requestId}/events`, method: 'get'
    },
      options);
    }
  

export const getUserListRequestEventsKey = (requestId: string,) => [`/api/v1/requests/${requestId}/events`];

    
export type UserListRequestEventsQueryResult = NonNullable<AsyncReturnType<typeof userListRequestEvents>>
export type UserListRequestEventsQueryError = ErrorType<ErrorResponseResponse>

export const useUserListRequestEvents = <TError = ErrorType<ErrorResponseResponse>>(
 requestId: string, options?: { swr?:SWRConfiguration<AsyncReturnType<typeof userListRequestEvents>, TError> & {swrKey: Key}, request?: SecondParameter<typeof customInstance> }

  ) => {

  const {swr: swrOptions, request: requestOptions} = options || {}

  const isEnable = !!(requestId)
  const swrKey = swrOptions?.swrKey ?? (() => isEnable ? getUserListRequestEventsKey(requestId) : null);
  const swrFn = () => userListRequestEvents(requestId, requestOptions);

  const query = useSwr<AsyncReturnType<typeof swrFn>, TError>(swrKey, swrFn, swrOptions)

  return {
    swrKey,
    ...query
  }
}

