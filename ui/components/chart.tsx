import ReactECharts from "echarts-for-react";
import { white } from "lib/color";
import useSeries from "../lib/api/series";
import Loading from "./loading";

interface SeriesChartProps {
  teams: string[];
}

const SeriesChart = ({ teams }: SeriesChartProps) => {
  const { data } = useSeries(teams);

  if (teams.length === 0) {
    return <></>;
  }

  if (!data) {
    return <Loading />;
  }

  const series = data.map((team, i) => ({
    name: teams[i],
    type: "line",
    showSymbol: false,
    data: team.map((e) => [e.time * 1000, e.score]),
  }));
  const labels = teams.map((team) => ({
    name: team,
    textStyle: {
      color: white,
    },
  }));
  return (
    <>
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
    </>
  );
};

export default SeriesChart;
