import {
  Box,
  Button,
  Flex,
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
import useMessage from "lib/useMessage";
import NextLink from "next/link";
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
  const { mutate } = useTasks();
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
    <Stack spacing={1} color={bgColor}>
      <Text fontSize="xl">
        {task.name} - {task.score}
      </Text>
      <Flex>
        <Stack w="70%" pl={1} spacing={1}>
          <Tags tags={[task.category, ...task.tags]} />
          <Box
            sx={{
              a: {
                textDecoration: "underline",
                textDecorationColor: pink,
                textDecorationThickness: "0.1em",
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
                <Tag variant="outline" maxW="10em">
                  <FontAwesomeIcon icon={faDownload} />
                  <TagLabel isTruncated>{a.name}</TagLabel>
                </Tag>
              </a>
            ))}
          </HStack>
          {task.author && <Right>author: {task.author}</Right>}
          <HStack>
            <Input
              id="flag"
              placeholder="Neko{...}"
              variant="flushed"
              {...register("flag", { required: true })}
            />
            <Button onClick={handleSubmit(onSubmit)}>Submit</Button>
          </HStack>
        </Stack>
        <Box w="30%" pl={1} sx={{ overflowY: "auto" }}>
          <Stack spacing={1} h="0">
            <Text size="lg">{task.solved_by.length}solves</Text>
            <Box>
              {task.solved_by.map((t) => (
                <Text fontSize="sm" key={t.team_name}>
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
