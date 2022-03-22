import { Stack, Table, Tbody, Text, Th, Thead, Tr } from "@chakra-ui/react";
import SeriesChart from "components/chart";
import CountryFlag from "components/countryflag";
import Loading from "components/loading";
import { CTF, fetchCTF } from "lib/api/ctf";
import { fetchSeries, SeriesEntry } from "lib/api/series";
import { isStaticMode } from "lib/static";
import { orderBy } from "lodash";
import { GetStaticPaths, GetStaticProps } from "next";
import useScoreboard, {
  fetchScoreboard,
  ScoreFeedEntry,
} from "../../lib/api/scoreboard";
import useTeam, { fetchTeam, Team } from "../../lib/api/team";
import { dateFormat } from "../../lib/date";
import { AllPageProps } from "../../lib/pages";

type TeamProps = {
  teamID: number;
  team: Team;
  scoreboard: ScoreFeedEntry[];
  series: SeriesEntry[][];
} & AllPageProps;

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
  const series = await fetchSeries([team.teamname]);
  const ctf = await fetchCTF();
  return {
    props: {
      teamID: Number(id),
      team: team,
      scoreboard: scoreboard,
      series: series,
      ctf: ctf,
      account: null,
    },
    revalidate: isStaticMode ? undefined : 1,
  };
};

export const getStaticPaths: GetStaticPaths = async () => {
  const scoreboard = await fetchScoreboard().catch(() => []);

  return {
    paths: scoreboard.map((entry) => ({
      params: { id: entry.team_id.toString() },
    })),
    fallback: isStaticMode ? false : "blocking",
  };
};

export default TeamPage;
