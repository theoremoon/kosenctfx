import { FormEventHandler } from "react";
import { UseFormRegister } from "react-hook-form";

import { Task } from "lib/api/tasks";

export type FlagSubmitParams = {
  flag: string;
};

export interface TaskProps {
  task: Task;
  tasks: Task[];
  tasksPath: string;

  registerFlag: UseFormRegister<FlagSubmitParams>;
  onFlagSubmit: FormEventHandler<HTMLFormElement>;

  filterTask: (tasks: Task[], filterText: string) => Task[];
  sortTask: (tasks: Task[]) => Task[];
  isSolved: (task: Task) => boolean;
}
