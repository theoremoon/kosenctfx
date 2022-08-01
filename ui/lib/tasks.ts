import { Task } from "lib/api/tasks";
import { Account } from "lib/api/account";
import { orderBy } from "lodash";

export const filterTask = (tasks: Task[], filterText: string) =>
  tasks.filter((t) => {
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

export const sortTask = (tasks: Task[]) =>
  orderBy(
    tasks,
    [(t: Task) => t.score, (t: Task) => t.solved_by.length],
    ["asc", "desc"]
  );

export const isSolved = (account: Account | null) => (task: Task) =>
  task.solved_by.map((by) => by.team_id).includes(account?.team_id || 0);
