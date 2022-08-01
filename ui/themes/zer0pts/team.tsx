import { Stack, Table, Tbody, Text, Th, Thead, Tr } from "@chakra-ui/react";
import { orderBy } from "lodash";
import { TeamProps } from "props/team";
import { dateFormat } from "lib/date";
import CountryFlag from "./components/countryflag";
import ChartView from "./components/chart";

const Team = ({ team, scorefeed, series }: TeamProps) => {
  const tasks = orderBy(
    Object.entries(scorefeed.taskStats).map(([key, taskStat]) => ({
      name: key,
      ...taskStat,
    })),
    ["solved_at"],
    ["desc"]
  );
  return (
    <Stack mt="10">
      <Text fontSize="2xl">
        {team.country && (
          <CountryFlag country={team.country} sx={{ display: "inline" }} />
        )}
        {team.teamname}
      </Text>
      <Text>
        Rank {scorefeed.pos} / {scorefeed.score} points
      </Text>

      <ChartView chartTeams={[team.teamname]} chartSeries={series} />

      {scorefeed && (
        <Table>
          <Thead>
            <Tr>
              <Th>Task</Th>
              <Th>Score</Th>
              <Th>Solved At</Th>
            </Tr>
          </Thead>
          <Tbody>
            {tasks.map((task) => (
              <tr key={task.name}>
                <td>{task.name}</td>
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
