import { fetchCTF } from "lib/api/ctf";
import { AllPageProps } from "lib/pages";
import { isStaticMode } from "lib/static";
import { GetStaticProps } from "next";
import Loading from "../../components/loading";
import useAccount, { fetchAccount } from "../../lib/api/account";
import useTasks, { fetchTasks, Task } from "../../lib/api/tasks";
import TasksView from "theme/tasks";
import { useCallback } from "react";
import { filterTask, sortTask, isSolved } from "lib/tasks";

type tasksProps = {
  tasks: Task[];
} & AllPageProps;

const TasksDefault = ({
  tasks: defaultTasks,
  account: defaultAccount,
}: tasksProps) => {
  const { data: tasks } = useTasks(defaultTasks);
  const { data: account } = useAccount(defaultAccount);
  const isSolvedByTeam = useCallback(isSolved(account || null), [account]);

  if (!tasks || account === undefined) {
    return <Loading />;
  }

  return (
    <>
      <TasksView
        tasks={tasks}
        filterTask={filterTask}
        sortTask={sortTask}
        isSolved={isSolvedByTeam}
      />
    </>
  );
};

export const getStaticProps: GetStaticProps<tasksProps> = async () => {
  const account = isStaticMode ? null : await fetchAccount().catch(() => null);
  const tasks = await fetchTasks().catch(() => []);
  const ctf = await fetchCTF();
  return {
    props: {
      tasks: tasks,
      account: account,
      ctf: ctf,
    },
  };
};

export default TasksDefault;
