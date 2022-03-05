import { ssrFetcher } from "lib/api";
import useSWR from "swr";

export interface Config {
  start_at: number;
  end_at: number;
  register_open: boolean;
  ctf_open: boolean;
  lock_second: number;
  lock_duration: number;
  lock_count: number;
  score_expr: string;
}

const useConfig = () => useSWR<Config>("/admin/get-config");

export const fetchConfig = () => ssrFetcher<Config>("/admin/get-config");

export default useConfig;
