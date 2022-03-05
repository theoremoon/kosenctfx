import { api, makeSWRResponse, ssrFetcher } from "lib/api";
import { isStaticMode } from "lib/static";
import useSWR from "swr";

export interface Team {
  teamname: string;
  team_id: number;
  country: string;
}

const useTeam = (teamID: string, staticValue: Team) => {
  return isStaticMode
    ? makeSWRResponse(staticValue)
    : useSWR<Team>(`/team/${teamID}`);
};

export const fetchTeam = (teamID: string) =>
  ssrFetcher<Team>(`/team/${teamID}`);

export default useTeam;
