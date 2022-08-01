import { GetStaticProps } from "next";
import { orderBy } from "lodash";

import Loading from "components/loading";
import { AllPageProps } from "lib/pages";
import { isStaticMode } from "lib/static";
import useAccount, { fetchAccount } from "lib/api/account";
import { fetchCTF } from "lib/api/ctf";
import useSeries, { fetchSeries } from "lib/api/series";
import useScoreboard, { fetchScoreboard } from "lib/api/scoreboard";
import useTasks, { fetchTasks, Task } from "lib/api/tasks";
import RankingView from "theme/ranking";
import { RankingProps } from "props/ranking";

type rankingProps = Omit<RankingProps & AllPageProps, "chartTeams">;
const Ranking = ({
  scoreboard: scoreboardDefault,
  tasks: tasksDefault,
  account: defaultAccount,
  chartSeries: defaultSeries,
}: rankingProps) => {
  const { data: account } = useAccount(defaultAccount);
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
  const account = isStaticMode ? null : await fetchAccount().catch(() => null);
  const topTeams = scoreboard.slice(0, 10).map((t) => t.team);
  const series = await fetchSeries(topTeams).catch(() => []);
  const ctf = await fetchCTF();
  return {
    props: {
      account: account,
      scoreboard: scoreboard,
      tasks: tasks,
      chartSeries: series,
      ctf: ctf,
    },
  };
};

export default Ranking;
