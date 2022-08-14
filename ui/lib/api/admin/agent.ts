import { ssrFetcher } from "lib/api";
import useSWR from "swr";

export interface Agent {
  agent_id: string;
  public_ip: string;
  last_activity_at: number;
}

const useAgents = () => useSWR<Agent[]>("/admin/list-agents");

export const fetchAgents = () => ssrFetcher<Agent[]>("/admin/list-agents");

export default useAgents;
