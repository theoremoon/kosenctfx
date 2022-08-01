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
  Modal,
  ModalBody,
  ModalContent,
  ModalOverlay,
  useDisclosure,
} from "@chakra-ui/react";
import { faDownload } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import NextLink from "next/link";

import TasksView from "./components/tasks";
import { TaskProps } from "props/task";
import { useRouter } from "next/router";
import Tags from "./components/tags";
import Right from "./components/right";

type taskModalBodyProps = Pick<
  TaskProps,
  "task" | "registerFlag" | "onFlagSubmit"
>;
const TaskModalBody = ({
  task,
  registerFlag,
  onFlagSubmit,
}: taskModalBodyProps) => {
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
          <form onSubmit={onFlagSubmit}>
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
                {...registerFlag("flag", { required: true })}
              />
              <Button type="submit" colorScheme="blue" variant="solid">
                Submit
              </Button>
            </HStack>
          </form>
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

const Task = ({
  task,
  tasksPath,
  tasks,
  filterTask,
  sortTask,
  isSolved,

  registerFlag,
  onFlagSubmit,
}: TaskProps) => {
  const { onClose } = useDisclosure();
  const router = useRouter();
  return (
    <>
      <TasksView
        tasks={tasks}
        filterTask={filterTask}
        sortTask={sortTask}
        isSolved={isSolved}
      />
      <Modal
        isOpen={true}
        onClose={() => {
          onClose();
          router.push(tasksPath, undefined, {
            scroll: false,
            shallow: true,
          });
        }}
        size="4xl"
      >
        <ModalOverlay />
        <ModalContent
          sx={{
            backgroundColor: "#ffffff",
          }}
        >
          <ModalBody>
            <TaskModalBody
              task={task}
              registerFlag={registerFlag}
              onFlagSubmit={onFlagSubmit}
            />
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
};
export default Task;
