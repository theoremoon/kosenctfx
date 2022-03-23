import { api } from "lib/api";
import useSWR from "swr";

export interface Team {
  teamname: string;
  team_id: number;
  country: string;
}

const useTeam = (teamid: string) =>
  useSWR<Team>(teamid, (team) => {
    return api.get<Team>(`/team/${team}`).then((r) => r.data);
  });
export default useTeam;
