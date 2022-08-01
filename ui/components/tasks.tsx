import { Task } from "lib/api/tasks";
import { orderBy } from "lodash";
import { useLocalStorage } from "usehooks-ts";
import { TeamInterface } from "../lib/team";

type TasksProps = {
  tasks: Task[];
  team: TeamInterface | null;
};

const Tasks = ({ tasks, team }: TasksProps) => {
  return <></>;
};

export default Tasks;
