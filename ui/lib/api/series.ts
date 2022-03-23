import { api } from "lib/api";
import useSWR from "swr";

export interface SeriesEntry {
  teamname: string;
  score: number;
  pos: number;
  time: number;
}

const useSeries = (teams: string[]) =>
  useSWR<SeriesEntry[][]>(teams.join(""), (url) => {
    return api
      .post<SeriesEntry[][]>("/series", {
        teams: teams,
      })
      .then((r) => r.data);
  });
export default useSeries;
