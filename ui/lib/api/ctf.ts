import useSWR from "swr";

export interface CTF {
  start_at: number;
  end_at: number;
  register_open: boolean;
  is_open: boolean;
  is_running: boolean;
  is_over: boolean;
}

const useCTF = () => useSWR<CTF>("/ctf");

export default useCTF;
