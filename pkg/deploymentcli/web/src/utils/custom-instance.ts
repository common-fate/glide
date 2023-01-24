// https://orval.dev/guides/custom-client
// https://orval.dev/reference/configuration/output
// custom-instance.ts
import Axios, { AxiosError, AxiosRequestConfig } from "axios";

export const customInstanceLocal = async <T>(
  config: AxiosRequestConfig,
  runtimeConfig?: AxiosRequestConfig
): Promise<T> => {
  const instance = Axios.create();
  const promise = instance({
    baseURL: "http://localhost:9000",
    ...config,
    ...runtimeConfig,
  }).then(({ data }) => data);

  return promise;
};

export const customInstanceRegistry = async <T>(
  config: AxiosRequestConfig,
  runtimeConfig?: AxiosRequestConfig
): Promise<T> => {
  const instance = Axios.create();
  const promise = instance({
    baseURL: "http://localhost:9001",
    ...config,
    ...runtimeConfig,
  }).then(({ data }) => data);

  return promise;
};
// In some case with react-query and swr you want to be able to override the return error type so you can also do it here like this
export type ErrorType<Error> = AxiosError<Error>;
// In case you want to wrap the body type (optional)
// (if the custom instance is processing data before sending it, like changing the case for example)
