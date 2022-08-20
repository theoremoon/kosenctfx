import useSWR from "swr";
import { Agent } from "./agent";
import { Task } from "./tasks";

export interface Deployment {
  id: number;
  challenge: Task | null;
  agent: Agent | null;
  port: number;
  status: string;

  requested_at: number;
  retires_at: number;
}

const useLivingDeployments = () =>
  useSWR<Deployment[]>("/admin/list-living-deployments");

export default useLivingDeployments;
