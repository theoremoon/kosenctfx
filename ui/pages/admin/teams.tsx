import AdminLayout from "components/adminLayout";
import useConfig from "lib/api/admin/config";
import useMessage from "lib/useMessage";
import { GetStaticProps } from "next";
import React, { useState } from "react";
import { api } from "../../lib/api";
import useScoreboard, {
  fetchScoreboard,
  ScoreFeedEntry,
} from "../../lib/api/scoreboard";
import { dateFormat } from "../../lib/date";

type Submission = {
  ChallengeID: number | undefined;
  TeamID: number;
  IsCorrect: boolean;
  IsValid: boolean;
  Flag: string;
  IPAddress: string;
  SubmittedAt: number;
};

type Team = {
  teamname: string;
  team_id: string;
  email: string;
  country: string;
  submissions: Submission[];
};

interface AdminTeamsProps {
  scoreboard: ScoreFeedEntry[];
}

const Teams = ({ scoreboard: defaultScoreboard }: AdminTeamsProps) => {
  const [teamName, setTeamName] = useState("");
  const [teamEmail, setTeamEmail] = useState("");
  const [team, setTeam] = useState<Team>();
  const { data: scoreboard } = useScoreboard(defaultScoreboard);

  const { message, error } = useMessage();
  const search = async (teamname: string) => {
    try {
      const data = await api
        .get<Team>("/admin/team", {
          params: {
            team: teamname,
          },
        })
        .then((r) => r.data);
      setTeam(data);
    } catch (e) {
      error(e);
    }
  };

  const updateEmail = async () => {
    if (!team) {
      return;
    }
    try {
      const data = await api.post<Response>("/admin/update-email", {
        id: team.team_id,
        email: teamEmail,
      });
      message(data);
      setTeamName(team.teamname);
      search(team.teamname);
    } catch (e) {
      error(e);
    }
  };

  // admin-check
  const { data: config, error: configError } = useConfig();
  if (config === undefined || configError !== undefined) {
    return <></>;
  }

  return (
    <>
      <h5 className="mt-4">Search Team by Name</h5>
      <div className="input-group">
        <input
          className="form-control"
          placeholder="teamname"
          value={teamName}
          onChange={(e) => setTeamName(e.target.value)}
        />

        <button className="btn btn-primary" onClick={() => search(teamName)}>
          Search
        </button>
      </div>

      {team && (
        <>
          <h5 className="mt-4">Team Profile</h5>
          <table className="table">
            <thead>
              <tr>
                <th>id</th>
                <th>name</th>
                <th>email</th>
                <th>country</th>
              </tr>
            </thead>
            <tbody style={{ maxHeight: "400px" }}>
              <tr>
                <td>{team.team_id}</td>
                <td>{team.teamname}</td>
                <td>{team.email}</td>
                <td>{team.country}</td>
              </tr>
            </tbody>
          </table>

          <div className="input-group">
            <input
              className="form-control"
              value={teamEmail}
              onChange={(e) => setTeamEmail(e.target.value)}
              placeholder={team.email}
            />

            <button className="btn btn-primary" onClick={() => updateEmail()}>
              Update
            </button>
          </div>
        </>
      )}

      {team && team.submissions.length > 0 && (
        <>
          <h5 className="mt-4">Submitted Flags</h5>
          <table className="table table-sm">
            <thead>
              <tr>
                <th>flag</th>
                <th>is_correct</th>
                <th>is_valid</th>
                <th>ip address</th>
                <th>time</th>
              </tr>
            </thead>
            <tbody>
              {team.submissions.map((s) => (
                <tr key={s.SubmittedAt}>
                  <td>
                    <code>
                      <pre
                        title={s.Flag}
                        className="text-truncate"
                        style={{ maxWidth: "200px" }}
                      >
                        {s.Flag}
                      </pre>
                    </code>
                  </td>
                  <td>{s.IsCorrect ? "⭕" : "❌"}</td>
                  <td>{s.IsValid ? "⭕" : "❌"}</td>
                  <td>{s.IPAddress}</td>
                  <td>{dateFormat(s.SubmittedAt)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}

      <table className="table">
        <thead>
          <tr>
            <th>Pos</th>
            <th>Score</th>
            <th>Team</th>
          </tr>
        </thead>
        <tbody>
          {scoreboard &&
            scoreboard.map((t) => (
              <tr key={t.team_id} onClick={() => search(t.team)}>
                <td>{t.pos}</td>
                <td>{t.score}</td>
                <td>{t.team}</td>
              </tr>
            ))}
        </tbody>
      </table>
    </>
  );
};

export const getStaticProps: GetStaticProps<AdminTeamsProps> = async () => {
  const scoreboard = await fetchScoreboard().catch(() => []);
  return {
    props: {
      scoreboard: scoreboard,
    },
  };
};

Teams.getLayout = AdminLayout;

export default Teams;
