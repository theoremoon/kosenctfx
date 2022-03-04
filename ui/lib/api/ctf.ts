import { makeSWRResponse, ssrFetcher } from "lib/api";
import { isStaticMode } from "lib/static";
import useSWR from "swr";

export interface CTF {
  start_at: number;
  end_at: number;
  register_open: boolean;
  is_open: boolean;
  is_running: boolean;
  is_over: boolean;
}

const useCTF = (fallback: CTF) => {
  return isStaticMode
    ? makeSWRResponse(fallback)
    : useSWR<CTF>("/ctf", null, {
        fallbackData: fallback,
      });
};
export const fetchCTF = () => ssrFetcher<CTF>("/ctf");

export default useCTF;
