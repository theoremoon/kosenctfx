import {
  ChakraProps,
  Icon,
  Link,
  Box,
  Spacer,
  Stack,
  Text,
  Flex,
  Heading,
} from "@chakra-ui/react";
import { faCheck } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Task } from "lib/api/tasks";
import { useRouter } from "next/router";
import React from "react";
import { pink } from "../lib/color";
import Tags from "./tags";

type TaskCardProps = {
  task: Task;
  isSolved: boolean;
} & React.ComponentPropsWithoutRef<"div"> &
  ChakraProps;

const TaskCard = ({ task, isSolved, ...props }: TaskCardProps) => {
  const router = useRouter();
  return (
    <Link
      onClick={() => {
        router.push("/tasks/" + task.id, undefined, {
          scroll: false,
          shallow: true,
        });
      }}
      sx={{
        borderRadius: "4px",
        backgroundColor: "#edf2f7",
        p: 2,
        filter: isSolved ? "brightness(0.7)" : "none",
        "&:hover": {
          textDecoration: "none",
          filter: "brightness(0.7)",
          cursor: "pointer",
        },
      }}
    >
      <Box>
        <Heading fontSize="xl" color="#000">
          {isSolved && (
            <Icon
              as={FontAwesomeIcon}
              icon={faCheck}
              sx={{
                color: "#88c4d7",
                fontSize: "1em",
              }}
            />
          )}{" "}
          {task.name}
        </Heading>
        <Flex justify="space-around" m={0}>
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
        </Flex>
        <Tags tags={[task.category, ...task.tags]} />
      </Box>
    </Link>
  );
};

export default TaskCard;
