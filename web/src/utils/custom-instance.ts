// https://orval.dev/guides/custom-client
// https://orval.dev/reference/configuration/output
// custom-instance.ts
import { Auth } from "@aws-amplify/auth";
import Axios, { AxiosError, AxiosRequestConfig } from "axios";

let apiURL: string;

/**
 * Used to update the API URL after we've fetched our runtime configuration.
 */
export const setAPIURL = (url: string) => {
  apiURL = url;
};

export const customInstance = async <T>(
  config: AxiosRequestConfig,
  runtimeConfig?: AxiosRequestConfig
): Promise<T> => {
  const instance = Axios.create();

  const session = await Auth.currentSession();
  const token = session.getIdToken().getJwtToken();

  instance.interceptors.request.use(
    async (config) => {
      if (token && config?.headers) {
        config.headers.Authorization = token;
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  const baseURL = apiURL;

  const promise = instance({
    baseURL,
    headers: {
      ...config.headers,
      Authorization: token,
    },
    ...config,
    ...runtimeConfig,
  }).then(({ data }) => data);

  return promise;
};

// In some case with react-query and swr you want to be able to override the return error type so you can also do it here like this
export type ErrorType<Error> = AxiosError<Error>;
// In case you want to wrap the body type (optional)
// (if the custom instance is processing data before sending it, like changing the case for example)
