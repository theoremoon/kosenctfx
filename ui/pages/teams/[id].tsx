import { GetStaticPaths, GetStaticProps } from "next";
import Loading from "components/loading";
import { isStaticMode } from "lib/static";
import { AllPageProps } from "lib/pages";
import { fetchCTF } from "lib/api/ctf";
import { fetchSeries } from "lib/api/series";
import useScoreboard, {
  fetchScoreboard,
  ScoreFeedEntry,
} from "lib/api/scoreboard";
import useTeam, { fetchTeam } from "lib/api/team";
import TeamView from "theme/team";
import { TeamProps } from "props/team";

type teamProps = Omit<TeamProps & AllPageProps, "scorefeed"> & {
  scoreboard: ScoreFeedEntry[];
};
const TeamPage = ({
  team: defaultTeam,
  scoreboard: defaultScoreboard,
  series: series,
}: teamProps) => {
  const { data: team } = useTeam(defaultTeam.team_id.toString(), defaultTeam);
  const { data: scoreboard } = useScoreboard(defaultScoreboard);

  if (!team || !scoreboard) {
    return <Loading />;
  }
  const scorefeed = scoreboard.filter(
    (score) => score.team_id === team.team_id
  )[0];

  return TeamView({ team, scorefeed, series });
};

export const getStaticProps: GetStaticProps<teamProps> = async (context) => {
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
