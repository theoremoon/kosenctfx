import { Task } from "lib/api/tasks";
import { ScoreFeedEntry } from "lib/api/scoreboard";
import { Account } from "lib/api/account";
import { SeriesEntry } from "lib/api/series";

export interface RankingProps {
  scoreboard: ScoreFeedEntry[];
  tasks: Task[];
  account: Account | null;

  chartTeams: string[];
  chartSeries: SeriesEntry[][];
}
