import {
  Box,
  Center,
  Code,
  Flex,
  Link,
  ListItem,
  Spacer,
  Stack,
  Text,
  UnorderedList,
} from "@chakra-ui/react";
import { fetchAccount } from "lib/api/account";
import { dateFormat } from "lib/date";
import { AllPageProps } from "lib/pages";
import { isStaticMode } from "lib/static";
import type { GetStaticProps, NextPage } from "next";
import useCTF, { fetchCTF } from "../lib/api/ctf";

type IndexPageProps = AllPageProps & {};
const Index: NextPage<IndexPageProps> = ({ ctf: fallbackCTF, ...props }: IndexPageProps) => {
  const { data: ctf } = useCTF(fallbackCTF);
  return (
    <Stack>
      <Text fontSize="4xl"> CTF Name</Text>

      {ctf && (
        <>
          <Text fontSize="xl" pt={10}>
            [ Time and Date ]
          </Text>

          <Text pl={4}>
            {dateFormat(ctf.start_at)} - {dateFormat(ctf.end_at)}
          </Text>
        </>
      )}

      <Text fontSize="xl" pt={10}>
        [ Contact ]
      </Text>

      <Text pl={4}>Discord: TBD</Text>

      <Text fontSize="xl" pt={4}>
        [ Prizes ]
      </Text>
      <Stack pl={4}>
        <Text fontWeight="bold">1st place</Text>
        <UnorderedList pl={4}>
          <ListItem>UOUO FISH LIFE</ListItem>
        </UnorderedList>
      </Stack>

      <Text fontSize="xl" pt={4}>
        [ Rules ]
      </Text>
      <Text pl={4}>
        <UnorderedList>
          <ListItem>Your team can be of any size.</ListItem>
          <ListItem>
            Anyone is allowed to participate: no restriction on age or
            nationality.
          </ListItem>
          <ListItem>
            Your position on the scoreboard depends on 2 things: 1) your total
            number of points (higher is better); 2) the timestamp of your last
            solved challenge (lower is better).
          </ListItem>
          <ListItem>
            The survey challenge is special: it does award you some points, but
            it doesn't update your "last solved challenge" timestamp. You can't
            get ahead simply by solving the survey faster.
          </ListItem>
          <ListItem>
            You can't brute-force flags. If you submit 5 incorrect flags in a
            short succession, the flag submission form will get locked for 5
            minutes.
          </ListItem>
          <ListItem>One person can participate in only one team.</ListItem>
          <ListItem>
            Sharing solutions, hints or flags with other teams during the
            competition is strictly forbidden.
          </ListItem>
          <ListItem>You are not allowed to attack the scoreserver.</ListItem>
          <ListItem>You are not allowed to attack other teams.</ListItem>
          <ListItem>
            You are not allowed to have multiple accounts. If you can't log in
            to your account, please contact us on Discord.
          </ListItem>
          <ListItem>
            We reserve the right to ban and disqualify any team that chooses to
            break any of these rules.
          </ListItem>
          <ListItem>
            The flag format is{" "}
            <Code variant="solid" colorScheme="gray">
              {"Neko\\{[\\x20-\\x7e]+\\}"}
            </Code>
            , unless specified otherwise.
          </ListItem>
          <ListItem>Most importantly: good luck and have fun!</ListItem>
        </UnorderedList>
      </Text>
    </Stack>
  );
};

export const getStaticProps: GetStaticProps<IndexPageProps> = async () => {
  const ctf = await fetchCTF();
  const account = (isStaticMode) ? null : await fetchAccount();
  return {
    props: {
      ctf: ctf,
      account: account,
    }
  }
};

export default Index;
