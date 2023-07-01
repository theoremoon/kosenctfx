import zer0pts_svg from "public/zer0pts_logo_white.svg";
import Image from "next/image";
import {
  Stack,
  Box,
  Center,
  Text,
  UnorderedList,
  ListItem,
  Code,
} from "@chakra-ui/react";
import { dateFormat } from "lib/date";
import { IndexProps } from "props/index";

const Index = ({ ctf, status }: IndexProps) => {
  return (
    <Stack mt={5}>
      <Box maxW="container.sm" mx="auto">
        <Box maxW="3xs" mx="auto">
          <Image unoptimized={true} src={zer0pts_svg} />
        </Box>
        <Center fontSize="4xl">zer0pts CTF 2023</Center>

        {ctf && (
          <>
            <Center>
              {dateFormat(ctf.start_at)} - {dateFormat(ctf.end_at)}
            </Center>
            <Center>
              <Text>{status}</Text>
            </Center>
          </>
        )}
      </Box>

      <Text fontSize="xl" mt={10}>
        [ About ]
      </Text>
      <Text pl={4}>
        Welcome to zer0pts CTF 2023! <br />
        zer0pts CTF is a jeopardy-style CTF.
        <br />
        We provide many fun challenges of varying difficulty and categories, and
        none of them require any guessing skills.
      </Text>

      <Text fontSize="xl" mt={10}>
        [ Contact ]
      </Text>
      <Text pl={4}>
        Discord:{" "}
        <a href="https://discord.gg/3QrDP2sMYd">
          https://discord.gg/3QrDP2sMYd
        </a>
      </Text>

      <Text fontSize="xl" mt={10}>
        [ Prizes ]
      </Text>
      <Stack pl={4}>
        <UnorderedList>
          <ListItem>
            <Text fontWeight="bold">1st: 800 USD</Text>
          </ListItem>
          <ListItem>
            <Text fontWeight="bold">2nd: 500 USD</Text>
          </ListItem>
          <ListItem>
            <Text fontWeight="bold">3rd: 300 USD</Text>
          </ListItem>
          <ListItem>
            <Text fontWeight="bold">4th: 200 USD</Text>
          </ListItem>
          <ListItem>
            <Text fontWeight="bold">5th: 200 USD</Text>
          </ListItem>
        </UnorderedList>
      </Stack>

      <Text fontSize="xl" mt={4}>
        [ Rules ]
      </Text>
      <Text pl={4}>
        <UnorderedList>
          <ListItem>No limit on your team size.</ListItem>
          <ListItem>
            Anyone can participate in this CTF: no restriction on your age or
            nationality.
          </ListItem>
          <ListItem>
            Your rank on the scoreboard depends on: 1) your total number of
            points (higher is better); 2) the timestamp of your last solved
            challenge (erlier is better).
          </ListItem>
          <ListItem>
            The survey challenge is special: it does award you some points, but
            it doesn't update your "last solved challenge" timestamp. You can't
            get ahead simply by solving the survey faster.
          </ListItem>
          <ListItem>
            You can't brute-force the flags. If you submit 5 incorrect flags in
            a short succession, the flag submission form will get locked for 5
            minutes.
          </ListItem>
          <ListItem>One person can participate in only one team.</ListItem>
          <ListItem>
            Sharing solutions, hints, or flags with other teams during the
            competition is strictly forbidden.
          </ListItem>
          <ListItem>You are not allowed to attack the scoreserver.</ListItem>
          <ListItem>You are not allowed to attack other teams.</ListItem>
          <ListItem>
            You are not allowed to have multiple accounts. If you can't log in
            to your account, please contact us on Discord.
          </ListItem>
          <ListItem>
            We reserve the right to ban and disqualify any teams breaking any of
            these rules.
          </ListItem>
          <ListItem>
            The flag format is{" "}
            <Code variant="solid" colorScheme="gray">
              {"zer0pts\\{[\\x20-\\x7e]+\\}"}
            </Code>
            , unless specified otherwise.
          </ListItem>
          <ListItem>Most importantly: good luck and have fun!</ListItem>
        </UnorderedList>
      </Text>

      <Text fontSize="xl" mt={10}>
        [ Sponsors ]
      </Text>
    </Stack>
  );
};
export default Index;
