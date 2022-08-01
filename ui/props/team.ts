import { ScoreFeedEntry } from "lib/api/scoreboard";
import { SeriesEntry } from "lib/api/series";
import { Team } from "lib/api/team";

export interface TeamProps {
  team: Team;
  scorefeed: ScoreFeedEntry;
  series: SeriesEntry[][];
}
