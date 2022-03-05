import { isStaticMode } from "lib/static";
import { GetStaticProps } from "next";
import Loading from "../../components/loading";
import Tasks from "../../components/tasks";
import useAccount, { Account, fetchAccount } from "../../lib/api/account";
import useTasks, { fetchTasks, Task } from "../../lib/api/tasks";

interface TasksProps {
  tasks: Task[];
  account: Account | null;
}

const TasksDefault = ({
  tasks: defaultTasks,
  account: defaultAccount,
}: TasksProps) => {
  const { data: tasks } = useTasks(defaultTasks);
  const { data: account } = useAccount(defaultAccount);
  if (!tasks || account === undefined) {
    return <Loading />;
  }
  return <Tasks tasks={tasks} team={account} />;
};

export const getStaticProps: GetStaticProps<TasksProps> = async () => {
  const account = isStaticMode ? null : await fetchAccount();
  const tasks = await fetchTasks();
  return {
    props: {
      tasks: tasks,
      account: account,
    },
  };
};

export default TasksDefault;
