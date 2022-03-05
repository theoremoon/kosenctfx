import {
  Box,
  FormControl,
  FormLabel,
  Input,
  SimpleGrid,
  Stack,
  Switch,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";
import TaskCard from "components/taskcard";
import { Task } from "lib/api/tasks";
import { orderBy } from "lodash";
import React, { useState } from "react";
import { useLocalStorage } from "usehooks-ts";
import { is_solved } from "../lib/is_solved";
import { TeamInterface } from "../lib/team";

type TasksProps = {
  tasks: Task[];
  team: TeamInterface | null;
};

const Tasks = ({ tasks, team }: TasksProps) => {
  const [showSolved, setShowSolved] = useLocalStorage("showSolved", true);
  const [filterText, setFilterText] = useLocalStorage("filterText", "");

  const filteredTasks = tasks.filter((t) => {
    if (filterText === "") {
      return true;
    }
    if (t.name.includes(filterText)) {
      return true;
    }
    if (t.category.includes(filterText)) {
      return true;
    }
    if (t.author.includes(filterText)) {
      return true;
    }
    if (t.tags.some((tag) => tag.includes(filterText))) {
      return true;
    }
    return false;
  });

  const sortedTasks = orderBy(
    filteredTasks,
    [(t: Task) => t.score, (t: Task) => t.solved_by.length],
    ["asc", "desc"]
  );
  return (
    <Box mt={4}>
      <Box>
        <FormControl display="flex">
          <FormLabel>Show solved tasks</FormLabel>
          <Switch
            isChecked={showSolved}
            onChange={() => {
              setShowSolved((oldValue) => !oldValue);
            }}
          />
        </FormControl>
        <Input
          placeholder="filter"
          value={filterText}
          onChange={(e) => setFilterText(e.target.value)}
        />
      </Box>
      <SimpleGrid columns={{ sm: 1, md: 3, lg: 5 }} spacing={4} mt={4}>
        {sortedTasks.map((t) => {
          const isSolved = team !== null && is_solved(team, t);
          if (isSolved && !showSolved) {
            return <></>;
          }
          return <TaskCard key={t.id} task={t} isSolved={isSolved} />;
        })}
      </SimpleGrid>
    </Box>
  );
};

export default Tasks;
