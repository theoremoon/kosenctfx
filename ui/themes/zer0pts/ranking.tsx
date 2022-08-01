import { Box } from "@chakra-ui/react";

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
      <Box mt={10}>
        <ChartView chartTeams={teams} chartSeries={chartSeries} />
      </Box>
      <RankingView tasks={tasks} scoreboard={scoreboard} account={account} />
    </>
  );
};

export default Ranking;
