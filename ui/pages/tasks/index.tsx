import { fetchCTF } from "lib/api/ctf";
import { AllPageProps } from "lib/pages";
import { GetStaticProps } from "next";
import Loading from "../../components/loading";
import useAccount from "../../lib/api/account";
import useTasks, { fetchTasks, Task } from "../../lib/api/tasks";
import { useCallback } from "react";
import { filterTask, sortTask, isSolved } from "lib/tasks";
import TasksView from "theme/tasks";
import { isStaticMode, revalidateInterval } from "lib/static";

type tasksProps = {
  tasks: Task[];
} & AllPageProps;

const TasksDefault = ({ tasks: defaultTasks }: tasksProps) => {
  const { data: tasks } = useTasks(defaultTasks);
  const { data: account } = useAccount(null);
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
  const tasks = await fetchTasks().catch(() => []);
  const ctf = await fetchCTF();
  return {
    props: {
      tasks: tasks,
      ctf: ctf,
    },
    revalidate: isStaticMode ? false : 10, // override revalidateInterval
  };
};

export default TasksDefault;
