import zer0pts_svg from "public/zer0pts_logo_white.svg";
import ottersec_png from "public/ottersec.png";
import htb_svg from "public/htb.svg";
import tw_png from "public/tokyowesterns.png";
import google_png from "public/google.png";
import Image from "next/image";
import NextLink from "next/link";
import {
  Stack,
  Box,
  Center,
  Text,
  UnorderedList,
  ListItem,
  Code,
  OrderedList,
  Flex,
  Link,
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
        We offer a diverse range of enjoyable challenges across various
        difficulty levels and categories, all without the need for any guessing
        skills.
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
        <OrderedList>
          <ListItem>
            <Text>There are prizes for the top-performing teams.</Text>
            <UnorderedList>
              <ListItem>
                <Text fontWeight="bold">
                  &#129351; 1000 USD + 1 year HTB voucher (VIP+) &#215; 4
                </Text>
              </ListItem>
              <ListItem>
                <Text fontWeight="bold">
                  &#129352; 500 USD + 1 year HTB voucher (VIP) &#215; 4
                </Text>
              </ListItem>
              <ListItem>
                <Text fontWeight="bold">
                  &#129353; 250 USD + 6 months HTB voucher (VIP) &#215; 4
                </Text>
              </ListItem>
            </UnorderedList>
          </ListItem>
          <ListItem>
            <Text>
              The team that secures the first place will qualify for the SECCON
              CTF 2023 Finals (International division).
            </Text>
          </ListItem>
        </OrderedList>
        <Text>
          The top 3 teams must submit writeups of some challenges to{" "}
          <code>zer0ptsctf@gmail.com</code> within 24h after the CTF ends. We
          will specify which challenges need writeups after the CTF.
        </Text>
      </Stack>

      <Text fontSize="xl" mt={4}>
        [ Rules ]
      </Text>
      <Text pl={4}>
        <UnorderedList>
          <ListItem>There is no limit on your team size.</ListItem>
          <ListItem>
            Anyone can participate in this CTF: there are no restrictions based
            on age or nationality.
          </ListItem>
          <ListItem>
            Your rank on the scoreboard depends on two factors: 1) your total
            number of points (higher is better); 2) the timestamp of your last
            solved challenge (erlier is better).
          </ListItem>
          <ListItem>
            The survey challenge is special: it awards you some points, but it
            doesn't update your "last solved challenge" timestamp. You can't get
            ahead simply by solving the survey faster.
          </ListItem>
          <ListItem>
            Brute-forcing the flags is not allowed. If you submit 5 incorrect
            flags in quick succession, the flag submission form will be locked
            for 5 minutes.
          </ListItem>
          <ListItem>Each person can participate in only one team.</ListItem>
          <ListItem>
            Sharing solutions, hints, or flags with other teams during the
            competition is strictly forbidden.
          </ListItem>
          <ListItem>You are not allowed to attack the scoreserver.</ListItem>
          <ListItem>You are not allowed to attack other teams.</ListItem>
          <ListItem>
            Having multiple accounts is not allowed. If you are unable to log in
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
      <Stack pl={4}>
        <Text>Generous sponsors:</Text>
        <Flex justifyContent={"center"} alignItems={"center"} h="4em">
          <Box flex="1" p={5} pos="relative">
            <Link target="_blank" href="https://osec.io/">
              <Image
                unoptimized={true}
                src={ottersec_png}
                alt="OtterSec"
                layout="fill"
                objectFit="contain"
              />
            </Link>
          </Box>

          <Box flex="1" p={5} pos="relative">
            <Link href="https://www.hackthebox.com/" target="_blank">
              <Image
                unoptimized={true}
                src={htb_svg}
                alt="HackTheBox"
                layout="fill"
                objectFit="contain"
              />
            </Link>
          </Box>

          <Box flex="1" p={10} pos="relative">
            <Link href="https://goo.gle/ctfsponsorship" target="_blank">
              <Image
                unoptimized={true}
                src={google_png}
                alt="Google CTF Sponsorship"
                layout="fill"
                objectFit="contain"
              />
            </Link>
          </Box>

          <Box flex="1" p={10} pos="relative">
            <Link href="https://twitter.com/tokyowesterns" target="_blank">
              <Image
                unoptimized={true}
                src={tw_png}
                alt="TokyoWesterns"
                layout="fill"
                objectFit="contain"
              />
            </Link>
          </Box>
        </Flex>
        <Center>
          Infra sponsored by{" "}
          <Link
            as={NextLink}
            href="https://goo.gle/ctfsponsorship"
            about="_blank"
          >
            Google CTF Sponsorship
          </Link>
          .
        </Center>
      </Stack>
    </Stack>
  );
};
export default Index;
