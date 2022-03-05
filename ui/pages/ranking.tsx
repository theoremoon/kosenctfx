import {
  Box,
  Icon,
  Link,
  Switch,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { faFlag, faMedal } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import SeriesChart from "components/chart";
import CountryFlag from "components/countryflag";
import Loading from "components/loading";
import { fetchSeries, SeriesEntry } from "lib/api/series";
import { white } from "lib/color";
import { isStaticMode } from "lib/static";
import { orderBy } from "lodash";
import { GetStaticProps } from "next";
import NextLink from "next/link";
import { useLocalStorage } from "usehooks-ts";
import useAccount, { Account, fetchAccount } from "../lib/api/account";
import useScoreboard, {
  fetchScoreboard,
  ScoreFeedEntry,
} from "../lib/api/scoreboard";
import useTasks, { fetchTasks, Task } from "../lib/api/tasks";
import { dateFormat } from "../lib/date";

interface RankingProps {
  scoreboard: ScoreFeedEntry[];
  tasks: Task[];
  account: Account | null;
  series: SeriesEntry[][];
}

const Ranking = ({
  scoreboard: scoreboardDefault,
  tasks: tasksDefault,
  account: defaultAccount,
  series: defaultSeries,
}: RankingProps) => {
  const { data: account } = useAccount(defaultAccount);
  const { data: scoreboard } = useScoreboard(scoreboardDefault);
  const { data: tasks } = useTasks(tasksDefault);
  const [selectedTeams, setSelectedTeams] = useLocalStorage<string[]>(
    "seriesTeamName",
    []
  );

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
  const topTeams = scoreboard.slice(0, 10).map((t) => t.team);
  const chartTeams = new Set<string>([...topTeams, ...selectedTeams]);
  if (account) {
    chartTeams.add(account.teamname);
  }

  return (
    <>
      <Box mt={10}>
        <SeriesChart teams={Array.from(chartTeams)} series={defaultSeries} />
      </Box>
      <Table size="sm" mt={20}>
        <Thead>
          <Tr>
            <Th>Rank</Th>
            <Th>Country</Th>
            {!isStaticMode && <Th>Chart</Th>}
            <Th>Team</Th>
            <Th>Score</Th>
            {orderedTasks.map((t) => (
              <Th
                key={t.name}
                sx={{
                  whiteSpace: "pre",
                  position: "relative",
                  borderRight: "none",
                  div: {
                    position: "absolute",
                    transformOrigin: "left",
                    transform: "translate(0.5em, -0.5em) rotate(-30deg)",
                  },
                }}
              >
                <div>{t.name}</div>
              </Th>
            ))}
          </Tr>
        </Thead>
        <Tbody>
          {scoreboard.map((t) => (
            <Tr
              key={t.team_id}
              sx={{
                background:
                  t.team_id === account?.team_id
                    ? "rgba(255,255,255, 0.2)"
                    : "",
              }}
            >
              <Td>
                {t.pos === 1 ? (
                  <Icon as={FontAwesomeIcon} icon={faMedal} />
                ) : (
                  t.pos
                )}
              </Td>
              <Td>
                <CountryFlag country={t.country} />
              </Td>
              {!isStaticMode && (
                <Td>
                  <Switch
                    isChecked={selectedTeams.includes(t.team)}
                    onChange={() => {
                      if (selectedTeams.includes(t.team)) {
                        setSelectedTeams((past) =>
                          past.filter((team) => team !== t.team)
                        );
                      } else {
                        setSelectedTeams((past) => [...past, t.team]);
                      }
                    }}
                  />
                </Td>
              )}
              <Td
                title={t.team}
                sx={{
                  maxW: "15em",
                  textOverflow: "ellipsis",
                  overflow: "hidden",
                  whiteSpace: "nowrap",
                }}
              >
                <Link as={NextLink} href={`/teams/${t.team_id}`}>
                  {t.team}
                </Link>
              </Td>
              <Td>{t.score}</Td>
              {orderedTasks.map((task) => {
                let flag: JSX.Element | string = " ";
                const solved = t.taskStats[task.name];
                if (solved) {
                  flag = (
                    <Box position={"absolute"}>
                      <Icon
                        as={FontAwesomeIcon}
                        icon={faFlag}
                        position={"relative"}
                        sx={{
                          top: -2,
                          left: -2,
                        }}
                        title={`${t.team} solved ${task.name} at ${dateFormat(
                          solved.time
                        )}`}
                      />
                    </Box>
                  );
                }
                return <Td key={"" + t.team_id + task.id}>{flag}</Td>;
              })}
            </Tr>
          ))}
        </Tbody>
      </Table>
    </>
  );
};

export const getStaticProps: GetStaticProps<RankingProps> = async () => {
  const scoreboard = await fetchScoreboard();
  const tasks = await fetchTasks();
  const account = isStaticMode ? null : await fetchAccount();
  const topTeams = scoreboard.slice(0, 10).map((t) => t.team);
  const series = await fetchSeries(topTeams);
  return {
    props: {
      account: account,
      scoreboard: scoreboard,
      tasks: tasks,
      series: series,
    },
  };
};

export default Ranking;
