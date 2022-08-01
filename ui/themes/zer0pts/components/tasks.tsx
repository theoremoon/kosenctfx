import {
  Box,
  FormControl,
  FormLabel,
  Input,
  SimpleGrid,
  Switch,
} from "@chakra-ui/react";
import { useLocalStorage } from "usehooks-ts";
import TaskCard from "./taskCard";
import { TasksProps } from "props/tasks";

const Tasks = ({ tasks, filterTask, sortTask, isSolved }: TasksProps) => {
  const [showSolved, setShowSolved] = useLocalStorage("showSolved", true);
  const [filterText, setFilterText] = useLocalStorage("filterText", "");

  return (
    <Box mt={4}>
      <Box>
        <FormControl display="flex">
          <FormLabel>Show solved tasks</FormLabel>
          <Switch
            isChecked={showSolved}
            onChange={() => {
              setShowSolved((oldValue: boolean) => !oldValue);
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
        {sortTask(filterTask(tasks, filterText)).map((t) => {
          const solved = isSolved(t);
          if (solved && !showSolved) {
            return <></>;
          }
          return <TaskCard key={t.id} task={t} isSolved={solved} />;
        })}
      </SimpleGrid>
    </Box>
  );
};

export default Tasks;
