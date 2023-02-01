/**
 * Generated by orval v6.10.3 🍺
 * Do not edit manually.
 * Common Fate
 * Common Fate API
 * OpenAPI spec version: 1.0
 */
import type {
  UpdateProviderV2
} from '.././types'
import { customInstance } from '../../custom-instance'


  
  // eslint-disable-next-line
  type SecondParameter<T extends (...args: any) => any> = T extends (
  config: any,
  args: infer P,
) => any
  ? P
  : never;

/**
 * @summary Delete providerv2
 */
export const adminDeleteProviderv2 = (
    providerId: string,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<void>(
      {url: `/api/v1/admin/providersv2/${providerId}`, method: 'delete'
    },
      options);
    }
  

/**
 * @summary Update providerv2
 */
export const adminUpdateProviderv2 = (
    providerId: string,
    updateProviderV2: UpdateProviderV2,
 options?: SecondParameter<typeof customInstance>) => {
      return customInstance<void>(
      {url: `/api/v1/admin/providersv2/${providerId}`, method: 'post',
      headers: {'Content-Type': 'application/json', },
      data: updateProviderV2
    },
      options);
    }
  

