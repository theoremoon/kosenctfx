import axios, { AxiosError } from "axios";

export interface APIError {
  message: string;
}

export const api = axios.create({
  baseURL: "/api",
});

/// swrから使うためにglobalなfetcher
export const defaultFetcher = (url: string) =>
  api.get(url).then((res) => res.data);

type Message = {
  message: string;
};

export type MessageResponse = AxiosError<Message>;
export type ErrorResponse = AxiosError<Message>;
