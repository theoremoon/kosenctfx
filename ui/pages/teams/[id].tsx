import { Stack, Table, Tbody, Text, Th, Thead, Tr } from "@chakra-ui/react";
import SeriesChart from "components/chart";
import CountryFlag from "components/countryflag";
import Loading from "components/loading";
import { fetchSeries, SeriesEntry } from "lib/api/series";
import { orderBy } from "lodash";
import { GetStaticPaths, GetStaticProps } from "next";
import useScoreboard, {
  fetchScoreboard,
  ScoreFeedEntry,
} from "../../lib/api/scoreboard";
import useTeam, { fetchTeam, Team } from "../../lib/api/team";
import { dateFormat } from "../../lib/date";

interface TeamProps {
  teamID: number;
  team: Team;
  scoreboard: ScoreFeedEntry[];
  series: SeriesEntry[][];
}

const TeamPage = ({
  teamID,
  team: defaultTeam,
  scoreboard: defaultScoreboard,
  series: defaultSeries,
}: TeamProps) => {
  const { data: team } = useTeam(teamID.toString(), defaultTeam);
  const { data: scoreboard } = useScoreboard(defaultScoreboard);

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

      <SeriesChart teams={[team.teamname]} series={defaultSeries} />

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

export const getStaticProps: GetStaticProps<TeamProps> = async (context) => {
  const id = context.params?.id;
  if (id === undefined) {
    return { notFound: true };
  }

  const scoreboard = await fetchScoreboard();
  const team = await fetchTeam(id.toString());
  const series = await fetchSeries([]);
  return {
    props: {
      teamID: Number(id),
      team: team,
      scoreboard: scoreboard,
      series: series,
    },
  };
};

export const getStaticPaths: GetStaticPaths = async () => {
  const scoreboard = await fetchScoreboard();

  return {
    paths: scoreboard.map((entry) => ({
      params: { id: entry.team_id.toString() },
    })),
    fallback: false,
  };
};

export default TeamPage;
