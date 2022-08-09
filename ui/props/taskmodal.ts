import { FormEventHandler } from "react";
import { UseFormRegister } from "react-hook-form";

import { Task } from "lib/api/tasks";

export type FlagSubmitParams = {
  flag: string;
};

export interface TaskModalProps {
  task: Task;
  onClose: () => void;
  registerFlag: UseFormRegister<FlagSubmitParams>;
  onFlagSubmit: FormEventHandler<HTMLFormElement>;
  isSolved: boolean;
}
