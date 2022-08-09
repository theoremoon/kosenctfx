import CountryFlag from "../components/countryflag";
import Link from "next/link";
import { dateFormat } from "lib/date";
import { RankingProps } from "props/ranking";
import styles from "./ranking.module.scss";

type rankingProps = Pick<RankingProps, "tasks" | "scoreboard" | "account">;
const Ranking = ({ tasks, scoreboard, account }: rankingProps) => {
  return (
    <table className={styles["ranking-table"]}>
      <thead>
        <tr>
          <th>Rank</th>
          <th>Country</th>
          <th>Team</th>
          <th>Score</th>
          {tasks.map((t) => (
            <th key={t.id} className={styles["task-name"]}>
              <div>{t.name}</div>
            </th>
          ))}
        </tr>
      </thead>
      <tbody>
        {scoreboard.map((t) => (
          <tr
            key={t.team_id}
            style={{
              background:
                t.team_id === account?.team_id ? "rgba(213,63,140,0.2)" : "",
            }}
          >
            <td className={styles["team-pos"]}>{t.pos === 1 ? "🎂" : t.pos}</td>
            <td className={styles["team-country"]}>
              {t.country ? <CountryFlag country={t.country} /> : ""}
            </td>
            <td className={styles["team-name"]} title={t.team}>
              <Link href={"/teams/" + t.team_id}>{t.team}</Link>
            </td>
            <td className={styles["team-score"]}>{t.score}</td>
            {tasks.map((task) => {
              const solved = t.taskStats[task.name];
              const id = t.team_id.toString() + task.id.toString();
              if (solved) {
                return (
                  <td
                    className={styles["team-solve-flag"]}
                    key={id}
                    title={`${t.team} solved ${task.name} at ${dateFormat(
                      solved.time
                    )}`}
                  >
                    🍰
                  </td>
                );
              } else {
                return <td className={styles["team-solve-flag"]} key={id}></td>;
              }
            })}
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default Ranking;
