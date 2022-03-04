import {
  Box,
  Center,
  Icon,
  Link,
  Switch,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { faFlag, faMedal } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import SeriesChart from "components/chart";
import CountryFlag from "components/countryflag";
import Loading from "components/loading";
import { white } from "lib/color";
import { orderBy } from "lodash";
import NextLink from "next/link";
import { useLocalStorage } from "usehooks-ts";
import useAccount from "../lib/api/account";
import useScoreboard, { ScoreFeedEntry } from "../lib/api/scoreboard";
import useTasks, { Task } from "../lib/api/tasks";
import { dateFormat } from "../lib/date";

interface RankingProps {
  scoreboard: ScoreFeedEntry[];
  tasks: Task[];
}

const Ranking = ({ scoreboard, tasks }: RankingProps) => {
  const { data: account } = useAccount();
  const orderedTasks = orderBy(
    tasks,
    [
      (t: Task) => t.category,
      (t: Task) => t.score,
      (t: Task) => t.solved_by.length,
    ],
    ["asc", "asc", "desc"]
  );
  const initialTeams = scoreboard.slice(0, 10).map((t) => t.team);
  if (account) {
    initialTeams.push(account.teamname);
  }

  const [seriesTeam, setSeriesTeam] = useLocalStorage(
    "seriesTeamName",
    initialTeams
  );

  return (
    <>
      <Box mt={20}>
        <SeriesChart teams={seriesTeam} />
      </Box>
      <Table size="sm" mt={20}>
        <Thead>
          <Tr>
            <Th>Rank</Th>
            <Th>Country</Th>
            <Th>Chart</Th>
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
              <Td>
                <Switch
                  isChecked={seriesTeam.includes(t.team)}
                  onChange={() => {
                    if (seriesTeam.includes(t.team)) {
                      setSeriesTeam((past) =>
                        past.filter((team) => team !== t.team)
                      );
                    } else {
                      setSeriesTeam((past) => [...past, t.team]);
                    }
                  }}
                />
              </Td>
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
                return (
                  <Td
                    key={"" + t.team_id + task.id}
                    sx={{
                      borderLeftColor: white,
                      borderRightColor: white,
                      borderLeftWidth: "1px",
                      borderRightWidth: "1px",
                    }}
                  >
                    {flag}
                  </Td>
                );
              })}
            </Tr>
          ))}
        </Tbody>
      </Table>
    </>
  );
};

const RankingDefault = () => {
  const { data: scoreboard } = useScoreboard();
  const { data: tasks } = useTasks();
  if (scoreboard === undefined) {
    return <Loading />;
  }
  return <Ranking scoreboard={scoreboard} tasks={tasks || []} />;
};
export default RankingDefault;
