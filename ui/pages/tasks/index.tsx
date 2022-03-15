import { fetchCTF } from "lib/api/ctf";
import { AllPageProps } from "lib/pages";
import { isStaticMode } from "lib/static";
import { GetStaticProps } from "next";
import Loading from "../../components/loading";
import Tasks from "../../components/tasks";
import useAccount, { Account, fetchAccount } from "../../lib/api/account";
import useTasks, { fetchTasks, Task } from "../../lib/api/tasks";

type TasksProps = {
  tasks: Task[];
  account: Account | null;
} & AllPageProps;

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
