import { Task } from "./api/tasks";
import { TeamInterface } from "./team";

export const is_solved = (
  team: TeamInterface | undefined,
  task: Task | undefined
): boolean => {
  if (!team || !task) {
    return false;
  }

  return (
    task.solved_by.filter((t) => t.team_id === Number(team.team_id)).length > 0
  );
};
