import { orderBy } from "lodash";
import { TeamProps } from "props/team";
import { dateFormat } from "lib/date";
import CountryFlag from "./components/countryflag";
import ChartView from "./components/chart";
import styles from "./team.module.scss";

const Team = ({ team, scorefeed, series }: TeamProps) => {
  const tasks = orderBy(
    Object.entries(scorefeed.taskStats).map(([key, taskStat]) => ({
      name: key,
      ...taskStat,
    })),
    ["time"],
    ["desc"]
  );
  return (
    <>
      <h1>
        {team.country ? <CountryFlag country={team.country} /> : ""}
        {team.teamname}
        <span className={styles["team-info"]}>
          Rank {scorefeed.pos} / {scorefeed.score} points
        </span>
      </h1>

      <ChartView chartTeams={[team.teamname]} chartSeries={series} />

      {scorefeed && (
        <table className={styles["score-table"]}>
          <thead>
            <tr>
              <th>Task</th>
              <th>Score</th>
              <th>Solved At</th>
            </tr>
          </thead>
          <tbody>
            {tasks.map((task) => (
              <tr key={task.name}>
                <td>{task.name}</td>
                <td>{task.points}</td>
                <td>{dateFormat(task.time)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </>
  );
};

export default Team;
