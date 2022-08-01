import {
  Box,
  Icon,
  Link,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { faFlag, faMedal } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import CountryFlag from "../components/countryflag";
import NextLink from "next/link";
import { dateFormat } from "lib/date";
import { RankingProps } from "props/ranking";

type rankingProps = Pick<RankingProps, "tasks" | "scoreboard" | "account">;
const Ranking = ({ tasks, scoreboard, account }: rankingProps) => {
  return (
    <Table size="sm" mt={20}>
      <Thead>
        <Tr>
          <Th>Rank</Th>
          <Th>Country</Th>
          <Th>Team</Th>
          <Th>Score</Th>
          {tasks.map((t) => (
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
                t.team_id === account?.team_id ? "rgba(255,255,255, 0.2)" : "",
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
            {tasks.map((task) => {
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
              return <Td key={`${t.team_id}-${task.id}`}>{flag}</Td>;
            })}
          </Tr>
        ))}
      </Tbody>
    </Table>
  );
};

export default Ranking;
