import { makeSWRResponse, ssrFetcher } from "lib/api";
import { isStaticMode } from "lib/static";
import useSWR from "swr";

export interface TaskStat {
  points: number;
  time: number;
}

export interface ScoreFeedEntry {
  pos: number;
  team: string;
  country: string;
  score: number;
  taskStats: { [taskName: string]: TaskStat };
  team_id: number;
  last_submission: number;
}

const useScoreboard = (staticValue: ScoreFeedEntry[]) => {
  return isStaticMode
    ? makeSWRResponse(staticValue)
    : useSWR<ScoreFeedEntry[]>("/scoreboard");
};

export const fetchScoreboard = () =>
  ssrFetcher<ScoreFeedEntry[]>("/scoreboard");
export default useScoreboard;
