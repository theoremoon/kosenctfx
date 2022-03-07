import {
  Box,
  Button,
  Flex,
  Heading,
  HStack,
  Input,
  Link,
  Stack,
  Tag,
  TagLabel,
  Text,
} from "@chakra-ui/react";
import { faDownload } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Task } from "lib/api/tasks";
import { isStaticMode } from "lib/static";
import useMessage from "lib/useMessage";
import NextLink from "next/link";
import React from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { api } from "../lib/api";
import useTasks from "../lib/api/tasks";
import { bgColor, pink } from "../lib/color";
import Right from "./right";
import Tags from "./tags";

type TaskModalBodyProps = {
  task: Task;
} & React.ComponentPropsWithoutRef<"div">;

type SubmitParams = {
  flag: string;
};

const TaskModalBody = ({ task, ...props }: TaskModalBodyProps) => {
  const { register, handleSubmit } = useForm<SubmitParams>();
  const { mutate } = useTasks([]);
  const { message, error } = useMessage();
  const onSubmit: SubmitHandler<SubmitParams> = async (values) => {
    try {
      const res = await api.post("/submit", {
        flag: values.flag,
      });
      message(res);
      mutate();
    } catch (e) {
      error(e);
    }
  };

  return (
    <Stack spacing={1} color="#000">
      <Heading as="h2" fontSize="3xl">
        {task.name}
      </Heading>
      <Flex>
        <Stack w="70%" pl={1} spacing={1}>
          <Tags
            tags={task.category ? [task.category, ...task.tags] : task.tags}
          />
          <HStack>
            <Box color="#000">
              <Box fontSize="2xl" sx={{ display: "inline" }}>
                {task.score}
              </Box>
              pts
            </Box>
            <Box color="#000">
              <Box fontSize="2xl" sx={{ display: "inline" }}>
                {task.solved_by.length}
              </Box>
              solves
            </Box>
          </HStack>
          <Box
            sx={{
              a: {
                color: "blue.300",
              },
              "a:focus": {
                outline: "none",
              },
            }}
            dangerouslySetInnerHTML={{ __html: task.description }}
          />
          <HStack minH="4em">
            {task.attachments.map((a) => (
              <a href={a.url} download key={a.url}>
                <Tag colorScheme="blackAlpha" variant="solid" maxW="10em">
                  <FontAwesomeIcon icon={faDownload} />
                  <TagLabel isTruncated>{a.name}</TagLabel>
                </Tag>
              </a>
            ))}
          </HStack>
          {task.author && <Right>author: {task.author}</Right>}
          {!isStaticMode && (
            <HStack>
              <Input
                id="flag"
                placeholder="zer0pts{...}"
                variant="flushed"
                sx={{
                  borderColor: "blue.300",
                  "&::placeholder": {
                    color: "#999",
                  },
                }}
                {...register("flag", { required: true })}
              />
              <Button
                onClick={handleSubmit(onSubmit)}
                colorScheme="blue"
                variant="solid"
              >
                Submit
              </Button>
            </HStack>
          )}
        </Stack>
        <Box w="30%" pl={1} sx={{ overflowY: "auto" }}>
          <Stack spacing={1} h="0">
            <Text fontSize="xl">solved by ({task.solved_by.length})</Text>
            <Box pl={2}>
              {task.solved_by.map((t) => (
                <Text
                  fontSize="sm"
                  key={t.team_name}
                  sx={{
                    "&:hover": {
                      textDecoration: "underline",
                    },
                  }}
                >
                  <Link as={NextLink} href={`/teams/${t.team_id}`}>
                    {t.team_name}
                  </Link>
                </Text>
              ))}
            </Box>
          </Stack>
        </Box>
      </Flex>
    </Stack>
  );
};

export default TaskModalBody;
