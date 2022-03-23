import { ChakraProps, Icon, Link, Spacer, Stack, Text } from "@chakra-ui/react";
import { faCheck } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Task } from "lib/api/tasks";
import { useRouter } from "next/router";
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
        borderWidth: "2px",
        borderColor: "white",
        borderRadius: "4px",
        p: 4,
        width: "200px",
        height: "calc(200px / 4 * 3)",
        filter: isSolved ? "brightness(0.7)" : "none",
        "&:hover": {
          textDecoration: "none",
          borderColor: pink,
          cursor: "pointer",
        },
      }}
    >
      <Stack>
        <Text fontSize="xl">
          {isSolved && (
            <Icon
              as={FontAwesomeIcon}
              icon={faCheck}
              sx={{
                color: "#00ff00",
                fontSize: "1em",
              }}
            />
          )}{" "}
          {task.name}
        </Text>
        <Spacer />
        <Text>
          {task.score} points / {task.solved_by.length} solves
        </Text>
        <Tags tags={task.tags} />
      </Stack>
    </Link>
  );
};

export default TaskCard;
