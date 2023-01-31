/**
 * Generated by orval v6.11.1 🍺
 * Do not edit manually.
 * Example API
 * Example API
 * OpenAPI spec version: 1.0
 */
import useSwr from 'swr'
import type {
  SWRConfiguration,
  Key
} from 'swr'
import type {
  ProviderV2,
  ErrorResponseResponse
} from '.././types/openapi.yml'
import type {
  CreateProviderDeployment,
  DeleteProvider200,
  UpdateProviderDeployment
} from '.././types'
import { customInstanceLocal } from '../../custom-instance';
import type { ErrorType } from '../../custom-instance';


  
  // eslint-disable-next-line
  type SecondParameter<T extends (...args: any) => any> = T extends (
  config: any,
  args: infer P,
) => any
  ? P
  : never;

/**
 * Lists the Providers installed to an org's Common Fate account
 * @summary List providers
 */
export const listProviders = (
    
 options?: SecondParameter<typeof customInstanceLocal>) => {
      return customInstanceLocal<ProviderV2[]>(
      {url: `/api/v1/providers`, method: 'get'
    },
      options);
    }
  

export const getListProvidersKey = () => [`/api/v1/providers`];

    
export type ListProvidersQueryResult = NonNullable<Awaited<ReturnType<typeof listProviders>>>
export type ListProvidersQueryError = ErrorType<ErrorResponseResponse>

export const useListProviders = <TError = ErrorType<ErrorResponseResponse>>(
  options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof listProviders>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstanceLocal> }

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
 * @summary Create provider
 */
export const createProvider = (
    createProviderDeployment: CreateProviderDeployment,
 options?: SecondParameter<typeof customInstanceLocal>) => {
      return customInstanceLocal<void>(
      {url: `/api/v1/providers`, method: 'post',
      headers: {'Content-Type': 'application/json', },
      data: createProviderDeployment
    },
      options);
    }
  

/**
 * Get provider by id
 * @summary Get provider detailed
 */
export const getProvider = (
    providerId: string,
 options?: SecondParameter<typeof customInstanceLocal>) => {
      return customInstanceLocal<ProviderV2>(
      {url: `/api/v1/providers/${providerId}`, method: 'get'
    },
      options);
    }
  

export const getGetProviderKey = (providerId: string,) => [`/api/v1/providers/${providerId}`];

    
export type GetProviderQueryResult = NonNullable<Awaited<ReturnType<typeof getProvider>>>
export type GetProviderQueryError = ErrorType<ErrorResponseResponse>

export const useGetProvider = <TError = ErrorType<ErrorResponseResponse>>(
 providerId: string, options?: { swr?:SWRConfiguration<Awaited<ReturnType<typeof getProvider>>, TError> & { swrKey?: Key, enabled?: boolean }, request?: SecondParameter<typeof customInstanceLocal> }

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
 * @summary Delete provider
 */
export const deleteProvider = (
    providerId: string,
 options?: SecondParameter<typeof customInstanceLocal>) => {
      return customInstanceLocal<DeleteProvider200>(
      {url: `/api/v1/providers/${providerId}`, method: 'delete'
    },
      options);
    }
  

/**
 * @summary Update provider
 */
export const updateProvider = (
    providerId: string,
    updateProviderDeployment: UpdateProviderDeployment,
 options?: SecondParameter<typeof customInstanceLocal>) => {
      return customInstanceLocal<ProviderV2 | void>(
      {url: `/api/v1/providers/${providerId}`, method: 'post',
      headers: {'Content-Type': 'application/json', },
      data: updateProviderDeployment
    },
      options);
    }
  

