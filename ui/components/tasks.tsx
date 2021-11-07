import { Box, Stack, Text, Wrap, WrapItem } from "@chakra-ui/react";
import TaskCard from "components/taskcard";
import { Task } from "lib/api/tasks";
import { groupBy, orderBy, sortBy, sortedUniq } from "lodash";
import is_solved from "../lib/is_solved";
import { TeamInterface } from "../lib/team";

type TasksProps = {
  tasks: Task[];
  team: TeamInterface | null;
};

const Tasks = ({ tasks, team }: TasksProps) => {
  const categories = sortedUniq(sortBy(tasks.map((t) => t.category)));
  const taskByCategories = groupBy(
    orderBy(
      tasks,
      [(t: Task) => t.score, (t: Task) => t.solved_by.length],
      ["asc", "desc"]
    ),
    (t) => t.category
  );
  return (
    <Box mt="10">
      {categories.map((c) => (
        <Stack key={c} mb={4}>
          <Text fontSize="3xl">{c}</Text>
          <Wrap pl={8}>
            {taskByCategories[c].map((t) => (
              <WrapItem key={t.id}>
                <TaskCard
                  task={t}
                  isSolved={team !== null && is_solved(team, t)}
                />
              </WrapItem>
            ))}
          </Wrap>
        </Stack>
      ))}
    </Box>
  );
};

export default Tasks;
