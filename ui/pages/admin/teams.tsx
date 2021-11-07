import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  Stack,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import Right from "components/right";
import useMessage from "lib/useMessage";
import { useState } from "react";
import { api } from "../../lib/api";
import useScoreboard from "../../lib/api/scoreboard";
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

const Teams = () => {
  const [teamName, setTeamName] = useState("");
  const [teamEmail, setTeamEmail] = useState("");
  const [team, setTeam] = useState<Team>();
  const { data: scoreboard } = useScoreboard();

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
      setTeamName(team.team_id);
      search(team.teamname);
    } catch (e) {
      error(e);
    }
  };

  return (
    <Stack>
      <Box w="md" mx="auto">
        <FormControl>
          <FormLabel htmlFor="teamname">teamname</FormLabel>
          <Input
            id="teamname"
            variant="flushed"
            value={teamName}
            onChange={(e) => setTeamName(e.target.value)}
          />
        </FormControl>
        <FormControl pt={4}>
          <Right>
            <Button onClick={() => search(teamName)}>Search</Button>
          </Right>
        </FormControl>
      </Box>

      {team && (
        <Box p={2}>
          <FormControl>
            <Input
              value={teamEmail}
              onChange={(e) => setTeamEmail(e.target.value)}
              placeholder={team.email}
            />
          </FormControl>
          <FormControl pt={4}>
            <Right>
              <Button
                onClick={() => {
                  updateEmail();
                }}
              >
                Update Team Email
              </Button>
            </Right>
          </FormControl>
          <Table>
            <Thead>
              <Tr>
                <Th>id</Th>
                <Th>name</Th>
                <Th>email</Th>
                <Th>country</Th>
              </Tr>
            </Thead>
            <Tbody>
              <Tr>
                <Td>{team.team_id}</Td>
                <Td>{team.teamname}</Td>
                <Td>{team.email}</Td>
                <Td>{team.country}</Td>
              </Tr>
            </Tbody>
          </Table>
        </Box>
      )}

      {team && team.submissions.length > 0 && (
        <Box pt={8}>
          <Table size="sm">
            <Thead>
              <Tr>
                <Th>flag</Th>
                <Th>is_correct</Th>
                <Th>is_valid</Th>
                <Th>ip address</Th>
                <Th>time</Th>
              </Tr>
            </Thead>
            <Tbody>
              {team.submissions.map((s) => (
                <Tr key={s.SubmittedAt}>
                  <Td maxW="200px" textOverflow="ellipsis" overflow="hidden">
                    <pre title={s.Flag}>{s.Flag}</pre>
                  </Td>
                  <Td>{s.IsCorrect ? "⭕" : "❌"}</Td>
                  <Td>{s.IsValid ? "⭕" : "❌"}</Td>
                  <Td>{s.IPAddress}</Td>
                  <Td>{dateFormat(s.SubmittedAt)}</Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </Box>
      )}

      <Box pt={8}>
        <Table size="sm">
          <Thead>
            <Tr>
              <Th>Pos</Th>
              <Th>Score</Th>
              <Th>Team</Th>
            </Tr>
          </Thead>
          <Tbody>
            {scoreboard &&
              scoreboard.map((t) => (
                <Tr key={t.team_id} onClick={() => search(t.team)}>
                  <Td>{t.pos}</Td>
                  <Td>{t.score}</Td>
                  <Td>{t.team}</Td>
                </Tr>
              ))}
          </Tbody>
        </Table>
      </Box>
    </Stack>
  );
};

export default Teams;
