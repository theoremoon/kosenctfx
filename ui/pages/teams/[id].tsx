import { Stack, Table, Tbody, Text, Th, Thead, Tr } from "@chakra-ui/react";
import SeriesChart from "components/chart";
import CountryFlag from "components/countryflag";
import Loading from "components/loading";
import { orderBy } from "lodash";
import { useRouter } from "next/router";
import useScoreboard from "../../lib/api/scoreboard";
import useTeam from "../../lib/api/team";
import { dateFormat } from "../../lib/date";
const Team = () => {
  const router = useRouter();
  const id = router.query.id as string;
  const { data: team } = useTeam(id);
  const { data: scoreboard } = useScoreboard();

  if (!team || !scoreboard) {
    return <Loading />;
  }
  const teamScore = scoreboard.filter(
    (score) => score.team_id === team.team_id
  )[0];

  return (
    <Stack mt="10">
      <Text fontSize="2xl">
        {team.country && (
          <>
            <CountryFlag country={team.country} sx={{ display: "inline" }} />{" "}
          </>
        )}
        {team.teamname}
      </Text>
      <Text>
        Rank {teamScore.pos} / {teamScore.score} points
      </Text>

      <SeriesChart teams={[team.teamname]} />

      {teamScore && (
        <Table>
          <Thead>
            <Tr>
              <Th>Task</Th>
              <Th>Score</Th>
              <Th>Solved At</Th>
            </Tr>
          </Thead>
          <Tbody>
            {orderBy(
              Object.entries(teamScore.taskStats).map(([key, taskStat]) => ({
                taskname: key,
                ...taskStat,
              })),
              ["solved_at"],
              ["desc"]
            ).map((task) => (
              <tr key={task.taskname}>
                <td>{task.taskname}</td>
                <td>{task.points}</td>
                <td>{dateFormat(task.time)}</td>
              </tr>
            ))}
          </Tbody>
        </Table>
      )}
    </Stack>
  );
};

export default Team;
