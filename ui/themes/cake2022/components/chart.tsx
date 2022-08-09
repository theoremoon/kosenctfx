import ReactECharts from "echarts-for-react";
import { RankingProps } from "props/ranking";

type chartProps = Pick<RankingProps, "chartTeams" | "chartSeries">;

const Chart = ({ chartTeams: teams, chartSeries }: chartProps) => {
  const series = chartSeries.map((team, i) => ({
    name: teams[i],
    type: "line",
    showSymbol: false,
    data: team.map((e) => [e.time * 1000, e.score]),
  }));
  const labels = teams.map((team) => ({
    name: team,
  }));
  return (
    <ReactECharts
      option={{
        tooltip: {
          trigger: "axis",
          axisPointer: {
            type: "cross",
            animation: "false",
          },
        },
        legend: {
          data: labels,
        },
        xAxis: {
          type: "time",
        },
        yAxis: {
          type: "value",
        },
        series: series,
      }}
      notMerge={true}
    />
  );
};

export default Chart;
