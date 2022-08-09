import { RankingProps } from "props/ranking";
import RankingView from "./components/ranking";
import ChartView from "./components/chart";

const Ranking = ({
  tasks,
  scoreboard,
  account,
  chartTeams: teams,
  chartSeries,
}: RankingProps) => {
  return (
    <>
      <div style={{ marginTop: "20px" }}>
        <ChartView chartTeams={teams} chartSeries={chartSeries} />
      </div>
      <RankingView tasks={tasks} scoreboard={scoreboard} account={account} />
    </>
  );
};

export default Ranking;
