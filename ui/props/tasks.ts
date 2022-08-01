import { Task } from "lib/api/tasks";

export interface TasksProps {
  tasks: Task[];
  filterTask: (tasks: Task[], filterText: string) => Task[];
  sortTask: (tasks: Task[]) => Task[];
  isSolved: (task: Task) => boolean;
}
