import axios, { AxiosError } from "axios";
import { SWRResponse } from "swr";

export interface APIError {
  message: string;
}

export const api = axios.create({
  baseURL: "/api",
});

export const ssrApi = axios.create({
  baseURL: "http://nginx:80/api/",
})

/// swrから使うためにglobalなfetcher
export const defaultFetcher = (url: string) =>
  api.get(url).then((res) => res.data);

/// ssrで使うためのfetcher
export const ssrFetcher = <T>(url: string) =>
  ssrApi.get(url).then((res) => res.data as T);

export const makeSWRResponse = <T, E = any>(data: T) => ({
  data: data,
  error: undefined,
  mutate: () => {},
  isValidating: false,
} as SWRResponse<T, E>);

type Message = {
  message: string;
};

export type MessageResponse = AxiosError<Message>;
export type ErrorResponse = AxiosError<Message>;
