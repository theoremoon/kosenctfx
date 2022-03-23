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

const useScoreboard = () => useSWR<ScoreFeedEntry[]>("/scoreboard");

export default useScoreboard;
