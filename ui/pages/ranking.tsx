import { GetStaticProps } from "next";
import { orderBy } from "lodash";

import Loading from "components/loading";
import { AllPageProps } from "lib/pages";
import useAccount from "lib/api/account";
import { fetchCTF } from "lib/api/ctf";
import useSeries, { fetchSeries } from "lib/api/series";
import useScoreboard, { fetchScoreboard } from "lib/api/scoreboard";
import useTasks, { fetchTasks, Task } from "lib/api/tasks";
import RankingView from "theme/ranking";
import { RankingProps } from "props/ranking";
import { isStaticMode, revalidateInterval } from "lib/static";

type rankingProps = Omit<RankingProps & AllPageProps, "chartTeams" | "account">;
const Ranking = ({
  scoreboard: scoreboardDefault,
  tasks: tasksDefault,
  chartSeries: defaultSeries,
}: rankingProps) => {
  const { data: account } = useAccount(null);
  const { data: scoreboard } = useScoreboard(scoreboardDefault);
  const { data: tasks } = useTasks(tasksDefault);

  const chartTeams = scoreboard?.slice(0, 10).map((t) => t.team) || [];
  if (account) {
    chartTeams.push(account.teamname);
  }
  const { data: series } = useSeries(chartTeams, defaultSeries);

  if (scoreboard === undefined) {
    return <Loading />;
  }

  const orderedTasks = orderBy(
    tasks,
    [
      (t: Task) => t.category,
      (t: Task) => t.score,
      (t: Task) => t.solved_by.length,
    ],
    ["asc", "asc", "desc"]
  );

  return RankingView({
    tasks: orderedTasks,
    scoreboard,
    account: account || null,
    chartTeams,
    chartSeries: series || defaultSeries,
  });
};

export const getStaticProps: GetStaticProps<rankingProps> = async () => {
  const scoreboard = await fetchScoreboard().catch(() => []);
  const tasks = await fetchTasks().catch(() => []);
  const topTeams = scoreboard.slice(0, 10).map((t) => t.team);
  const series = await fetchSeries(topTeams).catch(() => []);
  const ctf = await fetchCTF();
  return {
    props: {
      scoreboard: scoreboard,
      tasks: tasks,
      chartSeries: series,
      ctf: ctf,
    },
    // staticModeでないときはクライアントのSWRが使われるはずなので、情報の更新間隔はクライアントのSWRのrevalidate間隔に従う
    revalidate: isStaticMode ? false : revalidateInterval,
  };
};

export default Ranking;
